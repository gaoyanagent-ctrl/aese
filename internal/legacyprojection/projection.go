// Package legacyprojection projects the narrow HCTM M3 tracer onto the
// IAOS DES-047 scenario apply wire contract. It deliberately keeps business
// codes in the request; resolving codes to tenant-local UUIDs belongs to IAOS.
package legacyprojection

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/industrial-ai/iaos-aese/internal/scenariopack"
)

const (
	TracerOrderNo = "SO-202607-0001"
	// DemoUnitPrice is an explicit fixture value, not a price inferred by IAOS.
	DemoUnitPrice = "128.5000"
)

// ApplyRequest is the request body accepted by POST /api/v1/scenarios/apply.
type ApplyRequest struct {
	PackKey     string   `json:"pack_key"`
	PackVersion string   `json:"pack_version"`
	ScenarioKey string   `json:"scenario_key"`
	RunID       string   `json:"run_id"`
	DryRun      bool     `json:"dry_run"`
	Objects     []Object `json:"objects"`
}

type Object struct {
	Type           string         `json:"type"`
	NaturalKey     map[string]any `json:"natural_key"`
	Data           map[string]any `json:"data"`
	IdempotencyKey string         `json:"idempotency_key,omitempty"`
}

// Warning makes every intentionally omitted HCTM field visible to callers.
type Warning struct {
	SourceEntity  string         `json:"source_entity"`
	NaturalKey    map[string]any `json:"natural_key"`
	DroppedFields []string       `json:"dropped_fields"`
	Message       string         `json:"message"`
}

type Result struct {
	Request  ApplyRequest `json:"request"`
	Warnings []Warning    `json:"warnings"`
}

type Options struct {
	StoryKey string
	RunID    string
	DryRun   bool
}

// Project creates the single M3 legacy tracer request. Only the customer used
// by the tracer is imported; materials are limited to those referenced by the
// tracer order, BOM and opening inventory. The two HCTM order records are
// merged into one 12,000-unit draft order with one line.
func Project(pack *scenariopack.Pack, opts Options) (Result, error) {
	if pack == nil {
		return Result{}, fmt.Errorf("scenario pack is required")
	}
	if strings.TrimSpace(opts.RunID) == "" {
		return Result{}, fmt.Errorf("run id is required")
	}
	story, err := findStory(pack, opts.StoryKey)
	if err != nil {
		return Result{}, err
	}
	sets := indexSets(pack, story)
	orders, ok := sets["sales_order"]
	if !ok {
		return Result{}, fmt.Errorf("story %s: sales_order record set is required", story.Ref.Key)
	}

	merged, componentCodes, customerCode, err := mergeTracerOrder(orders)
	if err != nil {
		return Result{}, err
	}
	result := Result{Request: ApplyRequest{
		PackKey: pack.Manifest.PackKey, PackVersion: pack.Manifest.PackVersion,
		ScenarioKey: story.Ref.Key, RunID: opts.RunID, DryRun: opts.DryRun,
		Objects: []Object{},
	}, Warnings: []Warning{}}

	customers, ok := sets["customer"]
	if !ok {
		return Result{}, fmt.Errorf("customer record set is required")
	}
	customer, err := recordBy(customers, "customer_code", customerCode)
	if err != nil {
		return Result{}, err
	}
	result.add("customer", customers, customer,
		[]string{"customer_code"},
		map[string]string{"customer_code": "code", "customer_name": "name", "customer_type": "oem_category"})

	// BOM and inventory expand the material dependency set before products are
	// emitted, preserving dependency order expected by IAOS.
	boms, ok := sets["bom"]
	if !ok {
		return Result{}, fmt.Errorf("bom record set is required")
	}
	for _, record := range boms.Records {
		componentCodes[stringField(record, "parent_material_code")] = true
		componentCodes[stringField(record, "component_material_code")] = true
	}
	inventories, ok := sets["inventory_transaction"]
	if !ok {
		return Result{}, fmt.Errorf("inventory_transaction record set is required")
	}
	for _, record := range inventories.Records {
		if stringField(record, "transaction_type") == "opening_balance" {
			componentCodes[stringField(record, "material_code")] = true
		}
	}
	materials, ok := sets["material"]
	if !ok {
		return Result{}, fmt.Errorf("material record set is required")
	}
	for _, record := range materials.Records {
		code := stringField(record, "material_code")
		if !componentCodes[code] {
			continue
		}
		result.add("material", materials, record, []string{"material_code"},
			map[string]string{"material_code": "code", "material_name": "name"})
	}

	for _, record := range boms.Records {
		result.add("bom", boms, record,
			[]string{"parent_material_code", "component_material_code"},
			map[string]string{
				"parent_material_code":    "parent_material_code",
				"component_material_code": "component_material_code",
				"qty_per":                 "quantity_required",
			})
	}
	for _, record := range inventories.Records {
		if stringField(record, "transaction_type") != "opening_balance" {
			continue
		}
		result.add("inventory", inventories, record,
			[]string{"material_code", "warehouse_code", "lot_no"},
			map[string]string{
				"material_code": "material_code", "warehouse_code": "warehouse_name",
				"qty": "quantity", "lot_no": "batch_no",
			})
	}

	result.Request.Objects = append(result.Request.Objects, merged)
	for _, record := range orders.Records {
		if orderBelongsToTracer(record) {
			result.warn(orders, record, map[string]string{
				"order_no": "order_no", "customer_code": "customer_code", "material_code": "lines.material_code",
				"order_qty": "lines.quantity", "due_date": "required_date",
			})
		}
	}
	return result, nil
}

func mergeTracerOrder(set scenariopack.RecordSet) (Object, map[string]bool, string, error) {
	var records []map[string]any
	var customer, material, due string
	var quantity int64
	for _, record := range set.Records {
		if !orderBelongsToTracer(record) {
			continue
		}
		c, m, d := stringField(record, "customer_code"), stringField(record, "material_code"), stringField(record, "due_date")
		if c == "" || m == "" || d == "" {
			return Object{}, nil, "", fmt.Errorf("sales_order %s: customer_code, material_code and due_date are required", stringField(record, "order_no"))
		}
		if len(records) > 0 && (c != customer || m != material || d != due) {
			return Object{}, nil, "", fmt.Errorf("tracer sales orders must share customer, material and due_date")
		}
		q, err := integer(record["order_qty"])
		if err != nil || q <= 0 {
			return Object{}, nil, "", fmt.Errorf("sales_order %s: order_qty must be a positive integer", stringField(record, "order_no"))
		}
		customer, material, due = c, m, d
		quantity += q
		records = append(records, record)
	}
	if len(records) != 2 || quantity != 12000 {
		return Object{}, nil, "", fmt.Errorf("tracer requires two source sales orders totalling 12000; got %d records totalling %d", len(records), quantity)
	}
	data := map[string]any{
		"order_no": TracerOrderNo, "customer_code": customer, "required_date": shanghaiMidnight(due), "status": "draft",
		"lines": []map[string]any{{"material_code": material, "quantity": quantity, "unit_price": DemoUnitPrice}},
	}
	return Object{
		Type: "sales_order", NaturalKey: map[string]any{"order_no": TracerOrderNo}, Data: data,
		IdempotencyKey: "hctm:sales_order:" + TracerOrderNo + ":12000",
	}, map[string]bool{material: true}, customer, nil
}

func orderBelongsToTracer(record map[string]any) bool {
	return stringField(record, "order_no") == TracerOrderNo || stringField(record, "original_order_no") == TracerOrderNo
}

func (r *Result) add(target string, set scenariopack.RecordSet, record map[string]any, keyFields []string, mapping map[string]string) {
	data := make(map[string]any, len(mapping))
	for source, destination := range mapping {
		if value, ok := record[source]; ok {
			data[destination] = value
		}
	}
	naturalKey := make(map[string]any, len(keyFields))
	for _, source := range keyFields {
		if value, ok := record[source]; ok {
			naturalKey[mapping[source]] = value
		}
	}
	r.Request.Objects = append(r.Request.Objects, Object{
		Type: target, NaturalKey: naturalKey, Data: data,
		IdempotencyKey: idempotencyKey(target, naturalKey),
	})
	r.warn(set, record, mapping)
}

func shanghaiMidnight(date string) string {
	if _, err := time.Parse(time.DateOnly, date); err != nil {
		return date
	}
	return date + "T00:00:00+08:00"
}

func (r *Result) warn(set scenariopack.RecordSet, record map[string]any, retained map[string]string) {
	var dropped []string
	for field := range record {
		if _, ok := retained[field]; !ok {
			dropped = append(dropped, field)
		}
	}
	if len(dropped) == 0 {
		return
	}
	sort.Strings(dropped)
	r.Warnings = append(r.Warnings, Warning{
		SourceEntity: set.Entity, NaturalKey: fields(record, set.NaturalKey), DroppedFields: dropped,
		Message: "fields are not represented by the IAOS M3 legacy contract",
	})
}

func indexSets(pack *scenariopack.Pack, story scenariopack.Story) map[string]scenariopack.RecordSet {
	sets := make(map[string]scenariopack.RecordSet)
	for _, set := range pack.RecordSets {
		sets[set.Entity] = set
	}
	for _, set := range story.Initial.RecordSets {
		sets[set.Entity] = set
	}
	return sets
}

func findStory(pack *scenariopack.Pack, key string) (scenariopack.Story, error) {
	if key == "" && len(pack.Stories) == 1 {
		return pack.Stories[0], nil
	}
	for _, story := range pack.Stories {
		if story.Ref.Key == key {
			return story, nil
		}
	}
	return scenariopack.Story{}, fmt.Errorf("story %q not found", key)
}

func recordBy(set scenariopack.RecordSet, field, value string) (map[string]any, error) {
	for _, record := range set.Records {
		if stringField(record, field) == value {
			return record, nil
		}
	}
	return nil, fmt.Errorf("%s %s=%q not found", set.Entity, field, value)
}

func fields(record map[string]any, names []string) map[string]any {
	out := make(map[string]any, len(names))
	for _, name := range names {
		if value, ok := record[name]; ok {
			out[name] = value
		}
	}
	return out
}

func idempotencyKey(entity string, key map[string]any) string {
	names := make([]string, 0, len(key))
	for name := range key {
		names = append(names, name)
	}
	sort.Strings(names)
	parts := []string{"hctm", entity}
	for _, name := range names {
		parts = append(parts, name+"="+fmt.Sprint(key[name]))
	}
	return strings.Join(parts, ":")
}

func stringField(record map[string]any, field string) string {
	value, _ := record[field].(string)
	return value
}

func integer(value any) (int64, error) {
	switch v := value.(type) {
	case json.Number:
		return v.Int64()
	case float64:
		if v != float64(int64(v)) {
			return 0, fmt.Errorf("not an integer")
		}
		return int64(v), nil
	case int:
		return int64(v), nil
	case int64:
		return v, nil
	default:
		return 0, fmt.Errorf("not numeric")
	}
}
