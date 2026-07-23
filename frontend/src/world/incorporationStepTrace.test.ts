import { describe, expect, it } from "vitest";
import {
  assertIncorporationStepCoverage,
  buildIncorporationStepTrace,
  incorporationStepDefinitions,
} from "./incorporationStepTrace";

describe("incorporation step trace", () => {
  it("maps the eight frames to every committed transition exactly once", () => {
    expect(incorporationStepDefinitions).toHaveLength(8);
    expect(
      incorporationStepDefinitions.flatMap((x) => x.capabilities),
    ).toHaveLength(15);
    expect(assertIncorporationStepCoverage()).toEqual([]);
  });

  it("filters evidence to the selected frame and reports missing transitions", () => {
    const trace = buildIncorporationStepTrace(
      { step: 2 } as never,
      {
        process_runs: [
          {
            process_key: "enterprise.incorporation.lifecycle.v1",
            trace: [
              { capability: "registration.observation.commit" },
              { capability: "bank.account.opening.submit" },
            ],
          },
        ],
        journal: [
          { capability_code: "registration.observation.commit" },
          { capability_code: "bank.account.opening.submit" },
        ],
        approvals: [],
        decisions: [
          { operation: "registration.observation.commit" },
          { operation: "bank.account.opening.submit" },
        ],
        world_exchanges: [
          { payload_type: "registration.approved.v1" },
          { payload_type: "bank.account.opened.v1" },
        ],
      } as never,
    );
    expect(trace?.transitions).toHaveLength(1);
    expect(trace?.journal).toHaveLength(1);
    expect(trace?.decisions).toHaveLength(1);
    expect(trace?.worldExchanges).toHaveLength(1);
    expect(trace?.unmatchedCapabilities).toEqual([]);
  });
});
