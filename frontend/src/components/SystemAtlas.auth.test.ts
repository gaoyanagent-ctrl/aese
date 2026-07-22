import { beforeEach, describe, expect, it, vi } from 'vitest';
import { atlasFetch } from './SystemAtlas';

describe('System Atlas authentication recovery', () => {
  beforeEach(() => {
    localStorage.clear();
    vi.restoreAllMocks();
  });

  it('refreshes a stale stored token once after a 401 and retries the Atlas request', async () => {
    localStorage.setItem('iaos_token', 'stale-token');
    const fetcher = vi.spyOn(globalThis, 'fetch')
      .mockResolvedValueOnce(new Response(JSON.stringify({ error: 'expired' }), { status: 401 }))
      .mockResolvedValueOnce(new Response(JSON.stringify({ token: 'fresh-token' }), { status: 200 }))
      .mockResolvedValueOnce(new Response(JSON.stringify({ project: 'aese' }), { status: 200 }));

    const response = await atlasFetch('/api/v1/system-atlas?view=aese');

    expect(response.status).toBe(200);
    expect(localStorage.getItem('iaos_token')).toBe('fresh-token');
    expect(fetcher).toHaveBeenCalledTimes(3);
    expect(new Headers(fetcher.mock.calls[0][1]?.headers).get('Authorization')).toBe('Bearer stale-token');
    expect(fetcher.mock.calls[1][0]).toContain('/api/v1/dev/token?');
    expect(new Headers(fetcher.mock.calls[2][1]?.headers).get('Authorization')).toBe('Bearer fresh-token');
  });
});
