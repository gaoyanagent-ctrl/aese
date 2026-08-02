import { describe, expect, it, vi } from "vitest";
import { confirmSiteInvestigationReport, generateFacilityProjectOptions, generatePlantRequirementOptions, loadPlantBuild } from "./plantBuild";

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

describe("M10 minimal player commands", () => {
  it("does not let the browser submit external investigation facts", async () => {
    localStorage.setItem("iaos_token", "token");
    localStorage.setItem("aese_iaos_tenant_id", "tenant-a");
    const fetchMock = vi.fn(async (_input: RequestInfo | URL, _init?: RequestInit) => ({ ok: true, json: async () => ({ status: "committed" }) }));
    vi.stubGlobal("fetch", fetchMock);
    await confirmSiteInvestigationReport("INC-1", "REQ-1", "INV-1");
    const body = JSON.parse(fetchMock.mock.calls[0][1]?.body as string);
    expect(body).toEqual({ schema_version: "1.0", case_code: "INC-1", requirement_id: "REQ-1", investigation_request_id: "INV-1", action: "accept_report" });
    expect(body).not.toHaveProperty("observation");
    vi.unstubAllGlobals();
  });

  it("asks the adviser for options using only the case identity", async () => {
    localStorage.setItem("iaos_token", "token");
    localStorage.setItem("aese_iaos_tenant_id", "tenant-a");
    const fetchMock = vi.fn(async (_input: RequestInfo | URL, _init?: RequestInit) => ({ ok: true, json: async () => ({ schema_version: "1.0", options: [] }) }));
    vi.stubGlobal("fetch", fetchMock);
    await generatePlantRequirementOptions("INC-1");
    expect(JSON.parse(fetchMock.mock.calls[0][1]?.body as string)).toEqual({ case_code: "INC-1" });
    vi.unstubAllGlobals();
  });

  it("asks the facility project agent using only the case identity", async () => {
    localStorage.setItem("iaos_token", "token");
    localStorage.setItem("aese_iaos_tenant_id", "tenant-a");
    const fetchMock = vi.fn(async (_input: RequestInfo | URL, _init?: RequestInit) => ({ ok: true, json: async () => ({ schema_version: "1.0", options: [] }) }));
    vi.stubGlobal("fetch", fetchMock);
    await generateFacilityProjectOptions("INC-1");
    expect(JSON.parse(fetchMock.mock.calls[0][1]?.body as string)).toEqual({ case_code: "INC-1" });
    vi.unstubAllGlobals();
  });
});
