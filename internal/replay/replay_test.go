package replay

import (
	"context"
	"fmt"
	"path/filepath"
	"reflect"
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
	simulationCalls  []iaosclient.SimulationEventRequest
	simulationResult iaosclient.SimulationEventResult
	businessCalls    []iaosclient.ScenarioBusinessEventRequest
	businessResult   iaosclient.ScenarioBusinessEventResult
}

func (f *fakeIAOS) IngestScenarioBusinessEvent(_ context.Context, _, _ string, request iaosclient.ScenarioBusinessEventRequest) (iaosclient.ScenarioBusinessEventResult, error) {
	f.businessCalls = append(f.businessCalls, request)
	if f.businessResult.EventID != "" || f.businessResult.Duplicate || f.businessResult.Committed {
		return f.businessResult, nil
	}
	var result iaosclient.ScenarioBusinessEventResult
	result.EventID = request.EventID
	result.EventType = request.EventType
	result.Cursor = int64(len(f.businessCalls))
	result.Subject = "iaos.tenant-hctm." + request.EventType
	result.CorrelationID = request.CorrelationID
	result.Committed = true
	result.BusinessObject.Type = request.BusinessObject.Type
	result.BusinessObject.Code = request.BusinessObject.Code
	result.BusinessObject.ID = "business-object-uuid"
	return result, nil
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
	f.simulationCalls = append(f.simulationCalls, request)
	if f.simulationResult.EventID != "" || f.simulationResult.Duplicate || f.simulationResult.Committed {
		return f.simulationResult, nil
	}
	var result iaosclient.SimulationEventResult
	result.EventID = "evt-sim-1"
	result.Subject = "iaos.tenant-hctm." + request.EventType
	result.CorrelationID = request.CorrelationID
	result.Committed = true
	result.BusinessObject.Type = request.BusinessObject.Type
	result.BusinessObject.Code = request.BusinessObject.Code
	result.BusinessObject.ID = "equipment-uuid"
	return result, nil
}

type stagedFakeIAOS struct {
	*fakeIAOS
	scenarioBusinessResults []iaosclient.ScenarioBusinessEventResult
	simulationResults       []iaosclient.SimulationEventResult
	scenarioBusinessIndex   int
	simulationIndex         int
}

func (f *stagedFakeIAOS) IngestScenarioBusinessEvent(_ context.Context, pack, story string, request iaosclient.ScenarioBusinessEventRequest) (iaosclient.ScenarioBusinessEventResult, error) {
	if f.scenarioBusinessIndex < len(f.scenarioBusinessResults) {
		result := f.scenarioBusinessResults[f.scenarioBusinessIndex]
		f.scenarioBusinessIndex++
		return result, nil
	}
	return f.fakeIAOS.IngestScenarioBusinessEvent(context.Background(), pack, story, request)
}

func (f *stagedFakeIAOS) IngestSimulationEvent(_ context.Context, request iaosclient.SimulationEventRequest) (iaosclient.SimulationEventResult, error) {
	if f.simulationIndex < len(f.simulationResults) {
		result := f.simulationResults[f.simulationIndex]
		f.simulationIndex++
		return result, nil
	}
	return f.fakeIAOS.IngestSimulationEvent(context.Background(), request)
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

func TestReplayRoutesCanonicalSupplierAndQualityEventsToSimulationIngress(t *testing.T) {
	fake := &fakeIAOS{}
	runner, _ := New(fake)
	supplierPayload := map[string]any{"po_no": "PO-202607-0001", "delay_days": 3}
	qualityPayload := map[string]any{"inspection_no": "IQC-202607-0002", "rejected_qty": 300}
	story := scenariopack.Story{Ref: scenariopack.StoryRef{Key: "order-expedite-01"}, Events: scenariopack.EventSequence{Events: []scenariopack.Event{
		{EventID: "evt-delay", EventType: "o2d.supplier_delivery.delayed", Timestamp: "2026-07-08T09:00:00+08:00", Metadata: map[string]any{"business_object_type": "purchase_order", "business_object_id": "PO-202607-0001", "idempotency_key": "delay-1", "correlation_id": "corr-1"}, Payload: supplierPayload},
		{EventID: "evt-iqc", EventType: "qms.incoming_inspection.failed", Timestamp: "2026-07-12T15:00:00+08:00", Metadata: map[string]any{"business_object_type": "inspection_order", "business_object_id": "IQC-202607-0002", "idempotency_key": "iqc-1", "correlation_id": "corr-1"}, Payload: qualityPayload},
	}}}

	summary, err := runner.Replay(context.Background(), story, Options{Apply: true, PackKey: "hctm"})
	if err != nil {
		t.Fatal(err)
	}
	if summary.Triggered != 2 || fake.simulations != 2 {
		t.Fatalf("summary=%#v simulations=%d", summary, fake.simulations)
	}
	want := []struct {
		eventType  string
		objectType string
		objectCode string
		payload    map[string]any
	}{
		{"o2d.supplier_delivery.delayed", "purchase_order", "PO-202607-0001", supplierPayload},
		{"qms.incoming_inspection.failed", "inspection_order", "IQC-202607-0002", qualityPayload},
	}
	for index, request := range fake.simulationCalls {
		if request.EventType != want[index].eventType || request.BusinessObject.Type != want[index].objectType || request.BusinessObject.Code != want[index].objectCode {
			t.Fatalf("request[%d]=%#v", index, request)
		}
		if request.Source != "aese:hctm/order-expedite-01" || !reflect.DeepEqual(request.Payload, want[index].payload) {
			t.Fatalf("request[%d] source/payload=%q %#v", index, request.Source, request.Payload)
		}
	}
}

func TestReplaySimulationDryRunPlansAllGovernedEventsWithoutWrites(t *testing.T) {
	fake := &fakeIAOS{}
	runner, _ := New(fake)
	story := scenariopack.Story{Ref: scenariopack.StoryRef{Key: "story"}, Events: scenariopack.EventSequence{Events: []scenariopack.Event{
		{EventID: "evt-delay", EventType: "o2d.supplier_delivery.delayed", Metadata: map[string]any{"business_object_type": "purchase_order", "business_object_id": "PO-1"}},
		{EventID: "evt-down", EventType: "eam.machine.down", Metadata: map[string]any{"business_object_type": "equipment", "business_object_id": "EQ-1"}},
		{EventID: "evt-iqc", EventType: "qms.incoming_inspection.failed", Metadata: map[string]any{"business_object_type": "inspection_order", "business_object_id": "IQC-1"}},
	}}}

	summary, err := runner.Replay(context.Background(), story, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if summary.Mode != "dry-run" || fake.simulations != 0 || summary.Triggered != 0 || summary.Skipped != 0 {
		t.Fatalf("summary=%#v simulations=%d", summary, fake.simulations)
	}
	for _, impact := range summary.Impacts {
		if impact.Action != "simulation_ingress" {
			t.Fatalf("impact=%#v", impact)
		}
	}
}

func TestReplayImpactIncludesRecoveryMetadataForGovernedActions(t *testing.T) {
	businessObject := func(objectType, code, id string) struct {
		Type string `json:"type"`
		Code string `json:"code"`
		ID   string `json:"id"`
	} {
		return struct {
			Type string `json:"type"`
			Code string `json:"code"`
			ID   string `json:"id"`
		}{Type: objectType, Code: code, ID: id}
	}
	fake := &stagedFakeIAOS{
		fakeIAOS: &fakeIAOS{},
		scenarioBusinessResults: []iaosclient.ScenarioBusinessEventResult{
			{
				EventID:       "evt-wo",
				EventType:     "proc.production_order.released",
				Subject:       "iaos.tenant-hctm.proc.production_order.released",
				CorrelationID: "corr-m7-o2",
				BusinessObject: businessObject("production_order", "PO-202607-0001", "po-uuid"),
				Cursor:       401,
				OperationRef: "op-release-202607-0001",
				Committed:    true,
			},
		},
		simulationResults: []iaosclient.SimulationEventResult{
			{
				EventID:        "evt-delay",
				Subject:        "iaos.tenant-hctm.o2d.supplier_delivery.delayed",
				CorrelationID:  "corr-m7-o2",
				Duplicate:      true,
				Cursor:         402,
				OperationRef:   "op-delay-dup",
				Committed:      false,
				BusinessObject: businessObject("purchase_order", "PO-202607-0001", "po-uuid"),
			},
			{
				EventID:        "evt-down",
				Subject:        "iaos.tenant-hctm.eam.machine.down",
				CorrelationID:  "corr-m7-o2",
				Cursor:         403,
				OperationRef:   "op-down-202607-0001",
				Committed:      true,
				BusinessObject: businessObject("equipment", "LAS-WLD-02", "eq-uuid"),
			},
		},
	}
	runner, _ := New(fake)
	story := scenariopack.Story{Ref: scenariopack.StoryRef{Key: "order-expedite-01"}, Events: scenariopack.EventSequence{Events: []scenariopack.Event{
		{EventID: "evt-wo", EventType: "proc.production_order.released", Metadata: map[string]any{"business_object_type": "production_order", "business_object_id": "PO-202607-0001", "idempotency_key": "release-1", "correlation_id": "corr-m7-o2"}},
		{EventID: "evt-delay", EventType: "o2d.supplier_delivery.delayed", Metadata: map[string]any{"business_object_type": "purchase_order", "business_object_id": "PO-202607-0001", "idempotency_key": "delay-1", "correlation_id": "corr-m7-o2"}},
		{EventID: "evt-down", EventType: "eam.machine.down", Metadata: map[string]any{"business_object_type": "equipment", "business_object_id": "LAS-WLD-02", "idempotency_key": "down-1", "correlation_id": "corr-m7-o2"}},
	}}}

	summary, err := runner.Replay(context.Background(), story, Options{Apply: true, PackKey: "hctm", Tenant: "tenant-hctm"})
	if err != nil || summary.Triggered != 2 || summary.Skipped != 1 || summary.Failed != 0 || len(summary.Impacts) != 3 {
		t.Fatalf("summary=%#v err=%v", summary, err)
	}

	releaseImpact := summary.Impacts[0]
	if releaseImpact.Action != "scenario_business_event" || releaseImpact.RecordID != "po-uuid" || releaseImpact.OperationRef != "op-release-202607-0001" {
		t.Fatalf("unexpected release impact=%#v", releaseImpact)
	}
	if releaseImpact.CorrelationID != "corr-m7-o2" || releaseImpact.Cursor != 401 || !releaseImpact.Committed {
		t.Fatalf("release impact should carry recovery metadata: %#v", releaseImpact)
	}
	if releaseImpact.NoOp {
		t.Fatal("release impact should not be marked no-op")
	}

	delayImpact := summary.Impacts[1]
	if delayImpact.Action != "duplicate" || !delayImpact.NoOp || delayImpact.Cursor != 402 || delayImpact.Committed {
		t.Fatalf("delay impact should be duplicate with metadata: %#v", delayImpact)
	}
	if delayImpact.CorrelationID != "corr-m7-o2" || delayImpact.OperationRef != "op-delay-dup" {
		t.Fatalf("delay impact metadata=%#v", delayImpact)
	}

	downImpact := summary.Impacts[2]
	if downImpact.Action != "simulation_ingress" || downImpact.NoOp || downImpact.Cursor != 403 || !downImpact.Committed {
		t.Fatalf("down impact should be committed simulation ingress: %#v", downImpact)
	}
	if downImpact.CorrelationID != "corr-m7-o2" || downImpact.OperationRef != "op-down-202607-0001" {
		t.Fatalf("down impact metadata=%#v", downImpact)
	}
}

func TestReplayLeavesOtherCanonicalEventsUnsupported(t *testing.T) {
	fake := &fakeIAOS{}
	runner, _ := New(fake)
	story := scenariopack.Story{Events: scenariopack.EventSequence{Events: []scenariopack.Event{{EventID: "evt", EventType: "whs.material.received"}}}}

	summary, err := runner.Replay(context.Background(), story, Options{Apply: true})
	if err != nil || summary.Skipped != 1 || summary.Impacts[0].Action != "unsupported" || fake.simulations != 0 || fake.decomposes != 0 {
		t.Fatalf("summary=%#v simulations=%d decomposes=%d err=%v", summary, fake.simulations, fake.decomposes, err)
	}
}

func TestReplaySimulationRequiresMatchingBusinessObjectMetadata(t *testing.T) {
	tests := []struct {
		name      string
		eventType string
		metadata  map[string]any
	}{
		{"supplier missing code", "o2d.supplier_delivery.delayed", map[string]any{"business_object_type": "purchase_order"}},
		{"quality wrong type", "qms.incoming_inspection.failed", map[string]any{"business_object_type": "goods_receipt", "business_object_id": "GR-1"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fake := &fakeIAOS{}
			runner, _ := New(fake)
			story := scenariopack.Story{Events: scenariopack.EventSequence{Events: []scenariopack.Event{{EventID: "evt", EventType: tt.eventType, Metadata: tt.metadata}}}}
			summary, err := runner.Replay(context.Background(), story, Options{Apply: true})
			if err == nil || summary.Failed != 1 || fake.simulations != 0 {
				t.Fatalf("summary=%#v simulations=%d err=%v", summary, fake.simulations, err)
			}
		})
	}
}

func TestReplayAcceptsCompleteDuplicateSimulationResponse(t *testing.T) {
	fake := &fakeIAOS{}
	fake.simulationResult.EventID = "evt-existing"
	fake.simulationResult.Subject = "iaos.tenant-hctm.o2d.supplier_delivery.delayed"
	fake.simulationResult.CorrelationID = "corr-1"
	fake.simulationResult.Duplicate = true
	fake.simulationResult.BusinessObject.Type = "purchase_order"
	fake.simulationResult.BusinessObject.Code = "PO-1"
	fake.simulationResult.BusinessObject.ID = "po-uuid"
	runner, _ := New(fake)
	story := scenariopack.Story{Events: scenariopack.EventSequence{Events: []scenariopack.Event{{EventID: "evt", EventType: "o2d.supplier_delivery.delayed", Metadata: map[string]any{"business_object_type": "purchase_order", "business_object_id": "PO-1", "idempotency_key": "delay-1", "correlation_id": "corr-1"}}}}}

	summary, err := runner.Replay(context.Background(), story, Options{Apply: true, PackKey: "hctm"})
	if err != nil || summary.Skipped != 1 || summary.Triggered != 0 || summary.Impacts[0].Action != "duplicate" {
		t.Fatalf("summary=%#v err=%v", summary, err)
	}
}

func TestReplayRejectsMalformedDuplicateSimulationResponse(t *testing.T) {
	fake := &fakeIAOS{}
	fake.simulationResult.EventID = "evt-existing"
	fake.simulationResult.Duplicate = true
	fake.simulationResult.BusinessObject.ID = "po-uuid"
	runner, _ := New(fake)
	story := scenariopack.Story{Events: scenariopack.EventSequence{Events: []scenariopack.Event{{EventID: "evt", EventType: "o2d.supplier_delivery.delayed", Metadata: map[string]any{"business_object_type": "purchase_order", "business_object_id": "PO-1", "idempotency_key": "delay-1"}}}}}

	summary, err := runner.Replay(context.Background(), story, Options{Apply: true})
	if err == nil || summary.Failed != 1 || summary.Skipped != 0 {
		t.Fatalf("summary=%#v err=%v", summary, err)
	}
}

func TestReplayRejectsSimulationSubjectForWrongTenant(t *testing.T) {
	fake := &fakeIAOS{}
	fake.simulationResult.EventID = "evt-sim-1"
	fake.simulationResult.Subject = "iaos.tenant-other.eam.machine.down"
	fake.simulationResult.Committed = true
	fake.simulationResult.BusinessObject.Type = "equipment"
	fake.simulationResult.BusinessObject.Code = "EQ-1"
	fake.simulationResult.BusinessObject.ID = "eq-uuid"
	runner, _ := New(fake)
	story := scenariopack.Story{Events: scenariopack.EventSequence{Events: []scenariopack.Event{{EventID: "evt", EventType: "eam.machine.down", Metadata: map[string]any{"business_object_type": "equipment", "business_object_id": "EQ-1", "idempotency_key": "idem-1"}}}}}

	summary, err := runner.Replay(context.Background(), story, Options{Apply: true, Tenant: "tenant-hctm"})
	if err == nil || summary.Failed != 1 || summary.Triggered != 0 {
		t.Fatalf("summary=%#v err=%v", summary, err)
	}
}

func TestReplayRejectsMalformedSimulationSubjectWhenTenantIsUnknown(t *testing.T) {
	fake := &fakeIAOS{}
	fake.simulationResult.EventID = "evt-sim-1"
	fake.simulationResult.Subject = "evil.tenant-hctm.eam.machine.down"
	fake.simulationResult.Committed = true
	fake.simulationResult.BusinessObject.Type = "equipment"
	fake.simulationResult.BusinessObject.Code = "EQ-1"
	fake.simulationResult.BusinessObject.ID = "eq-uuid"
	runner, _ := New(fake)
	story := scenariopack.Story{Events: scenariopack.EventSequence{Events: []scenariopack.Event{{EventID: "evt", EventType: "eam.machine.down", Metadata: map[string]any{"business_object_type": "equipment", "business_object_id": "EQ-1", "idempotency_key": "idem-1"}}}}}

	summary, err := runner.Replay(context.Background(), story, Options{Apply: true})
	if err == nil || summary.Failed != 1 || summary.Triggered != 0 {
		t.Fatalf("summary=%#v err=%v", summary, err)
	}
}

func TestReplayRealPackRoutingContract(t *testing.T) {
	pack, err := scenariopack.Load(filepath.Join("..", "..", "scenario-packs", "hctm"))
	if err != nil {
		t.Fatal(err)
	}
	story := pack.Stories[0]

	dryFake := &fakeIAOS{records: map[string][]map[string]any{"sales_order": {{"id": "so-uuid", "status": "confirmed"}}}}
	dryRunner, _ := New(dryFake)
	dry, err := dryRunner.Replay(context.Background(), story, Options{PackKey: pack.Manifest.PackKey, OrderID: "so-uuid", Tenant: pack.Manifest.TenantTemplate})
	if err != nil || dryFake.simulations != 0 || dryFake.decomposes != 0 || len(dryFake.businessCalls) != 0 || countReplayAction(dry.Impacts, "simulation_ingress") != 3 || countReplayAction(dry.Impacts, "scenario_business_event") != 6 || countReplayAction(dry.Impacts, "decompose_sales_order") != 1 || dry.Skipped != 12 {
		t.Fatalf("dry=%#v simulations=%d decomposes=%d err=%v", dry, dryFake.simulations, dryFake.decomposes, err)
	}

	applyFake := &fakeIAOS{records: map[string][]map[string]any{"sales_order": {{"id": "so-uuid", "status": "confirmed"}}}}
	applyRunner, _ := New(applyFake)
	applied, err := applyRunner.Replay(context.Background(), story, Options{Apply: true, PackKey: pack.Manifest.PackKey, OrderID: "so-uuid", Tenant: pack.Manifest.TenantTemplate})
	if err != nil || applyFake.simulations != 3 || len(applyFake.businessCalls) != 6 || applyFake.decomposes != 1 || applied.Triggered != 10 || applied.Skipped != 12 {
		t.Fatalf("applied=%#v simulations=%d decomposes=%d err=%v", applied, applyFake.simulations, applyFake.decomposes, err)
	}
	want := []struct{ eventType, objectType, objectCode string }{
		{"o2d.supplier_delivery.delayed", "purchase_order", "PO-202607-0001"},
		{"eam.machine.down", "equipment", "LAS-WLD-02"},
		{"qms.incoming_inspection.failed", "inspection_order", "IQC-202607-0002"},
	}
	canonicalPayloads := map[string]map[string]any{}
	for _, event := range story.Events.Events {
		canonicalPayloads[event.EventType] = event.Payload
	}
	for index, request := range applyFake.simulationCalls {
		if request.EventType != want[index].eventType || request.BusinessObject.Type != want[index].objectType || request.BusinessObject.Code != want[index].objectCode || request.Source != "aese:hctm/order-expedite-01" || !reflect.DeepEqual(request.Payload, canonicalPayloads[request.EventType]) {
			t.Fatalf("simulation[%d]=%#v", index, request)
		}
	}
}

func countReplayAction(impacts []ReplayImpact, action string) int {
	count := 0
	for _, impact := range impacts {
		if impact.Action == action {
			count++
		}
	}
	return count
}

func TestReplayFailsClosedWhenSimulationResponseIsNotCommitted(t *testing.T) {
	fake := &fakeIAOS{records: map[string][]map[string]any{"sales_order": {{"id": "so-uuid", "status": "confirmed"}}}}
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

func TestReplaySkipsOrderDecompositionAfterFulfillment(t *testing.T) {
	fake := &fakeIAOS{records: map[string][]map[string]any{"sales_order": {{"id": "so-uuid", "status": "partially_shipped"}}}}
	runner, _ := New(fake)
	story := scenariopack.Story{Events: scenariopack.EventSequence{Events: []scenariopack.Event{{EventID: "evt-1", EventType: "o2d.order.confirmed", Payload: map[string]any{"order_no": "SO-1"}}}}}
	summary, err := runner.Replay(context.Background(), story, Options{OrderID: "so-uuid", Apply: true})
	if err != nil || summary.Triggered != 0 || summary.Skipped != 1 || fake.decomposes != 0 || summary.Impacts[0].Action != "already_fulfilled" {
		t.Fatalf("summary=%#v err=%v", summary, err)
	}
}

func TestReplayUsesScenarioApplyOrderID(t *testing.T) {
	fake := &fakeIAOS{records: map[string][]map[string]any{"sales_order": {{"id": "so-uuid", "status": "confirmed"}}}}
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
