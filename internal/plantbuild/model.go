// Package plantbuild implements the deterministic M10 facility campaign.
package plantbuild

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"time"
)

type Money struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
	Scale    int    `json:"scale"`
}
type SiteOption struct {
	Code, Mode, AvailableAt                                                        string
	TotalCash                                                                      Money
	AreaM2, ElectricityKVA, LogisticsScore, TalentScore, RiskScore, ExpansionScore int
}
type Assessment struct {
	SiteCode      string   `json:"site_code"`
	Feasible      bool     `json:"feasible"`
	HardFailures  []string `json:"hard_failures"`
	WeightedScore string   `json:"weighted_score"`
	Source        string   `json:"source"`
	Confidence    string   `json:"confidence"`
	Version       string   `json:"version"`
}
type Zone struct {
	Code    string `json:"code"`
	Parent  string `json:"parent"`
	Purpose string `json:"purpose"`
	Status  string `json:"status"`
	AreaM2  int    `json:"area_m2"`
}
type WorkPackage struct {
	Code         string   `json:"code"`
	Predecessors []string `json:"predecessors"`
	Status       string   `json:"status"`
	Cost         Money    `json:"cost"`
	Evidence     string   `json:"evidence,omitempty"`
}
type Knowledge struct {
	Actor      string `json:"actor"`
	Fact       string `json:"fact"`
	ObservedAt string `json:"observed_at"`
	Source     string `json:"source"`
	Visibility string `json:"visibility"`
}
type Frame struct {
	Step                    int               `json:"step"`
	Phase                   string            `json:"phase"`
	SimTime                 string            `json:"sim_time"`
	Title                   string            `json:"title"`
	CausationID             string            `json:"causation_id"`
	SelectedSite            string            `json:"selected_site"`
	Assessments             []Assessment      `json:"assessments"`
	Zones                   []Zone            `json:"zones"`
	WorkPackages            []WorkPackage     `json:"work_packages"`
	Utilities               map[string]string `json:"utilities"`
	Knowledge               []Knowledge       `json:"knowledge"`
	WorldProgress           int               `json:"world_progress"`
	IAOSPlanProgress        int               `json:"iaos_plan_progress"`
	Discrepancy             string            `json:"discrepancy"`
	Cash                    Money             `json:"cash"`
	Committed               Money             `json:"committed"`
	Payable                 Money             `json:"payable"`
	Paid                    Money             `json:"paid"`
	CapabilityBuildEligible bool              `json:"capability_build_eligible"`
	IAOSCursor              int64             `json:"iaos_cursor"`
}
type Trace struct {
	SchemaVersion  string  `json:"schema_version"`
	Campaign       string  `json:"campaign"`
	WorldRunID     string  `json:"world_run_id"`
	Timezone       string  `json:"timezone"`
	PolicyVersion  string  `json:"policy_version"`
	M9TerminalHash string  `json:"m9_terminal_hash"`
	Frames         []Frame `json:"frames"`
}

func money(v string) Money { return Money{v, "CNY", 2} }
func Candidates() []SiteOption {
	return []SiteOption{{"SITE-SZ-EAST-GREENFIELD", "greenfield", "2027-01-15T00:00:00+08:00", money("28000000.00"), 32000, 4000, 88, 78, 90, 100}, {"SITE-SZ-NORTH-LEASED-SHELL", "leased_shell", "2026-05-01T00:00:00+08:00", money("13500000.00"), 18000, 2400, 86, 82, 78, 70}, {"SITE-SZ-WEST-BUILD-TO-SUIT", "build_to_suit", "2026-08-15T00:00:00+08:00", money("16000000.00"), 22000, 3000, 80, 75, 82, 88}}
}
func Assess(s SiteOption) Assessment {
	fails := []string{}
	cost, _ := new(big.Rat).SetString(s.TotalCash.Value)
	limit, _ := new(big.Rat).SetString("15000000.00")
	if cost.Cmp(limit) > 0 {
		fails = append(fails, "budget_exceeded")
	}
	if s.ElectricityKVA < 2200 {
		fails = append(fails, "electricity_below_minimum")
	}
	deadline, _ := time.Parse(time.RFC3339, "2026-09-01T00:00:00+08:00")
	available, _ := time.Parse(time.RFC3339, s.AvailableAt)
	if available.After(deadline) {
		fails = append(fails, "available_after_deadline")
	}
	score := fmt.Sprintf("%.2f", float64(s.LogisticsScore)*.30+float64(s.TalentScore)*.20+float64(100-s.RiskScore)*.20+float64(s.ExpansionScore)*.30)
	return Assessment{s.Code, len(fails) == 0, fails, score, "fictional-site-survey-v1", "0.90", "site-score-v1"}
}
func Ranked() []Assessment {
	out := []Assessment{}
	for _, s := range Candidates() {
		out = append(out, Assess(s))
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Feasible != out[j].Feasible {
			return out[i].Feasible
		}
		return out[i].WeightedScore > out[j].WeightedScore
	})
	return out
}
func BuildTrace() Trace {
	assess := Ranked()
	zones := []Zone{{"SITE-HCTM-SZ-01", "", "site", "controlled", 18000}, {"BUILDING-HCTM-SZ-MAIN", "SITE-HCTM-SZ-01", "building", "accepted", 16000}, {"ZONE-HCTM-SZ-OFFICE", "BUILDING-HCTM-SZ-MAIN", "office", "accepted", 1800}, {"ZONE-HCTM-SZ-PRODUCTION", "BUILDING-HCTM-SZ-MAIN", "production", "accepted", 8500}, {"ZONE-HCTM-SZ-WAREHOUSE", "BUILDING-HCTM-SZ-MAIN", "warehouse", "accepted", 2800}, {"ZONE-HCTM-SZ-QUALITY", "BUILDING-HCTM-SZ-MAIN", "quality", "accepted", 1300}, {"ZONE-HCTM-SZ-UTILITY", "BUILDING-HCTM-SZ-MAIN", "utility", "accepted", 1600}}
	wbs := []WorkPackage{{"WP-DESIGN", []string{}, "completed", money("1000000.00"), "design-accepted"}, {"WP-PERMIT", []string{"WP-DESIGN"}, "completed", money("500000.00"), "permit-observation"}, {"WP-RENOVATION", []string{"WP-PERMIT"}, "completed", money("6500000.00"), "renovation-accepted"}, {"WP-UTILITY", []string{"WP-PERMIT"}, "completed", money("3000000.00"), "utility-accepted"}, {"WP-FIRE-EHS", []string{"WP-RENOVATION", "WP-UTILITY"}, "completed", money("1500000.00"), "inspection-passed"}, {"WP-NETWORK", []string{"WP-RENOVATION"}, "completed", money("1000000.00"), "network-accepted"}}
	base := Frame{0, "eligible", "2026-01-12T16:00:00+08:00", "消费 M9 机器资格", "m9-terminal:sha256:m9-incorporation", "", []Assessment{}, []Zone{}, []WorkPackage{}, map[string]string{}, []Knowledge{}, 0, 0, "none", money("20000000.00"), money("0.00"), money("0.00"), money("0.00"), false, 200}
	evaluating := base
	evaluating.Step = 1
	evaluating.Phase = "evaluating"
	evaluating.SimTime = "2026-01-15T10:00:00+08:00"
	evaluating.Title = "三个候选完成硬约束与解释性评分"
	evaluating.Assessments = assess
	evaluating.CausationID = "observation-site-assessments-001"
	evaluating.IAOSCursor = 204
	selected := evaluating
	selected.Step = 2
	selected.Phase = "site_selected"
	selected.SimTime = "2026-01-16T16:00:00+08:00"
	selected.Title = "CEO/CFO 批准租赁标准厂房方案"
	selected.SelectedSite = "SITE-SZ-NORTH-LEASED-SHELL"
	selected.CausationID = "outcome-site-selection-approved-001"
	selected.IAOSCursor = 208
	controlled := selected
	controlled.Step = 3
	controlled.Phase = "site_controlled"
	controlled.SimTime = "2026-05-01T09:00:00+08:00"
	controlled.Title = "园区实际交付场地控制权"
	controlled.CausationID = "observation-site-control-delivered-001"
	controlled.IAOSCursor = 212
	approved := controlled
	approved.Step = 4
	approved.Phase = "project_approved"
	approved.SimTime = "2026-05-02T10:00:00+08:00"
	approved.Title = "设施项目、WBS 与 1,350 万承诺获批"
	approved.Committed = money("13500000.00")
	approved.CausationID = "outcome-facility-project-approved-001"
	approved.IAOSCursor = 216
	constructing := approved
	constructing.Step = 5
	constructing.Phase = "constructing"
	constructing.SimTime = "2026-06-15T10:00:00+08:00"
	constructing.Title = "改造推进，公用工程服务方实际延期"
	constructing.WorldProgress = 48
	constructing.IAOSPlanProgress = 65
	constructing.Discrepancy = "utility_delay_open"
	constructing.CausationID = "world-utility-delay-001"
	delayed := constructing
	delayed.Step = 6
	delayed.SimTime = "2026-06-16T09:00:00+08:00"
	delayed.Title = "项目负责人收到延期 observation 并提交重排"
	delayed.Knowledge = []Knowledge{{"ACTOR-HCTM-PD-01", "utility-delay", "2026-06-16T09:00:00+08:00", "utility-observation", "project-team"}}
	delayed.CausationID = "intent-project-rebaseline-001"
	delayed.IAOSCursor = 220
	rebaseline := delayed
	rebaseline.Step = 7
	rebaseline.SimTime = "2026-07-20T10:00:00+08:00"
	rebaseline.Title = "IAOS 批准重排，World 计算并行施工后果"
	rebaseline.WorldProgress = 82
	rebaseline.IAOSPlanProgress = 82
	rebaseline.Discrepancy = "utility_delay_mitigated"
	rebaseline.CausationID = "outcome-project-rebaseline-approved-001"
	rebaseline.IAOSCursor = 224
	acceptance := rebaseline
	acceptance.Step = 8
	acceptance.Phase = "facility_acceptance"
	acceptance.SimTime = "2026-08-20T10:00:00+08:00"
	acceptance.Title = "空间、公用工程、消防与 EHS 全部验收"
	acceptance.WorldProgress = 100
	acceptance.IAOSPlanProgress = 100
	acceptance.Zones = zones
	acceptance.WorkPackages = wbs
	acceptance.Utilities = map[string]string{"electricity": "accepted:2400kVA", "water": "accepted", "compressed_air": "accepted", "fire": "accepted", "ehs": "accepted", "network": "accepted"}
	acceptance.Discrepancy = "closed"
	acceptance.Payable = money("10000000.00")
	acceptance.Paid = money("10000000.00")
	acceptance.Cash = money("10000000.00")
	acceptance.CausationID = "observation-facility-accepted-001"
	acceptance.IAOSCursor = 230
	eligible := acceptance
	eligible.Step = 9
	eligible.Phase = "capability_build_eligible"
	eligible.SimTime = "2026-08-20T16:00:00+08:00"
	eligible.Title = "设施载体就绪，获得 M11 能力建设资格"
	eligible.CapabilityBuildEligible = true
	eligible.CausationID = "world-capability-build-eligibility-001"
	return Trace{"1.0", "plant-build", "world-run-genesis-plant-build-001", "Asia/Shanghai", "plant-build-policy-v1", "sha256:m9-incorporation-terminal", []Frame{base, evaluating, selected, controlled, approved, constructing, delayed, rebaseline, acceptance, eligible}}
}
func Validate(t Trace) error {
	want := []string{"eligible", "evaluating", "site_selected", "site_controlled", "project_approved", "constructing", "constructing", "constructing", "facility_acceptance", "capability_build_eligible"}
	if t.M9TerminalHash == "" || len(t.Frames) != len(want) {
		return errors.New("invalid M9 input or frame count")
	}
	for i, f := range t.Frames {
		if f.Phase != want[i] {
			return fmt.Errorf("phase %d", i)
		}
		if _, e := time.Parse(time.RFC3339, f.SimTime); e != nil {
			return e
		}
	}
	rank := Ranked()
	if len(rank) != 3 || rank[0].SiteCode != "SITE-SZ-NORTH-LEASED-SHELL" || !rank[0].Feasible {
		return errors.New("invalid site decision")
	}
	last := t.Frames[9]
	cash, _ := new(big.Rat).SetString(last.Cash.Value)
	paid, _ := new(big.Rat).SetString(last.Paid.Value)
	opening, _ := new(big.Rat).SetString("20000000.00")
	if new(big.Rat).Add(cash, paid).Cmp(opening) != 0 {
		return errors.New("cash not conserved")
	}
	commit, _ := new(big.Rat).SetString(last.Committed.Value)
	budget, _ := new(big.Rat).SetString("15000000.00")
	if commit.Cmp(budget) > 0 || paid.Cmp(commit) > 0 {
		return errors.New("budget or payment violation")
	}
	if !last.CapabilityBuildEligible || len(last.Zones) != 7 || last.Discrepancy != "closed" {
		return errors.New("acceptance incomplete")
	}
	return nil
}
func Hash(t Trace) string {
	b, _ := json.Marshal(t)
	h := sha256.Sum256(b)
	return "sha256:" + hex.EncodeToString(h[:])
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
	if e := d.Decode(&t); e != nil {
		return t, e
	}
	return t, Validate(t)
}
func Reset() Trace { return BuildTrace() }

type Operator struct {
	Actor, Mode string
	Permissions map[string]bool
}

func Authorize(o Operator, action string, accepted bool, withinBudget bool) error {
	if o.Mode != "human" && o.Mode != "agent" {
		return errors.New("bad mode")
	}
	if !o.Permissions[action] {
		return errors.New("permission denied")
	}
	if action == "genesis.payment.approve" && !accepted {
		return errors.New("milestone not accepted")
	}
	if !withinBudget {
		return errors.New("budget exceeded")
	}
	return nil
}
