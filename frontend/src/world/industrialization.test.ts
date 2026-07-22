import { describe, expect, it, vi } from "vitest";
import { loadIndustrialization } from "./industrialization";
describe("industrialization API", () => {
  it("normalizes release collections", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn(async () => ({
        ok: true,
        json: async () => ({
          frames: [
            { releases: null, trials: null, knowledge: null, apqp_gates: null },
          ],
        }),
      })),
    );
    const t = await loadIndustrialization();
    expect(t.frames[0].releases).toEqual([]);
    expect(t.frames[0].apqp_gates).toEqual({});
    vi.unstubAllGlobals();
  });
});
