import type { SandboxScenario } from "../scenario/types";

export function createPlaybackFixture(eventCount = 3): SandboxScenario {
  const metric = (value: number) => ({
    value,
    unit: "件",
    risk: "normal" as const,
    trend: "flat" as const,
  });
  const kpis = (value: number) => ({
    orderDemand: metric(value),
    availableFinishedGoods: metric(900),
    materialShortageRisk: metric(0),
    capacityRisk: metric(0),
    deliveryRisk: metric(0),
  });

  return {
    key: "fixture",
    name: "播放测试",
    version: "0.1.0",
    dataSource: "preview",
    timezone: "Asia/Shanghai",
    startsAt: "2026-07-19T08:00:00+08:00",
    endsAt: "2026-07-19T09:00:00+08:00",
    defaultSpeed: 1,
    acts: [{ number: 1, title: "测试幕", eventRange: [1, eventCount] }],
    layout: {
      width: 100,
      height: 100,
      nodes: [{
        id: "node-a",
        businessCode: "NODE-A",
        label: "节点 A",
        kind: "process",
        position: { x: 0, y: 0 },
        status: "normal",
      }],
      edges: [{ id: "edge-a", source: "node-a", target: "node-a", status: "normal" }],
    },
    initialState: {
      entities: [{
        id: "order-a",
        type: "sales_order",
        businessCode: "SO-A",
        name: "订单 A",
        status: "draft",
        risk: "normal",
        attributes: { quantity: 1, untouched: true },
      }],
      kpis: kpis(1),
    },
    timeline: Array.from({ length: eventCount }, (_, index) => ({
      sequence: index + 1,
      id: `event-${index + 1}`,
      timestamp: `2026-07-19T08:0${index}:00+08:00`,
      eventType: "test.event",
      title: `事件 ${index + 1}`,
      description: "测试事件",
      act: 1,
      domain: "order",
      severity: index === 1 ? "critical" : "normal",
      relatedEntityIds: ["order-a"],
      delta: {
        nodeStatuses: [{ id: "node-a", status: index === 1 ? "critical" : "normal" }],
        entityUpdates: [{
          id: "order-a",
          status: `step-${index + 1}`,
          attributes: { quantity: index + 2 },
        }],
      },
      kpis: kpis(index + 2),
    })),
    agentOutputs: [],
  };
}
