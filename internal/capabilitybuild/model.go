package capabilitybuild

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type Money struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
	Scale    int    `json:"scale"`
}
type Option struct {
	Code     string   `json:"code"`
	Mode     string   `json:"mode"`
	Cash     Money    `json:"cash"`
	LeadDays int      `json:"lead_days"`
	Risk     int      `json:"risk"`
	Score    int      `json:"score"`
	Feasible bool     `json:"feasible"`
	Failures []string `json:"failures"`
}
type Equipment struct {
	Code     string `json:"code"`
	Zone     string `json:"zone"`
	Status   string `json:"status"`
	PowerKVA int    `json:"power_kva"`
	Evidence string `json:"evidence"`
}
type Worker struct {
	Code           string   `json:"code"`
	Position       string   `json:"position"`
	Status         string   `json:"status"`
	Skills         []string `json:"skills"`
	QualifiedUntil string   `json:"qualified_until"`
	Backup         bool     `json:"backup"`
}
type Knowledge struct {
	Actor      string `json:"actor"`
	Fact       string `json:"fact"`
	ObservedAt string `json:"observed_at"`
	Visibility string `json:"visibility"`
}
type Frame struct {
	Step                      int             `json:"step"`
	Phase                     string          `json:"phase"`
	SimTime                   string          `json:"sim_time"`
	Title                     string          `json:"title"`
	CausationID               string          `json:"causation_id"`
	Cash                      Money           `json:"cash"`
	Committed                 Money           `json:"committed"`
	Paid                      Money           `json:"paid"`
	WageReserve               Money           `json:"wage_reserve"`
	FacilityPayable           Money           `json:"facility_payable"`
	Equipment                 []Equipment     `json:"equipment"`
	Workers                   []Worker        `json:"workers"`
	Knowledge                 []Knowledge     `json:"knowledge"`
	WorldProgress             int             `json:"world_progress"`
	IAOSProgress              int             `json:"iaos_progress"`
	Discrepancy               string          `json:"discrepancy"`
	Gate                      map[string]bool `json:"gate"`
	IndustrializationEligible bool            `json:"industrialization_eligible"`
	IAOSCursor                int64           `json:"iaos_cursor"`
}
type Trace struct {
	SchemaVersion   string   `json:"schema_version"`
	Campaign        string   `json:"campaign"`
	WorldRunID      string   `json:"world_run_id"`
	Timezone        string   `json:"timezone"`
	PolicyVersion   string   `json:"policy_version"`
	M10TerminalHash string   `json:"m10_terminal_hash"`
	Options         []Option `json:"options"`
	Frames          []Frame  `json:"frames"`
}

func m(v string) Money { return Money{v, "CNY", 2} }
func Options() []Option {
	return []Option{{"ACQ-PURCHASE", "purchase", m("18000000.00"), 150, 35, 82, false, []string{"cash_buffer_breached"}}, {"ACQ-LEASE-MIX", "finance_lease", m("10500000.00"), 120, 22, 91, true, []string{}}, {"ACQ-VENDOR-CREDIT", "vendor_credit", m("12500000.00"), 180, 45, 76, false, []string{"available_after_deadline"}}}
}
func equipment() []Equipment {
	codes := []string{"EQ-FORM-01", "EQ-CNC-01", "EQ-LASER-WELD-01", "EQ-WASH-01", "EQ-LEAK-TEST-01", "EQ-ASSEMBLY-01", "LAB-QUALITY-01"}
	out := make([]Equipment, 0, len(codes))
	for _, c := range codes {
		zone := "ZONE-HCTM-SZ-PRODUCTION"
		if c == "LAB-QUALITY-01" {
			zone = "ZONE-HCTM-SZ-QUALITY"
		}
		out = append(out, Equipment{c, zone, "accepted", 120, "commissioning-calibration-safety-accepted"})
	}
	return out
}
func workers() []Worker {
	roles := []string{"PLANT-MANAGER", "PLANNING", "PROCUREMENT", "QUALITY", "PROCESS", "EQUIPMENT", "WAREHOUSE", "OPERATOR-A", "OPERATOR-B", "INSPECTOR"}
	out := make([]Worker, 0, len(roles))
	for i, r := range roles {
		out = append(out, Worker{fmt.Sprintf("WORKER-HCTM-%02d", i+1), "POS-" + r, "onboarded", []string{"SAFETY", "ROLE-PRACTICAL"}, "2027-12-31T23:59:59+08:00", r == "OPERATOR-B"})
	}
	return out
}
func BuildTrace() Trace {
	base := Frame{0, "eligible", "2026-08-21T09:00:00+08:00", "消费 M10 设施资格", "m10-terminal", m("10000000.00"), m("0.00"), m("0.00"), m("3000000.00"), m("3500000.00"), []Equipment{}, []Worker{}, []Knowledge{}, 0, 0, "none", map[string]bool{}, false, 300}
	frames := []Frame{base}
	add := func(phase, at, title, cause string, world, iaos int) {
		f := frames[len(frames)-1]
		f.Step = len(frames)
		f.Phase = phase
		f.SimTime = at
		f.Title = title
		f.CausationID = cause
		f.WorldProgress = world
		f.IAOSProgress = iaos
		frames = append(frames, f)
	}
	add("funded", "2026-08-25T10:00:00+08:00", "剩余认缴资本实际到账并受治理入账", "capital-received", 5, 5)
	frames[1].Cash = m("20000000.00")
	add("sourcing", "2026-08-28T10:00:00+08:00", "采购/租赁/账期方案完成硬约束比选", "acquisition-evaluated", 10, 10)
	add("ordered", "2026-09-02T10:00:00+08:00", "融资租赁组合订单获批", "equipment-order-approved", 20, 20)
	frames[3].Committed = m("14000000.00")
	add("installing", "2026-12-15T10:00:00+08:00", "设备到货落位安装，空间与 utility 守恒", "equipment-installed", 55, 60)
	add("commissioning", "2026-12-20T10:00:00+08:00", "检漏设备校准漂移，计划与现实产生差异", "leak-calibration-failed", 62, 80)
	frames[5].Discrepancy = "leak_calibration_drift_open"
	add("commissioning", "2026-12-21T10:00:00+08:00", "设备与质量负责人收到 observation 并申请整改", "remediation-intent", 62, 80)
	frames[6].Knowledge = []Knowledge{{"ACTOR-EQUIPMENT", "leak-calibration-drift", "2026-12-21T10:00:00+08:00", "equipment-quality"}, {"ACTOR-QUALITY", "leak-calibration-drift", "2026-12-21T10:00:00+08:00", "equipment-quality"}}
	add("staffing", "2027-01-15T10:00:00+08:00", "整改复验通过，核心团队到岗并完成实操资格", "remediation-accepted", 92, 92)
	frames[7].Equipment = equipment()
	frames[7].Workers = workers()
	frames[7].Discrepancy = "closed"
	add("capability_acceptance", "2027-01-20T10:00:00+08:00", "设备、实验室、仓储、人员、班次与安全联合验收", "capability-accepted", 100, 100)
	frames[8].Equipment = equipment()
	frames[8].Workers = workers()
	frames[8].Cash = m("8500000.00")
	frames[8].Committed = m("14000000.00")
	frames[8].Paid = m("11500000.00")
	frames[8].Gate = map[string]bool{"funding": true, "facility_balance": true, "equipment": true, "laboratory": true, "warehouse": true, "workforce": true, "qualification": true, "shift": true, "safety": true}
	add("industrialization_eligible", "2027-01-20T16:00:00+08:00", "通用生产能力就绪，可进入 M12 产品工业化", "industrialization-eligibility", 100, 100)
	frames[9].Equipment = equipment()
	frames[9].Workers = workers()
	frames[9].IndustrializationEligible = true
	return Trace{"1.0", "capability-build", "world-run-genesis-capability-build-001", "Asia/Shanghai", "capability-build-policy-v1", "sha256:m10-plant-build-terminal", Options(), frames}
}
func Validate(t Trace) error {
	want := []string{"eligible", "funded", "sourcing", "ordered", "installing", "commissioning", "commissioning", "staffing", "capability_acceptance", "industrialization_eligible"}
	if t.M10TerminalHash == "" || len(t.Frames) != len(want) {
		return errors.New("invalid terminal or frame count")
	}
	for i, f := range t.Frames {
		if f.Phase != want[i] {
			return fmt.Errorf("phase %d", i)
		}
		if _, e := time.Parse(time.RFC3339, f.SimTime); e != nil {
			return e
		}
	}
	if !t.Frames[9].IndustrializationEligible || len(t.Frames[9].Equipment) != 7 || len(t.Frames[9].Workers) != 10 {
		return errors.New("joint gate incomplete")
	}
	return nil
}
func Hash(t Trace) (string, error) {
	b, e := json.Marshal(t)
	if e != nil {
		return "", e
	}
	s := sha256.Sum256(b)
	return hex.EncodeToString(s[:]), nil
}
func Snapshot(t Trace) ([]byte, error) {
	if e := Validate(t); e != nil {
		return nil, e
	}
	return json.Marshal(t)
}
func Restore(b []byte) (Trace, error) {
	var t Trace
	d := json.NewDecoder(bytes.NewReader(b))
	d.DisallowUnknownFields()
	e := d.Decode(&t)
	if e == nil {
		e = Validate(t)
	}
	return t, e
}
