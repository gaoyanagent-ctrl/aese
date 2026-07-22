package aese3

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

type Evidence struct {
	Code   string `json:"code"`
	Value  int    `json:"value"`
	Unit   string `json:"unit"`
	Source string `json:"source_ref"`
}

type Milestone struct {
	Code           string     `json:"code"`
	Title          string     `json:"title"`
	Design         string     `json:"design"`
	Terminal       string     `json:"terminal"`
	TerminalReady  bool       `json:"terminal_ready"`
	WorldOwner     string     `json:"world_owner"`
	BusinessOwner  string     `json:"business_owner"`
	Evidence       []Evidence `json:"evidence"`
	BusinessWrites int        `json:"automatic_business_writes"`
	EvidenceHash   string     `json:"evidence_hash"`
}

type Program struct {
	SchemaVersion                   string      `json:"schema_version"`
	Code                            string      `json:"code"`
	Tenant                          string      `json:"tenant"`
	Timezone                        string      `json:"timezone"`
	ParentTerminal                  string      `json:"parent_terminal"`
	Milestones                      []Milestone `json:"milestones"`
	IndustrySimulationPlatformReady bool        `json:"industry_simulation_platform_ready"`
	ProgramHash                     string      `json:"program_hash"`
	Limitations                     string      `json:"limitations"`
}

func hash(v any) string {
	b, _ := json.Marshal(v)
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}

func BuildProgram() Program {
	rows := []struct {
		code, title, design, terminal string
		evidence                      []Evidence
	}{
		{"M17", "Rolling IBP & S&OP", "DES-019", "integrated_plan_cycle_closed", []Evidence{{"weekly_horizon", 13, "week", "M16:renewed"}, {"monthly_horizon", 12, "month", "M16:renewed"}, {"review_gates", 5, "gate", "IBP-GENESIS-M17-001"}, {"scenarios", 3, "scenario", "IBP-GENESIS-M17-001"}}},
		{"M18", "Product & Customer Portfolio", "DES-020", "portfolio_operating_model_validated", []Evidence{{"products", 2, "product", "PORT-GENESIS-M18-001"}, {"customers", 2, "customer", "PORT-GENESIS-M18-001"}, {"shared_capacity_violations", 0, "violation", "HCTM-A-LINE"}}},
		{"M19", "Multi-site Fulfillment Network", "DES-021", "network_operating_model_validated", []Evidence{{"network_nodes", 3, "node", "NET-GENESIS-M19-001"}, {"lanes", 2, "lane", "NET-GENESIS-M19-001"}, {"unreconciled_in_transit", 0, "shipment", "NET-GENESIS-M19-001"}}},
		{"M20", "After-sales & Closed-loop Quality", "DES-022", "customer_lifecycle_closed", []Evidence{{"complaints", 1, "case", "QUAL-GENESIS-M20-001"}, {"rma_units", 120, "unit", "RMA-GENESIS-M20-001"}, {"unreconciled_units", 0, "unit", "RMA-GENESIS-M20-001"}}},
		{"M21", "Plant Resource & EHS Resilience", "DES-023", "plant_resilience_cycle_closed", []Evidence{{"near_misses", 1, "event", "EHS-GENESIS-M21-001"}, {"utility_outages", 1, "event", "UTIL-GENESIS-M21-001"}, {"safety_bypass", 0, "event", "EHS-GENESIS-M21-001"}}},
		{"M22", "Group Finance & Investment", "DES-024", "group_value_cycle_closed", []Evidence{{"management_ledgers", 3, "view", "FIN-GENESIS-M22-001"}, {"cash_profit_conflations", 0, "violation", "FIN-GENESIS-M22-001"}, {"capex_decisions", 1, "decision", "CAPEX-GENESIS-M22-001"}}},
		{"M23", "Governed Multi-agent Organization", "DES-025", "agent_operating_model_qualified", []Evidence{{"governed_agents", 7, "agent", "AGENT-GENESIS-M23-001"}, {"benchmarks", 3, "benchmark", "AGENT-GENESIS-M23-001"}, {"unauthorized_writes", 0, "write", "AGENT-GENESIS-M23-001"}}},
		{"M24", "Scenario Platform Productization", "DES-026", "industry_simulation_platform_ready", []Evidence{{"certification_gates", 5, "gate", "PLATFORM-GENESIS-M24-001"}, {"reference_packs", 1, "pack", "hctm-genesis@1.0.0"}, {"failed_certification_gates", 0, "gate", "PLATFORM-GENESIS-M24-001"}}},
	}
	m := make([]Milestone, 0, len(rows))
	for _, row := range rows {
		x := Milestone{row.code, row.title, row.design, row.terminal, true, "AESE World", "IAOS governed runtime", row.evidence, 0, ""}
		x.EvidenceHash = hash(x)
		m = append(m, x)
	}
	p := Program{"1.0", "AESE3-GENESIS-M17-M24", "tenant-hctm", "Asia/Shanghai", "strategy_assurance_cycle_closed=true; disposition=renewed", m, true, "", "Synthetic reference evidence; no statutory accounting, real production target, autonomous business execution, or causal proof."}
	p.ProgramHash = hash(p)
	return p
}

func Validate(p Program) error {
	want := []string{"integrated_plan_cycle_closed", "portfolio_operating_model_validated", "network_operating_model_validated", "customer_lifecycle_closed", "plant_resilience_cycle_closed", "group_value_cycle_closed", "agent_operating_model_qualified", "industry_simulation_platform_ready"}
	if len(p.Milestones) != len(want) || !p.IndustrySimulationPlatformReady {
		return errors.New("AESE 3 completion gate is open")
	}
	for i, m := range p.Milestones {
		if m.Code != fmt.Sprintf("M%d", i+17) || m.Terminal != want[i] || !m.TerminalReady {
			return fmt.Errorf("milestone %d terminal mismatch", i+17)
		}
		if m.BusinessWrites != 0 || m.WorldOwner != "AESE World" || m.BusinessOwner != "IAOS governed runtime" {
			return fmt.Errorf("milestone %s ownership boundary violated", m.Code)
		}
		if len(m.Evidence) < 3 || m.EvidenceHash == "" {
			return fmt.Errorf("milestone %s evidence incomplete", m.Code)
		}
		claimed := m.EvidenceHash
		m.EvidenceHash = ""
		if hash(m) != claimed {
			return fmt.Errorf("milestone %s evidence hash mismatch", m.Code)
		}
	}
	claimed := p.ProgramHash
	p.ProgramHash = ""
	if hash(p) != claimed {
		return errors.New("program hash mismatch")
	}
	return nil
}

func ParseStrict(data []byte) (Program, error) {
	var p Program
	d := json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	if err := d.Decode(&p); err != nil {
		return Program{}, err
	}
	if err := d.Decode(&struct{}{}); err != io.EOF {
		return Program{}, errors.New("trailing JSON value")
	}
	if err := Validate(p); err != nil {
		return Program{}, err
	}
	return p, nil
}
