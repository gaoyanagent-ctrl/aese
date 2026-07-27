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
    outbox: unknown[];
    world_exchanges: Array<Record<string, unknown>>;
    process_runs: Array<Record<string, unknown>>;
    decisions: unknown[];
    runtime_artifact: Record<string, unknown>;
    lineage: Record<string, unknown>;
  };
  iaos_lifecycle_warning?: {
    code: "case_not_found";
    message: string;
    requested_case: string;
    available_cases: string[];
  };
};

export type IncorporationFactSource =
  | "iaos_committed"
  | "world_baseline"
  | "pending";

export type DisplayFact = {
  value: string;
  source: IncorporationFactSource;
  sourceLabel: string;
};

type LifecycleState = {
  proposed_company_name?: string;
  legal_entity_code?: string;
  bank_account_code?: string;
  commitment_minor?: number;
  contribution_minor?: number;
  budget_authorized_minor?: number;
};

const displayFact = (
  value: string,
  source: IncorporationFactSource,
  sourceLabel: string,
): DisplayFact => ({ value, source, sourceLabel });

const minorToMajor = (minor: number): string => (minor / 100).toFixed(2);

/**
 * IAOS committed state wins. Pack values remain visible only as explicitly
 * labelled simulation assumptions and never masquerade as active-case facts.
 */
export function resolveIncorporationFacts(
  frame: IncorporationFrame,
  lifecycle?: IncorporationTrace["iaos_lifecycle"],
) {
  const state = (lifecycle?.state ?? {}) as LifecycleState;
  const proposedName = state.proposed_company_name?.trim();
  const legalEntityCode = state.legal_entity_code?.trim();
  const bankAccountCode = state.bank_account_code?.trim();
  const hasCommittedCase = Boolean(lifecycle?.case_code);

  return {
    legalEntity: legalEntityCode
      ? displayFact(legalEntityCode, "iaos_committed", "IAOS 登记提交事实")
      : proposedName
        ? displayFact(
            `${proposedName}（待登记编码）`,
            "pending",
            "IAOS 设立案输入",
          )
        : displayFact(
            `${frame.company.owner}（场景目标）`,
            "world_baseline",
            "AESE 华辰场景假设",
          ),
    companyAccount: bankAccountCode
      ? displayFact(bankAccountCode, "iaos_committed", "IAOS 开户提交事实")
      : displayFact(
          frame.company.code,
          "world_baseline",
          "AESE 华辰场景假设",
        ),
    capitalCommitted:
      hasCommittedCase && Number(state.commitment_minor) > 0
        ? displayFact(
            minorToMajor(Number(state.commitment_minor)),
            "iaos_committed",
            "IAOS 设立案提交事实",
          )
        : displayFact(
            frame.capital_committed.value,
            "world_baseline",
            "AESE 华辰场景假设",
          ),
    capitalPaid:
      hasCommittedCase && Number(state.contribution_minor) > 0
        ? displayFact(
            minorToMajor(Number(state.contribution_minor)),
            "iaos_committed",
            "IAOS 资本核验提交事实",
          )
        : displayFact(
            frame.capital_paid.value,
            "world_baseline",
            "AESE 华辰场景假设",
          ),
    budget:
      hasCommittedCase && Number(state.budget_authorized_minor) > 0
        ? displayFact(
            minorToMajor(Number(state.budget_authorized_minor)),
            "iaos_committed",
            "IAOS 预算审批提交事实",
          )
        : displayFact(
            frame.budget.amount.value,
            "world_baseline",
            "AESE 华辰场景假设",
          ),
  };
}

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
      const headers = { Authorization: `Bearer ${token}`, "X-Tenant-ID": tenant };
      const lifecycle = await fetch(`${base}/api/v1/incorporations/${encodeURIComponent(caseCode)}/trace`, {
        signal,
        headers,
      });
      if (lifecycle.status === 404) {
        const recent = await fetch(`${base}/api/v1/incorporations/recent`, {
          signal,
          headers,
        });
        const body = recent.ok
          ? await recent.json() as { items?: Array<{ case_code?: string }> }
          : {};
        trace.iaos_lifecycle_warning = {
          code: "case_not_found",
          message: `IAOS 中不存在设立案 ${caseCode}，已保留 AESE World 本地基线。`,
          requested_case: caseCode,
          available_cases: (body.items ?? [])
            .map((item) => item.case_code?.trim() ?? "")
            .filter(Boolean),
        };
        return trace;
      }
      if (!lifecycle.ok) throw new Error(`IAOS lifecycle API ${lifecycle.status}`);
      trace.iaos_lifecycle = await lifecycle.json();
    }
  }
  return trace;
}

export async function submitIncorporationObservation(
  caseCode: string,
  payloadType: string,
  result: string,
): Promise<void> {
  const params = new URLSearchParams(window.location.hash.split("?")[1] ?? "");
  const tenant = params.get("tenant") ?? "tenant-hctm-genesis";
  const token = localStorage.getItem("iaos_token");
  if (!token) throw new Error("缺少 IAOS 登录凭据，请从 IAOS 重新打开 AESE");
  const correlation = params.get("correlation") || `corr-${caseCode}`;
  const nonce = `${caseCode}-${payloadType}-${Date.now()}`;
  const response = await fetch(
    `${resolveIaosLifecycleBase()}/api/v1/world-bridge/observations`,
    {
      method: "POST",
      headers: {
        Authorization: `Bearer ${token}`,
        "Content-Type": "application/json",
        "X-Tenant-ID": tenant,
      },
      body: JSON.stringify({
        schema_version: "1.0",
        message_id: `obs-${nonce}`,
        kind: "observation",
        tenant_id: tenant,
        world_pack_key: "hctm-genesis",
        world_pack_version: "1.0.0",
        world_run_id: correlation,
        branch_id: "main",
        sim_occurred_at: new Date().toISOString(),
        correlation_id: correlation,
        causation_id: `external-review-${caseCode}`,
        idempotency_key: `obs-${nonce}`,
        producer: {
          system: "aese",
          component: "world-runtime",
          version: "1.0.0",
        },
        subject_ref: {
          namespace: "hctm",
          type: "incorporation_case",
          code: caseCode,
        },
        payload_type: payloadType,
        payload: { result },
      }),
    },
  );
  if (!response.ok) {
    const detail = await response.text();
    throw new Error(`World Observation ${response.status}: ${detail}`);
  }
}
