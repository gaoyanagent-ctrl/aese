package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/industrial-ai/iaos-aese/internal/genesisworkspace"
)

func TestGenesisWorkspaceListReturnsSessionExpiredForUpstreamUnauthorized(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, `{"error":"invalid authorization token"}`, http.StatusUnauthorized)
	}))
	defer upstream.Close()

	control := &genesisworkspace.ControlPlaneClient{BaseURL: upstream.URL}
	service := &genesisworkspace.Service{
		Store:        genesisworkspace.NewStore(filepath.Join(t.TempDir(), "workspaces.json")),
		ControlPlane: control,
	}
	server := New(Config{GenesisWorkspaceService: service, AllowLocalGenesisAuth: true})
	request := httptest.NewRequest(http.MethodGet, "/api/aese/v1/genesis/workspaces", nil)
	request.Header.Set("Authorization", "Bearer expired-player-token")
	request.Header.Set("X-Genesis-Player-Id", "player-local-test")
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("status=%d body=%s", response.Code, response.Body.String())
	}
	var body errorResponse
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Code != "player_session_expired" || body.Retryable {
		t.Fatalf("unexpected error response: %#v", body)
	}
}

func TestGenesisWorkspaceRequiresVerifiedPlayerSessionByDefault(t *testing.T) {
	service := &genesisworkspace.Service{
		Store: genesisworkspace.NewStore(filepath.Join(t.TempDir(), "workspaces.json")),
	}
	server := New(Config{GenesisWorkspaceService: service})
	request := httptest.NewRequest(http.MethodGet, "/api/aese/v1/genesis/workspaces", nil)
	request.Header.Set("X-Genesis-Player-Id", "client-controlled-player")
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("status=%d body=%s", response.Code, response.Body.String())
	}
	var body errorResponse
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Code != "player_session_required" {
		t.Fatalf("unexpected error response: %#v", body)
	}
}
