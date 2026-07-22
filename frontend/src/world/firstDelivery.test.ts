import { describe, expect, it, vi } from "vitest";
import { loadFirstDelivery } from "./firstDelivery";
describe("first delivery", () => {
  it("normalizes shipments", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn(async () => ({
        ok: true,
        json: async () => ({ frames: [{ shipments: null }] }),
      })),
    );
    const t = await loadFirstDelivery();
    expect(t.frames[0].shipments).toEqual([]);
    vi.unstubAllGlobals();
  });
});
