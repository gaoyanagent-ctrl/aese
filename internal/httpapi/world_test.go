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
