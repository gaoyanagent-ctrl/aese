package experiment

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
)

const (
	RulesVersion = "experiment-rules-1.0.0"
	PRNGVersion  = "sha256-counter-v1"
)

var Streams = []string{"demand", "supplier", "equipment", "quality", "payment"}

type Definition struct {
	SchemaVersion string   `json:"schema_version"`
	Code          string   `json:"code"`
	Tenant        string   `json:"tenant"`
	Checkpoint    string   `json:"checkpoint"`
	ParentHash    string   `json:"parent_hash"`
	HorizonWeeks  int      `json:"horizon_weeks"`
	Profiles      []string `json:"profiles"`
	Policies      []Policy `json:"policies"`
	Seeds         []uint64 `json:"seeds"`
	RulesVersion  string   `json:"rules_version"`
	PRNGVersion   string   `json:"prng_version"`
}

type Policy struct {
	Code            string `json:"code"`
	SafetyStock     int    `json:"safety_stock"`
	DualSourcePct   int    `json:"dual_source_pct"`
	MaintenanceGate int    `json:"maintenance_gate"`
	OvertimeLimit   int    `json:"overtime_limit"`
	CashBuffer      int64  `json:"cash_buffer_cny"`
}

type Draw struct {
	Week int    `json:"week"`
	Name string `json:"stream"`
	Raw  uint64 `json:"raw"`
}

type Metrics struct {
	Demand        int   `json:"demand_units"`
	Accepted      int   `json:"accepted_units"`
	BacklogPeak   int   `json:"backlog_peak_units"`
	OTIFBP        int   `json:"otif_basis_points"`
	InventoryPeak int   `json:"inventory_peak_units"`
	CashTrough    int64 `json:"cash_trough_cny"`
	GrossMargin   int64 `json:"gross_margin_cny"`
	Overtime      int   `json:"overtime_units"`
	Scrap         int   `json:"scrap_units"`
	Expedite      int   `json:"expedite_units"`
	RecoveryWeeks int   `json:"recovery_weeks"`
}

type Run struct {
	RunID            string   `json:"run_id"`
	BranchID         string   `json:"branch_id"`
	ParentHash       string   `json:"parent_hash"`
	Profile          string   `json:"profile"`
	Policy           string   `json:"policy"`
	Seed             uint64   `json:"seed"`
	PairKey          string   `json:"pair_key"`
	DrawHash         string   `json:"draw_hash"`
	EventHash        string   `json:"event_log_hash"`
	StateHash        string   `json:"state_hash"`
	MetricsHash      string   `json:"metrics_hash"`
	Metrics          Metrics  `json:"metrics"`
	Violations       []string `json:"violations"`
	Status           string   `json:"status"`
	ProductionWrites int      `json:"production_writes"`
}

type Comparison struct {
	Policy               string `json:"policy"`
	Pairs                int    `json:"pairs"`
	AcceptedDelta        int    `json:"accepted_delta_units"`
	GrossMarginDelta     int64  `json:"gross_margin_delta_cny"`
	CashTroughDelta      int64  `json:"cash_trough_delta_cny"`
	ConstraintViolations int    `json:"constraint_violations"`
	Pareto               bool   `json:"pareto"`
}

type EvidenceBundle struct {
	SchemaVersion         string       `json:"schema_version"`
	ExperimentCode        string       `json:"experiment_code"`
	DefinitionHash        string       `json:"definition_hash"`
	Checkpoint            string       `json:"checkpoint"`
	ParentHash            string       `json:"parent_hash"`
	RulesVersion          string       `json:"rules_version"`
	PRNGVersion           string       `json:"prng_version"`
	Streams               []string     `json:"streams"`
	Runs                  []Run        `json:"runs"`
	Comparisons           []Comparison `json:"comparisons"`
	FailedRuns            []string     `json:"failed_runs"`
	CancelledRuns         []string     `json:"cancelled_runs"`
	RecommendationStatus  string       `json:"recommendation_status"`
	StrategyEvidenceReady bool         `json:"strategy_evidence_ready"`
	ConclusionLimit       string       `json:"conclusion_limit"`
	EvidenceHash          string       `json:"evidence_hash"`
}

func DefaultDefinition() Definition {
	return Definition{SchemaVersion: "1.0", Code: "EXP-GENESIS-M14-001", Tenant: "tenant-hctm", Checkpoint: "first_commercial_cycle_closed", ParentHash: "m13:9af75fd1cdaef617", HorizonWeeks: 12, Profiles: []string{"nominal", "demand-spike", "supplier-risk", "equipment-risk", "cash-delay"}, Policies: []Policy{{"baseline", 600, 0, 75, 1200, 3000000}, {"lean", 150, 0, 85, 800, 2500000}, {"resilient", 1200, 30, 65, 1800, 4000000}}, Seeds: []uint64{1103, 2207, 3301, 4409}, RulesVersion: RulesVersion, PRNGVersion: PRNGVersion}
}

func Validate(d Definition) error {
	if d.SchemaVersion != "1.0" || d.Code == "" || d.Tenant != "tenant-hctm" || d.Checkpoint != "first_commercial_cycle_closed" || d.ParentHash == "" {
		return errors.New("incompatible experiment identity or checkpoint")
	}
	if d.HorizonWeeks != 12 || d.RulesVersion != RulesVersion || d.PRNGVersion != PRNGVersion {
		return errors.New("unsupported horizon or rule version")
	}
	if len(d.Profiles) == 0 || len(d.Policies) < 3 || len(d.Seeds) == 0 || len(d.Profiles)*len(d.Policies)*len(d.Seeds) > 120 {
		return errors.New("matrix outside approved quota")
	}
	seen := map[string]bool{}
	for _, p := range d.Policies {
		if p.Code == "" || seen[p.Code] || p.SafetyStock < 0 || p.SafetyStock > 2400 || p.DualSourcePct < 0 || p.DualSourcePct > 50 || p.OvertimeLimit < 0 || p.OvertimeLimit > 2400 || p.CashBuffer < 0 {
			return errors.New("invalid policy bounds")
		}
		seen[p.Code] = true
	}
	return nil
}

func hash(v any) string {
	b, _ := json.Marshal(v)
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}

func draw(seed uint64, profile, stream string, week int) uint64 {
	h := sha256.Sum256([]byte(fmt.Sprintf("%s|%d|%s|%s|%d", PRNGVersion, seed, profile, stream, week)))
	return binary.BigEndian.Uint64(h[:8])
}

func draws(seed uint64, profile string, weeks int) []Draw {
	out := make([]Draw, 0, weeks*len(Streams))
	for w := 1; w <= weeks; w++ {
		for _, s := range Streams {
			out = append(out, Draw{w, s, draw(seed, profile, s, w)})
		}
	}
	return out
}

func execute(d Definition, profile string, p Policy, seed uint64) Run {
	ds := draws(seed, profile, d.HorizonWeeks)
	demand, accepted, backlog, peak, invPeak, overtime, scrap, expedite := 0, 0, 0, 0, 0, 0, 0, 0
	cash := int64(14572000)
	trough := cash
	margin := int64(0)
	for w := 1; w <= d.HorizonWeeks; w++ {
		get := func(s string) uint64 { return draw(seed, profile, s, w) }
		q := 900 + int(get("demand")%301)
		if profile == "demand-spike" && w >= 4 && w <= 7 {
			q += 350
		}
		demand += q
		supplyLoss := int(get("supplier") % 121)
		if profile == "supplier-risk" {
			supplyLoss += 120
		}
		machineLoss := int(get("equipment") % 91)
		if profile == "equipment-risk" {
			machineLoss += 130
		}
		qualityLoss := int(get("quality") % 61)
		capacity := 1050 - machineLoss - supplyLoss*(100-p.DualSourcePct)/100
		need := q + backlog + p.SafetyStock
		if need > capacity {
			use := need - capacity
			if use > p.OvertimeLimit/12 {
				use = p.OvertimeLimit / 12
			}
			capacity += use
			overtime += use
		}
		good := capacity - qualityLoss
		if good < 0 {
			good = 0
		}
		scrap += qualityLoss
		shipped := good
		if shipped > q+backlog {
			inv := shipped - (q + backlog)
			if inv > invPeak {
				invPeak = inv
			}
			shipped = q + backlog
		}
		accepted += shipped
		backlog = q + backlog - shipped
		if backlog > peak {
			peak = backlog
		}
		if backlog > p.SafetyStock {
			expedite += backlog - p.SafetyStock
		}
		revenue := int64(shipped) * 1200
		cost := int64(good)*760 + int64(overtime)*0 + int64(qualityLoss)*180
		margin += revenue - cost
		collection := revenue
		if profile == "cash-delay" && get("payment")%100 < 55 {
			collection = 0
		}
		cash += collection - cost
		if cash < trough {
			trough = cash
		}
	}
	otif := 10000
	if demand > 0 {
		otif = accepted * 10000 / demand
	}
	violations := []string{}
	if trough < p.CashBuffer {
		violations = append(violations, "cash_buffer")
	}
	if backlog > 0 {
		violations = append(violations, "unrecovered_backlog")
	}
	m := Metrics{demand, accepted, peak, otif, invPeak, trough, margin, overtime, scrap, expedite, 0}
	if backlog > 0 {
		m.RecoveryWeeks = 1 + backlog/1000
	}
	pair := fmt.Sprintf("%s:%d", profile, seed)
	base := map[string]any{"pair": pair, "policy": p.Code, "parent": d.ParentHash, "rules": d.RulesVersion}
	return Run{RunID: "run-" + hash(base)[:16], BranchID: "branch-" + hash(map[string]any{"base": base, "policy": p})[:16], ParentHash: d.ParentHash, Profile: profile, Policy: p.Code, Seed: seed, PairKey: pair, DrawHash: hash(ds), EventHash: hash(map[string]any{"draws": ds, "policy": p}), StateHash: hash(map[string]any{"metrics": m, "backlog": backlog}), MetricsHash: hash(m), Metrics: m, Violations: violations, Status: "completed", ProductionWrites: 0}
}

func BuildEvidence(d Definition) (EvidenceBundle, error) {
	if err := Validate(d); err != nil {
		return EvidenceBundle{}, err
	}
	runs := make([]Run, 0, len(d.Profiles)*len(d.Policies)*len(d.Seeds))
	for _, profile := range d.Profiles {
		for _, seed := range d.Seeds {
			expected := ""
			for _, p := range d.Policies {
				r := execute(d, profile, p, seed)
				if expected == "" {
					expected = r.DrawHash
				} else if expected != r.DrawHash {
					return EvidenceBundle{}, errors.New("common random numbers violated")
				}
				runs = append(runs, r)
			}
		}
	}
	sort.Slice(runs, func(i, j int) bool { return runs[i].RunID < runs[j].RunID })
	base := map[string]Run{}
	for _, r := range runs {
		if r.Policy == "baseline" {
			base[r.PairKey] = r
		}
	}
	comps := []Comparison{}
	for _, p := range d.Policies {
		if p.Code == "baseline" {
			continue
		}
		c := Comparison{Policy: p.Code}
		for _, r := range runs {
			if r.Policy != p.Code {
				continue
			}
			b := base[r.PairKey]
			c.Pairs++
			c.AcceptedDelta += r.Metrics.Accepted - b.Metrics.Accepted
			c.GrossMarginDelta += r.Metrics.GrossMargin - b.Metrics.GrossMargin
			c.CashTroughDelta += r.Metrics.CashTrough - b.Metrics.CashTrough
			c.ConstraintViolations += len(r.Violations)
		}
		c.Pareto = c.ConstraintViolations == 0 || c.AcceptedDelta > 0 || c.CashTroughDelta > 0
		comps = append(comps, c)
	}
	e := EvidenceBundle{SchemaVersion: "1.0", ExperimentCode: d.Code, DefinitionHash: hash(d), Checkpoint: d.Checkpoint, ParentHash: d.ParentHash, RulesVersion: d.RulesVersion, PRNGVersion: d.PRNGVersion, Streams: append([]string(nil), Streams...), Runs: runs, Comparisons: comps, FailedRuns: []string{}, CancelledRuns: []string{}, RecommendationStatus: "proposed_not_applied", StrategyEvidenceReady: true, ConclusionLimit: "Simulation decision evidence only; no production policy, order, budget, schedule or cash was modified."}
	e.EvidenceHash = hash(e)
	return e, nil
}

func ReplayHash(d Definition) (string, error) {
	e, err := BuildEvidence(d)
	if err != nil {
		return "", err
	}
	return e.EvidenceHash, nil
}
