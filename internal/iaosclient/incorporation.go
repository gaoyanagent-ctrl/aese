package iaosclient

import (
	"context"
	"encoding/json"
	"net/url"
)

type IncorporationState struct {
	CaseCode          string   `json:"case_code"`
	TenantID          string   `json:"tenant_id"`
	State             string   `json:"state"`
	LegalEntityCode   string   `json:"legal_entity_code"`
	BankAccountCode   string   `json:"bank_account_code"`
	OrganizationCode  string   `json:"organization_code"`
	AppointmentCode   string   `json:"appointment_code"`
	OperatingMandate  string   `json:"operating_mandate_code"`
	CommitmentMinor   int64    `json:"commitment_minor"`
	ContributionMinor int64    `json:"contribution_minor"`
	BudgetMinor       int64    `json:"budget_authorized_minor"`
	Currency          string   `json:"currency"`
	ProposedName      string   `json:"proposed_company_name"`
	RegisteredAddress string   `json:"registered_address"`
	BusinessScope     string   `json:"business_scope"`
	Discrepancies     []string `json:"discrepancies"`
}

type IncorporationJournalEntry struct {
	Sequence      int    `json:"sequence"`
	Capability    string `json:"capability_code"`
	ActorType     string `json:"actor_type"`
	ActorID       string `json:"actor_id"`
	CorrelationID string `json:"correlation_id"`
	CreatedAt     string `json:"created_at"`
}

type IncorporationWorldExchange struct {
	Cursor        int64  `json:"cursor"`
	MessageID     string `json:"message_id"`
	Kind          string `json:"kind"`
	CorrelationID string `json:"correlation_id"`
	PayloadType   string `json:"payload_type"`
	RecordedAt    string `json:"recorded_at"`
}

type IncorporationTrace struct {
	SchemaVersion  string                       `json:"schema_version"`
	CaseCode       string                       `json:"case_code"`
	State          IncorporationState           `json:"state"`
	Journal        []IncorporationJournalEntry  `json:"journal"`
	WorldExchanges []IncorporationWorldExchange `json:"world_exchanges"`
	AgentRuns      []IncorporationAgentRun      `json:"agent_runs"`
	Verified       bool                         `json:"verified"`
}

type IncorporationAgentRun struct {
	ID               string          `json:"id"`
	WorkItemSequence int             `json:"work_item_sequence"`
	AgentID          string          `json:"agent_id"`
	Capability       string          `json:"capability_code"`
	Output           json.RawMessage `json:"output"`
	Status           string          `json:"status"`
	CompletedAt      string          `json:"completed_at"`
}

type IncorporationWorkItem struct {
	Sequence    int            `json:"sequence"`
	Capability  string         `json:"capability_code"`
	TaskType    string         `json:"task_type"`
	Participant string         `json:"participant_id"`
	Gate        string         `json:"gate"`
	Status      string         `json:"status"`
	Effective   string         `json:"effective_status"`
	Correlation string         `json:"correlation_id"`
	AgentAuth   map[string]any `json:"agent_authorization"`
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

type FinanceOpening struct {
	CaseCode            string                     `json:"case_code"`
	FinanceOpeningReady bool                       `json:"finance_opening_ready"`
	OrganizationCode    string                     `json:"organization_code"`
	OrganizationStatus  string                     `json:"organization_status"`
	Roles               []string                   `json:"roles"`
	BookCode            string                     `json:"book_code"`
	BookName            string                     `json:"book_name"`
	FiscalYear          int                        `json:"fiscal_year"`
	AccountingPeriods   []FinanceAccountingPeriod  `json:"accounting_periods"`
	AccountingStandard  string                     `json:"accounting_standard"`
	FunctionalCurrency  string                     `json:"functional_currency"`
	PeriodCode          string                     `json:"period_code"`
	PeriodStatus        string                     `json:"period_status"`
	JournalEntryNo      string                     `json:"journal_entry_no"`
	JournalStatus       string                     `json:"journal_status"`
	DebitMinor          int64                      `json:"debit_minor"`
	CreditMinor         int64                      `json:"credit_minor"`
	TrialBalance        []FinanceTrialBalanceLine  `json:"trial_balance"`
	BankJournal         []FinanceBankJournalLine   `json:"bank_journal"`
	GeneralLedger       []FinanceGeneralLedgerLine `json:"general_ledger"`
	OpeningBalanceSheet FinanceOpeningBalanceSheet `json:"opening_balance_sheet"`
	EvidenceRef         string                     `json:"evidence_ref"`
}

type FinanceAccountingPeriod struct {
	PeriodCode string `json:"period_code"`
	StartsOn string `json:"starts_on"`
	EndsOn string `json:"ends_on"`
	Status string `json:"status"`
}

func (c *Client) IncorporationTrace(ctx context.Context, caseCode string) (IncorporationTrace, error) {
	var bundle struct {
		Verified bool               `json:"verified"`
		Trace    IncorporationTrace `json:"trace"`
	}
	err := c.request(ctx, "GET", "api/v1/incorporations/"+url.PathEscape(caseCode)+"/evidence", nil, &bundle)
	bundle.Trace.Verified = bundle.Verified
	return bundle.Trace, err
}

func (c *Client) IncorporationWorkItems(ctx context.Context, caseCode string) ([]IncorporationWorkItem, error) {
	var out struct {
		Items []IncorporationWorkItem `json:"items"`
	}
	err := c.request(ctx, "GET", "api/v1/incorporations/"+url.PathEscape(caseCode)+"/work-items", nil, &out)
	return out.Items, err
}

func (c *Client) FinanceOpening(ctx context.Context, caseCode string) (FinanceOpening, error) {
	var out FinanceOpening
	err := c.request(ctx, "GET", "api/v1/finance/opening/"+url.PathEscape(caseCode), nil, &out)
	return out, err
}
