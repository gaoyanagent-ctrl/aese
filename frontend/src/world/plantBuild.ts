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

export async function submitSiteInvestigationObservation(caseCode: string, worldRunID: string, observation: SiteInvestigationObservation) {
  const response = await fetch("/api/aese/v1/world/plant-build/observations", {
    method: "POST", headers: { "content-type": "application/json", ...iaosHeaders() },
    body: JSON.stringify({ case_code: caseCode, world_run_id: worldRunID, observation }),
  });
  if (!response.ok) throw new Error(`提交外部调研事实 ${response.status}: ${await response.text()}`);
  return response.json() as Promise<{ status: "committed"; world_message_id: string; observation: SiteInvestigationObservation }>;
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
