export type CapabilityTrace = {
  campaign: string;
  options: Array<{
    code: string;
    mode: string;
    feasible: boolean;
    failures: string[];
  }>;
  frames: Array<{
    step: number;
    phase: string;
    sim_time: string;
    title: string;
    cash: { value: string };
    committed: { value: string };
    paid: { value: string };
    equipment: Array<{ code: string; status: string; zone: string }>;
    workers: Array<{
      code: string;
      position: string;
      status: string;
      skills: string[];
    }>;
    knowledge: Array<{ actor: string; fact: string }>;
    world_progress: number;
    iaos_progress: number;
    discrepancy: string;
    gate: Record<string, boolean>;
    industrialization_eligible: boolean;
  }>;
};
export async function loadCapabilityBuild(
  signal?: AbortSignal,
): Promise<CapabilityTrace> {
  const r = await fetch("/api/aese/v1/world/capability-build", { signal });
  if (!r.ok) throw new Error(`Capability Build API ${r.status}`);
  const t = (await r.json()) as CapabilityTrace;
  t.frames = (t.frames ?? []).map((f) => ({
    ...f,
    equipment: f.equipment ?? [],
    workers: f.workers ?? [],
    knowledge: f.knowledge ?? [],
    gate: f.gate ?? {},
  }));
  return t;
}
