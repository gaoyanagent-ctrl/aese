import { afterEach, describe, expect, it, vi } from 'vitest';
import { loadGenesisTrace } from './api';

describe('loadGenesisTrace', () => {
  afterEach(() => vi.unstubAllGlobals());

  it('normalizes a nullable knowledge collection at the API boundary', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(new Response(JSON.stringify({
      schema_version: '1.0', world_run_id: 'r1', timezone: 'Asia/Shanghai', actor_ref: {},
      frames: [{ knowledge: null }],
    }), { status: 200 })));

    const trace = await loadGenesisTrace();
    expect(trace.frames[0].knowledge).toEqual([]);
  });
});
