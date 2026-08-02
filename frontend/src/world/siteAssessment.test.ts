import { describe, expect, it } from "vitest";
import { assessObservedSites } from "./siteAssessment";
import type { FacilityRequirement, ProposalSet, SiteInvestigationItem } from "./plantBuild";

const requirement = { investment_request: { value: "10000000.00", currency: "CNY", scale: 2 }, minimum_area_m2: 8000, minimum_electricity_kva: 2000, target_available_at: "2026-12-01T00:00:00Z" } as FacilityRequirement;
const proposalSet = { proposal_set_id: "SET-1", revision: 1, proposals: [{ proposal_id: "SITE-1", display_name: "Agent 候选 A", estimated_amount: { likely: { value: "7500000.00" } }, estimated_schedule: { likely: "2026-11-01T00:00:00Z" } }] } as ProposalSet;
const item = (overrides: Partial<NonNullable<SiteInvestigationItem["observation"]>> = {}): SiteInvestigationItem => ({
  request: { proposal_set_id: "SET-1", expected_revision: 1, proposal_id: "SITE-1" } as SiteInvestigationItem["request"], status: "observed", work_item_status: "completed",
  observation: { schema_version: "1.0", observation_id: "OBS-1", investigation_request_id: "INV-1", proposal_id: "SITE-1", result: "completed", ownership_status: "verified", available_area_m2: 9000, electricity_kva: 2500, quoted_amount: { value: "8000000.00", currency: "CNY", scale: 2 }, available_at: "2026-10-01T00:00:00Z", permit_status: "eligible", evidence_refs: ["world-document:1"], notes: "", external_actor_id: "park", observed_at: "2026-08-01T00:00:00Z", ...overrides },
});

describe("observed site assessment", () => {
  it("scores only delivered observations with explainable components", () => {
    const result = assessObservedSites(requirement, proposalSet, [item()], { cost: 35, schedule: 25, capacity: 20, control: 20 });
    expect(result[0].eligible).toBe(true);
    expect(result[0].total_score).toBeGreaterThan(0);
    expect(result[0].estimated).toEqual({ amount: "7500000.00", available_at: "2026-11-01T00:00:00Z" });
    expect(result[0].evidence_refs).toEqual(["world-document:1"]);
  });
  it("hard-fails capacity before weighted score", () => {
    const result = assessObservedSites(requirement, proposalSet, [item({ electricity_kva: 1200 })], { cost: 100, schedule: 0, capacity: 0, control: 0 });
    expect(result[0].eligible).toBe(false);
    expect(result[0].total_score).toBeNull();
    expect(result[0].hard_failures).toContain("可用电力低于最低要求");
    expect(result[0].criteria.find((criterion) => criterion.key === "electricity")).toMatchObject({
      required: "最低要求 2,000 kVA",
      observed: "实测 1,200 kVA",
      difference: "短缺 800 kVA",
      passed: false,
    });
  });
  it("excludes trusted observations from an older proposal-set revision", () => {
    const stale = item();
    stale.request.proposal_set_id = "SET-OLD";
    stale.request.expected_revision = 1;
    const result = assessObservedSites(requirement, proposalSet, [stale], { cost: 35, schedule: 25, capacity: 20, control: 20 });
    expect(result).toEqual([]);
  });
});
