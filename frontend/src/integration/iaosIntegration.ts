import type { IaosScenarioSnapshot } from '../scenario/types';

export const IAOS_CONNECTION_STORAGE = {
  token: 'iaos_token',
  tenantId: 'aese_iaos_tenant_id',
  baseUrl: 'aese_iaos_base_url',
  uiUrl: 'aese_iaos_ui_url',
  orchestratorUrl: 'aese_aese_orchestrator_url',
  activeRunId: 'aese_aese_active_run_id',
  activeRunContext: 'aese_aese_active_run_context',
} as const;

export interface IntegrationConfiguration {
  tenantId: string;
  packKey: string;
  scenarioKey: string;
  baseUrl: string;
  uiUrl: string;
  orchestratorUrl?: string;
}

export interface IntegrationCheckResult {
  snapshot: IaosScenarioSnapshot;
  tenantName: string;
  counts: Record<'sales_order' | 'work_order' | 'inventory' | 'equipment', number>;
  eventChannelAvailable: boolean;
}

export interface OrchestrationPlan {
  pack_key: string;
  pack_version: string;
  scenario_key: string;
  correlation_id: string;
  total_events: number;
  stages: {
    stage: string;
    event_ids: string[];
    event_types: string[];
    event_count: number;
    action_hints: string[];
  }[];
  act_count: number;
  allowable_run_actions: string[];
  plan_hash: string;
}

export interface OrchestrationRunOutcome {
  action?: string;
  actor?: string;
  summary?: Record<string, unknown>;
  plan_hash?: string;
  initialize_dry_run?: Record<string, unknown>;
  operation_ref?: string;
  run_id?: string;
  correlation?: string;
  applied?: boolean;
  run_version?: string;
  reset_confirmation_token?: string;
  confirmation_expires_at?: string;
  apply?: boolean;
  warnings?: unknown[];
  entity_contracts?: unknown[];
  impacts?: unknown[];
  event_count?: number;
}

export interface OrchestrationRun {
  run_id: string;
  run_version: string;
  pack_key: string;
  pack_version: string;
  scenario_key: string;
  plan_hash: string;
  status: string;
  current_act: number;
  total_acts: number;
  cursor: number;
  tenant: string;
  target: string;
  created_at: string;
  updated_at: string;
  allowed_actions: string[];
  last_error?: string;
  retryable?: boolean;
  outcome?: unknown;
  plan?: unknown;
  reset_confirmation_required?: boolean;
}

export interface OrchestrationRunContext {
  runId: string;
  runVersion: string;
  planHash: string;
  status: string;
  currentAct: number;
  totalActs: number;
  cursor: number;
  tenantId: string;
  packKey: string;
  scenarioKey: string;
  target: string;
  updatedAt: string;
}

export interface OrchestrationActionResponse {
  run: OrchestrationRun;
  action: string;
}

export interface OrchestrationApiErrorDetails {
  status: number;
  code: string;
  message: string;
  retryable: boolean;
  runId?: string;
  runVersion?: string;
  requiredPermission?: string;
}

export class OrchestrationError extends Error {
  status: number;
  code: string;
  retryable: boolean;
  runId: string;
  runVersion: string;
  requiredPermission: string;

  constructor(payload: OrchestrationApiErrorDetails) {
    super(payload.message);
    this.name = 'OrchestrationError';
    this.status = payload.status;
    this.code = payload.code;
    this.retryable = payload.retryable;
    this.runId = payload.runId ?? '';
    this.runVersion = payload.runVersion ?? '';
    this.requiredPermission = payload.requiredPermission ?? '';
  }
}

function normalizedBaseUrl(value: string) {
  return value.trim().replace(/\/$/, '');
}

function defaultIaosApiBaseUrl(): string {
  if (typeof window === 'undefined') {
    return 'http://127.0.0.1:8082';
  }
  return `${window.location.origin}/api`;
}

function defaultOrchestratorBaseUrl(): string {
  if (typeof window === 'undefined') {
    return 'http://127.0.0.1:8090';
  }
  const host = window.location.hostname || '127.0.0.1';
  return `${window.location.protocol}//${host}:8090`;
}

function apiUrl(baseUrl: string, path: string) {
  if (baseUrl === '') {
    return `${path}`;
  }
  return `${normalizedBaseUrl(baseUrl)}${path}`;
}

function normalizeBaseUrl(value: string | undefined, fallback: string) {
  const trimmed = value?.trim() ?? '';
  if (trimmed === '') {
    return fallback;
  }
  if (/^[a-zA-Z][a-zA-Z0-9+.-]*:\/\//.test(trimmed)) {
    return normalizedBaseUrl(trimmed);
  }
  if (trimmed.startsWith('/')) {
    if (typeof window === 'undefined') {
      return fallback;
    }
    return `${window.location.origin}${trimmed}`;
  }
  return normalizedBaseUrl(`http://${trimmed}`);
}

export function resolveIaosApiBaseUrl(input?: string): string {
  return normalizeBaseUrl(input, defaultIaosApiBaseUrl());
}

export function resolveOrchestratorBaseUrl(input?: string): string {
  return normalizeBaseUrl(input, defaultOrchestratorBaseUrl());
}

function orchestrationApiUrl(baseUrl: string, path: string) {
  return apiUrl(normalizedBaseUrl(baseUrl), `/api/aese/v1${path.startsWith('/') ? path : `/${path}`}`);
}

function defaultOrchestratorUrl(): string {
  if (typeof window === 'undefined') {
    return 'http://127.0.0.1:8090';
  }
  const host = window.location.hostname || '127.0.0.1';
  return `http://${host}:8090`;
}

function getStoredIaosToken(): string {
  const token = localStorage.getItem(IAOS_CONNECTION_STORAGE.token);
  if (!token) throw new Error('未检测到 IAOS 调用凭据，请先完成联动检查');
  return token;
}

async function responseJson<T>(response: Response): Promise<T> {
  if (response.ok) return response.json() as Promise<T>;
  const status = response.status;
  let message = `IAOS 请求失败 (${status})`;
  let code = 'api_request_failed';
  let retryable = false;
  let runID = '';
  let runVersion = '';
  let requiredPermission = '';
  try {
    const body = await response.json() as {
      error?: string;
      message?: string;
      code?: string;
      retryable?: boolean;
      run_id?: string;
      run_version?: string;
      required_permission?: string;
    };
    code = body.code || code;
    message = body.message ?? body.error ?? message;
    retryable = body.retryable ?? retryable;
    runID = body.run_id ?? runID;
    runVersion = body.run_version ?? runVersion;
    requiredPermission = body.required_permission ?? requiredPermission;
  } catch { /* response is not JSON */ }
  throw new OrchestrationError({
    status,
    code,
    message: sanitizeSensitiveMessage(message),
    retryable,
    runId: runID,
    runVersion,
    requiredPermission,
  });
}

function sanitizeSensitiveMessage(value: string): string {
  return value.replace(/\b(?:Bearer|BEARER|bearer)\s+[A-Za-z0-9._~+/=:-]+/g, 'Bearer [REDACTED]');
}

export function defaultIntegrationConfiguration(): IntegrationConfiguration {
  const host = typeof window === 'undefined' ? '127.0.0.1' : window.location.hostname;
  return {
    tenantId: localStorage.getItem(IAOS_CONNECTION_STORAGE.tenantId) ?? 'tenant-hctm',
    packKey: 'hctm',
    scenarioKey: 'order-expedite-01',
    baseUrl: localStorage.getItem(IAOS_CONNECTION_STORAGE.baseUrl) ?? '',
    uiUrl: localStorage.getItem(IAOS_CONNECTION_STORAGE.uiUrl) ?? `http://${host}:3000`,
    orchestratorUrl: localStorage.getItem(IAOS_CONNECTION_STORAGE.orchestratorUrl) ?? defaultOrchestratorUrl(),
  };
}

export async function connectAndInspectIAOS(configuration: IntegrationConfiguration): Promise<IntegrationCheckResult> {
  const iaosBase = resolveIaosApiBaseUrl(configuration.baseUrl);
  const tokenResponse = await fetch(apiUrl(iaosBase, `/api/v1/dev/token?tenant_id=${encodeURIComponent(configuration.tenantId)}&roles=admin`));
  const { token } = await responseJson<{ token: string }>(tokenResponse);
  if (!token) throw new Error('IAOS 未返回可用的本地演示身份');

  const headers = {
    Accept: 'application/json',
    Authorization: `Bearer ${token}`,
    'X-Tenant-ID': configuration.tenantId,
  };
  const request = <T,>(path: string) => fetch(apiUrl(iaosBase, path), { headers }).then(responseJson<T>);
  const scenarioPath = `/api/v1/scenarios/${encodeURIComponent(configuration.packKey)}/${encodeURIComponent(configuration.scenarioKey)}`;

  const [profile, snapshot, eventPage, ...entityPages] = await Promise.all([
    request<{ tenant_name?: string; tenant_id?: string }>('/api/v1/profile'),
    request<IaosScenarioSnapshot>(`${scenarioPath}/snapshot`),
    request<{ items: unknown[] }>(`${scenarioPath}/events?after=0&limit=1`),
    ...(['sales_order', 'work_order', 'inventory', 'equipment'] as const).map((entity) =>
      request<{ total: number }>(`/api/v1/entities/${entity}/records?page=1&page_size=1`),
    ),
  ]);

  localStorage.setItem(IAOS_CONNECTION_STORAGE.token, token);
  localStorage.setItem(IAOS_CONNECTION_STORAGE.tenantId, configuration.tenantId);
  localStorage.setItem(IAOS_CONNECTION_STORAGE.baseUrl, iaosBase);
  localStorage.setItem(IAOS_CONNECTION_STORAGE.uiUrl, configuration.uiUrl.trim().replace(/\/$/, ''));
  localStorage.setItem(IAOS_CONNECTION_STORAGE.orchestratorUrl, resolveOrchestratorBaseUrl(configuration.orchestratorUrl));

  const entityKeys = ['sales_order', 'work_order', 'inventory', 'equipment'] as const;
  return {
    snapshot,
    tenantName: profile.tenant_name ?? profile.tenant_id ?? configuration.tenantId,
    counts: Object.fromEntries(entityKeys.map((key, index) => [key, entityPages[index].total])) as IntegrationCheckResult['counts'],
    eventChannelAvailable: Array.isArray(eventPage.items),
  };
}

export function getStoredActiveRunId(): string {
  return localStorage.getItem(IAOS_CONNECTION_STORAGE.activeRunId) ?? '';
}

export function setStoredActiveRunId(value: string): void {
  localStorage.setItem(IAOS_CONNECTION_STORAGE.activeRunId, value);
}

export function clearStoredActiveRunId(): void {
  localStorage.removeItem(IAOS_CONNECTION_STORAGE.activeRunId);
}

function asString(value: unknown): string | undefined {
  return typeof value === 'string' && value.length > 0 ? value : undefined;
}

function asNonNegativeInt(value: unknown, fallback: number): number {
  const typed = typeof value === 'number' ? value : Number(value);
  return Number.isSafeInteger(typed) && typed >= 0 ? typed : fallback;
}

export function toRunContext(run: OrchestrationRun): OrchestrationRunContext {
  return {
    runId: run.run_id,
    runVersion: run.run_version,
    planHash: run.plan_hash,
    status: run.status,
    currentAct: run.current_act,
    totalActs: run.total_acts,
    cursor: run.cursor,
    tenantId: run.tenant,
    packKey: run.pack_key,
    scenarioKey: run.scenario_key,
    target: run.target,
    updatedAt: new Date().toISOString(),
  };
}

export function getStoredRunContext(): OrchestrationRunContext | null {
  const raw = localStorage.getItem(IAOS_CONNECTION_STORAGE.activeRunContext);
  if (!raw) return null;
  try {
    const parsed = JSON.parse(raw) as unknown;
    if (!parsed || typeof parsed !== 'object') return null;

    const record = parsed as Record<string, unknown>;
    const runId = asString(record.runId);
    const runVersion = asString(record.runVersion);
    const planHash = asString(record.planHash);
    const status = asString(record.status);
    const tenantId = asString(record.tenantId);
    const packKey = asString(record.packKey);
    const scenarioKey = asString(record.scenarioKey);
    const target = asString(record.target);
    const updatedAt = asString(record.updatedAt);

    if (!runId || !runVersion || !planHash || !status || !tenantId || !packKey || !scenarioKey || !target || !updatedAt) {
      return null;
    }

    return {
      runId,
      runVersion,
      planHash,
      status,
      tenantId,
      packKey,
      scenarioKey,
      target,
      cursor: asNonNegativeInt(record.cursor, 0),
      currentAct: asNonNegativeInt(record.currentAct, 0),
      totalActs: asNonNegativeInt(record.totalActs, 0),
      updatedAt,
    };
  } catch {
    return null;
  }
}

export function setStoredRunContext(context: OrchestrationRunContext): void {
  localStorage.setItem(IAOS_CONNECTION_STORAGE.activeRunContext, JSON.stringify(context));
}

export function clearStoredRunContext(): void {
  localStorage.removeItem(IAOS_CONNECTION_STORAGE.activeRunContext);
}

export function iaosBusinessUrl(uiUrl: string, menuKey: string) {
  return `${uiUrl.trim().replace(/\/$/, '')}/#${menuKey}`;
}

function buildOrchestrationHeaders(token: string) {
  return {
    Accept: 'application/json',
    Authorization: `Bearer ${token}`,
    'Content-Type': 'application/json',
  };
}

function orchestrationRequest<T>(
  configuration: IntegrationConfiguration,
  path: string,
  method: string,
  body?: Record<string, unknown> | undefined,
  headers: Record<string, string> = {},
) {
  const orchestratorBase = resolveOrchestratorBaseUrl(configuration.orchestratorUrl);
  const token = getStoredIaosToken();
  return fetch(orchestrationApiUrl(orchestratorBase, path), {
    method,
    headers: { ...buildOrchestrationHeaders(token), ...headers },
    body: body ? JSON.stringify(body) : undefined,
  }).then((response) => responseJson<T>(response));
}

export async function createScenarioRun(
  configuration: IntegrationConfiguration,
  runTarget: string,
  options: {
    runID?: string;
    planHash?: string;
    actor?: string;
  } = {},
): Promise<OrchestrationRun> {
  const effectiveTarget = resolveIaosApiBaseUrl(runTarget);
  const response = await orchestrationRequest<OrchestrationRun>(
    configuration,
    '/runs',
    'POST',
    {
      target: effectiveTarget,
      pack_dir: `scenario-packs/${configuration.packKey}`,
      tenant: configuration.tenantId,
      actor: options.actor,
      run_id: options.runID,
      plan_hash: options.planHash,
      story_key: configuration.scenarioKey,
    },
  );

  return response;
}

export async function getRunStatus(configuration: IntegrationConfiguration, runId: string): Promise<OrchestrationRun> {
  return orchestrationRequest<OrchestrationRun>(configuration, `/runs/${encodeURIComponent(runId)}`, 'GET');
}

export async function getRunPlan(configuration: IntegrationConfiguration, storyKey: string): Promise<OrchestrationPlan> {
  return orchestrationRequest<OrchestrationPlan>(
    configuration,
    '/runs/plan',
    'POST',
    {
      story_key: storyKey,
      pack_dir: `scenario-packs/${configuration.packKey}`,
    },
  );
}

export async function executeRunAction(
  configuration: IntegrationConfiguration,
  runId: string,
  action: string,
  options: {
    planHash?: string;
    runVersion?: string;
    expectedCursor?: number;
    apply?: boolean;
    dryRun?: boolean;
    idempotencyKey?: string;
    confirmationToken?: string;
  } = {},
): Promise<OrchestrationActionResponse> {
  const { planHash, runVersion, expectedCursor, apply, dryRun, idempotencyKey, confirmationToken } = options;
  const headers: Record<string, string> = {};
  if (idempotencyKey) {
    headers['Idempotency-Key'] = idempotencyKey;
  }
  if (confirmationToken) {
    headers['X-Aese-Reset-Token'] = confirmationToken;
  }

  return orchestrationRequest<OrchestrationActionResponse>(configuration, `/runs/${encodeURIComponent(runId)}/${action}`, 'POST', {
    plan_hash: planHash,
    run_version: runVersion,
    expected_cursor: expectedCursor,
    apply,
    dry_run: dryRun,
    confirmation_token: confirmationToken,
  }, headers);
}
