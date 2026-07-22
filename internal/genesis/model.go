// Package genesis implements the narrow LAS-WLD-02 three-state tracer.
package genesis

import (
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/industrial-ai/iaos-aese/internal/knowledge"
	"github.com/industrial-ai/iaos-aese/internal/worldcontract"
)

type EquipmentState struct {
	Code              string `json:"code"`
	Condition         string `json:"condition"`
	Vibration         string `json:"vibration"`
	Unit              string `json:"unit"`
	AvailableCapacity string `json:"available_capacity"`
}
type IAOSProjection struct {
	EquipmentCode    string `json:"equipment_code"`
	RegisteredStatus string `json:"registered_status"`
	MaintenanceOrder string `json:"maintenance_order,omitempty"`
	Cursor           int64  `json:"cursor"`
}
type Frame struct {
	Step        int                       `json:"step"`
	SimTime     string                    `json:"sim_time"`
	Title       string                    `json:"title"`
	World       EquipmentState            `json:"world"`
	IAOS        IAOSProjection            `json:"iaos"`
	Knowledge   []worldcontract.Knowledge `json:"knowledge"`
	Discrepancy worldcontract.Discrepancy `json:"discrepancy"`
	CausationID string                    `json:"causation_id"`
}
type Trace struct {
	SchemaVersion string                  `json:"schema_version"`
	WorldRunID    string                  `json:"world_run_id"`
	Timezone      string                  `json:"timezone"`
	ActorRef      worldcontract.StableRef `json:"actor_ref"`
	Frames        []Frame                 `json:"frames"`
}

func BuildTrace() Trace {
	actor := worldcontract.StableRef{Namespace: "hctm", Type: "position", Code: "EQUIPMENT-ENGINEER-SZ"}
	subject := worldcontract.StableRef{Namespace: "hctm", Type: "equipment", Code: "LAS-WLD-02"}
	base := worldcontract.Discrepancy{SchemaVersion: "1.0", DiscrepancyID: "disc-las-wld-02-001", TenantID: "tenant-hctm", WorldRunID: "world-run-genesis-001", BranchID: "main", SubjectRef: subject, Kind: "world_vs_iaos", WorldFactRef: "world-event-vibration-rise-01", Status: "open", DetectedAt: "2026-07-08T10:15:00+08:00", CorrelationID: "corr-equipment-degradation-01"}
	normal := EquipmentState{"LAS-WLD-02", "normal", "5.80", "mm/s", "1.00"}
	degraded := EquipmentState{"LAS-WLD-02", "degraded", "7.20", "mm/s", "0.70"}
	critical := EquipmentState{"LAS-WLD-02", "critical", "8.10", "mm/s", "0.00"}
	iaosNormal := IAOSProjection{"LAS-WLD-02", "normal", "", 40}
	unknown := Frame{0, "2026-07-08T10:15:00+08:00", "世界已退化，系统与角色尚未知", degraded, iaosNormal, []worldcontract.Knowledge{}, base, "world-event-vibration-rise-01"}
	k := worldcontract.Knowledge{SchemaVersion: "1.0", KnowledgeID: "knowledge-engineer-vibration-001", TenantID: "tenant-hctm", WorldRunID: "world-run-genesis-001", BranchID: "main", ActorRef: actor, FactRef: worldcontract.StableRef{Namespace: "hctm", Type: "world_event", Code: "world-event-vibration-rise-01"}, ObservedAt: "2026-07-08T10:16:00+08:00", ValidAt: "2026-07-08T10:15:00+08:00", SourceRef: "sensor:VIB-LAS-WLD-02-01", Confidence: "0.95", VisibilityScope: "assigned_recipients"}
	investigating := base
	investigating.Status = "investigating"
	investigating.KnowledgeRef = k.KnowledgeID
	closed := investigating
	closed.Status = "closed"
	closed.IAOSRecordRef = "maintenance_work_order:EAM-WO-20260708-001"
	closed.ClosedAt = "2026-07-08T11:05:00+08:00"
	return Trace{"1.0", "world-run-genesis-001", "Asia/Shanghai", actor, []Frame{unknown, {1, "2026-07-08T10:16:00+08:00", "传感观察使设备工程师获知", degraded, iaosNormal, []worldcontract.Knowledge{k}, investigating, "obs-vibration-001"}, {2, "2026-07-08T10:21:00+08:00", "IAOS 已提交检修工单", critical, IAOSProjection{"LAS-WLD-02", "inspection_scheduled", "EAM-WO-20260708-001", 43}, []worldcontract.Knowledge{k}, investigating, "outcome-inspection-001"}, {3, "2026-07-08T11:05:00+08:00", "检修后世界结果与管理记录对账关闭", normal, IAOSProjection{"LAS-WLD-02", "available", "EAM-WO-20260708-001", 44}, []worldcontract.Knowledge{k}, closed, "world-event-repair-completed-01"}}}
}
func ValidateTrace(trace Trace) error {
	if trace.SchemaVersion != "1.0" || trace.Timezone != "Asia/Shanghai" || len(trace.Frames) < 4 {
		return fmt.Errorf("invalid genesis trace")
	}
	store := knowledge.New()
	var previous time.Time
	for i, frame := range trace.Frames {
		at, err := time.Parse(time.RFC3339, frame.SimTime)
		if err != nil {
			return err
		}
		if i > 0 && at.Before(previous) {
			return fmt.Errorf("frame time moved backwards")
		}
		previous = at
		for _, k := range frame.Knowledge {
			if err := store.Learn(k); err != nil && i == 1 {
				return err
			}
		}
	}
	if len(trace.Frames[0].Knowledge) != 0 || trace.Frames[len(trace.Frames)-1].Discrepancy.Status != "closed" {
		return fmt.Errorf("three-state lifecycle incomplete")
	}
	return nil
}

type Conservation struct {
	Opening  string `json:"opening"`
	Inbound  string `json:"inbound"`
	Consumed string `json:"consumed"`
	Loss     string `json:"loss"`
	Closing  string `json:"closing"`
}

func ValidateConservation(c Conservation) error {
	parse := func(s string) (*big.Rat, error) {
		v, ok := new(big.Rat).SetString(s)
		if !ok {
			return nil, fmt.Errorf("invalid decimal %q", s)
		}
		return v, nil
	}
	o, e := parse(c.Opening)
	if e != nil {
		return e
	}
	in, _ := parse(c.Inbound)
	used, _ := parse(c.Consumed)
	loss, _ := parse(c.Loss)
	close, _ := parse(c.Closing)
	left := new(big.Rat).Add(o, in)
	right := new(big.Rat).Add(used, loss)
	right.Add(right, close)
	if left.Cmp(right) != 0 {
		return fmt.Errorf("conservation mismatch %s != %s", left.RatString(), right.RatString())
	}
	return nil
}
func JSON() []byte { data, _ := json.Marshal(BuildTrace()); return data }
