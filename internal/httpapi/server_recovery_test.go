package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/industrial-ai/iaos-aese/internal/application"
	"github.com/industrial-ai/iaos-aese/internal/iaosclient"
	"github.com/industrial-ai/iaos-aese/internal/scenariopack"
)

func TestRefreshRunFromFactsRecoversIncrementallyOnRestart(t *testing.T) {
	pack, err := loadPack(filepath.Join("..", "..", "scenario-packs", "hctm"))
	if err != nil {
		t.Fatalf("load pack: %v", err)
	}
	story, err := application.FindStory(pack, "order-expedite-01")
	if err != nil {
		t.Fatalf("find story: %v", err)
	}
	plan, err := application.CompilePlan(pack, "order-expedite-01")
	if err != nil {
		t.Fatalf("compile plan: %v", err)
	}

	baseAt := func(event scenariopack.Event, cursor int64) map[string]any {
		return map[string]any{
			"event_id":             event.EventID,
			"event_type":           event.EventType,
			"cursor":               cursor,
			"occurred_at":          "2026-07-20T10:00:00+08:00",
			"correlation_id":       plan.Correlation,
			"business_object_type": "sales_order",
			"business_object_code": "SO-202607-0001",
			"payload":              map[string]any{},
		}
	}

	asObserved := func(event scenariopack.Event, cursor int64) iaosclient.ScenarioObservedEvent {
		return iaosclient.ScenarioObservedEvent{
			Cursor:        cursor,
			EventID:       event.EventID,
			EventType:     event.EventType,
			OccurredAt:    time.Date(2026, 7, 20, 10, 0, 0, 0, time.FixedZone("CST", 8*3600)),
			CorrelationID: plan.Correlation,
			BusinessType:  "sales_order",
			BusinessCode:  "SO-202607-0001",
			Payload:       json.RawMessage(`{}`),
		}
	}

	snapshotPhases := [][]map[string]any{
		{baseAt(story.Events.Events[0], 1), baseAt(story.Events.Events[1], 2), baseAt(story.Events.Events[2], 3)},
		{baseAt(story.Events.Events[0], 1), baseAt(story.Events.Events[1], 2), baseAt(story.Events.Events[2], 3),
			baseAt(story.Events.Events[3], 4), baseAt(story.Events.Events[4], 5), baseAt(story.Events.Events[5], 6), baseAt(story.Events.Events[6], 7)},
	}
	streamPhases := [][]iaosclient.ScenarioObservedEvent{
		{asObserved(story.Events.Events[0], 1), asObserved(story.Events.Events[1], 2), asObserved(story.Events.Events[2], 3)},
		{asObserved(story.Events.Events[3], 4), asObserved(story.Events.Events[4], 5), asObserved(story.Events.Events[5], 6), asObserved(story.Events.Events[6], 7)},
	}

	var phase int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/api/v1/scenarios/hctm/order-expedite-01/snapshot") &&
			!strings.HasSuffix(r.URL.Path, "/api/v1/scenarios/hctm/order-expedite-01/events") {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if strings.HasSuffix(r.URL.Path, "/events") {
			afterRaw := r.URL.Query().Get("after")
			after, err := strconv.ParseInt(afterRaw, 10, 64)
			if err != nil {
				after = 0
			}
			eventIndex := 0
			if after >= 3 {
				eventIndex = 1
			}
			events := make([]iaosclient.ScenarioObservedEvent, 0)
			for _, item := range streamPhases[eventIndex] {
				if item.Cursor > after {
					events = append(events, item)
				}
			}
			_ = json.NewEncoder(w).Encode(iaosclient.ScenarioEventsResponse{Items: events, HasMore: false})
			return
		}
		if phase >= len(snapshotPhases) {
			phase = len(snapshotPhases) - 1
		}
		_ = json.NewEncoder(w).Encode(iaosclient.ScenarioSnapshot{PackKey: "hctm", ScenarioKey: "order-expedite-01", CorrelationID: plan.Correlation, Cursor: 0, Events: snapshotPhases[phase]})
		phase++
	}))
	defer server.Close()

	client, err := iaosclient.New(iaosclient.Config{BaseURL: server.URL, Token: "token", TenantID: pack.Manifest.TenantTemplate})
	if err != nil {
		t.Fatalf("iaos client: %v", err)
	}

	run := &runRecord{
		RunID:       "run-restart-1",
		PackKey:     pack.Manifest.PackKey,
		PackVersion: pack.Manifest.PackVersion,
		PackDir:     filepath.Join("..", "..", "scenario-packs", "hctm"),
		ScenarioKey: story.Ref.Key,
		Plan:        plan,
		Status:      application.RunStatusRunning,
		CurrentAct:  0,
		TenantID:    pack.Manifest.TenantTemplate,
		Target:      server.URL,
		Token:       "token",
		Actor:       "aese-user",
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		Cursor:      0,
		ActionCache: map[string]actionCache{},
	}

	s := New(Config{PackDir: filepath.Join("..", "..", "scenario-packs", "hctm")})
	if err := s.refreshRunFromFacts(context.Background(), run, client); err != nil {
		t.Fatalf("first refresh: %v", err)
	}
	if run.CurrentAct != 1 {
		t.Fatalf("current act after first refresh = %d, want 1", run.CurrentAct)
	}
	if run.Cursor != 3 {
		t.Fatalf("cursor after first refresh = %d, want 3", run.Cursor)
	}

	recoveredRun := *run
	if err := s.refreshRunFromFacts(context.Background(), &recoveredRun, client); err != nil {
		t.Fatalf("second refresh: %v", err)
	}
	if recoveredRun.CurrentAct != 2 {
		t.Fatalf("current act after second refresh = %d, want 2", recoveredRun.CurrentAct)
	}
	if recoveredRun.Cursor != 7 {
		t.Fatalf("cursor after second refresh = %d, want 7", recoveredRun.Cursor)
	}
}

func TestRefreshRunFromFactsFiltersMismatchedCorrelation(t *testing.T) {
	packDir := filepath.Join("..", "..", "scenario-packs", "hctm")
	pack, err := loadPack(packDir)
	if err != nil {
		t.Fatalf("load pack: %v", err)
	}
	story, err := application.FindStory(pack, "order-expedite-01")
	if err != nil {
		t.Fatalf("find story: %v", err)
	}
	plan, err := application.CompilePlan(pack, "order-expedite-01")
	if err != nil {
		t.Fatalf("compile plan: %v", err)
	}

	matchStageOne := []map[string]any{
		{"event_id": story.Events.Events[0].EventID, "event_type": story.Events.Events[0].EventType, "cursor": 1, "correlation_id": plan.Correlation},
		{"event_id": story.Events.Events[1].EventID, "event_type": story.Events.Events[1].EventType, "cursor": 2, "correlation_id": plan.Correlation},
		{"event_id": story.Events.Events[2].EventID, "event_type": story.Events.Events[2].EventType, "cursor": 3, "correlation_id": plan.Correlation},
	}
	mismatchedEvents := []iaosclient.ScenarioObservedEvent{
		{Cursor: 4, EventID: story.Events.Events[3].EventID, EventType: story.Events.Events[3].EventType, CorrelationID: "corr-cross-tenant", BusinessType: "sales_order", BusinessCode: "SO-202607-0001", Payload: json.RawMessage(`{}`), OccurredAt: time.Now().UTC()},
		{Cursor: 5, EventID: story.Events.Events[4].EventID, EventType: story.Events.Events[4].EventType, CorrelationID: "corr-cross-tenant", BusinessType: "sales_order", BusinessCode: "SO-202607-0001", Payload: json.RawMessage(`{}`), OccurredAt: time.Now().UTC()},
		{Cursor: 6, EventID: story.Events.Events[5].EventID, EventType: story.Events.Events[5].EventType, CorrelationID: "corr-cross-tenant", BusinessType: "sales_order", BusinessCode: "SO-202607-0001", Payload: json.RawMessage(`{}`), OccurredAt: time.Now().UTC()},
		{Cursor: 7, EventID: story.Events.Events[6].EventID, EventType: story.Events.Events[6].EventType, CorrelationID: "corr-cross-tenant", BusinessType: "sales_order", BusinessCode: "SO-202607-0001", Payload: json.RawMessage(`{}`), OccurredAt: time.Now().UTC()},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/api/v1/profile"):
			_ = json.NewEncoder(w).Encode(iaosclient.ProfileResponse{Username: "aese-user", TenantID: pack.Manifest.TenantTemplate})
		case strings.HasSuffix(r.URL.Path, "/api/v1/scenarios/"+pack.Manifest.PackKey+"/"+story.Ref.Key+"/snapshot"):
			_ = json.NewEncoder(w).Encode(iaosclient.ScenarioSnapshot{
				PackKey:       pack.Manifest.PackKey,
				ScenarioKey:   story.Ref.Key,
				CorrelationID: plan.Correlation,
				Cursor:        0,
				Events:        matchStageOne,
			})
		case strings.HasSuffix(r.URL.Path, "/api/v1/scenarios/"+pack.Manifest.PackKey+"/"+story.Ref.Key+"/events"):
			_ = json.NewEncoder(w).Encode(iaosclient.ScenarioEventsResponse{Items: mismatchedEvents, HasMore: false, NextCursor: 7})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	client, err := iaosclient.New(iaosclient.Config{BaseURL: server.URL, Token: "token", TenantID: pack.Manifest.TenantTemplate})
	if err != nil {
		t.Fatalf("iaos client: %v", err)
	}

	run := &runRecord{
		RunID:       "run-correlation-filter",
		PackKey:     pack.Manifest.PackKey,
		PackVersion: pack.Manifest.PackVersion,
		PackDir:     packDir,
		ScenarioKey: story.Ref.Key,
		Plan:        plan,
		Status:      application.RunStatusRunning,
		CurrentAct:  0,
		TenantID:    pack.Manifest.TenantTemplate,
		Target:      server.URL,
		Token:       "token",
		Actor:       "aese-user",
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		Cursor:      0,
		ActionCache: map[string]actionCache{},
	}

	s := New(Config{PackDir: packDir})
	if err := s.refreshRunFromFacts(context.Background(), run, client); err != nil {
		t.Fatalf("refresh run from facts: %v", err)
	}
	if run.CurrentAct != 1 {
		t.Fatalf("current act = %d, want 1", run.CurrentAct)
	}
	if run.Cursor != 3 {
		t.Fatalf("cursor = %d, want 3", run.Cursor)
	}
	if run.Status != application.RunStatusRunning {
		t.Fatalf("status = %s, want %s", run.Status, application.RunStatusRunning)
	}
}

func TestRunCreateRejectsWithoutSnapshotPermission(t *testing.T) {
	packDir := filepath.Join("..", "..", "scenario-packs", "hctm")
	pack, err := loadPack(packDir)
	if err != nil {
		t.Fatalf("load pack: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/api/v1/profile"):
			_ = json.NewEncoder(w).Encode(iaosclient.ProfileResponse{Username: "aese-user", TenantID: pack.Manifest.TenantTemplate})
		case strings.HasSuffix(r.URL.Path, "/api/v1/scenarios/"+pack.Manifest.PackKey+"/"+pack.Stories[0].Ref.Key+"/snapshot"):
			http.Error(w, `{"error":"permission denied"}`, http.StatusForbidden)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	payload := runCreateRequest{
		Target:   server.URL,
		Tenant:   pack.Manifest.TenantTemplate,
		StoryKey: pack.Stories[0].Ref.Key,
		Token:    "token",
		PackDir:  packDir,
		RunID:    "run-t11-create-1",
	}
	requestBody, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	api := httptest.NewServer(New(Config{PackDir: packDir}))
	defer api.Close()

	req, err := http.NewRequest(http.MethodPost, api.URL+"/api/aese/v1/runs", bytes.NewReader(requestBody))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	api.Config.Handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("create status = %d, want %d", recorder.Code, http.StatusForbidden)
	}

	var response errorResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Code != "forbidden" {
		t.Fatalf("error code = %q, want forbidden", response.Code)
	}
	if response.RequiredPermission != "scenario.run.read" {
		t.Fatalf("required_permission = %q, want scenario.run.read", response.RequiredPermission)
	}
}

func TestRunStatusRejectsWhenSnapshotPermissionDenied(t *testing.T) {
	packDir := filepath.Join("..", "..", "scenario-packs", "hctm")
	pack, err := loadPack(packDir)
	if err != nil {
		t.Fatalf("load pack: %v", err)
	}
	plan, err := application.CompilePlan(pack, pack.Stories[0].Ref.Key)
	if err != nil {
		t.Fatalf("compile plan: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/api/v1/profile"):
			_ = json.NewEncoder(w).Encode(iaosclient.ProfileResponse{Username: "aese-user", TenantID: pack.Manifest.TenantTemplate})
		case strings.HasSuffix(r.URL.Path, "/api/v1/scenarios/"+pack.Manifest.PackKey+"/"+pack.Stories[0].Ref.Key+"/snapshot"):
			http.Error(w, `{"error":"permission denied"}`, http.StatusForbidden)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	run := &runRecord{
		RunID:       "run-t11-status-1",
		PackKey:     pack.Manifest.PackKey,
		PackVersion: pack.Manifest.PackVersion,
		PackDir:     packDir,
		ScenarioKey: pack.Stories[0].Ref.Key,
		Plan:        plan,
		Status:      application.RunStatusPlanned,
		CurrentAct:  0,
		TenantID:    pack.Manifest.TenantTemplate,
		Target:      server.URL,
		Token:       "token",
		Actor:       "aese-user",
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		ActionCache: map[string]actionCache{},
	}

	api := New(Config{PackDir: packDir})
	api.runs[run.RunID] = run

	req, err := http.NewRequest(http.MethodGet, "http://example.com/api/aese/v1/runs/"+run.RunID, nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer token")
	recorder := httptest.NewRecorder()
	api.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("status code = %d, want %d", recorder.Code, http.StatusForbidden)
	}

	var response errorResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Code != "forbidden" {
		t.Fatalf("error code = %q, want forbidden", response.Code)
	}
	if response.RequiredPermission != "scenario.run.read" {
		t.Fatalf("required_permission = %q, want scenario.run.read", response.RequiredPermission)
	}
}

func TestRunCreateBlocksConcurrentWritableRunForSameTenantAndScenario(t *testing.T) {
	packDir := filepath.Join("..", "..", "scenario-packs", "hctm")
	pack, err := loadPack(packDir)
	if err != nil {
		t.Fatalf("load pack: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/api/v1/profile"):
			_ = json.NewEncoder(w).Encode(iaosclient.ProfileResponse{Username: "aese-user", TenantID: pack.Manifest.TenantTemplate})
		case strings.HasSuffix(r.URL.Path, "/api/v1/scenarios/"+pack.Manifest.PackKey+"/"+pack.Stories[0].Ref.Key+"/snapshot"):
			_ = json.NewEncoder(w).Encode(iaosclient.ScenarioSnapshot{PackKey: pack.Manifest.PackKey, ScenarioKey: pack.Stories[0].Ref.Key, CorrelationID: "corr", Cursor: 0})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	api := httptest.NewServer(New(Config{PackDir: packDir}))
	defer api.Close()

	payload := func(runID string) runCreateRequest {
		return runCreateRequest{
			Target:   server.URL,
			Tenant:   pack.Manifest.TenantTemplate,
			StoryKey: pack.Stories[0].Ref.Key,
			RunID:    runID,
			Token:    "token",
			PackDir:  packDir,
		}
	}

	firstBody, err := json.Marshal(payload("run-conflict-1"))
	if err != nil {
		t.Fatalf("marshal first payload: %v", err)
	}
	secondBody, err := json.Marshal(payload("run-conflict-2"))
	if err != nil {
		t.Fatalf("marshal second payload: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, api.URL+"/api/aese/v1/runs", bytes.NewReader(firstBody))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	firstRecorder := httptest.NewRecorder()
	api.Config.Handler.ServeHTTP(firstRecorder, req)
	if firstRecorder.Code != http.StatusCreated {
		t.Fatalf("create first status = %d, want %d", firstRecorder.Code, http.StatusCreated)
	}

	req, err = http.NewRequest(http.MethodPost, api.URL+"/api/aese/v1/runs", bytes.NewReader(secondBody))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	secondRecorder := httptest.NewRecorder()
	api.Config.Handler.ServeHTTP(secondRecorder, req)
	if secondRecorder.Code != http.StatusConflict {
		t.Fatalf("create second status = %d, want %d", secondRecorder.Code, http.StatusConflict)
	}

	var response errorResponse
	if err := json.Unmarshal(secondRecorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Code != "conflict" {
		t.Fatalf("error code = %q, want conflict", response.Code)
	}
}

func TestRunCreateAllowsAfterCompletedWritableRun(t *testing.T) {
	packDir := filepath.Join("..", "..", "scenario-packs", "hctm")
	pack, err := loadPack(packDir)
	if err != nil {
		t.Fatalf("load pack: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/api/v1/profile"):
			_ = json.NewEncoder(w).Encode(iaosclient.ProfileResponse{Username: "aese-user", TenantID: pack.Manifest.TenantTemplate})
		case strings.HasSuffix(r.URL.Path, "/api/v1/scenarios/"+pack.Manifest.PackKey+"/"+pack.Stories[0].Ref.Key+"/snapshot"):
			_ = json.NewEncoder(w).Encode(iaosclient.ScenarioSnapshot{PackKey: pack.Manifest.PackKey, ScenarioKey: pack.Stories[0].Ref.Key, CorrelationID: "corr", Cursor: 0})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	s := New(Config{PackDir: packDir})
	s.runs["run-completed"] = &runRecord{
		RunID:       "run-completed",
		PackKey:     pack.Manifest.PackKey,
		PackVersion: pack.Manifest.PackVersion,
		PackDir:     packDir,
		ScenarioKey: pack.Stories[0].Ref.Key,
		Plan:        application.Plan{},
		Status:      application.RunStatusCompleted,
		CurrentAct:  0,
		TenantID:    pack.Manifest.TenantTemplate,
	}

	api := httptest.NewServer(s)
	defer api.Close()

	request := runCreateRequest{
		Target:   server.URL,
		Tenant:   pack.Manifest.TenantTemplate,
		StoryKey: pack.Stories[0].Ref.Key,
		RunID:    "run-new-after-complete",
		Token:    "token",
		PackDir:  packDir,
	}
	requestBody, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, api.URL+"/api/aese/v1/runs", bytes.NewReader(requestBody))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	s.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want %d", recorder.Code, http.StatusCreated)
	}
	var response struct {
		RunID string `json:"run_id"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.RunID != request.RunID {
		t.Fatalf("run id = %q, want %q", response.RunID, request.RunID)
	}
}

func TestRunActionBlocksStaleCursor(t *testing.T) {
	packDir := filepath.Join("..", "..", "scenario-packs", "hctm")
	pack, err := loadPack(packDir)
	if err != nil {
		t.Fatalf("load pack: %v", err)
	}
	plan, err := application.CompilePlan(pack, pack.Stories[0].Ref.Key)
	if err != nil {
		t.Fatalf("compile plan: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/api/v1/profile"):
			_ = json.NewEncoder(w).Encode(iaosclient.ProfileResponse{Username: "aese-user", TenantID: pack.Manifest.TenantTemplate})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	api := New(Config{PackDir: packDir})
	api.runs["run-stale-cursor"] = &runRecord{
		RunID:       "run-stale-cursor",
		PackKey:     pack.Manifest.PackKey,
		PackVersion: pack.Manifest.PackVersion,
		PackDir:     packDir,
		ScenarioKey: pack.Stories[0].Ref.Key,
		Plan:        plan,
		Status:      application.RunStatusPlanned,
		CurrentAct:  0,
		TenantID:    pack.Manifest.TenantTemplate,
		Target:      server.URL,
		Token:       "token",
		Actor:       "aese-user",
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		Cursor:      3,
		ActionCache: map[string]actionCache{},
	}

	body, err := json.Marshal(map[string]any{"expected_cursor": 2})
	if err != nil {
		t.Fatalf("marshal body: %v", err)
	}
	req, err := http.NewRequest(http.MethodPost, "http://example.com/api/aese/v1/runs/run-stale-cursor/initialize", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer token")
	recorder := httptest.NewRecorder()
	api.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusConflict {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusConflict)
	}
	var response errorResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Code != "cursor_mismatch" {
		t.Fatalf("error code = %q, want cursor_mismatch", response.Code)
	}
}

func TestRunStatusRejectsTenantMismatch(t *testing.T) {
	packDir := filepath.Join("..", "..", "scenario-packs", "hctm")
	pack, err := loadPack(packDir)
	if err != nil {
		t.Fatalf("load pack: %v", err)
	}
	plan, err := application.CompilePlan(pack, pack.Stories[0].Ref.Key)
	if err != nil {
		t.Fatalf("compile plan: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/api/v1/profile") {
			_ = json.NewEncoder(w).Encode(iaosclient.ProfileResponse{Username: "aese-user", TenantID: "tenant-other"})
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	run := &runRecord{
		RunID:       "run-tenant-mismatch",
		PackKey:     pack.Manifest.PackKey,
		PackVersion: pack.Manifest.PackVersion,
		PackDir:     packDir,
		ScenarioKey: pack.Stories[0].Ref.Key,
		Plan:        plan,
		Status:      application.RunStatusPlanned,
		CurrentAct:  0,
		TenantID:    pack.Manifest.TenantTemplate,
		Target:      server.URL,
		Token:       "token",
		Actor:       "aese-user",
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		ActionCache: map[string]actionCache{},
	}

	api := New(Config{PackDir: packDir})
	api.runs[run.RunID] = run

	req, err := http.NewRequest(http.MethodGet, "http://example.com/api/aese/v1/runs/"+run.RunID, nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer token")
	recorder := httptest.NewRecorder()
	api.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("status code = %d, want %d", recorder.Code, http.StatusForbidden)
	}

	var response errorResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Code != "tenant_mismatch" {
		t.Fatalf("error code = %q, want tenant_mismatch", response.Code)
	}
}

func TestRunActionAdvanceIsIdempotentUnderConcurrentCalls(t *testing.T) {
	packDir := filepath.Join("..", "..", "scenario-packs", "hctm")
	pack, err := loadPack(packDir)
	if err != nil {
		t.Fatalf("load pack: %v", err)
	}
	plan, err := application.CompilePlan(pack, pack.Stories[0].Ref.Key)
	if err != nil {
		t.Fatalf("compile plan: %v", err)
	}

	var decomposeCalls int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/api/v1/profile"):
			_ = json.NewEncoder(w).Encode(iaosclient.ProfileResponse{Username: "aese-user", TenantID: pack.Manifest.TenantTemplate})
		case strings.HasSuffix(r.URL.Path, "/api/v1/scenarios/"+pack.Manifest.PackKey+"/"+pack.Stories[0].Ref.Key+"/snapshot"):
			_ = json.NewEncoder(w).Encode(iaosclient.ScenarioSnapshot{
				PackKey:       pack.Manifest.PackKey,
				ScenarioKey:   pack.Stories[0].Ref.Key,
				CorrelationID: plan.Correlation,
				Cursor:        0,
			})
		case strings.HasSuffix(r.URL.Path, "/api/v1/scenarios/"+pack.Manifest.PackKey+"/"+pack.Stories[0].Ref.Key+"/events"):
			_ = json.NewEncoder(w).Encode(iaosclient.ScenarioEventsResponse{Items: []iaosclient.ScenarioObservedEvent{}, HasMore: false, NextCursor: 0})
		case strings.HasSuffix(r.URL.Path, "/api/v1/entities/sales_order/records"):
			_ = json.NewEncoder(w).Encode(map[string]any{
				"total": 1,
				"data":  []map[string]any{{"id": "so-202607-0001", "order_no": "SO-202607-0001", "status": "new"}},
			})
		case strings.HasSuffix(r.URL.Path, "/api/v1/entities/sales_order/so-202607-0001/decompose"):
			atomic.AddInt64(&decomposeCalls, 1)
			_ = json.NewEncoder(w).Encode(iaosclient.DecomposeResult{
				Status:       "decompose_submitted",
				Decomposing:  true,
				SalesOrderID: "so-202607-0001",
				OrderNo:      "SO-202607-0001",
				Cursor:       10,
				OperationRef: "op-decompose-so-202607-0001",
				Committed:    true,
			})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	api := New(Config{PackDir: packDir})
	api.runs["run-concurrent-advance"] = &runRecord{
		RunID:       "run-concurrent-advance",
		PackKey:     pack.Manifest.PackKey,
		PackVersion: pack.Manifest.PackVersion,
		PackDir:     packDir,
		ScenarioKey: pack.Stories[0].Ref.Key,
		Plan:        plan,
		Status:      application.RunStatusReady,
		CurrentAct:  0,
		TenantID:    pack.Manifest.TenantTemplate,
		Target:      server.URL,
		Token:       "token",
		Actor:       "aese-user",
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		Cursor:      0,
		ActionCache: map[string]actionCache{},
	}

	payload := []byte(`{"apply":true,"idempotency_key":"idemp-advance-001"}`)
	type actionCallResult struct {
		status int
		body   []byte
		err    error
	}
	callAdvance := func(done chan<- actionCallResult) {
		req, reqErr := http.NewRequest(http.MethodPost, "http://example.com/api/aese/v1/runs/run-concurrent-advance/advance", bytes.NewReader(payload))
		if reqErr != nil {
			done <- actionCallResult{err: reqErr}
			return
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer token")
		recorder := httptest.NewRecorder()
		api.ServeHTTP(recorder, req)
		done <- actionCallResult{status: recorder.Code, body: append([]byte(nil), recorder.Body.Bytes()...)}
	}

	results := make(chan actionCallResult, 2)
	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			callAdvance(results)
		}()
	}
	wg.Wait()
	close(results)

	for result := range results {
		if result.err != nil {
			t.Fatalf("send request: %v", result.err)
		}
		if result.status != http.StatusOK {
			t.Fatalf("advance status = %d, want %d", result.status, http.StatusOK)
		}
		var response actionResponse
		if err := json.Unmarshal(result.body, &response); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		if response.Action != "advance" {
			t.Fatalf("action = %q, want advance", response.Action)
		}
		outcomeMap, ok := response.Run.Outcome.(map[string]any)
		if !ok {
			t.Fatal("outcome expected as object")
		}
		rawImpacts, ok := outcomeMap["impacts"].([]any)
		if !ok {
			t.Fatal("outcome.impacts expected as array")
		}
		foundEvidence := false
		for _, raw := range rawImpacts {
			impact, ok := raw.(map[string]any)
			if !ok {
				continue
			}
			cursor, hasCursor := impact["cursor"].(float64)
			operationRef, hasRef := impact["operation_ref"].(string)
			committed, hasCommitted := impact["committed"].(bool)
			correlationID, hasCorrelation := impact["correlation_id"].(string)
			_, hasNoOp := impact["no_op"]
			if hasCursor && hasRef && hasCommitted && hasCorrelation && hasNoOp && operationRef != "" && correlationID == plan.Correlation && cursor == 10 {
				if !committed {
					t.Fatal("replayed impact should report committed=true for decompose action")
				}
				foundEvidence = true
				break
			}
		}
		if !foundEvidence {
			t.Fatal("expected replay impact with cursor, operation_ref and committed")
		}
	}
	if got := atomic.LoadInt64(&decomposeCalls); got != 1 {
		t.Fatalf("decompose calls = %d, want 1", got)
	}
}

func TestRunActionResetConflictsAfterFirstReset(t *testing.T) {
	packDir := filepath.Join("..", "..", "scenario-packs", "hctm")
	pack, err := loadPack(packDir)
	if err != nil {
		t.Fatalf("load pack: %v", err)
	}
	plan, err := application.CompilePlan(pack, pack.Stories[0].Ref.Key)
	if err != nil {
		t.Fatalf("compile plan: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/api/v1/profile"):
			_ = json.NewEncoder(w).Encode(iaosclient.ProfileResponse{Username: "aese-user", TenantID: pack.Manifest.TenantTemplate})
		case strings.HasSuffix(r.URL.Path, "/api/v1/scenarios/"+pack.Manifest.PackKey+"/"+pack.Stories[0].Ref.Key+"/snapshot"):
			_ = json.NewEncoder(w).Encode(iaosclient.ScenarioSnapshot{
				PackKey:       pack.Manifest.PackKey,
				ScenarioKey:   pack.Stories[0].Ref.Key,
				CorrelationID: plan.Correlation,
				Cursor:        3,
			})
		case strings.HasSuffix(r.URL.Path, "/api/v1/scenarios/"+pack.Manifest.PackKey+"/"+pack.Stories[0].Ref.Key+"/events"):
			_ = json.NewEncoder(w).Encode(iaosclient.ScenarioEventsResponse{Items: []iaosclient.ScenarioObservedEvent{}, HasMore: false, NextCursor: 0})
		case strings.HasSuffix(r.URL.Path, "/api/v1/scenarios/reset"):
			_ = json.NewEncoder(w).Encode(iaosclient.ScenarioSummary{
				PackKey:       pack.Manifest.PackKey,
				PackVersion:   pack.Manifest.PackVersion,
				ScenarioKey:   pack.Stories[0].Ref.Key,
				RunID:         "run-reset-conflict",
				CorrelationID: plan.Correlation,
				TenantID:      pack.Manifest.TenantTemplate,
				NoOp:          1,
				Committed:     true,
			})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	api := New(Config{PackDir: packDir})
	api.runs["run-reset-conflict"] = &runRecord{
		RunID:       "run-reset-conflict",
		PackKey:     pack.Manifest.PackKey,
		PackVersion: pack.Manifest.PackVersion,
		PackDir:     packDir,
		ScenarioKey: pack.Stories[0].Ref.Key,
		Plan:        plan,
		Status:      application.RunStatusReady,
		CurrentAct:  0,
		TenantID:    pack.Manifest.TenantTemplate,
		Target:      server.URL,
		Token:       "token",
		Actor:       "aese-user",
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		Cursor:      0,
		ActionCache: map[string]actionCache{},
	}

	planBody := []byte(`{"plan_hash":"` + plan.PlanHash + `"}`)
	resetPlanReq, reqErr := http.NewRequest(http.MethodPost, "http://example.com/api/aese/v1/runs/run-reset-conflict/reset-plan", bytes.NewReader(planBody))
	if reqErr != nil {
		t.Fatalf("new request: %v", reqErr)
	}
	resetPlanReq.Header.Set("Content-Type", "application/json")
	resetPlanReq.Header.Set("Authorization", "Bearer token")
	resetPlanRecorder := httptest.NewRecorder()
	api.ServeHTTP(resetPlanRecorder, resetPlanReq)
	if resetPlanRecorder.Code != http.StatusOK {
		t.Fatalf("reset-plan status = %d, want %d", resetPlanRecorder.Code, http.StatusOK)
	}
	var planResp actionResponse
	if err := json.Unmarshal(resetPlanRecorder.Body.Bytes(), &planResp); err != nil {
		t.Fatalf("decode reset-plan response: %v", err)
	}
	planOutcome, ok := planResp.Run.Outcome.(map[string]any)
	if !ok {
		t.Fatal("reset-plan outcome expected as object")
	}
	resetToken, ok := planOutcome["reset_confirmation_token"].(string)
	if !ok || resetToken == "" {
		t.Fatal("reset confirmation token missing in response")
	}

	resetPayload, err := json.Marshal(map[string]any{
		"apply":              true,
		"idempotency_key":    "reset-1",
		"confirmation_token": resetToken,
		"plan_hash":          plan.PlanHash,
	})
	if err != nil {
		t.Fatalf("marshal reset payload: %v", err)
	}
	resetReq, err := http.NewRequest(http.MethodPost, "http://example.com/api/aese/v1/runs/run-reset-conflict/reset", bytes.NewReader(resetPayload))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	resetReq.Header.Set("Content-Type", "application/json")
	resetReq.Header.Set("Authorization", "Bearer token")
	resetRecorder := httptest.NewRecorder()
	api.ServeHTTP(resetRecorder, resetReq)
	if resetRecorder.Code != http.StatusOK {
		t.Fatalf("first reset status = %d, want %d", resetRecorder.Code, http.StatusOK)
	}

	conflictPayload, err := json.Marshal(map[string]any{
		"apply":              true,
		"idempotency_key":    "reset-2",
		"confirmation_token": resetToken,
		"plan_hash":          plan.PlanHash,
	})
	if err != nil {
		t.Fatalf("marshal conflict payload: %v", err)
	}
	conflictReq, err := http.NewRequest(http.MethodPost, "http://example.com/api/aese/v1/runs/run-reset-conflict/reset", bytes.NewReader(conflictPayload))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	conflictReq.Header.Set("Content-Type", "application/json")
	conflictReq.Header.Set("Authorization", "Bearer token")
	conflictRecorder := httptest.NewRecorder()
	api.ServeHTTP(conflictRecorder, conflictReq)
	if conflictRecorder.Code != http.StatusConflict {
		t.Fatalf("second reset status = %d, want %d", conflictRecorder.Code, http.StatusConflict)
	}

	var conflictResp errorResponse
	if err := json.Unmarshal(conflictRecorder.Body.Bytes(), &conflictResp); err != nil {
		t.Fatalf("decode conflict response: %v", err)
	}
	if conflictResp.Code != "invalid_state" {
		t.Fatalf("error code = %q, want invalid_state", conflictResp.Code)
	}
}

func TestRunActionRejectsTenantMismatchDuringAction(t *testing.T) {
	packDir := filepath.Join("..", "..", "scenario-packs", "hctm")
	pack, err := loadPack(packDir)
	if err != nil {
		t.Fatalf("load pack: %v", err)
	}
	plan, err := application.CompilePlan(pack, pack.Stories[0].Ref.Key)
	if err != nil {
		t.Fatalf("compile plan: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/api/v1/profile") {
			_ = json.NewEncoder(w).Encode(iaosclient.ProfileResponse{Username: "aese-user", TenantID: "tenant-other"})
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	api := New(Config{PackDir: packDir})
	api.runs["run-action-tenant-mismatch"] = &runRecord{
		RunID:       "run-action-tenant-mismatch",
		PackKey:     pack.Manifest.PackKey,
		PackVersion: pack.Manifest.PackVersion,
		PackDir:     packDir,
		ScenarioKey: pack.Stories[0].Ref.Key,
		Plan:        plan,
		Status:      application.RunStatusPlanned,
		CurrentAct:  0,
		TenantID:    pack.Manifest.TenantTemplate,
		Target:      server.URL,
		Token:       "token",
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		Cursor:      0,
		ActionCache: map[string]actionCache{},
	}

	req, err := http.NewRequest(http.MethodPost, "http://example.com/api/aese/v1/runs/run-action-tenant-mismatch/preflight", bytes.NewReader([]byte(`{}`)))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer token")
	recorder := httptest.NewRecorder()
	api.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusForbidden)
	}
	var response errorResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Code != "tenant_mismatch" {
		t.Fatalf("error code = %q, want tenant_mismatch", response.Code)
	}
}

func TestRunActionRejectsInsufficientPermission(t *testing.T) {
	packDir := filepath.Join("..", "..", "scenario-packs", "hctm")
	pack, err := loadPack(packDir)
	if err != nil {
		t.Fatalf("load pack: %v", err)
	}
	plan, err := application.CompilePlan(pack, pack.Stories[0].Ref.Key)
	if err != nil {
		t.Fatalf("compile plan: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/api/v1/profile"):
			_ = json.NewEncoder(w).Encode(iaosclient.ProfileResponse{Username: "aese-user", TenantID: pack.Manifest.TenantTemplate})
		case strings.HasSuffix(r.URL.Path, "/api/v1/scenarios/"+pack.Manifest.PackKey+"/"+pack.Stories[0].Ref.Key+"/snapshot"):
			_ = json.NewEncoder(w).Encode(iaosclient.ScenarioSnapshot{PackKey: pack.Manifest.PackKey, ScenarioKey: pack.Stories[0].Ref.Key, CorrelationID: plan.Correlation, Cursor: 0})
		case strings.HasSuffix(r.URL.Path, "/api/v1/scenarios/apply"):
			http.Error(w, `{"error":"permission denied"}`, http.StatusForbidden)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	api := New(Config{PackDir: packDir})
	api.runs["run-insufficient-permission"] = &runRecord{
		RunID:       "run-insufficient-permission",
		PackKey:     pack.Manifest.PackKey,
		PackVersion: pack.Manifest.PackVersion,
		PackDir:     packDir,
		ScenarioKey: pack.Stories[0].Ref.Key,
		Plan:        plan,
		Status:      application.RunStatusInitializing,
		CurrentAct:  0,
		TenantID:    pack.Manifest.TenantTemplate,
		Target:      server.URL,
		Token:       "token",
		Actor:       "aese-user",
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		Cursor:      0,
		ActionCache: map[string]actionCache{},
	}

	initializePayload, err := json.Marshal(map[string]any{
		"plan_hash":       plan.PlanHash,
		"apply":           true,
		"idempotency_key": "init-no-permission",
	})
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}
	req, err := http.NewRequest(http.MethodPost, "http://example.com/api/aese/v1/runs/run-insufficient-permission/initialize", bytes.NewReader(initializePayload))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer token")
	recorder := httptest.NewRecorder()
	api.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusForbidden)
	}
	var response errorResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Code != "forbidden" {
		t.Fatalf("error code = %q, want forbidden", response.Code)
	}
	if response.RequiredPermission != "scenario.run.execute" {
		t.Fatalf("required_permission = %q, want scenario.run.execute", response.RequiredPermission)
	}
}

func TestRunActionRejectsResetWithInsufficientPermission(t *testing.T) {
	packDir := filepath.Join("..", "..", "scenario-packs", "hctm")
	pack, err := loadPack(packDir)
	if err != nil {
		t.Fatalf("load pack: %v", err)
	}
	plan, err := application.CompilePlan(pack, pack.Stories[0].Ref.Key)
	if err != nil {
		t.Fatalf("compile plan: %v", err)
	}

	apiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/api/v1/profile"):
			_ = json.NewEncoder(w).Encode(iaosclient.ProfileResponse{Username: "aese-user", TenantID: pack.Manifest.TenantTemplate})
		case strings.HasSuffix(r.URL.Path, "/api/v1/scenarios/"+pack.Manifest.PackKey+"/"+pack.Stories[0].Ref.Key+"/snapshot"):
			_ = json.NewEncoder(w).Encode(iaosclient.ScenarioSnapshot{PackKey: pack.Manifest.PackKey, ScenarioKey: pack.Stories[0].Ref.Key, CorrelationID: plan.Correlation, Cursor: 0})
		case strings.HasSuffix(r.URL.Path, "/api/v1/scenarios/reset"):
			http.Error(w, `{"error":"permission denied"}`, http.StatusForbidden)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer apiServer.Close()

	api := New(Config{PackDir: packDir})
	api.runs["run-reset-no-permission"] = &runRecord{
		RunID:               "run-reset-no-permission",
		PackKey:             pack.Manifest.PackKey,
		PackVersion:         pack.Manifest.PackVersion,
		PackDir:             packDir,
		ScenarioKey:         pack.Stories[0].Ref.Key,
		Plan:                plan,
		Status:              application.RunStatusResetting,
		CurrentAct:          0,
		TenantID:            pack.Manifest.TenantTemplate,
		Target:              apiServer.URL,
		Token:               "token",
		Actor:               "aese-user",
		CreatedAt:           time.Now().UTC(),
		UpdatedAt:           time.Now().UTC(),
		Cursor:              0,
		ResetToken:          "reset-token-2026",
		ResetTokenExpiresAt: time.Now().UTC().Add(10 * time.Minute),
		ActionCache:         map[string]actionCache{},
	}

	payload, err := json.Marshal(map[string]any{
		"apply":              true,
		"idempotency_key":    "reset-no-permission",
		"confirmation_token": "reset-token-2026",
		"plan_hash":          plan.PlanHash,
	})
	if err != nil {
		t.Fatalf("marshal reset payload: %v", err)
	}
	req, err := http.NewRequest(http.MethodPost, "http://example.com/api/aese/v1/runs/run-reset-no-permission/reset", bytes.NewReader(payload))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer token")
	recorder := httptest.NewRecorder()
	api.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusForbidden)
	}
	var response errorResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Code != "forbidden" {
		t.Fatalf("error code = %q, want forbidden", response.Code)
	}
	if response.RequiredPermission != "scenario.run.reset" {
		t.Fatalf("required_permission = %q, want scenario.run.reset", response.RequiredPermission)
	}
}

func TestRunActionRejectsExpiredResetToken(t *testing.T) {
	packDir := filepath.Join("..", "..", "scenario-packs", "hctm")
	pack, err := loadPack(packDir)
	if err != nil {
		t.Fatalf("load pack: %v", err)
	}
	plan, err := application.CompilePlan(pack, pack.Stories[0].Ref.Key)
	if err != nil {
		t.Fatalf("compile plan: %v", err)
	}

	apiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/api/v1/profile"):
			_ = json.NewEncoder(w).Encode(iaosclient.ProfileResponse{Username: "aese-user", TenantID: pack.Manifest.TenantTemplate})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer apiServer.Close()

	api := New(Config{PackDir: packDir})
	api.runs["run-reset-expired-token"] = &runRecord{
		RunID:               "run-reset-expired-token",
		PackKey:             pack.Manifest.PackKey,
		PackVersion:         pack.Manifest.PackVersion,
		PackDir:             packDir,
		ScenarioKey:         pack.Stories[0].Ref.Key,
		Plan:                plan,
		Status:              application.RunStatusResetting,
		CurrentAct:          0,
		TenantID:            pack.Manifest.TenantTemplate,
		Target:              apiServer.URL,
		Token:               "token",
		Actor:               "aese-user",
		CreatedAt:           time.Now().UTC(),
		UpdatedAt:           time.Now().UTC(),
		Cursor:              0,
		ResetToken:          "reset-token-expired",
		ResetTokenExpiresAt:  time.Now().UTC().Add(-time.Minute),
		ActionCache:         map[string]actionCache{},
	}

	payload, err := json.Marshal(map[string]any{
		"apply":              true,
		"idempotency_key":    "reset-expired-token",
		"confirmation_token": "reset-token-expired",
	})
	if err != nil {
		t.Fatalf("marshal reset payload: %v", err)
	}
	req, err := http.NewRequest(http.MethodPost, "http://example.com/api/aese/v1/runs/run-reset-expired-token/reset", bytes.NewReader(payload))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer token")
	recorder := httptest.NewRecorder()
	api.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusForbidden)
	}

	var response errorResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Code != "reset_confirmation_invalid" {
		t.Fatalf("error code = %q, want reset_confirmation_invalid", response.Code)
	}
}

func TestRefreshRunFromFactsUsesRecommendationsToRecoverAwaitingVerification(t *testing.T) {
	packDir := filepath.Join("..", "..", "scenario-packs", "hctm")
	pack, err := loadPack(packDir)
	if err != nil {
		t.Fatalf("load pack: %v", err)
	}
	story, err := application.FindStory(pack, "order-expedite-01")
	if err != nil {
		t.Fatalf("find story: %v", err)
	}
	plan, err := application.CompilePlan(pack, "order-expedite-01")
	if err != nil {
		t.Fatalf("compile plan: %v", err)
	}

	events := make([]iaosclient.ScenarioObservedEvent, 0, len(story.Events.Events))
	for i, event := range story.Events.Events {
		events = append(events, iaosclient.ScenarioObservedEvent{
			Cursor:        int64(i + 1),
			EventID:       event.EventID,
			EventType:     event.EventType,
			OccurredAt:    time.Now().UTC().Add(time.Duration(i) * time.Minute),
			CorrelationID: plan.Correlation,
			BusinessType:  "sales_order",
			BusinessCode:  "SO-202607-0001",
			Payload:       json.RawMessage(`{}`),
		})
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/api/v1/profile"):
			_ = json.NewEncoder(w).Encode(iaosclient.ProfileResponse{Username: "aese-user", TenantID: pack.Manifest.TenantTemplate})
		case strings.HasSuffix(r.URL.Path, "/api/v1/scenarios/"+pack.Manifest.PackKey+"/"+story.Ref.Key+"/snapshot"):
			_ = json.NewEncoder(w).Encode(iaosclient.ScenarioSnapshot{
				PackKey:       pack.Manifest.PackKey,
				ScenarioKey:   story.Ref.Key,
				CorrelationID: plan.Correlation,
				Cursor:        int64(len(story.Events.Events)),
				Events: []map[string]any{
					{"event_id": story.Events.Events[0].EventID, "event_type": story.Events.Events[0].EventType, "cursor": 1, "correlation_id": plan.Correlation},
				},
			})
		case strings.HasSuffix(r.URL.Path, "/api/v1/scenarios/"+pack.Manifest.PackKey+"/"+story.Ref.Key+"/events"):
			_ = json.NewEncoder(w).Encode(iaosclient.ScenarioEventsResponse{Items: events, HasMore: false, NextCursor: int64(len(story.Events.Events))})
		case strings.HasSuffix(r.URL.Path, "/api/v1/scenarios/"+pack.Manifest.PackKey+"/"+story.Ref.Key+"/recommendations"):
			item := map[string]any{"run_id": "run-recommendation-recovery"}
			encoded, _ := json.Marshal(item)
			_ = json.NewEncoder(w).Encode(iaosclient.ScenarioRecommendationsResponse{Items: []json.RawMessage{encoded}})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	client, err := iaosclient.New(iaosclient.Config{BaseURL: server.URL, Token: "token", TenantID: pack.Manifest.TenantTemplate})
	if err != nil {
		t.Fatalf("iaos client: %v", err)
	}

	run := &runRecord{
		RunID:       "run-recommendation-recovery",
		PackKey:     pack.Manifest.PackKey,
		PackVersion: pack.Manifest.PackVersion,
		PackDir:     packDir,
		ScenarioKey: story.Ref.Key,
		Plan:        plan,
		Status:      application.RunStatusAnalyzing,
		CurrentAct:  0,
		TenantID:    pack.Manifest.TenantTemplate,
		Target:      server.URL,
		Token:       "token",
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		Cursor:      0,
		ActionCache: map[string]actionCache{},
	}

	s := New(Config{PackDir: packDir})
	if err := s.refreshRunFromFacts(context.Background(), run, client); err != nil {
		t.Fatalf("refresh run from facts: %v", err)
	}

	if run.Status != application.RunStatusAwaitingVerification {
		t.Fatalf("run status = %s, want %s", run.Status, application.RunStatusAwaitingVerification)
	}
	if run.CurrentAct != plan.ActCount {
		t.Fatalf("current act = %d, want %d", run.CurrentAct, plan.ActCount)
	}
}
