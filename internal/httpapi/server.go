package httpapi

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/industrial-ai/iaos-aese/internal/aese3"
	"github.com/industrial-ai/iaos-aese/internal/application"
	"github.com/industrial-ai/iaos-aese/internal/assurance"
	"github.com/industrial-ai/iaos-aese/internal/capabilitybuild"
	"github.com/industrial-ai/iaos-aese/internal/experiment"
	"github.com/industrial-ai/iaos-aese/internal/firstdelivery"
	"github.com/industrial-ai/iaos-aese/internal/genesis"
	"github.com/industrial-ai/iaos-aese/internal/iaosclient"
	"github.com/industrial-ai/iaos-aese/internal/incorporation"
	"github.com/industrial-ai/iaos-aese/internal/industrialization"
	"github.com/industrial-ai/iaos-aese/internal/legacyprojection"
	"github.com/industrial-ai/iaos-aese/internal/plantbuild"
	"github.com/industrial-ai/iaos-aese/internal/replay"
	"github.com/industrial-ai/iaos-aese/internal/scenariopack"
	"github.com/industrial-ai/iaos-aese/internal/strategyrelease"
)

const (
	defaultBodyLimit   = int64(1 << 20)
	defaultTimeout     = 30 * time.Second
	resetTokenTTL      = 10 * time.Minute
	resetTokenByteSize = 16
)

type Config struct {
	PackDir        string
	RequestTimeout time.Duration
	BodyLimit      int64
	Logger         *log.Logger
}

type Server struct {
	cfg   Config
	mux   *http.ServeMux
	logf  func(format string, args ...any)
	mu    sync.RWMutex
	runs  map[string]*runRecord
	order []string
}

type actionCache struct {
	HttpStatus         int       `json:"status"`
	Action             string    `json:"action"`
	Idempotency        string    `json:"idempotency_key"`
	ErrorCode          string    `json:"error_code,omitempty"`
	Error              string    `json:"error,omitempty"`
	ErrorRetryable     bool      `json:"retryable,omitempty"`
	RequiredPermission string    `json:"required_permission,omitempty"`
	Outcome            any       `json:"outcome,omitempty"`
	RunStatus          string    `json:"run_status"`
	Cursor             int64     `json:"cursor"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type runRecord struct {
	RunID               string
	PackKey             string
	PackVersion         string
	PackDir             string
	ScenarioKey         string
	Plan                application.Plan
	Status              application.RunStatus
	CurrentAct          int
	TenantID            string
	Target              string
	Token               string
	Actor               string
	CreatedAt           time.Time
	UpdatedAt           time.Time
	LastError           string
	Retryable           bool
	Cursor              int64
	ResetToken          string
	ResetTokenExpiresAt time.Time
	ActionCache         map[string]actionCache
}

type runCreateRequest struct {
	Target   string `json:"target"`
	Tenant   string `json:"tenant"`
	StoryKey string `json:"story_key"`
	PlanHash string `json:"plan_hash"`
	RunID    string `json:"run_id"`
	Actor    string `json:"actor"`
	Token    string `json:"token"`
	PackDir  string `json:"pack_dir"`
}

type runActionRequest struct {
	PlanHash          string `json:"plan_hash,omitempty"`
	RunVersion        string `json:"run_version,omitempty"`
	Apply             *bool  `json:"apply,omitempty"`
	DryRun            bool   `json:"dry_run,omitempty"`
	IdempotencyKey    string `json:"idempotency_key,omitempty"`
	ExpectedCursor    *int64 `json:"expected_cursor,omitempty"`
	ConfirmationToken string `json:"confirmation_token,omitempty"`
}

type errorResponse struct {
	Error              string `json:"error"`
	Code               string `json:"code"`
	Retryable          bool   `json:"retryable,omitempty"`
	RunID              string `json:"run_id,omitempty"`
	RunVersion         string `json:"run_version,omitempty"`
	Status             string `json:"status,omitempty"`
	RequiredPermission string `json:"required_permission,omitempty"`
}

type runResponse struct {
	RunID                     string   `json:"run_id"`
	RunVersion                string   `json:"run_version"`
	PackKey                   string   `json:"pack_key"`
	PackVersion               string   `json:"pack_version"`
	ScenarioKey               string   `json:"scenario_key"`
	PlanHash                  string   `json:"plan_hash"`
	Status                    string   `json:"status"`
	CurrentAct                int      `json:"current_act"`
	TotalActs                 int      `json:"total_acts"`
	Cursor                    int64    `json:"cursor"`
	TenantID                  string   `json:"tenant"`
	Target                    string   `json:"target"`
	CreatedAt                 string   `json:"created_at"`
	UpdatedAt                 string   `json:"updated_at"`
	AllowedActions            []string `json:"allowed_actions"`
	LastError                 string   `json:"last_error,omitempty"`
	Retryable                 bool     `json:"retryable,omitempty"`
	Outcome                   any      `json:"outcome,omitempty"`
	Plan                      any      `json:"plan,omitempty"`
	ResetConfirmationRequired bool     `json:"reset_confirmation_required,omitempty"`
}

type actionResponse struct {
	Run    runResponse `json:"run"`
	Action string      `json:"action"`
}

type apiError struct {
	statusCode         int
	code               string
	message            string
	retryable          bool
	requiredPermission string
}

func (e apiError) Error() string { return e.message }

func New(cfg Config) *Server {
	if cfg.PackDir == "" {
		cfg.PackDir = "scenario-packs/hctm"
	}
	if cfg.RequestTimeout <= 0 {
		cfg.RequestTimeout = defaultTimeout
	}
	if cfg.BodyLimit <= 0 {
		cfg.BodyLimit = defaultBodyLimit
	}
	if cfg.Logger == nil {
		cfg.Logger = log.Default()
	}
	server := &Server{
		cfg:  cfg,
		runs: map[string]*runRecord{},
		logf: cfg.Logger.Printf,
		mux:  http.NewServeMux(),
	}
	server.RegisterRoutes()
	return server
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.addCORSHeaders(w, r)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	s.mux.ServeHTTP(w, r)
}

func (s *Server) addCORSHeaders(w http.ResponseWriter, r *http.Request) {
	origin := strings.TrimSpace(r.Header.Get("Origin"))
	if origin == "" {
		origin = "*"
	}
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Vary", "Origin")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Idempotency-Key, X-Aese-Reset-Token")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Type, X-Request-ID")
	w.Header().Set("Access-Control-Allow-Credentials", "false")
}

func (s *Server) RegisterRoutes() {
	s.mux.HandleFunc("/health", s.handleHealth)
	s.mux.HandleFunc("/ready", s.handleReady)
	s.mux.HandleFunc("/api/aese/v1/", s.handleAPI)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed", false, "", "")
		return
	}
	s.writeJSON(w, http.StatusOK, map[string]any{"status": "UP"})
}

func (s *Server) handleReady(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed", false, "", "")
		return
	}
	s.writeJSON(w, http.StatusOK, map[string]any{
		"status": "OK",
		"routes": []string{
			"/health",
			"/ready",
			"/api/aese/v1/scenarios",
			"/api/aese/v1/runs/plan",
			"/api/aese/v1/runs",
			"/api/aese/v1/runs/:run_id",
			"/api/aese/v1/runs/:run_id/preflight",
			"/api/aese/v1/runs/:run_id/initialize",
			"/api/aese/v1/runs/:run_id/advance",
			"/api/aese/v1/runs/:run_id/run-to-end",
			"/api/aese/v1/runs/:run_id/analyze",
			"/api/aese/v1/runs/:run_id/verify",
			"/api/aese/v1/runs/:run_id/reset-plan",
			"/api/aese/v1/runs/:run_id/reset",
			"/api/aese/v1/world/genesis",
			"/api/aese/v1/world/incorporation",
			"/api/aese/v1/world/plant-build",
			"/api/aese/v1/world/capability-build",
			"/api/aese/v1/world/industrialization",
			"/api/aese/v1/world/first-delivery",
			"/api/aese/v1/world/experiments",
			"/api/aese/v1/world/strategy-control",
			"/api/aese/v1/world/strategy-assurance",
			"/api/aese/v1/world/aese3",
		},
	})
}

func (s *Server) handleAPI(w http.ResponseWriter, r *http.Request) {
	parts := splitPath(r.URL.Path)
	if len(parts) < 3 || parts[0] != "api" || parts[1] != "aese" || parts[2] != "v1" {
		s.writeError(w, http.StatusNotFound, "not_found", "route not found", false, "", "")
		return
	}
	rest := parts[3:]
	if len(rest) == 0 {
		s.writeError(w, http.StatusNotFound, "not_found", "route not found", false, "", "")
		return
	}

	ctx := context.Background()
	if s.cfg.RequestTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, s.cfg.RequestTimeout)
		defer cancel()
	}

	switch rest[0] {
	case "world":
		if len(rest) != 2 || r.Method != http.MethodGet {
			s.writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed", false, "", "")
			return
		}
		if rest[1] == "incorporation" {
			trace := incorporation.BuildTrace()
			if err := incorporation.Validate(trace); err != nil {
				s.writeError(w, http.StatusInternalServerError, "incorporation_invalid", err.Error(), false, "", "")
				return
			}
			s.writeJSON(w, http.StatusOK, trace)
			return
		}
		if rest[1] == "plant-build" {
			trace := plantbuild.BuildTrace()
			if err := plantbuild.Validate(trace); err != nil {
				s.writeError(w, http.StatusInternalServerError, "plant_build_invalid", err.Error(), false, "", "")
				return
			}
			s.writeJSON(w, http.StatusOK, trace)
			return
		}
		if rest[1] == "capability-build" {
			trace := capabilitybuild.BuildTrace()
			if err := capabilitybuild.Validate(trace); err != nil {
				s.writeError(w, http.StatusInternalServerError, "capability_build_invalid", err.Error(), false, "", "")
				return
			}
			s.writeJSON(w, http.StatusOK, trace)
			return
		}
		if rest[1] == "industrialization" {
			trace := industrialization.BuildTrace()
			if err := industrialization.Validate(trace); err != nil {
				s.writeError(w, http.StatusInternalServerError, "industrialization_invalid", err.Error(), false, "", "")
				return
			}
			s.writeJSON(w, http.StatusOK, trace)
			return
		}
		if rest[1] == "first-delivery" {
			trace := firstdelivery.BuildTrace()
			if err := firstdelivery.Validate(trace); err != nil {
				s.writeError(w, http.StatusInternalServerError, "first_delivery_invalid", err.Error(), false, "", "")
				return
			}
			s.writeJSON(w, http.StatusOK, trace)
			return
		}
		if rest[1] == "experiments" {
			evidence, err := experiment.BuildEvidence(experiment.DefaultDefinition())
			if err != nil {
				s.writeError(w, http.StatusInternalServerError, "experiment_invalid", err.Error(), false, "", "")
				return
			}
			s.writeJSON(w, http.StatusOK, evidence)
			return
		}
		if rest[1] == "strategy-control" {
			trace := strategyrelease.BuildTrace()
			if err := strategyrelease.Validate(trace); err != nil {
				s.writeError(w, http.StatusInternalServerError, "strategy_invalid", err.Error(), false, "", "")
				return
			}
			s.writeJSON(w, http.StatusOK, trace)
			return
		}
		if rest[1] == "strategy-assurance" {
			cycle := assurance.BuildCycle()
			if err := assurance.Validate(cycle); err != nil {
				s.writeError(w, http.StatusInternalServerError, "assurance_invalid", err.Error(), false, "", "")
				return
			}
			s.writeJSON(w, http.StatusOK, cycle)
			return
		}
		if rest[1] == "aese3" {
			program := aese3.BuildProgram()
			if err := aese3.Validate(program); err != nil {
				s.writeError(w, http.StatusInternalServerError, "aese3_invalid", err.Error(), false, "", "")
				return
			}
			s.writeJSON(w, http.StatusOK, program)
			return
		}
		if rest[1] != "genesis" {
			s.writeError(w, http.StatusNotFound, "not_found", "route not found", false, "", "")
			return
		}
		trace := genesis.BuildTrace()
		if err := genesis.ValidateTrace(trace); err != nil {
			s.writeError(w, http.StatusInternalServerError, "genesis_invalid", err.Error(), false, "", "")
			return
		}
		s.writeJSON(w, http.StatusOK, trace)
		return
	case "scenarios":
		if r.Method != http.MethodGet {
			s.writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed", false, "", "")
			return
		}
		s.handleScenarios(ctx, w, r)
		return

	case "runs":
		switch len(rest) {
		case 1:
			if r.Method != http.MethodPost {
				s.writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed", false, "", "")
				return
			}
			s.handleRunCreate(ctx, w, r)
			return
		case 2:
			if rest[1] == "plan" {
				if r.Method != http.MethodPost {
					s.writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed", false, "", "")
					return
				}
				s.handleRunPlan(ctx, w, r)
				return
			}
			if r.Method != http.MethodGet {
				s.writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed", false, "", "")
				return
			}
			s.handleRunStatus(ctx, w, r, rest[1])
			return
		case 3:
			if r.Method != http.MethodPost {
				s.writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed", false, "", "")
				return
			}
			s.handleRunAction(ctx, w, r, rest[1], rest[2])
			return
		}
	}

	s.writeError(w, http.StatusNotFound, "not_found", "route not found", false, "", "")
}

func (s *Server) handleScenarios(ctx context.Context, w http.ResponseWriter, _ *http.Request) {
	pack, err := loadPack(s.cfg.PackDir)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, "pack_load_failed", err.Error(), false, "", "")
		return
	}
	summary := scenariopack.Inspect(pack)
	s.writeJSON(w, http.StatusOK, map[string]any{
		"pack_key":        summary.PackKey,
		"pack_version":    summary.PackVersion,
		"tenant_template": summary.TenantTemplate,
		"stories":         summary.Stories,
		"entities":        summary.Entities,
		"pack_dir":        s.cfg.PackDir,
	})
	_ = ctx
}

func (s *Server) handleRunPlan(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var payload struct {
		StoryKey string `json:"story_key"`
		PackDir  string `json:"pack_dir"`
	}
	if err := decodeRequestBody(r, s.cfg.BodyLimit, &payload); err != nil {
		s.writeError(w, http.StatusBadRequest, "invalid_request", err.Error(), false, "", "")
		return
	}
	packDir := firstNonEmpty(payload.PackDir, s.cfg.PackDir)
	pack, err := loadPack(packDir)
	if err != nil {
		s.writeError(w, http.StatusBadRequest, "pack_load_failed", err.Error(), false, "", "")
		return
	}
	storyKey := firstNonEmpty(payload.StoryKey, firstStoryKey(pack))
	if storyKey == "" {
		s.writeError(w, http.StatusBadRequest, "invalid_request", "story_key is required", false, "", "")
		return
	}
	plan, err := application.CompilePlan(pack, storyKey)
	if err != nil {
		s.writeError(w, http.StatusBadRequest, "invalid_story", err.Error(), false, "", "")
		return
	}
	s.writeJSON(w, http.StatusOK, plan)
	_ = ctx
}

func (s *Server) handleRunCreate(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var payload runCreateRequest
	if err := decodeRequestBody(r, s.cfg.BodyLimit, &payload); err != nil {
		s.writeError(w, http.StatusBadRequest, "invalid_request", err.Error(), false, "", "")
		return
	}

	token, tokenErr := extractToken(r, payload.Token)
	if tokenErr != nil {
		s.writeError(w, http.StatusUnauthorized, "auth_required", tokenErr.Error(), false, "", "")
		return
	}
	runTarget := strings.TrimSpace(payload.Target)
	if runTarget == "" {
		s.writeError(w, http.StatusBadRequest, "invalid_request", "target is required", false, "", "")
		return
	}
	client, err := application.NewIAOSClient(application.ClientConfig{BaseURL: runTarget, Token: token})
	if err != nil {
		s.writeError(w, http.StatusUnauthorized, "auth_invalid", err.Error(), false, "", "")
		return
	}
	profile, err := client.Profile(ctx)
	if err != nil {
		s.writeErrorFromAPI(w, err, "")
		return
	}

	packDir := firstNonEmpty(payload.PackDir, s.cfg.PackDir)
	pack, err := loadPack(packDir)
	if err != nil {
		s.writeError(w, http.StatusBadRequest, "pack_load_failed", err.Error(), false, "", "")
		return
	}
	storyKey := firstNonEmpty(payload.StoryKey, firstStoryKey(pack))
	if storyKey == "" {
		s.writeError(w, http.StatusBadRequest, "invalid_request", "story_key is required", false, "", "")
		return
	}
	plan, err := application.CompilePlan(pack, storyKey)
	if err != nil {
		s.writeError(w, http.StatusBadRequest, "invalid_story", err.Error(), false, "", "")
		return
	}
	if payload.PlanHash != "" && payload.PlanHash != plan.PlanHash {
		s.writeError(w, http.StatusConflict, "plan_hash_mismatch", "plan_hash is stale", true, "", "")
		return
	}
	if strings.TrimSpace(payload.Target) == "" {
		s.writeError(w, http.StatusBadRequest, "invalid_request", "target is required", false, "", "")
		return
	}
	tenant := firstNonEmpty(payload.Tenant, pack.Manifest.TenantTemplate)
	if tenant == "" {
		s.writeError(w, http.StatusBadRequest, "invalid_request", "tenant is required", false, "", "")
		return
	}
	if strings.TrimSpace(profile.TenantID) != "" && profile.TenantID != tenant {
		s.writeError(w, http.StatusForbidden, "tenant_mismatch", "token tenant does not match request tenant", false, "", "")
		return
	}

	runID := firstNonEmpty(payload.RunID, application.EffectiveRunID("", "api"))
	now := time.Now().UTC()
	run := &runRecord{
		RunID:       runID,
		PackKey:     pack.Manifest.PackKey,
		PackVersion: pack.Manifest.PackVersion,
		PackDir:     packDir,
		ScenarioKey: storyKey,
		Plan:        plan,
		Status:      application.RunStatusPlanned,
		CurrentAct:  0,
		TenantID:    tenant,
		Target:      runTarget,
		Token:       token,
		Actor:       firstNonEmpty(payload.Actor, profile.Username, "aese-user"),
		CreatedAt:   now,
		UpdatedAt:   now,
		ActionCache: map[string]actionCache{},
	}
	initialCursor, cursorErr := s.refreshRunCursor(ctx, run.Target, run.Token, run.TenantID, run.PackKey, run.ScenarioKey)
	if cursorErr != nil {
		if apiErr := mapIAOSError(cursorErr, "read"); apiErr != nil {
			if apiErr.code != "action_failed" && apiErr.code != "not_found" {
				s.writeError(w, apiErr.statusCode, apiErr.code, apiErr.message, apiErr.retryable, runID, "", apiErr.requiredPermission)
				return
			}
			if apiErr.code == "action_failed" {
				s.writeError(w, http.StatusInternalServerError, "run_cursor_load_failed", cursorErr.Error(), false, runID, "")
				return
			}
			s.logf("scenario snapshot missing, create run with zero cursor: run=%s", run.RunID)
		}
	}
	run.Cursor = initialCursor

	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.claimRun(run); err != nil {
		s.writeErrorFromAPI(w, err, runID)
		return
	}
	if _, exists := s.runs[run.RunID]; exists {
		s.writeError(w, http.StatusConflict, "conflict", "run id already exists", false, run.RunID, "")
		return
	}
	s.runs[run.RunID] = run
	s.order = append(s.order, run.RunID)

	response := toRunResponse(run, nil)
	response.Plan = plan
	s.writeJSON(w, http.StatusCreated, response)
	_ = ctx
}

func (s *Server) handleRunStatus(ctx context.Context, w http.ResponseWriter, r *http.Request, runID string) {
	s.mu.Lock()
	run := s.runs[runID]
	if run == nil {
		s.mu.Unlock()
		s.writeError(w, http.StatusNotFound, "run_not_found", "run not found", false, "", "")
		return
	}

	tokenFromHeader, tokenErr := extractBearerToken(r)
	if tokenErr != nil {
		s.mu.Unlock()
		s.writeError(w, http.StatusUnauthorized, "auth_invalid", tokenErr.Error(), false, run.RunID, "")
		return
	}
	if tokenFromHeader == "" {
		s.mu.Unlock()
		s.writeError(w, http.StatusUnauthorized, "auth_required", "run has no operator token", false, run.RunID, "")
		return
	}
	if tokenFromHeader != run.Token {
		s.mu.Unlock()
		s.writeError(w, http.StatusForbidden, "token_mismatch", "token does not match run operator", false, run.RunID, "")
		return
	}

	client, err := application.NewIAOSClient(application.ClientConfig{BaseURL: run.Target, Token: run.Token, TenantID: run.TenantID})
	if err != nil {
		s.mu.Unlock()
		s.writeError(w, http.StatusUnauthorized, "auth_invalid", err.Error(), false, run.RunID, "")
		return
	}
	if _, apiErr := s.loadProfileForRun(ctx, client, run); apiErr != nil {
		s.mu.Unlock()
		s.writeError(w, apiErr.statusCode, apiErr.code, apiErr.message, apiErr.retryable, run.RunID, runVersion(run), apiErr.requiredPermission)
		return
	}
	if refreshErr := s.refreshRunFromFacts(ctx, run, client); refreshErr != nil {
		s.logf("refresh run state failed: run=%s err=%v", run.RunID, refreshErr)
		s.mu.Unlock()
		if apiErr := mapIAOSError(refreshErr, "read"); apiErr != nil && apiErr.code != "action_failed" {
			s.writeError(w, apiErr.statusCode, apiErr.code, apiErr.message, apiErr.retryable, run.RunID, runVersion(run), apiErr.requiredPermission)
			return
		}
		s.writeError(w, http.StatusInternalServerError, "run_refresh_failed", refreshErr.Error(), false, run.RunID, runVersion(run))
		return
	}
	response := toRunResponse(run, nil)
	s.mu.Unlock()

	s.writeJSON(w, http.StatusOK, response)
}

func (s *Server) handleRunAction(ctx context.Context, w http.ResponseWriter, r *http.Request, runID, action string) {
	var payload runActionRequest
	if err := decodeRequestBody(r, s.cfg.BodyLimit, &payload); err != nil {
		s.writeError(w, http.StatusBadRequest, "invalid_request", err.Error(), false, "", "")
		return
	}

	apply := true
	if payload.DryRun {
		apply = false
	}
	if payload.Apply != nil {
		apply = *payload.Apply
	}

	idempotency := strings.TrimSpace(payload.IdempotencyKey)
	if idempotency == "" {
		idempotency = strings.TrimSpace(r.Header.Get("Idempotency-Key"))
	}

	s.mu.Lock()
	run := s.runs[runID]
	if run == nil {
		s.mu.Unlock()
		s.writeError(w, http.StatusNotFound, "run_not_found", "run not found", false, "", "")
		return
	}

	if run.Token == "" {
		s.mu.Unlock()
		s.writeError(w, http.StatusUnauthorized, "auth_required", "run has no operator token", false, run.RunID, "")
		return
	}

	tokenFromHeader, tokenErr := extractBearerToken(r)
	if tokenErr != nil {
		s.mu.Unlock()
		s.writeError(w, http.StatusUnauthorized, "auth_invalid", tokenErr.Error(), false, run.RunID, "")
		return
	}
	if tokenFromHeader == "" {
		s.mu.Unlock()
		s.writeError(w, http.StatusUnauthorized, "auth_required", "run has no operator token", false, run.RunID, "")
		return
	}
	if tokenFromHeader != run.Token {
		s.mu.Unlock()
		s.writeError(w, http.StatusForbidden, "token_mismatch", "token does not match run operator", false, run.RunID, "")
		return
	}

	client, err := application.NewIAOSClient(application.ClientConfig{BaseURL: run.Target, Token: tokenFromHeader, TenantID: run.TenantID})
	if err != nil {
		s.mu.Unlock()
		s.writeError(w, http.StatusUnauthorized, "auth_invalid", err.Error(), false, run.RunID, runVersion(run))
		return
	}
	if _, apiErr := s.loadProfileForRun(ctx, client, run); apiErr != nil {
		s.mu.Unlock()
		s.writeError(w, apiErr.statusCode, apiErr.code, apiErr.message, apiErr.retryable, run.RunID, runVersion(run), apiErr.requiredPermission)
		return
	}

	if payload.RunVersion != "" && payload.RunVersion != runVersion(run) {
		s.mu.Unlock()
		s.writeError(w, http.StatusConflict, "run_version_mismatch", "run_version is stale", true, run.RunID, runVersion(run))
		return
	}
	if payload.PlanHash != "" && payload.PlanHash != run.Plan.PlanHash {
		s.mu.Unlock()
		s.writeError(w, http.StatusConflict, "plan_hash_mismatch", "plan_hash is stale", true, run.RunID, runVersion(run))
		return
	}
	if payload.ExpectedCursor != nil && run.Cursor != *payload.ExpectedCursor {
		s.mu.Unlock()
		s.writeError(w, http.StatusConflict, "cursor_mismatch", "run cursor is stale", true, run.RunID, runVersion(run))
		return
	}

	if actionRequiresIdempotency(action, apply) && idempotency == "" {
		s.mu.Unlock()
		s.writeError(w, http.StatusBadRequest, "idempotency_required", "idempotency key is required for write actions", true, run.RunID, runVersion(run))
		return
	}
	cacheKey := runActionCacheKey(action, idempotency)
	if cacheKey != "" {
		if cached, ok := run.ActionCache[cacheKey]; ok {
			runCopy := *run
			if cached.HttpStatus >= 400 {
				s.mu.Unlock()
				if cached.ErrorCode != "" {
					s.writeError(w, cached.HttpStatus, cached.ErrorCode, cached.Error, cached.ErrorRetryable, run.RunID, runVersion(run), cached.RequiredPermission)
					return
				}
				s.writeError(w, cached.HttpStatus, "action_failed", "previous action execution failed", false, run.RunID, runVersion(run))
				return
			}
			response := toRunResponse(&runCopy, cached.Outcome)
			response.Plan = run.Plan
			s.mu.Unlock()
			s.writeJSON(w, cached.HttpStatus, actionResponse{Run: response, Action: action})
			return
		}
	}

	pack, err := loadPack(run.PackDir)
	if err != nil {
		s.mu.Unlock()
		s.logf("load pack failed: %v", err)
		s.writeError(w, http.StatusInternalServerError, "pack_load_failed", err.Error(), false, run.RunID, runVersion(run))
		return
	}

	story, storyErr := application.FindStory(pack, run.ScenarioKey)
	if storyErr != nil {
		s.mu.Unlock()
		s.writeError(w, http.StatusInternalServerError, "story_missing", storyErr.Error(), false, run.RunID, runVersion(run))
		return
	}

	allowed := make(map[string]struct{}, 8)
	for _, allowedAction := range application.AllowedActions(run.Status) {
		allowed[string(allowedAction)] = struct{}{}
	}
	if _, ok := allowed[action]; !ok {
		s.mu.Unlock()
		s.writeError(w, http.StatusConflict, "invalid_state", fmt.Sprintf("action %q not allowed in status %q", action, run.Status), true, run.RunID, runVersion(run))
		return
	}

	var (
		outcome any
		runErr  *apiError
	)

	actor := firstNonEmpty(run.Actor, "aese-user")
	ctxWithCancel, cancel := context.WithTimeout(ctx, s.cfg.RequestTimeout)
	defer cancel()

	switch action {
	case "preflight":
		outcome, runErr = s.executePreflight(ctxWithCancel, client, pack, run, story, apply)
	case string(application.RunActionInitialize):
		outcome, runErr = s.executeInitialize(ctxWithCancel, client, pack, run, story, apply, actor)
	case string(application.RunActionAdvance):
		outcome, runErr = s.executeAdvance(ctxWithCancel, client, pack, run, story, apply, actor)
	case string(application.RunActionRunToEnd):
		outcome, runErr = s.executeRunToEnd(ctxWithCancel, client, pack, run, story, apply, actor)
	case string(application.RunActionAnalyze):
		outcome, runErr = s.executeAnalyze(ctxWithCancel, client, pack, run, story, apply, actor)
	case string(application.RunActionVerify):
		outcome, runErr = s.executeVerify(ctxWithCancel, client, pack, run, story, apply, actor)
	case "reset-plan":
		outcome, runErr = s.executeResetPlan(ctxWithCancel, client, pack, run, story)
	case string(application.RunActionReset):
		outcome, runErr = s.executeReset(ctxWithCancel, client, pack, run, story, apply, firstNonEmpty(payload.ConfirmationToken, r.Header.Get("X-Aese-Reset-Token")))
	default:
		s.mu.Unlock()
		s.writeError(w, http.StatusNotFound, "unsupported_action", "unsupported action", false, run.RunID, runVersion(run))
		return
	}

	httpStatus := http.StatusOK
	if runErr != nil {
		httpStatus = runErr.statusCode
		run.LastError = runErr.message
		run.Retryable = runErr.retryable
		run.UpdatedAt = time.Now().UTC()
		if cacheKey != "" {
			run.ActionCache[cacheKey] = actionCache{
				HttpStatus:         httpStatus,
				Action:             action,
				Idempotency:        idempotency,
				ErrorCode:          runErr.code,
				Error:              runErr.message,
				ErrorRetryable:     runErr.retryable,
				RequiredPermission: runErr.requiredPermission,
				UpdatedAt:          run.UpdatedAt,
				RunStatus:          string(run.Status),
				Cursor:             run.Cursor,
			}
		}
		s.mu.Unlock()
		s.writeError(w, httpStatus, runErr.code, runErr.message, runErr.retryable, run.RunID, runVersion(run), runErr.requiredPermission)
		return
	}

	run.LastError = ""
	run.Retryable = false
	run.UpdatedAt = time.Now().UTC()
	if refreshErr := s.refreshRunFromFacts(ctxWithCancel, run, client); refreshErr != nil {
		if snapshotErr := s.refreshSnapshot(ctxWithCancel, run); snapshotErr != nil {
			s.logf("refresh snapshot failed: run=%s err=%v", run.RunID, snapshotErr)
		}
		s.logf("refresh run facts failed: run=%s err=%v", run.RunID, refreshErr)
	}
	if cacheKey != "" {
		run.ActionCache[cacheKey] = actionCache{
			HttpStatus:  httpStatus,
			Action:      action,
			Idempotency: idempotency,
			Outcome:     outcome,
			UpdatedAt:   run.UpdatedAt,
			RunStatus:   string(run.Status),
			Cursor:      run.Cursor,
		}
	}
	s.mu.Unlock()

	response := toRunResponse(run, outcome)
	response.Plan = run.Plan
	s.writeJSON(w, http.StatusOK, actionResponse{Run: response, Action: action})
}

func (s *Server) loadProfileForRun(ctx context.Context, client *iaosclient.Client, run *runRecord) (*iaosclient.ProfileResponse, *apiError) {
	profile, err := client.Profile(ctx)
	if err != nil {
		return nil, mapIAOSError(err, "read")
	}
	if run.TenantID != "" {
		if profile.TenantID != "" && profile.TenantID != run.TenantID {
			return &profile, &apiError{
				statusCode: http.StatusForbidden,
				code:       "tenant_mismatch",
				message:    "token tenant does not match run tenant",
				retryable:  false,
			}
		}
		run.TenantID = firstNonEmpty(profile.TenantID, run.TenantID)
	}
	if profile.Username != "" {
		run.Actor = firstNonEmpty(profile.Username, run.Actor)
	}
	return &profile, nil
}

func (s *Server) refreshRunCursor(ctx context.Context, target, token, tenantID, packKey, scenarioKey string) (int64, error) {
	client, err := application.NewIAOSClient(application.ClientConfig{BaseURL: target, Token: token, TenantID: tenantID})
	if err != nil {
		return 0, err
	}
	snapshot, err := client.ScenarioSnapshot(ctx, packKey, scenarioKey)
	if err != nil {
		return 0, err
	}
	return snapshot.Cursor, nil
}

func (s *Server) refreshRunFromFacts(ctx context.Context, run *runRecord, client *iaosclient.Client) error {
	snapshot, err := client.ScenarioSnapshot(ctx, run.PackKey, run.ScenarioKey)
	if err != nil {
		if isIAOSNotFoundError(err) {
			s.logf("scenario snapshot missing on refresh, keep local run state: run=%s", run.RunID)
			return nil
		}
		return err
	}
	baseCursor := run.Cursor
	cursor := maxInt64(baseCursor, snapshot.Cursor)
	events := make([]iaosclient.ScenarioObservedEvent, 0, len(snapshot.Events))
	for _, raw := range snapshot.Events {
		event, ok := parseSnapshotObservedEvent(raw)
		if !ok {
			continue
		}
		if event.Cursor <= baseCursor {
			continue
		}
		if strings.TrimSpace(run.Plan.Correlation) != "" && event.CorrelationID != run.Plan.Correlation {
			continue
		}
		events = append(events, event)
		if event.Cursor > cursor {
			cursor = event.Cursor
		}
	}

	if run.Plan.Correlation != "" {
		pageAfter := baseCursor
		for {
			queryResp, queryErr := client.ScenarioEvents(ctx, run.PackKey, run.ScenarioKey, pageAfter, 200, run.Plan.Correlation)
			if queryErr != nil {
				s.logf("load scenario events failed: run=%s err=%v", run.RunID, queryErr)
				break
			}
			observed := queryResp.Items
			for _, event := range observed {
				if event.Cursor <= baseCursor {
					continue
				}
				if strings.TrimSpace(run.Plan.Correlation) != "" && event.CorrelationID != run.Plan.Correlation {
					continue
				}
				events = append(events, event)
				if event.Cursor > cursor {
					cursor = event.Cursor
				}
			}
			if !queryResp.HasMore {
				break
			}
			if queryResp.NextCursor <= pageAfter {
				break
			}
			pageAfter = queryResp.NextCursor
		}
	}
	completedActs := inferCompletedActsFromFacts(run.Plan, events, run.CurrentAct)
	if completedActs > run.Plan.ActCount {
		completedActs = run.Plan.ActCount
	}
	if completedActs > run.CurrentAct {
		run.CurrentAct = completedActs
	}

	hasRecommendations := false
	var rawRecommendations []json.RawMessage
	rawRecommendations = append(rawRecommendations, snapshot.Recommendations...)
	if recommendations, recErr := client.ScenarioRecommendations(ctx, run.PackKey, run.ScenarioKey); recErr == nil {
		rawRecommendations = append(rawRecommendations, recommendations.Items...)
	} else {
		s.logf("load scenario recommendations failed: run=%s err=%v", run.RunID, recErr)
	}
	if runHasRecommendationForRun(rawRecommendations, run.RunID) {
		hasRecommendations = true
	}

	run.Status = inferRunStatusFromFacts(run.Status, completedActs, run.Plan.ActCount, hasRecommendations)
	run.Cursor = cursor
	run.UpdatedAt = time.Now().UTC()
	return nil
}

func parseSnapshotObservedEvent(raw map[string]any) (iaosclient.ScenarioObservedEvent, bool) {
	if len(raw) == 0 {
		return iaosclient.ScenarioObservedEvent{}, false
	}
	var event iaosclient.ScenarioObservedEvent
	encoded, err := json.Marshal(raw)
	if err != nil {
		return iaosclient.ScenarioObservedEvent{}, false
	}
	if err := json.Unmarshal(encoded, &event); err != nil {
		return iaosclient.ScenarioObservedEvent{}, false
	}
	return event, true
}

func maxInt64(a, b int64) int64 {
	if a >= b {
		return a
	}
	return b
}

func (s *Server) executePreflight(ctx context.Context, client *iaosclient.Client, pack *scenariopack.Pack, run *runRecord, story scenariopack.Story, _ bool) (any, *apiError) {
	summary, warnings, err := application.ApplyScenario(ctx, client, pack, run.ScenarioKey, run.RunID, false)
	if err != nil {
		return nil, mapIAOSError(err, string(application.RunActionPreflight))
	}

	entityContracts := []map[string]any{}
	entities := map[string]struct{}{}
	for _, set := range pack.RecordSets {
		entities[set.Entity] = struct{}{}
	}
	for _, set := range story.Initial.RecordSets {
		entities[set.Entity] = struct{}{}
	}
	for entity := range entities {
		schema, err := client.Schema(ctx, entity)
		if err != nil {
			if isIAOSNotFoundError(err) {
				s.logf("scenario preflight schema missing, skipping entity contract: run=%s entity=%s", run.RunID, entity)
				warnings = append(warnings, legacyprojection.Warning{
					SourceEntity: entity,
					Message:      fmt.Sprintf("schema for entity %s is missing, skipping contract preview", entity),
				})
				continue
			}
			return nil, mapIAOSError(err, string(application.RunActionPreflight))
		}
		entityContracts = append(entityContracts, map[string]any{
			"entity":      entity,
			"schema":      schema.Entity,
			"fields":      schema.Fields,
			"permissions": schema.Permissions,
			"version":     schema.Version,
			"storage":     schema.StorageStrategy,
		})
	}

	transition, transErr := application.NextStatus(run.Status, application.RunActionPreflight, application.RunTransitionContext{})
	if transErr != nil {
		return nil, &apiError{statusCode: http.StatusConflict, code: "invalid_state", message: transErr.Error(), retryable: true}
	}
	run.Status = transition

	return map[string]any{
		"action":             string(application.RunActionPreflight),
		"plan_hash":          run.Plan.PlanHash,
		"initialize_dry_run": summary,
		"warnings":           warnings,
		"entity_contracts":   entityContracts,
		"story_key":          story.Ref.Key,
	}, nil
}

func (s *Server) executeInitialize(ctx context.Context, client *iaosclient.Client, pack *scenariopack.Pack, run *runRecord, _ scenariopack.Story, apply bool, actor string) (any, *apiError) {
	record, _, err := application.ApplyScenario(ctx, client, pack, run.ScenarioKey, run.RunID, apply)
	if err != nil {
		return nil, mapIAOSError(err, string(application.RunActionInitialize))
	}
	transition, transErr := application.NextStatus(run.Status, application.RunActionInitialize, application.RunTransitionContext{CurrentAct: run.CurrentAct, TotalActs: run.Plan.ActCount})
	if transErr != nil {
		return nil, &apiError{statusCode: http.StatusConflict, code: "invalid_state", message: transErr.Error(), retryable: true}
	}
	run.Status = transition
	_ = actor
	return map[string]any{
		"action": string(application.RunActionInitialize),
		"actor":  actor,
		"summary": map[string]any{
			"dry_run":     record.DryRun,
			"inserted":    record.Inserted,
			"updated":     record.Updated,
			"no_op":       record.NoOp,
			"conflicts":   record.Conflicts,
			"unsupported": record.Unsupported,
		},
		"run_id":      run.RunID,
		"correlation": record.CorrelationID,
	}, nil
}

func (s *Server) executeAdvance(ctx context.Context, client *iaosclient.Client, pack *scenariopack.Pack, run *runRecord, story scenariopack.Story, apply bool, actor string) (any, *apiError) {
	targetStage, stageErr := stageEventIDs(run.Plan, run.CurrentAct)
	if stageErr != nil {
		return nil, &apiError{statusCode: http.StatusConflict, code: "invalid_state", message: stageErr.Error(), retryable: true}
	}
	events, mapErr := filterEvents(story.Events.Events, targetStage)
	if mapErr != nil {
		return nil, &apiError{statusCode: http.StatusInternalServerError, code: "scenario_event_mismatch", message: mapErr.Error(), retryable: false}
	}
	newStory := story
	newStory.Events.Events = events
	summary, replayErr := replayEvents(ctx, client, newStory, run, actor, apply)
	if replayErr != nil {
		return nil, mapIAOSError(replayErr, string(application.RunActionAdvance))
	}
	run.CurrentAct++
	transition, transErr := application.NextStatus(run.Status, application.RunActionAdvance, application.RunTransitionContext{CurrentAct: run.CurrentAct - 1, TotalActs: run.Plan.ActCount})
	if transErr != nil {
		return nil, &apiError{statusCode: http.StatusConflict, code: "invalid_state", message: transErr.Error(), retryable: true}
	}
	run.Status = transition
	return summary, nil
}

func (s *Server) executeRunToEnd(ctx context.Context, client *iaosclient.Client, _ *scenariopack.Pack, run *runRecord, story scenariopack.Story, apply bool, actor string) (any, *apiError) {
	if run.Status != application.RunStatusReady && run.Status != application.RunStatusRunning {
		return nil, &apiError{statusCode: http.StatusConflict, code: "invalid_state", message: "run-to-end only valid when status is ready or running", retryable: true}
	}
	if run.CurrentAct >= run.Plan.ActCount {
		return nil, &apiError{statusCode: http.StatusConflict, code: "invalid_state", message: "all acts already completed", retryable: true}
	}

	targetEventIDs := make([]string, 0)
	for actIndex := run.CurrentAct; actIndex < run.Plan.ActCount; actIndex++ {
		eventIDs, err := stageEventIDs(run.Plan, actIndex)
		if err != nil {
			return nil, &apiError{statusCode: http.StatusConflict, code: "invalid_state", message: err.Error(), retryable: true}
		}
		targetEventIDs = append(targetEventIDs, eventIDs...)
	}
	events, mapErr := filterEvents(story.Events.Events, targetEventIDs)
	if mapErr != nil {
		return nil, &apiError{statusCode: http.StatusInternalServerError, code: "scenario_event_mismatch", message: mapErr.Error(), retryable: false}
	}
	newStory := story
	newStory.Events.Events = events
	summary, replayErr := replayEvents(ctx, client, newStory, run, actor, apply)
	if replayErr != nil {
		return nil, mapIAOSError(replayErr, string(application.RunActionRunToEnd))
	}
	run.CurrentAct = run.Plan.ActCount
	run.Status = application.RunStatusAwaitingAnalysis
	return summary, nil
}

func (s *Server) executeAnalyze(ctx context.Context, client *iaosclient.Client, pack *scenariopack.Pack, run *runRecord, _ scenariopack.Story, apply bool, actor string) (any, *apiError) {
	summary, err := application.RunAgents(ctx, client, pack, run.ScenarioKey, run.RunID, apply)
	if err != nil {
		return nil, mapIAOSError(err, string(application.RunActionAnalyze))
	}
	_ = actor
	transition, transErr := application.NextStatus(run.Status, application.RunActionAnalyze, application.RunTransitionContext{CurrentAct: run.CurrentAct, TotalActs: run.Plan.ActCount})
	if transErr != nil {
		return nil, &apiError{statusCode: http.StatusConflict, code: "invalid_state", message: transErr.Error(), retryable: true}
	}
	run.Status = transition
	return summary, nil
}

func (s *Server) executeVerify(ctx context.Context, client *iaosclient.Client, pack *scenariopack.Pack, run *runRecord, _ scenariopack.Story, apply bool, actor string) (any, *apiError) {
	summary, err := application.VerifyScenario(ctx, client, pack, run.ScenarioKey, application.VerifyOptions{
		Target: run.Target,
		Tenant: run.TenantID,
		Actor:  actor,
	})
	if err != nil {
		run.Retryable = true
		return summary, &apiError{statusCode: http.StatusConflict, code: "verify_failed", message: err.Error(), retryable: true}
	}
	if apply {
		transition, transErr := application.NextStatus(run.Status, application.RunActionVerify, application.RunTransitionContext{CurrentAct: run.CurrentAct, TotalActs: run.Plan.ActCount})
		if transErr != nil {
			return nil, &apiError{statusCode: http.StatusConflict, code: "invalid_state", message: transErr.Error(), retryable: true}
		}
		run.Status = transition
	}
	_ = apply
	return summary, nil
}

func (s *Server) executeResetPlan(ctx context.Context, client *iaosclient.Client, pack *scenariopack.Pack, run *runRecord, _ scenariopack.Story) (any, *apiError) {
	transition, transErr := application.NextStatus(run.Status, application.RunActionResetPlan, application.RunTransitionContext{})
	if transErr != nil {
		return nil, &apiError{statusCode: http.StatusConflict, code: "invalid_state", message: transErr.Error(), retryable: true}
	}
	summary, err := application.ResetScenario(ctx, client, pack, run.ScenarioKey, run.RunID, false)
	if err != nil {
		return nil, mapIAOSError(err, string(application.RunActionResetPlan))
	}
	token, tokenErr := resetToken()
	if tokenErr != nil {
		return nil, &apiError{statusCode: http.StatusInternalServerError, code: "internal_error", message: "failed to generate reset confirmation token", retryable: false}
	}
	run.ResetToken = token
	run.ResetTokenExpiresAt = time.Now().UTC().Add(resetTokenTTL)
	run.Status = transition
	return map[string]any{
		"action":                   string(application.RunActionResetPlan),
		"summary":                  summary,
		"reset_confirmation_token": token,
		"confirmation_expires_at":  run.ResetTokenExpiresAt.Format(time.RFC3339),
	}, nil
}

func (s *Server) executeReset(ctx context.Context, client *iaosclient.Client, pack *scenariopack.Pack, run *runRecord, _ scenariopack.Story, apply bool, confirmation string) (any, *apiError) {
	if run.Status != application.RunStatusResetting {
		return nil, &apiError{statusCode: http.StatusConflict, code: "invalid_state", message: "reset requires reset-plan to be requested first", retryable: true}
	}
	if run.ResetToken == "" || run.ResetTokenExpiresAt.IsZero() || run.ResetTokenExpiresAt.Before(time.Now().UTC()) {
		return nil, &apiError{statusCode: http.StatusForbidden, code: "reset_confirmation_invalid", message: "reset confirmation token is missing or expired", retryable: false}
	}
	if strings.TrimSpace(confirmation) == "" || strings.TrimSpace(confirmation) != run.ResetToken {
		return nil, &apiError{statusCode: http.StatusForbidden, code: "reset_confirmation_mismatch", message: "reset confirmation token mismatch", retryable: false}
	}
	summary, err := application.ResetScenario(ctx, client, pack, run.ScenarioKey, run.RunID, apply)
	if err != nil {
		return nil, mapIAOSError(err, string(application.RunActionReset))
	}
	if apply {
		transition, transErr := application.NextStatus(run.Status, application.RunActionReset, application.RunTransitionContext{})
		if transErr != nil {
			return nil, &apiError{statusCode: http.StatusConflict, code: "invalid_state", message: transErr.Error(), retryable: true}
		}
		run.Status = transition
		run.CurrentAct = 0
		run.ResetToken = ""
		run.ResetTokenExpiresAt = time.Time{}
	}
	return map[string]any{
		"action":      string(application.RunActionReset),
		"summary":     summary,
		"applied":     apply,
		"run_version": runVersion(run),
	}, nil
}

func (s *Server) writeErrorFromAPI(w http.ResponseWriter, err error, runID string) {
	var apiErr *apiError
	if errors.As(err, &apiErr) {
		s.writeError(w, apiErr.statusCode, apiErr.code, apiErr.message, apiErr.retryable, runID, "", apiErr.requiredPermission)
		return
	}
	s.writeError(w, http.StatusConflict, "conflict", err.Error(), false, runID, "")
}

func (s *Server) claimRun(run *runRecord) error {
	currentRunKey := runKey(run.TenantID, run.PackKey, run.ScenarioKey)
	for _, existing := range s.runs {
		existingKey := runKey(existing.TenantID, existing.PackKey, existing.ScenarioKey)
		if currentRunKey != existingKey {
			continue
		}
		if existing.RunID == run.RunID {
			continue
		}
		switch existing.Status {
		case application.RunStatusCompleted, application.RunStatusFailed:
			continue
		default:
			return &apiError{statusCode: http.StatusConflict, code: "conflict", message: fmt.Sprintf("another writable run exists for tenant/story %s", currentRunKey)}
		}
	}
	return nil
}

func (s *Server) refreshSnapshot(ctx context.Context, run *runRecord) error {
	client, err := application.NewIAOSClient(application.ClientConfig{BaseURL: run.Target, Token: run.Token, TenantID: run.TenantID})
	if err != nil {
		return err
	}
	snapshot, err := client.ScenarioSnapshot(ctx, run.PackKey, run.ScenarioKey)
	if err != nil {
		return err
	}
	run.Cursor = snapshot.Cursor
	return nil
}

func toRunResponse(run *runRecord, outcome any) runResponse {
	allowed := make([]string, 0, len(application.AllowedActions(run.Status)))
	for _, action := range application.AllowedActions(run.Status) {
		allowed = append(allowed, string(action))
	}
	response := runResponse{
		RunID:                     run.RunID,
		RunVersion:                runVersion(run),
		PackKey:                   run.PackKey,
		PackVersion:               run.PackVersion,
		ScenarioKey:               run.ScenarioKey,
		PlanHash:                  run.Plan.PlanHash,
		Status:                    string(run.Status),
		CurrentAct:                run.CurrentAct,
		TotalActs:                 run.Plan.ActCount,
		Cursor:                    run.Cursor,
		TenantID:                  run.TenantID,
		Target:                    run.Target,
		CreatedAt:                 run.CreatedAt.Format(time.RFC3339),
		UpdatedAt:                 run.UpdatedAt.Format(time.RFC3339),
		AllowedActions:            allowed,
		LastError:                 run.LastError,
		Retryable:                 run.Retryable,
		Outcome:                   outcome,
		Plan:                      nil,
		ResetConfirmationRequired: run.ResetToken != "" && !run.ResetTokenExpiresAt.IsZero() && run.ResetTokenExpiresAt.After(time.Now().UTC()),
	}
	return response
}

func runVersion(run *runRecord) string {
	return strconv.FormatInt(run.UpdatedAt.UnixNano(), 10)
}

func actionRequiresIdempotency(action string, apply bool) bool {
	if !apply {
		return false
	}
	switch action {
	case string(application.RunActionPreflight), string(application.RunActionResetPlan):
		return false
	default:
		return true
	}
}

func runActionCacheKey(action, idempotency string) string {
	idempotency = strings.TrimSpace(idempotency)
	if idempotency == "" {
		return ""
	}
	return action + ":" + idempotency
}

func inferCompletedActsFromFacts(plan application.Plan, events []iaosclient.ScenarioObservedEvent, alreadyCompleted int) int {
	observed := map[string]struct{}{}
	for _, event := range events {
		if strings.TrimSpace(event.EventID) != "" {
			observed[event.EventID] = struct{}{}
		}
	}
	completed := alreadyCompleted
	actIndex := 0
	for _, stage := range plan.Stages {
		stageID := string(stage.StageID)
		if !strings.HasPrefix(stageID, "act-") {
			continue
		}
		actIndex++
		if actIndex <= alreadyCompleted {
			continue
		}
		allDone := true
		for _, eventID := range stage.EventIDs {
			if _, ok := observed[eventID]; !ok {
				allDone = false
				break
			}
		}
		if allDone {
			completed = actIndex
			continue
		}
		break
	}
	if completed > len(plan.Stages) {
		completed = len(plan.Stages)
	}
	return completed
}

func runHasRecommendationForRun(items []json.RawMessage, runID string) bool {
	if runID == "" || len(items) == 0 {
		return false
	}
	type recommendationEnvelope struct {
		RunID string `json:"run_id"`
	}
	for _, raw := range items {
		var item recommendationEnvelope
		if err := json.Unmarshal(raw, &item); err != nil {
			continue
		}
		if strings.TrimSpace(item.RunID) == runID {
			return true
		}
	}
	return false
}

func inferRunStatusFromFacts(current application.RunStatus, completedActs, totalActs int, hasRunRecommendations bool) application.RunStatus {
	if current == application.RunStatusCompleted || current == application.RunStatusResetting || current == application.RunStatusReset {
		return current
	}
	if completedActs >= totalActs {
		if hasRunRecommendations {
			if current == application.RunStatusAnalyzing {
				return application.RunStatusAwaitingVerification
			}
			if current == application.RunStatusAwaitingVerification {
				return application.RunStatusAwaitingVerification
			}
			if current == application.RunStatusCompleted {
				return application.RunStatusCompleted
			}
			return application.RunStatusAwaitingAnalysis
		}
		switch current {
		case application.RunStatusAwaitingVerification, application.RunStatusAnalyzing:
			return application.RunStatusAnalyzing
		case application.RunStatusAwaitingAnalysis:
			return application.RunStatusAwaitingAnalysis
		case application.RunStatusReady, application.RunStatusRunning:
			return application.RunStatusAwaitingAnalysis
		default:
			return application.RunStatusAwaitingAnalysis
		}
	}
	if completedActs > 0 {
		switch current {
		case application.RunStatusPlanned, application.RunStatusInitializing:
			return application.RunStatusReady
		case application.RunStatusReady, application.RunStatusRunning, application.RunStatusAwaitingAnalysis, application.RunStatusAnalyzing, application.RunStatusAwaitingVerification:
			return current
		default:
			return current
		}
	}
	return current
}

func splitPath(value string) []string {
	path := strings.Trim(value, "/")
	if path == "" {
		return nil
	}
	return strings.Split(path, "/")
}

func decodeRequestBody(r *http.Request, limit int64, dst any) error {
	if r.Body == nil {
		return fmt.Errorf("request body is required")
	}
	defer r.Body.Close()
	reader := io.LimitReader(r.Body, limit+1)
	data, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("read request body: %w", err)
	}
	if int64(len(data)) > limit {
		return fmt.Errorf("request body exceeds limit of %d bytes", limit)
	}
	if len(bytesTrim(data)) == 0 {
		return fmt.Errorf("empty request body")
	}
	if err := json.Unmarshal(data, dst); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}
	return nil
}

func bytesTrim(data []byte) []byte {
	start := 0
	for start < len(data) && (data[start] == '\n' || data[start] == '\r' || data[start] == ' ' || data[start] == '\t') {
		start++
	}
	end := len(data)
	for end > start && (data[end-1] == '\n' || data[end-1] == '\r' || data[end-1] == ' ' || data[end-1] == '\t') {
		end--
	}
	return data[start:end]
}

func extractToken(r *http.Request, bodyToken string) (string, error) {
	if strings.TrimSpace(bodyToken) != "" {
		return bodyToken, nil
	}
	return extractBearerToken(r)
}

func extractBearerToken(r *http.Request) (string, error) {
	authorization := strings.TrimSpace(r.Header.Get("Authorization"))
	if authorization == "" {
		return "", nil
	}
	parts := strings.SplitN(authorization, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
		return "", fmt.Errorf("Authorization must be Bearer token")
	}
	return strings.TrimSpace(parts[1]), nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func stageEventIDs(plan application.Plan, actIndex int) ([]string, error) {
	if actIndex < 0 {
		return nil, fmt.Errorf("act index %d invalid", actIndex)
	}
	stageID := "act-" + strconv.Itoa(actIndex+1)
	for _, stage := range plan.Stages {
		if string(stage.StageID) == stageID {
			return stage.EventIDs, nil
		}
	}
	return nil, fmt.Errorf("stage %s not found", stageID)
}

func filterEvents(events []scenariopack.Event, ids []string) ([]scenariopack.Event, error) {
	index := make(map[string]scenariopack.Event, len(events))
	for _, event := range events {
		index[event.EventID] = event
	}
	filtered := make([]scenariopack.Event, 0, len(ids))
	for _, id := range ids {
		event, ok := index[id]
		if !ok {
			return nil, fmt.Errorf("event %s not found in story", id)
		}
		filtered = append(filtered, event)
	}
	if len(filtered) == 0 {
		return nil, fmt.Errorf("empty event list")
	}
	return filtered, nil
}

func loadPack(packDir string) (*scenariopack.Pack, error) {
	return scenariopack.Load(packDir)
}

func firstStoryKey(pack *scenariopack.Pack) string {
	if len(pack.Stories) == 0 {
		return ""
	}
	return pack.Stories[0].Ref.Key
}

func replayEvents(ctx context.Context, client *iaosclient.Client, story scenariopack.Story, run *runRecord, actor string, apply bool) (replay.ReplaySummary, error) {
	runner, err := replay.New(client)
	if err != nil {
		return replay.ReplaySummary{}, err
	}
	opts := replay.Options{Apply: apply, Target: run.Target, Tenant: run.TenantID, Actor: actor, PackKey: run.PackKey, OrderID: "", Entities: nil}
	return runner.Replay(ctx, story, opts)
}

func mapIAOSError(err error, action string) *apiError {
	if err == nil {
		return &apiError{statusCode: http.StatusInternalServerError, code: "internal_error", message: "unknown error"}
	}
	if errors.Is(err, context.Canceled) {
		return &apiError{statusCode: http.StatusRequestTimeout, code: "request_canceled", message: "request canceled", retryable: false}
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return &apiError{statusCode: http.StatusGatewayTimeout, code: "upstream_timeout", message: "upstream request timeout", retryable: true}
	}
	var apiErr *iaosclient.APIError
	if errors.As(err, &apiErr) {
		requiredPermission := apiErr.RequiredPermission
		if requiredPermission == "" {
			requiredPermission = scenarioRunPermission(action)
		}
		switch apiErr.StatusCode {
		case http.StatusUnauthorized:
			return &apiError{statusCode: http.StatusUnauthorized, code: "auth_invalid", message: apiErr.Error(), retryable: false, requiredPermission: requiredPermission}
		case http.StatusForbidden:
			return &apiError{statusCode: http.StatusForbidden, code: "forbidden", message: apiErr.Error(), retryable: false, requiredPermission: requiredPermission}
		case http.StatusNotFound:
			return &apiError{statusCode: http.StatusNotFound, code: "not_found", message: apiErr.Error(), retryable: true, requiredPermission: requiredPermission}
		case http.StatusConflict, http.StatusLocked:
			return &apiError{statusCode: apiErr.StatusCode, code: "conflict", message: apiErr.Error(), retryable: true, requiredPermission: requiredPermission}
		default:
			if apiErr.StatusCode >= 500 {
				return &apiError{statusCode: apiErr.StatusCode, code: "upstream_error", message: apiErr.Error(), retryable: true, requiredPermission: requiredPermission}
			}
			return &apiError{statusCode: apiErr.StatusCode, code: "bad_request", message: apiErr.Error(), retryable: false, requiredPermission: requiredPermission}
		}
	}
	return &apiError{statusCode: http.StatusBadRequest, code: "action_failed", message: err.Error(), retryable: false}
}

func isIAOSNotFoundError(err error) bool {
	var apiErr *iaosclient.APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusNotFound
	}
	return false
}

func scenarioRunPermission(action string) string {
	switch action {
	case "read":
		return "scenario.run.read"
	case string(application.RunActionReset):
		return "scenario.run.reset"
	case string(application.RunActionPreflight), string(application.RunActionInitialize), string(application.RunActionAdvance), string(application.RunActionRunToEnd),
		string(application.RunActionAnalyze), string(application.RunActionVerify), string(application.RunActionResetPlan):
		return "scenario.run.execute"
	default:
		return "scenario.run.read"
	}
}

func runKey(tenant, pack, scenario string) string {
	return strings.Join([]string{tenant, pack, scenario}, "/")
}

func resetToken() (string, error) {
	buffer := make([]byte, resetTokenByteSize)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}
	return hex.EncodeToString(buffer), nil
}

func (s *Server) writeJSON(w http.ResponseWriter, status int, payload any) {
	buf, err := json.Marshal(payload)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, "internal_error", "failed to encode response", false, "", "")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(buf)
}

func (s *Server) writeError(w http.ResponseWriter, status int, code string, message string, retryable bool, runID, runVersion string, requiredPermission ...string) {
	permission := firstNonEmpty(requiredPermission...)
	s.writeJSON(w, status, errorResponse{
		Error:              message,
		Code:               code,
		Retryable:          retryable,
		Status:             strconv.Itoa(status),
		RunID:              runID,
		RunVersion:         runVersion,
		RequiredPermission: permission,
	})
}
