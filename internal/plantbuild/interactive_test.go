package plantbuild

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"strings"
	"testing"
	"time"
)

func projectOptionFixture(title string) ProjectPlanOption {
	items := []ProjectWBSItem{
		{WBSCode: "WBS-01", Name: "工程设计", Phase: "design", Sequence: 1, OwnerPosition: "plant-project-lead", PlannedStartAt: "2026-09-01T00:00:00Z", PlannedFinishAt: "2026-09-30T00:00:00Z", BudgetShareBPS: 2000, AcceptanceCriteria: "设计包批准"},
		{WBSCode: "WBS-02", Name: "设备采购", Phase: "procurement", Sequence: 2, OwnerPosition: "procurement-lead", PlannedStartAt: "2026-10-01T00:00:00Z", PlannedFinishAt: "2026-11-30T00:00:00Z", BudgetShareBPS: 3000, AcceptanceCriteria: "设备到货"},
		{WBSCode: "WBS-03", Name: "现场施工", Phase: "construction", Sequence: 3, OwnerPosition: "construction-lead", PlannedStartAt: "2026-12-01T00:00:00Z", PlannedFinishAt: "2027-01-31T00:00:00Z", BudgetShareBPS: 3500, AcceptanceCriteria: "施工验收"},
		{WBSCode: "WBS-04", Name: "联调投产", Phase: "commissioning", Sequence: 4, OwnerPosition: "commissioning-lead", PlannedStartAt: "2027-02-01T00:00:00Z", PlannedFinishAt: "2027-02-28T00:00:00Z", BudgetShareBPS: 1500, AcceptanceCriteria: "联调通过"},
	}
	return ProjectPlanOption{Title: title, BusinessRationale: "在投资边界内形成可执行交付基线", ProjectName: title + "项目", DeliveryStrategy: "design_build", BudgetCeiling: Money{"18000000.00", "CNY", 2}, TargetStartAt: "2026-09-01T00:00:00Z", TargetReadyAt: "2027-02-28T00:00:00Z", WBSItems: items, Tradeoffs: []string{"周期与控制权平衡"}}
}

type planningCompleterStub struct{ content string }

func (s planningCompleterStub) CompleteJSON(context.Context, string, string, float64, int) (string, string, map[string]int, error) {
	return s.content, "request-1", map[string]int{"total_tokens": 42}, nil
}

type planningSequenceCompleterStub struct {
	contents  []string
	calls     int
	maxTokens []int
}

type planningTimeoutThenContentCompleterStub struct {
	content string
	calls   int
}

func (s *planningTimeoutThenContentCompleterStub) CompleteJSON(_ context.Context, _ string, _ string, _ float64, _ int) (string, string, map[string]int, error) {
	s.calls++
	if s.calls == 1 {
		return "", "request-timeout", nil, context.DeadlineExceeded
	}
	return s.content, "request-recovery", map[string]int{"total_tokens": 20}, nil
}

func (s *planningSequenceCompleterStub) CompleteJSON(_ context.Context, _ string, _ string, _ float64, maxTokens int) (string, string, map[string]int, error) {
	index := s.calls
	if index >= len(s.contents) {
		index = len(s.contents) - 1
	}
	s.calls++
	s.maxTokens = append(s.maxTokens, maxTokens)
	return s.contents[index], "request-" + string(rune('1'+index)), map[string]int{"total_tokens": 10 + index}, nil
}

func requirementFixture() FacilityRequirement {
	return FacilityRequirement{SchemaVersion: "1.0", RequirementID: "REQ-1", TenantID: "tenant-a", CaseCode: "INC-1", LegalEntityCode: "LE-1", TargetRegion: "华东", FacilityPurpose: "汽车零部件制造", MinimumAreaM2: 12000, MinimumElectricKVA: 2200, TargetAvailableAt: "2027-01-01T00:00:00+08:00", CandidateCount: 2, AllowedOptionTypes: []string{"leased_shell", "build_to_suit"}, InvestmentRequest: Money{"18000000.00", "CNY", 2}, MinimumCashReserve: Money{"5000000.00", "CNY", 2}, FinancialConstraint: FinancialConstraint{AvailableCash: Money{"30000000.00", "CNY", 2}, ApprovedBudget: Money{"20000000.00", "CNY", 2}, CashSourceRef: "ledger:CASH-1", BudgetSourceRef: "budget:BUD-1", SnapshotHash: "sha256:abc"}, Preferences: []string{"优先投产速度"}, Revision: 1, RevisionReason: "首次规划"}
}

func TestInteractiveRequirementRejectsHardcodedOrMalformedInputs(t *testing.T) {
	v := requirementFixture()
	if err := ValidateRequirement(v); err != nil {
		t.Fatal(err)
	}
	v.CandidateCount = 1
	if ValidateRequirement(v) == nil {
		t.Fatal("candidate count below two accepted")
	}
	v = requirementFixture()
	v.FinancialConstraint.CashSourceRef = ""
	if ValidateRequirement(v) == nil {
		t.Fatal("cash without authority source accepted")
	}
	v = requirementFixture()
	v.InvestmentRequest.Value = "not-money"
	if ValidateRequirement(v) == nil {
		t.Fatal("invalid decimal accepted")
	}
}

func TestUnconfiguredPlanningProviderFailsClosed(t *testing.T) {
	provider := UnconfiguredPlanningProvider{}
	if provider.Status().State != "not_configured" {
		t.Fatal("bad status")
	}
	_, err := provider.Generate(context.Background(), requirementFixture())
	if !errors.Is(err, ErrPlanningModelNotConfigured) {
		t.Fatalf("error=%v", err)
	}
}

func TestProposalReviewRequiresReasonAndHumanIdentity(t *testing.T) {
	valid := ProposalReview{ProposalSetID: "SET-1", ProposalID: "P-1", Action: "adopt_for_investigation", Reason: "符合一期投产要求", ReviewedBy: "project-lead", ReviewedAt: "2026-08-01T10:00:00Z", ExpectedRevision: 1}
	if err := ValidateReview(valid); err != nil {
		t.Fatal(err)
	}
	valid.Reason = "短"
	if ValidateReview(valid) == nil {
		t.Fatal("short reason accepted")
	}
}

func TestInvestigationContractsRequireTraceableWorldFacts(t *testing.T) {
	request := InvestigationRequest{SchemaVersion: "1.0", InvestigationRequestID: "INV-1", CaseCode: "INC-1", ProposalSetID: "SET-1", ProposalID: "P-1", ExpectedRevision: 1, WorldRunID: "world-1", Scope: []string{"ownership"}, RequestedBy: "project-lead", RequestedAt: "2026-08-01T10:00:00Z", Status: "waiting_world"}
	if err := ValidateInvestigationRequest(request); err != nil {
		t.Fatal(err)
	}
	observation := InvestigationObservation{SchemaVersion: "1.0", ObservationID: "OBS-1", InvestigationRequestID: "INV-1", ProposalID: "P-1", Result: "completed", OwnershipStatus: "verified", AvailableAreaM2: 9000, ElectricityKVA: 3000, QuotedAmount: Money{"9800000.00", "CNY", 2}, AvailableAt: "2026-10-01T00:00:00Z", PermitStatus: "eligible", EvidenceRefs: []string{"world-document:QUOTE-1"}, ExternalActorID: "park-operator", ObservedAt: "2026-08-01T11:00:00Z"}
	if err := ValidateInvestigationObservation(observation); err != nil {
		t.Fatal(err)
	}
	observation.EvidenceRefs = nil
	if ValidateInvestigationObservation(observation) == nil {
		t.Fatal("observation without evidence accepted")
	}
}

func TestSiteControlContractsRequireAgreementAndHandoverEvidence(t *testing.T) {
	request := SiteControlRequest{
		SchemaVersion: "1.0", ControlRequestID: "CTRL-1", SelectionID: "SEL-1",
		CaseCode: "INC-1", SelectedProposalID: "P-1", WorldRunID: "world-1",
		AgreementMode: "lease", RequestedHandover: "2026-10-01T00:00:00Z",
		RequiredEvidence: []string{"executed_agreement", "handover_record", "possession_authority"},
		RequestedBy:      "project-lead", RequestedAt: "2026-08-02T10:00:00Z", Status: "waiting_world",
	}
	if err := ValidateSiteControlRequest(request); err != nil {
		t.Fatal(err)
	}
	request.RequiredEvidence = []string{"executed_agreement", "handover_record"}
	if ValidateSiteControlRequest(request) == nil {
		t.Fatal("site control request without possession authority was accepted")
	}
	observation := SiteControlObservation{
		SchemaVersion: "1.0", ObservationID: "OBS-CTRL-1", ControlRequestID: "CTRL-1",
		SelectionID: "SEL-1", Result: "delivered", AgreementRef: "agreement:LEASE-1",
		HandoverRef: "handover:HANDOVER-1", EffectiveAt: "2026-10-01T00:00:00Z",
		EvidenceRefs:    []string{"world-document:LEASE-1", "world-document:HANDOVER-1"},
		ExternalActorID: "world-park-operator", ObservedAt: "2026-10-01T08:00:00Z",
	}
	if err := ValidateSiteControlObservation(observation); err != nil {
		t.Fatal(err)
	}
	observation.AgreementRef = ""
	if ValidateSiteControlObservation(observation) == nil {
		t.Fatal("delivered site control without agreement evidence was accepted")
	}
}

func TestWorldEngineGeneratesDeterministicSiteControlEvidence(t *testing.T) {
	request := SiteControlRequest{
		SchemaVersion: "1.0", ControlRequestID: "CTRL-1", SelectionID: "SEL-1",
		CaseCode: "INC-1", SelectedProposalID: "P-1", WorldRunID: "world-1",
		AgreementMode: "lease", RequestedHandover: "2026-10-01T00:00:00Z",
		RequiredEvidence: []string{"executed_agreement", "handover_record", "possession_authority"},
		RequestedBy:      "project-lead", RequestedAt: "2026-08-02T10:00:00Z", Status: "waiting_world",
	}
	first, err := GenerateSiteControlObservation(request)
	if err != nil {
		t.Fatal(err)
	}
	second, err := GenerateSiteControlObservation(request)
	if err != nil {
		t.Fatal(err)
	}
	if CanonicalHash(first) != CanonicalHash(second) {
		t.Fatalf("world delivery is not replayable: first=%+v second=%+v", first, second)
	}
	if first.EffectiveAt != request.RequestedHandover || first.ObservedAt != request.RequestedHandover ||
		!strings.HasPrefix(first.AgreementRef, "agreement:LEASE-") || !strings.HasPrefix(first.HandoverRef, "handover:HO-") {
		t.Fatalf("unexpected generated observation %+v", first)
	}
	if len(first.EvidenceRefs) != 3 || !strings.Contains(first.EvidenceRefs[2], request.ControlRequestID) {
		t.Fatalf("missing generated authority evidence %+v", first.EvidenceRefs)
	}
	if err := ValidateSiteControlConfirmation(SiteControlConfirmation{SchemaVersion: "1.0", CaseCode: "INC-1", ControlRequestID: "CTRL-1", Action: "accept_delivery"}); err != nil {
		t.Fatal(err)
	}
}

func TestWorldEngineGeneratesDeterministicInvestigationEvidence(t *testing.T) {
	requirement := requirementFixture()
	proposal := SiteOptionProposal{
		ProposalID: "P-1", OptionType: "leased_shell", DisplayName: "候选园区", BusinessRationale: "快速投产",
		EstimatedAmount:   AmountRange{Minimum: Money{"12000000.00", "CNY", 2}, Likely: Money{"15000000.00", "CNY", 2}, Maximum: Money{"18000000.00", "CNY", 2}, Basis: "Agent 估算"},
		EstimatedSchedule: ScheduleRange{Earliest: "2026-10-01T00:00:00Z", Likely: "2026-11-01T00:00:00Z", Latest: "2026-12-01T00:00:00Z"},
		Assumptions:       []string{"待调研"}, FactsRequired: []string{"园区报价"}, Risks: []string{"增容"}, SourceRefs: []string{"requirement:REQ-1"}, Confidence: "0.60", Status: "proposed",
	}
	request := InvestigationRequest{
		SchemaVersion: "1.0", InvestigationRequestID: "INV-1", CaseCode: requirement.CaseCode,
		ProposalSetID: "SET-1", ProposalID: proposal.ProposalID, ExpectedRevision: 1,
		WorldRunID: "world-1", Scope: []string{"ownership", "commercial_quote", "available_area", "electricity_capacity", "available_date", "permit"},
		RequestedBy: "project-lead", RequestedAt: "2026-08-02T10:00:00Z", Status: "waiting_world",
	}
	first, err := GenerateInvestigationObservation(request, requirement, proposal)
	if err != nil {
		t.Fatal(err)
	}
	second, err := GenerateInvestigationObservation(request, requirement, proposal)
	if err != nil {
		t.Fatal(err)
	}
	if CanonicalHash(first) != CanonicalHash(second) {
		t.Fatalf("world investigation is not replayable: first=%+v second=%+v", first, second)
	}
	if first.AvailableAreaM2 < requirement.MinimumAreaM2 || first.ElectricityKVA < requirement.MinimumElectricKVA ||
		first.ExternalActorID != "world-park-investigation-team" || len(first.EvidenceRefs) < 3 {
		t.Fatalf("unexpected generated observation %+v", first)
	}
	if err := ValidateInvestigationConfirmation(InvestigationConfirmation{SchemaVersion: "1.0", CaseCode: requirement.CaseCode, RequirementID: requirement.RequirementID, InvestigationRequestID: request.InvestigationRequestID, Action: "accept_report"}); err != nil {
		t.Fatal(err)
	}
}

func TestAIPlanningProviderProducesValidatedCandidateOnlySet(t *testing.T) {
	content := `{"proposals":[{"option_type":"leased_shell","display_name":"快速租赁改造","business_rationale":"缩短投产周期","estimated_amount":{"minimum":{"value":"12000000.00","currency":"CNY","scale":2},"likely":{"value":"15000000.00","currency":"CNY","scale":2},"maximum":{"value":"18000000.00","currency":"CNY","scale":2},"basis":"需求参数与待验证市场估算"},"estimated_schedule":{"earliest":"2026-10-01T00:00:00+08:00","likely":"2026-11-01T00:00:00+08:00","latest":"2026-12-01T00:00:00+08:00"},"assumptions":["存在标准厂房"],"facts_required":["租赁报价"],"risks":["增容延期"],"source_refs":["requirement:REQ-1"],"confidence":"0.62"},{"option_type":"build_to_suit","display_name":"定制代建","business_rationale":"平衡控制权和周期","estimated_amount":{"minimum":{"value":"15000000.00","currency":"CNY","scale":2},"likely":{"value":"17000000.00","currency":"CNY","scale":2},"maximum":{"value":"18000000.00","currency":"CNY","scale":2},"basis":"需求参数与待验证市场估算"},"estimated_schedule":{"earliest":"2026-11-01T00:00:00+08:00","likely":"2026-12-01T00:00:00+08:00","latest":"2027-01-01T00:00:00+08:00"},"assumptions":["园区可代建"],"facts_required":["交付承诺"],"risks":["承包商履约"],"source_refs":["requirement:REQ-1"],"confidence":"0.55"}]}`
	provider := AIPlanningProvider{Completer: planningCompleterStub{content}, Provider: "MiniMax", Model: "MiniMax-M3", Now: func() time.Time { return time.Date(2026, 8, 1, 10, 0, 0, 0, time.UTC) }}
	set, err := provider.Generate(context.Background(), requirementFixture())
	if err != nil {
		t.Fatal(err)
	}
	if len(set.Proposals) != 2 || set.Evidence.Provider != "MiniMax" || set.Evidence.TokenUsage["total_tokens"] != 42 {
		t.Fatalf("bad set %+v", set)
	}
	set.Proposals[1].ProposalID = set.Proposals[0].ProposalID
	if ValidateProposalSet(requirementFixture(), set) == nil {
		t.Fatal("duplicate proposal identity accepted")
	}
}

func TestAIPlanningProviderProducesRequirementOptionsWithinAuthorityLimits(t *testing.T) {
	content := `{"options":[{"title":"轻资产快速投产","business_rationale":"优先缩短投产周期","target_region":"江苏苏州及周边","facility_purpose":"电池冷却板制造与质检","minimum_area_m2":9000,"minimum_electricity_kva":1800,"target_available_at":"2026-12-01T00:00:00+08:00","candidate_count":3,"allowed_option_types":["lease_and_retrofit","build_to_suit"],"investment_request":{"value":"15000000.00","currency":"CNY","scale":2},"minimum_cash_reserve":{"value":"6000000.00","currency":"CNY","scale":2},"preferences":["快速投产"],"tradeoffs":["扩展空间有限"]},{"title":"均衡扩展","business_rationale":"兼顾周期与后续扩产","target_region":"长三角制造园区","facility_purpose":"冷却板制造、仓储与扩产预留","minimum_area_m2":12000,"minimum_electricity_kva":2200,"target_available_at":"2027-01-01T00:00:00+08:00","candidate_count":3,"allowed_option_types":["build_to_suit","existing_plant_purchase"],"investment_request":{"value":"18000000.00","currency":"CNY","scale":2},"minimum_cash_reserve":{"value":"7000000.00","currency":"CNY","scale":2},"preferences":["扩展性"],"tradeoffs":["准备周期较长"]}]}`
	provider := AIPlanningProvider{Completer: planningCompleterStub{content}, Provider: "MiniMax", Model: "MiniMax-M3", Now: func() time.Time { return time.Date(2026, 8, 2, 10, 0, 0, 0, time.UTC) }}
	requirement := requirementFixture()
	set, err := provider.GenerateRequirementOptions(context.Background(), RequirementOptionSeed{TenantID: requirement.TenantID, CaseCode: requirement.CaseCode, LegalEntityCode: requirement.LegalEntityCode, FinancialConstraint: requirement.FinancialConstraint})
	if err != nil {
		t.Fatal(err)
	}
	if len(set.Options) != 2 || set.Options[0].OptionID != "requirement-option-1" || set.Evidence.PromptVersion != RequirementPromptVersion {
		t.Fatalf("bad requirement options %+v", set)
	}
}

func TestAIPlanningProviderRepairsProposalAboveInvestmentCeiling(t *testing.T) {
	tooExpensive := `{"proposals":[{"option_type":"leased_shell","display_name":"越界方案一","business_rationale":"验证投资上限","estimated_amount":{"minimum":{"value":"17000000.00","currency":"CNY","scale":2},"likely":{"value":"19000000.00","currency":"CNY","scale":2},"maximum":{"value":"22000000.00","currency":"CNY","scale":2},"basis":"概念估算"},"estimated_schedule":{"earliest":"2026-10-01T00:00:00+08:00","likely":"2026-11-01T00:00:00+08:00","latest":"2026-12-01T00:00:00+08:00"},"assumptions":["假设"],"facts_required":["报价"],"risks":["超预算"],"source_refs":["requirement:REQ-1"],"confidence":"0.50"},{"option_type":"build_to_suit","display_name":"越界方案二","business_rationale":"验证投资上限","estimated_amount":{"minimum":{"value":"17000000.00","currency":"CNY","scale":2},"likely":{"value":"19000000.00","currency":"CNY","scale":2},"maximum":{"value":"22000000.00","currency":"CNY","scale":2},"basis":"概念估算"},"estimated_schedule":{"earliest":"2026-10-01T00:00:00+08:00","likely":"2026-11-01T00:00:00+08:00","latest":"2026-12-01T00:00:00+08:00"},"assumptions":["假设"],"facts_required":["报价"],"risks":["超预算"],"source_refs":["requirement:REQ-1"],"confidence":"0.50"}]}`
	valid := strings.ReplaceAll(tooExpensive, "22000000.00", "18000000.00")
	valid = strings.ReplaceAll(valid, "19000000.00", "17500000.00")
	completer := &planningSequenceCompleterStub{contents: []string{tooExpensive, valid}}
	provider := AIPlanningProvider{Completer: completer, Provider: "MiniMax", Model: "MiniMax-M3"}
	set, err := provider.Generate(context.Background(), requirementFixture())
	if err != nil {
		t.Fatal(err)
	}
	if completer.calls != 2 {
		t.Fatalf("calls=%d", completer.calls)
	}
	if set.Evidence.TokenUsage["total_tokens"] != 21 {
		t.Fatalf("usage=%v", set.Evidence.TokenUsage)
	}
}

func TestAIPlanningProviderRejectsMalformedOrUngroundedOutput(t *testing.T) {
	for name, content := range map[string]string{
		"bad_json":   `{not-json}`,
		"no_sources": `{"proposals":[{"option_type":"leased_shell","display_name":"方案一","business_rationale":"用于验证无来源失败关闭","estimated_amount":{"minimum":{"value":"1.00","currency":"CNY","scale":2},"likely":{"value":"2.00","currency":"CNY","scale":2},"maximum":{"value":"3.00","currency":"CNY","scale":2},"basis":"估算"},"estimated_schedule":{"earliest":"2026-10-01T00:00:00+08:00","likely":"2026-11-01T00:00:00+08:00","latest":"2026-12-01T00:00:00+08:00"},"assumptions":["假设"],"facts_required":["报价"],"risks":["风险"],"source_refs":[],"confidence":"0.50"},{"option_type":"build_to_suit","display_name":"方案二","business_rationale":"用于验证无来源失败关闭","estimated_amount":{"minimum":{"value":"1.00","currency":"CNY","scale":2},"likely":{"value":"2.00","currency":"CNY","scale":2},"maximum":{"value":"3.00","currency":"CNY","scale":2},"basis":"估算"},"estimated_schedule":{"earliest":"2026-10-01T00:00:00+08:00","likely":"2026-11-01T00:00:00+08:00","latest":"2026-12-01T00:00:00+08:00"},"assumptions":["假设"],"facts_required":["报价"],"risks":["风险"],"source_refs":[],"confidence":"0.50"}]}`,
	} {
		t.Run(name, func(t *testing.T) {
			provider := AIPlanningProvider{Completer: planningCompleterStub{content}, Provider: "MiniMax", Model: "MiniMax-M3"}
			if _, err := provider.Generate(context.Background(), requirementFixture()); err == nil {
				t.Fatal("invalid provider output accepted")
			}
		})
	}
}

func TestAIPlanningProviderProducesGovernedProjectWBSOptions(t *testing.T) {
	first := projectOptionFixture("快速总承包")
	second := projectOptionFixture("分段受控交付")
	raw, err := json.Marshal(map[string]any{"options": []ProjectPlanOption{first, second}})
	if err != nil {
		t.Fatal(err)
	}
	provider := AIPlanningProvider{Completer: planningCompleterStub{content: string(raw)}, Provider: "MiniMax", Model: "MiniMax-M3", Now: func() time.Time { return time.Date(2026, 8, 2, 10, 0, 0, 0, time.UTC) }}
	requirement := requirementFixture()
	set, err := provider.GenerateProjectPlanOptions(context.Background(), ProjectPlanSeed{TenantID: requirement.TenantID, CaseCode: requirement.CaseCode, SelectionID: "SEL-1", ControlObservationID: "OBS-CTRL-1", Requirement: requirement})
	if err != nil {
		t.Fatal(err)
	}
	if len(set.Options) != 2 || set.Options[0].OptionID != "project-option-1" || set.Evidence.PromptVersion != ProjectPromptVersion {
		t.Fatalf("bad project option set %+v", set)
	}
	invalid := set.Options[0]
	invalid.WBSItems[0].BudgetShareBPS++
	if ValidateProjectPlanOption(invalid, requirement) == nil {
		t.Fatal("unbalanced WBS budget shares accepted")
	}
}

func TestAIPlanningProviderRepairsTruncatedProjectOptionsOnce(t *testing.T) {
	first := projectOptionFixture("快速总承包")
	second := projectOptionFixture("分段受控交付")
	valid, err := json.Marshal(map[string]any{"options": []ProjectPlanOption{first, second}})
	if err != nil {
		t.Fatal(err)
	}
	completer := &planningSequenceCompleterStub{contents: []string{
		`{"options":[{"title":"被截断`,
		string(valid),
	}}
	provider := AIPlanningProvider{Completer: completer, Provider: "MiniMax", Model: "MiniMax-M3"}
	requirement := requirementFixture()
	set, err := provider.GenerateProjectPlanOptions(context.Background(), ProjectPlanSeed{
		TenantID: requirement.TenantID, CaseCode: requirement.CaseCode,
		SelectionID: "SEL-1", ControlObservationID: "OBS-CTRL-1", Requirement: requirement,
	})
	if err != nil {
		t.Fatal(err)
	}
	if completer.calls != 2 || len(set.Options) != 2 || set.Evidence.RequestID != "request-2" {
		t.Fatalf("calls=%d set=%+v", completer.calls, set)
	}
}

func TestAIPlanningProviderSeparatesCompletionAndGovernanceRepairs(t *testing.T) {
	first := projectOptionFixture("快速总承包")
	second := projectOptionFixture("分段受控交付")
	invalid := first
	invalid.WBSItems = append([]ProjectWBSItem(nil), first.WBSItems...)
	invalid.WBSItems[1].Sequence = 1
	invalidRaw, err := json.Marshal(map[string]any{"options": []ProjectPlanOption{invalid, second}})
	if err != nil {
		t.Fatal(err)
	}
	validRaw, err := json.Marshal(map[string]any{"options": []ProjectPlanOption{first, second}})
	if err != nil {
		t.Fatal(err)
	}
	completer := &planningSequenceCompleterStub{contents: []string{
		`{"options":[{"title":"被截断`,
		string(invalidRaw),
		string(validRaw),
	}}
	provider := AIPlanningProvider{Completer: completer, Provider: "MiniMax", Model: "MiniMax-M3"}
	requirement := requirementFixture()
	set, err := provider.GenerateProjectPlanOptions(context.Background(), ProjectPlanSeed{
		TenantID: requirement.TenantID, CaseCode: requirement.CaseCode,
		SelectionID: "SEL-1", ControlObservationID: "OBS-CTRL-1", Requirement: requirement,
	})
	if err != nil {
		t.Fatal(err)
	}
	if completer.calls != 3 || len(set.Options) != 2 || set.Evidence.RequestID != "request-3" {
		t.Fatalf("calls=%d set=%+v", completer.calls, set)
	}
	for _, maxTokens := range completer.maxTokens {
		if maxTokens > 4096 {
			t.Fatalf("project completion max_tokens=%d, want <=4096", maxTokens)
		}
	}
}

func TestAIPlanningProviderUsesCompletionRepairForTimeout(t *testing.T) {
	first := projectOptionFixture("快速总承包")
	second := projectOptionFixture("分段受控交付")
	validRaw, err := json.Marshal(map[string]any{"options": []ProjectPlanOption{first, second}})
	if err != nil {
		t.Fatal(err)
	}
	completer := &planningTimeoutThenContentCompleterStub{content: string(validRaw)}
	provider := AIPlanningProvider{Completer: completer, Provider: "MiniMax", Model: "MiniMax-M3"}
	requirement := requirementFixture()
	set, err := provider.GenerateProjectPlanOptions(context.Background(), ProjectPlanSeed{
		TenantID: requirement.TenantID, CaseCode: requirement.CaseCode,
		SelectionID: "SEL-1", ControlObservationID: "OBS-CTRL-1", Requirement: requirement,
	})
	if err != nil {
		t.Fatal(err)
	}
	if completer.calls != 2 || len(set.Options) != 2 || set.Evidence.RequestID != "request-recovery" {
		t.Fatalf("calls=%d set=%+v", completer.calls, set)
	}
}

func TestValidateProjectPlanOptionExplainsInvalidWBSField(t *testing.T) {
	option := projectOptionFixture("可解释校验")
	option.WBSItems[1].Sequence = 1
	err := ValidateProjectPlanOption(option, requirementFixture())
	if err == nil || !strings.Contains(err.Error(), `WBS item "WBS-02" sequence=1, want 2`) {
		t.Fatalf("error=%v", err)
	}
}

func TestGenerateContractBidObservationIsReplayStableAndBounded(t *testing.T) {
	rfq := ContractRFQ{SchemaVersion: "1.0", RFQID: "RFQ-1", CaseCode: "INC-1", ProjectID: "PROJECT-1",
		PackageCode: "WBS-02", PackageName: "厂房施工", SourcingStrategy: "specialist_packages", BidCount: 3,
		ContractCeiling: Money{Value: "12000000.00", Currency: "CNY", Scale: 2}, RequiredReadyAt: "2027-05-01T00:00:00Z",
		WorldRunID: "world-run-1", RequestedBy: "founder-principal", RequestedAt: "2026-08-02T12:00:00Z", Status: "waiting_world"}
	first, err := GenerateContractBidObservation(rfq)
	if err != nil {
		t.Fatal(err)
	}
	second, err := GenerateContractBidObservation(rfq)
	if err != nil {
		t.Fatal(err)
	}
	if CanonicalHash(first) != CanonicalHash(second) || len(first.Bids) != rfq.BidCount {
		t.Fatalf("contractor market fact is not replay stable: first=%+v second=%+v", first, second)
	}
	ceiling, _ := new(big.Rat).SetString(rfq.ContractCeiling.Value)
	for _, bid := range first.Bids {
		amount, _ := new(big.Rat).SetString(bid.QuotedAmount.Value)
		if amount.Cmp(ceiling) > 0 || bid.Qualification != "eligible" || len(bid.EvidenceRefs) < 2 || !strings.Contains(bid.ContractorName, "虚构") {
			t.Fatalf("generated bid is not governed and bounded: %+v", bid)
		}
	}
}

func TestGenerateConstructionProgressObservationIsReplayStableAndWorldOwned(t *testing.T) {
	execution := ConstructionExecution{SchemaVersion: "1.0", ExecutionID: "construction-1", CaseCode: "INC-1", ProjectID: "PROJECT-1", ContractID: "CONTRACT-1", PackageCode: "WBS-1", PackageName: "厂房工程", ContractorCode: "WORLD-CONTRACTOR-01", ContractorName: "澄岳工程（虚构）", WorldRunID: "world-run-1", StartedBy: "founder-principal", StartedAt: "2026-08-02T00:00:00Z", Status: "waiting_world"}
	first, err := GenerateConstructionProgressObservation(execution)
	if err != nil {
		t.Fatal(err)
	}
	second, err := GenerateConstructionProgressObservation(execution)
	if err != nil {
		t.Fatal(err)
	}
	if CanonicalHash(first) != CanonicalHash(second) {
		t.Fatalf("construction fact is not replay stable: first=%+v second=%+v", first, second)
	}
	if first.ExternalActorID != "world-construction-contractor" || first.ProgressBPS != 10000 || first.QualityStatus != "passed" || first.SafetyStatus != "clear" || len(first.EvidenceRefs) < 3 {
		t.Fatalf("invalid governed construction fact: %+v", first)
	}
}

func TestContractRecommendationEvidenceMatchesIAOSCanonicalContract(t *testing.T) {
	rfq := ContractRFQ{SchemaVersion: "1.0", RFQID: "RFQ-2", CaseCode: "INC-2", ProjectID: "PROJECT-2", PackageCode: "WBS-03", PackageName: "现场施工", SourcingStrategy: "specialist_packages", BidCount: 3, ContractCeiling: Money{Value: "9000000.00", Currency: "CNY", Scale: 2}, RequiredReadyAt: "2027-03-01T00:00:00Z", WorldRunID: "world-run-2", RequestedBy: "founder-principal", RequestedAt: "2026-08-02T12:00:00Z", Status: "waiting_world"}
	observation, err := GenerateContractBidObservation(rfq)
	if err != nil {
		t.Fatal(err)
	}
	selected := observation.Bids[0].BidID
	content := `{"selected_bid_id":"` + selected + `","recommendation_reason":"成本、工期和质保条件综合最优","alternative_comparison":"其余方案报价或交付条件相对较弱"}`
	provider := AIPlanningProvider{Completer: planningCompleterStub{content: content}, Provider: "MiniMax", Model: "MiniMax-M3", Now: func() time.Time { return time.Date(2026, 8, 2, 16, 0, 0, 0, time.UTC) }}
	seed := ContractRecommendationSeed{RFQ: rfq, Observation: observation}
	advice, err := provider.GenerateContractRecommendation(context.Background(), seed)
	if err != nil {
		t.Fatal(err)
	}
	want := CanonicalHash(map[string]any{"rfq": rfq, "observation": observation})
	if advice.Evidence.InputHash != want || advice.Evidence.PromptVersion != ContractPromptVersion {
		t.Fatalf("evidence=%+v want_input_hash=%s", advice.Evidence, want)
	}
}
