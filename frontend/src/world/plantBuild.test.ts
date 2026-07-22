import { describe, expect, it, vi } from "vitest";
import { loadPlantBuild } from "./plantBuild";

describe("loadPlantBuild", () => {
  it("normalizes nullable campaign collections", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn(async () => ({
        ok: true,
        json: async () => ({
          frames: [
            {
              assessments: null,
              zones: null,
              work_packages: null,
              knowledge: null,
              utilities: null,
            },
          ],
        }),
      })),
    );
    const trace = await loadPlantBuild();
    expect(trace.frames[0].assessments).toEqual([]);
    expect(trace.frames[0].zones).toEqual([]);
    expect(trace.frames[0].utilities).toEqual({});
    vi.unstubAllGlobals();
  });
});
