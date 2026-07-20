import type { IaosScenarioSnapshot } from '../scenario/types';

export const IAOS_CONNECTION_STORAGE = {
  token: 'iaos_token',
  tenantId: 'aese_iaos_tenant_id',
  baseUrl: 'aese_iaos_base_url',
  uiUrl: 'aese_iaos_ui_url',
} as const;

export interface IntegrationConfiguration {
  tenantId: string;
  packKey: string;
  scenarioKey: string;
  baseUrl: string;
  uiUrl: string;
}

export interface IntegrationCheckResult {
  snapshot: IaosScenarioSnapshot;
  tenantName: string;
  counts: Record<'sales_order' | 'work_order' | 'inventory' | 'equipment', number>;
  eventChannelAvailable: boolean;
}

function normalizedBaseUrl(value: string) {
  return value.trim().replace(/\/$/, '');
}

function apiUrl(baseUrl: string, path: string) {
  return `${normalizedBaseUrl(baseUrl)}${path}`;
}

async function responseJson<T>(response: Response): Promise<T> {
  if (response.ok) return response.json() as Promise<T>;
  let message = `IAOS 请求失败 (${response.status})`;
  try {
    const body = await response.json() as { error?: string; message?: string };
    message = body.message ?? body.error ?? message;
  } catch { /* response is not JSON */ }
  throw new Error(message);
}

export function defaultIntegrationConfiguration(): IntegrationConfiguration {
  const host = typeof window === 'undefined' ? '127.0.0.1' : window.location.hostname;
  return {
    tenantId: localStorage.getItem(IAOS_CONNECTION_STORAGE.tenantId) ?? 'tenant-hctm',
    packKey: 'hctm',
    scenarioKey: 'order-expedite-01',
    baseUrl: localStorage.getItem(IAOS_CONNECTION_STORAGE.baseUrl) ?? '',
    uiUrl: localStorage.getItem(IAOS_CONNECTION_STORAGE.uiUrl) ?? `http://${host}:3000`,
  };
}

export async function connectAndInspectIAOS(configuration: IntegrationConfiguration): Promise<IntegrationCheckResult> {
  const tokenResponse = await fetch(apiUrl(configuration.baseUrl, `/api/v1/dev/token?tenant_id=${encodeURIComponent(configuration.tenantId)}&roles=admin`));
  const { token } = await responseJson<{ token: string }>(tokenResponse);
  if (!token) throw new Error('IAOS 未返回可用的本地演示身份');

  const headers = {
    Accept: 'application/json',
    Authorization: `Bearer ${token}`,
    'X-Tenant-ID': configuration.tenantId,
  };
  const request = <T,>(path: string) => fetch(apiUrl(configuration.baseUrl, path), { headers }).then(responseJson<T>);
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
  localStorage.setItem(IAOS_CONNECTION_STORAGE.baseUrl, normalizedBaseUrl(configuration.baseUrl));
  localStorage.setItem(IAOS_CONNECTION_STORAGE.uiUrl, configuration.uiUrl.trim().replace(/\/$/, ''));

  const entityKeys = ['sales_order', 'work_order', 'inventory', 'equipment'] as const;
  return {
    snapshot,
    tenantName: profile.tenant_name ?? profile.tenant_id ?? configuration.tenantId,
    counts: Object.fromEntries(entityKeys.map((key, index) => [key, entityPages[index].total])) as IntegrationCheckResult['counts'],
    eventChannelAvailable: Array.isArray(eventPage.items),
  };
}

export function iaosBusinessUrl(uiUrl: string, menuKey: string) {
  return `${uiUrl.trim().replace(/\/$/, '')}/#${menuKey}`;
}
