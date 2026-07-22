import { beforeEach, describe, expect, it, vi } from 'vitest';
import {
  IAOS_CONNECTION_STORAGE,
  clearStoredRunContext,
  connectAndInspectIAOS,
  executeRunAction,
  getStoredRunContext,
  iaosBusinessUrl,
  setStoredRunContext,
} from './iaosIntegration';

describe('IAOS visual integration', () => {
  beforeEach(() => localStorage.clear());

  it('acquires the HCTM identity and verifies all linked business sources', async () => {
    const fetcher = vi.fn()
      .mockResolvedValueOnce(new Response(JSON.stringify({ token: 'hctm-token' }), { status: 200 }))
      .mockResolvedValueOnce(new Response(JSON.stringify({ tenant_name: '华辰热管理系统集团' }), { status: 200 }))
      .mockResolvedValueOnce(new Response(JSON.stringify({ cursor: 6, entities: [], events: [], recommendations: [], gaps: [], kpis: {} }), { status: 200 }))
      .mockResolvedValueOnce(new Response(JSON.stringify({ items: [] }), { status: 200 }))
      .mockResolvedValueOnce(new Response(JSON.stringify({ total: 1 }), { status: 200 }))
      .mockResolvedValueOnce(new Response(JSON.stringify({ total: 4 }), { status: 200 }))
      .mockResolvedValueOnce(new Response(JSON.stringify({ total: 6 }), { status: 200 }))
      .mockResolvedValueOnce(new Response(JSON.stringify({ total: 1 }), { status: 200 }));
    vi.stubGlobal('fetch', fetcher);

    const result = await connectAndInspectIAOS({ tenantId: 'tenant-hctm', packKey: 'hctm', scenarioKey: 'order-expedite-01', baseUrl: '', uiUrl: 'http://localhost:3000' });

    expect(result.tenantName).toBe('华辰热管理系统集团');
    expect(result.counts).toEqual({ sales_order: 1, work_order: 4, inventory: 6, equipment: 1 });
    expect(fetcher.mock.calls[0][0]).toContain('/api/v1/dev/token');
    expect(localStorage.getItem('iaos_token')).toBe('hctm-token');
    vi.unstubAllGlobals();
  });

  it('builds a direct IAOS business menu link', () => {
    expect(iaosBusinessUrl('http://localhost:3000/', 'sales_order')).toBe('http://localhost:3000/#sales_order');
  });

  it('redacts token in IAOS error responses', async () => {
    const fetcher = vi.fn(async () => new Response(JSON.stringify({ error: 'upstream token validation failed: Bearer secret.token.value' }), { status: 401 })) as unknown as typeof fetch;
    vi.stubGlobal('fetch', fetcher);

    await expect(connectAndInspectIAOS({ tenantId: 'tenant-hctm', packKey: 'hctm', scenarioKey: 'order-expedite-01', baseUrl: '', uiUrl: 'http://localhost:3000' })).rejects.toThrowError('Bearer [REDACTED]');
    vi.unstubAllGlobals();
  });

  it('persists and clears run context round-trip', () => {
    const context = {
      runId: 'run-hctm-20260721',
      runVersion: 'v1',
      planHash: 'plan-hash-2026-07-21',
      status: 'ready',
      currentAct: 2,
      totalActs: 7,
      cursor: 12,
      tenantId: 'tenant-hctm',
      packKey: 'hctm',
      scenarioKey: 'order-expedite-01',
      target: 'http://tenant-hctm-ui',
      updatedAt: '2026-07-21T00:00:00Z',
    };

    setStoredRunContext(context);
    expect(getStoredRunContext()).toEqual(context);

    clearStoredRunContext();
    expect(getStoredRunContext()).toBeNull();
  });

  it('adds orchestration auth/idempotency/reset headers and action payload', async () => {
    const calls: string[] = [];
    const fetcher = vi.fn(async (input: string, init: RequestInit | undefined) => {
      calls.push(`${input}|${JSON.stringify(init?.headers ?? {})}|${init?.body ?? ''}`);
      return new Response(
        JSON.stringify({
          run: {
            run_id: 'run-hctm-20260721',
            run_version: 'v2',
            pack_key: 'hctm',
            pack_version: '0.1.0',
            scenario_key: 'order-expedite-01',
            plan_hash: 'plan-hash-2026-07-21',
            status: 'ready',
            current_act: 1,
            total_acts: 7,
            cursor: 1,
            tenant: 'tenant-hctm',
            target: 'http://127.0.0.1:3000',
            created_at: '2026-07-21T00:00:00Z',
            updated_at: '2026-07-21T00:00:00Z',
            allowed_actions: ['reset'],
          },
        }),
        { status: 200 },
      );
    }) as unknown as typeof fetch;
    vi.stubGlobal('fetch', fetcher);
    localStorage.setItem(IAOS_CONNECTION_STORAGE.token, 'demo-token');

    await executeRunAction(
      {
        tenantId: 'tenant-hctm',
        packKey: 'hctm',
        scenarioKey: 'order-expedite-01',
        baseUrl: '',
        uiUrl: 'http://localhost:3000',
        orchestratorUrl: 'http://localhost:8090',
      },
      'run-hctm-20260721',
      'reset',
      {
        planHash: 'plan-hash-2026-07-21',
        runVersion: 'v1',
        expectedCursor: 9,
        apply: true,
        dryRun: false,
        idempotencyKey: 'idem-reset-001',
        confirmationToken: 'confirm-reset-token',
      },
    );

    const [url, rawHeaders, rawBody] = calls[0].split('|');
    expect(url).toContain('http://localhost:8090/api/aese/v1/runs/run-hctm-20260721/reset');

    const headers = new Headers(JSON.parse(rawHeaders));
    expect(headers.get('Authorization')).toBe('Bearer demo-token');
    expect(headers.get('Idempotency-Key')).toBe('idem-reset-001');
    expect(headers.get('X-Aese-Reset-Token')).toBe('confirm-reset-token');

    const body = JSON.parse(rawBody);
    expect(body).toEqual({
      plan_hash: 'plan-hash-2026-07-21',
      run_version: 'v1',
      expected_cursor: 9,
      apply: true,
      dry_run: false,
      confirmation_token: 'confirm-reset-token',
    });

    vi.unstubAllGlobals();
  });
});
