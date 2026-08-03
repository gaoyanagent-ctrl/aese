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
	PlanningPromptVersion    = "plant-planning-v2"
	RequirementPromptVersion = "plant-requirement-adviser-v1"
	ProjectPromptVersion     = "facility-project-wbs-v2"
	ContractPromptVersion    = "facility-contract-award-v1"
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

type RequirementOptionSeed struct {
	TenantID            string              `json:"tenant_id"`
	CaseCode            string              `json:"case_code"`
	LegalEntityCode     string              `json:"legal_entity_code"`
	FinancialConstraint FinancialConstraint `json:"financial_constraint"`
}

type RequirementOption struct {
	OptionID           string   `json:"option_id"`
	Title              string   `json:"title"`
	BusinessRationale  string   `json:"business_rationale"`
	TargetRegion       string   `json:"target_region"`
	FacilityPurpose    string   `json:"facility_purpose"`
	MinimumAreaM2      int      `json:"minimum_area_m2"`
	MinimumElectricKVA int      `json:"minimum_electricity_kva"`
	TargetAvailableAt  string   `json:"target_available_at"`
	CandidateCount     int      `json:"candidate_count"`
	AllowedOptionTypes []string `json:"allowed_option_types"`
	InvestmentRequest  Money    `json:"investment_request"`
	MinimumCashReserve Money    `json:"minimum_cash_reserve"`
	Preferences        []string `json:"preferences"`
	Tradeoffs          []string `json:"tradeoffs"`
}

type RequirementOptionSet struct {
	SchemaVersion string              `json:"schema_version"`
	Options       []RequirementOption `json:"options"`
	Evidence      ProposalEvidence    `json:"evidence"`
}

type ProjectPlanSeed struct {
	TenantID             string              `json:"tenant_id"`
	CaseCode             string              `json:"case_code"`
	SelectionID          string              `json:"selection_id"`
	ControlObservationID string              `json:"control_observation_id"`
	Requirement          FacilityRequirement `json:"facility_requirement"`
}

type ProjectWBSItem struct {
	WBSCode            string `json:"wbs_code"`
	Name               string `json:"name"`
	Phase              string `json:"phase"`
	Sequence           int    `json:"sequence"`
	OwnerPosition      string `json:"owner_position"`
	PlannedStartAt     string `json:"planned_start_at"`
	PlannedFinishAt    string `json:"planned_finish_at"`
	BudgetShareBPS     int    `json:"budget_share_bps"`
	AcceptanceCriteria string `json:"acceptance_criteria"`
}

type ProjectPlanOption struct {
	OptionID          string           `json:"option_id"`
	Title             string           `json:"title"`
	BusinessRationale string           `json:"business_rationale"`
	ProjectName       string           `json:"project_name"`
	DeliveryStrategy  string           `json:"delivery_strategy"`
	BudgetCeiling     Money            `json:"budget_ceiling"`
	TargetStartAt     string           `json:"target_start_at"`
	TargetReadyAt     string           `json:"target_ready_at"`
	WBSItems          []ProjectWBSItem `json:"wbs_items"`
	Tradeoffs         []string         `json:"tradeoffs"`
}

type ProjectPlanOptionSet struct {
	SchemaVersion string              `json:"schema_version"`
	Options       []ProjectPlanOption `json:"options"`
	Evidence      ProposalEvidence    `json:"evidence"`
}

type ContractRFQ struct {
	SchemaVersion    string `json:"schema_version"`
	RFQID            string `json:"rfq_id"`
	CaseCode         string `json:"case_code"`
	ProjectID        string `json:"project_id"`
	PackageCode      string `json:"package_code"`
	PackageName      string `json:"package_name"`
	SourcingStrategy string `json:"sourcing_strategy"`
	BidCount         int    `json:"bid_count"`
	ContractCeiling  Money  `json:"contract_ceiling"`
	RequiredReadyAt  string `json:"required_ready_at"`
	WorldRunID       string `json:"world_run_id"`
	RequestedBy      string `json:"requested_by"`
	RequestedAt      string `json:"requested_at"`
	Status           string `json:"status"`
}

type ContractBid struct {
	SchemaVersion   string   `json:"schema_version"`
	BidID           string   `json:"bid_id"`
	RFQID           string   `json:"rfq_id"`
	ContractorCode  string   `json:"contractor_code"`
	ContractorName  string   `json:"contractor_name"`
	QuotedAmount    Money    `json:"quoted_amount"`
	PromisedReadyAt string   `json:"promised_ready_at"`
	Qualification   string   `json:"qualification"`
	WarrantyMonths  int      `json:"warranty_months"`
	MilestoneCount  int      `json:"milestone_count"`
	EvidenceRefs    []string `json:"evidence_refs"`
	ObservedAt      string   `json:"observed_at"`
}

type ContractBidObservation struct {
	SchemaVersion   string        `json:"schema_version"`
	ObservationID   string        `json:"observation_id"`
	RFQID           string        `json:"rfq_id"`
	ExternalActorID string        `json:"external_actor_id"`
	Bids            []ContractBid `json:"bids"`
	ObservedAt      string        `json:"observed_at"`
}

type ContractRecommendationSeed struct {
	RFQ         ContractRFQ            `json:"rfq"`
	Observation ContractBidObservation `json:"observation"`
}
type ContractRecommendationAdvice struct {
	SelectedBidID         string           `json:"selected_bid_id"`
	RecommendationReason  string           `json:"recommendation_reason"`
	AlternativeComparison string           `json:"alternative_comparison"`
	Evidence              ProposalEvidence `json:"evidence"`
}

type FacilityProject struct {
	ProjectID     string           `json:"project_id"`
	CaseCode      string           `json:"case_code"`
	ProjectName   string           `json:"project_name"`
	BudgetCeiling Money            `json:"budget_ceiling"`
	TargetReadyAt string           `json:"target_ready_at"`
	WBSItems      []ProjectWBSItem `json:"wbs_items"`
	Status        string           `json:"status"`
}

type FacilityProjectItem struct {
	Project FacilityProject `json:"project"`
}

type ContractAwardItem struct {
	RFQ               ContractRFQ            `json:"rfq"`
	Status            string                 `json:"status"`
	Observation       ContractBidObservation `json:"observation"`
	Recommendation    map[string]any         `json:"recommendation"`
	ApprovalRequestID string                 `json:"approval_request_id"`
	ApprovalStatus    string                 `json:"approval_status"`
	Contract          map[string]any         `json:"contract"`
}

type ConstructionExecution struct {
	SchemaVersion  string `json:"schema_version"`
	ExecutionID    string `json:"execution_id"`
	CaseCode       string `json:"case_code"`
	ProjectID      string `json:"project_id"`
	ContractID     string `json:"contract_id"`
	PackageCode    string `json:"package_code"`
	PackageName    string `json:"package_name"`
	ContractorCode string `json:"contractor_code"`
	ContractorName string `json:"contractor_name"`
	WorldRunID     string `json:"world_run_id"`
	StartedBy      string `json:"started_by"`
	StartedAt      string `json:"started_at"`
	Status         string `json:"status"`
}

type ConstructionProgressObservation struct {
	SchemaVersion   string   `json:"schema_version"`
	ObservationID   string   `json:"observation_id"`
	ExecutionID     string   `json:"execution_id"`
	ContractID      string   `json:"contract_id"`
	PackageCode     string   `json:"package_code"`
	Result          string   `json:"result"`
	ProgressBPS     int      `json:"progress_bps"`
	QualityStatus   string   `json:"quality_status"`
	SafetyStatus    string   `json:"safety_status"`
	PunchItems      []string `json:"punch_items"`
	EvidenceRefs    []string `json:"evidence_refs"`
	ExternalActorID string   `json:"external_actor_id"`
	ObservedAt      string   `json:"observed_at"`
}

type ConstructionMilestoneItem struct {
	Execution   ConstructionExecution           `json:"execution"`
	Status      string                          `json:"status"`
	Observation ConstructionProgressObservation `json:"observation"`
	Acceptance  map[string]any                  `json:"acceptance"`
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

// InvestigationConfirmation is the player's minimal acknowledgement that the
// park investigation team may publish its report. All external facts are
// generated server-side from the authoritative request, requirement and Agent
// proposal; the browser cannot supply quote, capacity or evidence values.
type InvestigationConfirmation struct {
	SchemaVersion          string `json:"schema_version"`
	CaseCode               string `json:"case_code"`
	RequirementID          string `json:"requirement_id"`
	InvestigationRequestID string `json:"investigation_request_id"`
	Action                 string `json:"action"`
}

type SiteControlRequest struct {
	SchemaVersion      string   `json:"schema_version"`
	ControlRequestID   string   `json:"control_request_id"`
	SelectionID        string   `json:"selection_id"`
	CaseCode           string   `json:"case_code"`
	SelectedProposalID string   `json:"selected_proposal_id"`
	WorldRunID         string   `json:"world_run_id"`
	AgreementMode      string   `json:"agreement_mode"`
	RequestedHandover  string   `json:"requested_handover_at"`
	RequiredEvidence   []string `json:"required_evidence"`
	RequestedBy        string   `json:"requested_by"`
	RequestedAt        string   `json:"requested_at"`
	Status             string   `json:"status"`
}

type SiteControlObservation struct {
	SchemaVersion    string   `json:"schema_version"`
	ObservationID    string   `json:"observation_id"`
	ControlRequestID string   `json:"control_request_id"`
	SelectionID      string   `json:"selection_id"`
	Result           string   `json:"result"`
	AgreementRef     string   `json:"agreement_ref"`
	HandoverRef      string   `json:"handover_ref"`
	EffectiveAt      string   `json:"effective_at,omitempty"`
	EvidenceRefs     []string `json:"evidence_refs"`
	Notes            string   `json:"notes"`
	ExternalActorID  string   `json:"external_actor_id"`
	ObservedAt       string   `json:"observed_at"`
}

// SiteControlConfirmation is the only input a player supplies when the World
// rights holder is ready to hand over a selected site. Business evidence is
// deliberately absent: the World engine derives it from the authoritative
// SiteControlRequest so a browser cannot invent agreement or handover facts.
type SiteControlConfirmation struct {
	SchemaVersion    string `json:"schema_version"`
	CaseCode         string `json:"case_code"`
	ControlRequestID string `json:"control_request_id"`
	Action           string `json:"action"`
}

type SiteControlItem struct {
	Request     SiteControlRequest      `json:"request"`
	Status      string                  `json:"status"`
	Observation *SiteControlObservation `json:"observation,omitempty"`
}

type PlanningProvider interface {
	Status() PlanningProviderStatus
	Generate(context.Context, FacilityRequirement) (ProposalSet, error)
}

type RequirementAdviser interface {
	GenerateRequirementOptions(context.Context, RequirementOptionSeed) (RequirementOptionSet, error)
}

type ProjectBaselinePlanner interface {
	GenerateProjectPlanOptions(context.Context, ProjectPlanSeed) (ProjectPlanOptionSet, error)
}

type ContractAwardAdviser interface {
	GenerateContractRecommendation(context.Context, ContractRecommendationSeed) (ContractRecommendationAdvice, error)
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

func (p AIPlanningProvider) GenerateRequirementOptions(ctx context.Context, seed RequirementOptionSeed) (RequirementOptionSet, error) {
	if p.Completer == nil {
		return RequirementOptionSet{}, ErrPlanningModelNotConfigured
	}
	if strings.TrimSpace(seed.TenantID) == "" || strings.TrimSpace(seed.CaseCode) == "" || strings.TrimSpace(seed.LegalEntityCode) == "" ||
		strings.TrimSpace(seed.FinancialConstraint.SnapshotHash) == "" {
		return RequirementOptionSet{}, errors.New("requirement adviser seed is incomplete")
	}
	input, _ := json.Marshal(seed)
	user := `你是制造企业设施需求顾问。根据企业身份和 IAOS 权威资金/预算边界，生成 3 个明显不同、但都可执行的设施需求草案。只返回严格 JSON：{"options":[{"title":"","business_rationale":"","target_region":"","facility_purpose":"","minimum_area_m2":1,"minimum_electricity_kva":1,"target_available_at":"RFC3339","candidate_count":3,"allowed_option_types":["lease_and_retrofit"],"investment_request":{"value":"0.00","currency":"CNY","scale":2},"minimum_cash_reserve":{"value":"0.00","currency":"CNY","scale":2},"preferences":[""],"tradeoffs":[""]}]}。方案类型只能来自 lease_and_retrofit、greenfield_build、build_to_suit、existing_plant_purchase。金额必须根据本企业 available_cash 和 approved_budget 推导且可由用户调整；investment_request 不得超过 approved_budget，也不得超过 available_cash 减 minimum_cash_reserve。不得声称已取得场地、报价、权属或许可。不要 Markdown。prompt_version=` + RequirementPromptVersion + `\nseed=` + string(input)
	content, requestID, usage, err := p.Completer.CompleteJSON(ctx, "Return strict JSON only. Use authority financial limits and never invent verified external facts.", user, 0.55, 8192)
	if err != nil {
		return RequirementOptionSet{}, err
	}
	var decoded struct {
		Options []RequirementOption `json:"options"`
	}
	if err := json.Unmarshal([]byte(content), &decoded); err != nil {
		return RequirementOptionSet{}, fmt.Errorf("decode requirement options: %w", err)
	}
	if len(decoded.Options) < 2 || len(decoded.Options) > 3 {
		return RequirementOptionSet{}, errors.New("requirement adviser must return two or three options")
	}
	for index := range decoded.Options {
		decoded.Options[index].OptionID = fmt.Sprintf("requirement-option-%d", index+1)
		candidate := FacilityRequirement{
			SchemaVersion: InteractiveSchemaVersion, RequirementID: "adviser-validation", TenantID: seed.TenantID,
			CaseCode: seed.CaseCode, LegalEntityCode: seed.LegalEntityCode, TargetRegion: decoded.Options[index].TargetRegion,
			FacilityPurpose: decoded.Options[index].FacilityPurpose, MinimumAreaM2: decoded.Options[index].MinimumAreaM2,
			MinimumElectricKVA: decoded.Options[index].MinimumElectricKVA, TargetAvailableAt: decoded.Options[index].TargetAvailableAt,
			CandidateCount: decoded.Options[index].CandidateCount, AllowedOptionTypes: decoded.Options[index].AllowedOptionTypes,
			InvestmentRequest: decoded.Options[index].InvestmentRequest, MinimumCashReserve: decoded.Options[index].MinimumCashReserve,
			FinancialConstraint: seed.FinancialConstraint, Preferences: decoded.Options[index].Preferences, Revision: 1, RevisionReason: "Agent 需求草案",
		}
		if strings.TrimSpace(decoded.Options[index].Title) == "" || strings.TrimSpace(decoded.Options[index].BusinessRationale) == "" || len(decoded.Options[index].Tradeoffs) == 0 {
			return RequirementOptionSet{}, fmt.Errorf("requirement option %d lacks explanation", index+1)
		}
		if err := ValidateRequirement(candidate); err != nil {
			return RequirementOptionSet{}, fmt.Errorf("requirement option %d: %w", index+1, err)
		}
		investment, _ := new(big.Rat).SetString(candidate.InvestmentRequest.Value)
		budget, _ := new(big.Rat).SetString(seed.FinancialConstraint.ApprovedBudget.Value)
		cash, _ := new(big.Rat).SetString(seed.FinancialConstraint.AvailableCash.Value)
		reserve, _ := new(big.Rat).SetString(candidate.MinimumCashReserve.Value)
		if investment.Cmp(budget) > 0 || investment.Cmp(new(big.Rat).Sub(cash, reserve)) > 0 {
			return RequirementOptionSet{}, fmt.Errorf("requirement option %d exceeds authority financial limits", index+1)
		}
	}
	now := time.Now().UTC()
	if p.Now != nil {
		now = p.Now().UTC()
	}
	set := RequirementOptionSet{SchemaVersion: InteractiveSchemaVersion, Options: decoded.Options}
	set.Evidence = ProposalEvidence{Provider: p.Provider, Model: p.Model, PromptVersion: RequirementPromptVersion, RequestID: requestID, InputHash: CanonicalHash(seed), OutputHash: CanonicalHash(decoded.Options), TokenUsage: usage, ValidatedAt: now.Format(time.RFC3339)}
	return set, nil
}

func (p AIPlanningProvider) GenerateProjectPlanOptions(ctx context.Context, seed ProjectPlanSeed) (ProjectPlanOptionSet, error) {
	if p.Completer == nil {
		return ProjectPlanOptionSet{}, ErrPlanningModelNotConfigured
	}
	if seed.TenantID == "" || seed.CaseCode == "" || seed.SelectionID == "" || seed.ControlObservationID == "" {
		return ProjectPlanOptionSet{}, errors.New("facility project planning seed is incomplete")
	}
	input, _ := json.Marshal(seed)
	schema := `{"options":[{"title":"","business_rationale":"","project_name":"","delivery_strategy":"design_build","budget_ceiling":{"value":"0.00","currency":"CNY","scale":2},"target_start_at":"RFC3339","target_ready_at":"RFC3339","wbs_items":[{"wbs_code":"WBS-01","name":"","phase":"design","sequence":1,"owner_position":"plant-project-lead","planned_start_at":"RFC3339","planned_finish_at":"RFC3339","budget_share_bps":2500,"acceptance_criteria":""}],"tradeoffs":[""]}]}`
	user := `你是制造企业设施项目总师。基于已交付场址和设施需求，生成 2–3 个不同的项目/WBS 管理方案，每个方案优先使用 4–6 个精简工作包。只返回严格 JSON，形状为：` + schema + `。delivery_strategy 只能是 design_bid_build、design_build、epcm。每个方案 4–12 个 WBS，phase 只能是 design、procurement、construction、commissioning，sequence 从 1 连续，budget_share_bps 合计 10000。文字字段应简洁，日期必须在项目开始和投产之间；budget_ceiling 不得超过 facility_requirement.investment_request。不要 Markdown。prompt_version=` + ProjectPromptVersion + "\nseed=" + string(input)
	var decoded struct {
		Options []ProjectPlanOption `json:"options"`
	}
	var requestID string
	var validationErr error
	totalUsage := map[string]int{}
	for attempt := 0; attempt < 2; attempt++ {
		attemptUser := user
		if attempt > 0 {
			attemptUser = `上一次设施项目方案输出被截断或未通过治理校验。重新返回完整严格 JSON，不要解释、Markdown 或思考过程。恰好生成 2 个方案，每个方案恰好 4 个 WBS，四个 phase 各一个，budget_share_bps 合计 10000，文字保持精简。错误摘要：` + validationErr.Error() + "\nJSON形状：" + schema + "\n权威输入：" + string(input)
		}
		content, currentRequestID, usage, err := p.Completer.CompleteJSON(ctx, "Return one complete strict JSON object only. Keep the WBS concise and inside authority limits.", attemptUser, 0.35, 8192)
		requestID = currentRequestID
		for key, value := range usage {
			totalUsage[key] += value
		}
		if err != nil {
			validationErr = err
			if attempt == 0 && (strings.Contains(err.Error(), "truncated") || strings.Contains(err.Error(), "empty completion")) {
				continue
			}
			return ProjectPlanOptionSet{}, fmt.Errorf("facility project agent generation failed: %w", err)
		}
		decoded.Options = nil
		decoder := json.NewDecoder(strings.NewReader(content))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&decoded); err != nil {
			validationErr = fmt.Errorf("decode facility project options: %w", err)
			continue
		}
		if len(decoded.Options) < 2 || len(decoded.Options) > 3 {
			validationErr = errors.New("facility project planner must return two or three options")
			continue
		}
		validationErr = nil
		for index := range decoded.Options {
			decoded.Options[index].OptionID = fmt.Sprintf("project-option-%d", index+1)
			if err := ValidateProjectPlanOption(decoded.Options[index], seed.Requirement); err != nil {
				validationErr = fmt.Errorf("project option %d: %w", index+1, err)
				break
			}
		}
		if validationErr == nil {
			break
		}
	}
	if validationErr != nil {
		return ProjectPlanOptionSet{}, fmt.Errorf("facility project output invalid after one repair: %w", validationErr)
	}
	now := time.Now().UTC()
	if p.Now != nil {
		now = p.Now().UTC()
	}
	return ProjectPlanOptionSet{SchemaVersion: InteractiveSchemaVersion, Options: decoded.Options,
		Evidence: ProposalEvidence{Provider: p.Provider, Model: p.Model, PromptVersion: ProjectPromptVersion,
			RequestID: requestID, InputHash: CanonicalHash(seed), OutputHash: CanonicalHash(decoded.Options), TokenUsage: totalUsage, ValidatedAt: now.Format(time.RFC3339)}}, nil
}

func (p AIPlanningProvider) GenerateContractRecommendation(ctx context.Context, seed ContractRecommendationSeed) (ContractRecommendationAdvice, error) {
	if p.Completer == nil {
		return ContractRecommendationAdvice{}, ErrPlanningModelNotConfigured
	}
	if seed.RFQ.RFQID == "" || len(seed.Observation.Bids) < 2 {
		return ContractRecommendationAdvice{}, errors.New("contract recommendation seed is incomplete")
	}
	input, _ := json.Marshal(seed)
	user := `你是制造企业工程采购评审 Agent。只比较输入中的可信投标，推荐一个 bid_id，并解释成本、交付、质保和履约权衡。只返回严格 JSON：{"selected_bid_id":"","recommendation_reason":"","alternative_comparison":""}。不得发明承包商、报价或证据，不得自行批准。prompt_version=` + ContractPromptVersion + "\nseed=" + string(input)
	content, requestID, usage, err := p.Completer.CompleteJSON(ctx, "Return strict JSON only. Use only trusted bids and never approve the recommendation.", user, 0.25, 2048)
	if err != nil {
		return ContractRecommendationAdvice{}, err
	}
	var out ContractRecommendationAdvice
	decoder := json.NewDecoder(strings.NewReader(content))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&out); err != nil {
		return out, fmt.Errorf("decode contract recommendation: %w", err)
	}
	found := false
	for _, bid := range seed.Observation.Bids {
		if bid.BidID == out.SelectedBidID {
			found = true
			break
		}
	}
	if !found || len(strings.TrimSpace(out.RecommendationReason)) < 6 || len(strings.TrimSpace(out.AlternativeComparison)) < 6 {
		return out, errors.New("contract recommendation is incomplete or selects an unknown bid")
	}
	now := time.Now().UTC()
	if p.Now != nil {
		now = p.Now().UTC()
	}
	out.Evidence = ProposalEvidence{Provider: p.Provider, Model: p.Model, PromptVersion: ContractPromptVersion, RequestID: requestID, InputHash: CanonicalHash(map[string]any{"rfq": seed.RFQ, "observation": seed.Observation}), OutputHash: CanonicalHash(map[string]any{"selected_bid_id": out.SelectedBidID, "reason": out.RecommendationReason, "alternatives": out.AlternativeComparison}), TokenUsage: usage, ValidatedAt: now.Format(time.RFC3339)}
	return out, nil
}

func ValidateProjectPlanOption(option ProjectPlanOption, requirement FacilityRequirement) error {
	strategies := map[string]bool{"design_bid_build": true, "design_build": true, "epcm": true}
	if option.Title == "" || option.BusinessRationale == "" || option.ProjectName == "" || !strategies[option.DeliveryStrategy] || len(option.Tradeoffs) == 0 {
		return errors.New("project option purpose, strategy or tradeoff is incomplete")
	}
	budget, ok := new(big.Rat).SetString(option.BudgetCeiling.Value)
	if !ok || budget.Sign() <= 0 {
		return errors.New("project budget is invalid")
	}
	limit, ok := new(big.Rat).SetString(requirement.InvestmentRequest.Value)
	if !ok || budget.Cmp(limit) > 0 || option.BudgetCeiling.Currency != requirement.InvestmentRequest.Currency || option.BudgetCeiling.Scale != requirement.InvestmentRequest.Scale {
		return errors.New("project budget exceeds facility investment authority")
	}
	start, err1 := time.Parse(time.RFC3339, option.TargetStartAt)
	ready, err2 := time.Parse(time.RFC3339, option.TargetReadyAt)
	if err1 != nil || err2 != nil || !ready.After(start) {
		return errors.New("project dates are invalid")
	}
	if len(option.WBSItems) < 4 || len(option.WBSItems) > 12 {
		return errors.New("project WBS must contain 4–12 work packages")
	}
	phases := map[string]bool{"design": true, "procurement": true, "construction": true, "commissioning": true}
	share := 0
	seen := map[string]bool{}
	for index, item := range option.WBSItems {
		itemStart, e1 := time.Parse(time.RFC3339, item.PlannedStartAt)
		itemFinish, e2 := time.Parse(time.RFC3339, item.PlannedFinishAt)
		if item.WBSCode == "" || item.Name == "" || seen[item.WBSCode] || !phases[item.Phase] || item.Sequence != index+1 || item.OwnerPosition == "" || item.AcceptanceCriteria == "" || e1 != nil || e2 != nil || itemFinish.Before(itemStart) || itemStart.Before(start) || itemFinish.After(ready) || item.BudgetShareBPS <= 0 {
			return fmt.Errorf("WBS item %q is invalid", item.WBSCode)
		}
		seen[item.WBSCode] = true
		share += item.BudgetShareBPS
	}
	if share != 10000 {
		return fmt.Errorf("WBS budget shares total %d, want 10000", share)
	}
	return nil
}

func (p AIPlanningProvider) Generate(ctx context.Context, requirement FacilityRequirement) (ProposalSet, error) {
	if err := ValidateRequirement(requirement); err != nil {
		return ProposalSet{}, err
	}
	if p.Completer == nil {
		return ProposalSet{}, ErrPlanningModelNotConfigured
	}
	input, _ := json.Marshal(requirement)
	user := `你是制造企业设施规划 Agent。只返回严格 JSON：{"proposals":[{"option_type":"","display_name":"","business_rationale":"","estimated_amount":{"minimum":{"value":"0.00","currency":"CNY","scale":2},"likely":{"value":"0.00","currency":"CNY","scale":2},"maximum":{"value":"0.00","currency":"CNY","scale":2},"basis":""},"estimated_schedule":{"earliest":"RFC3339","likely":"RFC3339","latest":"RFC3339"},"assumptions":[""],"facts_required":[""],"risks":[""],"source_refs":["requirement:<id>"],"confidence":"0.00"}]}。候选数量必须等于 candidate_count，option_type 只能来自 allowed_option_types。每个候选的 minimum、likely、maximum 必须使用 investment_request 的 currency/scale，且 maximum 不得超过 investment_request.value；超出上限的方案不能作为候选返回，应改为上限内的可行方案并在 risks 说明范围取舍。估算必须明确依据；不得声称已取得报价、权属、许可或容量证明。不要 Markdown 或思考过程。prompt_version=` + PlanningPromptVersion + `\nrequirement=` + string(input)
	inputHash := CanonicalHash(requirement)
	now := time.Now().UTC()
	if p.Now != nil {
		now = p.Now().UTC()
	}
	totalUsage := map[string]int{}
	var previousContent string
	var validationErr error
	for attempt := 0; attempt < 2; attempt++ {
		attemptUser := user
		if attempt > 0 {
			attemptUser += "\n上一版输出未通过治理校验：" + validationErr.Error() + "。请重新生成完整 JSON，不要解释。\nprevious_invalid_output=" + previousContent
		}
		content, requestID, usage, err := p.Completer.CompleteJSON(ctx, "Return strict JSON only. Never invent verified external facts or exceed the investment ceiling.", attemptUser, 0.6, 8192)
		if err != nil {
			return ProposalSet{}, err
		}
		for key, value := range usage {
			totalUsage[key] += value
		}
		previousContent = content
		var decoded struct {
			Proposals []SiteOptionProposal `json:"proposals"`
		}
		decoder := json.NewDecoder(strings.NewReader(content))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&decoded); err != nil {
			validationErr = fmt.Errorf("invalid plant planning JSON: %w", err)
			continue
		}
		for i := range decoded.Proposals {
			decoded.Proposals[i].ProposalID = fmt.Sprintf("site-%s-r%d-%02d", stableInteractiveCode(requirement.RequirementID), requirement.Revision, i+1)
			decoded.Proposals[i].Status = "proposed"
		}
		set := ProposalSet{SchemaVersion: InteractiveSchemaVersion, ProposalSetID: "proposal-set-" + strings.TrimPrefix(inputHash, "sha256:")[:20], RequirementID: requirement.RequirementID, Revision: 1, Status: "candidate_only", Proposals: decoded.Proposals}
		set.Evidence = ProposalEvidence{Provider: p.Provider, Model: p.Model, PromptVersion: PlanningPromptVersion, RequestID: requestID, InputHash: inputHash, OutputHash: CanonicalHash(set.Proposals), TokenUsage: totalUsage, ValidatedAt: now.Format(time.RFC3339)}
		if err := ValidateProposalSet(requirement, set); err != nil {
			validationErr = err
			continue
		}
		return set, nil
	}
	return ProposalSet{}, fmt.Errorf("plant planning output invalid after one repair: %w", validationErr)
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
func (UnconfiguredPlanningProvider) GenerateRequirementOptions(context.Context, RequirementOptionSeed) (RequirementOptionSet, error) {
	return RequirementOptionSet{}, ErrPlanningModelNotConfigured
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
		if proposal.EstimatedAmount.Maximum.Currency != requirement.InvestmentRequest.Currency ||
			proposal.EstimatedAmount.Maximum.Scale != requirement.InvestmentRequest.Scale {
			return fmt.Errorf("proposal %s amount currency or scale differs from investment request", proposal.ProposalID)
		}
		maximum, _ := new(big.Rat).SetString(proposal.EstimatedAmount.Maximum.Value)
		ceiling, _ := new(big.Rat).SetString(requirement.InvestmentRequest.Value)
		if maximum.Cmp(ceiling) > 0 {
			return fmt.Errorf("proposal %s maximum amount exceeds investment request", proposal.ProposalID)
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

func ValidateInvestigationConfirmation(v InvestigationConfirmation) error {
	if v.SchemaVersion != InteractiveSchemaVersion || strings.TrimSpace(v.CaseCode) == "" ||
		strings.TrimSpace(v.RequirementID) == "" || strings.TrimSpace(v.InvestigationRequestID) == "" {
		return errors.New("investigation confirmation is incomplete")
	}
	if v.Action != "accept_report" {
		return errors.New("investigation confirmation action is invalid")
	}
	return nil
}

func ValidateSiteControlRequest(v SiteControlRequest) error {
	if v.SchemaVersion != InteractiveSchemaVersion || strings.TrimSpace(v.ControlRequestID) == "" ||
		strings.TrimSpace(v.SelectionID) == "" || strings.TrimSpace(v.CaseCode) == "" ||
		strings.TrimSpace(v.SelectedProposalID) == "" || strings.TrimSpace(v.WorldRunID) == "" ||
		strings.TrimSpace(v.RequestedBy) == "" || v.Status != "waiting_world" {
		return errors.New("site control request is incomplete")
	}
	allowedModes := map[string]bool{"lease": true, "purchase": true, "build_to_suit": true, "use_agreement": true}
	if !allowedModes[v.AgreementMode] {
		return errors.New("site control agreement mode is invalid")
	}
	if _, err := time.Parse(time.RFC3339, v.RequestedHandover); err != nil {
		return fmt.Errorf("requested_handover_at must be RFC3339: %w", err)
	}
	if _, err := time.Parse(time.RFC3339, v.RequestedAt); err != nil {
		return fmt.Errorf("requested_at must be RFC3339: %w", err)
	}
	required := map[string]bool{"executed_agreement": false, "handover_record": false, "possession_authority": false}
	for _, item := range v.RequiredEvidence {
		if _, ok := required[item]; !ok || required[item] {
			return errors.New("site control evidence scope is invalid")
		}
		required[item] = true
	}
	for _, present := range required {
		if !present {
			return errors.New("site control requires agreement, handover and possession evidence")
		}
	}
	return nil
}

func ValidateSiteControlObservation(v SiteControlObservation) error {
	if v.SchemaVersion != InteractiveSchemaVersion || strings.TrimSpace(v.ObservationID) == "" ||
		strings.TrimSpace(v.ControlRequestID) == "" || strings.TrimSpace(v.SelectionID) == "" ||
		strings.TrimSpace(v.ExternalActorID) == "" || len(v.EvidenceRefs) == 0 {
		return errors.New("site control observation is incomplete")
	}
	if v.Result != "delivered" && v.Result != "delayed" && v.Result != "rejected" {
		return errors.New("site control result is invalid")
	}
	if _, err := time.Parse(time.RFC3339, v.ObservedAt); err != nil {
		return fmt.Errorf("observed_at must be RFC3339: %w", err)
	}
	if v.Result == "delivered" {
		if strings.TrimSpace(v.AgreementRef) == "" || strings.TrimSpace(v.HandoverRef) == "" {
			return errors.New("delivered site control requires agreement and handover references")
		}
		if _, err := time.Parse(time.RFC3339, v.EffectiveAt); err != nil {
			return fmt.Errorf("effective_at must be RFC3339: %w", err)
		}
	}
	return nil
}

func ValidateSiteControlConfirmation(v SiteControlConfirmation) error {
	if v.SchemaVersion != InteractiveSchemaVersion || strings.TrimSpace(v.CaseCode) == "" ||
		strings.TrimSpace(v.ControlRequestID) == "" {
		return errors.New("site control confirmation is incomplete")
	}
	if v.Action != "accept_delivery" {
		return errors.New("site control confirmation action is invalid")
	}
	return nil
}

// GenerateSiteControlObservation turns an authoritative waiting request into a
// deterministic World fact. The request's simulation handover time is used as
// both occurrence and effective time, making retries and replay byte-stable.
func GenerateSiteControlObservation(request SiteControlRequest) (SiteControlObservation, error) {
	if err := ValidateSiteControlRequest(request); err != nil {
		return SiteControlObservation{}, err
	}
	suffix := strings.ToUpper(strings.TrimPrefix(CanonicalHash(request), "sha256:")[:12])
	mode := strings.ToUpper(strings.ReplaceAll(request.AgreementMode, "_", "-"))
	agreementRef := "agreement:" + mode + "-" + suffix
	handoverRef := "handover:HO-" + suffix
	observation := SiteControlObservation{
		SchemaVersion:    InteractiveSchemaVersion,
		ObservationID:    "site-control-observation-" + strings.ToLower(suffix),
		ControlRequestID: request.ControlRequestID,
		SelectionID:      request.SelectionID,
		Result:           "delivered",
		AgreementRef:     agreementRef,
		HandoverRef:      handoverRef,
		EffectiveAt:      request.RequestedHandover,
		EvidenceRefs: []string{
			"world-document:" + agreementRef,
			"world-document:" + handoverRef,
			"world-evidence:possession-authority:" + request.ControlRequestID,
		},
		Notes:           "园区权利方已完成协议签署、现场交接和占有权限移交；事实由 AESE World 交付策略生成。",
		ExternalActorID: "world-park-rights-holder",
		ObservedAt:      request.RequestedHandover,
	}
	if err := ValidateSiteControlObservation(observation); err != nil {
		return SiteControlObservation{}, err
	}
	return observation, nil
}

// GenerateContractBidObservation models the external contractor market. It
// consumes only the authoritative RFQ and deterministically creates fictional
// quotations below its ceiling. The browser can acknowledge receipt but cannot
// provide contractor names, prices, dates, qualifications or evidence.
func GenerateContractBidObservation(rfq ContractRFQ) (ContractBidObservation, error) {
	if rfq.SchemaVersion != InteractiveSchemaVersion || rfq.RFQID == "" || rfq.CaseCode == "" ||
		rfq.BidCount < 2 || rfq.BidCount > 5 || rfq.ContractCeiling.Currency == "" {
		return ContractBidObservation{}, errors.New("contract RFQ is incomplete")
	}
	ceiling, ok := new(big.Rat).SetString(rfq.ContractCeiling.Value)
	if !ok || ceiling.Sign() <= 0 {
		return ContractBidObservation{}, errors.New("contract RFQ ceiling is invalid")
	}
	ready, err := time.Parse(time.RFC3339, rfq.RequiredReadyAt)
	if err != nil {
		return ContractBidObservation{}, errors.New("contract RFQ required_ready_at is invalid")
	}
	requested, err := time.Parse(time.RFC3339, rfq.RequestedAt)
	if err != nil {
		return ContractBidObservation{}, errors.New("contract RFQ requested_at is invalid")
	}
	hash := strings.TrimPrefix(CanonicalHash(rfq), "sha256:")
	suffix := strings.ToUpper(hash[:12])
	names := []struct{ code, name string }{
		{"WORLD-CONTRACTOR-01", "澄岳工程（虚构）"},
		{"WORLD-CONTRACTOR-02", "远拓建设（虚构）"},
		{"WORLD-CONTRACTOR-03", "岚川工业工程（虚构）"},
		{"WORLD-CONTRACTOR-04", "衡筑设施（虚构）"},
		{"WORLD-CONTRACTOR-05", "启辰项目服务（虚构）"},
	}
	ratios := []int64{92, 96, 89, 98, 94}
	bids := make([]ContractBid, 0, rfq.BidCount)
	for index := 0; index < rfq.BidCount; index++ {
		amount := new(big.Rat).Mul(ceiling, big.NewRat(ratios[(index+int(hash[0]))%len(ratios)], 100))
		promised := ready.AddDate(0, 0, (index%3)-1)
		observed := requested.Add(time.Duration(index+1) * time.Hour)
		bidID := fmt.Sprintf("bid-%s-%02d", strings.ToLower(suffix), index+1)
		bids = append(bids, ContractBid{
			SchemaVersion: InteractiveSchemaVersion, BidID: bidID, RFQID: rfq.RFQID,
			ContractorCode: names[index].code, ContractorName: names[index].name,
			QuotedAmount:    Money{Value: amount.FloatString(rfq.ContractCeiling.Scale), Currency: rfq.ContractCeiling.Currency, Scale: rfq.ContractCeiling.Scale},
			PromisedReadyAt: promised.Format(time.RFC3339), Qualification: "eligible",
			WarrantyMonths: 18 + index*6, MilestoneCount: 4 + index,
			EvidenceRefs: []string{"world-document:sealed-bid-" + bidID, "world-evidence:qualification-" + names[index].code},
			ObservedAt:   observed.Format(time.RFC3339),
		})
	}
	return ContractBidObservation{SchemaVersion: InteractiveSchemaVersion, ObservationID: "contract-bid-observation-" + strings.ToLower(suffix), RFQID: rfq.RFQID, ExternalActorID: "world-contractor-market", Bids: bids, ObservedAt: requested.Add(6 * time.Hour).Format(time.RFC3339)}, nil
}

// GenerateConstructionProgressObservation is the deterministic fictional
// construction-site actor. The browser only confirms advancing time; it never
// types progress, quality, safety or evidence values.
func GenerateConstructionProgressObservation(execution ConstructionExecution) (ConstructionProgressObservation, error) {
	if execution.SchemaVersion != InteractiveSchemaVersion || execution.ExecutionID == "" || execution.ContractID == "" || execution.PackageCode == "" || execution.Status != "waiting_world" {
		return ConstructionProgressObservation{}, errors.New("construction execution is incomplete or not waiting for World")
	}
	started, err := time.Parse(time.RFC3339, execution.StartedAt)
	if err != nil {
		return ConstructionProgressObservation{}, errors.New("construction started_at is invalid")
	}
	suffix := strings.ToLower(strings.TrimPrefix(CanonicalHash(execution), "sha256:")[:12])
	return ConstructionProgressObservation{SchemaVersion: InteractiveSchemaVersion, ObservationID: "construction-observation-" + suffix, ExecutionID: execution.ExecutionID, ContractID: execution.ContractID, PackageCode: execution.PackageCode, Result: "ready_for_inspection", ProgressBPS: 10000, QualityStatus: "passed", SafetyStatus: "clear", PunchItems: []string{}, EvidenceRefs: []string{"world-document:completion-report-" + suffix, "world-evidence:quality-inspection-" + suffix, "world-evidence:safety-clearance-" + suffix}, ExternalActorID: "world-construction-contractor", ObservedAt: started.AddDate(0, 1, 0).Format(time.RFC3339)}, nil
}

// GenerateInvestigationObservation produces a replay-stable external report
// from authority data. Values deliberately remain distinct from the Agent's
// estimates while satisfying the verified virtual park offer. This is a World
// fact generator, not a browser convenience default.
func GenerateInvestigationObservation(request InvestigationRequest, requirement FacilityRequirement, proposal SiteOptionProposal) (InvestigationObservation, error) {
	if err := ValidateInvestigationRequest(request); err != nil {
		return InvestigationObservation{}, err
	}
	if err := ValidateRequirement(requirement); err != nil {
		return InvestigationObservation{}, err
	}
	if request.CaseCode != requirement.CaseCode || request.ProposalID != proposal.ProposalID {
		return InvestigationObservation{}, errors.New("investigation authority references do not match")
	}
	hash := strings.TrimPrefix(CanonicalHash(struct {
		Request     InvestigationRequest
		Requirement FacilityRequirement
		Proposal    SiteOptionProposal
	}{request, requirement, proposal}), "sha256:")
	suffix := strings.ToUpper(hash[:12])
	seed := int(hash[0]) + int(hash[1])
	area := requirement.MinimumAreaM2 * (105 + seed%16) / 100
	electricity := requirement.MinimumElectricKVA * (105 + (seed/3)%16) / 100
	quoted := proposal.EstimatedAmount.Likely
	amount, ok := new(big.Rat).SetString(quoted.Value)
	if !ok {
		return InvestigationObservation{}, errors.New("proposal likely amount is invalid")
	}
	amount.Mul(amount, big.NewRat(int64(96+seed%5), 100))
	quoted.Value = amount.FloatString(quoted.Scale)
	availableAt := proposal.EstimatedSchedule.Likely
	available, err := time.Parse(time.RFC3339, availableAt)
	if err != nil {
		return InvestigationObservation{}, errors.New("proposal likely schedule is invalid")
	}
	target, _ := time.Parse(time.RFC3339, requirement.TargetAvailableAt)
	if available.After(target) {
		availableAt = target.Format(time.RFC3339)
	}
	observed, _ := time.Parse(time.RFC3339, request.RequestedAt)
	observation := InvestigationObservation{
		SchemaVersion: InteractiveSchemaVersion, ObservationID: "site-observation-" + strings.ToLower(suffix),
		InvestigationRequestID: request.InvestigationRequestID, ProposalID: request.ProposalID,
		Result: "completed", OwnershipStatus: "verified", AvailableAreaM2: area, ElectricityKVA: electricity,
		QuotedAmount: quoted, AvailableAt: availableAt, PermitStatus: "eligible",
		EvidenceRefs: []string{
			"world-document:site-title-report-" + suffix,
			"world-document:utility-capacity-report-" + suffix,
			"world-document:commercial-quote-" + suffix,
		},
		Notes:           "园区调研团队已核验权属、空间、公用工程、报价、交付日期和许可条件。",
		ExternalActorID: "world-park-investigation-team", ObservedAt: observed.Add(4 * time.Hour).Format(time.RFC3339),
	}
	if err := ValidateInvestigationObservation(observation); err != nil {
		return InvestigationObservation{}, err
	}
	return observation, nil
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
