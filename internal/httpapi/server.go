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
	"github.com/industrial-ai/iaos-aese/internal/creative"
	"github.com/industrial-ai/iaos-aese/internal/experiment"
	"github.com/industrial-ai/iaos-aese/internal/firstdelivery"
	"github.com/industrial-ai/iaos-aese/internal/gameprojection"
	"github.com/industrial-ai/iaos-aese/internal/genesis"
	"github.com/industrial-ai/iaos-aese/internal/genesisworkspace"
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
	PackDir                  string
	IAOSBaseURL              string
	RequestTimeout           time.Duration
	BodyLimit                int64
	Logger                   *log.Logger
	CreativeProvider         creative.Provider
	CreativeJobStore         *creative.JobStore
	PlantPlanningProvider    plantbuild.PlanningProvider
	GenesisWorkspaceService  *genesisworkspace.Service
	GenesisPlayerAuth        *genesisworkspace.PlayerAuthClient
	AllowLocalGenesisAuth    bool
	AllowLocalPlantAuthority bool
}

type Server struct {
	cfg                     Config
	mux                     *http.ServeMux
	logf                    func(format string, args ...any)
	mu                      sync.RWMutex
	runs                    map[string]*runRecord
	order                   []string
	creativeProvider        creative.Provider
	creativeJobStore        *creative.JobStore
	plantPlanningProvider   plantbuild.PlanningProvider
	genesisWorkspaceService *genesisworkspace.Service
	genesisPlayerAuth       *genesisworkspace.PlayerAuthClient
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
	if cfg.CreativeProvider == nil {
		cfg.CreativeProvider = creative.DeterministicProvider{}
	}
	if cfg.PlantPlanningProvider == nil {
		cfg.PlantPlanningProvider = plantbuild.UnconfiguredPlanningProvider{}
	}
	server := &Server{
		cfg:                     cfg,
		runs:                    map[string]*runRecord{},
		logf:                    cfg.Logger.Printf,
		mux:                     http.NewServeMux(),
		creativeProvider:        cfg.CreativeProvider,
		creativeJobStore:        cfg.CreativeJobStore,
		plantPlanningProvider:   cfg.PlantPlanningProvider,
		genesisWorkspaceService: cfg.GenesisWorkspaceService,
		genesisPlayerAuth:       cfg.GenesisPlayerAuth,
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
			"/api/aese/v1/game/incorporation/:case/projection",
			"/api/aese/v1/game/creative/intent",
			"/api/aese/v1/game/creative/names",
			"/api/aese/v1/genesis/workspaces",
			"/api/aese/v1/genesis/workspaces/:workspace/session",
			"/api/aese/v1/world/plant-build",
			"/api/aese/v1/world/plant-build/planning-status",
			"/api/aese/v1/world/plant-build/financial-constraints",
			"/api/aese/v1/world/plant-build/requirements/:requirement_id",
			"/api/aese/v1/world/plant-build/proposals",
			"/api/aese/v1/world/plant-build/reviews",
			"/api/aese/v1/world/plant-build/investigations",
			"/api/aese/v1/world/plant-build/observations",
			"/api/aese/v1/world/plant-build/site-selections",
			"/api/aese/v1/world/plant-build/site-selections/finalize",
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
	case "commands":
		s.handleIAOSCommandGateway(ctx, w, r, rest)
		return
	case "auth":
		if s.genesisPlayerAuth == nil || len(rest) != 2 {
			s.writeError(w, http.StatusNotFound, "not_found", "route not found", false, "", "")
			return
		}
		switch {
		case rest[1] == "register" && r.Method == http.MethodPost:
			var request genesisworkspace.PlayerRegisterRequest
			if err := decodeStrictRequestBody(r, s.cfg.BodyLimit, &request); err != nil {
				s.writeError(w, http.StatusBadRequest, "invalid_registration_request", err.Error(), false, "", "")
				return
			}
			result, err := s.genesisPlayerAuth.Register(ctx, request)
			if err != nil {
				s.writeGenesisAuthError(w, err, "registration_failed")
				return
			}
			s.writeJSON(w, http.StatusCreated, result)
			return
		case rest[1] == "login" && r.Method == http.MethodPost:
			var request genesisworkspace.PlayerLoginRequest
			if err := decodeStrictRequestBody(r, s.cfg.BodyLimit, &request); err != nil {
				s.writeError(w, http.StatusBadRequest, "invalid_login_request", err.Error(), false, "", "")
				return
			}
			result, err := s.genesisPlayerAuth.Login(ctx, request)
			if err != nil {
				s.writeGenesisAuthError(w, err, "login_failed")
				return
			}
			s.writeJSON(w, http.StatusOK, result)
			return
		case rest[1] == "session" && r.Method == http.MethodGet:
			token := bearerToken(r.Header.Get("Authorization"))
			if token == "" {
				s.writeError(w, http.StatusUnauthorized, "player_session_required", "Genesis Player session required", false, "", "")
				return
			}
			result, err := s.genesisPlayerAuth.Session(ctx, token)
			if err != nil {
				s.writeGenesisAuthError(w, err, "session_lookup_failed")
				return
			}
			s.writeJSON(w, http.StatusOK, result)
			return
		default:
			s.writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed", false, "", "")
			return
		}
	case "genesis":
		if len(rest) < 2 || rest[1] != "workspaces" || s.genesisWorkspaceService == nil {
			s.writeError(w, http.StatusNotFound, "not_found", "route not found", false, "", "")
			return
		}
		token := bearerToken(r.Header.Get("Authorization"))
		playerID := ""
		if token != "" && s.genesisPlayerAuth != nil {
			session, err := s.genesisPlayerAuth.Session(ctx, token)
			if err != nil {
				s.writeGenesisAuthError(w, err, "player_session_expired")
				return
			}
			playerID = strings.TrimSpace(session.Player.SubjectID)
			ctx = genesisworkspace.WithIAOSToken(ctx, token)
		} else if s.cfg.AllowLocalGenesisAuth {
			playerID = strings.TrimSpace(r.Header.Get("X-Genesis-Player-Id"))
			if token != "" {
				ctx = genesisworkspace.WithIAOSToken(ctx, token)
			}
		}
		if playerID == "" {
			s.writeError(w, http.StatusUnauthorized, "player_session_required", "authenticated Genesis Player session required", false, "", "")
			return
		}
		if len(rest) == 4 && rest[3] == "session" && r.Method == http.MethodPost {
			result, err := s.genesisWorkspaceService.RefreshSession(ctx, playerID, rest[2])
			if err != nil {
				s.writeGenesisWorkspaceError(w, err, "workspace_session_failed")
				return
			}
			s.writeJSON(w, http.StatusOK, result)
			return
		}
		if len(rest) != 2 {
			s.writeError(w, http.StatusNotFound, "not_found", "route not found", false, "", "")
			return
		}
		switch r.Method {
		case http.MethodGet:
			items, err := s.genesisWorkspaceService.List(ctx, playerID)
			if err != nil {
				s.writeGenesisWorkspaceError(w, err, "workspace_list_failed")
				return
			}
			s.writeJSON(w, http.StatusOK, map[string]any{"items": items})
			return
		case http.MethodPost:
			var request genesisworkspace.CreateRequest
			if err := decodeStrictRequestBody(r, s.cfg.BodyLimit, &request); err != nil {
				s.writeError(w, http.StatusBadRequest, "invalid_workspace_request", err.Error(), false, "", "")
				return
			}
			request.OwnerPlayerID = playerID
			result, err := s.genesisWorkspaceService.Create(ctx, request)
			if err != nil {
				s.writeGenesisWorkspaceError(w, err, "workspace_provisioning_failed")
				return
			}
			s.writeJSON(w, http.StatusCreated, result)
			return
		default:
			s.writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed", false, "", "")
			return
		}
	case "game":
		if len(rest) == 3 && rest[1] == "creative" && rest[2] == "status" && r.Method == http.MethodGet {
			status := creative.ProviderStatus{State: "not_configured", Provider: "none", PromptVersion: "genesis-naming-v1"}
			if provider, ok := s.creativeProvider.(creative.StatusProvider); ok {
				status = provider.ProviderStatus()
			}
			s.writeJSON(w, http.StatusOK, status)
			return
		}
		if len(rest) == 3 && rest[1] == "creative" && rest[2] == "jobs" && r.Method == http.MethodGet {
			if strings.TrimSpace(r.Header.Get("Authorization")) == "" || s.creativeJobStore == nil {
				s.writeError(w, http.StatusUnauthorized, "creative_job_access_denied", "authenticated workspace session required", false, "", "")
				return
			}
			tenantID := strings.TrimSpace(r.Header.Get("X-IAOS-Tenant-Id"))
			if tenantID == "" {
				s.writeError(w, http.StatusBadRequest, "creative_job_tenant_required", "X-IAOS-Tenant-Id is required", false, "", "")
				return
			}
			if err := s.validateCreativeSession(ctx, r, tenantID); err != nil {
				s.writeError(w, http.StatusForbidden, "creative_job_access_denied", err.Error(), false, "", "")
				return
			}
			items, err := s.creativeJobStore.List(tenantID, strings.TrimSpace(r.URL.Query().Get("case")))
			if err != nil {
				s.writeError(w, http.StatusInternalServerError, "creative_job_read_failed", err.Error(), true, "", "")
				return
			}
			s.writeJSON(w, http.StatusOK, map[string]any{"items": items})
			return
		}
		if len(rest) == 3 && rest[1] == "creative" && r.Method == http.MethodPost {
			provider := s.creativeProvider
			switch rest[2] {
			case "intent":
				var request creative.FounderIntentRequest
				if err := decodeStrictRequestBody(r, s.cfg.BodyLimit, &request); err != nil {
					s.writeError(w, http.StatusBadRequest, "invalid_creative_request", err.Error(), false, "", "")
					return
				}
				if err := s.validateCreativeSession(ctx, r, request.TenantID); err != nil {
					s.writeError(w, http.StatusForbidden, "creative_session_invalid", err.Error(), false, "", "")
					return
				}
				intent, err := provider.AnalyzeIntent(ctx, request)
				if err != nil {
					s.writeError(w, http.StatusBadGateway, "creative_provider_failed", err.Error(), true, "", "")
					return
				}
				s.writeJSON(w, http.StatusOK, intent)
				return
			case "names":
				var intent creative.FounderIntent
				if err := decodeStrictRequestBody(r, s.cfg.BodyLimit, &intent); err != nil {
					s.writeError(w, http.StatusBadRequest, "invalid_creative_request", err.Error(), false, "", "")
					return
				}
				if err := s.validateCreativeSession(ctx, r, intent.TenantID); err != nil {
					s.writeError(w, http.StatusForbidden, "creative_session_invalid", err.Error(), false, "", "")
					return
				}
				providerStatus := creative.ProviderStatus{State: "not_configured", Provider: "unknown", PromptVersion: "genesis-naming-v1"}
				if statusProvider, ok := provider.(creative.StatusProvider); ok {
					providerStatus = statusProvider.ProviderStatus()
				}
				inputHash := creative.Hash(intent)
				jobID := "creative-" + strings.TrimPrefix(inputHash, "sha256:")[:24]
				started := time.Now().UTC()
				job := creative.CreativeJob{
					JobID: jobID, TenantID: intent.TenantID, CaseCode: intent.CaseCode,
					WorkspaceID:   strings.TrimSpace(r.Header.Get("X-Genesis-Workspace-Id")),
					CorrelationID: "corr-" + intent.CaseCode,
					Kind:          "company_naming", Status: "running", Provider: providerStatus.Provider,
					Model: providerStatus.Model, ModelVersion: providerStatus.Model,
					Prompt: providerStatus.PromptVersion, PromptVersion: providerStatus.PromptVersion,
					BaseURLHost: providerStatus.BaseURLHost, InputHash: inputHash,
					RequestID: jobID, CreatedAt: started.Format(time.RFC3339),
					TokenUsage: map[string]int{}, ValidationResult: "pending",
				}
				if providerStatus.State == "fallback" {
					job.FallbackReason = "external_model_not_configured"
				}
				if s.creativeJobStore != nil {
					existing, replay, err := s.creativeJobStore.Begin(job)
					if errors.Is(err, creative.ErrJobRunning) {
						s.writeError(w, http.StatusConflict, "creative_job_running", "a naming job is already running for this input", true, "", "")
						return
					}
					if err != nil {
						s.writeError(w, http.StatusInternalServerError, "creative_job_write_failed", err.Error(), true, "", "")
						return
					}
					if replay {
						s.writeJSON(w, http.StatusOK, map[string]any{
							"intent_id": intent.IntentID, "status": "candidate_only",
							"proposals": existing.Parameters["proposals"], "creative_job": existing,
							"provider_status": providerStatus, "idempotent_replay": true,
						})
						return
					}
				}
				generationEvidence := creative.GenerationEvidence{TokenUsage: map[string]int{}}
				generationCtx := creative.WithGenerationEvidence(ctx, func(e creative.GenerationEvidence) {
					if e.RequestID != "" {
						generationEvidence.RequestID = e.RequestID
					}
					if e.FinishReason != "" {
						generationEvidence.FinishReason = e.FinishReason
					}
					for key, value := range e.TokenUsage {
						generationEvidence.TokenUsage[key] += value
					}
				})
				proposals, err := provider.GenerateNames(generationCtx, intent)
				if err != nil {
					job.Status = "failed"
					job.Error = err.Error()
					job.ValidationResult = "failed"
					job.LatencyMS = time.Since(started).Milliseconds()
					job.CompletedAt = time.Now().UTC().Format(time.RFC3339)
					if s.creativeJobStore != nil {
						_ = s.creativeJobStore.Save(job)
					}
					s.writeError(w, http.StatusBadGateway, "creative_provider_failed", err.Error(), true, "", "")
					return
				}
				job.Status = "completed"
				if generationEvidence.RequestID != "" {
					job.RequestID = generationEvidence.RequestID
				}
				job.FinishReason = generationEvidence.FinishReason
				if job.FinishReason == "" {
					job.FinishReason = "stop"
				}
				job.TokenUsage = generationEvidence.TokenUsage
				job.ValidationResult = "valid"
				job.LatencyMS = time.Since(started).Milliseconds()
				job.ContentHash = creative.Hash(proposals)
				job.Parameters = map[string]any{"proposals": proposals}
				job.CompletedAt = time.Now().UTC().Format(time.RFC3339)
				if s.creativeJobStore != nil {
					if err := s.creativeJobStore.Save(job); err != nil {
						s.writeError(w, http.StatusInternalServerError, "creative_job_write_failed", err.Error(), true, "", "")
						return
					}
				}
				s.writeJSON(w, http.StatusOK, map[string]any{
					"intent_id":       intent.IntentID,
					"status":          "candidate_only",
					"proposals":       proposals,
					"creative_job":    job,
					"provider_status": providerStatus,
				})
				return
			}
		}
		if len(rest) != 4 || rest[1] != "incorporation" || rest[3] != "projection" || r.Method != http.MethodGet {
			s.writeError(w, http.StatusNotFound, "not_found", "route not found", false, "", "")
			return
		}
		if s.cfg.IAOSBaseURL != "" {
			token := strings.TrimSpace(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer "))
			if token == "" {
				s.writeError(w, http.StatusUnauthorized, "iaos_token_required", "Bearer token is required for the live IAOS projection", false, "", "")
				return
			}
			client, err := iaosclient.New(iaosclient.Config{
				BaseURL: s.cfg.IAOSBaseURL, Token: token, TenantID: r.Header.Get("X-IAOS-Tenant-Id"),
			})
			if err != nil {
				s.writeError(w, http.StatusInternalServerError, "iaos_projection_config_invalid", err.Error(), false, "", "")
				return
			}
			trace, err := client.IncorporationTrace(ctx, rest[2])
			if err != nil {
				if iaosclient.IsStatus(err, http.StatusNotFound) {
					s.writeError(w, http.StatusNotFound, "incorporation_case_not_found", "incorporation case not found", false, "", "")
					return
				}
				s.writeError(w, http.StatusBadGateway, "iaos_trace_unavailable", err.Error(), true, "", "")
				return
			}
			items, err := client.IncorporationWorkItems(ctx, rest[2])
			if err != nil {
				s.writeError(w, http.StatusBadGateway, "iaos_work_items_unavailable", err.Error(), true, "", "")
				return
			}
			projection, err := gameprojection.FromIAOS(trace, items)
			if err != nil {
				s.writeError(w, http.StatusUnprocessableEntity, "iaos_projection_invalid", err.Error(), false, "", "")
				return
			}
			opening, err := client.FinanceOpening(ctx, rest[2])
			if err != nil {
				s.writeError(w, http.StatusBadGateway, "iaos_finance_opening_unavailable", err.Error(), true, "", "")
				return
			}
			if err := gameprojection.AttachFinanceOpening(&projection, opening); err != nil {
				s.writeError(w, http.StatusUnprocessableEntity, "iaos_finance_projection_invalid", err.Error(), false, "", "")
				return
			}
			s.writeJSON(w, http.StatusOK, projection)
			return
		}
		frame := len(incorporation.BuildTrace().Frames) - 1
		if raw := r.URL.Query().Get("frame"); raw != "" {
			value, err := strconv.Atoi(raw)
			if err != nil {
				s.writeError(w, http.StatusBadRequest, "invalid_frame", "frame must be an integer", false, "", "")
				return
			}
			frame = value
		}
		projection, err := gameprojection.FromIncorporationTrace(incorporation.BuildTrace(), rest[2], frame)
		if err != nil {
			s.writeError(w, http.StatusBadRequest, "projection_invalid", err.Error(), false, "", "")
			return
		}
		s.writeJSON(w, http.StatusOK, projection)
		return
	case "world":
		if len(rest) == 3 && rest[1] == "plant-build" && rest[2] == "planning-status" && r.Method == http.MethodGet {
			s.writeJSON(w, http.StatusOK, s.plantPlanningProvider.Status())
			return
		}
		if len(rest) == 3 && rest[1] == "plant-build" && rest[2] == "financial-constraints" && r.Method == http.MethodGet {
			tenantID := firstNonEmptyString(r.Header.Get("X-IAOS-Tenant-Id"), r.Header.Get("X-Tenant-ID"))
			caseCode := strings.TrimSpace(r.URL.Query().Get("case_code"))
			token := bearerToken(r.Header.Get("Authorization"))
			if token == "" || tenantID == "" {
				s.writeError(w, http.StatusUnauthorized, "plant_financial_identity_required", "IAOS bearer token and tenant context are required", false, "", "")
				return
			}
			client, err := iaosclient.New(iaosclient.Config{BaseURL: s.cfg.IAOSBaseURL, Token: token, TenantID: tenantID})
			if err != nil {
				s.writeError(w, http.StatusServiceUnavailable, "plant_financial_gateway_unavailable", err.Error(), true, "", "")
				return
			}
			result, err := client.PlantFinancialConstraint(ctx, caseCode)
			if err != nil {
				var apiErr *iaosclient.APIError
				if errors.As(err, &apiErr) {
					s.writeError(w, apiErr.StatusCode, firstNonEmptyString(apiErr.ErrorCode, "plant_financial_constraint_unavailable"), apiErr.Message, apiErr.StatusCode >= 500, "", apiErr.RequiredPermission)
					return
				}
				s.writeError(w, http.StatusBadGateway, "plant_financial_constraint_unavailable", err.Error(), true, "", "")
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(result)
			return
		}
		if len(rest) == 4 && rest[1] == "plant-build" && rest[2] == "requirements" && r.Method == http.MethodGet {
			tenantID := firstNonEmptyString(r.Header.Get("X-IAOS-Tenant-Id"), r.Header.Get("X-Tenant-ID"))
			token := bearerToken(r.Header.Get("Authorization"))
			if token == "" || tenantID == "" {
				s.writeError(w, http.StatusUnauthorized, "plant_requirement_identity_required", "IAOS bearer token and tenant context are required", false, "", "")
				return
			}
			client, err := iaosclient.New(iaosclient.Config{BaseURL: s.cfg.IAOSBaseURL, Token: token, TenantID: tenantID})
			if err != nil {
				s.writeError(w, http.StatusServiceUnavailable, "plant_requirement_gateway_unavailable", err.Error(), true, "", "")
				return
			}
			result, err := client.PlantRequirement(ctx, rest[3])
			if err != nil {
				var apiErr *iaosclient.APIError
				if errors.As(err, &apiErr) {
					s.writeError(w, apiErr.StatusCode, firstNonEmptyString(apiErr.ErrorCode, "plant_requirement_unavailable"), apiErr.Message, apiErr.StatusCode >= 500, "", apiErr.RequiredPermission)
					return
				}
				s.writeError(w, http.StatusBadGateway, "plant_requirement_unavailable", err.Error(), true, "", "")
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(result)
			return
		}
		if len(rest) == 3 && rest[1] == "plant-build" && rest[2] == "proposals" && r.Method == http.MethodGet {
			tenantID := firstNonEmptyString(r.Header.Get("X-IAOS-Tenant-Id"), r.Header.Get("X-Tenant-ID"))
			token := bearerToken(r.Header.Get("Authorization"))
			if token == "" || tenantID == "" {
				s.writeError(w, http.StatusUnauthorized, "plant_proposal_identity_required", "IAOS bearer token and tenant context are required", false, "", "")
				return
			}
			client, err := iaosclient.New(iaosclient.Config{BaseURL: s.cfg.IAOSBaseURL, Token: token, TenantID: tenantID})
			if err != nil {
				s.writeError(w, http.StatusServiceUnavailable, "plant_proposal_gateway_unavailable", err.Error(), true, "", "")
				return
			}
			result, err := client.PlantProposalSet(ctx, r.URL.Query().Get("requirement_id"))
			if err != nil {
				var apiErr *iaosclient.APIError
				if errors.As(err, &apiErr) {
					s.writeError(w, apiErr.StatusCode, firstNonEmptyString(apiErr.ErrorCode, "plant_proposal_unavailable"), apiErr.Message, apiErr.StatusCode >= 500, "", apiErr.RequiredPermission)
					return
				}
				s.writeError(w, http.StatusBadGateway, "plant_proposal_unavailable", err.Error(), true, "", "")
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(result)
			return
		}
		if len(rest) == 4 && rest[1] == "plant-build" && rest[2] == "proposals" && rest[3] == "manual" && r.Method == http.MethodPost {
			var input struct {
				RequirementID    string                        `json:"requirement_id"`
				ProposalSetID    string                        `json:"proposal_set_id"`
				ExpectedRevision int                           `json:"expected_revision"`
				Proposal         plantbuild.SiteOptionProposal `json:"proposal"`
			}
			if err := decodeStrictRequestBody(r, s.cfg.BodyLimit, &input); err != nil {
				s.writeError(w, http.StatusBadRequest, "invalid_manual_site_proposal", err.Error(), false, "", "")
				return
			}
			tenantID := firstNonEmptyString(r.Header.Get("X-IAOS-Tenant-Id"), r.Header.Get("X-Tenant-ID"))
			authority, actorID, err := s.plantAuthorityClient(ctx, r, tenantID)
			if err != nil || authority == nil {
				s.writeError(w, http.StatusUnauthorized, "plant_planning_identity_invalid", firstNonEmptyString(errorString(err), "IAOS authority is required"), false, "", "")
				return
			}
			requirementRaw, err := authority.PlantRequirement(ctx, input.RequirementID)
			if err != nil {
				s.writePlantAuthorityError(w, err, "facility_requirement_unavailable")
				return
			}
			var requirement plantbuild.FacilityRequirement
			if err := json.Unmarshal(requirementRaw, &requirement); err != nil {
				s.writeError(w, http.StatusBadGateway, "facility_requirement_invalid", err.Error(), false, "", "")
				return
			}
			var previous plantbuild.ProposalSet
			hasPrevious := false
			setRaw, err := authority.PlantProposalSet(ctx, input.RequirementID)
			if err != nil {
				var apiErr *iaosclient.APIError
				if !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusNotFound || input.ExpectedRevision != 0 {
					s.writePlantAuthorityError(w, err, "site_proposal_previous_revision_unavailable")
					return
				}
			} else {
				if err := json.Unmarshal(setRaw, &previous); err != nil || previous.ProposalSetID != input.ProposalSetID || previous.Revision != input.ExpectedRevision || input.ExpectedRevision < 1 {
					s.writeError(w, http.StatusConflict, "site_proposal_revision_conflict", "候选集已变化，请刷新后重新提交人工候选", false, "", "")
					return
				}
				hasPrevious = true
			}
			proposalID, err := resetToken()
			if err != nil {
				s.writeError(w, http.StatusInternalServerError, "manual_site_proposal_id_failed", err.Error(), true, "", "")
				return
			}
			input.Proposal.ProposalID = "manual-site-" + proposalID
			input.Proposal.Status = "proposed"
			input.Proposal.SourceRefs = append(input.Proposal.SourceRefs, "human:"+actorID)
			next := plantbuild.ProposalSet{
				SchemaVersion: "1.0", ProposalSetID: "manual-proposal-set-" + requirement.RequirementID,
				RequirementID: requirement.RequirementID, Revision: 1, Status: "candidate_only",
				Proposals: []plantbuild.SiteOptionProposal{input.Proposal},
			}
			inputHash := plantbuild.CanonicalHash(requirement)
			if hasPrevious {
				next = previous
				next.Revision = previous.Revision + 1
				next.Proposals = append(append([]plantbuild.SiteOptionProposal{}, previous.Proposals...), input.Proposal)
				inputHash = plantbuild.CanonicalHash(previous)
			}
			next.Evidence = plantbuild.ProposalEvidence{
				Provider: "human", Model: "manual-entry", PromptVersion: "manual-candidate-v1",
				SourceType: "human_manual", ParentRevision: previous.Revision,
				InputHash: inputHash, OutputHash: plantbuild.CanonicalHash(next.Proposals),
				ValidatedAt: time.Now().UTC().Format(time.RFC3339),
			}
			if err := plantbuild.ValidateProposalSet(requirement, next); err != nil {
				s.writeError(w, http.StatusUnprocessableEntity, "invalid_manual_site_proposal", err.Error(), false, "", "")
				return
			}
			if err := postPlantAuthorityCommand(ctx, authority, "site.proposal.record", actorID, requirement.CaseCode, next.ProposalSetID, next.Revision, next); err != nil {
				s.writePlantAuthorityError(w, err, "manual_site_proposal_not_committed")
				return
			}
			s.writeJSON(w, http.StatusCreated, map[string]any{"status": "committed", "proposal_set": next, "manual_proposal_id": input.Proposal.ProposalID})
			return
		}
		if len(rest) == 3 && rest[1] == "plant-build" && rest[2] == "proposals" && r.Method == http.MethodPost {
			var request plantbuild.FacilityRequirement
			if err := decodeStrictRequestBody(r, s.cfg.BodyLimit, &request); err != nil {
				s.writeError(w, http.StatusBadRequest, "invalid_facility_requirement", err.Error(), false, "", "")
				return
			}
			if err := plantbuild.ValidateRequirement(request); err != nil {
				s.writeError(w, http.StatusUnprocessableEntity, "invalid_facility_requirement", err.Error(), false, "", "")
				return
			}
			if err := s.validateCreativeSession(ctx, r, request.TenantID); err != nil {
				s.writeError(w, http.StatusForbidden, "plant_planning_session_invalid", err.Error(), false, "", "")
				return
			}
			authority, actorID, err := s.plantAuthorityClient(ctx, r, request.TenantID)
			if err != nil {
				s.writeError(w, http.StatusUnauthorized, "plant_planning_identity_invalid", err.Error(), false, "", "")
				return
			}
			if authority != nil {
				if err := postPlantAuthorityCommand(ctx, authority, "facility.requirement.define", actorID, request.CaseCode, request.RequirementID, request.Revision, request); err != nil {
					s.writePlantAuthorityError(w, err, "facility_requirement_not_committed")
					return
				}
			}
			providerStatus := s.plantPlanningProvider.Status()
			inputHash := plantbuild.CanonicalHash(request)
			jobID := "plant-planning-" + strings.TrimPrefix(inputHash, "sha256:")[:24]
			started := time.Now().UTC()
			job := creative.CreativeJob{
				JobID: jobID, TenantID: request.TenantID, CaseCode: request.CaseCode,
				WorkspaceID: strings.TrimSpace(r.Header.Get("X-Genesis-Workspace-Id")), CorrelationID: "corr-" + request.CaseCode,
				Kind: "facility_planning", Status: "running", Provider: providerStatus.Provider,
				Model: providerStatus.Model, ModelVersion: providerStatus.Model,
				Prompt: providerStatus.PromptVersion, PromptVersion: providerStatus.PromptVersion,
				InputHash: inputHash, RequestID: jobID, CreatedAt: started.Format(time.RFC3339),
				TokenUsage: map[string]int{}, ValidationResult: "pending",
			}
			if s.creativeJobStore != nil {
				existing, replay, beginErr := s.creativeJobStore.Begin(job)
				if errors.Is(beginErr, creative.ErrJobRunning) {
					s.writeError(w, http.StatusConflict, "plant_planning_job_running", "同一设施需求的 Agent 任务正在运行", true, "", "")
					return
				}
				if beginErr != nil {
					s.writeError(w, http.StatusInternalServerError, "plant_planning_evidence_write_failed", beginErr.Error(), true, "", "")
					return
				}
				if replay {
					set, ok := decodeStoredProposalSet(existing.Parameters["proposal_set"])
					if ok && plantbuild.ValidateProposalSet(request, set) == nil {
						if authority != nil {
							if err := postPlantAuthorityProposalCommand(ctx, authority, actorID, request.CaseCode, set, existing); err != nil {
								s.writePlantAuthorityError(w, err, "site_proposal_not_committed")
								return
							}
						}
						s.writeJSON(w, http.StatusOK, map[string]any{"proposal_set": set, "agent_job": existing, "idempotent_replay": true, "authority_status": "committed"})
						return
					}
					// A completed job created by an older contract may contain a
					// proposal that the current IAOS authority must reject. Replace
					// it with a new running attempt instead of replaying it forever.
					if err := s.creativeJobStore.Save(job); err != nil {
						s.writeError(w, http.StatusInternalServerError, "plant_planning_evidence_write_failed", err.Error(), true, "", "")
						return
					}
				}
			}
			set, err := s.plantPlanningProvider.Generate(ctx, request)
			if errors.Is(err, plantbuild.ErrPlanningModelNotConfigured) {
				job.Status, job.Error, job.ValidationResult = "failed", err.Error(), "failed"
				job.CompletedAt = time.Now().UTC().Format(time.RFC3339)
				if s.creativeJobStore != nil {
					_ = s.creativeJobStore.Save(job)
				}
				s.writeError(w, http.StatusServiceUnavailable, "plant_planning_model_not_configured", "外部设施规划模型未启用；请配置模型或使用人工新增候选表单", false, "", "")
				return
			}
			if err != nil {
				job.Status, job.Error, job.ValidationResult = "failed", err.Error(), "failed"
				job.CompletedAt = time.Now().UTC().Format(time.RFC3339)
				if s.creativeJobStore != nil {
					_ = s.creativeJobStore.Save(job)
				}
				s.writeError(w, http.StatusBadGateway, "plant_planning_provider_failed", err.Error(), true, "", "")
				return
			}
			job.Status, job.ValidationResult = "completed", "valid"
			job.RequestID, job.TokenUsage = set.Evidence.RequestID, set.Evidence.TokenUsage
			job.ContentHash, job.Parameters = set.Evidence.OutputHash, map[string]any{"proposal_set": set}
			job.LatencyMS, job.CompletedAt = time.Since(started).Milliseconds(), time.Now().UTC().Format(time.RFC3339)
			if authority != nil {
				if err := postPlantAuthorityProposalCommand(ctx, authority, actorID, request.CaseCode, set, job); err != nil {
					job.Status = "failed"
					job.Error = "IAOS business commit failed: " + err.Error()
					if s.creativeJobStore != nil {
						_ = s.creativeJobStore.Save(job)
					}
					s.writePlantAuthorityError(w, err, "site_proposal_not_committed")
					return
				}
			}
			if s.creativeJobStore != nil {
				if err := s.creativeJobStore.Save(job); err != nil {
					s.writeError(w, http.StatusInternalServerError, "plant_planning_evidence_write_failed", err.Error(), true, "", "")
					return
				}
			}
			s.writeJSON(w, http.StatusOK, map[string]any{"proposal_set": set, "agent_job": job, "idempotent_replay": false, "authority_status": "committed"})
			return
		}
		if len(rest) == 3 && rest[1] == "plant-build" && rest[2] == "reviews" && r.Method == http.MethodGet {
			tenantID := firstNonEmptyString(r.Header.Get("X-IAOS-Tenant-Id"), r.Header.Get("X-Tenant-ID"))
			authority, _, err := s.plantAuthorityClient(ctx, r, tenantID)
			if err != nil || authority == nil {
				s.writeError(w, http.StatusUnauthorized, "plant_review_identity_invalid", firstNonEmptyString(errorString(err), "IAOS authority is required"), false, "", "")
				return
			}
			result, err := authority.PlantProposalReviews(ctx, r.URL.Query().Get("proposal_set_id"))
			if err != nil {
				s.writePlantAuthorityError(w, err, "plant_reviews_unavailable")
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(result)
			return
		}
		if len(rest) == 3 && rest[1] == "plant-build" && rest[2] == "reviews" && r.Method == http.MethodPost {
			var review plantbuild.ProposalReview
			if err := decodeStrictRequestBody(r, s.cfg.BodyLimit, &review); err != nil {
				s.writeError(w, http.StatusBadRequest, "invalid_site_proposal_review", err.Error(), false, "", "")
				return
			}
			tenantID := firstNonEmptyString(r.Header.Get("X-IAOS-Tenant-Id"), r.Header.Get("X-Tenant-ID"))
			authority, actorID, err := s.plantAuthorityClient(ctx, r, tenantID)
			if err != nil || authority == nil {
				s.writeError(w, http.StatusUnauthorized, "plant_review_identity_invalid", firstNonEmptyString(errorString(err), "IAOS authority is required"), false, "", "")
				return
			}
			review.ReviewedBy = actorID
			review.ReviewedAt = time.Now().UTC().Format(time.RFC3339)
			if err := plantbuild.ValidateReview(review); err != nil {
				s.writeError(w, http.StatusUnprocessableEntity, "invalid_site_proposal_review", err.Error(), false, "", "")
				return
			}
			existingRaw, err := authority.PlantProposalReviews(ctx, review.ProposalSetID)
			if err != nil {
				s.writePlantAuthorityError(w, err, "plant_reviews_unavailable")
				return
			}
			var existingEnvelope struct {
				Items []plantbuild.ProposalReview `json:"items"`
			}
			if err := json.Unmarshal(existingRaw, &existingEnvelope); err != nil {
				s.writeError(w, http.StatusBadGateway, "plant_reviews_invalid", "IAOS review projection is invalid", true, "", "")
				return
			}
			for _, existing := range existingEnvelope.Items {
				if existing.ProposalID != review.ProposalID || existing.ExpectedRevision != review.ExpectedRevision {
					continue
				}
				if existing.Action == review.Action && strings.TrimSpace(existing.Reason) == strings.TrimSpace(review.Reason) {
					s.writeJSON(w, http.StatusOK, map[string]any{"status": "already_committed", "proposal_review": existing, "idempotent_replay": true})
					return
				}
				s.writeError(w, http.StatusConflict, "site_proposal_review_immutable", "该候选版本已提交审阅；如需改变决定，请创建新的候选集修订", false, "", "")
				return
			}
			if err := postPlantAuthorityCommand(ctx, authority, "site.proposal.review", actorID, review.ProposalSetID, review.ProposalID, review.ExpectedRevision, review); err != nil {
				s.writePlantAuthorityError(w, err, "site_proposal_review_not_committed")
				return
			}
			s.writeJSON(w, http.StatusCreated, map[string]any{"status": "committed", "proposal_review": review})
			return
		}
		if len(rest) == 3 && rest[1] == "plant-build" && rest[2] == "investigations" && r.Method == http.MethodGet {
			tenantID := firstNonEmptyString(r.Header.Get("X-IAOS-Tenant-Id"), r.Header.Get("X-Tenant-ID"))
			authority, _, err := s.plantAuthorityClient(ctx, r, tenantID)
			if err != nil || authority == nil {
				s.writeError(w, http.StatusUnauthorized, "plant_investigation_identity_invalid", firstNonEmptyString(errorString(err), "IAOS authority is required"), false, "", "")
				return
			}
			result, err := authority.PlantInvestigations(ctx, r.URL.Query().Get("case_code"))
			if err != nil {
				s.writePlantAuthorityError(w, err, "plant_investigations_unavailable")
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(result)
			return
		}
		if len(rest) == 3 && rest[1] == "plant-build" && rest[2] == "investigations" && r.Method == http.MethodPost {
			var request plantbuild.InvestigationRequest
			if err := decodeStrictRequestBody(r, s.cfg.BodyLimit, &request); err != nil {
				s.writeError(w, http.StatusBadRequest, "invalid_site_investigation_request", err.Error(), false, "", "")
				return
			}
			tenantID := firstNonEmptyString(r.Header.Get("X-IAOS-Tenant-Id"), r.Header.Get("X-Tenant-ID"))
			authority, actorID, err := s.plantAuthorityClient(ctx, r, tenantID)
			if err != nil || authority == nil {
				s.writeError(w, http.StatusUnauthorized, "plant_investigation_identity_invalid", firstNonEmptyString(errorString(err), "IAOS authority is required"), false, "", "")
				return
			}
			request.RequestedBy = actorID
			request.Status = "waiting_world"
			if err := plantbuild.ValidateInvestigationRequest(request); err != nil {
				s.writeError(w, http.StatusUnprocessableEntity, "invalid_site_investigation_request", err.Error(), false, "", "")
				return
			}
			if err := postPlantAuthorityCommand(ctx, authority, "site.investigation.request", actorID, request.CaseCode, request.InvestigationRequestID, request.ExpectedRevision, request); err != nil {
				s.writePlantAuthorityError(w, err, "site_investigation_request_not_committed")
				return
			}
			s.writeJSON(w, http.StatusCreated, map[string]any{"status": "waiting_world", "investigation_request": request})
			return
		}
		if len(rest) == 3 && rest[1] == "plant-build" && rest[2] == "observations" && r.Method == http.MethodPost {
			var input struct {
				CaseCode    string                              `json:"case_code"`
				WorldRunID  string                              `json:"world_run_id"`
				Observation plantbuild.InvestigationObservation `json:"observation"`
			}
			if err := decodeStrictRequestBody(r, s.cfg.BodyLimit, &input); err != nil {
				s.writeError(w, http.StatusBadRequest, "invalid_site_investigation_observation", err.Error(), false, "", "")
				return
			}
			if err := plantbuild.ValidateInvestigationObservation(input.Observation); err != nil {
				s.writeError(w, http.StatusUnprocessableEntity, "invalid_site_investigation_observation", err.Error(), false, "", "")
				return
			}
			tenantID := firstNonEmptyString(r.Header.Get("X-IAOS-Tenant-Id"), r.Header.Get("X-Tenant-ID"))
			authority, actorID, err := s.plantAuthorityClient(ctx, r, tenantID)
			if err != nil || authority == nil {
				s.writeError(w, http.StatusUnauthorized, "plant_observation_identity_invalid", firstNonEmptyString(errorString(err), "IAOS authority is required"), false, "", "")
				return
			}
			correlationID := "corr-m10-" + strings.TrimSpace(input.CaseCode)
			messageID := "world-" + input.Observation.ObservationID
			envelope := map[string]any{
				"schema_version": "1.0", "message_id": messageID, "kind": "observation", "tenant_id": tenantID,
				"world_pack_key": "genesis-plant-planning", "world_pack_version": "1.0.0", "world_run_id": input.WorldRunID,
				"branch_id": "main", "sim_occurred_at": input.Observation.ObservedAt, "correlation_id": correlationID,
				"idempotency_key": messageID, "producer": map[string]string{"system": "aese", "component": "plant-investigation-world"},
				"subject_ref":  map[string]string{"type": "site_investigation_request", "code": input.Observation.InvestigationRequestID},
				"payload_type": "site.investigation.observed.v1", "payload": input.Observation,
			}
			rawEnvelope, _ := json.Marshal(envelope)
			if _, err := authority.PostGovernedCommand(ctx, "api/v1/world-bridge/observations", rawEnvelope); err != nil {
				s.writePlantAuthorityError(w, err, "site_investigation_world_observation_not_accepted")
				return
			}
			commit := map[string]any{"schema_version": "1.0", "investigation_request_id": input.Observation.InvestigationRequestID, "world_message_id": messageID}
			if err := postPlantAuthorityCommand(ctx, authority, "site.investigation.observation.commit", actorID, input.CaseCode, input.Observation.ObservationID, 1, commit); err != nil {
				s.writePlantAuthorityError(w, err, "site_investigation_observation_not_committed")
				return
			}
			s.writeJSON(w, http.StatusCreated, map[string]any{"status": "committed", "world_message_id": messageID, "observation": input.Observation})
			return
		}
		if len(rest) == 3 && rest[1] == "plant-build" && rest[2] == "site-selections" && r.Method == http.MethodGet {
			tenantID := firstNonEmptyString(r.Header.Get("X-IAOS-Tenant-Id"), r.Header.Get("X-Tenant-ID"))
			authority, _, err := s.plantAuthorityClient(ctx, r, tenantID)
			if err != nil || authority == nil {
				s.writeError(w, http.StatusUnauthorized, "plant_selection_identity_invalid", firstNonEmptyString(errorString(err), "IAOS authority is required"), false, "", "")
				return
			}
			result, err := authority.PlantSiteSelections(ctx, r.URL.Query().Get("case_code"))
			if err != nil {
				s.writePlantAuthorityError(w, err, "plant_selections_unavailable")
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(result)
			return
		}
		if len(rest) == 3 && rest[1] == "plant-build" && rest[2] == "site-controls" && r.Method == http.MethodGet {
			tenantID := firstNonEmptyString(r.Header.Get("X-IAOS-Tenant-Id"), r.Header.Get("X-Tenant-ID"))
			authority, _, err := s.plantAuthorityClient(ctx, r, tenantID)
			if err != nil || authority == nil {
				s.writeError(w, http.StatusUnauthorized, "plant_site_control_identity_invalid", firstNonEmptyString(errorString(err), "IAOS authority is required"), false, "", "")
				return
			}
			result, err := authority.PlantSiteControls(ctx, r.URL.Query().Get("case_code"))
			if err != nil {
				s.writePlantAuthorityError(w, err, "plant_site_controls_unavailable")
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(result)
			return
		}
		if len(rest) == 3 && rest[1] == "plant-build" && rest[2] == "site-controls" && r.Method == http.MethodPost {
			var request plantbuild.SiteControlRequest
			if err := decodeStrictRequestBody(r, s.cfg.BodyLimit, &request); err != nil {
				s.writeError(w, http.StatusBadRequest, "invalid_site_control_request", err.Error(), false, "", "")
				return
			}
			tenantID := firstNonEmptyString(r.Header.Get("X-IAOS-Tenant-Id"), r.Header.Get("X-Tenant-ID"))
			authority, actorID, err := s.plantAuthorityClient(ctx, r, tenantID)
			if err != nil || authority == nil {
				s.writeError(w, http.StatusUnauthorized, "plant_site_control_identity_invalid", firstNonEmptyString(errorString(err), "IAOS authority is required"), false, "", "")
				return
			}
			request.RequestedBy = actorID
			request.Status = "waiting_world"
			if err := plantbuild.ValidateSiteControlRequest(request); err != nil {
				s.writeError(w, http.StatusUnprocessableEntity, "invalid_site_control_request", err.Error(), false, "", "")
				return
			}
			if err := postPlantAuthorityCommand(ctx, authority, "site.control.request", actorID, request.CaseCode, request.ControlRequestID, 1, request); err != nil {
				s.writePlantAuthorityError(w, err, "site_control_request_not_committed")
				return
			}
			s.writeJSON(w, http.StatusCreated, map[string]any{"status": "waiting_world", "site_control_request": request})
			return
		}
		if len(rest) == 4 && rest[1] == "plant-build" && rest[2] == "site-controls" && rest[3] == "observations" && r.Method == http.MethodPost {
			var input plantbuild.SiteControlConfirmation
			if err := decodeStrictRequestBody(r, s.cfg.BodyLimit, &input); err != nil {
				s.writeError(w, http.StatusBadRequest, "invalid_site_control_confirmation", err.Error(), false, "", "")
				return
			}
			if err := plantbuild.ValidateSiteControlConfirmation(input); err != nil {
				s.writeError(w, http.StatusUnprocessableEntity, "invalid_site_control_confirmation", err.Error(), false, "", "")
				return
			}
			tenantID := firstNonEmptyString(r.Header.Get("X-IAOS-Tenant-Id"), r.Header.Get("X-Tenant-ID"))
			authority, actorID, err := s.plantAuthorityClient(ctx, r, tenantID)
			if err != nil || authority == nil {
				s.writeError(w, http.StatusUnauthorized, "plant_site_control_identity_invalid", firstNonEmptyString(errorString(err), "IAOS authority is required"), false, "", "")
				return
			}
			rawControls, err := authority.PlantSiteControls(ctx, input.CaseCode)
			if err != nil {
				s.writePlantAuthorityError(w, err, "plant_site_controls_unavailable")
				return
			}
			var controls struct {
				Items []plantbuild.SiteControlItem `json:"items"`
			}
			if err := json.Unmarshal(rawControls, &controls); err != nil {
				s.writeError(w, http.StatusBadGateway, "plant_site_controls_invalid", "IAOS site control projection is invalid", true, "", "")
				return
			}
			var selected *plantbuild.SiteControlItem
			for index := range controls.Items {
				if controls.Items[index].Request.ControlRequestID == input.ControlRequestID {
					selected = &controls.Items[index]
					break
				}
			}
			if selected == nil || selected.Request.CaseCode != input.CaseCode {
				s.writeError(w, http.StatusNotFound, "site_control_request_not_found", "authoritative site control request was not found", false, "", "")
				return
			}
			if selected.Status == "controlled" && selected.Observation != nil {
				s.writeJSON(w, http.StatusOK, map[string]any{"status": "committed", "idempotent_replay": true, "world_message_id": "world-" + selected.Observation.ObservationID, "observation": selected.Observation})
				return
			}
			if selected.Status != "waiting_world" {
				s.writeError(w, http.StatusConflict, "site_control_request_not_waiting", "site control request is not waiting for World delivery", false, "", "")
				return
			}
			observation, err := plantbuild.GenerateSiteControlObservation(selected.Request)
			if err != nil {
				s.writeError(w, http.StatusUnprocessableEntity, "site_control_world_generation_failed", err.Error(), false, "", "")
				return
			}
			correlationID := "corr-m10-" + strings.TrimSpace(input.CaseCode)
			messageID := "world-" + observation.ObservationID
			envelope := map[string]any{
				"schema_version": "1.0", "message_id": messageID, "kind": "observation", "tenant_id": tenantID,
				"world_pack_key": "genesis-plant-delivery", "world_pack_version": "1.1.0", "world_run_id": selected.Request.WorldRunID,
				"branch_id": "main", "sim_occurred_at": observation.ObservedAt, "correlation_id": correlationID,
				"idempotency_key": messageID, "producer": map[string]string{"system": "aese", "component": "plant-site-control-world"},
				"subject_ref":  map[string]string{"type": "site_control_request", "code": observation.ControlRequestID},
				"payload_type": "site.control.delivered.v1", "payload": observation,
			}
			rawEnvelope, _ := json.Marshal(envelope)
			if _, err := authority.PostGovernedCommand(ctx, "api/v1/world-bridge/observations", rawEnvelope); err != nil {
				s.writePlantAuthorityError(w, err, "site_control_world_observation_not_accepted")
				return
			}
			commit := map[string]any{"schema_version": "1.0", "control_request_id": observation.ControlRequestID, "world_message_id": messageID}
			if err := postPlantAuthorityCommand(ctx, authority, "site.control.observation.commit", actorID, input.CaseCode, observation.ObservationID, 1, commit); err != nil {
				s.writePlantAuthorityError(w, err, "site_control_observation_not_committed")
				return
			}
			s.writeJSON(w, http.StatusCreated, map[string]any{"status": "committed", "idempotent_replay": false, "world_message_id": messageID, "observation": observation})
			return
		}
		if len(rest) == 4 && rest[1] == "plant-build" && rest[2] == "approvals" && r.Method == http.MethodGet {
			tenantID := firstNonEmptyString(r.Header.Get("X-IAOS-Tenant-Id"), r.Header.Get("X-Tenant-ID"))
			authority, _, err := s.plantAuthorityClient(ctx, r, tenantID)
			if err != nil || authority == nil {
				s.writeError(w, http.StatusUnauthorized, "plant_approval_identity_invalid", firstNonEmptyString(errorString(err), "IAOS authority is required"), false, "", "")
				return
			}
			result, err := authority.ApprovalDetail(ctx, rest[3])
			if err != nil {
				s.writePlantAuthorityError(w, err, "plant_approval_unavailable")
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(result)
			return
		}
		if len(rest) == 3 && rest[1] == "plant-build" && rest[2] == "site-selections" && r.Method == http.MethodPost {
			var input struct {
				SchemaVersion               string         `json:"schema_version"`
				RecommendationID            string         `json:"recommendation_id"`
				CaseCode                    string         `json:"case_code"`
				ProposalSetID               string         `json:"proposal_set_id"`
				ProposalSetRevision         int            `json:"proposal_set_revision"`
				SelectedProposalID          string         `json:"selected_proposal_id"`
				AssessmentPolicyVersion     string         `json:"assessment_policy_version"`
				Weights                     map[string]int `json:"weights"`
				RecommendationReason        string         `json:"recommendation_reason"`
				AlternativeComparison       string         `json:"alternative_comparison"`
				SingleSourceExceptionReason string         `json:"single_source_exception_reason,omitempty"`
				RecommendedAt               string         `json:"recommended_at"`
			}
			if err := decodeStrictRequestBody(r, s.cfg.BodyLimit, &input); err != nil {
				s.writeError(w, http.StatusBadRequest, "invalid_site_selection_recommendation", err.Error(), false, "", "")
				return
			}
			tenantID := firstNonEmptyString(r.Header.Get("X-IAOS-Tenant-Id"), r.Header.Get("X-Tenant-ID"))
			authority, actorID, err := s.plantAuthorityClient(ctx, r, tenantID)
			if err != nil || authority == nil {
				s.writeError(w, http.StatusUnauthorized, "plant_selection_identity_invalid", firstNonEmptyString(errorString(err), "IAOS authority is required"), false, "", "")
				return
			}
			result, err := postPlantAuthorityCommandResult(ctx, authority, "site.selection.recommend", actorID, input.CaseCode, input.RecommendationID, input.ProposalSetRevision, input)
			if err != nil {
				s.writePlantAuthorityError(w, err, "site_selection_recommendation_not_committed")
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write(result)
			return
		}
		if len(rest) == 4 && rest[1] == "plant-build" && rest[2] == "site-selections" && rest[3] == "finalize" && r.Method == http.MethodPost {
			var input struct {
				SchemaVersion     string `json:"schema_version"`
				RecommendationID  string `json:"recommendation_id"`
				ApprovalRequestID string `json:"approval_request_id"`
				FormalizedAt      string `json:"formalized_at"`
				CaseCode          string `json:"case_code"`
			}
			if err := decodeStrictRequestBody(r, s.cfg.BodyLimit, &input); err != nil {
				s.writeError(w, http.StatusBadRequest, "invalid_site_selection_formalization", err.Error(), false, "", "")
				return
			}
			tenantID := firstNonEmptyString(r.Header.Get("X-IAOS-Tenant-Id"), r.Header.Get("X-Tenant-ID"))
			authority, actorID, err := s.plantAuthorityClient(ctx, r, tenantID)
			if err != nil || authority == nil {
				s.writeError(w, http.StatusUnauthorized, "plant_selection_identity_invalid", firstNonEmptyString(errorString(err), "IAOS authority is required"), false, "", "")
				return
			}
			payload := map[string]any{"schema_version": input.SchemaVersion, "recommendation_id": input.RecommendationID, "approval_request_id": input.ApprovalRequestID, "formalized_at": input.FormalizedAt}
			result, err := postPlantAuthorityCommandResult(ctx, authority, "site.selection.formalize", actorID, input.CaseCode, input.RecommendationID, 1, payload)
			if err != nil {
				s.writePlantAuthorityError(w, err, "site_selection_not_formalized")
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write(result)
			return
		}
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

func (s *Server) handleIAOSCommandGateway(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
	rest []string,
) {
	if r.Method != http.MethodPost || len(rest) < 3 || rest[1] != "iaos" {
		s.writeError(w, http.StatusNotFound, "command_route_not_found", "command route not found", false, "", "")
		return
	}
	pathParts := rest[2:]
	if !allowedIAOSCommandPath(pathParts) {
		s.writeError(w, http.StatusForbidden, "command_route_not_allowed", "IAOS command route is not allow-listed", false, "", "")
		return
	}
	token := bearerToken(r.Header.Get("Authorization"))
	tenantID := strings.TrimSpace(r.Header.Get("X-IAOS-Tenant-Id"))
	if tenantID == "" {
		tenantID = strings.TrimSpace(r.Header.Get("X-Tenant-ID"))
	}
	if token == "" || tenantID == "" {
		s.writeError(w, http.StatusUnauthorized, "command_identity_required", "IAOS bearer token and tenant context are required", false, "", "")
		return
	}
	var body json.RawMessage
	if err := decodeStrictRequestBody(r, s.cfg.BodyLimit, &body); err != nil {
		s.writeError(w, http.StatusBadRequest, "invalid_command_body", err.Error(), false, "", "")
		return
	}
	client, err := iaosclient.New(iaosclient.Config{
		BaseURL: s.cfg.IAOSBaseURL, Token: token, TenantID: tenantID,
	})
	if err != nil {
		s.writeError(w, http.StatusServiceUnavailable, "command_gateway_unavailable", err.Error(), true, "", "")
		return
	}
	result, err := client.PostGovernedCommand(
		ctx, "api/v1/"+strings.Join(pathParts, "/"), body,
	)
	if err != nil {
		var apiErr *iaosclient.APIError
		if errors.As(err, &apiErr) {
			s.writeError(w, apiErr.StatusCode, firstNonEmptyString(apiErr.ErrorCode, "iaos_command_rejected"), apiErr.Message, apiErr.StatusCode >= 500, "", apiErr.RequiredPermission)
			return
		}
		s.writeError(w, http.StatusBadGateway, "iaos_command_failed", err.Error(), true, "", "")
		return
	}
	if len(result) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(result)
}

func allowedIAOSCommandPath(parts []string) bool {
	if len(parts) == 2 && parts[0] == "incorporations" && parts[1] == "cases" {
		return true
	}
	if len(parts) == 5 && parts[0] == "incorporations" &&
		parts[2] == "work-items" &&
		(parts[4] == "execute" || parts[4] == "dispatch-agent") {
		_, err := strconv.Atoi(parts[3])
		return parts[1] != "" && err == nil
	}
	if len(parts) == 5 && parts[0] == "incorporations" &&
		parts[2] == "gates" && parts[4] == "submit" {
		return parts[1] != "" && parts[3] != ""
	}
	if len(parts) == 3 && parts[0] == "approvals" && (parts[2] == "approve" || parts[2] == "reject") {
		return parts[1] != ""
	}
	return len(parts) == 2 && parts[0] == "world-bridge" && parts[1] == "observations"
}

func firstNonEmptyString(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func errorString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func (s *Server) plantAuthorityClient(ctx context.Context, r *http.Request, tenantID string) (*iaosclient.Client, string, error) {
	if strings.TrimSpace(s.cfg.IAOSBaseURL) == "" {
		if s.cfg.AllowLocalPlantAuthority {
			return nil, "local-test-actor", nil
		}
		return nil, "", fmt.Errorf("IAOS plant authority is not configured")
	}
	token := bearerToken(r.Header.Get("Authorization"))
	tenantID = strings.TrimSpace(tenantID)
	if token == "" || tenantID == "" {
		return nil, "", fmt.Errorf("IAOS bearer token and tenant context are required")
	}
	client, err := iaosclient.New(iaosclient.Config{BaseURL: s.cfg.IAOSBaseURL, Token: token, TenantID: tenantID})
	if err != nil {
		return nil, "", err
	}
	profile, err := client.Profile(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("load IAOS actor identity: %w", err)
	}
	if profile.TenantID != tenantID || strings.TrimSpace(profile.Username) == "" {
		return nil, "", fmt.Errorf("IAOS actor identity does not match tenant")
	}
	return client, profile.Username, nil
}

func postPlantAuthorityCommand(ctx context.Context, client *iaosclient.Client, capabilityCode, actorID, scopeID, objectID string, revision int, payload any) error {
	_, err := postPlantAuthorityCommandResult(ctx, client, capabilityCode, actorID, scopeID, objectID, revision, payload)
	return err
}

func postPlantAuthorityProposalCommand(ctx context.Context, client *iaosclient.Client, actorID, caseCode string, set plantbuild.ProposalSet, job creative.CreativeJob) error {
	run := plantbuild.AgentRunEvidence{
		AgentRunID: job.JobID, CaseCode: caseCode, AgentID: "plant-planning-agent", Status: job.Status,
		Provider: job.Provider, Model: job.Model, ModelVersion: job.ModelVersion,
		PromptVersion: job.PromptVersion, RequestID: job.RequestID, InputHash: job.InputHash,
		OutputHash: job.ContentHash, TokenUsage: job.TokenUsage, ValidationResult: job.ValidationResult,
		LatencyMS: job.LatencyMS, StartedAt: job.CreatedAt, CompletedAt: job.CompletedAt,
	}
	_, err := postPlantAuthorityCommandResultWithAgentRun(ctx, client, "site.proposal.record", actorID, caseCode, set.ProposalSetID, set.Revision, set, &run)
	return err
}

func postPlantAuthorityCommandResult(ctx context.Context, client *iaosclient.Client, capabilityCode, actorID, scopeID, objectID string, revision int, payload any) (json.RawMessage, error) {
	return postPlantAuthorityCommandResultWithAgentRun(ctx, client, capabilityCode, actorID, scopeID, objectID, revision, payload, nil)
}

func postPlantAuthorityCommandResultWithAgentRun(ctx context.Context, client *iaosclient.Client, capabilityCode, actorID, scopeID, objectID string, revision int, payload any, agentRun *plantbuild.AgentRunEvidence) (json.RawMessage, error) {
	field := ""
	switch capabilityCode {
	case "facility.requirement.define":
		field = "facility_requirement"
	case "site.proposal.record":
		field = "proposal_set"
	case "site.proposal.review":
		field = "proposal_review"
	case "site.investigation.request":
		field = "investigation_request"
	case "site.investigation.observation.commit":
		field = "investigation_observation"
	case "site.selection.recommend":
		field = "site_selection_recommendation"
	case "site.selection.formalize":
		field = "site_selection_formalization"
	case "site.control.request":
		field = "site_control_request"
	case "site.control.observation.commit":
		field = "site_control_observation"
	default:
		return nil, fmt.Errorf("unsupported M10 capability %q", capabilityCode)
	}
	command := map[string]any{
		"capability_code": capabilityCode,
		"actor_id":        actorID,
		"correlation_id":  "corr-m10-" + strings.TrimSpace(scopeID),
		"idempotency_key": fmt.Sprintf("m10:%s:%s:r%d", capabilityCode, strings.TrimSpace(objectID), revision),
		field:             payload,
	}
	if agentRun != nil {
		command["agent_run"] = agentRun
	}
	raw, err := json.Marshal(command)
	if err != nil {
		return nil, err
	}
	return client.PostGovernedCommand(ctx, "api/v1/genesis/plant/interactive/actions", raw)
}

func decodeStoredProposalSet(value any) (plantbuild.ProposalSet, bool) {
	var set plantbuild.ProposalSet
	raw, err := json.Marshal(value)
	if err != nil || json.Unmarshal(raw, &set) != nil || set.ProposalSetID == "" {
		return set, false
	}
	return set, true
}

func (s *Server) writePlantAuthorityError(w http.ResponseWriter, err error, fallbackCode string) {
	var apiErr *iaosclient.APIError
	if errors.As(err, &apiErr) {
		s.writeError(w, apiErr.StatusCode, firstNonEmptyString(apiErr.ErrorCode, fallbackCode), apiErr.Message, apiErr.StatusCode >= 500, "", apiErr.RequiredPermission)
		return
	}
	s.writeError(w, http.StatusBadGateway, fallbackCode, err.Error(), true, "", "")
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

func (s *Server) writeGenesisWorkspaceError(w http.ResponseWriter, err error, fallbackCode string) {
	var upstream *genesisworkspace.ControlPlaneHTTPError
	if !errors.As(err, &upstream) {
		s.writeError(w, http.StatusBadGateway, fallbackCode, err.Error(), true, "", "")
		return
	}
	switch upstream.StatusCode {
	case http.StatusUnauthorized:
		s.writeError(w, http.StatusUnauthorized, "player_session_expired",
			"IAOS player session expired; sign in to IAOS again and reopen Enterprise Genesis",
			false, "", "")
	case http.StatusForbidden:
		s.writeError(w, http.StatusForbidden, "workspace_access_denied",
			"current IAOS account cannot access this Genesis workspace",
			false, "", "")
	default:
		s.writeError(w, http.StatusBadGateway, fallbackCode, err.Error(), upstream.StatusCode >= 500, "", "")
	}
}

func (s *Server) writeGenesisAuthError(w http.ResponseWriter, err error, fallbackCode string) {
	var upstream *genesisworkspace.PlayerAuthHTTPError
	if !errors.As(err, &upstream) {
		s.writeError(w, http.StatusBadGateway, fallbackCode, "IAOS authentication service is unavailable", true, "", "")
		return
	}
	payload := struct {
		Error string `json:"error"`
		Code  string `json:"code"`
	}{}
	_ = json.Unmarshal([]byte(upstream.Body), &payload)
	if strings.TrimSpace(payload.Error) == "" {
		payload.Error = "IAOS authentication request failed"
	}
	if strings.TrimSpace(payload.Code) == "" {
		payload.Code = fallbackCode
	}
	status := upstream.StatusCode
	if status < 400 || status > 599 {
		status = http.StatusBadGateway
	}
	s.writeError(w, status, payload.Code, payload.Error, status >= 500, "", "")
}

func bearerToken(authorization string) string {
	fields := strings.Fields(strings.TrimSpace(authorization))
	if len(fields) != 2 || !strings.EqualFold(fields[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(fields[1])
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

func decodeStrictRequestBody(r *http.Request, limit int64, dst any) error {
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
	decoder := json.NewDecoder(strings.NewReader(string(data)))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dst); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}
	if decoder.Decode(&struct{}{}) != io.EOF {
		return fmt.Errorf("invalid JSON: multiple values")
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

func (s *Server) validateCreativeSession(ctx context.Context, r *http.Request, tenantID string) error {
	tenantID = strings.TrimSpace(tenantID)
	if tenantID == "" {
		return fmt.Errorf("creative request tenant is required")
	}
	if strings.TrimSpace(s.cfg.IAOSBaseURL) == "" {
		return nil
	}
	token, err := extractBearerToken(r)
	if err != nil || token == "" {
		return fmt.Errorf("authenticated workspace session required")
	}
	client, err := iaosclient.New(iaosclient.Config{
		BaseURL: s.cfg.IAOSBaseURL, Token: token, TenantID: tenantID,
	})
	if err != nil {
		return fmt.Errorf("invalid IAOS session configuration")
	}
	profile, err := client.Profile(ctx)
	if err != nil {
		return fmt.Errorf("IAOS workspace session validation failed")
	}
	if profile.TenantID != tenantID {
		return fmt.Errorf("workspace session tenant does not match creative request")
	}
	return nil
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
