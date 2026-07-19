package validate

import (
	"strings"
	"testing"

	"github.com/industrial-ai/iaos-aese/internal/scenariopack"
)

func TestPackValid(t *testing.T) {
	pack := fixturePack()
	result := Pack(pack)
	if !result.Valid() {
		t.Fatalf("expected valid pack, got:\n%s", issueText(result))
	}
}

func TestPackReportsBrokenFixtures(t *testing.T) {
	tests := []struct {
		name, want string
		breakPack  func(*scenariopack.Pack)
	}{
		{"duplicate natural key", "duplicate natural key", func(p *scenariopack.Pack) {
			p.RecordSets[0].Records = append(p.RecordSets[0].Records, p.RecordSets[0].Records[0])
		}},
		{"missing reference", "references missing customer", func(p *scenariopack.Pack) { p.RecordSets[2].Records[0]["customer_code"] = "missing" }},
		{"inspection purchase reference", "references missing purchase_order", func(p *scenariopack.Pack) {
			p.RecordSets = append(p.RecordSets,
				scenariopack.RecordSet{SchemaVersion: "1.0.0", Entity: "purchase_order", NaturalKey: []string{"po_no"}, Source: "initial.json", Records: []map[string]any{{"po_no": "PO-1", "material_code": "M-1"}}},
				scenariopack.RecordSet{SchemaVersion: "1.0.0", Entity: "inspection_order", NaturalKey: []string{"inspection_no"}, Source: "initial.json", Records: []map[string]any{{"inspection_no": "IQC-1", "po_no": "PO-MISSING", "material_code": "M-1"}}},
			)
		}},
		{"out of order", "timeline is out of order", func(p *scenariopack.Pack) { p.Stories[0].Events.Events[1].Timestamp = "2026-06-30T00:00:00Z" }},
		{"duplicate idempotency", "duplicate idempotency key", func(p *scenariopack.Pack) { p.Stories[0].Events.Events[1].IdempotencyKey = "idem-1" }},
		{"negative quantity", "quantity must be non-negative", func(p *scenariopack.Pack) { p.RecordSets[2].Records[0]["quantity"] = -1 }},
		{"shipment invariant", "exceeds available finished goods", func(p *scenariopack.Pack) { p.Stories[0].Events.Events[1].Payload["quantity"] = 200.0 }},
		{"expected outcome", "assertion event_count evaluated", func(p *scenariopack.Pack) { p.Stories[0].Expected.Assertions[0].Expected = 3 }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := fixturePack()
			tt.breakPack(p)
			got := issueText(Pack(p))
			if !strings.Contains(got, tt.want) {
				t.Fatalf("want %q in:\n%s", tt.want, got)
			}
		})
	}
}

func fixturePack() *scenariopack.Pack {
	sets := []scenariopack.RecordSet{
		{SchemaVersion: "1.0.0", Entity: "customer", NaturalKey: []string{"customer_code"}, Source: "master.json", Records: []map[string]any{{"customer_code": "C-1"}}},
		{SchemaVersion: "1.0.0", Entity: "material", NaturalKey: []string{"material_code"}, Source: "master.json", Records: []map[string]any{{"material_code": "M-1"}}},
		{SchemaVersion: "1.0.0", Entity: "sales_order", NaturalKey: []string{"sales_order_no"}, Source: "initial.json", Records: []map[string]any{{"sales_order_no": "SO-1", "customer_code": "C-1", "material_code": "M-1", "quantity": 100.0}}},
		{SchemaVersion: "1.0.0", Entity: "inventory_transaction", NaturalKey: []string{"transaction_no"}, Source: "initial.json", Records: []map[string]any{{"transaction_no": "INV-1", "material_code": "M-1", "quantity": 100.0, "direction": "in"}}},
	}
	story := scenariopack.Story{Ref: scenariopack.StoryRef{Key: "story"}, Initial: scenariopack.InitialState{SchemaVersion: "1.0.0", StoryKey: "story", Source: "initial.json"}, Events: scenariopack.EventSequence{SchemaVersion: "1.0.0", StoryKey: "story", CorrelationID: "corr-1", Source: "events.json", Events: []scenariopack.Event{
		{EventID: "evt-1", EventType: "o2d.order.confirmed", Timestamp: "2026-07-01T00:00:00Z", CorrelationID: "corr-1", IdempotencyKey: "idem-1", Payload: map[string]any{"quantity": 100.0}},
		{EventID: "evt-2", EventType: "o2d.shipment.dispatched", Timestamp: "2026-07-02T00:00:00Z", CorrelationID: "corr-1", CausationID: "evt-1", IdempotencyKey: "idem-2", Payload: map[string]any{"material_code": "M-1", "quantity": 100.0}},
	}}, Expected: scenariopack.ExpectedOutcomes{SchemaVersion: "1.0.0", StoryKey: "story", Source: "expected.json", Assertions: []scenariopack.Assertion{{Key: "event_count", Type: "event_count", Field: "events", Operator: "count_equals", Expected: 2}}}}
	return &scenariopack.Pack{Manifest: scenariopack.Manifest{SchemaVersion: "1.0.0", PackKey: "test", PackVersion: "0.1.0", TenantTemplate: "tenant-test"}, RecordSets: sets, Stories: []scenariopack.Story{story}}
}

func issueText(result Result) string {
	values := make([]string, len(result.Issues))
	for i, v := range result.Issues {
		values[i] = v.Error()
	}
	return strings.Join(values, "\n")
}
