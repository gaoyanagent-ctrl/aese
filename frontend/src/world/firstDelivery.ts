export type FirstDeliveryTrace = {
  frames: Array<{
    step: number;
    phase: string;
    sim_time: string;
    title: string;
    demand: number;
    inventory: number;
    shipped: number;
    accepted: number;
    world_progress: number;
    iaos_progress: number;
    discrepancy: string;
    shipments: Array<{ code: string; quantity: number; accepted: number }>;
    cash: { value: string };
    invoice_gross: { value: string };
    ar: { value: string };
    collected: { value: string };
    revenue: { value: string };
    actual_cost: { value: string };
    gross_margin: { value: string };
    first_commercial_cycle_closed: boolean;
  }>;
};
export async function loadFirstDelivery(
  signal?: AbortSignal,
): Promise<FirstDeliveryTrace> {
  const r = await fetch("/api/aese/v1/world/first-delivery", { signal });
  if (!r.ok) throw new Error(`First Delivery API ${r.status}`);
  const t = (await r.json()) as FirstDeliveryTrace;
  t.frames = (t.frames ?? []).map((f) => ({
    ...f,
    shipments: f.shipments ?? [],
  }));
  return t;
}
