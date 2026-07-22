// Package incorporation implements the deterministic M9 Genesis incorporation campaign.
package incorporation

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"
)

type Money struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
	Scale    int    `json:"scale"`
}
type CashAccount struct {
	Code    string `json:"code"`
	Owner   string `json:"owner"`
	Status  string `json:"status"`
	Balance Money  `json:"balance"`
}
type Appointment struct {
	Position   string `json:"position"`
	Assignee   string `json:"assignee"`
	Resolution string `json:"resolution,omitempty"`
	Status     string `json:"status"`
	AcceptedAt string `json:"accepted_at,omitempty"`
}
type Knowledge struct {
	Actor      string `json:"actor"`
	Fact       string `json:"fact"`
	ObservedAt string `json:"observed_at"`
	Source     string `json:"source"`
	Confidence string `json:"confidence"`
	Visibility string `json:"visibility"`
}
type Governance struct {
	CEO             Appointment `json:"ceo"`
	CFO             Appointment `json:"cfo"`
	ProjectDirector Appointment `json:"project_director"`
	MandateActive   bool        `json:"mandate_active"`
}
type Budget struct {
	Code   string `json:"code"`
	Status string `json:"status"`
	Amount Money  `json:"amount"`
	Owner  string `json:"owner"`
}
type Frame struct {
	Step                 int         `json:"step"`
	Phase                string      `json:"phase"`
	SimTime              string      `json:"sim_time"`
	Title                string      `json:"title"`
	CausationID          string      `json:"causation_id"`
	LegalEntityStatus    string      `json:"legal_entity_status"`
	RegistrationStatus   string      `json:"registration_status"`
	Investor             CashAccount `json:"investor"`
	Company              CashAccount `json:"company"`
	CapitalCommitted     Money       `json:"capital_committed"`
	CapitalPaid          Money       `json:"capital_paid"`
	Governance           Governance  `json:"governance"`
	Budget               Budget      `json:"budget"`
	Knowledge            []Knowledge `json:"knowledge"`
	IAOSCursor           int64       `json:"iaos_cursor"`
	PlantProjectEligible bool        `json:"plant_project_eligible"`
}
type Trace struct {
	SchemaVersion string  `json:"schema_version"`
	Campaign      string  `json:"campaign"`
	WorldRunID    string  `json:"world_run_id"`
	Timezone      string  `json:"timezone"`
	PolicyVersion string  `json:"policy_version"`
	Frames        []Frame `json:"frames"`
}

func money(v string) Money { return Money{v, "CNY", 2} }
func emptyAppointment(position, assignee string) Appointment {
	return Appointment{Position: position, Assignee: assignee, Status: "pending"}
}

func BuildTrace() Trace {
	inv := CashAccount{"INV-HCTM-FOUNDERS-CASH", "INV-HCTM-FOUNDERS", "open", money("50000000.00")}
	company := CashAccount{"BANK-HCTM-SZ-OPERATING-01", "HCTM-SZ-MFG", "not_opened", money("0.00")}
	gov := Governance{emptyAppointment("POS-HCTM-CEO", "ACTOR-HCTM-CEO-01"), emptyAppointment("POS-HCTM-CFO", "ACTOR-HCTM-CFO-01"), emptyAppointment("POS-HCTM-SZ-PROJECT-DIRECTOR", "ACTOR-HCTM-PD-01"), false}
	budget := Budget{"BUDGET-HCTM-SZ-Y1", "not_submitted", money("15000000.00"), "HCTM-SZ-MFG"}
	base := Frame{0, "pre_incorporation", "2026-01-05T09:00:00+08:00", "投资人形成设立意图", "intent-incorporation-001", "not_registered", "not_submitted", inv, company, money("30000000.00"), money("0.00"), gov, budget, []Knowledge{}, 100, false}
	registering := base
	registering.Step = 1
	registering.Phase = "registering"
	registering.SimTime = "2026-01-05T10:00:00+08:00"
	registering.Title = "IAOS 已批准设立方案，外部登记处理中"
	registering.RegistrationStatus = "submitted"
	registering.CausationID = "outcome-incorporation-approved-001"
	registering.IAOSCursor = 102
	registered := registering
	registered.Step = 2
	registered.Phase = "registered"
	registered.SimTime = "2026-01-08T10:00:00+08:00"
	registered.Title = "虚构监管机构批准法人登记"
	registered.RegistrationStatus = "approved"
	registered.LegalEntityStatus = "registered"
	registered.Investor.Balance = money("49990000.00")
	registered.CausationID = "observation-registration-approved-001"
	registered.IAOSCursor = 104
	capitalizing := registered
	capitalizing.Step = 3
	capitalizing.Phase = "capitalizing"
	capitalizing.SimTime = "2026-01-10T10:00:00+08:00"
	capitalizing.Title = "运营账户开户并完成首期资本到账"
	capitalizing.Company.Status = "open"
	capitalizing.Company.Balance = money("20000000.00")
	capitalizing.Investor.Balance = money("29990000.00")
	capitalizing.CapitalPaid = money("20000000.00")
	capitalizing.CausationID = "observation-capital-received-001"
	capitalizing.IAOSCursor = 108
	organizing := capitalizing
	organizing.Step = 4
	organizing.Phase = "organizing"
	organizing.SimTime = "2026-01-10T14:00:00+08:00"
	organizing.Title = "董事会任命管理层，等待岗位接受"
	organizing.Governance.CEO.Resolution = "RES-HCTM-BOARD-CEO-001"
	organizing.Governance.CFO.Resolution = "RES-HCTM-BOARD-CFO-001"
	organizing.Governance.ProjectDirector.Resolution = "RES-HCTM-BOARD-PD-001"
	organizing.CausationID = "outcome-appointments-committed-001"
	organizing.IAOSCursor = 112
	accepted := organizing
	accepted.Step = 5
	accepted.Phase = "organizing"
	accepted.SimTime = "2026-01-11T09:00:00+08:00"
	accepted.Title = "CEO、CFO 与项目负责人接受任命并获知职责"
	accepted.Governance.CEO.Status = "accepted"
	accepted.Governance.CEO.AcceptedAt = accepted.SimTime
	accepted.Governance.CFO.Status = "accepted"
	accepted.Governance.CFO.AcceptedAt = accepted.SimTime
	accepted.Governance.ProjectDirector.Status = "accepted"
	accepted.Governance.ProjectDirector.AcceptedAt = accepted.SimTime
	accepted.Governance.MandateActive = true
	accepted.Knowledge = []Knowledge{{"ACTOR-HCTM-CEO-01", "appointment:POS-HCTM-CEO", "2026-01-11T09:05:00+08:00", "board-resolution", "1.00", "assignee"}, {"ACTOR-HCTM-CFO-01", "cash:capital-received", "2026-01-11T09:10:00+08:00", "bank-observation", "1.00", "finance"}, {"ACTOR-HCTM-PD-01", "appointment:POS-HCTM-SZ-PROJECT-DIRECTOR", "2026-01-11T09:15:00+08:00", "board-resolution", "1.00", "assignee"}}
	accepted.CausationID = "observation-appointments-accepted-001"
	accepted.IAOSCursor = 116
	budgeted := accepted
	budgeted.Step = 6
	budgeted.Phase = "budgeted"
	budgeted.SimTime = "2026-01-12T15:00:00+08:00"
	budgeted.Title = "CEO 提交、CFO 审查并由独立审批人批准启动预算"
	budgeted.Budget.Status = "approved"
	budgeted.CausationID = "outcome-budget-approved-001"
	budgeted.IAOSCursor = 120
	eligible := budgeted
	eligible.Step = 7
	eligible.Phase = "plant_project_eligible"
	eligible.SimTime = "2026-01-12T16:00:00+08:00"
	eligible.Title = "法人、资本、岗位、mandate 与预算齐备，获得 M10 资格"
	eligible.PlantProjectEligible = true
	eligible.CausationID = "world-eligibility-evaluated-001"
	return Trace{"1.0", "incorporation", "world-run-genesis-incorporation-001", "Asia/Shanghai", "genesis-incorporation-policy-v1", []Frame{base, registering, registered, capitalizing, organizing, accepted, budgeted, eligible}}
}

func rat(v string) (*big.Rat, error) {
	r, ok := new(big.Rat).SetString(v)
	if !ok {
		return nil, fmt.Errorf("invalid decimal %q", v)
	}
	return r, nil
}
func Validate(t Trace) error {
	want := []string{"pre_incorporation", "registering", "registered", "capitalizing", "organizing", "organizing", "budgeted", "plant_project_eligible"}
	if t.SchemaVersion != "1.0" || t.Timezone != "Asia/Shanghai" || len(t.Frames) != len(want) {
		return fmt.Errorf("invalid campaign envelope")
	}
	var previous time.Time
	for i, f := range t.Frames {
		at, e := time.Parse(time.RFC3339, f.SimTime)
		if e != nil {
			return e
		}
		if i > 0 && at.Before(previous) {
			return fmt.Errorf("time moved backwards")
		}
		previous = at
		if f.Phase != want[i] {
			return fmt.Errorf("phase %d: %s", i, f.Phase)
		}
		if f.Investor.Owner == "" || f.Company.Owner == "" {
			return fmt.Errorf("cash without owner")
		}
		if f.Budget.Amount.Currency != "CNY" || f.Budget.Amount.Scale != 2 {
			return fmt.Errorf("invalid budget unit")
		}
	}
	last := t.Frames[len(t.Frames)-1]
	opening, _ := rat("50000000.00")
	paid, _ := rat(last.CapitalPaid.Value)
	fee, _ := rat("10000.00")
	closing, _ := rat(last.Investor.Balance.Value)
	right := new(big.Rat).Add(paid, fee)
	right.Add(right, closing)
	if opening.Cmp(right) != 0 {
		return fmt.Errorf("investor cash not conserved")
	}
	company, _ := rat(last.Company.Balance.Value)
	if company.Cmp(paid) != 0 {
		return fmt.Errorf("company cash not conserved")
	}
	if !last.PlantProjectEligible || last.LegalEntityStatus != "registered" || last.Company.Status != "open" || !last.Governance.MandateActive || last.Budget.Status != "approved" {
		return fmt.Errorf("eligibility prerequisites incomplete")
	}
	return nil
}
func JSON() []byte { b, _ := json.Marshal(BuildTrace()); return b }
func Hash(t Trace) string {
	b, _ := json.Marshal(t)
	h := sha256.Sum256(b)
	return "sha256:" + hex.EncodeToString(h[:])
}
func DecideRegistration(complete bool) string {
	if complete {
		return "approved"
	}
	return "rejected"
}
func CanSubmitBudget(appointed, accepted, mandate bool) bool { return appointed && accepted && mandate }

type Operator struct { Actor string; Mode string; Permissions map[string]bool }
func Authorize(operator Operator, action string, frame Frame) error {
	if operator.Mode != "human" && operator.Mode != "agent" { return errors.New("invalid operator mode") }
	if !operator.Permissions[action] { return errors.New("permission denied") }
	if action == "genesis.budget.submit" && (!frame.Governance.MandateActive || frame.Governance.CEO.Status != "accepted") { return errors.New("CEO appointment or mandate inactive") }
	if action == "genesis.budget.approve" && operator.Actor == frame.Governance.CEO.Assignee { return errors.New("self approval forbidden") }
	return nil
}
func Snapshot(t Trace) ([]byte, error) { if err:=Validate(t);err!=nil{return nil,err};return json.Marshal(t) }
func Restore(data []byte) (Trace, error) { var t Trace;dec:=json.NewDecoder(bytes.NewReader(data));dec.DisallowUnknownFields();if err:=dec.Decode(&t);err!=nil{return t,err};return t,Validate(t) }
func Reset() Trace { return BuildTrace() }
