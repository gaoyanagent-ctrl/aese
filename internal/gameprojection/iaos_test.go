package gameprojection

import (
	"encoding/json"
	"testing"

	"github.com/industrial-ai/iaos-aese/internal/iaosclient"
)

func TestFromIAOSUsesCommittedWorkItemsAndWorldCursor(t *testing.T) {
	trace := iaosclient.IncorporationTrace{
		SchemaVersion: "1.0", CaseCode: "INC-LIVE-001", Verified: true,
		State: iaosclient.IncorporationState{
			CaseCode: "INC-LIVE-001", TenantID: "tenant-live", State: "registration_submitted",
			CommitmentMinor: 100000000, Currency: "CNY", ProposedName: "澄流热管理有限公司",
		},
		Journal:        []iaosclient.IncorporationJournalEntry{{Sequence: 1, CorrelationID: "corr-live", CreatedAt: "2026-07-27T10:00:00Z"}},
		WorldExchanges: []iaosclient.IncorporationWorldExchange{{Cursor: 42, MessageID: "intent-1", Kind: "intent", CorrelationID: "corr-live", RecordedAt: "2026-07-27T10:00:01Z"}},
	}
	items := []iaosclient.IncorporationWorkItem{
		{Sequence: 1, Capability: "incorporation.case.open", TaskType: "human_task", Participant: "founder-principal", Effective: "completed"},
		{Sequence: 2, Capability: "founder.resolution.prepare", TaskType: "agent_task", Participant: "incorporation-agent", Effective: "ready"},
	}
	projection, err := FromIAOS(trace, items)
	if err != nil {
		t.Fatal(err)
	}
	if projection.Cursor != 42 || projection.Chapter != "registration" || projection.Lifecycle.Progress != 50 {
		t.Fatalf("unexpected projection %#v", projection)
	}
	if len(projection.WorkItems) != 2 || projection.WorkItems[1].EvidenceRef == "" || projection.Brand.CompanyName != "澄流热管理有限公司" {
		t.Fatalf("projection lost governed facts %#v", projection)
	}
}

func TestFromIAOSExposesFounderResolutionAsG1ApprovalSubject(t *testing.T) {
	output, _ := json.Marshal(map[string]any{"governance_resolution": map[string]any{
		"title": "关于澄流公司设立的创始人决议草案", "prepared_by": "incorporation-agent",
		"resolution_objective": "建立工业热管理企业", "key_proposals": "注册、开户并建立治理团队", "risk_notes": "外部结果须可信回传",
	}})
	trace := iaosclient.IncorporationTrace{
		SchemaVersion: "1.0", CaseCode: "INC-REVIEW-001", Verified: true,
		State:     iaosclient.IncorporationState{CaseCode: "INC-REVIEW-001", TenantID: "tenant-review", State: "incorporation_case_opened", ProposedName: "澄流热管理有限公司", RegisteredAddress: "苏州市工业园区", BusinessScope: "热管理产品研发"},
		Journal:   []iaosclient.IncorporationJournalEntry{{Sequence: 1, CorrelationID: "corr-review", CreatedAt: "2026-07-28T10:00:00Z"}},
		AgentRuns: []iaosclient.IncorporationAgentRun{{ID: "run-1", Capability: "founder.resolution.prepare", Status: "completed", Output: output}},
	}
	items := []iaosclient.IncorporationWorkItem{
		{Sequence: 1, Capability: "incorporation.case.open", TaskType: "human_task", Participant: "founder-principal", Effective: "completed"},
		{Sequence: 2, Capability: "founder.resolution.prepare", TaskType: "agent_task", Participant: "incorporation-agent", Effective: "completed"},
		{Sequence: 3, Capability: "founder.resolution.approve", TaskType: "approval", Participant: "founder-principal", Gate: "G1", Effective: "waiting_approval"},
	}
	projection, err := FromIAOS(trace, items)
	if err != nil {
		t.Fatal(err)
	}
	review := projection.WorkItems[2].Review
	if review == nil || review.Title != "关于澄流公司设立的创始人决议草案" || review.PreparedBy != "incorporation-agent" || review.EvidenceRef != "iaos:agent-run:run-1" {
		t.Fatalf("founder resolution not projected: %#v", review)
	}
	if len(review.Fields) < 5 || review.Fields[3].Value != "建立工业热管理企业" || review.Risks[0] != "外部结果须可信回传" {
		t.Fatalf("founder resolution details lost: %#v", review)
	}
}

func TestFromIAOSUnlocksHeadquartersForOrganizationMission(t *testing.T) {
	trace := iaosclient.IncorporationTrace{
		SchemaVersion: "1.0", CaseCode: "INC-ORG-001", Verified: true,
		State: iaosclient.IncorporationState{
			CaseCode: "INC-ORG-001", TenantID: "tenant-org", State: "capital_contribution_verified",
		},
		Journal: []iaosclient.IncorporationJournalEntry{{Sequence: 10, CorrelationID: "corr-org", CreatedAt: "2026-07-28T10:00:00Z"}},
	}
	items := []iaosclient.IncorporationWorkItem{{
		Sequence: 11, Capability: "organization.establish", TaskType: "system_task",
		Participant: "iaos-runtime", Effective: "ready",
	}}
	projection, err := FromIAOS(trace, items)
	if err != nil {
		t.Fatal(err)
	}
	for _, building := range projection.Buildings {
		if building.Kind == "headquarters" {
			if !building.Available {
				t.Fatal("organization mission points to a locked headquarters")
			}
			return
		}
	}
	t.Fatal("headquarters missing from game projection")
}

func TestFromIAOSProvidesReviewSubjectForEveryApprovalGate(t *testing.T) {
	for _, gate := range []string{"G1", "G2", "G3", "G4", "G5", "G6", "G7"} {
		trace := iaosclient.IncorporationTrace{
			SchemaVersion: "1.0", CaseCode: "INC-" + gate, Verified: true,
			State:   iaosclient.IncorporationState{TenantID: "tenant-review", State: "incorporation_case_opened", ProposedName: "测试企业", Currency: "CNY"},
			Journal: []iaosclient.IncorporationJournalEntry{{Sequence: 1, CreatedAt: "2026-07-28T10:00:00Z"}},
		}
		items := []iaosclient.IncorporationWorkItem{{Sequence: 1, Capability: "test.approve", TaskType: "approval", Participant: "founder-principal", Gate: gate, Effective: "waiting_approval"}}
		projection, err := FromIAOS(trace, items)
		if err != nil {
			t.Fatalf("%s: %v", gate, err)
		}
		if projection.WorkItems[0].Review == nil || projection.WorkItems[0].Review.Title == "待审批事项" {
			t.Fatalf("%s has no specific approval subject: %#v", gate, projection.WorkItems[0].Review)
		}
	}
}

func TestFromIAOSRejectsUnverifiedTrace(t *testing.T) {
	_, err := FromIAOS(iaosclient.IncorporationTrace{SchemaVersion: "1.0", CaseCode: "INC-1"}, nil)
	if err == nil {
		t.Fatal("unverified trace accepted")
	}
}

func TestAttachFinanceOpeningRequiresBalancedCommittedJournal(t *testing.T) {
	projection := Projection{
		SchemaVersion: SchemaVersion, ProjectionID: "gp-finance", TenantID: "tenant-finance",
		CaseCode: "INC-FINANCE", WorldRunID: "world-finance", SimTime: "2026-07-28T10:00:00Z",
		Scene:     Scene{SceneID: "headquarters-v1", Mode: "2.5d", Theme: "industrial-warm"},
		Lifecycle: Lifecycle{State: "capital_contribution_verified", CurrentStep: "organization.establish", Progress: 55},
		Finance: FinanceOpening{
			Roles: []string{}, TrialBalance: []FinanceTrialBalanceLine{},
			BankJournal: []FinanceBankJournalLine{}, GeneralLedger: []FinanceGeneralLedgerLine{},
			OpeningBalanceSheet: FinanceOpeningBalanceSheet{
				Assets: []FinanceStatementLine{}, Liabilities: []FinanceStatementLine{}, Equity: []FinanceStatementLine{},
			},
		},
	}
	opening := iaosclient.FinanceOpening{
		CaseCode: "INC-FINANCE", FinanceOpeningReady: true,
		OrganizationCode: "FIN-INC-FINANCE", OrganizationStatus: "active",
		Roles: []string{"cfo", "general_ledger_accountant"}, BookCode: "BOOK-INC-FINANCE",
		AccountingStandard: "CAS-BE", FunctionalCurrency: "CNY", PeriodCode: "2026-07",
		PeriodStatus: "open", JournalEntryNo: "OPEN-INC-FINANCE", JournalStatus: "posted",
		DebitMinor: 100_000_000, CreditMinor: 100_000_000,
		TrialBalance: []iaosclient.FinanceTrialBalanceLine{
			{AccountCode: "1002", AccountName: "银行存款", AccountClass: "asset", NormalBalance: "debit", DebitMinor: 100_000_000, BalanceMinor: 100_000_000},
			{AccountCode: "4001", AccountName: "实收资本", AccountClass: "equity", NormalBalance: "credit", CreditMinor: 100_000_000, BalanceMinor: -100_000_000},
		},
		BankJournal: []iaosclient.FinanceBankJournalLine{{
			EntryNo: "OPEN-INC-FINANCE", BusinessDate: "2026-07-28",
			Description: "实收资本到账", DebitMinor: 100_000_000,
			BalanceMinor: 100_000_000, SourceType: "capital_contribution",
			SourceRef: "INC-FINANCE", EvidenceRef: "iaos:incorporation:INC-FINANCE:capital_contribution",
		}},
		GeneralLedger: []iaosclient.FinanceGeneralLedgerLine{
			{AccountCode: "1002", AccountName: "银行存款", DebitMinor: 100_000_000, ClosingBalanceMinor: 100_000_000},
			{AccountCode: "4001", AccountName: "实收资本", CreditMinor: 100_000_000, ClosingBalanceMinor: 100_000_000},
		},
		OpeningBalanceSheet: iaosclient.FinanceOpeningBalanceSheet{
			AsOf: "2026-07-28", Currency: "CNY",
			Assets:           []iaosclient.FinanceStatementLine{{AccountCode: "1002", AccountName: "银行存款", AmountMinor: 100_000_000}},
			Liabilities:      []iaosclient.FinanceStatementLine{},
			Equity:           []iaosclient.FinanceStatementLine{{AccountCode: "4001", AccountName: "实收资本", AmountMinor: 100_000_000}},
			TotalAssetsMinor: 100_000_000, TotalEquityMinor: 100_000_000, Balanced: true,
		},
		EvidenceRef: "iaos:finance:INC-FINANCE:OPEN-INC-FINANCE",
	}
	if err := AttachFinanceOpening(&projection, opening); err != nil {
		t.Fatal(err)
	}
	if !projection.Finance.Ready || len(projection.Finance.TrialBalance) != 2 ||
		len(projection.Finance.BankJournal) != 1 || len(projection.Finance.GeneralLedger) != 2 ||
		!projection.Finance.OpeningBalanceSheet.Balanced ||
		projection.Finance.DebitMinor != projection.Finance.CreditMinor {
		t.Fatalf("finance opening facts lost: %#v", projection.Finance)
	}
	opening.CreditMinor--
	if err := AttachFinanceOpening(&projection, opening); err == nil {
		t.Fatal("unbalanced finance opening accepted")
	}
}
