package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGenesisWorldAPI(t *testing.T) {
	server := New(Config{})
	req := httptest.NewRequest(http.MethodGet, "/api/aese/v1/world/genesis", nil)
	res := httptest.NewRecorder()
	server.ServeHTTP(res, req)
	if res.Code != 200 {
		t.Fatalf("status=%d body=%s", res.Code, res.Body.String())
	}
	if !strings.Contains(res.Body.String(), "世界已退化") || !strings.Contains(res.Body.String(), "closed") {
		t.Fatalf("incomplete trace %s", res.Body.String())
	}
}

func TestIncorporationWorldAPI(t *testing.T) {
	server := New(Config{})
	req := httptest.NewRequest(http.MethodGet, "/api/aese/v1/world/incorporation", nil)
	res := httptest.NewRecorder()
	server.ServeHTTP(res, req)
	if res.Code != 200 {
		t.Fatalf("status=%d body=%s", res.Code, res.Body.String())
	}
	if !strings.Contains(res.Body.String(), `"campaign":"incorporation"`) || !strings.Contains(res.Body.String(), `"plant_project_eligible":true`) {
		t.Fatalf("incomplete campaign %s", res.Body.String())
	}
}

func TestPlantBuildWorldAPI(t *testing.T) {
	server := New(Config{})
	req := httptest.NewRequest(http.MethodGet, "/api/aese/v1/world/plant-build", nil)
	res := httptest.NewRecorder()
	server.ServeHTTP(res, req)
	if res.Code != 200 {
		t.Fatalf("status=%d body=%s", res.Code, res.Body.String())
	}
	if !strings.Contains(res.Body.String(), `"campaign":"plant-build"`) || !strings.Contains(res.Body.String(), `"capability_build_eligible":true`) {
		t.Fatalf("incomplete campaign %s", res.Body.String())
	}
}

func TestAESE3CompletionAPI(t *testing.T) {
	server := New(Config{})
	req := httptest.NewRequest(http.MethodGet, "/api/aese/v1/world/aese3", nil)
	res := httptest.NewRecorder()
	server.ServeHTTP(res, req)
	if res.Code != http.StatusOK { t.Fatalf("status=%d body=%s", res.Code, res.Body.String()) }
	body := res.Body.String()
	for _, want := range []string{`"code":"M17"`, `"code":"M24"`, `"industry_simulation_platform_ready":true`, `"automatic_business_writes":0`} {
		if !strings.Contains(body, want) { t.Fatalf("missing %s in %s", want, body) }
	}
}
