package agenttrace

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/industrial-ai/iaos-aese/internal/iaosclient"
)

type ToolEvidence struct {
	ToolKey string `json:"tool_key"`
	CallID  string `json:"call_id"`
	Count   int    `json:"count"`
}

type Fact struct {
	Key      string `json:"key"`
	Value    any    `json:"value"`
	Evidence string `json:"evidence"`
}

type Recommendation struct {
	AgentKey                  string   `json:"agent_key"`
	CorrelationID             string   `json:"correlation_id"`
	Summary                   string   `json:"summary"`
	Facts                     []Fact   `json:"facts"`
	Risks                     []string `json:"risks"`
	Recommendations           []string `json:"recommendations"`
	ObjectRefs                []string `json:"object_refs"`
	ToolCallIDs               []string `json:"tool_call_ids"`
	Completeness              string   `json:"completeness"`
	DataGaps                  []string `json:"data_gaps,omitempty"`
	Confidence                string   `json:"confidence"`
	Status                    string   `json:"status"`
	RequiresHumanConfirmation bool     `json:"requires_human_confirmation"`
}

type RunSummary struct {
	RunID         string           `json:"run_id"`
	Mode          string           `json:"mode"`
	PackKey       string           `json:"pack_key"`
	StoryKey      string           `json:"story_key"`
	CorrelationID string           `json:"correlation_id"`
	ToolEvidence  []ToolEvidence   `json:"tool_evidence"`
	Agents        []Recommendation `json:"agents"`
}

type queryOutput struct {
	EntityCode string           `json:"entity_code"`
	Records    []map[string]any `json:"records"`
	Count      int              `json:"count"`
}

type collectedContext struct {
	outputs map[string]queryOutput
	calls   map[string]string
}

func Run(ctx context.Context, client *iaosclient.Client, packKey, storyKey, correlationID, runID string, apply bool) (RunSummary, error) {
	summary := RunSummary{RunID: runID, Mode: "dry-run", PackKey: packKey, StoryKey: storyKey, CorrelationID: correlationID}
	if !apply {
		for _, key := range requiredToolKeys() {
			summary.ToolEvidence = append(summary.ToolEvidence, ToolEvidence{ToolKey: key})
		}
		return summary, nil
	}
	summary.Mode = "apply"
	collected := collectedContext{outputs: map[string]queryOutput{}, calls: map[string]string{}}
	queries := []struct {
		key   string
		input map[string]any
	}{
		{"hctm.product.read", map[string]any{"limit": 20}},
		{"hctm.sales_order.read", map[string]any{"filters": map[string]any{"order_no": "SO-202607-0001"}}},
		{"hctm.sales_order_line.read", map[string]any{"limit": 100}},
		{"hctm.inventory.read", map[string]any{"limit": 100}},
		{"hctm.bom.read", map[string]any{"limit": 100}},
		{"hctm.purchase_order.read", map[string]any{"limit": 50}},
		{"hctm.equipment.read", map[string]any{"filters": map[string]any{"code": "LAS-WLD-02"}}},
		{"hctm.inspection_order.read", map[string]any{"filters": map[string]any{"inspection_no": "IQC-202607-0002"}}},
		{"hctm.work_order.read", map[string]any{"limit": 100}},
	}
	for _, query := range queries {
		result, err := client.CallAITool(ctx, query.key, correlationID, runID, query.input)
		if err != nil {
			return summary, fmt.Errorf("call %s: %w", query.key, err)
		}
		var output queryOutput
		if err := json.Unmarshal(result.Output, &output); err != nil {
			return summary, fmt.Errorf("decode %s output: %w", query.key, err)
		}
		if output.EntityCode == "" || output.Records == nil {
			return summary, fmt.Errorf("tool %s returned malformed entity output", query.key)
		}
		collected.outputs[query.key] = output
		collected.calls[query.key] = result.CallID
		summary.ToolEvidence = append(summary.ToolEvidence, ToolEvidence{ToolKey: query.key, CallID: result.CallID, Count: output.Count})
	}

	planning, err := buildPlanning(correlationID, collected)
	if err != nil {
		return summary, err
	}
	quality, err := buildQuality(correlationID, collected)
	if err != nil {
		return summary, err
	}
	business, err := buildBusiness(correlationID, collected, planning, quality)
	if err != nil {
		return summary, err
	}
	summary.Agents = []Recommendation{planning, quality, business}
	return summary, nil
}

func requiredToolKeys() []string {
	return []string{"hctm.product.read", "hctm.sales_order.read", "hctm.sales_order_line.read", "hctm.inventory.read", "hctm.bom.read", "hctm.purchase_order.read", "hctm.equipment.read", "hctm.inspection_order.read", "hctm.work_order.read"}
}

func buildPlanning(correlationID string, c collectedContext) (Recommendation, error) {
	order, err := exactlyOne(c.outputs["hctm.sales_order.read"].Records, "sales order")
	if err != nil {
		return Recommendation{}, err
	}
	orderID := text(order["id"])
	products := indexBy(c.outputs["hctm.product.read"].Records, "code")
	finished := products["HCTM-BCP-A01"]
	aluminum := products["AL-PLATE-6061-T6"]
	if finished == nil || aluminum == nil {
		return Recommendation{}, fmt.Errorf("planning context is missing required products")
	}
	demand := sumWhere(c.outputs["hctm.sales_order_line.read"].Records, "sales_order_id", orderID, "quantity")
	fg := sumWhere(c.outputs["hctm.inventory.read"].Records, "product_id", text(finished["id"]), "quantity")
	alInventory := sumWhere(c.outputs["hctm.inventory.read"].Records, "product_id", text(aluminum["id"]), "quantity")
	usage := sumPair(c.outputs["hctm.bom.read"].Records, text(finished["id"]), text(aluminum["id"]), "quantity_required")
	required := multiply(demand, usage)
	net := subtract(demand, fg)
	po1 := find(c.outputs["hctm.purchase_order.read"].Records, "po_no", "PO-202607-0001")
	equipment, err := exactlyOne(c.outputs["hctm.equipment.read"].Records, "equipment")
	if err != nil {
		return Recommendation{}, err
	}
	if po1 == nil {
		return Recommendation{}, fmt.Errorf("planning context is missing PO-202607-0001")
	}
	po1Qty := quantity(po1["order_qty"])
	confirmedSupply := new(big.Rat).Add(new(big.Rat).Set(alInventory), po1Qty)
	materialGap := subtract(required, confirmedSupply)
	materialAssessment := "铝板现有库存与 PO-202607-0001 数量合计可覆盖需求，但该采购单延期，供料时点仍有风险"
	risks := []string{"supplier_eta_timing_risk", "welding_capacity_risk"}
	if materialGap.Sign() > 0 {
		materialAssessment = fmt.Sprintf("铝板现有库存与 PO-202607-0001 数量合计仍缺 %s，且该采购单延期", decimal(materialGap))
		risks = append([]string{"material_quantity_shortfall"}, risks...)
	} else {
		materialGap = new(big.Rat)
	}
	calls := callIDs(c, "hctm.sales_order.read", "hctm.sales_order_line.read", "hctm.product.read", "hctm.inventory.read", "hctm.bom.read", "hctm.purchase_order.read", "hctm.equipment.read", "hctm.work_order.read")
	return Recommendation{
		AgentKey: "planning", CorrelationID: correlationID,
		Summary:         fmt.Sprintf("订单需求 %s，成品库存 %s，净生产需求 %s；%s，LAS-WLD-02 为 %s，交付与焊接产能仍有风险。", decimal(demand), decimal(fg), decimal(net), materialAssessment, text(equipment["status"])),
		Facts:           []Fact{{"demand_qty", decimal(demand), "SO-202607-0001"}, {"finished_goods_qty", decimal(fg), "inventory:HCTM-BCP-A01"}, {"net_production_qty", decimal(net), "calculated"}, {"aluminum_required_qty", decimal(required), "BOM:AL-PLATE-6061-T6"}, {"aluminum_inventory_qty", decimal(alInventory), "inventory:AL-PLATE-6061-T6"}, {"delayed_po_qty", text(po1["order_qty"]), "PO-202607-0001"}, {"material_supply_gap_qty", decimal(materialGap), "calculated from governed inventory and PO-202607-0001"}, {"equipment_status", text(equipment["status"]), "LAS-WLD-02"}},
		Risks:           risks,
		Recommendations: []string{"评估并经人工批准后释放备选供应商 PO-202607-0002", "为电池冷却板 A 线准备加班方案", "保留分批发运和客户重承诺草稿"},
		ObjectRefs:      []string{"SO-202607-0001", "PO-202607-0001", "PO-202607-0002", "LAS-WLD-02", "HCTM-BCP-A01", "AL-PLATE-6061-T6"}, ToolCallIDs: calls,
		Completeness: "complete_for_current_risk", Confidence: "high", Status: "suggested", RequiresHumanConfirmation: true,
	}, nil
}

func buildQuality(correlationID string, c collectedContext) (Recommendation, error) {
	iqc, err := exactlyOne(c.outputs["hctm.inspection_order.read"].Records, "inspection order")
	if err != nil {
		return Recommendation{}, err
	}
	po := find(c.outputs["hctm.purchase_order.read"].Records, "po_no", text(iqc["po_no"]))
	if po == nil {
		return Recommendation{}, fmt.Errorf("quality context is missing related purchase order")
	}
	if text(iqc["status"]) != "failed" {
		return Recommendation{}, fmt.Errorf("quality tracer requires a failed inspection")
	}
	gaps := []string{}
	if text(iqc["inspection_level"]) != "tightened" {
		gaps = append(gaps, "tightened_inspection_not_recorded")
	}
	if quantity(iqc["accepted_qty"]).Sign() == 0 {
		gaps = append(gaps, "accepted_quantity_not_released")
	}
	return Recommendation{
		AgentKey: "quality", CorrelationID: correlationID,
		Summary:         fmt.Sprintf("%s 的批次 %s 在 %s 发现 %s 张 %s；不合格数量不得直接投产，根因证据尚不足。", text(po["supplier_code"]), text(iqc["lot_no"]), text(iqc["inspection_no"]), decimal(quantity(iqc["rejected_qty"])), text(iqc["defect_code"])),
		Facts:           []Fact{{"supplier_code", text(po["supplier_code"]), text(po["po_no"])}, {"lot_no", text(iqc["lot_no"]), text(iqc["inspection_no"])}, {"rejected_qty", text(iqc["rejected_qty"]), text(iqc["inspection_no"])}, {"defect_code", text(iqc["defect_code"]), text(iqc["inspection_no"])}, {"severity", text(iqc["severity"]), text(iqc["inspection_no"])}},
		Risks:           []string{"incoming_material_quality_risk", "root_cause_insufficient_evidence"},
		Recommendations: []string{"隔离不合格数量并阻止直接投产", "保留供应商、收货号和批次追溯", "合格数量仅在质量放行后使用", "发起供应商纠正措施草稿"},
		ObjectRefs:      []string{text(iqc["inspection_no"]), text(iqc["po_no"]), text(iqc["receipt_no"]), text(iqc["lot_no"]), text(po["supplier_code"])}, ToolCallIDs: callIDs(c, "hctm.inspection_order.read", "hctm.purchase_order.read"),
		Completeness: "partial", DataGaps: gaps, Confidence: "high", Status: "suggested", RequiresHumanConfirmation: true,
	}, nil
}

func buildBusiness(correlationID string, c collectedContext, planning, quality Recommendation) (Recommendation, error) {
	order, err := exactlyOne(c.outputs["hctm.sales_order.read"].Records, "sales order")
	if err != nil {
		return Recommendation{}, err
	}
	demand := "unknown"
	for _, fact := range planning.Facts {
		if fact.Key == "demand_qty" {
			demand = fmt.Sprint(fact.Value)
		}
	}
	return Recommendation{
		AgentKey: "business_analysis", CorrelationID: correlationID,
		Summary:         fmt.Sprintf("订单 %s 的需求为 %s。供应延期、设备停机和来料不良均已形成可查询风险，但 IAOS 当前没有完工入库、发运和实际成本事实，因此不能判断最终交付数量、缺口或利润影响。", text(order["order_no"]), demand),
		Facts:           []Fact{{"demand_qty", demand, text(order["order_no"])}, {"delivery_status", "not_determinable", "missing shipment facts"}, {"cost_impact", "qualitative_only", "missing cost actuals"}},
		Risks:           []string{"delivery_commitment_open", "procurement_cost_pressure", "overtime_cost_pressure", "quality_cost_pressure"},
		Recommendations: []string{"补齐受治理的完工入库和发运事实后再计算交付缺口", "补齐加急采购、加班和检验实际成本后再计算利润影响", "保留 300 件补产或客户重承诺作为待决策草稿"},
		ObjectRefs:      []string{text(order["order_no"]), "PO-202607-0001", "LAS-WLD-02", "IQC-202607-0002"}, ToolCallIDs: callIDs(c, requiredToolKeys()...),
		Completeness: "partial", DataGaps: []string{"finished_goods_receipt", "shipment_dispatch", "cost_actuals"}, Confidence: "medium", Status: "suggested", RequiresHumanConfirmation: true,
	}, nil
}

func exactlyOne(records []map[string]any, label string) (map[string]any, error) {
	if len(records) != 1 {
		return nil, fmt.Errorf("%s query returned %d records", label, len(records))
	}
	return records[0], nil
}
func find(records []map[string]any, key, value string) map[string]any {
	for _, r := range records {
		if text(r[key]) == value {
			return r
		}
	}
	return nil
}
func indexBy(records []map[string]any, key string) map[string]map[string]any {
	out := map[string]map[string]any{}
	for _, r := range records {
		out[text(r[key])] = r
	}
	return out
}
func text(v any) string {
	if v == nil {
		return ""
	}
	return fmt.Sprint(v)
}
func quantity(v any) *big.Rat {
	r := new(big.Rat)
	if _, ok := r.SetString(text(v)); !ok {
		return new(big.Rat)
	}
	return r
}
func sumWhere(records []map[string]any, key, want, value string) *big.Rat {
	out := new(big.Rat)
	for _, r := range records {
		if text(r[key]) == want {
			out.Add(out, quantity(r[value]))
		}
	}
	return out
}
func sumPair(records []map[string]any, parent, child, value string) *big.Rat {
	out := new(big.Rat)
	for _, r := range records {
		if text(r["parent_product_id"]) == parent && text(r["child_product_id"]) == child {
			out.Add(out, quantity(r[value]))
		}
	}
	return out
}
func subtract(a, b *big.Rat) *big.Rat { return new(big.Rat).Sub(a, b) }
func multiply(a, b *big.Rat) *big.Rat { return new(big.Rat).Mul(a, b) }
func decimal(v *big.Rat) string {
	if v == nil {
		return "0"
	}
	s := v.FloatString(4)
	return strings.TrimRight(strings.TrimRight(s, "0"), ".")
}
func callIDs(c collectedContext, keys ...string) []string {
	out := []string{}
	for _, k := range keys {
		if id := c.calls[k]; id != "" {
			out = append(out, id)
		}
	}
	return out
}
