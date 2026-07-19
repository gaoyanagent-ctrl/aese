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

func testClient(t *testing.T, fn roundTripFunc) *Client {
	t.Helper()
	client, err := New(Config{BaseURL: "http://iaos.test", Token: "secret", TenantID: "tenant-hctm", HTTP: fn})
	if err != nil {
		t.Fatal(err)
	}
	return client
}
