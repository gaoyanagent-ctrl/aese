import { afterEach, describe, expect, it, vi } from "vitest";
import { loadIncorporation, resolveIaosLifecycleBase } from "./incorporation";

describe("incorporation API", () => {
  afterEach(() => {
    vi.unstubAllGlobals();
    localStorage.clear();
    window.location.hash = "";
  });

  it("normalizes actor knowledge", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue(
        new Response(JSON.stringify({ frames: [{ knowledge: null }] }), {
          status: 200,
        }),
      ),
    );
    const trace = await loadIncorporation();
    expect(trace.frames[0].knowledge).toEqual([]);
  });

  it("rejects a stale AESE-origin IAOS base and uses port 8082", async () => {
    localStorage.setItem("aese_iaos_base_url", window.location.origin);
    localStorage.setItem("aese_iaos_tenant_id", "tenant-hctm-genesis");
    localStorage.setItem("iaos_token", "test-token");
    window.location.hash =
      "#world-incorporation?tenant=tenant-hctm-genesis&case=INC-E2E-1";
    const fetcher = vi
      .fn()
      .mockResolvedValueOnce(
        new Response(JSON.stringify({ frames: [] }), { status: 200 }),
      )
      .mockResolvedValueOnce(
        new Response(
          JSON.stringify({
            case_code: "INC-E2E-1",
            state: { state: "enterprise_operational_ready" },
          }),
          { status: 200 },
        ),
      );
    vi.stubGlobal("fetch", fetcher);

    expect(resolveIaosLifecycleBase()).toBe(
      `http://${window.location.hostname}:8082`,
    );
    const trace = await loadIncorporation();
    expect(fetcher.mock.calls[1][0]).toBe(
      `http://${window.location.hostname}:8082/api/v1/incorporations/INC-E2E-1/trace`,
    );
    expect(trace.iaos_lifecycle?.case_code).toBe("INC-E2E-1");
  });

  it("accepts a lifecycle token handoff and removes it from the URL", async () => {
    localStorage.setItem("iaos_token", "stale-other-tenant-token");
    window.location.hash =
      "#world-incorporation?tenant=tenant-hctm-genesis&case=INC-E2E-1&auth_token=fresh-genesis-token";
    const fetcher = vi
      .fn()
      .mockResolvedValueOnce(
        new Response(JSON.stringify({ frames: [] }), { status: 200 }),
      )
      .mockResolvedValueOnce(
        new Response(JSON.stringify({ case_code: "INC-E2E-1" }), {
          status: 200,
        }),
      );
    vi.stubGlobal("fetch", fetcher);

    await loadIncorporation();

    expect(fetcher.mock.calls[1][1]?.headers).toMatchObject({
      Authorization: "Bearer fresh-genesis-token",
    });
    expect(localStorage.getItem("iaos_token")).toBe("fresh-genesis-token");
    expect(window.location.hash).not.toContain("auth_token");
  });
});
