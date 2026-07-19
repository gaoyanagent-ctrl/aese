import { describe, expect, it } from "vitest";
import preview from "../../../scenario-packs/hctm/stories/order-expedite-01/preview.json";

import { parseSandboxScenario } from "../scenario/validation";
import {
  createInitialPlaybackState,
  playbackReducer,
  replayToStep,
} from "./reducer";
import { createPlaybackFixture } from "./testFixture";

describe("playback reducer", () => {
  it("replays deltas deterministically without mutating the scenario", () => {
    const scenario = createPlaybackFixture();
    const original = structuredClone(scenario);

    const atStepTwo = replayToStep(scenario, 2);

    expect(atStepTwo.entities[0]).toMatchObject({
      status: "step-2",
      attributes: { quantity: 3, untouched: true },
    });
    expect(atStepTwo.nodeStatuses["node-a"]).toBe("critical");
    expect(atStepTwo.kpis.orderDemand.value).toBe(3);
    expect(scenario).toEqual(original);
    expect(replayToStep(scenario, 2)).toEqual(atStepTwo);
  });

  it("supports bounds, previous, seek and a repeatable reset", () => {
    const scenario = createPlaybackFixture();
    let state = createInitialPlaybackState(scenario);

    state = playbackReducer(state, { type: "previous" });
    expect(state.currentStep).toBe(0);

    state = playbackReducer(state, { type: "seek", step: 99 });
    expect(state.currentStep).toBe(3);
    expect(state.status).toBe("paused");

    state = playbackReducer(state, { type: "previous" });
    expect(state.currentStep).toBe(2);
    expect(state.viewState).toEqual(replayToStep(scenario, 2));

    state = playbackReducer(state, { type: "reset" });
    expect(state).toEqual(createInitialPlaybackState(scenario));
  });

  it("automatically pauses at the final event and cannot play past it", () => {
    const scenario = createPlaybackFixture(22);
    let state = playbackReducer(createInitialPlaybackState(scenario), { type: "play" });
    for (let index = 0; index < 22; index += 1) {
      state = playbackReducer(state, { type: "tick" });
    }

    expect(state.currentStep).toBe(22);
    expect(state.status).toBe("paused");
    expect(playbackReducer(state, { type: "play" }).status).toBe("paused");
    expect(playbackReducer(state, { type: "tick" }).currentStep).toBe(22);
  });

  it("accepts only the supported playback speeds", () => {
    const initial = createInitialPlaybackState(createPlaybackFixture());
    expect(playbackReducer(initial, { type: "set-speed", speed: 4 }).speed).toBe(4);
  });

  it("replays the canonical 22-event story to its shipment outcome", () => {
    const scenario = parseSandboxScenario(preview);
    const finalState = replayToStep(scenario, scenario.timeline.length);
    const order = finalState.entities.find(({ id }) => id === "order-main");
    const secondShipment = finalState.entities.find(({ id }) => id === "shipment-2");

    expect(scenario.timeline).toHaveLength(22);
    expect(order?.attributes).toMatchObject({ shippedQty: 11_700, shortageQty: 300 });
    expect(secondShipment?.attributes).toMatchObject({ shippedQty: 2_700, shortageQty: 300 });
    expect(finalState.kpis.deliveryRisk.risk).toBe("critical");
  });
});
