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
    [{ hasDecision: true }, "selected"],
    [{ hasDecision: true, hasSiteControlRequest: true }, "site_control"],
    [{ hasDecision: true, hasSiteControlRequest: true, hasSiteControl: true }, "project_ready"],
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
    }).key).toBe("selected");
  });

  it("does not treat a formal site decision as physical site control", () => {
    expect(derivePlantGameStage({ ...empty, hasDecision: true }).key).toBe("selected");
    expect(derivePlantGameStage({ ...empty, hasDecision: true, hasSiteControlRequest: true }).key).toBe("site_control");
  });
});
