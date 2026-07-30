package gameprojection

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/industrial-ai/iaos-aese/internal/iaosclient"
)

var stateChapter = map[string]string{
	"draft": "founder_intent", "incorporation_opened": "founder_office", "incorporation_case_opened": "founder_office",
	"founder_resolution_approved": "founder_office", "capital_commitments_confirmed": "registration",
	"registration_submitted": "registration", "legal_entity_registered": "banking",
	"bank_account_opening_submitted": "banking", "bank_account_opened": "banking",
	"capital_contribution_verified": "talent_governance", "organization_established": "talent_governance",
	"executive_appointment_proposed": "talent_governance", "executive_appointment_accepted": "talent_governance",
	"executive_appointment_approved": "talent_governance", "executive_appointments_accepted": "talent_governance",
	"operating_mandate_active": "opening", "operating_mandates_activated": "opening",
	"initial_budget_approved": "opening", "enterprise_operational_ready": "operating_world",
}

// FromIAOS builds the game view exclusively from IAOS committed trace and
// persistent work items. It never executes a capability or advances a process.
func FromIAOS(trace iaosclient.IncorporationTrace, items []iaosclient.IncorporationWorkItem) (Projection, error) {
	if trace.SchemaVersion != "1.0" || !trace.Verified || trace.CaseCode == "" || trace.State.TenantID == "" {
		return Projection{}, fmt.Errorf("IAOS incorporation trace is incomplete or unverified")
	}
	chapter, ok := stateChapter[trace.State.State]
	if !ok {
		return Projection{}, fmt.Errorf("IAOS state %q has no game chapter", trace.State.State)
	}
	simTime := ""
	correlation := ""
	for _, entry := range trace.Journal {
		if _, err := time.Parse(time.RFC3339, entry.CreatedAt); err == nil {
			simTime = entry.CreatedAt
		}
		if entry.CorrelationID != "" {
			correlation = entry.CorrelationID
		}
	}
	if simTime == "" {
		return Projection{}, fmt.Errorf("IAOS trace has no RFC3339 journal timestamp")
	}
	cursor := int64(len(trace.Journal))
	exchanges := make([]Exchange, 0, len(trace.WorldExchanges))
	for _, item := range trace.WorldExchanges {
		if item.Cursor > cursor {
			cursor = item.Cursor
		}
		exchanges = append(exchanges, Exchange{
			ExchangeID: item.MessageID, Kind: item.Kind, Status: "committed",
			Correlation: item.CorrelationID, EvidenceRef: "world:" + item.MessageID, OccurredAt: item.RecordedAt,
		})
	}
	completed := 0
	workItems := make([]WorkItem, 0, len(items))
	actors := map[string]Actor{
		"founder-principal":      {ActorID: "founder-principal", ActorType: "human", DisplayName: "创始治理者", Position: "董事长", State: "available"},
		"incorporation-agent":    {ActorID: "incorporation-agent", ActorType: "agent", DisplayName: agentLabel("incorporation-agent"), Position: agentLabel("incorporation-agent"), State: "available"},
		"governance-agent":       {ActorID: "governance-agent", ActorType: "agent", DisplayName: agentLabel("governance-agent"), Position: agentLabel("governance-agent"), State: "available"},
		"legal-compliance-agent": {ActorID: "legal-compliance-agent", ActorType: "agent", DisplayName: agentLabel("legal-compliance-agent"), Position: agentLabel("legal-compliance-agent"), State: "available"},
		"finance-agent":          {ActorID: "finance-agent", ActorType: "agent", DisplayName: agentLabel("finance-agent"), Position: agentLabel("finance-agent"), State: "available"},
		"audit-agent":            {ActorID: "audit-agent", ActorType: "agent", DisplayName: agentLabel("audit-agent"), Position: agentLabel("audit-agent"), State: "available"},
	}
	for _, item := range items {
		status := item.Effective
		if status == "" {
			status = item.Status
		}
		if status == "completed" {
			completed++
		}
		ownerType := ownerType(item.Participant)
		workItems = append(workItems, WorkItem{
			WorkItemID: fmt.Sprintf("WI-%02d", item.Sequence), Title: item.Capability,
			Kind: item.TaskType, Status: status, OwnerType: ownerType, OwnerID: item.Participant,
			Capability: item.Capability, Gate: item.Gate, RequiresMe: ownerType == "human" && status != "completed" && status != "locked",
			EvidenceRef: fmt.Sprintf("iaos:work-item:%s:%d", trace.CaseCode, item.Sequence),
		})
		if item.TaskType == "approval" && status != "locked" {
			workItems[len(workItems)-1].Review = approvalReview(trace, item)
		}
		if ownerType == "agent" {
			state := "available"
			if status == "ready" || status == "running" {
				state = "working"
			}
			if configured, ok := item.AgentAuth["status"].(string); ok && configured != "" && configured != "active" {
				state = configured
			}
			actors[item.Participant] = Actor{ActorID: item.Participant, ActorType: "agent", DisplayName: agentLabel(item.Participant), Position: agentLabel(item.Participant), State: state, WorkItemID: fmt.Sprintf("WI-%02d", item.Sequence)}
		}
	}
	actorList := make([]Actor, 0, len(actors))
	for _, id := range []string{"founder-principal", "incorporation-agent", "governance-agent", "legal-compliance-agent", "finance-agent", "audit-agent"} {
		if actor, exists := actors[id]; exists {
			actorList = append(actorList, actor)
		}
	}
	progress := 0
	if len(items) > 0 {
		progress = completed * 100 / len(items)
	}
	currency := trace.State.Currency
	if currency == "" {
		currency = "CNY"
	}
	projection := Projection{
		SchemaVersion: SchemaVersion, ProjectionID: fmt.Sprintf("gp-%s-%d", strings.ToLower(trace.CaseCode), cursor),
		TenantID: trace.State.TenantID, CaseCode: trace.CaseCode, WorldRunID: correlation,
		Chapter: chapter, SimTime: simTime, TimeScale: 0, Paused: true, Cursor: cursor,
		Scene:     Scene{SceneID: sceneForChapter(chapter), Mode: "2.5d", Theme: "industrial-warm"},
		Lifecycle: Lifecycle{State: trace.State.State, CurrentStep: currentStep(workItems), Progress: progress},
		Buildings: buildingsForChapter(chapter), Actors: actorList, WorkItems: workItems,
		Resources: Resources{
			FounderCash:      Money{Value: "0", Currency: currency, Scale: 2},
			CompanyCash:      Money{Value: fmt.Sprint(trace.State.ContributionMinor), Currency: currency, Scale: 2},
			CapitalCommitted: Money{Value: fmt.Sprint(trace.State.CommitmentMinor), Currency: currency, Scale: 2},
			CapitalPaid:      Money{Value: fmt.Sprint(trace.State.ContributionMinor), Currency: currency, Scale: 2},
			BudgetAuthorized: Money{Value: fmt.Sprint(trace.State.BudgetMinor), Currency: currency, Scale: 2}, RiskLevel: "low",
		},
		Finance: FinanceOpening{
			Roles: []string{}, TrialBalance: []FinanceTrialBalanceLine{},
			BankJournal: []FinanceBankJournalLine{}, GeneralLedger: []FinanceGeneralLedgerLine{},
			OpeningBalanceSheet: FinanceOpeningBalanceSheet{
				Assets: []FinanceStatementLine{}, Liabilities: []FinanceStatementLine{}, Equity: []FinanceStatementLine{},
			},
		},
		Exchanges: exchanges, Brand: Brand{Status: "selected", CompanyName: trace.State.ProposedName, PrimaryColor: "#167C80"},
		Notifications: []Notification{{NotificationID: fmt.Sprintf("notice-%d", cursor), Severity: "info", Message: currentStep(workItems)}},
		EvidenceRefs:  []EvidenceRef{{Ref: "iaos:trace:" + trace.CaseCode, Kind: "incorporation_trace"}, {Ref: "iaos:work-items:" + trace.CaseCode, Kind: "process_work_items"}},
	}
	if trace.State.State == "enterprise_operational_ready" {
		projection.Brand.Status = "approved"
	}
	for index, discrepancy := range trace.State.Discrepancies {
		projection.Notifications = append(projection.Notifications, Notification{
			NotificationID: fmt.Sprintf("discrepancy-%d", index+1), Severity: "warning",
			Message: discrepancy, ActionRef: "iaos:discrepancy:" + trace.CaseCode,
		})
	}
	if projection.WorldRunID == "" {
		projection.WorldRunID = "case-" + trace.CaseCode
	}
	return projection, projection.Validate()
}

// AttachFinanceOpening enriches a read-only game projection with committed
// finance facts. It never infers a journal from company cash.
func AttachFinanceOpening(projection *Projection, opening iaosclient.FinanceOpening) error {
	if projection == nil {
		return fmt.Errorf("projection is required")
	}
	if opening.CaseCode != "" && opening.CaseCode != projection.CaseCode {
		return fmt.Errorf("finance opening case %q does not match projection case %q", opening.CaseCode, projection.CaseCode)
	}
	lines := make([]FinanceTrialBalanceLine, 0, len(opening.TrialBalance))
	for _, line := range opening.TrialBalance {
		lines = append(lines, FinanceTrialBalanceLine{
			AccountCode: line.AccountCode, AccountName: line.AccountName,
			AccountClass: line.AccountClass, NormalBalance: line.NormalBalance,
			DebitMinor: line.DebitMinor, CreditMinor: line.CreditMinor,
			BalanceMinor: line.BalanceMinor,
		})
	}
	bankJournal := make([]FinanceBankJournalLine, 0, len(opening.BankJournal))
	for _, line := range opening.BankJournal {
		bankJournal = append(bankJournal, FinanceBankJournalLine{
			EntryNo: line.EntryNo, BusinessDate: line.BusinessDate,
			Description: line.Description, DebitMinor: line.DebitMinor,
			CreditMinor: line.CreditMinor, BalanceMinor: line.BalanceMinor,
			SourceType: line.SourceType, SourceRef: line.SourceRef,
			EvidenceRef: line.EvidenceRef,
		})
	}
	generalLedger := make([]FinanceGeneralLedgerLine, 0, len(opening.GeneralLedger))
	for _, line := range opening.GeneralLedger {
		generalLedger = append(generalLedger, FinanceGeneralLedgerLine{
			AccountCode: line.AccountCode, AccountName: line.AccountName,
			OpeningBalanceMinor: line.OpeningBalanceMinor, DebitMinor: line.DebitMinor,
			CreditMinor: line.CreditMinor, ClosingBalanceMinor: line.ClosingBalanceMinor,
		})
	}
	statementLines := func(lines []iaosclient.FinanceStatementLine) []FinanceStatementLine {
		out := make([]FinanceStatementLine, 0, len(lines))
		for _, line := range lines {
			out = append(out, FinanceStatementLine{
				AccountCode: line.AccountCode,
				AccountName: line.AccountName,
				AmountMinor: line.AmountMinor,
			})
		}
		return out
	}
	projection.Finance = FinanceOpening{
		Ready: opening.FinanceOpeningReady, OrganizationCode: opening.OrganizationCode,
		OrganizationStatus: opening.OrganizationStatus, Roles: append([]string(nil), opening.Roles...),
		BookCode: opening.BookCode, AccountingStandard: opening.AccountingStandard,
		BookName: opening.BookName, FiscalYear: opening.FiscalYear,
		AccountingPeriods: func() []FinanceAccountingPeriod {
			out:=make([]FinanceAccountingPeriod,0,len(opening.AccountingPeriods))
			for _, period:=range opening.AccountingPeriods { out=append(out,FinanceAccountingPeriod{PeriodCode:period.PeriodCode,StartsOn:period.StartsOn,EndsOn:period.EndsOn,Status:period.Status}) }
			return out
		}(),
		FunctionalCurrency: opening.FunctionalCurrency, PeriodCode: opening.PeriodCode,
		PeriodStatus: opening.PeriodStatus, JournalEntryNo: opening.JournalEntryNo,
		JournalStatus: opening.JournalStatus, DebitMinor: opening.DebitMinor,
		CreditMinor: opening.CreditMinor, TrialBalance: lines,
		BankJournal: bankJournal, GeneralLedger: generalLedger,
		OpeningBalanceSheet: FinanceOpeningBalanceSheet{
			AsOf: opening.OpeningBalanceSheet.AsOf, Currency: opening.OpeningBalanceSheet.Currency,
			Assets:                statementLines(opening.OpeningBalanceSheet.Assets),
			Liabilities:           statementLines(opening.OpeningBalanceSheet.Liabilities),
			Equity:                statementLines(opening.OpeningBalanceSheet.Equity),
			TotalAssetsMinor:      opening.OpeningBalanceSheet.TotalAssetsMinor,
			TotalLiabilitiesMinor: opening.OpeningBalanceSheet.TotalLiabilitiesMinor,
			TotalEquityMinor:      opening.OpeningBalanceSheet.TotalEquityMinor,
			Balanced:              opening.OpeningBalanceSheet.Balanced,
		},
		EvidenceRef: opening.EvidenceRef,
	}
	if opening.EvidenceRef != "" {
		projection.EvidenceRefs = append(projection.EvidenceRefs, EvidenceRef{Ref: opening.EvidenceRef, Kind: "finance_opening"})
	}
	return projection.Validate()
}

func approvalReview(trace iaosclient.IncorporationTrace, item iaosclient.IncorporationWorkItem) *ApprovalReview {
	ref := fmt.Sprintf("iaos:approval-subject:%s:%s", trace.CaseCode, item.Gate)
	base := &ApprovalReview{
		DocumentType: "governed_action", Title: "待审批事项", Summary: "请核对事项内容及执行影响后作出决定。",
		PreparedBy: "IAOS Runtime", Status: "待批准", Risks: []string{"批准后将形成受审计的业务事实并推进企业设立进程。"},
		Effect: "批准并执行 " + item.Capability, EvidenceRef: ref,
	}
	money := func(v int64) string { return fmt.Sprintf("¥%.2f", float64(v)/100) }
	switch item.Gate {
	case "G1":
		base.DocumentType, base.Title = "founder_resolution", "创始人设立决议草案"
		base.Summary = "确认创业构想转为正式企业设立行动，并授权进入登记准备。"
		base.PreparedBy = "企业设立专员"
		base.Fields = []ReviewField{{Label: "拟设企业", Value: trace.State.ProposedName}, {Label: "注册地址", Value: trace.State.RegisteredAddress}, {Label: "经营范围", Value: trace.State.BusinessScope}}
		base.Effect = "决议生效，允许登记资本承诺并准备设立材料。"
		for _, run := range trace.AgentRuns {
			if run.Capability != "founder.resolution.prepare" || run.Status != "completed" {
				continue
			}
			var output struct {
				Resolution struct {
					Title               string `json:"title"`
					PreparedBy          string `json:"prepared_by"`
					ResolutionObjective string `json:"resolution_objective"`
					KeyProposals        string `json:"key_proposals"`
					RiskNotes           string `json:"risk_notes"`
				} `json:"governance_resolution"`
			}
			if json.Unmarshal(run.Output, &output) == nil && output.Resolution.Title != "" {
				base.Title, base.PreparedBy = output.Resolution.Title, output.Resolution.PreparedBy
				base.Fields = append(base.Fields,
					ReviewField{Label: "决议目标", Value: output.Resolution.ResolutionObjective},
					ReviewField{Label: "核心提案", Value: output.Resolution.KeyProposals},
				)
				if output.Resolution.RiskNotes != "" {
					base.Risks = []string{output.Resolution.RiskNotes}
				}
				base.EvidenceRef = "iaos:agent-run:" + run.ID
			}
		}
	case "G2":
		base.DocumentType, base.Title = "registration_application", "企业设立登记申请"
		base.Summary, base.PreparedBy = "核准向模拟登记机构提交的企业法定身份资料。", "法务合规专员"
		base.Fields = []ReviewField{{Label: "企业名称", Value: trace.State.ProposedName}, {Label: "注册地址", Value: trace.State.RegisteredAddress}, {Label: "经营范围", Value: trace.State.BusinessScope}, {Label: "认缴资本", Value: money(trace.State.CommitmentMinor)}}
		base.Effect = "向外部登记世界提交申请；批准不等于登记成功。"
	case "G3":
		base.DocumentType, base.Title = "bank_account_application", "企业银行账户开户申请"
		base.Summary, base.PreparedBy = "以已登记法律主体身份申请企业基本账户。", "财务负责人"
		base.Fields = []ReviewField{{Label: "法律主体", Value: trace.State.LegalEntityCode}, {Label: "企业名称", Value: trace.State.ProposedName}}
		base.Effect = "向模拟银行提交开户申请；账户须等待银行回传后生效。"
	case "G4":
		base.DocumentType, base.Title = "capital_contribution_verification", "实缴资本到账核验"
		base.Summary, base.PreparedBy = "确认实际到账资金并与创始人认缴承诺核对。", "财务负责人"
		base.Fields = []ReviewField{{Label: "认缴资本", Value: money(trace.State.CommitmentMinor)}, {Label: "收款账户", Value: trace.State.BankAccountCode}}
		base.Risks = []string{"输入的实缴金额必须与银行到账事实一致，且不得超过认缴承诺。"}
		base.Effect = "把本次输入金额确认为公司实缴资本和可用现金。"
	case "G5":
		base.DocumentType, base.Title = "executive_appointment", "首届管理层任命议案"
		base.Summary, base.PreparedBy = "审议已获候选人接受的首届管理团队任命。", "治理组织专员"
		base.Fields = []ReviewField{{Label: "组织", Value: trace.State.OrganizationCode}, {Label: "拟任岗位", Value: "CEO、CFO、工厂项目负责人"}, {Label: "接受状态", Value: "候选人已接受"}}
		base.Effect = "批准后生成正式管理层任命记录。"
	case "G6":
		base.DocumentType, base.Title = "operating_mandate", "管理层经营授权书"
		base.Summary, base.PreparedBy = "向已任命管理层授予受范围、限额和期限约束的经营权限。", "治理组织专员"
		base.Fields = []ReviewField{{Label: "任命记录", Value: trace.State.AppointmentCode}, {Label: "授权对象", Value: "CEO、CFO、工厂项目负责人"}, {Label: "治理原则", Value: "岗位职责内执行，超限事项升级审批"}}
		base.Effect = "经营授权生效，管理层可在 mandate 范围内执行工作。"
	case "G7":
		base.DocumentType, base.Title = "initial_budget", "首年启动预算审批案"
		base.Summary, base.PreparedBy = "审议企业开业阶段的首年预算授权上限。", "财务负责人"
		base.Fields = []ReviewField{{Label: "已核验公司现金", Value: money(trace.State.ContributionMinor)}, {Label: "当前预算授权", Value: money(trace.State.BudgetMinor)}}
		base.Risks = []string{"本次输入的预算不得超过已核验公司现金；预算授权不等于立即支出。"}
		base.Effect = "把本次输入金额确认为首年预算授权上限。"
	}
	return base
}

func agentLabel(id string) string {
	return map[string]string{"incorporation-agent": "企业设立专员", "governance-agent": "治理组织专员", "legal-compliance-agent": "法务合规专员", "finance-agent": "财务负责人", "audit-agent": "独立审计专员"}[id]
}

func currentStep(items []WorkItem) string {
	for _, item := range items {
		if item.Status != "completed" && item.Status != "locked" {
			return item.Capability
		}
	}
	return "enterprise_operational_ready"
}
