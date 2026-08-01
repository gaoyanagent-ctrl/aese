package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/industrial-ai/iaos-aese/internal/creative"
	"github.com/industrial-ai/iaos-aese/internal/plantbuild"
)

type planningProviderStub struct{ calls int }

func (p *planningProviderStub) Status() plantbuild.PlanningProviderStatus {
	return plantbuild.PlanningProviderStatus{State: "connected", Provider: "MiniMax", Model: "MiniMax-M3", PromptVersion: plantbuild.PlanningPromptVersion}
}
func (p *planningProviderStub) Generate(_ context.Context, request plantbuild.FacilityRequirement) (plantbuild.ProposalSet, error) {
	p.calls++
	return plantbuild.ProposalSet{SchemaVersion: "1.0", ProposalSetID: "SET-1", RequirementID: request.RequirementID, Revision: 1, Status: "candidate_only", Evidence: plantbuild.ProposalEvidence{Provider: "MiniMax", Model: "MiniMax-M3", PromptVersion: plantbuild.PlanningPromptVersion, RequestID: "req-1", InputHash: "sha256:input", OutputHash: "sha256:output", TokenUsage: map[string]int{"total_tokens": 12}, ValidatedAt: "2026-08-01T10:00:00Z"}}, nil
}

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
		case strings.Contains(r.URL.Path, "/finance/opening/"):
			_, _ = w.Write([]byte(`{"case_code":"INC-LIVE","finance_opening_ready":false,"roles":[],"debit_minor":0,"credit_minor":0,"trial_balance":[]}`))
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

func TestGameProjectionPreservesMissingCaseAsNotFound(t *testing.T) {
	iaos := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}))
	defer iaos.Close()
	server := New(Config{IAOSBaseURL: iaos.URL})
	req := httptest.NewRequest(http.MethodGet, "/api/aese/v1/game/incorporation/INC-NEW/projection", nil)
	req.Header.Set("Authorization", "Bearer founder-token")
	req.Header.Set("X-IAOS-Tenant-Id", "tenant-new")
	res := httptest.NewRecorder()
	server.ServeHTTP(res, req)
	if res.Code != http.StatusNotFound {
		t.Fatalf("status=%d body=%s", res.Code, res.Body.String())
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

func TestPlantPlanningStatusFailsClosedWithoutExternalModel(t *testing.T) {
	server := New(Config{})
	req := httptest.NewRequest(http.MethodGet, "/api/aese/v1/world/plant-build/planning-status", nil)
	res := httptest.NewRecorder()
	server.ServeHTTP(res, req)
	if res.Code != http.StatusOK || !strings.Contains(res.Body.String(), `"state":"not_configured"`) {
		t.Fatalf("status=%d body=%s", res.Code, res.Body.String())
	}
}

func TestPlantFinancialConstraintUsesNarrowAuthenticatedIAOSRead(t *testing.T) {
	iaos := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/genesis/plant/interactive/financial-constraints" || r.URL.Query().Get("case_code") != "INC-LIVE" {
			t.Fatalf("unexpected IAOS request %s?%s", r.URL.Path, r.URL.RawQuery)
		}
		if r.Header.Get("Authorization") != "Bearer founder-token" || r.Header.Get("X-IAOS-Tenant-Id") != "tenant-live" {
			t.Fatalf("identity not forwarded: %#v", r.Header)
		}
		_, _ = w.Write([]byte(`{"case_code":"INC-LIVE","legal_entity_code":"LE-LIVE","financial_constraint":{"available_cash":{"value":"30000000.00","currency":"CNY","scale":2},"approved_budget":{"value":"20000000.00","currency":"CNY","scale":2},"cash_source_ref":"gl:BOOK-INC-LIVE:1002","budget_source_ref":"budget:BUDGET-INC-LIVE","snapshot_hash":"sha256:authority"}}`))
	}))
	defer iaos.Close()
	server := New(Config{IAOSBaseURL: iaos.URL})
	req := httptest.NewRequest(http.MethodGet, "/api/aese/v1/world/plant-build/financial-constraints?case_code=INC-LIVE", nil)
	req.Header.Set("Authorization", "Bearer founder-token")
	req.Header.Set("X-IAOS-Tenant-Id", "tenant-live")
	res := httptest.NewRecorder()
	server.ServeHTTP(res, req)
	if res.Code != http.StatusOK || !strings.Contains(res.Body.String(), `"snapshot_hash":"sha256:authority"`) {
		t.Fatalf("status=%d body=%s", res.Code, res.Body.String())
	}
}

func TestPlantProposalSetUsesNarrowAuthenticatedIAOSRead(t *testing.T) {
	iaos := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/genesis/plant/interactive/proposal-sets" || r.URL.Query().Get("requirement_id") != "REQ-LIVE" {
			t.Fatalf("unexpected IAOS request %s?%s", r.URL.Path, r.URL.RawQuery)
		}
		if r.Header.Get("Authorization") != "Bearer founder-token" || r.Header.Get("X-IAOS-Tenant-Id") != "tenant-live" {
			t.Fatalf("identity not forwarded: %#v", r.Header)
		}
		_, _ = w.Write([]byte(`{"schema_version":"1.0","proposal_set_id":"SET-LIVE","requirement_id":"REQ-LIVE","revision":1,"status":"candidate_only","proposals":[],"evidence":{"provider":"MiniMax","model":"M3","prompt_version":"v1","input_hash":"sha256:in","output_hash":"sha256:out","validated_at":"2026-08-01T00:00:00Z"}}`))
	}))
	defer iaos.Close()
	server := New(Config{IAOSBaseURL: iaos.URL})
	req := httptest.NewRequest(http.MethodGet, "/api/aese/v1/world/plant-build/proposals?requirement_id=REQ-LIVE", nil)
	req.Header.Set("Authorization", "Bearer founder-token")
	req.Header.Set("X-IAOS-Tenant-Id", "tenant-live")
	res := httptest.NewRecorder()
	server.ServeHTTP(res, req)
	if res.Code != http.StatusOK || !strings.Contains(res.Body.String(), `"proposal_set_id":"SET-LIVE"`) {
		t.Fatalf("status=%d body=%s", res.Code, res.Body.String())
	}
}

func TestPlantPlanningProposalPersistsEvidenceAndReplaysIdempotently(t *testing.T) {
	provider := &planningProviderStub{}
	server := New(Config{PlantPlanningProvider: provider, CreativeJobStore: creative.NewJobStore(t.TempDir() + "/jobs.json"), AllowLocalPlantAuthority: true})
	body := `{"schema_version":"1.0","requirement_id":"REQ-1","tenant_id":"tenant-a","case_code":"INC-1","legal_entity_code":"LE-1","target_region":"华东","facility_purpose":"汽车零部件制造","minimum_area_m2":12000,"minimum_electricity_kva":2200,"target_available_at":"2027-01-01T00:00:00+08:00","candidate_count":2,"allowed_option_types":["leased_shell","build_to_suit"],"investment_request":{"value":"18000000.00","currency":"CNY","scale":2},"minimum_cash_reserve":{"value":"5000000.00","currency":"CNY","scale":2},"financial_constraint":{"available_cash":{"value":"30000000.00","currency":"CNY","scale":2},"approved_budget":{"value":"20000000.00","currency":"CNY","scale":2},"cash_source_ref":"ledger:CASH-1","budget_source_ref":"budget:BUD-1","snapshot_hash":"sha256:abc"},"preferences":["优先投产速度"],"revision":1,"revision_reason":"首次规划"}`
	for attempt := 0; attempt < 2; attempt++ {
		req := httptest.NewRequest(http.MethodPost, "/api/aese/v1/world/plant-build/proposals", strings.NewReader(body))
		res := httptest.NewRecorder()
		server.ServeHTTP(res, req)
		if res.Code != http.StatusOK || !strings.Contains(res.Body.String(), `"kind":"facility_planning"`) {
			t.Fatalf("attempt=%d status=%d body=%s", attempt, res.Code, res.Body.String())
		}
	}
	if provider.calls != 1 {
		t.Fatalf("provider calls=%d", provider.calls)
	}
}

func TestPlantPlanningCommitsRequirementAndAgentProposalThroughIAOSCapabilities(t *testing.T) {
	capabilities := []string{}
	iaos := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/profile" {
			_, _ = w.Write([]byte(`{"username":"founder-principal","tenant_id":"tenant-a"}`))
			return
		}
		if r.URL.Path != "/api/v1/genesis/plant/interactive/actions" || r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		var command struct {
			CapabilityCode string `json:"capability_code"`
			ActorID        string `json:"actor_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&command); err != nil {
			t.Fatal(err)
		}
		if command.ActorID != "founder-principal" {
			t.Fatalf("actor=%s", command.ActorID)
		}
		capabilities = append(capabilities, command.CapabilityCode)
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"status":"committed"}`))
	}))
	defer iaos.Close()
	provider := &planningProviderStub{}
	server := New(Config{IAOSBaseURL: iaos.URL, PlantPlanningProvider: provider, CreativeJobStore: creative.NewJobStore(t.TempDir() + "/jobs.json")})
	body := `{"schema_version":"1.0","requirement_id":"REQ-1","tenant_id":"tenant-a","case_code":"INC-1","legal_entity_code":"LE-1","target_region":"华东","facility_purpose":"汽车零部件制造","minimum_area_m2":12000,"minimum_electricity_kva":2200,"target_available_at":"2027-01-01T00:00:00+08:00","candidate_count":2,"allowed_option_types":["leased_shell","build_to_suit"],"investment_request":{"value":"18000000.00","currency":"CNY","scale":2},"minimum_cash_reserve":{"value":"5000000.00","currency":"CNY","scale":2},"financial_constraint":{"available_cash":{"value":"30000000.00","currency":"CNY","scale":2},"approved_budget":{"value":"20000000.00","currency":"CNY","scale":2},"cash_source_ref":"ledger:CASH-1","budget_source_ref":"budget:BUD-1","snapshot_hash":"sha256:abc"},"preferences":["优先投产速度"],"revision":1,"revision_reason":"首次规划"}`
	req := httptest.NewRequest(http.MethodPost, "/api/aese/v1/world/plant-build/proposals", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer founder-token")
	req.Header.Set("X-IAOS-Tenant-Id", "tenant-a")
	res := httptest.NewRecorder()
	server.ServeHTTP(res, req)
	if res.Code != http.StatusOK || !strings.Contains(res.Body.String(), `"authority_status":"committed"`) {
		t.Fatalf("status=%d body=%s", res.Code, res.Body.String())
	}
	if strings.Join(capabilities, ",") != "facility.requirement.define,site.proposal.record" {
		t.Fatalf("capabilities=%v", capabilities)
	}
}

func TestManualPlantProposalAppendsAuthorityRevisionWithAuthenticatedActor(t *testing.T) {
	var command map[string]any
	proposalSet := `{"schema_version":"1.0","proposal_set_id":"SET-1","requirement_id":"REQ-1","revision":1,"status":"candidate_only","proposals":[{"proposal_id":"P-1","option_type":"lease_and_retrofit","display_name":"候选一","business_rationale":"满足初始规划需求","estimated_amount":{"minimum":{"value":"7000000.00","currency":"CNY","scale":2},"likely":{"value":"8000000.00","currency":"CNY","scale":2},"maximum":{"value":"9000000.00","currency":"CNY","scale":2},"basis":"概念估算"},"estimated_schedule":{"earliest":"2026-09-01T00:00:00Z","likely":"2026-10-01T00:00:00Z","latest":"2026-11-01T00:00:00Z"},"assumptions":["可租赁"],"facts_required":["权属"],"risks":["延期"],"source_refs":["requirement:REQ-1"],"confidence":"0.7","status":"proposed"},{"proposal_id":"P-2","option_type":"build_to_suit","display_name":"候选二","business_rationale":"满足初始规划需求","estimated_amount":{"minimum":{"value":"7000000.00","currency":"CNY","scale":2},"likely":{"value":"8000000.00","currency":"CNY","scale":2},"maximum":{"value":"9000000.00","currency":"CNY","scale":2},"basis":"概念估算"},"estimated_schedule":{"earliest":"2026-09-01T00:00:00Z","likely":"2026-10-01T00:00:00Z","latest":"2026-11-01T00:00:00Z"},"assumptions":["可建设"],"facts_required":["许可"],"risks":["延期"],"source_refs":["requirement:REQ-1"],"confidence":"0.6","status":"proposed"}],"evidence":{"provider":"MiniMax","model":"M3","prompt_version":"plant-planning-v1","input_hash":"sha256:in","output_hash":"sha256:out","validated_at":"2026-08-01T00:00:00Z"}}`
	requirement := `{"schema_version":"1.0","requirement_id":"REQ-1","tenant_id":"tenant-a","case_code":"INC-1","legal_entity_code":"LE-1","target_region":"华东","facility_purpose":"制造","minimum_area_m2":5000,"minimum_electricity_kva":2000,"target_available_at":"2027-01-01T00:00:00Z","candidate_count":2,"allowed_option_types":["lease_and_retrofit","build_to_suit"],"investment_request":{"value":"12000000.00","currency":"CNY","scale":2},"minimum_cash_reserve":{"value":"3000000.00","currency":"CNY","scale":2},"financial_constraint":{"available_cash":{"value":"20000000.00","currency":"CNY","scale":2},"approved_budget":{"value":"12000000.00","currency":"CNY","scale":2},"cash_source_ref":"ledger:CASH","budget_source_ref":"budget:B1","snapshot_hash":"sha256:authority"},"preferences":["工期"],"revision":1,"revision_reason":"首次规划"}`
	iaos := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/api/v1/profile":
			_, _ = w.Write([]byte(`{"username":"project-owner","tenant_id":"tenant-a"}`))
		case r.URL.Path == "/api/v1/genesis/plant/interactive/proposal-sets":
			_, _ = w.Write([]byte(proposalSet))
		case r.URL.Path == "/api/v1/genesis/plant/interactive/requirements/REQ-1":
			_, _ = w.Write([]byte(requirement))
		case r.URL.Path == "/api/v1/genesis/plant/interactive/actions":
			if err := json.NewDecoder(r.Body).Decode(&command); err != nil {
				t.Fatal(err)
			}
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"status":"committed"}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer iaos.Close()
	server := New(Config{IAOSBaseURL: iaos.URL})
	body := `{"requirement_id":"REQ-1","proposal_set_id":"SET-1","expected_revision":1,"proposal":{"option_type":"lease_and_retrofit","display_name":"人工园区候选","business_rationale":"项目负责人根据招商线索补充并要求正式调研","estimated_amount":{"minimum":{"value":"7000000.00","currency":"CNY","scale":2},"likely":{"value":"8000000.00","currency":"CNY","scale":2},"maximum":{"value":"9000000.00","currency":"CNY","scale":2},"basis":"人工概念估算"},"estimated_schedule":{"earliest":"2026-10-01T00:00:00Z","likely":"2026-10-01T00:00:00Z","latest":"2026-10-01T00:00:00Z"},"assumptions":["园区存在可租赁空间"],"facts_required":["权属与正式报价"],"risks":["交付日期未核验"],"source_refs":["manual:user-input"],"confidence":"0.40"}}`
	req := httptest.NewRequest(http.MethodPost, "/api/aese/v1/world/plant-build/proposals/manual", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer owner-token")
	req.Header.Set("X-IAOS-Tenant-Id", "tenant-a")
	res := httptest.NewRecorder()
	server.ServeHTTP(res, req)
	if res.Code != http.StatusCreated {
		t.Fatalf("status=%d body=%s", res.Code, res.Body.String())
	}
	if command["actor_id"] != "project-owner" || command["capability_code"] != "site.proposal.record" {
		t.Fatalf("command=%v", command)
	}
	set := command["proposal_set"].(map[string]any)
	if set["revision"] != float64(2) {
		t.Fatalf("set=%v", set)
	}
	evidence := set["evidence"].(map[string]any)
	if evidence["source_type"] != "human_manual" || evidence["parent_revision"] != float64(1) {
		t.Fatalf("evidence=%v", evidence)
	}
	proposals := set["proposals"].([]any)
	added := proposals[2].(map[string]any)
	refs := added["source_refs"].([]any)
	if refs[len(refs)-1] != "human:project-owner" {
		t.Fatalf("refs=%v", refs)
	}
}

func TestManualPlantProposalCreatesFirstAuthorityRevisionWithoutAgentSet(t *testing.T) {
	var command map[string]any
	requirement := `{"schema_version":"1.0","requirement_id":"REQ-1","tenant_id":"tenant-a","case_code":"INC-1","legal_entity_code":"LE-1","target_region":"华东","facility_purpose":"制造","minimum_area_m2":5000,"minimum_electricity_kva":2000,"target_available_at":"2027-01-01T00:00:00Z","candidate_count":2,"allowed_option_types":["lease_and_retrofit","build_to_suit"],"investment_request":{"value":"12000000.00","currency":"CNY","scale":2},"minimum_cash_reserve":{"value":"3000000.00","currency":"CNY","scale":2},"financial_constraint":{"available_cash":{"value":"20000000.00","currency":"CNY","scale":2},"approved_budget":{"value":"12000000.00","currency":"CNY","scale":2},"cash_source_ref":"ledger:CASH","budget_source_ref":"budget:B1","snapshot_hash":"sha256:authority"},"preferences":["工期"],"revision":1,"revision_reason":"首次规划"}`
	iaos := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/api/v1/profile":
			_, _ = w.Write([]byte(`{"username":"project-owner","tenant_id":"tenant-a"}`))
		case r.URL.Path == "/api/v1/genesis/plant/interactive/requirements/REQ-1":
			_, _ = w.Write([]byte(requirement))
		case r.URL.Path == "/api/v1/genesis/plant/interactive/proposal-sets":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error":"proposal set not found","code":"not_found"}`))
		case r.URL.Path == "/api/v1/genesis/plant/interactive/actions":
			if err := json.NewDecoder(r.Body).Decode(&command); err != nil {
				t.Fatal(err)
			}
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"status":"committed"}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer iaos.Close()
	server := New(Config{IAOSBaseURL: iaos.URL})
	body := `{"requirement_id":"REQ-1","proposal_set_id":"","expected_revision":0,"proposal":{"option_type":"lease_and_retrofit","display_name":"人工首个候选","business_rationale":"项目负责人在模型不可用时建立首个权威候选","estimated_amount":{"minimum":{"value":"7000000.00","currency":"CNY","scale":2},"likely":{"value":"8000000.00","currency":"CNY","scale":2},"maximum":{"value":"9000000.00","currency":"CNY","scale":2},"basis":"人工概念估算"},"estimated_schedule":{"earliest":"2026-10-01T00:00:00Z","likely":"2026-10-01T00:00:00Z","latest":"2026-10-01T00:00:00Z"},"assumptions":["园区存在可租赁空间"],"facts_required":["权属与正式报价"],"risks":["交付日期未核验"],"source_refs":["manual:user-input"],"confidence":"0.40"}}`
	req := httptest.NewRequest(http.MethodPost, "/api/aese/v1/world/plant-build/proposals/manual", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer owner-token")
	req.Header.Set("X-IAOS-Tenant-Id", "tenant-a")
	res := httptest.NewRecorder()
	server.ServeHTTP(res, req)
	if res.Code != http.StatusCreated {
		t.Fatalf("status=%d body=%s", res.Code, res.Body.String())
	}
	set := command["proposal_set"].(map[string]any)
	if set["revision"] != float64(1) || set["proposal_set_id"] != "manual-proposal-set-REQ-1" || len(set["proposals"].([]any)) != 1 {
		t.Fatalf("set=%v", set)
	}
	evidence := set["evidence"].(map[string]any)
	if evidence["parent_revision"] != float64(0) || !strings.HasPrefix(evidence["input_hash"].(string), "sha256:") {
		t.Fatalf("evidence=%v", evidence)
	}
}

func TestPlantProposalReviewOverwritesActorFromAuthenticatedProfile(t *testing.T) {
	var received map[string]any
	iaos := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/profile" {
			_, _ = w.Write([]byte(`{"username":"project-owner","tenant_id":"tenant-a"}`))
			return
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatal(err)
		}
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"status":"committed"}`))
	}))
	defer iaos.Close()
	server := New(Config{IAOSBaseURL: iaos.URL})
	req := httptest.NewRequest(http.MethodPost, "/api/aese/v1/world/plant-build/reviews", strings.NewReader(`{"proposal_set_id":"SET-1","proposal_id":"SITE-1","action":"adopt_for_investigation","reason":"进入外部调研并核验权属","expected_revision":1}`))
	req.Header.Set("Authorization", "Bearer owner-token")
	req.Header.Set("X-IAOS-Tenant-Id", "tenant-a")
	res := httptest.NewRecorder()
	server.ServeHTTP(res, req)
	if res.Code != http.StatusCreated {
		t.Fatalf("status=%d body=%s", res.Code, res.Body.String())
	}
	if received["actor_id"] != "project-owner" {
		t.Fatalf("command=%v", received)
	}
	review := received["proposal_review"].(map[string]any)
	if review["reviewed_by"] != "project-owner" || review["reviewed_at"] == "" {
		t.Fatalf("review=%v", review)
	}
}

func TestPlantInvestigationObservationUsesWorldBridgeBeforeCapabilityCommit(t *testing.T) {
	paths := []string{}
	iaos := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/profile" {
			_, _ = w.Write([]byte(`{"username":"project-owner","tenant_id":"tenant-a"}`))
			return
		}
		paths = append(paths, r.URL.Path)
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"status":"committed"}`))
	}))
	defer iaos.Close()
	server := New(Config{IAOSBaseURL: iaos.URL})
	body := `{"case_code":"INC-1","world_run_id":"plant-build-INC-1","observation":{"schema_version":"1.0","observation_id":"OBS-1","investigation_request_id":"INV-1","proposal_id":"SITE-1","result":"completed","ownership_status":"verified","available_area_m2":9000,"electricity_kva":3000,"quoted_amount":{"value":"9800000.00","currency":"CNY","scale":2},"available_at":"2026-10-01T00:00:00Z","permit_status":"eligible","evidence_refs":["world-document:QUOTE-1"],"notes":"现场核验完成","external_actor_id":"virtual-park-operator","observed_at":"2026-08-01T11:00:00Z"}}`
	req := httptest.NewRequest(http.MethodPost, "/api/aese/v1/world/plant-build/observations", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer owner-token")
	req.Header.Set("X-IAOS-Tenant-Id", "tenant-a")
	res := httptest.NewRecorder()
	server.ServeHTTP(res, req)
	if res.Code != http.StatusCreated {
		t.Fatalf("status=%d body=%s", res.Code, res.Body.String())
	}
	want := "/api/v1/world-bridge/observations,/api/v1/genesis/plant/interactive/actions"
	if strings.Join(paths, ",") != want {
		t.Fatalf("paths=%v, want governed World-first order", paths)
	}
}

func TestPlantSiteRecommendationUsesAuthenticatedCommandGateway(t *testing.T) {
	var command map[string]any
	iaos := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/profile" {
			_, _ = w.Write([]byte(`{"username":"project-owner","tenant_id":"tenant-a"}`))
			return
		}
		if r.URL.Path != "/api/v1/genesis/plant/interactive/actions" {
			http.NotFound(w, r)
			return
		}
		if err := json.NewDecoder(r.Body).Decode(&command); err != nil {
			t.Fatal(err)
		}
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"status":"committed","result":{"approval_request_id":"APR-1"}}`))
	}))
	defer iaos.Close()
	server := New(Config{IAOSBaseURL: iaos.URL})
	body := `{"schema_version":"1.0","recommendation_id":"REC-1","case_code":"INC-1","proposal_set_id":"SET-1","proposal_set_revision":1,"selected_proposal_id":"SITE-1","assessment_policy_version":"site-assessment-v1","weights":{"cost":35,"schedule":25,"capacity":20,"control":20},"recommendation_reason":"综合可信现场事实推荐候选一","alternative_comparison":"候选二在工期和产能方面劣于候选一","recommended_at":"2026-08-01T12:00:00Z"}`
	req := httptest.NewRequest(http.MethodPost, "/api/aese/v1/world/plant-build/site-selections", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer owner-token")
	req.Header.Set("X-IAOS-Tenant-Id", "tenant-a")
	res := httptest.NewRecorder()
	server.ServeHTTP(res, req)
	if res.Code != http.StatusCreated || !strings.Contains(res.Body.String(), `"approval_request_id":"APR-1"`) {
		t.Fatalf("status=%d body=%s", res.Code, res.Body.String())
	}
	if command["capability_code"] != "site.selection.recommend" || command["actor_id"] != "project-owner" {
		t.Fatalf("command=%v", command)
	}
	if command["site_selection_recommendation"].(map[string]any)["selected_proposal_id"] != "SITE-1" {
		t.Fatalf("recommendation=%v", command)
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
