package iaosclient

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) Do(r *http.Request) (*http.Response, error) { return f(r) }

func response(status int, body string) *http.Response {
	return &http.Response{StatusCode: status, Body: io.NopCloser(strings.NewReader(body)), Header: make(http.Header)}
}

func TestPlanUpsertIsReadOnlyAndExact(t *testing.T) {
	var methods []string
	client := testClient(t, func(r *http.Request) (*http.Response, error) {
		methods = append(methods, r.Method)
		if got := r.Header.Get("Authorization"); got != "Bearer secret" {
			t.Fatalf("authorization = %q", got)
		}
		return response(200, `{"total":2,"data":[{"id":"id-1","customer_code":"CUST-1","name":"old"},{"id":"id-2","customer_code":"CUST-10","name":"other"}]}`), nil
	})
	plan, err := client.PlanUpsert(context.Background(), UpsertRequest{
		Entity: "customer", NaturalKey: []string{"customer_code"},
		Record: map[string]any{"customer_code": "CUST-1", "name": "new"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if plan.Action != ActionUpdate || plan.RecordID != "id-1" {
		t.Fatalf("plan = %#v", plan)
	}
	if len(methods) != 1 || methods[0] != http.MethodGet {
		t.Fatalf("dry-run methods = %v", methods)
	}
}

func TestUpsertCreatesAndUpdates(t *testing.T) {
	tests := []struct {
		name       string
		listBody   string
		wantMethod string
		wantAction UpsertAction
	}{
		{"create", `{"total":0,"data":[]}`, http.MethodPost, ActionCreate},
		{"update", `{"total":1,"data":[{"id":"id-1","customer_code":"CUST-1","name":"old"}]}`, http.MethodPut, ActionUpdate},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			client := testClient(t, func(r *http.Request) (*http.Response, error) {
				calls++
				if calls == 1 {
					return response(200, tt.listBody), nil
				}
				if r.Method != tt.wantMethod {
					t.Fatalf("write method = %s", r.Method)
				}
				var body map[string]any
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Fatal(err)
				}
				if body["customer_code"] != "CUST-1" {
					t.Fatalf("body = %#v", body)
				}
				return response(map[bool]int{true: 201, false: 200}[r.Method == http.MethodPost], `{"id":"new-id"}`), nil
			})
			result, err := client.Upsert(context.Background(), UpsertRequest{Entity: "customer", NaturalKey: []string{"customer_code"}, Record: map[string]any{"customer_code": "CUST-1", "name": "new"}})
			if err != nil {
				t.Fatal(err)
			}
			if result.Action != tt.wantAction || !result.Applied || calls != 2 {
				t.Fatalf("result=%#v calls=%d", result, calls)
			}
		})
	}
}

func TestUpsertUnchangedDoesNotWrite(t *testing.T) {
	calls := 0
	client := testClient(t, func(r *http.Request) (*http.Response, error) {
		calls++
		return response(200, `{"total":1,"data":[{"id":"id-1","customer_code":"CUST-1","qty":12}]}`), nil
	})
	result, err := client.Upsert(context.Background(), UpsertRequest{Entity: "customer", NaturalKey: []string{"customer_code"}, Record: map[string]any{"customer_code": "CUST-1", "qty": 12}})
	if err != nil {
		t.Fatal(err)
	}
	if result.Action != ActionUnchanged || result.Applied || calls != 1 {
		t.Fatalf("result=%#v calls=%d", result, calls)
	}
}

func TestAPIErrorDoesNotExposeToken(t *testing.T) {
	client := testClient(t, func(r *http.Request) (*http.Response, error) {
		return response(403, `{"error":"forbidden"}`), nil
	})
	_, err := client.Schema(context.Background(), "customer")
	if err == nil || !IsStatus(err, 403) {
		t.Fatalf("err = %v", err)
	}
	if strings.Contains(err.Error(), "secret") {
		t.Fatalf("token leaked: %v", err)
	}
}

func TestDecomposeUsesGovernedEndpoint(t *testing.T) {
	client := testClient(t, func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/entities/sales_order/id-1/decompose" {
			t.Fatalf("request = %s %s", r.Method, r.URL.Path)
		}
		if r.Header.Get("X-Correlation-ID") != "corr-1" || r.Header.Get("Idempotency-Key") != "idem-1" {
			t.Fatalf("trace headers missing")
		}
		return response(200, `{"status":"confirmed","decomposing":true,"sales_order_id":"id-1","order_no":"SO-1"}`), nil
	})
	got, err := client.DecomposeSalesOrderTrace(context.Background(), "id-1", "corr-1", "idem-1")
	if err != nil || !got.Decomposing || got.OrderNo != "SO-1" {
		t.Fatalf("got=%#v err=%v", got, err)
	}
}

func TestSimulationIngressUsesGovernedEndpointAndWireContract(t *testing.T) {
	client := testClient(t, func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/simulation/events" {
			t.Fatalf("request = %s %s", r.Method, r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer secret" || r.Header.Get("X-IAOS-Tenant-ID") != "tenant-hctm" {
			t.Fatalf("governance headers missing: %#v", r.Header)
		}
		var body SimulationEventRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body.EventType != "eam.machine.down" || body.Source != "aese:hctm/story" || body.BusinessObject.Code != "LAS-WLD-02" || body.IdempotencyKey != "idem-1" {
			t.Fatalf("body = %#v", body)
		}
		return response(200, `{"event_id":"evt-1","subject":"iaos.tenant-hctm.eam.machine.down","correlation_id":"corr-1","committed":true,"business_object":{"type":"equipment","code":"LAS-WLD-02","id":"equipment-uuid"}}`), nil
	})

	got, err := client.IngestSimulationEvent(context.Background(), SimulationEventRequest{
		EventType: "eam.machine.down", Source: "aese:hctm/story", OccurredAt: "2026-07-08T09:30:00+08:00",
		CorrelationID: "corr-1", IdempotencyKey: "idem-1",
		BusinessObject: SimulationBusinessObject{Type: "equipment", Code: "LAS-WLD-02"},
		Payload:        map[string]any{"equipment_code": "LAS-WLD-02"},
	})
	if err != nil || !got.Committed || got.EventID != "evt-1" || got.BusinessObject.ID != "equipment-uuid" {
		t.Fatalf("got=%#v err=%v", got, err)
	}
}

func TestScenarioApplyUsesGovernedEndpointAndCorrelation(t *testing.T) {
	client := testClient(t, func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/scenarios/apply" {
			t.Fatalf("request=%s %s", r.Method, r.URL.Path)
		}
		if got := r.Header.Get("X-Correlation-ID"); got != "corr-hctm" {
			t.Fatalf("correlation=%q", got)
		}
		return response(200, `{"pack_key":"hctm","run_id":"run-1","dry_run":true,"committed":false,"inserted":18,"results":[]}`), nil
	})
	got, err := client.ApplyScenario(context.Background(), map[string]any{"pack_key": "hctm"}, "corr-hctm")
	if err != nil || got.Inserted != 18 || got.Committed {
		t.Fatalf("got=%#v err=%v", got, err)
	}
}

func TestScenarioResetUsesGovernedEndpoint(t *testing.T) {
	client := testClient(t, func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/api/v1/scenarios/reset" {
			t.Fatalf("path=%s", r.URL.Path)
		}
		return response(200, `{"pack_key":"hctm","deleted":2,"preserved_l1":12,"results":[]}`), nil
	})
	got, err := client.ResetScenario(context.Background(), ScenarioResetRequest{PackKey: "hctm", PackVersion: "0.1.0", ScenarioKey: "story", RunID: "reset-1", DryRun: true}, "corr")
	if err != nil || got.Deleted != 2 || got.PreservedL1 != 12 {
		t.Fatalf("got=%#v err=%v", got, err)
	}
}

func TestUpsertMetadataSchemaUsesGovernedEndpointAndWireContract(t *testing.T) {
	client := testClient(t, func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/metadata/schema/planning_context" {
			t.Fatalf("request = %s %s", r.Method, r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer secret" || r.Header.Get("X-IAOS-Tenant-ID") != "tenant-hctm" {
			t.Fatalf("governance headers missing: %#v", r.Header)
		}
		var body struct {
			DisplayName string           `json:"display_name"`
			Fields      []map[string]any `json:"fields"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body.DisplayName != "计划上下文" || len(body.Fields) != 1 || body.Fields[0]["name"] != "order_no" {
			t.Fatalf("body = %#v", body)
		}
		return response(http.StatusOK, `{}`), nil
	})

	err := client.UpsertMetadataSchema(context.Background(), "planning_context", MetadataSchemaRequest{
		DisplayName: "计划上下文",
		Fields:      json.RawMessage(`[{"name":"order_no","type":"string","required":true}]`),
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestAIToolLifecycleUsesRegistryEndpoints(t *testing.T) {
	manifest := AIToolManifest{
		ToolKey: "planning.context.read", DisplayName: "读取计划上下文",
		ToolType: "query", SourceRef: "entity.records",
		InputSchema: json.RawMessage(`{"type":"object"}`), RiskLevel: "low",
		ConfirmationMode: "none", PermissionResource: "planning.context.read",
		Examples: json.RawMessage(`[{"order_no":"SO-1"}]`), Metadata: json.RawMessage(`{"entity":"sales_order"}`),
	}
	want := []struct {
		method string
		path   string
		query  string
	}{
		{http.MethodGet, "/api/v1/ai/tools", "limit=200&include_disabled=true"},
		{http.MethodPost, "/api/v1/ai/tools", ""},
		{http.MethodPatch, "/api/v1/ai/tools/planning.context.read", ""},
		{http.MethodPost, "/api/v1/ai/tools/planning.context.read/enable", ""},
	}
	calls := 0
	client := testClient(t, func(r *http.Request) (*http.Response, error) {
		if calls >= len(want) {
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.String())
		}
		current := want[calls]
		calls++
		if r.Method != current.method || r.URL.Path != current.path || r.URL.RawQuery != current.query {
			t.Fatalf("request %d = %s %s?%s, want %s %s?%s", calls, r.Method, r.URL.Path, r.URL.RawQuery, current.method, current.path, current.query)
		}
		if calls == 1 {
			return response(http.StatusOK, `{"items":[{"tool_key":"planning.context.read"}]}`), nil
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if calls == 2 && (body["tool_key"] != manifest.ToolKey || body["source_ref"] != manifest.SourceRef) {
			t.Fatalf("create body = %#v", body)
		}
		if calls == 3 {
			if body["display_name"] != manifest.DisplayName || body["risk_level"] != "low" {
				t.Fatalf("patch body = %#v", body)
			}
			if _, exists := body["source_ref"]; exists {
				t.Fatalf("patch unexpectedly changes immutable source_ref: %#v", body)
			}
		}
		return response(http.StatusOK, `{}`), nil
	})

	tools, err := client.ListAITools(context.Background(), true)
	if err != nil || len(tools) != 1 || tools[0].ToolKey != manifest.ToolKey {
		t.Fatalf("tools=%#v err=%v", tools, err)
	}
	if err := client.CreateAITool(context.Background(), manifest); err != nil {
		t.Fatal(err)
	}
	if err := client.UpdateAITool(context.Background(), manifest); err != nil {
		t.Fatal(err)
	}
	if err := client.EnableAITool(context.Background(), manifest.ToolKey); err != nil {
		t.Fatal(err)
	}
	if calls != len(want) {
		t.Fatalf("calls = %d, want %d", calls, len(want))
	}
}

func TestCallAIToolUsesGovernedEndpointAndReturnsTrace(t *testing.T) {
	client := testClient(t, func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/ai/tools/planning.context.read/call" {
			t.Fatalf("request = %s %s", r.Method, r.URL.Path)
		}
		if r.Header.Get("X-Correlation-ID") != "corr-plan-1" {
			t.Fatalf("correlation = %q", r.Header.Get("X-Correlation-ID"))
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		input, _ := body["input"].(map[string]any)
		if body["session_id"] != "session-plan-1" || input["order_no"] != "SO-202607-0001" {
			t.Fatalf("body = %#v", body)
		}
		return response(http.StatusOK, `{"call_id":"call-1","status":"succeeded","output":{"total_demand":12000},"execution_ref":{"source_ref":"entity.records"}}`), nil
	})

	got, err := client.CallAITool(context.Background(), "planning.context.read", "corr-plan-1", "session-plan-1", map[string]any{"order_no": "SO-202607-0001"})
	if err != nil || got.CallID != "call-1" || got.Status != "succeeded" || !strings.Contains(string(got.Output), `"total_demand":12000`) {
		t.Fatalf("got=%#v err=%v", got, err)
	}
}

func TestCallAIToolFailsClosedOnUnsuccessfulResponse(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantStatus int
	}{
		{name: "incomplete 200", statusCode: http.StatusOK, body: `{"call_id":"call-1","status":"approval_required"}`},
		{name: "HTTP conflict", statusCode: http.StatusConflict, body: `{"error":"approval_required"}`, wantStatus: http.StatusConflict},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := testClient(t, func(r *http.Request) (*http.Response, error) {
				return response(tt.statusCode, tt.body), nil
			})
			_, err := client.CallAITool(context.Background(), "planning.context.read", "corr", "session", map[string]any{})
			if err == nil {
				t.Fatal("CallAITool returned nil error")
			}
			if tt.wantStatus != 0 && !IsStatus(err, tt.wantStatus) {
				t.Fatalf("err = %v, want HTTP status %d", err, tt.wantStatus)
			}
			if strings.Contains(err.Error(), "secret") {
				t.Fatalf("token leaked: %v", err)
			}
		})
	}
}

func testClient(t *testing.T, fn roundTripFunc) *Client {
	t.Helper()
	client, err := New(Config{BaseURL: "http://iaos.test", Token: "secret", TenantID: "tenant-hctm", HTTP: fn})
	if err != nil {
		t.Fatal(err)
	}
	return client
}
