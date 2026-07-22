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
