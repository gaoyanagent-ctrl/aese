// Package gameprojection defines the read-only contract consumed by the
// Enterprise Genesis 2D/2.5D presentation layer.
package gameprojection

import (
	"fmt"
	"time"
)

const SchemaVersion = "1.0"

type Projection struct {
	SchemaVersion string         `json:"schema_version"`
	ProjectionID  string         `json:"projection_id"`
	TenantID      string         `json:"tenant_id"`
	CaseCode      string         `json:"case_code"`
	WorldRunID    string         `json:"world_run_id"`
	Chapter       string         `json:"chapter"`
	SimTime       string         `json:"sim_time"`
	TimeScale     int            `json:"time_scale"`
	Paused        bool           `json:"paused"`
	Cursor        int64          `json:"cursor"`
	Scene         Scene          `json:"world_scene"`
	Lifecycle     Lifecycle      `json:"lifecycle"`
	Buildings     []Building     `json:"buildings"`
	Actors        []Actor        `json:"actors"`
	WorkItems     []WorkItem     `json:"work_items"`
	Resources     Resources      `json:"resources"`
	Finance       FinanceOpening `json:"finance_opening"`
	Exchanges     []Exchange     `json:"exchanges"`
	Brand         Brand          `json:"brand"`
	Notifications []Notification `json:"notifications"`
	EvidenceRefs  []EvidenceRef  `json:"evidence_refs"`
}

type Scene struct {
	SceneID string `json:"scene_id"`
	Mode    string `json:"mode"`
	Theme   string `json:"theme"`
}

type Lifecycle struct {
	State       string `json:"state"`
	CurrentStep string `json:"current_step"`
	Progress    int    `json:"progress"`
	BlockedBy   string `json:"blocked_by,omitempty"`
}

type Building struct {
	Code      string `json:"code"`
	Kind      string `json:"kind"`
	Label     string `json:"label"`
	State     string `json:"state"`
	X         int    `json:"x"`
	Y         int    `json:"y"`
	Evidence  string `json:"evidence_ref,omitempty"`
	Available bool   `json:"available"`
}

type Actor struct {
	ActorID     string `json:"actor_id"`
	ActorType   string `json:"actor_type"`
	DisplayName string `json:"display_name"`
	Position    string `json:"position"`
	State       string `json:"state"`
	WorkItemID  string `json:"work_item_id,omitempty"`
}

type WorkItem struct {
	WorkItemID  string          `json:"work_item_id"`
	Title       string          `json:"title"`
	Kind        string          `json:"kind"`
	Status      string          `json:"status"`
	OwnerType   string          `json:"owner_type"`
	OwnerID     string          `json:"owner_id"`
	Capability  string          `json:"capability"`
	Gate        string          `json:"gate,omitempty"`
	RequiresMe  bool            `json:"requires_me"`
	EvidenceRef string          `json:"evidence_ref"`
	Review      *ApprovalReview `json:"approval_review,omitempty"`
}

type ApprovalReview struct {
	DocumentType string        `json:"document_type"`
	Title        string        `json:"title"`
	Summary      string        `json:"summary"`
	PreparedBy   string        `json:"prepared_by"`
	Status       string        `json:"status"`
	Fields       []ReviewField `json:"fields"`
	Risks        []string      `json:"risks"`
	Effect       string        `json:"approval_effect"`
	EvidenceRef  string        `json:"evidence_ref"`
}

type ReviewField struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type Money struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
	Scale    int    `json:"scale"`
}

type Resources struct {
	FounderCash      Money  `json:"founder_cash"`
	CompanyCash      Money  `json:"company_cash"`
	CapitalCommitted Money  `json:"capital_committed"`
	CapitalPaid      Money  `json:"capital_paid"`
	BudgetAuthorized Money  `json:"budget_authorized"`
	RiskLevel        string `json:"risk_level"`
}

type FinanceOpening struct {
	Ready               bool                       `json:"ready"`
	OrganizationCode    string                     `json:"organization_code,omitempty"`
	OrganizationStatus  string                     `json:"organization_status,omitempty"`
	Roles               []string                   `json:"roles"`
	BookCode            string                     `json:"book_code,omitempty"`
	BookName            string                     `json:"book_name,omitempty"`
	FiscalYear          int                        `json:"fiscal_year,omitempty"`
	AccountingPeriods   []FinanceAccountingPeriod  `json:"accounting_periods"`
	AccountingStandard  string                     `json:"accounting_standard,omitempty"`
	FunctionalCurrency  string                     `json:"functional_currency,omitempty"`
	PeriodCode          string                     `json:"period_code,omitempty"`
	PeriodStatus        string                     `json:"period_status,omitempty"`
	JournalEntryNo      string                     `json:"journal_entry_no,omitempty"`
	JournalStatus       string                     `json:"journal_status,omitempty"`
	DebitMinor          int64                      `json:"debit_minor"`
	CreditMinor         int64                      `json:"credit_minor"`
	TrialBalance        []FinanceTrialBalanceLine  `json:"trial_balance"`
	BankJournal         []FinanceBankJournalLine   `json:"bank_journal"`
	GeneralLedger       []FinanceGeneralLedgerLine `json:"general_ledger"`
	OpeningBalanceSheet FinanceOpeningBalanceSheet `json:"opening_balance_sheet"`
	EvidenceRef         string                     `json:"evidence_ref,omitempty"`
}

type FinanceAccountingPeriod struct {
	PeriodCode string `json:"period_code"`
	StartsOn string `json:"starts_on"`
	EndsOn string `json:"ends_on"`
	Status string `json:"status"`
}

type FinanceTrialBalanceLine struct {
	AccountCode   string `json:"account_code"`
	AccountName   string `json:"account_name"`
	AccountClass  string `json:"account_class"`
	NormalBalance string `json:"normal_balance"`
	DebitMinor    int64  `json:"debit_minor"`
	CreditMinor   int64  `json:"credit_minor"`
	BalanceMinor  int64  `json:"balance_minor"`
}

type FinanceBankJournalLine struct {
	EntryNo      string `json:"entry_no"`
	BusinessDate string `json:"business_date"`
	Description  string `json:"description"`
	DebitMinor   int64  `json:"debit_minor"`
	CreditMinor  int64  `json:"credit_minor"`
	BalanceMinor int64  `json:"balance_minor"`
	SourceType   string `json:"source_type"`
	SourceRef    string `json:"source_ref"`
	EvidenceRef  string `json:"evidence_ref"`
}

type FinanceGeneralLedgerLine struct {
	AccountCode         string `json:"account_code"`
	AccountName         string `json:"account_name"`
	OpeningBalanceMinor int64  `json:"opening_balance_minor"`
	DebitMinor          int64  `json:"debit_minor"`
	CreditMinor         int64  `json:"credit_minor"`
	ClosingBalanceMinor int64  `json:"closing_balance_minor"`
}

type FinanceStatementLine struct {
	AccountCode string `json:"account_code"`
	AccountName string `json:"account_name"`
	AmountMinor int64  `json:"amount_minor"`
}

type FinanceOpeningBalanceSheet struct {
	AsOf                  string                 `json:"as_of"`
	Currency              string                 `json:"currency"`
	Assets                []FinanceStatementLine `json:"assets"`
	Liabilities           []FinanceStatementLine `json:"liabilities"`
	Equity                []FinanceStatementLine `json:"equity"`
	TotalAssetsMinor      int64                  `json:"total_assets_minor"`
	TotalLiabilitiesMinor int64                  `json:"total_liabilities_minor"`
	TotalEquityMinor      int64                  `json:"total_equity_minor"`
	Balanced              bool                   `json:"balanced"`
}

type Exchange struct {
	ExchangeID  string `json:"exchange_id"`
	Kind        string `json:"kind"`
	Status      string `json:"status"`
	Correlation string `json:"correlation_id"`
	EvidenceRef string `json:"evidence_ref"`
	OccurredAt  string `json:"occurred_at"`
}

type Brand struct {
	Status       string `json:"status"`
	CompanyName  string `json:"company_name,omitempty"`
	LogoAssetID  string `json:"logo_asset_id,omitempty"`
	PrimaryColor string `json:"primary_color,omitempty"`
}

type Notification struct {
	NotificationID string `json:"notification_id"`
	Severity       string `json:"severity"`
	Message        string `json:"message"`
	ActionRef      string `json:"action_ref,omitempty"`
}

type EvidenceRef struct {
	Ref  string `json:"ref"`
	Kind string `json:"kind"`
}

func (p *Projection) Validate() error {
	if p.SchemaVersion != SchemaVersion {
		return fmt.Errorf("unsupported schema_version %q", p.SchemaVersion)
	}
	if p.ProjectionID == "" || p.TenantID == "" || p.CaseCode == "" || p.WorldRunID == "" {
		return fmt.Errorf("projection identity is incomplete")
	}
	if _, err := time.Parse(time.RFC3339, p.SimTime); err != nil {
		return fmt.Errorf("sim_time must be RFC3339: %w", err)
	}
	if p.TimeScale != 0 && p.TimeScale != 1 && p.TimeScale != 2 && p.TimeScale != 4 {
		return fmt.Errorf("time_scale must be 0, 1, 2, or 4")
	}
	if p.Cursor < 0 || p.Lifecycle.Progress < 0 || p.Lifecycle.Progress > 100 {
		return fmt.Errorf("cursor or progress is outside its allowed range")
	}
	if p.Scene.Mode != "2d" && p.Scene.Mode != "2.5d" {
		return fmt.Errorf("world_scene.mode must be 2d or 2.5d")
	}
	for _, item := range p.WorkItems {
		if item.WorkItemID == "" || item.OwnerID == "" || item.Capability == "" || item.EvidenceRef == "" {
			return fmt.Errorf("work item %q is not traceable", item.WorkItemID)
		}
	}
	if p.Finance.Ready {
		if p.Finance.OrganizationCode == "" || p.Finance.BookCode == "" || p.Finance.JournalEntryNo == "" {
			return fmt.Errorf("ready finance opening is missing organization, book, or journal identity")
		}
		if p.Finance.DebitMinor <= 0 || p.Finance.DebitMinor != p.Finance.CreditMinor {
			return fmt.Errorf("ready finance opening is not balanced")
		}
		if len(p.Finance.BankJournal) == 0 || len(p.Finance.GeneralLedger) < 2 ||
			!p.Finance.OpeningBalanceSheet.Balanced {
			return fmt.Errorf("ready finance opening is missing reconciled journal, ledger, or balance sheet")
		}
	}
	return nil
}
