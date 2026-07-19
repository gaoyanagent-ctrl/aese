package agenttrace

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAgentRecommendationsUseRuntimeFactsAndBusinessFailsClosed(t *testing.T) {
	c := collectedContext{
		outputs: map[string]queryOutput{
			"hctm.product.read":          {Records: []map[string]any{{"id": "p-fg", "code": "HCTM-BCP-A01"}, {"id": "p-al", "code": "AL-PLATE-6061-T6"}}},
			"hctm.sales_order.read":      {Records: []map[string]any{{"id": "so-id", "order_no": "SO-202607-0001"}}},
			"hctm.sales_order_line.read": {Records: []map[string]any{{"sales_order_id": "so-id", "quantity": 12000}}},
			"hctm.inventory.read":        {Records: []map[string]any{{"product_id": "p-fg", "quantity": "1200.0000"}, {"product_id": "p-al", "quantity": "8000.0000"}}},
			"hctm.bom.read":              {Records: []map[string]any{{"parent_product_id": "p-fg", "child_product_id": "p-al", "quantity_required": "1.0500"}}},
			"hctm.purchase_order.read":   {Records: []map[string]any{{"po_no": "PO-202607-0001", "order_qty": "5000.0000", "status": "delayed"}, {"po_no": "PO-202607-0002", "supplier_code": "SUP-BETA-AL"}}},
			"hctm.equipment.read":        {Records: []map[string]any{{"code": "LAS-WLD-02", "status": "maintenance"}}},
			"hctm.inspection_order.read": {Records: []map[string]any{{"inspection_no": "IQC-202607-0002", "po_no": "PO-202607-0002", "receipt_no": "GR-202607-0002", "lot_no": "BETA-20260712-01", "rejected_qty": "300.0000", "accepted_qty": "0.0000", "defect_code": "SURFACE_SCRATCH", "severity": "major", "inspection_level": "normal", "status": "failed"}}},
			"hctm.work_order.read":       {Records: []map[string]any{{"wo_no": "WO-1"}}},
		},
		calls: map[string]string{"hctm.sales_order.read": "call-so", "hctm.inspection_order.read": "call-iqc"},
	}
	planning, err := buildPlanning("corr-so-202607-0001", c)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(planning.Summary, "12,000") && !strings.Contains(planning.Summary, "12000") {
		t.Fatalf("planning summary lacks demand: %s", planning.Summary)
	}
	if !strings.Contains(planning.Summary, "10800") && !strings.Contains(planning.Summary, "10,800") {
		t.Fatalf("planning summary lacks net demand: %s", planning.Summary)
	}
	quality, err := buildQuality("corr-so-202607-0001", c)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(quality.Summary, "SURFACE_SCRATCH") || quality.Completeness != "partial" {
		t.Fatalf("unexpected quality recommendation: %+v", quality)
	}
	business, err := buildBusiness("corr-so-202607-0001", c, planning, quality)
	if err != nil {
		t.Fatal(err)
	}
	if business.Completeness != "partial" || !strings.Contains(business.Summary, "不能判断") {
		t.Fatalf("business recommendation did not fail closed: %+v", business)
	}
	if strings.Contains(business.Summary, "11,700") || strings.Contains(business.Summary, "11700") {
		t.Fatalf("business recommendation fabricated shipment: %s", business.Summary)
	}
}

func TestLoadCanonicalAgentToolBundle(t *testing.T) {
	bundle, err := LoadBundle("../../scenario-packs/hctm")
	if err != nil {
		t.Fatal(err)
	}
	if len(bundle.MetadataSchemas) != 5 || len(bundle.Tools) != 9 {
		t.Fatalf("unexpected bundle counts: schemas=%d tools=%d", len(bundle.MetadataSchemas), len(bundle.Tools))
	}
}

func TestAgentToolBundleFailsClosedForWrongPackOrMissingTool(t *testing.T) {
	bundle, err := LoadBundle("../../scenario-packs/hctm")
	if err != nil {
		t.Fatal(err)
	}
	if err := bundle.ValidatePack("other"); err == nil {
		t.Fatal("expected mismatched pack to fail")
	}

	data, err := os.ReadFile("../../scenario-packs/hctm/agent-tools.json")
	if err != nil {
		t.Fatal(err)
	}
	data = []byte(strings.Replace(string(data), `"tool_key":"hctm.work_order.read"`, `"tool_key":"hctm.work_order.missing"`, 1))
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "agent-tools.json"), data, 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadBundle(dir); err == nil || !strings.Contains(err.Error(), "missing required tool") {
		t.Fatalf("expected missing required tool error, got %v", err)
	}
}

func TestPlanningReportsRuntimeMaterialShortfall(t *testing.T) {
	c := collectedContext{
		outputs: map[string]queryOutput{
			"hctm.product.read":          {Records: []map[string]any{{"id": "p-fg", "code": "HCTM-BCP-A01"}, {"id": "p-al", "code": "AL-PLATE-6061-T6"}}},
			"hctm.sales_order.read":      {Records: []map[string]any{{"id": "so-id", "order_no": "SO-202607-0001"}}},
			"hctm.sales_order_line.read": {Records: []map[string]any{{"sales_order_id": "so-id", "quantity": 12000}}},
			"hctm.inventory.read":        {Records: []map[string]any{{"product_id": "p-fg", "quantity": "1200.0000"}, {"product_id": "p-al", "quantity": "0.0000"}}},
			"hctm.bom.read":              {Records: []map[string]any{{"parent_product_id": "p-fg", "child_product_id": "p-al", "quantity_required": "1.0500"}}},
			"hctm.purchase_order.read":   {Records: []map[string]any{{"po_no": "PO-202607-0001", "order_qty": "5000.0000", "status": "delayed"}}},
			"hctm.equipment.read":        {Records: []map[string]any{{"code": "LAS-WLD-02", "status": "maintenance"}}},
			"hctm.work_order.read":       {Records: []map[string]any{{"wo_no": "WO-1"}}},
		},
		calls: map[string]string{},
	}
	planning, err := buildPlanning("corr-so-202607-0001", c)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(planning.Summary, "仍缺 7600") || !contains(planning.Risks, "material_quantity_shortfall") {
		t.Fatalf("planning recommendation hid current material shortfall: %+v", planning)
	}
}

func contains(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
