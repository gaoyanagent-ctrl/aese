package replay

import (
	"context"
	"fmt"
	"testing"

	"github.com/industrial-ai/iaos-aese/internal/iaosclient"
	"github.com/industrial-ai/iaos-aese/internal/scenariopack"
)

type fakeIAOS struct {
	plans            map[string]iaosclient.UpsertPlan
	records          map[string][]map[string]any
	planCalls        int
	writeCalls       int
	decomposes       int
	simulations      int
	lastSimulation   iaosclient.SimulationEventRequest
	simulationResult iaosclient.SimulationEventResult
}

func (f *fakeIAOS) Schema(context.Context, string) (iaosclient.Schema, error) {
	schema := iaosclient.Schema{StorageStrategy: "TABLE"}
	schema.Permissions.Create = true
	schema.Permissions.Update = true
	return schema, nil
}
func (f *fakeIAOS) PlanUpsert(_ context.Context, req iaosclient.UpsertRequest) (iaosclient.UpsertPlan, error) {
	f.planCalls++
	key := fmt.Sprint(req.Record[req.NaturalKey[0]])
	if plan, ok := f.plans[key]; ok {
		return plan, nil
	}
	return iaosclient.UpsertPlan{Action: iaosclient.ActionCreate, Entity: req.Entity}, nil
}
func (f *fakeIAOS) Upsert(ctx context.Context, req iaosclient.UpsertRequest) (iaosclient.UpsertResult, error) {
	f.writeCalls++
	plan, err := f.PlanUpsert(ctx, req)
	return iaosclient.UpsertResult{UpsertPlan: plan, Applied: true}, err
}
func (f *fakeIAOS) FindExact(_ context.Context, entity string, _ map[string]any) ([]map[string]any, error) {
	return f.records[entity], nil
}
func (f *fakeIAOS) DecomposeSalesOrder(context.Context, string) (iaosclient.DecomposeResult, error) {
	f.decomposes++
	return iaosclient.DecomposeResult{Status: "confirmed", Decomposing: true}, nil
}
func (f *fakeIAOS) DecomposeSalesOrderTrace(ctx context.Context, id, correlation, idempotency string) (iaosclient.DecomposeResult, error) {
	return f.DecomposeSalesOrder(ctx, id)
}
func (f *fakeIAOS) IngestSimulationEvent(_ context.Context, request iaosclient.SimulationEventRequest) (iaosclient.SimulationEventResult, error) {
	f.simulations++
	f.lastSimulation = request
	if f.simulationResult.EventID != "" || f.simulationResult.Duplicate || f.simulationResult.Committed {
		return f.simulationResult, nil
	}
	var result iaosclient.SimulationEventResult
	result.EventID = "evt-sim-1"
	result.Subject = "iaos.tenant-hctm.eam.machine.down"
	result.CorrelationID = request.CorrelationID
	result.Committed = true
	result.BusinessObject.Type = request.BusinessObject.Type
	result.BusinessObject.Code = request.BusinessObject.Code
	result.BusinessObject.ID = "equipment-uuid"
	return result, nil
}

func TestApplyRecordSetsDryRunNeverWrites(t *testing.T) {
	fake := &fakeIAOS{plans: map[string]iaosclient.UpsertPlan{"CUST-1": {Action: iaosclient.ActionCreate, NaturalKey: map[string]any{"customer_code": "CUST-1"}}}}
	runner, _ := New(fake)
	sets := []scenariopack.RecordSet{{Entity: "customer", NaturalKey: []string{"customer_code"}, Records: []map[string]any{{"customer_code": "CUST-1"}}}}
	summary, err := runner.ApplyRecordSets(context.Background(), sets, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if summary.Mode != "dry-run" || summary.Created != 1 || fake.writeCalls != 0 {
		t.Fatalf("summary=%#v writes=%d", summary, fake.writeCalls)
	}
}

func TestApplyRecordSetsRequiresExplicitApply(t *testing.T) {
	fake := &fakeIAOS{plans: map[string]iaosclient.UpsertPlan{"CUST-1": {Action: iaosclient.ActionCreate, NaturalKey: map[string]any{"customer_code": "CUST-1"}}}}
	runner, _ := New(fake)
	sets := []scenariopack.RecordSet{{Entity: "customer", NaturalKey: []string{"customer_code"}, Records: []map[string]any{{"customer_code": "CUST-1"}}}}
	summary, err := runner.ApplyRecordSets(context.Background(), sets, Options{Apply: true})
	if err != nil {
		t.Fatal(err)
	}
	if summary.Mode != "apply" || fake.writeCalls != 1 || !summary.Impacts[0].Applied {
		t.Fatalf("summary=%#v writes=%d", summary, fake.writeCalls)
	}
}

func TestReplayDryRunResolvesButDoesNotTrigger(t *testing.T) {
	fake := &fakeIAOS{records: map[string][]map[string]any{"sales_order": {{"id": "so-id", "order_no": "SO-1"}}}}
	runner, _ := New(fake)
	story := scenariopack.Story{Ref: scenariopack.StoryRef{Key: "story"}, Events: scenariopack.EventSequence{Events: []scenariopack.Event{{EventID: "evt-1", EventType: "o2d.order.confirmed", Payload: map[string]any{"order_no": "SO-1"}}}}}
	summary, err := runner.Replay(context.Background(), story, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if summary.Impacts[0].Action != "decompose_sales_order" || fake.decomposes != 0 {
		t.Fatalf("summary=%#v decomposes=%d", summary, fake.decomposes)
	}
}

func TestReplayApplyUsesDecomposeAndSkipsUnsupported(t *testing.T) {
	fake := &fakeIAOS{records: map[string][]map[string]any{"sales_order": {{"id": "so-id", "order_no": "SO-1"}}}}
	runner, _ := New(fake)
	story := scenariopack.Story{Ref: scenariopack.StoryRef{Key: "story"}, Events: scenariopack.EventSequence{Events: []scenariopack.Event{
		{EventID: "evt-x", EventType: "eam.machine.down", Timestamp: "2026-07-08T09:30:00+08:00", Metadata: map[string]any{"business_object_type": "equipment", "business_object_id": "LAS-WLD-02", "idempotency_key": "machine-down-1", "correlation_id": "corr-1"}, Payload: map[string]any{"equipment_code": "LAS-WLD-02"}},
		{EventID: "evt-1", EventType: "o2d.order.confirmed", Payload: map[string]any{"order_no": "SO-1"}},
	}}}
	summary, err := runner.Replay(context.Background(), story, Options{Apply: true, PackKey: "hctm"})
	if err != nil {
		t.Fatal(err)
	}
	if summary.Triggered != 2 || summary.Skipped != 0 || fake.decomposes != 1 || fake.simulations != 1 {
		t.Fatalf("summary=%#v decomposes=%d simulations=%d", summary, fake.decomposes, fake.simulations)
	}
	if fake.lastSimulation.EventType != "eam.machine.down" || fake.lastSimulation.BusinessObject.Code != "LAS-WLD-02" || fake.lastSimulation.IdempotencyKey != "machine-down-1" {
		t.Fatalf("simulation request=%#v", fake.lastSimulation)
	}
	if fake.lastSimulation.Source != "aese:hctm/story" {
		t.Fatalf("simulation source=%q", fake.lastSimulation.Source)
	}
}

func TestReplayFailsClosedWhenSimulationResponseIsNotCommitted(t *testing.T) {
	fake := &fakeIAOS{}
	fake.simulationResult.EventID = "evt-sim-1"
	fake.simulationResult.BusinessObject.ID = "equipment-uuid"
	runner, _ := New(fake)
	story := scenariopack.Story{Ref: scenariopack.StoryRef{Key: "story"}, Events: scenariopack.EventSequence{Events: []scenariopack.Event{{
		EventID: "evt-x", EventType: "eam.machine.down", Timestamp: "2026-07-08T09:30:00+08:00",
		Metadata: map[string]any{"business_object_type": "equipment", "business_object_id": "LAS-WLD-02", "idempotency_key": "machine-down-1"},
	}}}}

	summary, err := runner.Replay(context.Background(), story, Options{Apply: true, PackKey: "hctm"})
	if err == nil || summary.Failed != 1 || summary.Triggered != 0 {
		t.Fatalf("summary=%#v err=%v", summary, err)
	}
}

func TestReplayUsesScenarioApplyOrderID(t *testing.T) {
	fake := &fakeIAOS{}
	runner, _ := New(fake)
	story := scenariopack.Story{Events: scenariopack.EventSequence{Events: []scenariopack.Event{{EventID: "evt-1", EventType: "o2d.order.confirmed", Payload: map[string]any{"order_no": "SO-1"}}}}}
	summary, err := runner.Replay(context.Background(), story, Options{OrderID: "so-uuid", Apply: true})
	if err != nil || summary.Triggered != 1 || summary.Impacts[0].RecordID != "so-uuid" || fake.decomposes != 1 {
		t.Fatalf("summary=%#v err=%v", summary, err)
	}
}

func TestVerifyMinimalAssertions(t *testing.T) {
	fake := &fakeIAOS{records: map[string][]map[string]any{"work_order": {{"id": "wo-1", "quantity": 10800, "status": "pending"}}}}
	runner, _ := New(fake)
	summary, err := runner.Verify(context.Background(), []Assertion{
		{Key: "wo-exists", Entity: "work_order", Match: map[string]any{"id": "wo-1"}, Operator: "exists"},
		{Key: "wo-qty", Entity: "work_order", Match: map[string]any{"id": "wo-1"}, Field: "quantity", Operator: "gte", Expected: 10000},
	}, Options{})
	if err != nil || summary.Passed != 2 {
		t.Fatalf("summary=%#v err=%v", summary, err)
	}
}
