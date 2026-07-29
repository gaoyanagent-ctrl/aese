import { describe, expect, it } from "vitest";
import { missionPlace } from "./LocationScene";

describe("Enterprise Genesis finance mission routing", () => {
  it.each([
    "finance.organization.configure",
    "accounting.book.activate",
    "chart.of.accounts.activate",
    "finance.opening.readiness.evaluate",
  ])("routes %s to enterprise headquarters", capability => {
    expect(missionPlace(capability)).toBe("headquarters");
  });

  it("routes capital posting to the finance center after bank evidence is committed", () => {
    expect(missionPlace("capital.contribution.post")).toBe("headquarters");
  });
});
