package industrialization

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
type Release struct {
	Code       string `json:"code"`
	Revision   string `json:"revision"`
	Status     string `json:"status"`
	Hash       string `json:"hash"`
	Supersedes string `json:"supersedes,omitempty"`
}
type Trial struct {
	Code                     string `json:"code"`
	Revision                 string `json:"revision"`
	Input, Good, Scrap       int
	Yield, Cpk, CycleSeconds string
	LeakFailures             int  `json:"leak_failures"`
	Traceable                bool `json:"traceable"`
}
type Knowledge struct{ Actor, Fact, ObservedAt, Visibility string }
type Frame struct {
	Step                     int             `json:"step"`
	Phase                    string          `json:"phase"`
	SimTime                  string          `json:"sim_time"`
	Title                    string          `json:"title"`
	CausationID              string          `json:"causation_id"`
	Cash                     Money           `json:"cash"`
	ContractLiability        Money           `json:"contract_liability"`
	TrialCost                Money           `json:"trial_cost"`
	Releases                 []Release       `json:"releases"`
	Trials                   []Trial         `json:"trials"`
	Knowledge                []Knowledge     `json:"knowledge"`
	APQPGates                map[string]bool `json:"apqp_gates"`
	Discrepancy              string          `json:"discrepancy"`
	WorldProgress            int             `json:"world_progress"`
	IAOSProgress             int             `json:"iaos_progress"`
	PPAPStatus               string          `json:"ppap_status"`
	Compatibility            string          `json:"compatibility"`
	SerialProductionEligible bool            `json:"serial_production_eligible"`
	IAOSCursor               int64           `json:"iaos_cursor"`
}
type Trace struct {
	SchemaVersion                   string `json:"schema_version"`
	Campaign                        string `json:"campaign"`
	WorldRunID                      string `json:"world_run_id"`
	Timezone                        string `json:"timezone"`
	PolicyVersion                   string `json:"policy_version"`
	M11TerminalHash                 string `json:"m11_terminal_hash"`
	Customer, Product, RFQ, Project string
	Frames                          []Frame `json:"frames"`
}

func m(v string) Money { return Money{v, "CNY", 2} }
func releases(rev string) []Release {
	codes := []string{"PRODUCT-HCTM-BCP-A01", "BOM-HCTM-BCP-A01-V1", "RT-HCTM-BCP-A01-V1", "PFMEA-HCTM-BCP-A01-V1", "CP-HCTM-BCP-A01-V1"}
	r := make([]Release, 0, len(codes))
	for _, c := range codes {
		h := sha256.Sum256([]byte(c + rev))
		r = append(r, Release{c, rev, "released", hex.EncodeToString(h[:]), ""})
	}
	return r
}
func BuildTrace() Trace {
	b := Frame{0, "eligible", "2027-01-21T09:00:00+08:00", "消费 M11 工业化资格", "m11-terminal", m("8500000.00"), m("0.00"), m("0.00"), []Release{}, []Trial{}, []Knowledge{}, map[string]bool{}, "none", 0, 0, "not_submitted", "pending", false, 400}
	fs := []Frame{b}
	add := func(p, at, title, c string, w, i int) {
		f := fs[len(fs)-1]
		f.Step = len(fs)
		f.Phase = p
		f.SimTime = at
		f.Title = title
		f.CausationID = c
		f.WorldProgress = w
		f.IAOSProgress = i
		fs = append(fs, f)
	}
	add("rfq", "2027-01-25T10:00:00+08:00", "收到虚构客户 RFQ 与关键要求", "rfq-received", 5, 5)
	add("nominated", "2027-02-10T10:00:00+08:00", "可行性与报价获批，客户定点及开发预付款到账", "nomination-received", 15, 15)
	fs[2].Cash = m("10500000.00")
	fs[2].ContractLiability = m("2000000.00")
	add("product_design", "2027-03-01T10:00:00+08:00", "产品与 EBOM revision A 发布", "product-revision-a", 25, 25)
	fs[3].Releases = releases("A")
	add("process_design", "2027-03-20T10:00:00+08:00", "MBOM、routing、PFMEA 与控制计划一致发布", "process-revision-a", 35, 35)
	fs[4].Releases = releases("A")
	add("supplier_tooling", "2027-05-01T10:00:00+08:00", "供应商、工装、量检具和首批可追溯物料就绪", "supplier-tooling-ready", 50, 50)
	add("trial_1", "2027-05-10T10:00:00+08:00", "首轮试制泄漏超限且 Cpk 不足", "trial-1-failed", 62, 78)
	fs[6].Trials = []Trial{{"TRIAL-HCTM-BCP-A01-T1", "A", 100, 82, 18, "82.00", "0.91", "68.00", 7, true}}
	fs[6].Discrepancy = "weld_leak_cpk_open"
	add("remediation", "2027-05-12T10:00:00+08:00", "质量 observation 送达，revision B 变更与遏制获批", "engineering-change-b", 68, 72)
	fs[7].Knowledge = []Knowledge{{"ACTOR-QUALITY", "trial-1-leak-cpk-failed", "2027-05-12T10:00:00+08:00", "project-quality"}}
	fs[7].Releases = releases("B")
	add("trial_2", "2027-06-01T10:00:00+08:00", "第二轮试制良率、节拍、MSA 与 Cpk 通过", "trial-2-passed", 90, 90)
	fs[8].Trials = []Trial{{"TRIAL-HCTM-BCP-A01-T1", "A", 100, 82, 18, "82.00", "0.91", "68.00", 7, true}, {"TRIAL-HCTM-BCP-A01-T2", "B", 120, 116, 4, "96.67", "1.67", "54.00", 0, true}}
	fs[8].Releases = releases("B")
	fs[8].Discrepancy = "closed"
	fs[8].TrialCost = m("1200000.00")
	add("ppap", "2027-06-15T10:00:00+08:00", "PPAP 包完整并由客户实际批准", "ppap-approved", 100, 100)
	fs[9].Trials = fs[8].Trials
	fs[9].Releases = releases("B")
	fs[9].APQPGates = map[string]bool{"rfq_feasibility": true, "product_design": true, "process_design": true, "supplier_tooling": true, "validation": true, "ppap": true}
	fs[9].PPAPStatus = "approved"
	fs[9].Compatibility = "hctm-stable-codes-compatible"
	add("serial_production_eligible", "2027-06-15T16:00:00+08:00", "量产质量与治理门通过，可进入 M13 正式订单", "production-release-approved", 100, 100)
	fs[10].Trials = fs[9].Trials
	fs[10].Releases = fs[9].Releases
	fs[10].SerialProductionEligible = true
	return Trace{"1.0", "industrialization", "world-run-genesis-industrialization-001", "Asia/Shanghai", "industrialization-policy-v1", "sha256:m11-capability-terminal", "CUST-SGNEV", "HCTM-BCP-A01", "RFQ-SGNEV-BCP-A01-01", "PROJECT-SGNEV-BCP-A01", fs}
}
func Validate(t Trace) error {
	want := []string{"eligible", "rfq", "nominated", "product_design", "process_design", "supplier_tooling", "trial_1", "remediation", "trial_2", "ppap", "serial_production_eligible"}
	if t.M11TerminalHash == "" || len(t.Frames) != len(want) {
		return errors.New("terminal/frame invalid")
	}
	for i, f := range t.Frames {
		if f.Phase != want[i] {
			return fmt.Errorf("phase %d", i)
		}
		if _, e := time.Parse(time.RFC3339, f.SimTime); e != nil {
			return e
		}
	}
	last := t.Frames[10]
	if !last.SerialProductionEligible || last.PPAPStatus != "approved" || len(last.APQPGates) != 6 {
		return errors.New("release gate incomplete")
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
