import { describe, expect, it } from "vitest";
import { derivePlantGameStage, type PlantGameFacts } from "./plantGame";

const empty: PlantGameFacts = {
  hasRequirement: false,
  proposalCount: 0,
  adoptedReviewCount: 0,
  investigationCount: 0,
  observedCount: 0,
  hasRecommendation: false,
  hasDecision: false,
  hasSiteControlRequest: false,
  hasSiteControl: false,
  hasProjectPlan: false,
  projectApprovalStatus: "",
  hasActiveProject: false,
  hasContractRFQ: false,
  hasContractBids: false,
  hasContractRecommendation: false,
  contractApprovalStatus: "",
  hasAwardedContract: false,
};

describe("M10 game projection stage", () => {
  it.each([
    [{}, "requirement"],
    [{ hasRequirement: true }, "proposal"],
    [{ proposalCount: 3 }, "proposal"],
    [{ adoptedReviewCount: 1 }, "investigation"],
    [{ investigationCount: 1 }, "investigation"],
    [{ observedCount: 1 }, "comparison"],
    [{ hasRecommendation: true }, "governance"],
    [{ hasDecision: true }, "site_control"],
    [{ hasDecision: true, hasSiteControlRequest: true }, "site_control"],
    [{ hasDecision: true, hasSiteControlRequest: true, hasSiteControl: true }, "project_planning"],
    [{ hasProjectPlan: true }, "project_approval"],
    [{ hasProjectPlan: true, projectApprovalStatus: "rejected" }, "project_planning"],
    [{ hasActiveProject: true }, "contract_sourcing"],
    [{ hasContractRFQ: true }, "contract_world"],
    [{ hasContractBids: true }, "contract_comparison"],
    [{ hasContractRecommendation: true }, "contract_approval"],
    [{ hasAwardedContract: true }, "contract_awarded"],
  ])("derives committed facts %o as %s", (partial, expected) => {
    expect(derivePlantGameStage({ ...empty, ...partial }).key).toBe(expected);
  });

  it("lets a formal decision dominate all earlier facts", () => {
    expect(derivePlantGameStage({
      ...empty,
      proposalCount: 4,
      investigationCount: 2,
      observedCount: 2,
      hasRecommendation: true,
      hasDecision: true,
    }).key).toBe("site_control");
  });

  it("opens the site-control task without treating the decision as delivered control", () => {
    expect(derivePlantGameStage({ ...empty, hasDecision: true }).key).toBe("site_control");
    expect(derivePlantGameStage({ ...empty, hasDecision: true, hasSiteControlRequest: true }).key).toBe("site_control");
    expect(derivePlantGameStage({ ...empty, hasDecision: true, hasSiteControl: true }).key).toBe("project_planning");
  });
});
