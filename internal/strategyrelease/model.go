package strategyrelease

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
)

const EvidenceHash = "af4b3a564635baf0f5d0000220a049967e64db962056f38e986a913b75936ce7"

type Thresholds struct {
	SafetyStock     int   `json:"safety_stock_units"`
	DualSourcePct   int   `json:"dual_source_pct"`
	MaintenanceGate int   `json:"maintenance_gate_pct"`
	OvertimeLimit   int   `json:"overtime_limit_units"`
	CashBuffer      int64 `json:"cash_buffer_cny"`
}
type Release struct {
	Code         string     `json:"code"`
	Version      string     `json:"version"`
	EvidenceHash string     `json:"evidence_hash"`
	PriorRelease string     `json:"prior_release"`
	Scope        []string   `json:"scope"`
	WindowWeeks  int        `json:"window_weeks"`
	Thresholds   Thresholds `json:"thresholds"`
	Hash         string     `json:"hash"`
}
type Guardrail struct {
	Code      string `json:"code"`
	Severity  string `json:"severity"`
	Threshold string `json:"threshold"`
	Status    string `json:"status"`
	Owner     string `json:"owner"`
}
type Commitment struct {
	Code         string `json:"code"`
	Kind         string `json:"kind"`
	Quantity     int    `json:"quantity"`
	Status       string `json:"status"`
	Compensation string `json:"compensation"`
}
type Frame struct {
	Step                      int          `json:"step"`
	Phase                     string       `json:"phase"`
	Title                     string       `json:"title"`
	Actor                     string       `json:"actor"`
	Approver                  string       `json:"approver"`
	ShadowWeek                int          `json:"shadow_week"`
	PilotWeek                 int          `json:"pilot_week"`
	BusinessWrites            int          `json:"business_writes"`
	CandidateActions          int          `json:"candidate_actions"`
	ActiveRelease             string       `json:"active_release"`
	Guardrails                []Guardrail  `json:"guardrails"`
	Commitments               []Commitment `json:"commitments"`
	Disposition               string       `json:"disposition"`
	StrategyChangeCycleClosed bool         `json:"strategy_change_cycle_closed"`
	Note                      string       `json:"note"`
}
type Trace struct {
	SchemaVersion  string            `json:"schema_version"`
	Candidate      string            `json:"candidate"`
	EvidenceHash   string            `json:"evidence_hash"`
	Release        Release           `json:"release"`
	SemanticDiff   map[string]string `json:"semantic_diff"`
	Frames         []Frame           `json:"frames"`
	RollbackFrames []Frame           `json:"rollback_frames"`
	TraceHash      string            `json:"trace_hash"`
	CausalClaim    string            `json:"causal_claim"`
}

func canonical(v any) string {
	b, _ := json.Marshal(v)
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}
func BuildTrace() Trace {
	r := Release{Code: "STR-GENESIS-M15-RESILIENT", Version: "1.0.0", EvidenceHash: EvidenceHash, PriorRelease: "STR-GENESIS-BASELINE@1.0.0", Scope: []string{"HCTM-SZ/A-LINE/HCTM-BCP-A01"}, WindowWeeks: 4, Thresholds: Thresholds{1200, 30, 65, 1800, 4000000}}
	r.Hash = canonical(r)
	g := []Guardrail{{"quality_gate", "hard_stop", "accepted_defect=0", "ok", "Quality"}, {"cash_buffer", "pause", "CNY 4,000,000", "ok", "CFO"}, {"otif", "pause", "98%", "ok", "Operations"}, {"monitor_freshness", "hard_stop", "<=15m", "ok", "Risk"}}
	frames := []Frame{{1, "candidate", "M14 evidence verified", "Strategy Analyst", "", 0, 0, 0, 0, r.PriorRelease, g, nil, "", false, "Pareto candidate; limitations retained"}, {2, "review", "Independent owners approve shadow", "Strategy Analyst", "CFO + Risk + Operations", 0, 0, 0, 0, r.PriorRelease, g, nil, "", false, "proposer cannot approve"}, {3, "shadow_running", "Four-week shadow begins", "Policy Evaluator", "Risk", 1, 0, 0, 3, r.PriorRelease, g, nil, "", false, "candidate decisions only"}, {4, "shadow_complete", "Shadow closes with zero writes", "Policy Evaluator", "Risk", 4, 0, 0, 12, r.PriorRelease, g, nil, "", false, "no intent, capability write, reservation or outcome"}, {5, "pilot_approved", "Independent pilot gate", "Operations", "CFO + Risk", 4, 0, 0, 12, r.PriorRelease, g, nil, "", false, "exact release hash approved"}, {6, "pilot_running", "Bounded canonical pilot", "Operations", "Risk", 4, 1, 3, 12, r.Code + "@" + r.Version, g, []Commitment{{"PO-PILOT-001", "purchase_order", 600, "open", "consume_under_prior_release"}}, "", false, "all actions via IAOS governance"}, {7, "pilot_running", "Pilot monitoring", "Operations", "Risk", 4, 4, 9, 12, r.Code + "@" + r.Version, g, []Commitment{{"PO-PILOT-001", "purchase_order", 600, "consumed", "consumed"}}, "", false, "no unknown breach"}, {8, "adoption_review", "Evidence and commitments reconciled", "CFO", "CEO + Risk", 4, 4, 9, 12, r.Code + "@" + r.Version, g, nil, "", false, "pilot is not a randomized causal estimate"}, {9, "closed", "Strategy change cycle closed", "CEO", "Risk", 4, 4, 9, 12, r.Code + "@" + r.Version, g, nil, "adopted", true, "review date and expiry required"}}
	bad := append([]Guardrail(nil), g...)
	bad[0].Status = "breached"
	rollback := []Frame{{1, "pilot_running", "Injected quality breach", "Quality", "", 0, 2, 4, 0, r.Code + "@" + r.Version, bad, []Commitment{{"WO-PILOT-009", "work_order", 300, "open", "isolate_and_inspect"}}, "", false, "hard stop"}, {2, "paused", "Kill switch stops future intent", "Risk", "CFO", 0, 2, 4, 0, r.Code + "@" + r.Version, bad, nil, "", false, "committed facts retained"}, {3, "rollback", "Prior release restored by CAS", "CFO", "CEO + Risk", 0, 2, 4, 0, r.PriorRelease, bad, []Commitment{{"WO-PILOT-009", "work_order", 300, "compensated", "isolated_and_inspected"}}, "", false, "rollback does not delete work order"}, {4, "closed", "Rollback cycle reconciled", "CEO", "Risk", 0, 2, 4, 0, r.PriorRelease, g, nil, "rolled_back", true, "all commitments resolved"}}
	t := Trace{SchemaVersion: "1.0", Candidate: "resilient", EvidenceHash: EvidenceHash, Release: r, SemanticDiff: map[string]string{"safety_stock": "600 -> 1200 units", "dual_source": "0% -> 30%", "maintenance_gate": "75% -> 65%", "overtime_limit": "1200 -> 1800 units", "cash_buffer": "CNY 3,000,000 -> 4,000,000"}, Frames: frames, RollbackFrames: rollback, CausalClaim: "Bounded pilot observation only; no statistical causal or permanent-optimum claim."}
	t.TraceHash = canonical(t)
	return t
}
func Validate(t Trace) error {
	if t.EvidenceHash != EvidenceHash || t.Release.Hash == "" || len(t.Frames) != 9 || len(t.RollbackFrames) != 4 {
		return errors.New("incomplete strategy trace")
	}
	if t.Frames[3].BusinessWrites != 0 || t.Frames[3].CandidateActions == 0 {
		return errors.New("shadow contract violated")
	}
	end := t.Frames[len(t.Frames)-1]
	rb := t.RollbackFrames[len(t.RollbackFrames)-1]
	if !end.StrategyChangeCycleClosed || end.Disposition != "adopted" || !rb.StrategyChangeCycleClosed || rb.Disposition != "rolled_back" {
		return errors.New("cycle not honestly closed")
	}
	if t.RollbackFrames[2].Commitments[0].Status != "compensated" {
		return errors.New("unresolved commitment")
	}
	return nil
}
