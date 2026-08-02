import type { FacilityRequirement, ProposalSet, SiteInvestigationItem } from "./plantBuild";

export type SiteAssessmentWeights = { cost: number; schedule: number; capacity: number; control: number };
export const DEFAULT_SITE_ASSESSMENT_WEIGHTS: SiteAssessmentWeights = { cost: 35, schedule: 25, capacity: 20, control: 20 };
export const SITE_ASSESSMENT_POLICY_VERSION = "site-assessment-v1";
export type SiteAssessmentCriterion = {
  key: "investment" | "area" | "electricity" | "available_at" | "ownership" | "permit";
  label: string;
  required: string;
  observed: string;
  difference: string;
  passed: boolean;
  failure?: string;
};
export type SiteAssessment = {
  proposal_id: string;
  display_name: string;
  observation_id: string;
  eligible: boolean;
  hard_failures: string[];
  total_score: number | null;
  component_scores: SiteAssessmentWeights;
  criteria: SiteAssessmentCriterion[];
  observed: { quote: string; available_at: string; area_m2: number; electricity_kva: number; ownership: string; permit: string };
  estimated?: { amount: string; available_at: string };
  evidence_refs: string[];
};

const clamp = (value: number) => Math.max(0, Math.min(100, value));
const ratioScore = (actual: number, minimum: number) => clamp(50 + (actual / minimum) * 25);
const quantity = (value: number) => new Intl.NumberFormat("en-US", { maximumFractionDigits: 2 }).format(value);
const money = (value: number, currency: string) => `${currency} ${new Intl.NumberFormat("en-US", { minimumFractionDigits: 2, maximumFractionDigits: 2 }).format(value)}`;
const day = (value: string) => value.slice(0, 10) || "无有效日期";

export function currentProposalInvestigations(proposalSet: ProposalSet | null, items: SiteInvestigationItem[]): SiteInvestigationItem[] {
  if (!proposalSet) return [];
  const proposalIDs = new Set(proposalSet.proposals.map((proposal) => proposal.proposal_id));
  return items.filter((item) =>
    item.request.proposal_set_id === proposalSet.proposal_set_id &&
    item.request.expected_revision === proposalSet.revision &&
    proposalIDs.has(item.request.proposal_id),
  );
}

export function assessObservedSites(requirement: FacilityRequirement, proposalSet: ProposalSet | null, items: SiteInvestigationItem[], weights: SiteAssessmentWeights): SiteAssessment[] {
  if (!proposalSet) return [];
  const totalWeight = Object.values(weights).reduce((sum, value) => sum + Math.max(0, value), 0) || 1;
  const investmentLimit = Number(requirement.investment_request.value);
  const targetTime = Date.parse(requirement.target_available_at);
  return currentProposalInvestigations(proposalSet, items).filter((item) => item.observation).map((item) => {
    const fact = item.observation!;
    const quote = Number(fact.quoted_amount.value);
    const availableTime = Date.parse(fact.available_at);
    const sameCurrency = fact.quoted_amount.currency === requirement.investment_request.currency;
    const quotePassed = sameCurrency && quote <= investmentLimit;
    const schedulePassed = Number.isFinite(availableTime) && availableTime <= targetTime;
    const scheduleDays = Number.isFinite(availableTime) && Number.isFinite(targetTime) ? Math.ceil(Math.abs(targetTime - availableTime) / 86_400_000) : 0;
    const criteria: SiteAssessmentCriterion[] = [
      {
        key: "investment", label: "投资金额", required: `不超过 ${money(investmentLimit, requirement.investment_request.currency)}`,
        observed: `正式报价 ${money(quote, fact.quoted_amount.currency)}`,
        difference: !sameCurrency ? "币种不一致" : quotePassed ? `低于上限 ${money(investmentLimit - quote, fact.quoted_amount.currency)}` : `超出上限 ${money(quote - investmentLimit, fact.quoted_amount.currency)}`,
        passed: quotePassed, failure: quotePassed ? undefined : "正式报价超过投资申请额或币种不一致",
      },
      {
        key: "area", label: "可用面积", required: `最低要求 ${quantity(requirement.minimum_area_m2)} m²`,
        observed: `实测 ${quantity(fact.available_area_m2)} m²`,
        difference: fact.available_area_m2 >= requirement.minimum_area_m2 ? `高于门槛 ${quantity(fact.available_area_m2 - requirement.minimum_area_m2)} m²` : `短缺 ${quantity(requirement.minimum_area_m2 - fact.available_area_m2)} m²`,
        passed: fact.available_area_m2 >= requirement.minimum_area_m2, failure: fact.available_area_m2 >= requirement.minimum_area_m2 ? undefined : "实际可用面积低于最低要求",
      },
      {
        key: "electricity", label: "可用电力", required: `最低要求 ${quantity(requirement.minimum_electricity_kva)} kVA`,
        observed: `实测 ${quantity(fact.electricity_kva)} kVA`,
        difference: fact.electricity_kva >= requirement.minimum_electricity_kva ? `高于门槛 ${quantity(fact.electricity_kva - requirement.minimum_electricity_kva)} kVA` : `短缺 ${quantity(requirement.minimum_electricity_kva - fact.electricity_kva)} kVA`,
        passed: fact.electricity_kva >= requirement.minimum_electricity_kva, failure: fact.electricity_kva >= requirement.minimum_electricity_kva ? undefined : "可用电力低于最低要求",
      },
      {
        key: "available_at", label: "可用日期", required: `不晚于 ${day(requirement.target_available_at)}`,
        observed: `实测 ${day(fact.available_at)}`,
        difference: !Number.isFinite(availableTime) ? "日期格式无效" : schedulePassed ? `提前 ${scheduleDays} 天` : `晚于目标 ${scheduleDays} 天`,
        passed: schedulePassed, failure: schedulePassed ? undefined : "实际可用日期晚于目标日期",
      },
      {
        key: "ownership", label: "权属核验", required: "必须提供核验结论", observed: fact.ownership_status || "未提供",
        difference: fact.ownership_status ? "结论已提供" : "缺少结论", passed: Boolean(fact.ownership_status), failure: fact.ownership_status ? undefined : "缺少权属核验结论",
      },
      {
        key: "permit", label: "许可条件", required: "必须提供许可结论", observed: fact.permit_status || "未提供",
        difference: fact.permit_status ? "结论已提供" : "缺少结论", passed: Boolean(fact.permit_status), failure: fact.permit_status ? undefined : "缺少许可结论",
      },
    ];
    const failures = criteria.flatMap((criterion) => criterion.failure ? [criterion.failure] : []);
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
      total_score: failures.length === 0 ? Math.round(score * 10) / 10 : null, component_scores: components, criteria,
      observed: { quote: fact.quoted_amount.value, available_at: fact.available_at, area_m2: fact.available_area_m2, electricity_kva: fact.electricity_kva, ownership: fact.ownership_status, permit: fact.permit_status },
      estimated: proposal ? { amount: proposal.estimated_amount.likely.value, available_at: proposal.estimated_schedule.likely } : undefined,
      evidence_refs: fact.evidence_refs,
    };
  }).sort((a, b) => Number(b.eligible) - Number(a.eligible) || (b.total_score ?? -1) - (a.total_score ?? -1));
}
