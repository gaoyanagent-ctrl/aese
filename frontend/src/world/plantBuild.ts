export type PBMoney = { value: string; currency: "CNY"; scale: 2 };
export type PlantBuildFrame = {
  step: number;
  phase: string;
  sim_time: string;
  title: string;
  causation_id: string;
  selected_site: string;
  assessments: Array<{
    site_code: string;
    feasible: boolean;
    hard_failures: string[];
    weighted_score: string;
    source: string;
    confidence: string;
  }>;
  zones: Array<{
    code: string;
    parent: string;
    purpose: string;
    status: string;
    area_m2: number;
  }>;
  work_packages: Array<{
    code: string;
    status: string;
    cost: PBMoney;
    evidence?: string;
  }>;
  utilities: Record<string, string>;
  knowledge: Array<{ actor: string; fact: string }>;
  world_progress: number;
  iaos_plan_progress: number;
  discrepancy: string;
  cash: PBMoney;
  committed: PBMoney;
  payable: PBMoney;
  paid: PBMoney;
  capability_build_eligible: boolean;
  iaos_cursor: number;
};
export type PlantBuildTrace = {
  schema_version: string;
  campaign: "plant-build";
  world_run_id: string;
  timezone: string;
  policy_version: string;
  m9_terminal_hash: string;
  frames: PlantBuildFrame[];
};

export type PlantPlanningProviderStatus = {
  state: "connected" | "not_configured" | "degraded";
  provider: string;
  model?: string;
  prompt_version: string;
};

export type FacilityMoney = {
  value: string;
  currency: "CNY";
  scale: 2;
};

export type FacilityRequirement = {
  schema_version: "1.0";
  requirement_id: string;
  tenant_id: string;
  case_code: string;
  legal_entity_code: string;
  target_region: string;
  facility_purpose: string;
  minimum_area_m2: number;
  minimum_electricity_kva: number;
  target_available_at: string;
  candidate_count: number;
  allowed_option_types: string[];
  investment_request: FacilityMoney;
  minimum_cash_reserve: FacilityMoney;
  financial_constraint: {
    available_cash: FacilityMoney;
    approved_budget: FacilityMoney;
    cash_source_ref: string;
    budget_source_ref: string;
    snapshot_hash: string;
  };
  preferences: string[];
  revision: number;
  revision_reason: string;
};

export type PlantFinancialConstraint = {
  case_code: string;
  legal_entity_code: string;
  financial_constraint: FacilityRequirement["financial_constraint"];
};

export type PlantRequirementOption = {
  option_id: string;
  title: string;
  business_rationale: string;
  target_region: string;
  facility_purpose: string;
  minimum_area_m2: number;
  minimum_electricity_kva: number;
  target_available_at: string;
  candidate_count: number;
  allowed_option_types: string[];
  investment_request: FacilityMoney;
  minimum_cash_reserve: FacilityMoney;
  preferences: string[];
  tradeoffs: string[];
};

export type PlantRequirementOptionSet = {
  schema_version: "1.0";
  options: PlantRequirementOption[];
  evidence: ProposalSet["evidence"];
};

export type SiteOptionProposal = {
  proposal_id: string;
  option_type: string;
  display_name: string;
  business_rationale: string;
  estimated_amount: {
    minimum: FacilityMoney;
    likely: FacilityMoney;
    maximum: FacilityMoney;
    basis: string;
  };
  estimated_schedule: {
    earliest: string;
    likely: string;
    latest: string;
  };
  assumptions: string[];
  facts_required: string[];
  risks: string[];
  source_refs: string[];
  confidence: string;
  status: string;
};

export type ProposalSet = {
  schema_version: string;
  proposal_set_id: string;
  requirement_id: string;
  revision: number;
  status: "candidate_only";
  proposals: SiteOptionProposal[];
  evidence: {
    provider: string;
    model: string;
    prompt_version: string;
    source_type?: string;
    parent_revision?: number;
    request_id?: string;
    input_hash: string;
    output_hash: string;
    token_usage?: Record<string, number>;
    validated_at: string;
  };
};

export type PlantPlanningResult = {
  proposal_set: ProposalSet;
  agent_job: unknown;
  idempotent_replay: boolean;
};

export type SiteProposalReviewInput = {
  proposal_set_id: string;
  proposal_id: string;
  action: "adopt_for_investigation" | "request_revision" | "discard";
  reason: string;
  expected_revision: number;
};

export type SiteInvestigationRequest = {
  schema_version: "1.0";
  investigation_request_id: string;
  case_code: string;
  proposal_set_id: string;
  proposal_id: string;
  expected_revision: number;
  world_run_id: string;
  scope: string[];
  requested_by: string;
  requested_at: string;
  status: "waiting_world";
};

export type SiteInvestigationObservation = {
  schema_version: "1.0";
  observation_id: string;
  investigation_request_id: string;
  proposal_id: string;
  result: "completed";
  ownership_status: string;
  available_area_m2: number;
  electricity_kva: number;
  quoted_amount: FacilityMoney;
  available_at: string;
  permit_status: string;
  evidence_refs: string[];
  notes: string;
  external_actor_id: string;
  observed_at: string;
};

export type SiteInvestigationItem = {
  request: SiteInvestigationRequest;
  status: "waiting_world" | "observed" | "cancelled";
  work_item_status: "waiting_world" | "completed" | "cancelled";
  observation?: SiteInvestigationObservation;
};

export type SiteSelectionRecommendationInput = {
  schema_version: "1.0";
  recommendation_id: string;
  case_code: string;
  proposal_set_id: string;
  proposal_set_revision: number;
  selected_proposal_id: string;
  assessment_policy_version: "site-assessment-v1";
  weights: { cost: number; schedule: number; capacity: number; control: number };
  recommendation_reason: string;
  alternative_comparison: string;
  single_source_exception_reason?: string;
  recommended_at: string;
};

export type SiteSelectionItem = {
  recommendation: SiteSelectionRecommendationInput & {
    requirement_id: string;
    input_hash: string;
    approval_flow_key: string;
    approval_request_id: string;
    eligible_count: number;
    status: string;
    recommended_by: string;
    assessments: Array<{
      proposal_id: string;
      display_name: string;
      observation_id: string;
      eligible: boolean;
      hard_failures: string[];
      total_score: number | null;
      evidence_refs: string[];
    }>;
  };
  record_status: string;
  approval_status: "pending" | "approved" | "rejected" | "consumed";
  decision?: { selection_id: string; selected_proposal_id: string; formalized_at: string };
};

export type SiteControlRequest = {
  schema_version: "1.0";
  control_request_id: string;
  selection_id: string;
  case_code: string;
  selected_proposal_id: string;
  world_run_id: string;
  agreement_mode: "lease" | "purchase" | "build_to_suit" | "use_agreement";
  requested_handover_at: string;
  required_evidence: Array<"executed_agreement" | "handover_record" | "possession_authority">;
  requested_by: string;
  requested_at: string;
  status: "waiting_world";
};

export type SiteControlObservation = {
  schema_version: "1.0";
  observation_id: string;
  control_request_id: string;
  selection_id: string;
  result: "delivered" | "delayed" | "rejected";
  agreement_ref: string;
  handover_ref: string;
  effective_at?: string;
  evidence_refs: string[];
  notes: string;
  external_actor_id: string;
  observed_at: string;
};

export type SiteControlItem = {
  request: SiteControlRequest;
  status: "waiting_world" | "controlled" | "delayed" | "rejected" | "cancelled";
  observation?: SiteControlObservation;
};

export type ProjectWBSItem = {
  wbs_code: string; name: string; phase: "design" | "procurement" | "construction" | "commissioning";
  sequence: number; owner_position: string; planned_start_at: string; planned_finish_at: string;
  budget_share_bps: number; acceptance_criteria: string;
};

export type ProjectPlanOption = {
  option_id: string; title: string; business_rationale: string; project_name: string;
  delivery_strategy: "design_bid_build" | "design_build" | "epcm";
  budget_ceiling: FacilityMoney; target_start_at: string; target_ready_at: string;
  wbs_items: ProjectWBSItem[]; tradeoffs: string[];
};

export type ProjectPlanOptionSet = { schema_version: "1.0"; options: ProjectPlanOption[]; evidence: ProposalSet["evidence"] };

export type FacilityProjectItem = {
  plan: {
    schema_version: "1.0";
    plan_id: string;
    case_code: string;
    project_name: string;
    delivery_strategy: ProjectPlanOption["delivery_strategy"];
    budget_ceiling: PBMoney;
    target_start_at: string;
    target_ready_at: string;
    wbs_items: ProjectWBSItem[];
    status: string;
    agent_evidence: Record<string, string>;
  };
  plan_hash: string; status: string; approval_request_id: string;
  approval_status: "" | "pending" | "approved" | "rejected" | "consumed";
  project?: { project_id?: string; project_name?: string; status?: string; wbs_items?: ProjectWBSItem[] };
};

export type ContractRFQ = {
  schema_version: "1.0"; rfq_id: string; case_code: string; project_id: string;
  package_code: string; package_name: string;
  sourcing_strategy: "general_contract" | "specialist_packages" | "epcm_managed";
  bid_count: number; contract_ceiling: PBMoney; required_ready_at: string;
  world_run_id: string; requested_by: string; requested_at: string; status: string;
};

export type ContractBid = {
  schema_version: "1.0"; bid_id: string; rfq_id: string; contractor_code: string;
  contractor_name: string; quoted_amount: PBMoney; promised_ready_at: string;
  qualification: string; warranty_months: number; milestone_count: number;
  evidence_refs: string[]; observed_at: string;
};

export type ContractBidObservation = {
  schema_version: "1.0"; observation_id: string; rfq_id: string;
  external_actor_id: string; bids: ContractBid[]; observed_at: string;
};

export type ContractRecommendationAdvice = {
  selected_bid_id: string; recommendation_reason: string; alternative_comparison: string;
  evidence: ProposalSet["evidence"];
};

export type ContractAwardRecommendation = {
  schema_version: "1.0"; recommendation_id: string; rfq_id: string;
  selected_bid_id: string; recommendation_reason: string; alternative_comparison: string;
  recommended_at: string;
};

export type ContractAwardItem = {
  rfq: ContractRFQ; status: string; observation?: ContractBidObservation;
  recommendation?: ContractAwardRecommendation; approval_request_id: string;
  approval_status: "" | "pending" | "approved" | "rejected" | "consumed";
  contract?: { contract_id?: string; contractor_name?: string; committed_amount?: PBMoney; status?: string; package_name?: string };
};

export type ConstructionExecution = {
  schema_version: "1.0"; execution_id: string; case_code: string; project_id: string;
  contract_id: string; package_code: string; package_name: string; contractor_code: string;
  contractor_name: string; world_run_id: string; started_by: string; started_at: string; status: string;
};
export type ConstructionProgressObservation = {
  schema_version: "1.0"; observation_id: string; execution_id: string; contract_id: string;
  package_code: string; result: string; progress_bps: number; quality_status: string;
  safety_status: string; punch_items: string[]; evidence_refs: string[];
  external_actor_id: string; observed_at: string;
};
export type ConstructionMilestoneItem = {
  execution: ConstructionExecution; status: string; observation?: ConstructionProgressObservation;
  acceptance?: { acceptance_id?: string; decision?: string; accepted_by?: string; accepted_at?: string; payment_status?: string };
};

export type PlantApprovalDetail = {
  item: {
    id: string;
    status: "pending" | "approved" | "rejected" | "cancelled" | "expired" | "consumed";
    requester_id: string;
    requester_display_name?: string;
    decision_note?: string;
    decided_by?: string;
    decided_at?: string;
  };
  detail: {
    flow_key: string;
    flow_version: number;
    flow_name: string;
    subject: {
      title: string;
      summary: string;
      operation: string;
      fields?: Record<string, unknown>;
      evidence?: string[];
    };
    requester_subject_id?: string;
    requester_display_name?: string;
    assignments: Array<{
      id: string;
      stage_code: string;
      stage_name: string;
      mode: string;
      display_name: string;
      selector_type: string;
      selector_value: string;
      status: string;
      decision?: string;
      decision_note?: string;
    }>;
    can_decide: boolean;
  };
};

export async function submitPlantProposalReview(input: SiteProposalReviewInput) {
  const response = await fetch("/api/aese/v1/world/plant-build/reviews", {
    method: "POST",
    headers: { "content-type": "application/json", ...iaosHeaders() },
    body: JSON.stringify(input),
  });
  if (!response.ok) {
    throw new Error(`保存候选审阅 ${response.status}: ${await response.text()}`);
  }
  return response.json() as Promise<{ status: "committed"; proposal_review: SiteProposalReviewInput & { reviewed_by: string; reviewed_at: string } }>;
}

export async function loadPlantProposalReviews(proposalSetID: string, signal?: AbortSignal) {
  const response = await fetch(`/api/aese/v1/world/plant-build/reviews?proposal_set_id=${encodeURIComponent(proposalSetID)}`, {
    headers: iaosHeaders(), signal,
  });
  if (!response.ok) throw new Error(`候选审阅记录 ${response.status}: ${await response.text()}`);
  return response.json() as Promise<{ items: Array<SiteProposalReviewInput & { reviewed_by: string; reviewed_at: string }> }>;
}

export async function loadSiteInvestigations(caseCode: string, signal?: AbortSignal) {
  const response = await fetch(`/api/aese/v1/world/plant-build/investigations?case_code=${encodeURIComponent(caseCode)}`, {
    headers: iaosHeaders(), signal,
  });
  if (!response.ok) throw new Error(`场址调研工作项 ${response.status}: ${await response.text()}`);
  return response.json() as Promise<{ items: SiteInvestigationItem[] }>;
}

export async function requestSiteInvestigation(input: SiteInvestigationRequest) {
  const response = await fetch("/api/aese/v1/world/plant-build/investigations", {
    method: "POST", headers: { "content-type": "application/json", ...iaosHeaders() }, body: JSON.stringify(input),
  });
  if (!response.ok) throw new Error(`发起场址调研 ${response.status}: ${await response.text()}`);
  return response.json() as Promise<{ status: "waiting_world"; investigation_request: SiteInvestigationRequest }>;
}

export async function confirmSiteInvestigationReport(caseCode: string, requirementID: string, investigationRequestID: string) {
  const response = await fetch("/api/aese/v1/world/plant-build/observations", {
    method: "POST", headers: { "content-type": "application/json", ...iaosHeaders() },
    body: JSON.stringify({ schema_version: "1.0", case_code: caseCode, requirement_id: requirementID, investigation_request_id: investigationRequestID, action: "accept_report" }),
  });
  if (!response.ok) throw new Error(`确认接收园区调研报告 ${response.status}: ${await response.text()}`);
  return response.json() as Promise<{ status: "committed"; idempotent_replay: boolean; world_message_id: string; observation: SiteInvestigationObservation }>;
}

export async function loadSiteSelections(caseCode: string, signal?: AbortSignal) {
  const response = await fetch(`/api/aese/v1/world/plant-build/site-selections?case_code=${encodeURIComponent(caseCode)}`, {
    headers: iaosHeaders(), signal,
  });
  if (!response.ok) throw new Error(`场址推荐与审批 ${response.status}: ${await response.text()}`);
  return response.json() as Promise<{ items: SiteSelectionItem[] }>;
}

export async function submitSiteSelectionRecommendation(input: SiteSelectionRecommendationInput) {
  const response = await fetch("/api/aese/v1/world/plant-build/site-selections", {
    method: "POST", headers: { "content-type": "application/json", ...iaosHeaders() }, body: JSON.stringify(input),
  });
  if (!response.ok) throw new Error(`提交正式推荐 ${response.status}: ${await response.text()}`);
  return response.json() as Promise<{ status: "committed"; result: SiteSelectionItem["recommendation"] }>;
}

export async function loadPlantApprovalDetail(approvalRequestID: string, signal?: AbortSignal) {
  const response = await fetch(`/api/aese/v1/world/plant-build/approvals/${encodeURIComponent(approvalRequestID)}`, {
    headers: iaosHeaders(), signal,
  });
  if (!response.ok) throw new Error(`读取游戏内审批事项 ${response.status}: ${await response.text()}`);
  return response.json() as Promise<PlantApprovalDetail>;
}

export async function decidePlantApproval(approvalRequestID: string, decision: "approve" | "reject", note: string) {
  const response = await fetch(`/api/aese/v1/commands/iaos/approvals/${encodeURIComponent(approvalRequestID)}/${decision}`, {
    method: "POST", headers: { "content-type": "application/json", ...iaosHeaders() }, body: JSON.stringify({ note }),
  });
  if (!response.ok) throw new Error(`${decision === "approve" ? "批准" : "驳回"}选址审批 ${response.status}: ${await response.text()}`);
  return response.json() as Promise<{ item: PlantApprovalDetail["item"] }>;
}

export async function finalizeSiteSelection(caseCode: string, recommendationID: string, approvalRequestID: string) {
  const response = await fetch("/api/aese/v1/world/plant-build/site-selections/finalize", {
    method: "POST", headers: { "content-type": "application/json", ...iaosHeaders() },
    body: JSON.stringify({ schema_version: "1.0", case_code: caseCode, recommendation_id: recommendationID, approval_request_id: approvalRequestID, formalized_at: new Date().toISOString() }),
  });
  if (!response.ok) throw new Error(`正式选址落地 ${response.status}: ${await response.text()}`);
  return response.json() as Promise<{ status: "committed"; result: SiteSelectionItem["decision"] }>;
}

export async function loadSiteControls(caseCode: string, signal?: AbortSignal) {
  const response = await fetch(`/api/aese/v1/world/plant-build/site-controls?case_code=${encodeURIComponent(caseCode)}`, {
    headers: iaosHeaders(), signal,
  });
  if (!response.ok) throw new Error(`场地控制工作项 ${response.status}: ${await response.text()}`);
  return response.json() as Promise<{ items: SiteControlItem[] }>;
}

export async function requestSiteControl(input: SiteControlRequest) {
  const response = await fetch("/api/aese/v1/world/plant-build/site-controls", {
    method: "POST", headers: { "content-type": "application/json", ...iaosHeaders() }, body: JSON.stringify(input),
  });
  if (!response.ok) throw new Error(`发起场地控制交付 ${response.status}: ${await response.text()}`);
  return response.json() as Promise<{ status: "waiting_world"; site_control_request: SiteControlRequest }>;
}

export async function confirmSiteControlDelivery(caseCode: string, controlRequestID: string) {
  const response = await fetch("/api/aese/v1/world/plant-build/site-controls/observations", {
    method: "POST", headers: { "content-type": "application/json", ...iaosHeaders() },
    body: JSON.stringify({ schema_version: "1.0", case_code: caseCode, control_request_id: controlRequestID, action: "accept_delivery" }),
  });
  if (!response.ok) throw new Error(`确认接收场地 ${response.status}: ${await response.text()}`);
  return response.json() as Promise<{ status: "committed"; idempotent_replay: boolean; world_message_id: string; observation: SiteControlObservation }>;
}

export async function generateFacilityProjectOptions(caseCode: string) {
  const response = await fetch("/api/aese/v1/world/plant-build/project-options", {
    method: "POST", headers: { "content-type": "application/json", ...iaosHeaders() }, body: JSON.stringify({ case_code: caseCode }),
  });
  if (!response.ok) throw new Error(`项目 Agent 准备方案 ${response.status}: ${await response.text()}`);
  return response.json() as Promise<ProjectPlanOptionSet>;
}

export async function submitFacilityProjectOption(caseCode: string, option: ProjectPlanOption, evidence: ProposalSet["evidence"]) {
  const response = await fetch("/api/aese/v1/world/plant-build/facility-projects", {
    method: "POST", headers: { "content-type": "application/json", ...iaosHeaders() }, body: JSON.stringify({ case_code: caseCode, option, evidence }),
  });
  if (!response.ok) throw new Error(`提交项目基线审批 ${response.status}: ${await response.text()}`);
  return response.json() as Promise<{ status: "waiting_approval"; plan_id: string; result: { result: { approval_request_id: string } } }>;
}

export async function loadFacilityProjects(caseCode: string, signal?: AbortSignal) {
  const response = await fetch(`/api/aese/v1/world/plant-build/facility-projects?case_code=${encodeURIComponent(caseCode)}`, { headers: iaosHeaders(), signal });
  if (!response.ok) throw new Error(`设施项目与 WBS ${response.status}: ${await response.text()}`);
  return response.json() as Promise<{ items: FacilityProjectItem[] }>;
}

export async function activateFacilityProject(caseCode: string, planID: string, approvalRequestID: string) {
  const response = await fetch("/api/aese/v1/world/plant-build/facility-projects/activate", {
    method: "POST", headers: { "content-type": "application/json", ...iaosHeaders() }, body: JSON.stringify({ case_code: caseCode, plan_id: planID, approval_request_id: approvalRequestID }),
  });
  if (!response.ok) throw new Error(`激活设施项目基线 ${response.status}: ${await response.text()}`);
  return response.json();
}

export async function submitFacilityProjectBaseline(caseCode: string, planID: string, expectedHash: string) {
  const response = await fetch("/api/aese/v1/world/plant-build/facility-projects/submit", {
    method: "POST", headers: { "content-type": "application/json", ...iaosHeaders() },
    body: JSON.stringify({ case_code: caseCode, plan_id: planID, expected_hash: expectedHash }),
  });
  if (!response.ok) throw new Error(`提交设施项目审批 ${response.status}: ${await response.text()}`);
  return response.json();
}

export async function loadContractAwards(caseCode: string, signal?: AbortSignal) {
  const response = await fetch(`/api/aese/v1/world/plant-build/contract-awards?case_code=${encodeURIComponent(caseCode)}`, { headers: iaosHeaders(), signal });
  if (!response.ok) throw new Error(`承包商寻源档案 ${response.status}: ${await response.text()}`);
  return response.json() as Promise<{ items: ContractAwardItem[] }>;
}

export async function issueContractRFQ(caseCode: string, packageCode: string, sourcingStrategy: ContractRFQ["sourcing_strategy"]) {
  const response = await fetch("/api/aese/v1/world/plant-build/contract-rfqs", {
    method: "POST", headers: { "content-type": "application/json", ...iaosHeaders() },
    body: JSON.stringify({ case_code: caseCode, package_code: packageCode, sourcing_strategy: sourcingStrategy }),
  });
  if (!response.ok) throw new Error(`发布工程采购邀请 ${response.status}: ${await response.text()}`);
  return response.json();
}

export async function confirmContractBids(caseCode: string, rfqID: string) {
  const response = await fetch("/api/aese/v1/world/plant-build/contract-bids/confirm", {
    method: "POST", headers: { "content-type": "application/json", ...iaosHeaders() },
    body: JSON.stringify({ case_code: caseCode, rfq_id: rfqID, action: "accept_bids" }),
  });
  if (!response.ok) throw new Error(`收取承包商正式投标 ${response.status}: ${await response.text()}`);
  return response.json();
}

export async function generateContractRecommendation(caseCode: string, rfqID: string) {
  const response = await fetch("/api/aese/v1/world/plant-build/contract-recommendations/agent", {
    method: "POST", headers: { "content-type": "application/json", ...iaosHeaders() }, body: JSON.stringify({ case_code: caseCode, rfq_id: rfqID }),
  });
  if (!response.ok) throw new Error(`工程采购 Agent 评审 ${response.status}: ${await response.text()}`);
  return response.json() as Promise<ContractRecommendationAdvice>;
}

export async function submitContractRecommendation(caseCode: string, rfqID: string, advice: ContractRecommendationAdvice) {
  const response = await fetch("/api/aese/v1/world/plant-build/contract-recommendations", {
    method: "POST", headers: { "content-type": "application/json", ...iaosHeaders() }, body: JSON.stringify({ case_code: caseCode, rfq_id: rfqID, advice }),
  });
  if (!response.ok) throw new Error(`提交合同授予审批 ${response.status}: ${await response.text()}`);
  return response.json();
}

export async function awardContract(caseCode: string, recommendationID: string, approvalRequestID: string) {
  const response = await fetch("/api/aese/v1/world/plant-build/contracts/award", {
    method: "POST", headers: { "content-type": "application/json", ...iaosHeaders() },
    body: JSON.stringify({ case_code: caseCode, recommendation_id: recommendationID, approval_request_id: approvalRequestID }),
  });
  if (!response.ok) throw new Error(`归档正式工程合同 ${response.status}: ${await response.text()}`);
  return response.json();
}

export async function loadConstructionMilestones(caseCode: string, signal?: AbortSignal) {
  const response = await fetch(`/api/aese/v1/world/plant-build/construction-milestones?case_code=${encodeURIComponent(caseCode)}`, { headers: iaosHeaders(), signal });
  if (!response.ok) throw new Error(`施工里程碑档案 ${response.status}: ${await response.text()}`);
  return response.json() as Promise<{ items: ConstructionMilestoneItem[] }>;
}
export async function startConstructionPackage(caseCode: string, contractID: string) {
  const response = await fetch("/api/aese/v1/world/plant-build/construction-packages", { method: "POST", headers: { "content-type": "application/json", ...iaosHeaders() }, body: JSON.stringify({ case_code: caseCode, contract_id: contractID }) });
  if (!response.ok) throw new Error(`启动施工包 ${response.status}: ${await response.text()}`);
  return response.json();
}
export async function confirmConstructionProgress(caseCode: string, executionID: string) {
  const response = await fetch("/api/aese/v1/world/plant-build/construction-observations/confirm", { method: "POST", headers: { "content-type": "application/json", ...iaosHeaders() }, body: JSON.stringify({ case_code: caseCode, execution_id: executionID, action: "advance_construction" }) });
  if (!response.ok) throw new Error(`推进施工现场 ${response.status}: ${await response.text()}`);
  return response.json();
}
export async function acceptConstructionMilestone(caseCode: string, executionID: string, observationID: string) {
  const response = await fetch("/api/aese/v1/world/plant-build/construction-milestones/accept", { method: "POST", headers: { "content-type": "application/json", ...iaosHeaders() }, body: JSON.stringify({ case_code: caseCode, execution_id: executionID, observation_id: observationID }) });
  if (!response.ok) throw new Error(`验收施工里程碑 ${response.status}: ${await response.text()}`);
  return response.json();
}
export async function loadPlantBuild(
  signal?: AbortSignal,
): Promise<PlantBuildTrace> {
  const r = await fetch("/api/aese/v1/world/plant-build", { signal });
  if (!r.ok) throw new Error(`Plant Build API ${r.status}`);
  const t = (await r.json()) as PlantBuildTrace;
  t.frames = (t.frames ?? []).map((f) => ({
    ...f,
    assessments: f.assessments ?? [],
    zones: f.zones ?? [],
    work_packages: f.work_packages ?? [],
    knowledge: f.knowledge ?? [],
    utilities: f.utilities ?? {},
  }));
  return t;
}

export async function loadPlantPlanningStatus(
  signal?: AbortSignal,
): Promise<PlantPlanningProviderStatus> {
  const response = await fetch(
    "/api/aese/v1/world/plant-build/planning-status",
    { signal },
  );
  if (!response.ok) throw new Error(`设施规划模型状态 ${response.status}`);
  return response.json() as Promise<PlantPlanningProviderStatus>;
}

const iaosHeaders = () => {
  const token = localStorage.getItem("iaos_token")?.trim();
  const tenant = (
    localStorage.getItem("aese_iaos_tenant_id") ??
    localStorage.getItem("iaos_tenant_id")
  )?.trim();
  if (!token || !tenant) throw new Error("缺少 IAOS 企业会话，请重新从企业入口进入");
  return { Authorization: `Bearer ${token}`, "X-IAOS-Tenant-Id": tenant };
};

export async function loadPlantFinancialConstraint(
  caseCode: string,
  signal?: AbortSignal,
): Promise<PlantFinancialConstraint> {
  if (!caseCode.trim()) throw new Error("当前企业会话缺少设立案编码");
  const response = await fetch(
    `/api/aese/v1/world/plant-build/financial-constraints?case_code=${encodeURIComponent(caseCode)}`,
    { headers: iaosHeaders(), signal },
  );
  if (!response.ok) {
    const detail = await response.text();
    throw new Error(`权威资金与预算快照 ${response.status}: ${detail}`);
  }
  return response.json() as Promise<PlantFinancialConstraint>;
}

export async function generatePlantRequirementOptions(caseCode: string) {
  const response = await fetch("/api/aese/v1/world/plant-build/requirement-options", {
    method: "POST", headers: { "content-type": "application/json", ...iaosHeaders() },
    body: JSON.stringify({ case_code: caseCode }),
  });
  if (!response.ok) throw new Error(`Agent 准备设施需求草案 ${response.status}: ${await response.text()}`);
  return response.json() as Promise<PlantRequirementOptionSet>;
}

export async function loadPlantRequirement(
  requirementID: string,
  signal?: AbortSignal,
): Promise<FacilityRequirement | null> {
  const response = await fetch(
    `/api/aese/v1/world/plant-build/requirements/${encodeURIComponent(requirementID)}`,
    { headers: iaosHeaders(), signal },
  );
  if (response.status === 404) return null;
  if (!response.ok) throw new Error(`设施需求历史 ${response.status}: ${await response.text()}`);
  return response.json() as Promise<FacilityRequirement>;
}

export async function loadPlantProposalSet(
  requirementID: string,
  signal?: AbortSignal,
): Promise<ProposalSet | null> {
  const response = await fetch(
    `/api/aese/v1/world/plant-build/proposals?requirement_id=${encodeURIComponent(requirementID)}`,
    { headers: iaosHeaders(), signal },
  );
  if (response.status === 404) return null;
  if (!response.ok) throw new Error(`设施候选历史 ${response.status}: ${await response.text()}`);
  return response.json() as Promise<ProposalSet>;
}

export async function generatePlantProposals(
  requirement: FacilityRequirement,
): Promise<PlantPlanningResult> {
  const workspace = localStorage.getItem("aese_genesis_workspace_id")?.trim();
  const response = await fetch("/api/aese/v1/world/plant-build/proposals", {
    method: "POST",
    headers: {
      "content-type": "application/json",
      ...iaosHeaders(),
      ...(workspace ? { "X-Genesis-Workspace-Id": workspace } : {}),
    },
    body: JSON.stringify(requirement),
  });
  if (!response.ok) {
    let message = `设施规划 Agent ${response.status}`;
    try {
      const payload = (await response.json()) as { error?: string };
      message = payload.error || message;
    } catch {
      // Non-JSON proxy failures still retain the status-oriented message.
    }
    throw new Error(message);
  }
  return response.json() as Promise<PlantPlanningResult>;
}

export async function submitManualPlantProposal(input: {
  requirement_id: string;
  proposal_set_id: string;
  expected_revision: number;
  proposal: Omit<SiteOptionProposal, "proposal_id" | "status">;
}) {
  const response = await fetch("/api/aese/v1/world/plant-build/proposals/manual", {
    method: "POST",
    headers: { "content-type": "application/json", ...iaosHeaders() },
    body: JSON.stringify(input),
  });
  if (!response.ok) throw new Error(`提交人工候选 ${response.status}: ${await response.text()}`);
  return response.json() as Promise<{ status: "committed"; proposal_set: ProposalSet; manual_proposal_id: string }>;
}
