import type {
  FounderIntent,
  FounderIntentRequest,
  GameProjection,
  GenesisWorkspaceResult,
  NamingProposal,
} from "./types";
import {
  resolveIaosLifecycleBase,
  submitIncorporationObservation,
} from "../world/incorporation";
export class GameApiError extends Error {
  constructor(
    public status: number,
    message: string,
  ) {
    super(message);
    this.name = "GameApiError";
  }
}
export async function loadGameProjection(
  caseCode: string,
  frame: number,
  signal?: AbortSignal,
) {
  const request = () => {
    const token = localStorage.getItem("iaos_token") ?? "";
    const tenant =
      localStorage.getItem("aese_iaos_tenant_id") ??
      localStorage.getItem("iaos_tenant_id") ??
      "";
    return fetch(
      `/api/aese/v1/game/incorporation/${encodeURIComponent(caseCode)}/projection?frame=${frame}`,
      {
        signal,
        headers: {
          ...(token ? { Authorization: `Bearer ${token}` } : {}),
          ...(tenant ? { "X-IAOS-Tenant-Id": tenant } : {}),
        },
      },
    );
  };
  let response = await request();
  if (
    (response.status === 401 || response.status === 404) &&
    localStorage.getItem("aese_genesis_workspace_id") &&
    genesisPlayerToken()
  ) {
    await refreshGenesisWorkspaceSession();
    response = await request();
  }
  if (!response.ok)
    throw new GameApiError(
      response.status,
      `GameProjection API ${response.status}: ${await response.text()}`,
    );
  return response.json() as Promise<GameProjection>;
}
async function postJSON<T>(path: string, body: unknown) {
  const token = localStorage.getItem("iaos_token") ?? "";
  const tenant =
    localStorage.getItem("aese_iaos_tenant_id") ??
    localStorage.getItem("iaos_tenant_id") ??
    "";
  const workspace = localStorage.getItem("aese_genesis_workspace_id") ?? "";
  const response = await fetch(path, {
    method: "POST",
    headers: {
      "content-type": "application/json",
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
      ...(tenant ? { "X-IAOS-Tenant-Id": tenant } : {}),
      ...(workspace ? { "X-Genesis-Workspace-Id": workspace } : {}),
    },
    body: JSON.stringify(body),
  });
  if (!response.ok) throw new Error(`Creative API ${response.status}`);
  return response.json() as Promise<T>;
}
export const analyzeFounderIntent = (request: FounderIntentRequest) =>
  postJSON<FounderIntent>("/api/aese/v1/game/creative/intent", request);
export type CreativeProviderStatus = {
  state: "connected" | "degraded" | "fallback" | "not_configured";
  provider: string;
  model: string;
  base_url_host?: string;
  prompt_version: string;
};
export async function loadCreativeProviderStatus() {
  const response = await fetch("/api/aese/v1/game/creative/status");
  if (!response.ok)
    throw new Error(`Creative Provider 状态加载失败 ${response.status}`);
  return response.json() as Promise<CreativeProviderStatus>;
}
export async function generateCompanyNames(intent: FounderIntent) {
  const result = await postJSON<{
    status: "candidate_only";
    proposals: NamingProposal[];
  }>("/api/aese/v1/game/creative/names", intent);
  return result.proposals;
}
export type GenesisPlayerProfile = {
  subject_id: string;
  username: string;
  display_name: string;
};
export type GenesisPlayerSession = {
  status: string;
  token: string;
  expires_at: string;
  player: GenesisPlayerProfile;
};
const genesisPlayerToken = () =>
  sessionStorage.getItem("aese_genesis_player_token")?.trim() ?? "";
async function genesisAuthFetch(
  path: string,
  init: RequestInit,
  action: string,
) {
  const response = await fetch(path, init);
  if (response.ok) return response;
  let message = `${action}失败`;
  try {
    const payload = (await response.json()) as {
      error?: string;
      code?: string;
    };
    const translations: Record<string, string> = {
      invalid_credentials: "用户名或密码错误",
      username_exists: "该用户名已经注册，请直接登录",
      account_locked: "连续登录失败次数过多，账号已临时锁定 15 分钟",
      account_disabled: "账号已停用，请联系系统管理员",
      invalid_registration: "注册信息不符合安全要求",
    };
    message = translations[payload.code ?? ""] ?? payload.error ?? message;
  } catch {
    /* response is not JSON */
  }
  throw new GameApiError(response.status, message);
}
function persistGenesisPlayerSession(session: GenesisPlayerSession) {
  sessionStorage.setItem("aese_genesis_player_token", session.token);
  localStorage.setItem("aese_genesis_username", session.player.username);
  localStorage.setItem(
    "aese_genesis_player_display_name",
    session.player.display_name,
  );
  localStorage.setItem("aese_genesis_player_id", session.player.subject_id);
  return session;
}
export async function loginGenesisPlayer(input: {
  username: string;
  password: string;
}) {
  const response = await genesisAuthFetch(
    "/api/aese/v1/auth/login",
    {
      method: "POST",
      headers: { "content-type": "application/json" },
      body: JSON.stringify(input),
    },
    "登录",
  );
  return persistGenesisPlayerSession(
    (await response.json()) as GenesisPlayerSession,
  );
}
export async function registerGenesisPlayer(input: {
  username: string;
  password: string;
  display_name: string;
}) {
  const response = await genesisAuthFetch(
    "/api/aese/v1/auth/register",
    {
      method: "POST",
      headers: { "content-type": "application/json" },
      body: JSON.stringify(input),
    },
    "注册",
  );
  return persistGenesisPlayerSession(
    (await response.json()) as GenesisPlayerSession,
  );
}
async function genesisWorkspaceFetch(
  path: string,
  init: RequestInit = {},
  action = "企业操作",
) {
  const token = genesisPlayerToken();
  if (!token) throw new Error(`创始人登录已过期，无法${action}，请重新登录`);
  const response = await fetch(path, {
    ...init,
    headers: {
      ...(init.headers as Record<string, string> | undefined),
      Authorization: `Bearer ${token}`,
    },
  });
  if (response.status === 401) {
    sessionStorage.removeItem("aese_genesis_player_token");
    throw new Error(`创始人登录已过期，无法${action}，请重新登录`);
  }
  return response;
}
export const currentGenesisUsername = () =>
  genesisPlayerToken()
    ? (localStorage.getItem("aese_genesis_username") ?? "")
    : "";
export function signOutGenesisPlayer() {
  localStorage.removeItem("aese_genesis_username");
  localStorage.removeItem("aese_genesis_player_id");
  localStorage.removeItem("iaos_token");
  localStorage.removeItem("aese_iaos_tenant_id");
  localStorage.removeItem("iaos_tenant_id");
  localStorage.removeItem("aese_genesis_workspace_id");
  localStorage.removeItem("aese_genesis_case_code");
  localStorage.removeItem("aese_genesis_player_display_name");
  localStorage.removeItem("aese_genesis_player_token");
  sessionStorage.removeItem("aese_genesis_player_token");
}
export async function listGenesisWorkspaces() {
  const response = await genesisWorkspaceFetch(
    "/api/aese/v1/genesis/workspaces",
    {},
    "加载企业列表",
  );
  if (!response.ok)
    throw new Error(
      `企业列表加载失败 ${response.status}: ${await response.text()}`,
    );
  const result = (await response.json()) as { items: GenesisWorkspaceResult[] };
  return result.items;
}
export async function createGenesisWorkspace(input: {
  display_name: string;
  idempotency_key: string;
  template_key: string;
  region: string;
  timezone: string;
  realism_level: "standard" | "strict";
  data_retention_confirmed: boolean;
}) {
  const response = await genesisWorkspaceFetch(
    "/api/aese/v1/genesis/workspaces",
    {
      method: "POST",
      headers: { "content-type": "application/json" },
      body: JSON.stringify(input),
    },
    "创建企业",
  );
  if (!response.ok)
    throw new Error(
      `创业空间创建失败 ${response.status}: ${await response.text()}`,
    );
  const result = (await response.json()) as GenesisWorkspaceResult;
  localStorage.setItem("iaos_token", result.tenant_token);
  localStorage.setItem("aese_iaos_tenant_id", result.tenant_id);
  localStorage.setItem("iaos_tenant_id", result.tenant_id);
  localStorage.setItem("aese_genesis_workspace_id", result.workspace_id);
  localStorage.setItem("aese_genesis_case_code", result.case_code);
  return result;
}
export async function refreshGenesisWorkspaceSession() {
  const workspace = localStorage.getItem("aese_genesis_workspace_id");
  if (!workspace) throw new Error("缺少创业空间标识，请从游戏主页重新进入");
  const response = await genesisWorkspaceFetch(
    `/api/aese/v1/genesis/workspaces/${encodeURIComponent(workspace)}/session`,
    { method: "POST" },
    "刷新 Founder 会话",
  );
  if (!response.ok)
    throw new Error(
      `Founder 会话刷新失败 ${response.status}: ${await response.text()}`,
  );
  const result = (await response.json()) as GenesisWorkspaceResult;
  localStorage.setItem("iaos_token", result.tenant_token);
  localStorage.setItem("aese_iaos_tenant_id", result.tenant_id);
  localStorage.setItem("iaos_tenant_id", result.tenant_id);
  localStorage.setItem("aese_genesis_workspace_id", result.workspace_id);
  localStorage.setItem("aese_genesis_case_code", result.case_code);
  return result.tenant_token;
}
export async function resumeGenesisWorkspace(
  workspace: GenesisWorkspaceResult,
) {
  const response = await genesisWorkspaceFetch(
    `/api/aese/v1/genesis/workspaces/${encodeURIComponent(workspace.workspace_id)}/session`,
    { method: "POST" },
    "进入企业",
  );
  if (!response.ok)
    throw new Error(
      `进入企业失败 ${response.status}: ${await response.text()}`,
    );
  const result = (await response.json()) as GenesisWorkspaceResult;
  localStorage.setItem("iaos_token", result.tenant_token);
  localStorage.setItem("aese_iaos_tenant_id", result.tenant_id);
  localStorage.setItem("iaos_tenant_id", result.tenant_id);
  localStorage.setItem("aese_genesis_workspace_id", result.workspace_id);
  localStorage.setItem("aese_genesis_case_code", result.case_code);
  return result;
}
export async function createIncorporationCase(input: {
  case_code: string;
  case_name: string;
  proposed_company_name: string;
  registered_address: string;
  business_scope: string;
}) {
  let token = localStorage.getItem("iaos_token");
  const tenant =
    localStorage.getItem("aese_iaos_tenant_id") ??
    localStorage.getItem("iaos_tenant_id") ??
    "tenant-hctm-genesis";
  if (!token) throw new Error("缺少 IAOS 登录凭据，请从 IAOS 重新进入企业创生");
  const submit = () =>
    fetch(`${resolveIaosLifecycleBase()}/api/v1/incorporations/cases`, {
      method: "POST",
      headers: {
        Authorization: `Bearer ${token}`,
        "X-Tenant-ID": tenant,
        "content-type": "application/json",
      },
      body: JSON.stringify(input),
    });
  let response = await submit();
  if (response.status === 422) {
    const detail = await response.clone().text();
    if (detail.includes("authenticated principal does not match")) {
      token = await refreshGenesisWorkspaceSession();
      response = await submit();
    }
  }
  if (!response.ok)
    throw new Error(
      `IAOS incorporation.case.open ${response.status}: ${await response.text()}`,
    );
  return response.json();
}
const iaosHeaders = () => {
  const token = localStorage.getItem("iaos_token");
  const tenant =
    localStorage.getItem("aese_iaos_tenant_id") ??
    localStorage.getItem("iaos_tenant_id") ??
    "tenant-hctm-genesis";
  if (!token) throw new Error("缺少 IAOS 登录凭据，请从 IAOS 重新进入企业创生");
  return {
    Authorization: `Bearer ${token}`,
    "X-Tenant-ID": tenant,
    "content-type": "application/json",
  };
};
async function iaosPost(path: string, body: unknown) {
  const submit = () =>
    fetch(`${resolveIaosLifecycleBase()}${path}`, {
      method: "POST",
      headers: iaosHeaders(),
      body: JSON.stringify(body),
    });
  let response = await submit();
  if (response.status === 422) {
    const detail = await response.clone().text();
    if (
      detail.includes("effective runtime artifact stale") &&
      localStorage.getItem("aese_genesis_workspace_id") &&
      genesisPlayerToken()
    ) {
      await refreshGenesisWorkspaceSession();
      response = await submit();
    }
  }
  if (!response.ok)
    throw new Error(`IAOS ${response.status}: ${await response.text()}`);
  return response.json();
}
export type WorkItemInput = {
  amount_minor?: number;
  currency?: string;
  business_note?: string;
  resolution_objective?: string;
  key_proposals?: string;
  risk_notes?: string;
  accounting_book_name?: string;
  fiscal_year?: number;
  accounting_periods?: Array<{
    code: string;
    starts_on: string;
    ends_on: string;
  }>;
};
export async function executeWorkItem(
  caseCode: string,
  sequence: number,
  input: WorkItemInput = {},
) {
  return iaosPost(
    `/api/v1/incorporations/${encodeURIComponent(caseCode)}/work-items/${sequence}/execute`,
    input,
  );
}
export async function dispatchAgentWorkItem(
  caseCode: string,
  sequence: number,
  input: WorkItemInput,
) {
  return iaosPost(
    `/api/v1/incorporations/${encodeURIComponent(caseCode)}/work-items/${sequence}/dispatch-agent`,
    input,
  );
}
export async function approveAndExecuteWorkItem(
  caseCode: string,
  sequence: number,
  capability: string,
  gate: string,
  input: WorkItemInput = {},
) {
  const materialSuffix =
    input.amount_minor !== undefined ? `-${input.amount_minor}` : "";
  const correlation = `corr-game-${caseCode}-${sequence}${materialSuffix}`;
  const submitted = (await iaosPost(
    `/api/v1/incorporations/${encodeURIComponent(caseCode)}/gates/${encodeURIComponent(gate)}/submit`,
    {
      capability,
      correlation_id: correlation,
      intent: { case_code: caseCode, ...input },
    },
  )) as { approval?: { id?: string; status?: string } };
  const approvalId = submitted.approval?.id;
  if (!approvalId) throw new Error(`IAOS ${gate} 未返回审批请求`);
  if (submitted.approval?.status !== "approved") {
    await iaosPost(
      `/api/v1/approvals/${encodeURIComponent(approvalId)}/approve`,
      {
        note: `Enterprise Genesis 玩家批准 ${capability}`,
      },
    );
  }
  const executeInput: WorkItemInput & { correlation_id: string } = {
    correlation_id: correlation,
  };
  if (input.amount_minor !== undefined)
    executeInput.amount_minor = input.amount_minor;
  if (input.currency !== undefined) executeInput.currency = input.currency;
  return iaosPost(
    `/api/v1/incorporations/${encodeURIComponent(caseCode)}/work-items/${sequence}/execute`,
    executeInput,
  );
}
export async function resolveWorldWorkItem(
  caseCode: string,
  sequence: number,
  capability: string,
  correlation: string,
  result: string,
) {
  const contracts: Record<string, string> = {
    "registration.observation.commit": "registration.decision.observed.v1",
    "bank.account.observation.commit": "bank.account.decision.observed.v1",
    "executive.appointment.acceptance.commit":
      "appointment.acceptance.observed.v1",
  };
  const payloadType = contracts[capability];
  if (!payloadType)
    throw new Error(`没有 World Observation 合同: ${capability}`);
  await submitIncorporationObservation(
    caseCode,
    payloadType,
    result,
    correlation,
  );
  await iaosPost(
    `/api/v1/incorporations/${encodeURIComponent(caseCode)}/work-items/${sequence}/execute`,
    { correlation_id: correlation },
  );
}
