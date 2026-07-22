import { describe, expect, it, vi } from "vitest";
import { loadCapabilityBuild } from "./capabilityBuild";
describe("capability build API", () => {
  it("normalizes private world collections", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn(async () => ({
        ok: true,
        json: async () => ({
          frames: [
            { equipment: null, workers: null, knowledge: null, gate: null },
          ],
        }),
      })),
    );
    const t = await loadCapabilityBuild();
    expect(t.frames[0].workers).toEqual([]);
    expect(t.frames[0].gate).toEqual({});
    vi.unstubAllGlobals();
  });
});
