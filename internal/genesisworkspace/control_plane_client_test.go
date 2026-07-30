package genesisworkspace

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestControlPlanePreservesUpstreamUnauthorizedStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, `{"error":"invalid authorization token"}`, http.StatusUnauthorized)
	}))
	defer server.Close()

	client := ControlPlaneClient{BaseURL: server.URL}
	_, err := client.List(WithIAOSToken(context.Background(), "expired-token"), "player-local-001")
	if err == nil {
		t.Fatal("expected upstream authorization error")
	}
	var statusError interface{ UpstreamStatusCode() int }
	if !errors.As(err, &statusError) {
		t.Fatalf("error does not preserve upstream status: %T %v", err, err)
	}
	if statusError.UpstreamStatusCode() != http.StatusUnauthorized {
		t.Fatalf("status=%d want=%d", statusError.UpstreamStatusCode(), http.StatusUnauthorized)
	}
}

func TestControlPlaneCreateConfirmsWorldAndExchangesSession(t *testing.T) {
	var paths []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer player-token" {
			t.Fatalf("missing player bearer token")
		}
		paths = append(paths, r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.URL.Path == "/api/v1/genesis/workspaces":
			_, _ = w.Write([]byte(`{"workspace_id":"gxw-1","tenant_id":"tenant-gx-1","world_run_id":"world-gx-1","case_code":"INC-GX-1","display_name":"测试企业","status":"awaiting_world","current_checkpoint":"runtime_installed","created_at":"2026-07-29T00:00:00Z","updated_at":"2026-07-29T00:00:00Z"}`))
		case strings.HasSuffix(r.URL.Path, "/world-ready"):
			_, _ = w.Write([]byte(`{"workspace_id":"gxw-1","tenant_id":"tenant-gx-1","world_run_id":"world-gx-1","case_code":"INC-GX-1","display_name":"测试企业","status":"active","current_checkpoint":"tenant_active","created_at":"2026-07-29T00:00:00Z","updated_at":"2026-07-29T00:00:00Z"}`))
		case strings.HasSuffix(r.URL.Path, "/session"):
			_, _ = w.Write([]byte(`{"workspace":{"workspace_id":"gxw-1","tenant_id":"tenant-gx-1","world_run_id":"world-gx-1","case_code":"INC-GX-1","display_name":"测试企业","status":"active","current_checkpoint":"tenant_active","created_at":"2026-07-29T00:00:00Z","updated_at":"2026-07-29T00:00:00Z"},"token":"tenant-session"}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	client := ControlPlaneClient{BaseURL: server.URL}
	ctx := WithIAOSToken(context.Background(), "player-token")
	result, err := client.Create(ctx, "player-local-001", CreateRequest{
		OwnerPlayerID: "player-local-001", DisplayName: "测试企业",
		IdempotencyKey: "workspace-create-001",
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.TenantToken != "tenant-session" || result.Status != StatusActive {
		t.Fatalf("unexpected result: %#v", result)
	}
	if len(paths) != 3 || !strings.HasSuffix(paths[1], "/world-ready") || !strings.HasSuffix(paths[2], "/session") {
		t.Fatalf("unexpected control-plane sequence: %#v", paths)
	}
}

func TestControlPlaneIsDisabledWithoutPlayerToken(t *testing.T) {
	client := ControlPlaneClient{BaseURL: "http://127.0.0.1:8082"}
	if client.enabled(context.Background()) {
		t.Fatal("control plane enabled without authenticated player token")
	}
}

func TestControlPlaneAdoptsLegacyWorkspaceThenExchangesSession(t *testing.T) {
	var adoptionSeen bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/genesis/workspaces/legacy-adoptions":
			adoptionSeen = true
			if r.Header.Get("Authorization") != "Bearer player-token" {
				t.Fatalf("missing player token")
			}
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatal(err)
			}
			if body["workspace_id"] != "gxw-legacy1234" || body["case_code"] != "INC-LEGACY-1" {
				t.Fatalf("unexpected adoption body: %#v", body)
			}
			if _, exists := body["tenant_id"]; exists {
				t.Fatal("AESE must not submit a tenant claim during adoption")
			}
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"workspace_id":"gxw-legacy1234","tenant_id":"tenant-session","world_run_id":"world-legacy1234","case_code":"INC-LEGACY-1","display_name":"旧企业","status":"active","current_checkpoint":"tenant_active"}`)
		case "/api/v1/genesis/workspaces/gxw-legacy1234/session":
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"workspace":{"workspace_id":"gxw-legacy1234","tenant_id":"tenant-session","world_run_id":"world-legacy1234","case_code":"INC-LEGACY-1","display_name":"旧企业","status":"active","current_checkpoint":"session_issued"},"token":"tenant-token"}`)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	client := ControlPlaneClient{BaseURL: server.URL}
	ctx := WithIAOSToken(context.Background(), "player-token")
	result, err := client.AdoptLegacy(ctx, "player-1", Workspace{
		WorkspaceID: "gxw-legacy1234", WorldRunID: "world-legacy1234",
		CaseCode: "INC-LEGACY-1", DisplayName: "旧企业",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !adoptionSeen || result.TenantToken != "tenant-token" {
		t.Fatalf("adoption/session was not completed: %#v", result)
	}
}
