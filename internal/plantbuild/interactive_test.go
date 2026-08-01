package plantbuild

import (
	"context"
	"errors"
	"testing"
	"time"
)

type planningCompleterStub struct{ content string }

func (s planningCompleterStub) CompleteJSON(context.Context, string, string, float64, int) (string, string, map[string]int, error) {
	return s.content, "request-1", map[string]int{"total_tokens": 42}, nil
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

func TestAIPlanningProviderProducesValidatedCandidateOnlySet(t *testing.T) {
	content := `{"proposals":[{"option_type":"leased_shell","display_name":"快速租赁改造","business_rationale":"缩短投产周期","estimated_amount":{"minimum":{"value":"12000000.00","currency":"CNY","scale":2},"likely":{"value":"15000000.00","currency":"CNY","scale":2},"maximum":{"value":"18000000.00","currency":"CNY","scale":2},"basis":"需求参数与待验证市场估算"},"estimated_schedule":{"earliest":"2026-10-01T00:00:00+08:00","likely":"2026-11-01T00:00:00+08:00","latest":"2026-12-01T00:00:00+08:00"},"assumptions":["存在标准厂房"],"facts_required":["租赁报价"],"risks":["增容延期"],"source_refs":["requirement:REQ-1"],"confidence":"0.62"},{"option_type":"build_to_suit","display_name":"定制代建","business_rationale":"平衡控制权和周期","estimated_amount":{"minimum":{"value":"16000000.00","currency":"CNY","scale":2},"likely":{"value":"18000000.00","currency":"CNY","scale":2},"maximum":{"value":"20000000.00","currency":"CNY","scale":2},"basis":"需求参数与待验证市场估算"},"estimated_schedule":{"earliest":"2026-11-01T00:00:00+08:00","likely":"2026-12-01T00:00:00+08:00","latest":"2027-01-01T00:00:00+08:00"},"assumptions":["园区可代建"],"facts_required":["交付承诺"],"risks":["承包商履约"],"source_refs":["requirement:REQ-1"],"confidence":"0.55"}]}`
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
