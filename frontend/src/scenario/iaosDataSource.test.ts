import { describe, expect, it, vi } from 'vitest';
import { IaosScenarioDataSource } from './iaosDataSource';

function source(fetcher: typeof fetch) {
  return new IaosScenarioDataSource({ baseUrl: 'http://iaos.test/', token: 'secret', tenantId: 'tenant-hctm', fetch: fetcher });
}

describe('IaosScenarioDataSource', () => {
  it('loads an authenticated tenant snapshot', async () => {
    const fetcher = vi.fn(async (_input: RequestInfo | URL, init?: RequestInit) => {
      expect(init?.headers).toMatchObject({ Authorization: 'Bearer secret', 'X-Tenant-ID': 'tenant-hctm' });
      return new Response(JSON.stringify({ cursor: 6, events: [] }), { status: 200 });
    }) as unknown as typeof fetch;
    await expect(source(fetcher).snapshot('hctm', 'order-expedite-01')).resolves.toMatchObject({ cursor: 6 });
    expect(fetcher).toHaveBeenCalledWith('http://iaos.test/api/v1/scenarios/hctm/order-expedite-01/snapshot', expect.anything());
  });

  it('maps authorization failures without falling back to Preview', async () => {
    const fetcher = vi.fn(async () => new Response(JSON.stringify({ error: 'scenario.read permission required' }), { status: 403 })) as unknown as typeof fetch;
    await expect(source(fetcher).snapshot('hctm', 'order-expedite-01')).rejects.toMatchObject({ status: 403, message: 'scenario.read permission required' });
  });

  it('parses SSE frames and preserves the persistent cursor', async () => {
    const encoder = new TextEncoder();
    const body = new ReadableStream({ start(controller) { controller.enqueue(encoder.encode(': heartbeat\n\nid: 7\nevent: scenario\ndata: {"cursor":7,"event_id":"evt-7"}\n\n')); controller.close(); } });
    const fetchMock = vi.fn(async (input: RequestInfo | URL) => {
      expect(String(input)).toContain('events/stream');
      return new Response(body, { status: 200, headers: { 'Content-Type': 'text/event-stream' } });
    });
    const fetcher = fetchMock as unknown as typeof fetch;
    const events: Array<{ cursor: number }> = [];
    await source(fetcher).stream('hctm', 'order-expedite-01', 6, new AbortController().signal, (event) => events.push(event));
    expect(events).toEqual([{ cursor: 7, event_id: 'evt-7' }]);
    expect(String(fetchMock.mock.calls[0][0])).toContain('events/stream?after=6');
  });
});
