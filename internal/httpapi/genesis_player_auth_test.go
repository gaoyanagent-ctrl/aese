package httpapi

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/industrial-ai/iaos-aese/internal/genesisworkspace"
)

func TestGenesisPlayerLoginProxyPreservesCredentialsAndDoesNotLogPassword(t *testing.T) {
	const password = "SafePassword123"
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/genesis/auth/login" {
			t.Fatalf("path=%s", r.URL.Path)
		}
		raw, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		var request genesisworkspace.PlayerLoginRequest
		if err = json.Unmarshal(raw, &request); err != nil {
			t.Fatal(err)
		}
		if request.Username != "founder-principal" || request.Password != password {
			t.Fatalf("unexpected request: %#v", request)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"status":"success","token":"player-token","player":{"subject_id":"founder-principal","username":"founder-principal","display_name":"Founder"}}`)
	}))
	defer upstream.Close()

	var logs bytes.Buffer
	server := New(Config{
		Logger:            log.New(&logs, "", 0),
		GenesisPlayerAuth: &genesisworkspace.PlayerAuthClient{BaseURL: upstream.URL},
	})
	request := httptest.NewRequest(http.MethodPost, "/api/aese/v1/auth/login",
		strings.NewReader(`{"username":"founder-principal","password":"`+password+`"}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", response.Code, response.Body.String())
	}
	if strings.Contains(logs.String(), password) {
		t.Fatal("password leaked to server logs")
	}
}

func TestGenesisPlayerRegistrationProxyPreservesValidationStatus(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = io.WriteString(w, `{"error":"password must contain 10-128 characters","code":"invalid_registration"}`)
	}))
	defer upstream.Close()

	server := New(Config{GenesisPlayerAuth: &genesisworkspace.PlayerAuthClient{BaseURL: upstream.URL}})
	request := httptest.NewRequest(http.MethodPost, "/api/aese/v1/auth/register",
		strings.NewReader(`{"username":"new-founder","password":"short","display_name":"New Founder"}`))
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	if response.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status=%d body=%s", response.Code, response.Body.String())
	}
	var body errorResponse
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Code != "invalid_registration" {
		t.Fatalf("unexpected error response: %#v", body)
	}
}
