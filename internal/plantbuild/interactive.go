package plantbuild

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"strings"
	"time"
)

const (
	InteractiveSchemaVersion = "1.0"
	PlanningPromptVersion    = "plant-planning-v1"
)

var ErrPlanningModelNotConfigured = errors.New("external planning model not configured")

type FacilityRequirement struct {
	SchemaVersion       string              `json:"schema_version"`
	RequirementID       string              `json:"requirement_id"`
	TenantID            string              `json:"tenant_id"`
	CaseCode            string              `json:"case_code"`
	LegalEntityCode     string              `json:"legal_entity_code"`
	TargetRegion        string              `json:"target_region"`
	FacilityPurpose     string              `json:"facility_purpose"`
	MinimumAreaM2       int                 `json:"minimum_area_m2"`
	MinimumElectricKVA  int                 `json:"minimum_electricity_kva"`
	TargetAvailableAt   string              `json:"target_available_at"`
	CandidateCount      int                 `json:"candidate_count"`
	AllowedOptionTypes  []string            `json:"allowed_option_types"`
	InvestmentRequest   Money               `json:"investment_request"`
	MinimumCashReserve  Money               `json:"minimum_cash_reserve"`
	FinancialConstraint FinancialConstraint `json:"financial_constraint"`
	Preferences         []string            `json:"preferences"`
	Revision            int                 `json:"revision"`
	RevisionReason      string              `json:"revision_reason"`
}

type FinancialConstraint struct {
	AvailableCash   Money  `json:"available_cash"`
	ApprovedBudget  Money  `json:"approved_budget"`
	CashSourceRef   string `json:"cash_source_ref"`
	BudgetSourceRef string `json:"budget_source_ref"`
	SnapshotHash    string `json:"snapshot_hash"`
}

type AmountRange struct {
	Minimum Money  `json:"minimum"`
	Likely  Money  `json:"likely"`
	Maximum Money  `json:"maximum"`
	Basis   string `json:"basis"`
}

type ScheduleRange struct {
	Earliest string `json:"earliest"`
	Likely   string `json:"likely"`
	Latest   string `json:"latest"`
}

type SiteOptionProposal struct {
	ProposalID        string        `json:"proposal_id"`
	OptionType        string        `json:"option_type"`
	DisplayName       string        `json:"display_name"`
	BusinessRationale string        `json:"business_rationale"`
	EstimatedAmount   AmountRange   `json:"estimated_amount"`
	EstimatedSchedule ScheduleRange `json:"estimated_schedule"`
	Assumptions       []string      `json:"assumptions"`
	FactsRequired     []string      `json:"facts_required"`
	Risks             []string      `json:"risks"`
	SourceRefs        []string      `json:"source_refs"`
	Confidence        string        `json:"confidence"`
	Status            string        `json:"status"`
}

type ProposalEvidence struct {
	Provider       string         `json:"provider"`
	Model          string         `json:"model"`
	PromptVersion  string         `json:"prompt_version"`
	SourceType     string         `json:"source_type,omitempty"`
	ParentRevision int            `json:"parent_revision"`
	RequestID      string         `json:"request_id,omitempty"`
	InputHash      string         `json:"input_hash"`
	OutputHash     string         `json:"output_hash"`
	TokenUsage     map[string]int `json:"token_usage,omitempty"`
	ValidatedAt    string         `json:"validated_at"`
}

type ProposalSet struct {
	SchemaVersion string               `json:"schema_version"`
	ProposalSetID string               `json:"proposal_set_id"`
	RequirementID string               `json:"requirement_id"`
	Revision      int                  `json:"revision"`
	Status        string               `json:"status"`
	Proposals     []SiteOptionProposal `json:"proposals"`
	Evidence      ProposalEvidence     `json:"evidence"`
}

type AgentRunEvidence struct {
	AgentRunID       string         `json:"agent_run_id"`
	CaseCode         string         `json:"case_code"`
	AgentID          string         `json:"agent_id"`
	Status           string         `json:"status"`
	Provider         string         `json:"provider"`
	Model            string         `json:"model"`
	ModelVersion     string         `json:"model_version"`
	PromptVersion    string         `json:"prompt_version"`
	RequestID        string         `json:"request_id"`
	InputHash        string         `json:"input_hash"`
	OutputHash       string         `json:"output_hash"`
	TokenUsage       map[string]int `json:"token_usage"`
	ValidationResult string         `json:"validation_result"`
	LatencyMS        int64          `json:"latency_ms"`
	StartedAt        string         `json:"started_at"`
	CompletedAt      string         `json:"completed_at"`
}

type ProposalReview struct {
	ProposalSetID    string `json:"proposal_set_id"`
	ProposalID       string `json:"proposal_id"`
	Action           string `json:"action"`
	Reason           string `json:"reason"`
	ReviewedBy       string `json:"reviewed_by"`
	ReviewedAt       string `json:"reviewed_at"`
	ExpectedRevision int    `json:"expected_revision"`
}

type InvestigationRequest struct {
	SchemaVersion          string   `json:"schema_version"`
	InvestigationRequestID string   `json:"investigation_request_id"`
	CaseCode               string   `json:"case_code"`
	ProposalSetID          string   `json:"proposal_set_id"`
	ProposalID             string   `json:"proposal_id"`
	ExpectedRevision       int      `json:"expected_revision"`
	WorldRunID             string   `json:"world_run_id"`
	Scope                  []string `json:"scope"`
	RequestedBy            string   `json:"requested_by"`
	RequestedAt            string   `json:"requested_at"`
	Status                 string   `json:"status"`
}

type InvestigationObservation struct {
	SchemaVersion          string   `json:"schema_version"`
	ObservationID          string   `json:"observation_id"`
	InvestigationRequestID string   `json:"investigation_request_id"`
	ProposalID             string   `json:"proposal_id"`
	Result                 string   `json:"result"`
	OwnershipStatus        string   `json:"ownership_status"`
	AvailableAreaM2        int      `json:"available_area_m2"`
	ElectricityKVA         int      `json:"electricity_kva"`
	QuotedAmount           Money    `json:"quoted_amount"`
	AvailableAt            string   `json:"available_at"`
	PermitStatus           string   `json:"permit_status"`
	EvidenceRefs           []string `json:"evidence_refs"`
	Notes                  string   `json:"notes"`
	ExternalActorID        string   `json:"external_actor_id"`
	ObservedAt             string   `json:"observed_at"`
}

type PlanningProvider interface {
	Status() PlanningProviderStatus
	Generate(context.Context, FacilityRequirement) (ProposalSet, error)
}

type JSONCompleter interface {
	CompleteJSON(context.Context, string, string, float64, int) (string, string, map[string]int, error)
}

type AIPlanningProvider struct {
	Completer JSONCompleter
	Provider  string
	Model     string
	Now       func() time.Time
}

func (p AIPlanningProvider) Status() PlanningProviderStatus {
	state := "connected"
	if p.Completer == nil {
		state = "not_configured"
	}
	return PlanningProviderStatus{State: state, Provider: p.Provider, Model: p.Model, PromptVersion: PlanningPromptVersion}
}

func (p AIPlanningProvider) Generate(ctx context.Context, requirement FacilityRequirement) (ProposalSet, error) {
	if err := ValidateRequirement(requirement); err != nil {
		return ProposalSet{}, err
	}
	if p.Completer == nil {
		return ProposalSet{}, ErrPlanningModelNotConfigured
	}
	input, _ := json.Marshal(requirement)
	user := `你是制造企业设施规划 Agent。只返回严格 JSON：{"proposals":[{"option_type":"","display_name":"","business_rationale":"","estimated_amount":{"minimum":{"value":"0.00","currency":"CNY","scale":2},"likely":{"value":"0.00","currency":"CNY","scale":2},"maximum":{"value":"0.00","currency":"CNY","scale":2},"basis":""},"estimated_schedule":{"earliest":"RFC3339","likely":"RFC3339","latest":"RFC3339"},"assumptions":[""],"facts_required":[""],"risks":[""],"source_refs":["requirement:<id>"],"confidence":"0.00"}]}。候选数量必须等于 candidate_count，option_type 只能来自 allowed_option_types。估算必须明确依据；不得声称已取得报价、权属、许可或容量证明。不要 Markdown 或思考过程。prompt_version=` + PlanningPromptVersion + `\nrequirement=` + string(input)
	content, requestID, usage, err := p.Completer.CompleteJSON(ctx, "Return strict JSON only. Never invent verified external facts.", user, 0.6, 8192)
	if err != nil {
		return ProposalSet{}, err
	}
	var decoded struct {
		Proposals []SiteOptionProposal `json:"proposals"`
	}
	decoder := json.NewDecoder(strings.NewReader(content))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&decoded); err != nil {
		return ProposalSet{}, fmt.Errorf("invalid plant planning JSON: %w", err)
	}
	inputHash := CanonicalHash(requirement)
	for i := range decoded.Proposals {
		decoded.Proposals[i].ProposalID = fmt.Sprintf("site-%s-r%d-%02d", stableInteractiveCode(requirement.RequirementID), requirement.Revision, i+1)
		decoded.Proposals[i].Status = "proposed"
	}
	now := time.Now().UTC()
	if p.Now != nil {
		now = p.Now().UTC()
	}
	set := ProposalSet{SchemaVersion: InteractiveSchemaVersion, ProposalSetID: "proposal-set-" + strings.TrimPrefix(inputHash, "sha256:")[:20], RequirementID: requirement.RequirementID, Revision: 1, Status: "candidate_only", Proposals: decoded.Proposals}
	set.Evidence = ProposalEvidence{Provider: p.Provider, Model: p.Model, PromptVersion: PlanningPromptVersion, RequestID: requestID, InputHash: inputHash, OutputHash: CanonicalHash(set.Proposals), TokenUsage: usage, ValidatedAt: now.Format(time.RFC3339)}
	if err := ValidateProposalSet(requirement, set); err != nil {
		return ProposalSet{}, err
	}
	return set, nil
}

func stableInteractiveCode(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = strings.NewReplacer(" ", "-", "_", "-", "/", "-").Replace(value)
	if value == "" {
		return "unnamed"
	}
	return value
}

type PlanningProviderStatus struct {
	State         string `json:"state"`
	Provider      string `json:"provider"`
	Model         string `json:"model,omitempty"`
	PromptVersion string `json:"prompt_version"`
}

type UnconfiguredPlanningProvider struct{}

func (UnconfiguredPlanningProvider) Status() PlanningProviderStatus {
	return PlanningProviderStatus{State: "not_configured", Provider: "none", PromptVersion: PlanningPromptVersion}
}
func (UnconfiguredPlanningProvider) Generate(context.Context, FacilityRequirement) (ProposalSet, error) {
	return ProposalSet{}, ErrPlanningModelNotConfigured
}

func ValidateRequirement(v FacilityRequirement) error {
	if v.SchemaVersion != InteractiveSchemaVersion || strings.TrimSpace(v.RequirementID) == "" || strings.TrimSpace(v.TenantID) == "" || strings.TrimSpace(v.CaseCode) == "" || strings.TrimSpace(v.LegalEntityCode) == "" {
		return errors.New("requirement identity is incomplete")
	}
	if strings.TrimSpace(v.TargetRegion) == "" || strings.TrimSpace(v.FacilityPurpose) == "" || v.MinimumAreaM2 <= 0 || v.MinimumElectricKVA <= 0 {
		return errors.New("facility requirement is incomplete")
	}
	if _, err := time.Parse(time.RFC3339, v.TargetAvailableAt); err != nil {
		return fmt.Errorf("target_available_at must be RFC3339: %w", err)
	}
	if v.CandidateCount < 2 || v.CandidateCount > 8 {
		return errors.New("candidate_count must be between 2 and 8")
	}
	if len(v.AllowedOptionTypes) == 0 || v.Revision < 1 || strings.TrimSpace(v.RevisionReason) == "" {
		return errors.New("option types and revision reason are required")
	}
	for name, amount := range map[string]Money{"investment_request": v.InvestmentRequest, "minimum_cash_reserve": v.MinimumCashReserve, "available_cash": v.FinancialConstraint.AvailableCash, "approved_budget": v.FinancialConstraint.ApprovedBudget} {
		if err := validateMoney(amount); err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
	}
	if v.InvestmentRequest.Currency != v.FinancialConstraint.AvailableCash.Currency || v.InvestmentRequest.Currency != v.FinancialConstraint.ApprovedBudget.Currency || v.MinimumCashReserve.Currency != v.InvestmentRequest.Currency {
		return errors.New("all financial constraints must use the same currency")
	}
	if strings.TrimSpace(v.FinancialConstraint.CashSourceRef) == "" || strings.TrimSpace(v.FinancialConstraint.BudgetSourceRef) == "" || !strings.HasPrefix(v.FinancialConstraint.SnapshotHash, "sha256:") {
		return errors.New("authoritative financial source refs and snapshot hash are required")
	}
	return nil
}

func ValidateProposalSet(requirement FacilityRequirement, set ProposalSet) error {
	if err := ValidateRequirement(requirement); err != nil {
		return err
	}
	if set.SchemaVersion != InteractiveSchemaVersion || set.RequirementID != requirement.RequirementID || set.Revision < 1 || set.Status != "candidate_only" {
		return errors.New("proposal set identity or status is invalid")
	}
	manualRevision := set.Evidence.Provider == "human"
	if (!manualRevision && len(set.Proposals) != requirement.CandidateCount) ||
		(manualRevision && (len(set.Proposals) < 1 || len(set.Proposals) > 12)) {
		return fmt.Errorf("proposal count %d does not match requested %d", len(set.Proposals), requirement.CandidateCount)
	}
	seen := map[string]bool{}
	allowed := map[string]bool{}
	for _, value := range requirement.AllowedOptionTypes {
		allowed[value] = true
	}
	for i, proposal := range set.Proposals {
		if strings.TrimSpace(proposal.ProposalID) == "" || strings.TrimSpace(proposal.DisplayName) == "" || strings.TrimSpace(proposal.BusinessRationale) == "" || !allowed[proposal.OptionType] || proposal.Status != "proposed" {
			return fmt.Errorf("proposal %d is incomplete or uses a forbidden option type", i+1)
		}
		if seen[proposal.ProposalID] || seen[proposal.DisplayName] {
			return errors.New("duplicate proposal identity")
		}
		seen[proposal.ProposalID], seen[proposal.DisplayName] = true, true
		if len(proposal.Assumptions) == 0 || len(proposal.FactsRequired) == 0 || len(proposal.Risks) == 0 || len(proposal.SourceRefs) == 0 {
			return fmt.Errorf("proposal %s must expose assumptions, unknown facts, risks, and sources", proposal.ProposalID)
		}
		if _, ok := new(big.Rat).SetString(proposal.Confidence); !ok {
			return fmt.Errorf("proposal %s confidence is invalid", proposal.ProposalID)
		}
		if err := validateAmountRange(proposal.EstimatedAmount); err != nil {
			return fmt.Errorf("proposal %s: %w", proposal.ProposalID, err)
		}
		if err := validateScheduleRange(proposal.EstimatedSchedule); err != nil {
			return fmt.Errorf("proposal %s: %w", proposal.ProposalID, err)
		}
	}
	if set.Evidence.Provider == "" || set.Evidence.Model == "" || !strings.HasPrefix(set.Evidence.InputHash, "sha256:") || !strings.HasPrefix(set.Evidence.OutputHash, "sha256:") {
		return errors.New("agent generation evidence is incomplete")
	}
	if manualRevision {
		if set.Revision < 1 || set.Evidence.Model != "manual-entry" || set.Evidence.PromptVersion != "manual-candidate-v1" ||
			set.Evidence.SourceType != "human_manual" || set.Evidence.ParentRevision != set.Revision-1 {
			return errors.New("manual proposal revision evidence is incomplete")
		}
	} else if set.Evidence.PromptVersion != PlanningPromptVersion {
		return errors.New("agent generation evidence is incomplete")
	}
	return nil
}

func ValidateReview(v ProposalReview) error {
	allowed := map[string]bool{"adopt_for_investigation": true, "request_revision": true, "add_manual_option": true, "discard": true}
	if v.ProposalSetID == "" || v.ProposalID == "" || !allowed[v.Action] || len(strings.TrimSpace(v.Reason)) < 6 || v.ReviewedBy == "" || v.ExpectedRevision < 1 {
		return errors.New("proposal review is incomplete")
	}
	if _, err := time.Parse(time.RFC3339, v.ReviewedAt); err != nil {
		return fmt.Errorf("reviewed_at must be RFC3339: %w", err)
	}
	return nil
}

func ValidateInvestigationRequest(v InvestigationRequest) error {
	if v.SchemaVersion != InteractiveSchemaVersion || v.InvestigationRequestID == "" || v.CaseCode == "" ||
		v.ProposalSetID == "" || v.ProposalID == "" || v.ExpectedRevision < 1 || v.WorldRunID == "" ||
		v.RequestedBy == "" || v.Status != "waiting_world" || len(v.Scope) == 0 {
		return errors.New("investigation request is incomplete")
	}
	if _, err := time.Parse(time.RFC3339, v.RequestedAt); err != nil {
		return fmt.Errorf("requested_at must be RFC3339: %w", err)
	}
	allowed := map[string]bool{"ownership": true, "commercial_quote": true, "available_area": true, "electricity_capacity": true, "available_date": true, "permit": true}
	seen := map[string]bool{}
	for _, item := range v.Scope {
		if !allowed[item] || seen[item] {
			return fmt.Errorf("unsupported or duplicated investigation scope: %s", item)
		}
		seen[item] = true
	}
	return nil
}

func ValidateInvestigationObservation(v InvestigationObservation) error {
	if v.SchemaVersion != InteractiveSchemaVersion || v.ObservationID == "" || v.InvestigationRequestID == "" ||
		v.ProposalID == "" || v.Result != "completed" || v.OwnershipStatus == "" || v.AvailableAreaM2 <= 0 ||
		v.ElectricityKVA <= 0 || v.PermitStatus == "" || len(v.EvidenceRefs) == 0 || v.ExternalActorID == "" {
		return errors.New("completed investigation observation is incomplete")
	}
	if err := validateMoney(v.QuotedAmount); err != nil {
		return fmt.Errorf("quoted_amount: %w", err)
	}
	if _, err := time.Parse(time.RFC3339, v.AvailableAt); err != nil {
		return fmt.Errorf("available_at must be RFC3339: %w", err)
	}
	if _, err := time.Parse(time.RFC3339, v.ObservedAt); err != nil {
		return fmt.Errorf("observed_at must be RFC3339: %w", err)
	}
	return nil
}

func validateMoney(v Money) error {
	if strings.TrimSpace(v.Currency) == "" || v.Scale < 0 || v.Scale > 6 {
		return errors.New("currency and supported scale are required")
	}
	n, ok := new(big.Rat).SetString(v.Value)
	if !ok || n.Sign() < 0 {
		return errors.New("amount must be a non-negative decimal string")
	}
	return nil
}

func validateAmountRange(v AmountRange) error {
	for _, amount := range []Money{v.Minimum, v.Likely, v.Maximum} {
		if err := validateMoney(amount); err != nil {
			return err
		}
	}
	if v.Minimum.Currency != v.Likely.Currency || v.Minimum.Currency != v.Maximum.Currency || strings.TrimSpace(v.Basis) == "" {
		return errors.New("amount range currency or basis is invalid")
	}
	min, _ := new(big.Rat).SetString(v.Minimum.Value)
	likely, _ := new(big.Rat).SetString(v.Likely.Value)
	max, _ := new(big.Rat).SetString(v.Maximum.Value)
	if min.Cmp(likely) > 0 || likely.Cmp(max) > 0 {
		return errors.New("amount range must be minimum <= likely <= maximum")
	}
	return nil
}

func validateScheduleRange(v ScheduleRange) error {
	values := make([]time.Time, 3)
	for i, raw := range []string{v.Earliest, v.Likely, v.Latest} {
		parsed, err := time.Parse(time.RFC3339, raw)
		if err != nil {
			return errors.New("schedule values must be RFC3339")
		}
		values[i] = parsed
	}
	if values[0].After(values[1]) || values[1].After(values[2]) {
		return errors.New("schedule must be earliest <= likely <= latest")
	}
	return nil
}

func CanonicalHash(value any) string {
	data, _ := json.Marshal(value)
	sum := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(sum[:])
}

func CanonicalizeProposalSet(set ProposalSet) ProposalSet {
	sort.Slice(set.Proposals, func(i, j int) bool { return set.Proposals[i].ProposalID < set.Proposals[j].ProposalID })
	return set
}
