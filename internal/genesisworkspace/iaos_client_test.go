package genesisworkspace

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIAOSClientReturnsFounderSessionForIncorporation(t *testing.T) {
	var runtimeAuthorization string
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/dev/token", func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]string{"token": "platform-token"})
	})
	mux.HandleFunc("/api/v1/platform/tenants/tenant-gx-test", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/api/v1/platform/tenants", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"status":"provisioning"}`))
	})
	mux.HandleFunc("/api/v1/platform-identities/founder/bootstrap", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})
	mux.HandleFunc("/api/v1/auth/login", func(w http.ResponseWriter, r *http.Request) {
		var input map[string]string
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			t.Fatal(err)
		}
		if input["tenant_id"] != "tenant-gx-test" || input["username"] != "founder-principal" {
			t.Fatalf("unexpected founder login: %#v", input)
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"token": "founder-token"})
	})
	mux.HandleFunc("/api/v1/incorporations/runtime/install", func(w http.ResponseWriter, r *http.Request) {
		runtimeAuthorization = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusCreated)
	})
	mux.HandleFunc("/api/v1/platform/tenants/tenant-gx-test/activate", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	token, err := (IAOSClient{BaseURL: server.URL, PlatformToken: "platform-token"}).Provision(
		context.Background(),
		Workspace{
			WorkspaceID: "gxw-test", TenantID: "tenant-gx-test",
			DisplayName: "测试企业", OwnerPlayerID: "player-test", WorldRunID: "world-test",
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	if token != "founder-token" {
		t.Fatalf("browser session = %q, want founder token", token)
	}
	if runtimeAuthorization != "Bearer founder-token" {
		t.Fatalf("runtime authorization = %q, want founder session", runtimeAuthorization)
	}
}
