import type { FacilityRequirement, ProposalSet, SiteInvestigationItem } from "./plantBuild";

export type SiteAssessmentWeights = { cost: number; schedule: number; capacity: number; control: number };
export type SiteAssessment = {
  proposal_id: string;
  display_name: string;
  observation_id: string;
  eligible: boolean;
  hard_failures: string[];
  total_score: number | null;
  component_scores: SiteAssessmentWeights;
  observed: { quote: string; available_at: string; area_m2: number; electricity_kva: number; ownership: string; permit: string };
  estimated?: { amount: string; available_at: string };
  evidence_refs: string[];
};

const clamp = (value: number) => Math.max(0, Math.min(100, value));
const ratioScore = (actual: number, minimum: number) => clamp(50 + (actual / minimum) * 25);

export function assessObservedSites(requirement: FacilityRequirement, proposalSet: ProposalSet | null, items: SiteInvestigationItem[], weights: SiteAssessmentWeights): SiteAssessment[] {
  const totalWeight = Object.values(weights).reduce((sum, value) => sum + Math.max(0, value), 0) || 1;
  const investmentLimit = Number(requirement.investment_request.value);
  const targetTime = Date.parse(requirement.target_available_at);
  return items.filter((item) => item.observation).map((item) => {
    const fact = item.observation!;
    const quote = Number(fact.quoted_amount.value);
    const availableTime = Date.parse(fact.available_at);
    const failures: string[] = [];
    if (fact.quoted_amount.currency !== requirement.investment_request.currency || quote > investmentLimit) failures.push("正式报价超过投资申请额或币种不一致");
    if (fact.available_area_m2 < requirement.minimum_area_m2) failures.push("实际可用面积低于最低要求");
    if (fact.electricity_kva < requirement.minimum_electricity_kva) failures.push("可用电力低于最低要求");
    if (!Number.isFinite(availableTime) || availableTime > targetTime) failures.push("实际可用日期晚于目标日期");
    if (!fact.ownership_status) failures.push("缺少权属核验结论");
    if (!fact.permit_status) failures.push("缺少许可结论");
    const daysEarly = Math.max(0, (targetTime - availableTime) / 86_400_000);
    const components: SiteAssessmentWeights = {
      cost: clamp((1 - quote / investmentLimit) * 100),
      schedule: clamp(50 + Math.min(daysEarly, 180) / 180 * 50),
      capacity: (ratioScore(fact.available_area_m2, requirement.minimum_area_m2) + ratioScore(fact.electricity_kva, requirement.minimum_electricity_kva)) / 2,
      control: ((fact.ownership_status === "verified" ? 100 : 60) + (fact.permit_status === "eligible" ? 100 : 60)) / 2,
    };
    const score = Object.entries(components).reduce((sum, [key, value]) => sum + value * Math.max(0, weights[key as keyof SiteAssessmentWeights]), 0) / totalWeight;
    const proposal = proposalSet?.proposals.find((candidate) => candidate.proposal_id === fact.proposal_id);
    return {
      proposal_id: fact.proposal_id, display_name: proposal?.display_name ?? fact.proposal_id,
      observation_id: fact.observation_id, eligible: failures.length === 0, hard_failures: failures,
      total_score: failures.length === 0 ? Math.round(score * 10) / 10 : null, component_scores: components,
      observed: { quote: fact.quoted_amount.value, available_at: fact.available_at, area_m2: fact.available_area_m2, electricity_kva: fact.electricity_kva, ownership: fact.ownership_status, permit: fact.permit_status },
      estimated: proposal ? { amount: proposal.estimated_amount.likely.value, available_at: proposal.estimated_schedule.likely } : undefined,
      evidence_refs: fact.evidence_refs,
    };
  }).sort((a, b) => Number(b.eligible) - Number(a.eligible) || (b.total_score ?? -1) - (a.total_score ?? -1));
}
