package httpapi

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIAOSCommandGatewayForwardsAllowListedCommandAndIdentity(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/incorporations/INC-GX-1/work-items/2/dispatch-agent" {
			t.Fatalf("path=%s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer founder-token" {
			t.Fatalf("authorization=%q", r.Header.Get("Authorization"))
		}
		if r.Header.Get("X-IAOS-Tenant-Id") != "tenant-gx-1" {
			t.Fatalf("tenant=%q", r.Header.Get("X-IAOS-Tenant-Id"))
		}
		body, _ := io.ReadAll(r.Body)
		if string(body) != `{"business_note":"prepare resolution"}` {
			t.Fatalf("body=%s", body)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"accepted"}`))
	}))
	defer upstream.Close()

	server := New(Config{IAOSBaseURL: upstream.URL})
	request := httptest.NewRequest(
		http.MethodPost,
		"/api/aese/v1/commands/iaos/incorporations/INC-GX-1/work-items/2/dispatch-agent",
		strings.NewReader(`{"business_note":"prepare resolution"}`),
	)
	request.Header.Set("Authorization", "Bearer founder-token")
	request.Header.Set("X-IAOS-Tenant-Id", "tenant-gx-1")
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)

	if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), `"accepted"`) {
		t.Fatalf("status=%d body=%s", response.Code, response.Body.String())
	}
}

func TestIAOSCommandGatewayRejectsArbitraryIAOSRoute(t *testing.T) {
	server := New(Config{IAOSBaseURL: "http://127.0.0.1:1"})
	request := httptest.NewRequest(
		http.MethodPost,
		"/api/aese/v1/commands/iaos/entities/account/records",
		strings.NewReader(`{}`),
	)
	request.Header.Set("Authorization", "Bearer founder-token")
	request.Header.Set("X-IAOS-Tenant-Id", "tenant-gx-1")
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)

	if response.Code != http.StatusForbidden ||
		!strings.Contains(response.Body.String(), "command_route_not_allowed") {
		t.Fatalf("status=%d body=%s", response.Code, response.Body.String())
	}
}

func TestIAOSCommandGatewayAllowsAssignedApprovalRejection(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/approvals/approval-1/reject" {
			t.Fatalf("request=%s %s", r.Method, r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		if string(body) != `{"note":"现场条件不满足投产要求"}` {
			t.Fatalf("body=%s", body)
		}
		_, _ = w.Write([]byte(`{"item":{"id":"approval-1","status":"rejected"}}`))
	}))
	defer upstream.Close()

	server := New(Config{IAOSBaseURL: upstream.URL})
	request := httptest.NewRequest(http.MethodPost, "/api/aese/v1/commands/iaos/approvals/approval-1/reject", strings.NewReader(`{"note":"现场条件不满足投产要求"}`))
	request.Header.Set("Authorization", "Bearer founder-token")
	request.Header.Set("X-IAOS-Tenant-Id", "tenant-gx-1")
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)
	if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), `"rejected"`) {
		t.Fatalf("status=%d body=%s", response.Code, response.Body.String())
	}
}
