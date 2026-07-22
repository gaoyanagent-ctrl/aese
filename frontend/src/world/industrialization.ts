export type IndustrializationTrace = {
  frames: Array<{
    step: number;
    phase: string;
    sim_time: string;
    title: string;
    world_progress: number;
    iaos_progress: number;
    discrepancy: string;
    ppap_status: string;
    compatibility: string;
    serial_production_eligible: boolean;
    releases: Array<{
      code: string;
      revision: string;
      status: string;
      hash: string;
    }>;
    trials: Array<{
      code: string;
      revision: string;
      Cpk: string;
      Yield: string;
      leak_failures: number;
    }>;
    knowledge: unknown[];
    apqp_gates: Record<string, boolean>;
  }>;
};
export async function loadIndustrialization(
  signal?: AbortSignal,
): Promise<IndustrializationTrace> {
  const r = await fetch("/api/aese/v1/world/industrialization", { signal });
  if (!r.ok) throw new Error(`Industrialization API ${r.status}`);
  const t = (await r.json()) as IndustrializationTrace;
  t.frames = (t.frames ?? []).map((f) => ({
    ...f,
    releases: f.releases ?? [],
    trials: f.trials ?? [],
    knowledge: f.knowledge ?? [],
    apqp_gates: f.apqp_gates ?? {},
  }));
  return t;
}
