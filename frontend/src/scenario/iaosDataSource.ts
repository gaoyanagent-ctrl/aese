import type { IaosScenarioEvent, IaosScenarioSnapshot } from './types';

export class IaosScenarioError extends Error {
  constructor(message: string, readonly status?: number) { super(message); }
}

interface IaosScenarioDataSourceOptions {
  baseUrl: string;
  token: string;
  tenantId: string;
  fetch?: typeof globalThis.fetch;
}

export class IaosScenarioDataSource {
  private readonly fetcher: typeof globalThis.fetch;
  private readonly baseUrl: string;
  constructor(private readonly options: IaosScenarioDataSourceOptions) {
    this.baseUrl = options.baseUrl.replace(/\/$/, '');
    this.fetcher = options.fetch ?? globalThis.fetch.bind(globalThis);
  }
  private headers(accept = 'application/json'): HeadersInit {
    return { Accept: accept, Authorization: `Bearer ${this.options.token}`, 'X-Tenant-ID': this.options.tenantId };
  }
  private path(pack: string, story: string, suffix: string) {
    return `${this.baseUrl}/api/v1/scenarios/${encodeURIComponent(pack)}/${encodeURIComponent(story)}/${suffix}`;
  }
  async snapshot(pack: string, story: string, signal?: AbortSignal): Promise<IaosScenarioSnapshot> {
    const response = await this.fetcher(this.path(pack, story, 'snapshot'), { headers: this.headers(), signal });
    if (!response.ok) throw await this.error(response);
    return response.json() as Promise<IaosScenarioSnapshot>;
  }
  async eventsAfter(pack: string, story: string, cursor: number, signal?: AbortSignal): Promise<IaosScenarioEvent[]> {
    const response = await this.fetcher(`${this.path(pack, story, 'events')}?after=${cursor}&limit=200`, { headers: this.headers(), signal });
    if (!response.ok) throw await this.error(response);
    return ((await response.json()) as { items: IaosScenarioEvent[] }).items;
  }
  async stream(pack: string, story: string, cursor: number, signal: AbortSignal, onEvent: (event: IaosScenarioEvent) => void): Promise<void> {
    const response = await this.fetcher(`${this.path(pack, story, 'events/stream')}?after=${cursor}`, { headers: this.headers('text/event-stream'), signal });
    if (!response.ok) throw await this.error(response);
    if (!response.body) throw new IaosScenarioError('IAOS 未返回事件流');
    const reader = response.body.getReader();
    const decoder = new TextDecoder();
    let buffer = '';
    while (true) {
      const { done, value } = await reader.read();
      if (done) return;
      buffer += decoder.decode(value, { stream: true }).replace(/\r\n/g, '\n');
      let boundary = buffer.indexOf('\n\n');
      while (boundary >= 0) {
        const frame = buffer.slice(0, boundary);
        buffer = buffer.slice(boundary + 2);
        const data = frame.split('\n').filter((line) => line.startsWith('data:')).map((line) => line.slice(5).trim()).join('\n');
        if (data) onEvent(JSON.parse(data) as IaosScenarioEvent);
        boundary = buffer.indexOf('\n\n');
      }
    }
  }
  private async error(response: Response) {
    let message = `IAOS 请求失败 (${response.status})`;
    try {
      const body = await response.json() as { error?: string; message?: string };
      message = body.message ?? body.error ?? message;
    } catch { /* non-JSON response */ }
    return new IaosScenarioError(message, response.status);
  }
}

export function liveDataSourceFromEnvironment() {
  return new IaosScenarioDataSource({
    baseUrl: import.meta.env.VITE_IAOS_BASE_URL ?? '',
    tenantId: import.meta.env.VITE_IAOS_TENANT_ID ?? 'tenant-hctm',
    token: localStorage.getItem('iaos_token') ?? import.meta.env.VITE_IAOS_TOKEN ?? 'dev-token',
  });
}
