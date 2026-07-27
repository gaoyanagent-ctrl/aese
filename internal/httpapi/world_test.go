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

func TestIncorporationGameProjectionAPI(t *testing.T) {
	server := New(Config{})
	req := httptest.NewRequest(http.MethodGet, "/api/aese/v1/game/incorporation/INC-DEMO-001/projection?frame=2", nil)
	res := httptest.NewRecorder()
	server.ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", res.Code, res.Body.String())
	}
	for _, want := range []string{`"chapter":"banking"`, `"case_code":"INC-DEMO-001"`, `"evidence_refs"`} {
		if !strings.Contains(res.Body.String(), want) {
			t.Fatalf("missing %s in %s", want, res.Body.String())
		}
	}
}

func TestCreativeIntentAndNamingAPIs(t *testing.T) {
	server := New(Config{})
	intentBody := `{"tenant_id":"tenant-demo","case_code":"INC-DEMO-001","raw_idea":"为工业客户创建可靠且可持续的热管理技术公司","industry":"热管理","customers":["新能源汽车制造商"],"offerings":["电池冷却板"],"brand_traits":["可靠","工程","长期主义"],"capital_minor":"100000000","risk_appetite":"balanced"}`
	req := httptest.NewRequest(http.MethodPost, "/api/aese/v1/game/creative/intent", strings.NewReader(intentBody))
	res := httptest.NewRecorder()
	server.ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("intent status=%d body=%s", res.Code, res.Body.String())
	}
	if !strings.Contains(res.Body.String(), `"intent_id":"intent-inc-demo-001"`) {
		t.Fatalf("unexpected intent %s", res.Body.String())
	}

	req = httptest.NewRequest(http.MethodPost, "/api/aese/v1/game/creative/names", strings.NewReader(res.Body.String()))
	res = httptest.NewRecorder()
	server.ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("names status=%d body=%s", res.Code, res.Body.String())
	}
	for _, want := range []string{`"status":"candidate_only"`, `"chinese_name":"澄流热管理有限公司"`, `"risk_hints"`} {
		if !strings.Contains(res.Body.String(), want) {
			t.Fatalf("missing %s in %s", want, res.Body.String())
		}
	}
}

func TestCreativeAPIRejectsUnknownFields(t *testing.T) {
	server := New(Config{})
	body := `{"tenant_id":"t","case_code":"c","raw_idea":"足够长的企业创业构想内容用于校验","industry":"工业","brand_traits":["可靠"],"risk_appetite":"balanced","untrusted_fact":true}`
	req := httptest.NewRequest(http.MethodPost, "/api/aese/v1/game/creative/intent", strings.NewReader(body))
	res := httptest.NewRecorder()
	server.ServeHTTP(res, req)
	if res.Code != http.StatusBadRequest || !strings.Contains(res.Body.String(), "unknown field") {
		t.Fatalf("status=%d body=%s", res.Code, res.Body.String())
	}
}

func TestGameProjectionConsumesLiveIAOSCommittedState(t *testing.T) {
	iaos := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer founder-token" {
			t.Errorf("authorization not forwarded: %q", r.Header.Get("Authorization"))
		}
		switch {
		case strings.HasSuffix(r.URL.Path, "/evidence"):
			_, _ = w.Write([]byte(`{"verified":true,"trace":{"schema_version":"1.0","case_code":"INC-LIVE","state":{"case_code":"INC-LIVE","tenant_id":"tenant-live","state":"incorporation_case_opened","commitment_minor":0,"contribution_minor":0,"budget_authorized_minor":0,"currency":"CNY","proposed_company_name":"澄流热管理有限公司"},"journal":[{"sequence":1,"capability_code":"incorporation.case.open","actor_type":"human","actor_id":"founder-principal","correlation_id":"corr-live","created_at":"2026-07-27T10:00:00Z"}],"world_exchanges":[]}}`))
		case strings.HasSuffix(r.URL.Path, "/work-items"):
			_, _ = w.Write([]byte(`{"items":[{"sequence":1,"capability_code":"incorporation.case.open","task_type":"human_task","participant_id":"founder-principal","status":"completed","effective_status":"completed"},{"sequence":2,"capability_code":"founder.resolution.prepare","task_type":"agent_task","participant_id":"incorporation-agent","status":"ready","effective_status":"ready"}]}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer iaos.Close()
	server := New(Config{IAOSBaseURL: iaos.URL})
	req := httptest.NewRequest(http.MethodGet, "/api/aese/v1/game/incorporation/INC-LIVE/projection", nil)
	req.Header.Set("Authorization", "Bearer founder-token")
	req.Header.Set("X-IAOS-Tenant-Id", "tenant-live")
	res := httptest.NewRecorder()
	server.ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", res.Code, res.Body.String())
	}
	for _, want := range []string{`"projection_id":"gp-inc-live-1"`, `"company_name":"澄流热管理有限公司"`, `"work_item_id":"WI-02"`} {
		if !strings.Contains(res.Body.String(), want) {
			t.Fatalf("missing %s in %s", want, res.Body.String())
		}
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
	if res.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", res.Code, res.Body.String())
	}
	body := res.Body.String()
	for _, want := range []string{`"code":"M17"`, `"code":"M24"`, `"industry_simulation_platform_ready":true`, `"automatic_business_writes":0`} {
		if !strings.Contains(body, want) {
			t.Fatalf("missing %s in %s", want, body)
		}
	}
}
