import { beforeEach, describe, expect, it, vi } from 'vitest';
import { connectAndInspectIAOS, iaosBusinessUrl } from './iaosIntegration';

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
    expect(localStorage.getItem('iaos_token')).toBe('hctm-token');
    vi.unstubAllGlobals();
  });

  it('builds a direct IAOS business menu link', () => {
    expect(iaosBusinessUrl('http://localhost:3000/', 'sales_order')).toBe('http://localhost:3000/#sales_order');
  });
});
