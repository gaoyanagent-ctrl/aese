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
  iaos_lifecycle?: {
    case_code: string;
    state: Record<string, unknown>;
    journal: unknown[];
    approvals: unknown[];
    world_exchanges: Array<Record<string, unknown>>;
    process_runs: Array<Record<string, unknown>>;
    decisions: unknown[];
    runtime_artifact: Record<string, unknown>;
    lineage: Record<string, unknown>;
  };
};

export function resolveIaosLifecycleBase(): string {
  const fallback = `http://${window.location.hostname || "127.0.0.1"}:8082`;
  const configured = localStorage.getItem("aese_iaos_base_url")?.trim();
  if (!configured) return fallback;
  try {
    const candidate = new URL(configured, window.location.origin);
    // A historical integration default stored the AESE/Vite origin here.
    // That origin cannot serve IAOS /api/v1 routes in the standalone stack.
    if (candidate.origin === window.location.origin) return fallback;
    return candidate.toString().replace(/\/$/, "");
  } catch {
    return fallback;
  }
}

function acceptLifecycleToken(params: URLSearchParams): string | null {
  const handedOff = params.get("auth_token")?.trim();
  if (!handedOff) return localStorage.getItem("iaos_token");
  localStorage.setItem("iaos_token", handedOff);
  params.delete("auth_token");
  const route = window.location.hash.split("?")[0] || "#world-incorporation";
  const query = params.toString();
  window.history.replaceState(
    null,
    "",
    `${window.location.pathname}${window.location.search}${route}${query ? `?${query}` : ""}`,
  );
  return handedOff;
}

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
  if (typeof window !== "undefined") {
    const params = new URLSearchParams(window.location.hash.split("?")[1] ?? "");
    const caseCode = params.get("case");
    const tenant = params.get("tenant") ?? localStorage.getItem("aese_iaos_tenant_id") ?? "tenant-hctm-genesis";
    const token = acceptLifecycleToken(params);
    const base = resolveIaosLifecycleBase();
    if (caseCode && token) {
      const lifecycle = await fetch(`${base}/api/v1/incorporations/${encodeURIComponent(caseCode)}/trace`, {
        signal,
        headers: { Authorization: `Bearer ${token}`, "X-Tenant-ID": tenant },
      });
      if (!lifecycle.ok) throw new Error(`IAOS lifecycle API ${lifecycle.status}`);
      trace.iaos_lifecycle = await lifecycle.json();
    }
  }
  return trace;
}
