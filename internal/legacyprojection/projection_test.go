package legacyprojection

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"github.com/industrial-ai/iaos-aese/internal/scenariopack"
)

func TestProjectRealHCTMPack(t *testing.T) {
	pack, err := scenariopack.Load(filepath.Join("..", "..", "scenario-packs", "hctm"))
	if err != nil {
		t.Fatal(err)
	}
	result, err := Project(pack, Options{StoryKey: "order-expedite-01", RunID: "run-test-001", DryRun: true})
	if err != nil {
		t.Fatal(err)
	}
	if result.Request.PackKey != "hctm" || result.Request.PackVersion != "0.1.0" || !result.Request.DryRun {
		t.Fatalf("unexpected request header: %+v", result.Request)
	}
	counts := map[string]int{}
	for _, object := range result.Request.Objects {
		counts[object.Type]++
		if _, hasID := object.Data["id"]; hasID {
			t.Fatalf("projection must not contain an environment UUID: %+v", object)
		}
	}
	wantCounts := map[string]int{"customer": 1, "material": 6, "bom": 5, "inventory": 5, "purchase_order": 2, "inspection_order": 1, "sales_order": 1}
	if len(result.Request.Objects) != 21 {
		t.Fatalf("object count = %d, want 21", len(result.Request.Objects))
	}
	for entity, want := range wantCounts {
		if counts[entity] != want {
			t.Errorf("%s count = %d, want %d", entity, counts[entity], want)
		}
	}

	order := onlyObject(t, result, "sales_order")
	if got := order.NaturalKey["order_no"]; got != TracerOrderNo {
		t.Fatalf("order natural key = %v", got)
	}
	if got := order.Data["status"]; got != "draft" {
		t.Fatalf("order status = %v, want draft", got)
	}
	if got := order.Data["required_date"]; got != "2026-07-20T00:00:00+08:00" {
		t.Fatalf("required_date = %v, want explicit Asia/Shanghai RFC3339", got)
	}
	lines, ok := order.Data["lines"].([]map[string]any)
	if !ok || len(lines) != 1 {
		t.Fatalf("order lines = %#v", order.Data["lines"])
	}
	if lines[0]["material_code"] != "HCTM-BCP-A01" || lines[0]["quantity"] != int64(12000) || lines[0]["unit_price"] != DemoUnitPrice {
		t.Fatalf("unexpected merged order line: %#v", lines[0])
	}

	bom := firstObject(t, result, "bom")
	if bom.Data["parent_material_code"] != "HCTM-BCP-A01" || bom.Data["component_material_code"] == "" {
		t.Fatalf("BOM must preserve stable business codes: %#v", bom.Data)
	}
	if _, ok := bom.Data["parent_product_id"]; ok {
		t.Fatal("BOM projection resolved an IAOS UUID client-side")
	}
	customer := onlyObject(t, result, "customer")
	if customer.NaturalKey["code"] != "CUST-SGNEV" {
		t.Fatalf("customer target natural key = %#v", customer.NaturalKey)
	}
	inventory := firstObject(t, result, "inventory")
	if inventory.NaturalKey["warehouse_name"] == nil || inventory.NaturalKey["batch_no"] == nil {
		t.Fatalf("inventory target natural key = %#v", inventory.NaturalKey)
	}
	if !hasDroppedField(result.Warnings, "bom", "scrap_rate") {
		t.Fatal("expected an explicit warning for dropped BOM scrap_rate")
	}
	if !hasDroppedField(result.Warnings, "sales_order", "legal_entity_code") {
		t.Fatal("expected an explicit warning for dropped sales order legal_entity_code")
	}

	firstPO := objectByNaturalKey(t, result, "purchase_order", "po_no", "PO-202607-0001")
	if firstPO.Data["supplier_code"] != "SUP-ALPHA-AL" || firstPO.Data["material_code"] != "AL-PLATE-6061-T6" || firstPO.Data["order_qty"] != json.Number("5000") {
		t.Fatalf("unexpected first purchase order: %#v", firstPO.Data)
	}
	if firstPO.Data["promised_date"] != "2026-07-08" || firstPO.Data["latest_eta"] != "2026-07-08" || firstPO.Data["status"] != "in_transit" {
		t.Fatalf("purchase order dates/status must preserve canonical values: %#v", firstPO.Data)
	}
	secondPO := objectByNaturalKey(t, result, "purchase_order", "po_no", "PO-202607-0002")
	if secondPO.Data["supplier_code"] != "SUP-BETA-AL" || secondPO.Data["order_qty"] != json.Number("2000") || secondPO.Data["promised_date"] != "2026-07-12" || secondPO.Data["latest_eta"] != "2026-07-12" || secondPO.Data["status"] != "draft" {
		t.Fatalf("unexpected second purchase order: %#v", secondPO.Data)
	}
	for _, field := range []string{"po_no", "supplier_code", "material_code", "order_qty", "promised_date", "latest_eta", "status"} {
		if hasDroppedField(result.Warnings, "purchase_order", field) {
			t.Fatalf("mapped purchase order field %s reported as dropped", field)
		}
	}
	if !hasDroppedField(result.Warnings, "purchase_order", "legal_entity_code") || !hasDroppedField(result.Warnings, "purchase_order", "plant_code") {
		t.Fatal("expected explicit warnings for unsupported purchase order fields")
	}
	inspection := onlyObject(t, result, "inspection_order")
	if inspection.NaturalKey["inspection_no"] != "IQC-202607-0002" || inspection.Data["po_no"] != "PO-202607-0002" || inspection.Data["material_code"] != "AL-PLATE-6061-T6" {
		t.Fatalf("unexpected inspection identity: %#v", inspection)
	}
	if inspection.Data["inspection_type"] != "incoming" || inspection.Data["sample_qty"] != json.Number("80") || inspection.Data["accepted_qty"] != json.Number("0") || inspection.Data["rejected_qty"] != json.Number("0") || inspection.Data["status"] != "pending" {
		t.Fatalf("unexpected pending inspection data: %#v", inspection.Data)
	}
	if _, exists := inspection.Data["receipt_no"]; exists {
		t.Fatalf("pending inspection must not invent a receipt: %#v", inspection.Data)
	}
	if _, exists := inspection.Data["lot_no"]; exists {
		t.Fatalf("pending inspection must not invent a lot: %#v", inspection.Data)
	}
	for _, field := range []string{"inspection_no", "po_no", "material_code", "inspection_type", "sample_qty", "accepted_qty", "rejected_qty", "status"} {
		if hasDroppedField(result.Warnings, "inspection_order", field) {
			t.Fatalf("mapped inspection order field %s reported as dropped", field)
		}
	}

	wire, err := json.Marshal(result.Request)
	if err != nil {
		t.Fatal(err)
	}
	text := string(wire)
	for _, fragment := range []string{`"objects"`, `"type":"sales_order"`, `"natural_key":{"order_no":"SO-202607-0001"}`, `"type":"purchase_order"`, `"natural_key":{"po_no":"PO-202607-0001"}`, `"type":"inspection_order"`, `"natural_key":{"inspection_no":"IQC-202607-0002"}`, `"unit_price":"128.5000"`} {
		if !strings.Contains(text, fragment) {
			t.Errorf("wire request missing %s: %s", fragment, text)
		}
	}
}

func TestProjectRequiresExplicitRunID(t *testing.T) {
	_, err := Project(&scenariopack.Pack{}, Options{})
	if err == nil || !strings.Contains(err.Error(), "run id") {
		t.Fatalf("expected run id error, got %v", err)
	}
}

func TestProjectRejectsNonTracerQuantity(t *testing.T) {
	pack, err := scenariopack.Load(filepath.Join("..", "..", "scenario-packs", "hctm"))
	if err != nil {
		t.Fatal(err)
	}
	pack.Stories[0].Initial.RecordSets[1].Records[1]["order_qty"] = json.Number("1999")
	_, err = Project(pack, Options{StoryKey: "order-expedite-01", RunID: "run-test-002"})
	if err == nil || !strings.Contains(err.Error(), "totalling 12000") {
		t.Fatalf("expected tracer total error, got %v", err)
	}
}

func onlyObject(t *testing.T, result Result, entity string) Object {
	t.Helper()
	var found []Object
	for _, object := range result.Request.Objects {
		if object.Type == entity {
			found = append(found, object)
		}
	}
	if len(found) != 1 {
		t.Fatalf("got %d %s objects", len(found), entity)
	}
	return found[0]
}

func firstObject(t *testing.T, result Result, entity string) Object {
	t.Helper()
	for _, object := range result.Request.Objects {
		if object.Type == entity {
			return object
		}
	}
	t.Fatalf("no %s object", entity)
	return Object{}
}

func objectByNaturalKey(t *testing.T, result Result, entity, field, value string) Object {
	t.Helper()
	for _, object := range result.Request.Objects {
		if object.Type == entity && object.NaturalKey[field] == value {
			return object
		}
	}
	t.Fatalf("no %s object with %s=%s", entity, field, value)
	return Object{}
}

func hasDroppedField(warnings []Warning, entity, field string) bool {
	for _, warning := range warnings {
		if warning.SourceEntity != entity {
			continue
		}
		for _, dropped := range warning.DroppedFields {
			if dropped == field {
				return true
			}
		}
	}
	return false
}
