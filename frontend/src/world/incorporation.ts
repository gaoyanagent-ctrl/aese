export type Money = { value: string; currency: "CNY"; scale: 2 };
export type Appointment = {
  position: string;
  assignee: string;
  resolution?: string;
  status: string;
  accepted_at?: string;
};
export type IncorporationFrame = {
  step: number;
  phase: string;
  sim_time: string;
  title: string;
  causation_id: string;
  legal_entity_status: string;
  registration_status: string;
  investor: { code: string; owner: string; status: string; balance: Money };
  company: { code: string; owner: string; status: string; balance: Money };
  capital_committed: Money;
  capital_paid: Money;
  governance: {
    ceo: Appointment;
    cfo: Appointment;
    project_director: Appointment;
    mandate_active: boolean;
  };
  budget: { code: string; status: string; amount: Money; owner: string };
  knowledge: Array<{
    actor: string;
    fact: string;
    observed_at: string;
    source: string;
    confidence: string;
    visibility: string;
  }>;
  iaos_cursor: number;
  plant_project_eligible: boolean;
};
export type IncorporationTrace = {
  schema_version: string;
  campaign: "incorporation";
  world_run_id: string;
  timezone: "Asia/Shanghai";
  policy_version: string;
  frames: IncorporationFrame[];
};
export async function loadIncorporation(
  signal?: AbortSignal,
): Promise<IncorporationTrace> {
  const response = await fetch("/api/aese/v1/world/incorporation", { signal });
  if (!response.ok) throw new Error(`Incorporation API ${response.status}`);
  const trace = (await response.json()) as IncorporationTrace;
  trace.frames = (trace.frames ?? []).map((frame) => ({
    ...frame,
    knowledge: frame.knowledge ?? [],
  }));
  return trace;
}
