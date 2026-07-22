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
