package assurance

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
)

const Release = "STR-GENESIS-M15-RESILIENT@1.0.0"

type Observation struct {
	Week              int    `json:"week"`
	Demand            int    `json:"demand_units"`
	SupplierDelay     int    `json:"supplier_delay_days"`
	EquipmentDowntime int    `json:"equipment_downtime_hours"`
	YieldBP           int    `json:"yield_basis_points"`
	PaymentDelay      int    `json:"payment_delay_days"`
	PolicyActions     int    `json:"policy_actions"`
	SourceRef         string `json:"source_ref"`
}
type Finding struct {
	Code          string `json:"code"`
	Severity      string `json:"severity"`
	Status        string `json:"status"`
	AffectedRange string `json:"affected_range"`
}
type Drift struct {
	Domain     string `json:"domain"`
	Method     string `json:"method"`
	Baseline   string `json:"baseline"`
	Current    string `json:"current"`
	Support    string `json:"support"`
	SampleSize int    `json:"sample_size"`
}
type Calibration struct {
	Code          string            `json:"code"`
	Parent        string            `json:"parent"`
	FitWeeks      string            `json:"fit_weeks"`
	HoldoutLocked bool              `json:"holdout_locked"`
	Diff          map[string]string `json:"diff"`
	Hash          string            `json:"hash"`
}
type Validation struct {
	Code                 string `json:"code"`
	HoldoutWeeks         string `json:"holdout_weeks"`
	CandidateFrozen      bool   `json:"candidate_frozen"`
	OriginalEvidenceHash string `json:"original_evidence_hash"`
	ReplayRuns           int    `json:"replay_runs"`
	FailedRuns           int    `json:"failed_runs"`
	Result               string `json:"result"`
	Hash                 string `json:"hash"`
}
type Decision struct {
	Disposition   string   `json:"disposition"`
	Reason        string   `json:"reason"`
	Approvers     []string `json:"approvers"`
	ReleaseEffect string   `json:"release_effect"`
	NextReview    string   `json:"next_review"`
	Closed        bool     `json:"strategy_assurance_cycle_closed"`
}
type Cycle struct {
	SchemaVersion     string        `json:"schema_version"`
	Code              string        `json:"code"`
	Tenant            string        `json:"tenant"`
	Release           string        `json:"release"`
	WindowStart       string        `json:"window_start"`
	WindowEnd         string        `json:"window_end"`
	CutoffAt          string        `json:"cutoff_at"`
	Timezone          string        `json:"timezone"`
	CursorRange       string        `json:"cursor_range"`
	Observations      []Observation `json:"observations"`
	DatasetHash       string        `json:"dataset_hash"`
	Missingness       int           `json:"missingness"`
	Corrections       []string      `json:"corrections"`
	Findings          []Finding     `json:"findings"`
	Drift             []Drift       `json:"drift"`
	Calibration       Calibration   `json:"calibration"`
	Validation        Validation    `json:"validation"`
	Decision          Decision      `json:"decision"`
	InjectedDecisions []Decision    `json:"injected_decisions"`
	CycleHash         string        `json:"cycle_hash"`
	Limitations       string        `json:"limitations"`
}

func hash(v any) string {
	b, _ := json.Marshal(v)
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}
func BuildCycle() Cycle {
	o := make([]Observation, 12)
	for i := range o {
		o[i] = Observation{i + 1, 980 + (i*37)%180, (i * 2) % 5, (i * 3) % 7, 9850 - (i%3)*12, (i * 4) % 9, 2 + i%3, fmt.Sprintf("world:week-%02d|iaos:cursor-%d", i+1, 100+i)}
	}
	d := []Drift{{"demand", "bounded-ratio-v1", "900..1250 units", "980..1165 units", "supported", 12}, {"supplier", "bounded-delay-v1", "0..7 days", "0..4 days", "supported", 12}, {"equipment", "bounded-hours-v1", "0..10 h", "0..6 h", "supported", 12}, {"quality", "bounded-yield-v1", "9800..10000 bp", "9826..9850 bp", "supported", 12}, {"payment", "bounded-delay-v1", "0..14 days", "0..8 days", "supported", 12}, {"policy-action", "bounded-frequency-v1", "0..5/week", "2..4/week", "supported", 12}}
	cal := Calibration{"CAL-GENESIS-M16-001", "M14-assumptions@1.0.0", "weeks-01..08", true, map[string]string{"demand_mean": "1075 -> 1068 units", "supplier_delay": "3 -> 2 days"}, ""}
	cal.Hash = hash(cal)
	v := Validation{"VAL-GENESIS-M16-001", "weeks-09..12", true, "af4b3a564635baf0f5d0000220a049967e64db962056f38e986a913b75936ce7", 60, 0, "candidate improves holdout without changing resilient Pareto support", ""}
	v.Hash = hash(v)
	c := Cycle{"1.0", "ASSURE-GENESIS-M16-001", "tenant-hctm", Release, "2026-09-01T00:00:00+08:00", "2026-11-23T23:59:59+08:00", "2026-11-30T00:00:00+08:00", "Asia/Shanghai", "100..111", o, "", 0, []string{}, []Finding{}, d, cal, v, Decision{"renewed", "data complete; assumptions and resilient support remain valid", []string{"Operations", "Risk", "CFO"}, "review expiry extended; thresholds unchanged", "2027-02-22", true}, []Decision{{"reexperiment_required", "injected demand leaves M14 support", []string{"Risk", "Operations"}, "release paused pending new M14 request", "", true}, {"retired", "injected persistent hard quality risk", []string{"Quality", "Risk", "CEO"}, "M15 rollback and commitment handling required", "", true}}, "", "Short synthetic window; assurance, not causal proof or real probability calibration."}
	c.DatasetHash = hash(map[string]any{"cutoff": c.CutoffAt, "cursor": c.CursorRange, "observations": o})
	c.CycleHash = hash(c)
	return c
}
func Validate(c Cycle) error {
	if len(c.Observations) != 12 || c.Missingness != 0 || len(c.Findings) != 0 {
		return errors.New("data quality gate failed")
	}
	if c.Calibration.FitWeeks != "weeks-01..08" || !c.Calibration.HoldoutLocked || c.Validation.HoldoutWeeks != "weeks-09..12" || !c.Validation.CandidateFrozen {
		return errors.New("holdout leakage")
	}
	if c.Validation.OriginalEvidenceHash == "" || c.Validation.ReplayRuns != 60 || c.Validation.FailedRuns != 0 {
		return errors.New("incomplete validation")
	}
	if !c.Decision.Closed || len(c.InjectedDecisions) != 2 {
		return errors.New("decision incomplete")
	}
	return nil
}
