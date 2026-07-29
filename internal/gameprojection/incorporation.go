package gameprojection

import (
	"fmt"
	"strings"

	"github.com/industrial-ai/iaos-aese/internal/incorporation"
)

var chapterByPhase = map[string]string{
	"pre_incorporation":      "founder_intent",
	"registering":            "registration",
	"registered":             "banking",
	"capitalizing":           "banking",
	"organizing":             "talent_governance",
	"budgeted":               "opening",
	"plant_project_eligible": "operating_world",
}

// FromIncorporationTrace projects an already validated business trace into a
// presentation contract. It never advances or mutates the source lifecycle.
func FromIncorporationTrace(trace incorporation.Trace, caseCode string, frameIndex int) (Projection, error) {
	if err := incorporation.Validate(trace); err != nil {
		return Projection{}, fmt.Errorf("source trace: %w", err)
	}
	if caseCode == "" || frameIndex < 0 || frameIndex >= len(trace.Frames) {
		return Projection{}, fmt.Errorf("case code or frame index is invalid")
	}
	frame := trace.Frames[frameIndex]
	chapter, ok := chapterByPhase[frame.Phase]
	if !ok {
		return Projection{}, fmt.Errorf("phase %q has no game chapter", frame.Phase)
	}
	progress := frameIndex * 100 / (len(trace.Frames) - 1)
	projection := Projection{
		SchemaVersion: SchemaVersion,
		ProjectionID:  fmt.Sprintf("gp-%s-%d", strings.ToLower(caseCode), frame.IAOSCursor),
		TenantID:      "tenant-hctm-genesis",
		CaseCode:      caseCode,
		WorldRunID:    trace.WorldRunID,
		Chapter:       chapter,
		SimTime:       frame.SimTime,
		TimeScale:     0,
		Paused:        true,
		Cursor:        frame.IAOSCursor,
		Scene:         Scene{SceneID: sceneForChapter(chapter), Mode: "2.5d", Theme: "industrial-warm"},
		Lifecycle:     Lifecycle{State: frame.Phase, CurrentStep: frame.Title, Progress: progress},
		Buildings:     buildingsForChapter(chapter),
		Actors:        actorsForFrame(frame, frameIndex),
		WorkItems:     workItemsForFrame(frame, caseCode, frameIndex),
		Resources: Resources{
			FounderCash:      fromMoney(frame.Investor.Balance),
			CompanyCash:      fromMoney(frame.Company.Balance),
			CapitalCommitted: fromMoney(frame.CapitalCommitted),
			CapitalPaid:      fromMoney(frame.CapitalPaid),
			BudgetAuthorized: fromMoney(frame.Budget.Amount),
			RiskLevel:        "low",
		},
		Finance: FinanceOpening{
			Roles: []string{}, TrialBalance: []FinanceTrialBalanceLine{},
			BankJournal: []FinanceBankJournalLine{}, GeneralLedger: []FinanceGeneralLedgerLine{},
			OpeningBalanceSheet: FinanceOpeningBalanceSheet{
				Assets: []FinanceStatementLine{}, Liabilities: []FinanceStatementLine{}, Equity: []FinanceStatementLine{},
			},
		},
		Exchanges:     exchangesForFrame(frame, caseCode),
		Brand:         Brand{Status: "candidate", CompanyName: "华辰热管理系统集团有限公司", PrimaryColor: "#2563EB"},
		Notifications: []Notification{{NotificationID: fmt.Sprintf("notice-%d", frame.IAOSCursor), Severity: "info", Message: frame.Title}},
		EvidenceRefs: []EvidenceRef{
			{Ref: "trace:" + caseCode, Kind: "incorporation_trace"},
			{Ref: "process:enterprise.incorporation.lifecycle.v1", Kind: "process"},
		},
	}
	if frame.PlantProjectEligible {
		projection.Brand.Status = "approved"
	}
	if err := projection.Validate(); err != nil {
		return Projection{}, err
	}
	return projection, nil
}

func fromMoney(value incorporation.Money) Money {
	parts := strings.SplitN(value.Value, ".", 2)
	minor := parts[0] + "00"
	if len(parts) == 2 {
		fraction := parts[1] + "00"
		minor = parts[0] + fraction[:2]
	}
	return Money{Value: minor, Currency: value.Currency, Scale: value.Scale}
}

func sceneForChapter(chapter string) string {
	return map[string]string{
		"founder_intent": "founder-office-v1", "registration": "civic-center-v1",
		"banking": "banking-district-v1", "talent_governance": "headquarters-v1",
		"opening": "headquarters-opening-v1", "operating_world": "enterprise-city-v1",
	}[chapter]
}

func buildingsForChapter(chapter string) []Building {
	all := []Building{
		{Code: "BLD-FOUNDER-OFFICE", Kind: "office", Label: "创始办公室", State: "active", X: 2, Y: 4, Available: true},
		{Code: "BLD-CIVIC-CENTER", Kind: "government", Label: "政务服务中心", State: "locked", X: 6, Y: 2},
		{Code: "BLD-BANK", Kind: "bank", Label: "合作银行", State: "locked", X: 8, Y: 5},
		{Code: "BLD-HEADQUARTERS", Kind: "headquarters", Label: "企业总部", State: "locked", X: 5, Y: 8},
	}
	rank := map[string]int{"founder_intent": 0, "registration": 1, "banking": 2, "talent_governance": 3, "opening": 3, "operating_world": 3}[chapter]
	for i := range all {
		if i <= rank {
			all[i].Available, all[i].State = true, "active"
		}
	}
	return all
}

func actorsForFrame(_ incorporation.Frame, frameIndex int) []Actor {
	actors := []Actor{
		{ActorID: "founder-principal", ActorType: "human", DisplayName: "创始治理者", Position: "董事长", State: "available"},
		{ActorID: "incorporation-agent", ActorType: "agent", DisplayName: "企业设立专员", Position: "设立与登记", State: "available"},
		{ActorID: "legal-compliance-agent", ActorType: "agent", DisplayName: "法务合规专员", Position: "法务与治理", State: "available"},
		{ActorID: "finance-agent", ActorType: "agent", DisplayName: "财务负责人", Position: "资本与预算", State: "available"},
		{ActorID: "governance-agent", ActorType: "agent", DisplayName: "治理组织专员", Position: "组织与任命", State: "available"},
		{ActorID: "audit-agent", ActorType: "agent", DisplayName: "独立审计专员", Position: "独立检查", State: "available"},
	}
	active := []int{1, 1, 2, 2, 4, 4, 3, 5}[frameIndex]
	actors[active].State = "working"
	actors[active].WorkItemID = fmt.Sprintf("WI-%02d", workItemRangeByFrame[frameIndex][0])
	return actors
}

type workItemSpec struct{ capability, kind, owner, gate string }

var workItemRangeByFrame = [8][2]int{{1, 3}, {4, 7}, {8, 9}, {10, 10}, {11, 14}, {15, 15}, {16, 17}, {18, 18}}

var incorporationWorkItems = []workItemSpec{
	{"incorporation.case.open", "human_task", "founder-principal", ""},
	{"founder.resolution.prepare", "agent_task", "incorporation-agent", ""},
	{"founder.resolution.approve", "approval", "founder-principal", "G1"},
	{"capital.commitment.record", "agent_task", "finance-agent", ""},
	{"registration.package.validate", "agent_task", "legal-compliance-agent", ""},
	{"registration.submit", "approval", "founder-principal", "G2"},
	{"registration.observation.commit", "world_wait", "world-registry", ""},
	{"bank.account.opening.submit", "approval", "founder-principal", "G3"},
	{"bank.account.observation.commit", "world_wait", "world-bank", ""},
	{"capital.contribution.verify", "agent_task", "finance-agent", "G4"},
	{"organization.establish", "capability", "iaos-runtime", ""},
	{"executive.appointment.propose", "agent_task", "governance-agent", ""},
	{"executive.appointment.acceptance.commit", "world_wait", "world-talent", ""},
	{"executive.appointment.approve", "approval", "founder-principal", "G5"},
	{"operating.mandate.grant", "approval", "founder-principal", "G6"},
	{"initial.budget.prepare", "agent_task", "finance-agent", ""},
	{"initial.budget.approve", "approval", "founder-principal", "G7"},
	{"enterprise.readiness.evaluate", "agent_task", "audit-agent", ""},
}

func workItemsForFrame(frame incorporation.Frame, caseCode string, frameIndex int) []WorkItem {
	activeStart, activeEnd := workItemRangeByFrame[frameIndex][0], workItemRangeByFrame[frameIndex][1]
	out := make([]WorkItem, 0, len(incorporationWorkItems))
	for index, spec := range incorporationWorkItems {
		sequence := index + 1
		status := "locked"
		if sequence < activeStart || frameIndex == len(workItemRangeByFrame)-1 {
			status = "completed"
		} else if sequence <= activeEnd {
			status = "active"
		}
		id := fmt.Sprintf("WI-%02d", sequence)
		out = append(out, WorkItem{
			WorkItemID: id, Title: spec.capability, Kind: spec.kind, Status: status,
			OwnerType: ownerType(spec.owner), OwnerID: spec.owner, Capability: spec.capability, Gate: spec.gate,
			RequiresMe:  status == "active" && spec.owner == "founder-principal",
			EvidenceRef: "trace:" + caseCode + "#" + id,
		})
	}
	_ = frame
	return out
}

func ownerType(owner string) string {
	if owner == "founder-principal" {
		return "human"
	}
	if strings.HasSuffix(owner, "-agent") {
		return "agent"
	}
	if strings.HasPrefix(owner, "world-") {
		return "world"
	}
	return "runtime"
}

func exchangesForFrame(frame incorporation.Frame, caseCode string) []Exchange {
	if frame.CausationID == "" {
		return []Exchange{}
	}
	return []Exchange{{ExchangeID: frame.CausationID, Kind: "committed_outcome", Status: "committed", Correlation: "corr-" + caseCode, EvidenceRef: "trace:" + caseCode + "#" + frame.CausationID, OccurredAt: frame.SimTime}}
}
