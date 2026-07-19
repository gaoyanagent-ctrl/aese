// Package replay coordinates dry-run/apply, the governed O2D tracer, and
// read-only outcome verification. It contains no direct database or NATS path.
package replay

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/industrial-ai/iaos-aese/internal/iaosclient"
	"github.com/industrial-ai/iaos-aese/internal/scenariopack"
)

type IAOS interface {
	Schema(context.Context, string) (iaosclient.Schema, error)
	PlanUpsert(context.Context, iaosclient.UpsertRequest) (iaosclient.UpsertPlan, error)
	Upsert(context.Context, iaosclient.UpsertRequest) (iaosclient.UpsertResult, error)
	FindExact(context.Context, string, map[string]any) ([]map[string]any, error)
	DecomposeSalesOrder(context.Context, string) (iaosclient.DecomposeResult, error)
	DecomposeSalesOrderTrace(context.Context, string, string, string) (iaosclient.DecomposeResult, error)
	IngestSimulationEvent(context.Context, iaosclient.SimulationEventRequest) (iaosclient.SimulationEventResult, error)
}

type Runner struct {
	client IAOS
	now    func() time.Time
}

func New(client IAOS) (*Runner, error) {
	if client == nil {
		return nil, fmt.Errorf("IAOS client is required")
	}
	return &Runner{client: client, now: time.Now}, nil
}

type Options struct {
	// Apply must be explicitly true before any external write is attempted.
	Apply  bool
	Target string
	Tenant string
	Actor  string
	// PackKey identifies the scenario package in IAOS audit records.
	PackKey string
	// OrderID is the governed UUID returned by scenario apply. It avoids
	// requiring a legacy sales_order metadata schema only for ID resolution.
	OrderID string
	// Entities limits imported record sets. Empty means all supplied sets.
	Entities map[string]bool
}

type Impact struct {
	Entity     string                  `json:"entity"`
	NaturalKey map[string]any          `json:"natural_key"`
	Action     iaosclient.UpsertAction `json:"action"`
	RecordID   string                  `json:"record_id,omitempty"`
	Changed    []string                `json:"changed_fields,omitempty"`
	Applied    bool                    `json:"applied"`
	Error      string                  `json:"error,omitempty"`
}

type RunSummary struct {
	RunID      string    `json:"run_id"`
	Mode       string    `json:"mode"`
	Target     string    `json:"target,omitempty"`
	Tenant     string    `json:"tenant,omitempty"`
	Actor      string    `json:"actor,omitempty"`
	StartedAt  time.Time `json:"started_at"`
	FinishedAt time.Time `json:"finished_at"`
	Planned    int       `json:"planned"`
	Created    int       `json:"created"`
	Updated    int       `json:"updated"`
	Unchanged  int       `json:"unchanged"`
	Failed     int       `json:"failed"`
	Impacts    []Impact  `json:"impacts"`
}

// ApplyRecordSets produces the same impact shape for dry-run and apply. It
// stops after the first write failure so a failure cannot silently fan out.
// IAOS currently exposes per-record CRUD, not an atomic batch endpoint, so the
// summary deliberately preserves any successfully applied prefix.
func (r *Runner) ApplyRecordSets(ctx context.Context, sets []scenariopack.RecordSet, opts Options) (RunSummary, error) {
	summary := r.start(opts, map[bool]string{true: "apply", false: "dry-run"}[opts.Apply])

	for _, set := range sets {
		if len(opts.Entities) > 0 && !opts.Entities[set.Entity] {
			continue
		}
		schema, err := r.client.Schema(ctx, set.Entity)
		if err != nil {
			summary.Failed++
			summary.Impacts = append(summary.Impacts, Impact{Entity: set.Entity, Error: err.Error()})
			summary.FinishedAt = r.now().UTC()
			return summary, fmt.Errorf("load IAOS schema for %s: %w", set.Entity, err)
		}
		if strings.EqualFold(schema.StorageStrategy, "EVENT") {
			err := fmt.Errorf("entity %s uses transient EVENT storage and cannot be upserted", set.Entity)
			summary.Failed++
			summary.Impacts = append(summary.Impacts, Impact{Entity: set.Entity, Error: err.Error()})
			summary.FinishedAt = r.now().UTC()
			return summary, err
		}
		for _, record := range set.Records {
			req := iaosclient.UpsertRequest{Entity: set.Entity, NaturalKey: set.NaturalKey, Record: record}
			plan, err := r.client.PlanUpsert(ctx, req)
			if err != nil {
				impact := Impact{Entity: set.Entity, NaturalKey: naturalKey(set.NaturalKey, record), Error: err.Error()}
				summary.Failed++
				summary.Impacts = append(summary.Impacts, impact)
				summary.FinishedAt = r.now().UTC()
				return summary, fmt.Errorf("plan %s upsert: %w", set.Entity, err)
			}
			impact := Impact{Entity: set.Entity, NaturalKey: plan.NaturalKey, Action: plan.Action, RecordID: plan.RecordID, Changed: plan.Changed}
			summary.Planned++
			if (plan.Action == iaosclient.ActionCreate && !schema.Permissions.Create) || (plan.Action == iaosclient.ActionUpdate && !schema.Permissions.Update) {
				impact.Error = "authenticated actor lacks required IAOS entity permission"
				summary.Failed++
				summary.Impacts = append(summary.Impacts, impact)
				summary.FinishedAt = r.now().UTC()
				return summary, fmt.Errorf("%s %s: %s", plan.Action, set.Entity, impact.Error)
			}
			if opts.Apply && plan.Action != iaosclient.ActionUnchanged {
				result, err := r.client.Upsert(ctx, req)
				if err != nil {
					impact.Error = err.Error()
					summary.Failed++
					summary.Impacts = append(summary.Impacts, impact)
					summary.FinishedAt = r.now().UTC()
					return summary, fmt.Errorf("apply %s upsert: %w", set.Entity, err)
				}
				impact.Applied = result.Applied
				impact.RecordID = result.RecordID
			}
			summary.Impacts = append(summary.Impacts, impact)
			switch impact.Action {
			case iaosclient.ActionCreate:
				summary.Created++
			case iaosclient.ActionUpdate:
				summary.Updated++
			case iaosclient.ActionUnchanged:
				summary.Unchanged++
			}
		}
	}
	summary.FinishedAt = r.now().UTC()
	return summary, nil
}

type ReplayImpact struct {
	EventID       string `json:"event_id"`
	EventType     string `json:"event_type"`
	CorrelationID string `json:"correlation_id,omitempty"`
	Action        string `json:"action"`
	RecordID      string `json:"record_id,omitempty"`
	Applied       bool   `json:"applied"`
	Error         string `json:"error,omitempty"`
}

type ReplaySummary struct {
	RunID      string         `json:"run_id"`
	Mode       string         `json:"mode"`
	Target     string         `json:"target,omitempty"`
	Tenant     string         `json:"tenant,omitempty"`
	Actor      string         `json:"actor,omitempty"`
	StoryKey   string         `json:"story_key"`
	StartedAt  time.Time      `json:"started_at"`
	FinishedAt time.Time      `json:"finished_at"`
	Triggered  int            `json:"triggered"`
	Skipped    int            `json:"skipped"`
	Failed     int            `json:"failed"`
	Impacts    []ReplayImpact `json:"impacts"`
}

var simulationBusinessObjectTypes = map[string]string{
	"o2d.supplier_delivery.delayed":  "purchase_order",
	"eam.machine.down":               "equipment",
	"qms.incoming_inspection.failed": "inspection_order",
}

// Replay translates supported canonical events to governed IAOS business and
// simulation ingress endpoints. Unsupported events remain visible and are
// never published directly to NATS.
func (r *Runner) Replay(ctx context.Context, story scenariopack.Story, opts Options) (ReplaySummary, error) {
	started := r.now().UTC()
	summary := ReplaySummary{RunID: newRunID(), Mode: map[bool]string{true: "apply", false: "dry-run"}[opts.Apply], Target: sanitizedTarget(opts.Target), Tenant: opts.Tenant, Actor: opts.Actor, StoryKey: story.Ref.Key, StartedAt: started, Impacts: []ReplayImpact{}}
	for _, event := range story.Events.Events {
		impact := ReplayImpact{EventID: event.EventID, EventType: event.EventType, CorrelationID: event.Correlation(), Action: "unsupported"}
		if objectType, supported := simulationBusinessObjectTypes[event.EventType]; supported {
			request, err := simulationEventRequest(event, objectType, opts.PackKey, story.Ref.Key)
			if err != nil {
				impact.Error = err.Error()
				summary.Failed++
				summary.Impacts = append(summary.Impacts, impact)
				summary.FinishedAt = r.now().UTC()
				return summary, fmt.Errorf("event %s: %s", event.EventID, impact.Error)
			}
			impact.Action, impact.RecordID = "simulation_ingress", request.BusinessObject.Code
			if opts.Apply {
				result, err := r.client.IngestSimulationEvent(ctx, request)
				if err != nil {
					impact.Error = err.Error()
					summary.Failed++
					summary.Impacts = append(summary.Impacts, impact)
					summary.FinishedAt = r.now().UTC()
					return summary, fmt.Errorf("event %s: %w", event.EventID, err)
				}
				if err := validateSimulationResult(request, result, opts.Tenant); err != nil {
					impact.Error = err.Error()
					summary.Failed++
					summary.Impacts = append(summary.Impacts, impact)
					summary.FinishedAt = r.now().UTC()
					return summary, fmt.Errorf("event %s: %w", event.EventID, err)
				}
				impact.RecordID = result.BusinessObject.ID
				if result.Duplicate {
					impact.Action = "duplicate"
					summary.Skipped++
				} else {
					impact.Applied = result.Committed
					summary.Triggered++
				}
			}
			summary.Impacts = append(summary.Impacts, impact)
			continue
		}
		if event.EventType != "o2d.order.confirmed" {
			summary.Skipped++
			summary.Impacts = append(summary.Impacts, impact)
			continue
		}
		orderNo := stringField(event.Payload, "order_no")
		if orderNo == "" && event.AggregateType == "sales_order" {
			orderNo = event.AggregateID
		}
		if orderNo == "" {
			impact.Error = "o2d.order.confirmed is missing order_no"
			summary.Failed++
			summary.Impacts = append(summary.Impacts, impact)
			summary.FinishedAt = r.now().UTC()
			return summary, fmt.Errorf("event %s: %s", event.EventID, impact.Error)
		}
		id := strings.TrimSpace(opts.OrderID)
		if id == "" {
			records, err := r.client.FindExact(ctx, "sales_order", map[string]any{"order_no": orderNo})
			if err != nil || len(records) != 1 {
				if err == nil {
					err = fmt.Errorf("expected one sales_order for order_no %s, got %d", orderNo, len(records))
				}
				impact.Error = err.Error()
				summary.Failed++
				summary.Impacts = append(summary.Impacts, impact)
				summary.FinishedAt = r.now().UTC()
				return summary, fmt.Errorf("event %s: %w", event.EventID, err)
			}
			id, _ = records[0]["id"].(string)
		}
		if id == "" {
			impact.Error = "matched sales_order has no id"
			summary.Failed++
			summary.Impacts = append(summary.Impacts, impact)
			summary.FinishedAt = r.now().UTC()
			return summary, fmt.Errorf("event %s: %s", event.EventID, impact.Error)
		}
		impact.Action, impact.RecordID = "decompose_sales_order", id
		if opts.Apply {
			result, err := r.client.DecomposeSalesOrderTrace(ctx, id, event.Correlation(), event.Idempotency())
			if err != nil {
				impact.Error = err.Error()
				summary.Failed++
				summary.Impacts = append(summary.Impacts, impact)
				summary.FinishedAt = r.now().UTC()
				return summary, fmt.Errorf("event %s: %w", event.EventID, err)
			}
			if result.Decomposing {
				impact.Applied = true
				summary.Triggered++
			} else {
				impact.Action = result.Status
				summary.Skipped++
			}
		}
		summary.Impacts = append(summary.Impacts, impact)
	}
	summary.FinishedAt = r.now().UTC()
	return summary, nil
}

func simulationEventRequest(event scenariopack.Event, expectedObjectType, packKey, storyKey string) (iaosclient.SimulationEventRequest, error) {
	objectType := stringField(event.Metadata, "business_object_type")
	objectCode := stringField(event.Metadata, "business_object_id")
	if objectType != expectedObjectType || objectCode == "" {
		return iaosclient.SimulationEventRequest{}, fmt.Errorf("%s requires %s business object metadata", event.EventType, expectedObjectType)
	}
	packKey = strings.TrimSpace(packKey)
	if packKey == "" {
		packKey = "unknown-pack"
	}
	return iaosclient.SimulationEventRequest{
		EventType: event.EventType, Source: "aese:" + packKey + "/" + storyKey,
		OccurredAt: event.Timestamp, CorrelationID: event.Correlation(), CausationID: event.Causation(),
		IdempotencyKey: event.Idempotency(),
		BusinessObject: iaosclient.SimulationBusinessObject{Type: objectType, Code: objectCode}, Payload: event.Payload,
	}, nil
}

func validateSimulationResult(request iaosclient.SimulationEventRequest, result iaosclient.SimulationEventResult, expectedTenant string) error {
	if strings.TrimSpace(result.EventID) == "" || strings.TrimSpace(result.Subject) == "" || strings.TrimSpace(result.BusinessObject.ID) == "" {
		return fmt.Errorf("simulation ingress returned an incomplete success response")
	}
	expectedTenant = strings.TrimSpace(expectedTenant)
	if expectedTenant != "" {
		expectedSubject := "iaos." + expectedTenant + "." + request.EventType
		if result.Subject != expectedSubject {
			return fmt.Errorf("simulation ingress returned a different event subject")
		}
	} else if !validSimulationSubject(result.Subject, request.EventType) {
		return fmt.Errorf("simulation ingress returned a different event subject")
	}
	if result.BusinessObject.Type != request.BusinessObject.Type || result.BusinessObject.Code != request.BusinessObject.Code {
		return fmt.Errorf("simulation ingress returned a different business object")
	}
	if request.CorrelationID != "" && result.CorrelationID != request.CorrelationID {
		return fmt.Errorf("simulation ingress returned a different correlation_id")
	}
	if !result.Duplicate && !result.Committed {
		return fmt.Errorf("simulation ingress did not commit the event")
	}
	return nil
}

func validSimulationSubject(subject, eventType string) bool {
	suffix := "." + eventType
	if !strings.HasPrefix(subject, "iaos.") || !strings.HasSuffix(subject, suffix) {
		return false
	}
	tenant := strings.TrimSuffix(strings.TrimPrefix(subject, "iaos."), suffix)
	return tenant != "" && !strings.Contains(tenant, ".")
}

type Assertion struct {
	Key      string         `json:"key"`
	Entity   string         `json:"entity"`
	Match    map[string]any `json:"match"`
	Field    string         `json:"field,omitempty"`
	Operator string         `json:"operator"`
	Expected any            `json:"expected,omitempty"`
}

type AssertionResult struct {
	Key      string `json:"key"`
	Passed   bool   `json:"passed"`
	Actual   any    `json:"actual,omitempty"`
	Expected any    `json:"expected,omitempty"`
	Error    string `json:"error,omitempty"`
}

type VerifySummary struct {
	RunID      string            `json:"run_id"`
	Target     string            `json:"target,omitempty"`
	Tenant     string            `json:"tenant,omitempty"`
	StartedAt  time.Time         `json:"started_at"`
	FinishedAt time.Time         `json:"finished_at"`
	Passed     int               `json:"passed"`
	Failed     int               `json:"failed"`
	Assertions []AssertionResult `json:"assertions"`
}

// Verify is read-only and implements the minimum M3 assertion vocabulary:
// count_eq, exists, eq, gte, and lte.
func (r *Runner) Verify(ctx context.Context, assertions []Assertion, opts Options) (VerifySummary, error) {
	summary := VerifySummary{RunID: newRunID(), Target: sanitizedTarget(opts.Target), Tenant: opts.Tenant, StartedAt: r.now().UTC(), Assertions: []AssertionResult{}}
	for _, assertion := range assertions {
		result := AssertionResult{Key: assertion.Key, Expected: assertion.Expected}
		records, err := r.client.FindExact(ctx, assertion.Entity, assertion.Match)
		if err != nil {
			result.Error = err.Error()
		} else {
			result.Passed, result.Actual, result.Error = evaluate(assertion, records)
		}
		if result.Passed {
			summary.Passed++
		} else {
			summary.Failed++
		}
		summary.Assertions = append(summary.Assertions, result)
	}
	summary.FinishedAt = r.now().UTC()
	if summary.Failed > 0 {
		return summary, fmt.Errorf("%d of %d assertions failed", summary.Failed, len(assertions))
	}
	return summary, nil
}

func evaluate(assertion Assertion, records []map[string]any) (bool, any, string) {
	switch assertion.Operator {
	case "count_eq":
		wanted, ok := number(assertion.Expected)
		if !ok {
			return false, len(records), "count_eq expected must be numeric"
		}
		return float64(len(records)) == wanted, len(records), ""
	case "exists":
		return len(records) > 0, len(records) > 0, ""
	case "eq", "gte", "lte":
		if len(records) != 1 {
			return false, len(records), fmt.Sprintf("expected exactly one matching record, got %d", len(records))
		}
		actual, ok := records[0][assertion.Field]
		if !ok {
			return false, nil, fmt.Sprintf("field %q is absent", assertion.Field)
		}
		if assertion.Operator == "eq" {
			return fmt.Sprint(actual) == fmt.Sprint(assertion.Expected), actual, ""
		}
		av, aok := number(actual)
		want, wok := number(assertion.Expected)
		if !aok || !wok {
			return false, actual, assertion.Operator + " requires numeric values"
		}
		if assertion.Operator == "gte" {
			return av >= want, actual, ""
		}
		return av <= want, actual, ""
	default:
		return false, nil, fmt.Sprintf("unsupported operator %q", assertion.Operator)
	}
}

func (r *Runner) start(opts Options, mode string) RunSummary {
	return RunSummary{RunID: newRunID(), Mode: mode, Target: sanitizedTarget(opts.Target), Tenant: opts.Tenant, Actor: opts.Actor, StartedAt: r.now().UTC(), Impacts: []Impact{}}
}

// RecordSets returns the L1 sets plus one story's typed L2 sets in stable pack
// order. Untyped InitialState.Records cannot be safely imported and are not
// guessed here.
func RecordSets(pack *scenariopack.Pack, storyKey string) ([]scenariopack.RecordSet, error) {
	if pack == nil {
		return nil, fmt.Errorf("scenario pack is required")
	}
	sets := append([]scenariopack.RecordSet(nil), pack.RecordSets...)
	if storyKey == "" {
		return sets, nil
	}
	for _, story := range pack.Stories {
		if story.Ref.Key == storyKey || story.Initial.StoryKey == storyKey {
			return append(sets, story.Initial.RecordSets...), nil
		}
	}
	return nil, fmt.Errorf("story %q not found", storyKey)
}

func naturalKey(fields []string, record map[string]any) map[string]any {
	key := make(map[string]any, len(fields))
	for _, field := range fields {
		key[field] = record[field]
	}
	return key
}

func stringField(values map[string]any, key string) string {
	value, _ := values[key].(string)
	return value
}

func number(value any) (float64, bool) {
	switch v := value.(type) {
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	case float64:
		return v, true
	case fmt.Stringer:
		var out float64
		_, err := fmt.Sscan(v.String(), &out)
		return out, err == nil
	default:
		var out float64
		_, err := fmt.Sscan(fmt.Sprint(value), &out)
		return out, err == nil
	}
}

func newRunID() string {
	var suffix [8]byte
	if _, err := rand.Read(suffix[:]); err != nil {
		return fmt.Sprintf("run-%d", time.Now().UnixNano())
	}
	return "run-" + hex.EncodeToString(suffix[:])
}

func sanitizedTarget(raw string) string {
	parsed, err := url.Parse(raw)
	if err != nil {
		return ""
	}
	parsed.User = nil
	parsed.RawQuery = ""
	parsed.Fragment = ""
	return parsed.String()
}

// SortedEntityAllowlist is a presentation helper for stable impact output.
func SortedEntityAllowlist(values map[string]bool) []string {
	var out []string
	for value, enabled := range values {
		if enabled {
			out = append(out, value)
		}
	}
	sort.Strings(out)
	return out
}
