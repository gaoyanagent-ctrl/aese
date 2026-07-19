package validate

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/industrial-ai/iaos-aese/internal/scenariopack"
)

type Issue struct {
	File    string `json:"file"`
	Record  string `json:"record,omitempty"`
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

func (i Issue) Error() string {
	where := i.File
	if i.Record != "" {
		where += ": record " + i.Record
	}
	if i.Field != "" {
		where += ": field " + i.Field
	}
	return where + ": " + i.Message
}

type Result struct {
	Issues []Issue `json:"issues"`
}

func (r Result) Valid() bool { return len(r.Issues) == 0 }

func Pack(pack *scenariopack.Pack) Result {
	var result Result
	add := func(file, record, field, message string) {
		result.Issues = append(result.Issues, Issue{file, record, field, message})
	}
	if pack.Manifest.PackKey == "" {
		add("manifest.json", "", "pack_key", "is required")
	}
	if pack.Manifest.PackVersion == "" {
		add("manifest.json", "", "pack_version", "is required")
	}
	if pack.Manifest.TenantTemplate == "" {
		add("manifest.json", "", "tenant_template", "is required")
	}

	sets := append([]scenariopack.RecordSet(nil), pack.RecordSets...)
	for _, story := range pack.Stories {
		sets = append(sets, story.Initial.RecordSets...)
	}
	indexes := map[string]map[string]struct{}{}
	for _, set := range sets {
		if set.Source == "" {
			set.Source = "initial-state.json"
		}
		validateSet(set, indexes, add)
	}
	validateReferences(sets, indexes, add)
	for _, story := range pack.Stories {
		validateStory(story, sets, add)
	}
	sort.SliceStable(result.Issues, func(i, j int) bool { return result.Issues[i].Error() < result.Issues[j].Error() })
	return result
}

func validateSet(set scenariopack.RecordSet, indexes map[string]map[string]struct{}, add func(string, string, string, string)) {
	file := set.Source
	if set.SchemaVersion != scenariopack.SupportedSchemaVersion {
		add(file, "", "schema_version", fmt.Sprintf("unsupported %q", set.SchemaVersion))
	}
	if set.Entity == "" {
		add(file, "", "entity", "is required")
		return
	}
	if len(set.NaturalKey) == 0 {
		add(file, set.Entity, "natural_key", "must contain at least one field")
		return
	}
	if indexes[set.Entity] == nil {
		indexes[set.Entity] = map[string]struct{}{}
	}
	for n, record := range set.Records {
		label := fmt.Sprintf("%s[%d]", set.Entity, n)
		parts := make([]string, 0, len(set.NaturalKey))
		missing := false
		for _, key := range set.NaturalKey {
			value, ok := record[key]
			if !ok || fmt.Sprint(value) == "" {
				add(file, label, key, "natural key value is required")
				missing = true
				continue
			}
			parts = append(parts, fmt.Sprint(value))
		}
		if !missing {
			key := strings.Join(parts, "\x1f")
			if _, exists := indexes[set.Entity][key]; exists {
				add(file, label, strings.Join(set.NaturalKey, ","), "duplicate natural key "+strings.Join(parts, "/"))
			}
			indexes[set.Entity][key] = struct{}{}
			label = set.Entity + ":" + strings.Join(parts, "/")
		}
		for field, value := range record {
			if quantityField(field) {
				if number, ok := asFloat(value); ok && number < 0 {
					add(file, label, field, "quantity must be non-negative")
				}
			}
		}
	}
}

type reference struct{ field, target string }

var references = map[string][]reference{
	"business_unit":         {{"group_code", "enterprise_group"}},
	"legal_entity":          {{"group_code", "enterprise_group"}},
	"plant":                 {{"business_unit_code", "business_unit"}, {"legal_entity_code", "legal_entity"}},
	"department":            {{"plant_code", "plant"}},
	"production_team":       {{"plant_code", "plant"}, {"department_code", "department"}, {"shift_code", "shift"}},
	"bom":                   {{"parent_material_code", "material"}, {"component_material_code", "material"}},
	"routing":               {{"material_code", "material"}, {"plant_code", "plant"}},
	"operation":             {{"routing_code", "routing"}, {"work_center_code", "work_center"}},
	"work_center":           {{"plant_code", "plant"}},
	"equipment":             {{"work_center_code", "work_center"}},
	"storage_location":      {{"warehouse_code", "warehouse"}},
	"sales_order":           {{"customer_code", "customer"}, {"legal_entity_code", "legal_entity"}, {"material_code", "material"}},
	"purchase_order":        {{"supplier_code", "supplier"}, {"material_code", "material"}, {"plant_code", "plant"}},
	"inspection_order":      {{"po_no", "purchase_order"}, {"material_code", "material"}},
	"production_order":      {{"sales_order_no", "sales_order"}, {"material_code", "material"}, {"routing_code", "routing"}, {"plant_code", "plant"}},
	"shipment":              {{"sales_order_no", "sales_order"}, {"material_code", "material"}},
	"inventory_transaction": {{"material_code", "material"}, {"warehouse_code", "warehouse"}, {"storage_location_code", "storage_location"}},
}

func validateReferences(sets []scenariopack.RecordSet, indexes map[string]map[string]struct{}, add func(string, string, string, string)) {
	for _, set := range sets {
		for n, record := range set.Records {
			label := fmt.Sprintf("%s[%d]", set.Entity, n)
			for _, ref := range references[set.Entity] {
				value, exists := record[ref.field]
				if !exists || fmt.Sprint(value) == "" {
					continue
				}
				if _, ok := indexes[ref.target][fmt.Sprint(value)]; !ok {
					add(set.Source, label, ref.field, fmt.Sprintf("references missing %s %q", ref.target, value))
				}
			}
		}
	}
}

func validateStory(story scenariopack.Story, sets []scenariopack.RecordSet, add func(string, string, string, string)) {
	if story.Initial.SchemaVersion != scenariopack.SupportedSchemaVersion {
		add(story.Initial.Source, "", "schema_version", fmt.Sprintf("unsupported %q", story.Initial.SchemaVersion))
	}
	if story.Events.SchemaVersion != scenariopack.SupportedSchemaVersion {
		add(story.Events.Source, "", "schema_version", fmt.Sprintf("unsupported %q", story.Events.SchemaVersion))
	}
	if story.Expected.SchemaVersion != scenariopack.SupportedSchemaVersion {
		add(story.Expected.Source, "", "schema_version", fmt.Sprintf("unsupported %q", story.Expected.SchemaVersion))
	}
	if story.Ref.Key == "" {
		add("manifest.json", "", "stories.key", "is required")
	}
	if story.Initial.StoryKey != story.Ref.Key {
		add(story.Initial.Source, "", "story_key", fmt.Sprintf("%q does not match manifest key %q", story.Initial.StoryKey, story.Ref.Key))
	}
	if story.Events.StoryKey != story.Ref.Key {
		add(story.Events.Source, "", "story_key", fmt.Sprintf("%q does not match manifest key %q", story.Events.StoryKey, story.Ref.Key))
	}
	if story.Expected.StoryKey != story.Ref.Key {
		add(story.Expected.Source, "", "story_key", fmt.Sprintf("%q does not match manifest key %q", story.Expected.StoryKey, story.Ref.Key))
	}
	ids, keys := map[string]struct{}{}, map[string]struct{}{}
	var previous time.Time
	correlation := story.Events.CorrelationID
	for n, event := range story.Events.Events {
		label := fmt.Sprintf("events[%d]", n)
		if event.EventID == "" {
			add(story.Events.Source, label, "event_id", "is required")
		} else if _, ok := ids[event.EventID]; ok {
			add(story.Events.Source, label, "event_id", "duplicate event ID "+event.EventID)
		}
		idempotency := event.Idempotency()
		if idempotency == "" {
			add(story.Events.Source, label, "idempotency_key", "is required")
		} else if _, ok := keys[idempotency]; ok {
			add(story.Events.Source, label, "idempotency_key", "duplicate idempotency key "+idempotency)
		}
		if event.Correlation() == "" {
			add(story.Events.Source, label, "correlation_id", "is required")
		} else if correlation == "" {
			correlation = event.Correlation()
		} else if event.Correlation() != correlation {
			add(story.Events.Source, label, "correlation_id", fmt.Sprintf("%q does not match story correlation %q", event.Correlation(), correlation))
		}
		at, err := event.Time()
		if err != nil {
			add(story.Events.Source, label, "timestamp", "must be RFC3339: "+err.Error())
		} else if !previous.IsZero() && at.Before(previous) && event.ConcurrentGroup == "" {
			add(story.Events.Source, label, "timestamp", "event timeline is out of order")
		} else {
			previous = at
		}
		if cause := event.Causation(); cause != "" {
			if _, ok := ids[cause]; !ok {
				add(story.Events.Source, label, "causation_id", "must reference an earlier event: "+cause)
			}
		}
		if event.EventType == "" {
			add(story.Events.Source, label, "event_type", "is required")
		}
		if event.Payload == nil {
			add(story.Events.Source, label, "payload", "is required")
		}
		ids[event.EventID], keys[idempotency] = struct{}{}, struct{}{}
	}
	if len(story.Expected.Assertions) == 0 {
		add(story.Expected.Source, "", "assertions", "must contain at least one machine assertion")
	}
	validateExpectedOutcomes(story, sets, add)
	validateMRPInvariant(story, sets, add)
	validateShipmentInvariant(story, sets, add)
}

func validateExpectedOutcomes(story scenariopack.Story, sets []scenariopack.RecordSet, add func(string, string, string, string)) {
	values := map[string]any{}
	for i, a := range story.Expected.Assertions {
		var actual any
		switch a.Type {
		case "event_count":
			actual = len(story.Events.Events)
		case "event_sequence":
			if a.Field == "correlation_id" {
				actual = story.Events.CorrelationID
			}
		case "record_aggregate":
			total := 0.0
			orderedMaterial := ""
			for _, set := range sets {
				if set.Entity == "sales_order" && len(set.Records) > 0 {
					orderedMaterial = fmt.Sprint(set.Records[0]["material_code"])
				}
			}
			for _, set := range sets {
				if set.Entity != a.Entity {
					continue
				}
				for _, record := range set.Records {
					if a.Operator == "sum_where_material_equals" && fmt.Sprint(record["material_code"]) != orderedMaterial {
						continue
					}
					n, _ := asFloat(record[a.Field])
					total += n
				}
			}
			actual = total
		case "event_payload":
			for _, event := range story.Events.Events {
				if event.EventType == a.Entity {
					if value, ok := payloadPath(event.Payload, a.Field); ok {
						actual = value
					}
				}
			}
		case "event_aggregate":
			total := 0.0
			for _, event := range story.Events.Events {
				if event.EventType == a.Entity {
					if value, ok := payloadPath(event.Payload, a.Field); ok {
						n, _ := asFloat(value)
						total += n
					}
				}
			}
			actual = total
		case "derived":
			switch a.Field {
			case "opening_finished_goods + finished_goods_received":
				actual = sumValues(values["opening_finished_goods"], values["finished_goods_received_qty"])
			case "total_customer_demand - actual_shipped_qty":
				actual = subValues(values["total_customer_demand"], values["actual_shipped_qty"])
			case "sales_order.delivery_status":
				actual = story.Expected.Summary["delivery_status"]
			}
		case "invariant":
			if a.Field == "actual_shipped_qty <= available_to_ship_qty" {
				left, _ := asFloat(values["actual_shipped_qty"])
				right, _ := asFloat(values["available_to_ship_qty"])
				actual = left <= right
			}
		}
		values[a.Key] = actual
		if actual == nil {
			add(story.Expected.Source, fmt.Sprintf("assertions[%d]", i), "field", "cannot evaluate "+a.Field)
			continue
		}
		if !equalValue(actual, a.Expected) {
			add(story.Expected.Source, fmt.Sprintf("assertions[%d]", i), "expected", fmt.Sprintf("assertion %s evaluated to %v, expected %v", a.Key, actual, a.Expected))
		}
	}
}

func payloadPath(payload map[string]any, path string) (any, bool) {
	parts := strings.Split(path, ".")
	if len(parts) == 1 {
		value, ok := payload[path]
		return value, ok
	}
	if parts[0] == "material_shortages" && len(parts) == 3 {
		items, _ := payload[parts[0]].([]any)
		for _, item := range items {
			row, _ := item.(map[string]any)
			if fmt.Sprint(row["material_code"]) == parts[1] {
				value, ok := row[parts[2]]
				return value, ok
			}
		}
	}
	return nil, false
}
func sumValues(a, b any) float64 { x, _ := asFloat(a); y, _ := asFloat(b); return x + y }
func subValues(a, b any) float64 { x, _ := asFloat(a); y, _ := asFloat(b); return x - y }
func equalValue(a, b any) bool {
	if x, ok := asFloat(a); ok {
		if y, ok := asFloat(b); ok {
			return abs(x-y) < 0.0001
		}
	}
	return fmt.Sprint(a) == fmt.Sprint(b)
}

func validateMRPInvariant(story scenariopack.Story, sets []scenariopack.RecordSet, add func(string, string, string, string)) {
	bom := map[string]float64{}
	for _, set := range sets {
		if set.Entity == "bom" {
			for _, record := range set.Records {
				qty, ok := asFloat(record["qty_per"])
				component := fmt.Sprint(record["component_material_code"])
				if !ok || qty <= 0 {
					add(set.Source, "bom:"+component, "qty_per", "must be greater than zero")
				} else {
					bom[component] = qty
				}
				if fmt.Sprint(record["parent_material_code"]) == component {
					add(set.Source, "bom:"+component, "component_material_code", "BOM component cannot equal parent")
				}
			}
		}
	}
	for n, event := range story.Events.Events {
		if !strings.Contains(event.EventType, "mrp.generated") {
			continue
		}
		label := fmt.Sprintf("events[%d]", n)
		demand, demandOK := firstNumber(event.Payload, "demand_qty", "order_qty")
		available, availableOK := firstNumber(event.Payload, "available_finished_goods_qty", "available_qty")
		net, netOK := firstNumber(event.Payload, "net_production_qty", "net_qty")
		if demandOK && availableOK && netOK && abs((demand-available)-net) > 0.0001 {
			add(story.Events.Source, label, "payload.net_production_qty", fmt.Sprintf("must equal demand minus available finished goods (%.3f)", demand-available))
		}
		shortages, _ := event.Payload["material_shortages"].([]any)
		for i, item := range shortages {
			row, ok := item.(map[string]any)
			if !ok {
				continue
			}
			component := fmt.Sprint(row["material_code"])
			required, ok := asFloat(row["required_qty"])
			usage, known := bom[component]
			if demandOK && ok && known {
				expected := demand * usage
				if abs(required-expected) > 0.0001 {
					add(story.Events.Source, label, fmt.Sprintf("payload.material_shortages[%d].required_qty", i), fmt.Sprintf("%.3f does not equal demand %.3f × BOM qty_per %.4f", required, demand, usage))
				}
			}
		}
	}
}

func validateShipmentInvariant(story scenariopack.Story, sets []scenariopack.RecordSet, add func(string, string, string, string)) {
	available, shipped := map[string]float64{}, map[string]float64{}
	for _, set := range sets {
		if set.Entity != "inventory_transaction" {
			continue
		}
		for _, record := range set.Records {
			if fmt.Sprint(record["direction"]) == "in" {
				q, _ := asFloat(record["quantity"])
				available[fmt.Sprint(record["material_code"])] += q
			}
		}
	}
	for n, event := range story.Events.Events {
		q, _ := firstNumber(event.Payload, "shipped_quantity", "ship_qty", "quantity", "received_qty", "received_quantity", "completed_qty", "completed_quantity", "shipment_quantity")
		material := fmt.Sprint(event.Payload["material_code"])
		if strings.Contains(event.EventType, "finished_goods") && strings.Contains(event.EventType, "received") {
			available[material] += q
		}
		if strings.Contains(event.EventType, "shipment") && strings.Contains(event.EventType, "dispatched") {
			requested, requestedOK := firstNumber(event.Payload, "requested_quantity", "requested_qty")
			shortage, shortageOK := firstNumber(event.Payload, "shortage_quantity", "shortage_qty")
			if requestedOK && shortageOK && abs(requested-(q+shortage)) > 0.0001 {
				add(story.Events.Source, fmt.Sprintf("events[%d]", n), "payload.shortage_quantity", "requested quantity must equal shipped plus shortage")
			}
			shipped[material] += q
			if shipped[material] > available[material] {
				add(story.Events.Source, fmt.Sprintf("events[%d]", n), "payload.shipped_quantity", fmt.Sprintf("cumulative shipment %.3f exceeds available finished goods %.3f", shipped[material], available[material]))
			}
		}
	}
}

func quantityField(field string) bool {
	field = strings.ToLower(field)
	return strings.Contains(field, "quantity") || strings.HasSuffix(field, "_qty") || strings.Contains(field, "stock")
}
func asFloat(value any) (float64, bool) {
	switch v := value.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case json.Number:
		n, e := v.Float64()
		return n, e == nil
	case string:
		n, e := strconv.ParseFloat(v, 64)
		return n, e == nil
	}
	return 0, false
}
func firstNumber(values map[string]any, keys ...string) (float64, bool) {
	for _, key := range keys {
		if n, ok := asFloat(values[key]); ok {
			return n, true
		}
	}
	return 0, false
}
func abs(value float64) float64 {
	if value < 0 {
		return -value
	}
	return value
}
