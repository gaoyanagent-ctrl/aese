package iaosclient

import (
	"context"
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
	Verified       bool                         `json:"verified"`
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
