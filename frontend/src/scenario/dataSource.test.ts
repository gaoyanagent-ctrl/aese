import { describe, expect, it } from "vitest";
import events from "../../../scenario-packs/hctm/stories/order-expedite-01/events.json";
import preview from "../../../scenario-packs/hctm/stories/order-expedite-01/preview.json";
import { StaticScenarioDataSource } from "./dataSource";
import { ScenarioValidationError, parseSandboxScenario } from "./validation";

describe("StaticScenarioDataSource", () => {
  it("loads the deterministic seven-act, 22-event preview", async () => {
    const source = new StaticScenarioDataSource({ "order-expedite-01": preview });
    const scenario = await source.loadScenario("order-expedite-01");

    expect(scenario.acts).toHaveLength(7);
    expect(scenario.timeline).toHaveLength(22);
    expect(scenario.timeline.map((event) => event.sequence)).toEqual(
      Array.from({ length: 22 }, (_, index) => index + 1),
    );
    expect(new Set(scenario.agentOutputs.map((output) => output.agent))).toEqual(
      new Set(["planning", "quality", "business_analysis"]),
    );
  });

  it("preserves the final shipment contract", async () => {
    const source = new StaticScenarioDataSource({ "order-expedite-01": preview });
    const scenario = await source.loadScenario("order-expedite-01");
    const finalEvent = scenario.timeline.at(-1)!;
    const finalOrder = finalEvent.delta.entityUpdates?.find(({ id }) => id === "order-main");
    const secondShipment = finalEvent.delta.entityUpdates?.find(({ id }) => id === "shipment-2");

    expect(finalOrder?.attributes).toMatchObject({ shippedQty: 11_700, shortageQty: 300 });
    expect(secondShipment?.attributes).toMatchObject({ shippedQty: 2_700, shortageQty: 300 });
    expect(finalEvent.description).toContain("累计实发 11,700 件");
  });

  it("adapts the canonical event identity, order and timestamps without drift", async () => {
    const source = new StaticScenarioDataSource({ "order-expedite-01": preview });
    const scenario = await source.loadScenario("order-expedite-01");

    expect(
      scenario.timeline.map(({ sequence, id, timestamp, eventType }) => ({ sequence, id, timestamp, eventType })),
    ).toEqual(
      events.events.map(({ sequence, event_id: id, timestamp, event_type: eventType }) => ({ sequence, id, timestamp, eventType })),
    );
    expect(scenario.timeline.every(({ delta, kpis }) => delta !== undefined && kpis !== undefined)).toBe(true);
  });

  it("fails clearly for a missing scenario key", async () => {
    const source = new StaticScenarioDataSource({});
    await expect(source.loadScenario("missing")).rejects.toThrow("Scenario not found: missing");
  });
});

describe("parseSandboxScenario", () => {
  it("rejects an empty timeline and incomplete act mapping", () => {
    const broken = structuredClone(preview) as Record<string, unknown>;
    broken.timeline = [];
    broken.acts = [];

    expect(() => parseSandboxScenario(broken)).toThrow(ScenarioValidationError);
    expect(() => parseSandboxScenario(broken)).toThrow("exactly 22 events");
  });

  it("rejects deltas that point to an unknown canvas node", () => {
    const broken = structuredClone(preview);
    broken.timeline[0].delta.nodeStatuses = [{ id: "unknown-node", status: "critical" }];

    expect(() => parseSandboxScenario(broken)).toThrow("references unknown node unknown-node");
  });
});
