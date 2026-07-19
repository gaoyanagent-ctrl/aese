import type {
  AgentOutput,
  KpiMetric,
  KpiSnapshot,
  RiskLevel,
  SandboxEdge,
  SandboxEntity,
  SandboxNode,
  SandboxScenario,
  ScenarioAct,
  ScenarioEvent,
  VisualDelta,
} from "./types";

export class ScenarioValidationError extends Error {
  constructor(readonly issues: string[]) {
    super(`Invalid sandbox scenario:\n- ${issues.join("\n- ")}`);
    this.name = "ScenarioValidationError";
  }
}

type JsonObject = Record<string, unknown>;
const risks: RiskLevel[] = ["normal", "watch", "critical"];
const domains = ["order", "planning", "supply", "equipment", "quality", "production", "logistics"] as const;
const agents = ["planning", "quality", "business_analysis"] as const;

function object(value: unknown, path: string, issues: string[]): JsonObject {
  if (typeof value !== "object" || value === null || Array.isArray(value)) {
    issues.push(`${path} must be an object`);
    return {};
  }
  return value as JsonObject;
}

function array(value: unknown, path: string, issues: string[]): unknown[] {
  if (!Array.isArray(value)) {
    issues.push(`${path} must be an array`);
    return [];
  }
  return value;
}

function string(value: unknown, path: string, issues: string[]): string {
  if (typeof value !== "string" || value.length === 0) {
    issues.push(`${path} must be a non-empty string`);
    return "";
  }
  return value;
}

function number(value: unknown, path: string, issues: string[]): number {
  if (typeof value !== "number" || !Number.isFinite(value)) {
    issues.push(`${path} must be a finite number`);
    return 0;
  }
  return value;
}

function oneOf<T extends string | number>(value: unknown, allowed: readonly T[], path: string, issues: string[]): T {
  if (!allowed.includes(value as T)) {
    issues.push(`${path} must be one of ${allowed.join(", ")}`);
    return allowed[0];
  }
  return value as T;
}

function rfc3339(value: unknown, path: string, issues: string[]): string {
  const result = string(value, path, issues);
  if (result && !/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?(?:Z|[+-]\d{2}:\d{2})$/.test(result)) {
    issues.push(`${path} must be RFC 3339 with an explicit offset`);
  }
  return result;
}

function metric(value: unknown, path: string, issues: string[]): KpiMetric {
  const item = object(value, path, issues);
  return {
    value: number(item.value, `${path}.value`, issues),
    unit: string(item.unit, `${path}.unit`, issues),
    risk: oneOf(item.risk, risks, `${path}.risk`, issues),
    trend: oneOf(item.trend, ["up", "down", "flat"] as const, `${path}.trend`, issues),
  };
}

function kpis(value: unknown, path: string, issues: string[]): KpiSnapshot {
  const item = object(value, path, issues);
  return {
    orderDemand: metric(item.orderDemand, `${path}.orderDemand`, issues),
    availableFinishedGoods: metric(item.availableFinishedGoods, `${path}.availableFinishedGoods`, issues),
    materialShortageRisk: metric(item.materialShortageRisk, `${path}.materialShortageRisk`, issues),
    capacityRisk: metric(item.capacityRisk, `${path}.capacityRisk`, issues),
    deliveryRisk: metric(item.deliveryRisk, `${path}.deliveryRisk`, issues),
  };
}

function node(value: unknown, path: string, issues: string[]): SandboxNode {
  const item = object(value, path, issues);
  const position = object(item.position, `${path}.position`, issues);
  return {
    id: string(item.id, `${path}.id`, issues),
    businessCode: string(item.businessCode, `${path}.businessCode`, issues),
    label: string(item.label, `${path}.label`, issues),
    kind: oneOf(item.kind, ["supplier", "warehouse", "quality", "process", "shipping", "customer"] as const, `${path}.kind`, issues),
    position: { x: number(position.x, `${path}.position.x`, issues), y: number(position.y, `${path}.position.y`, issues) },
    status: oneOf(item.status, risks, `${path}.status`, issues),
    ...(item.entityId === undefined ? {} : { entityId: string(item.entityId, `${path}.entityId`, issues) }),
  };
}

function edge(value: unknown, path: string, issues: string[]): SandboxEdge {
  const item = object(value, path, issues);
  return {
    id: string(item.id, `${path}.id`, issues),
    source: string(item.source, `${path}.source`, issues),
    target: string(item.target, `${path}.target`, issues),
    ...(item.label === undefined ? {} : { label: string(item.label, `${path}.label`, issues) }),
    status: oneOf(item.status, risks, `${path}.status`, issues),
  };
}

function entity(value: unknown, path: string, issues: string[]): SandboxEntity {
  const item = object(value, path, issues);
  const attributes = object(item.attributes, `${path}.attributes`, issues);
  for (const [key, attribute] of Object.entries(attributes)) {
    if (!["string", "number", "boolean"].includes(typeof attribute) && attribute !== null) {
      issues.push(`${path}.attributes.${key} must be a scalar`);
    }
  }
  return {
    id: string(item.id, `${path}.id`, issues),
    type: string(item.type, `${path}.type`, issues),
    businessCode: string(item.businessCode, `${path}.businessCode`, issues),
    name: string(item.name, `${path}.name`, issues),
    status: string(item.status, `${path}.status`, issues),
    risk: oneOf(item.risk, risks, `${path}.risk`, issues),
    attributes: attributes as SandboxEntity["attributes"],
  };
}

function delta(value: unknown, path: string, issues: string[]): VisualDelta {
  const item = object(value, path, issues);
  return {
    nodeStatuses: array(item.nodeStatuses ?? [], `${path}.nodeStatuses`, issues).map((entry, index) => {
      const update = object(entry, `${path}.nodeStatuses[${index}]`, issues);
      return { id: string(update.id, `${path}.nodeStatuses[${index}].id`, issues), status: oneOf(update.status, risks, `${path}.nodeStatuses[${index}].status`, issues) };
    }),
    edgeStatuses: array(item.edgeStatuses ?? [], `${path}.edgeStatuses`, issues).map((entry, index) => {
      const update = object(entry, `${path}.edgeStatuses[${index}]`, issues);
      return { id: string(update.id, `${path}.edgeStatuses[${index}].id`, issues), status: oneOf(update.status, risks, `${path}.edgeStatuses[${index}].status`, issues) };
    }),
    entityUpdates: array(item.entityUpdates ?? [], `${path}.entityUpdates`, issues).map((entry, index) => {
      const update = object(entry, `${path}.entityUpdates[${index}]`, issues);
      const attributes = update.attributes === undefined ? undefined : object(update.attributes, `${path}.entityUpdates[${index}].attributes`, issues);
      return {
        id: string(update.id, `${path}.entityUpdates[${index}].id`, issues),
        ...(update.status === undefined ? {} : { status: string(update.status, `${path}.entityUpdates[${index}].status`, issues) }),
        ...(update.risk === undefined ? {} : { risk: oneOf(update.risk, risks, `${path}.entityUpdates[${index}].risk`, issues) }),
        ...(attributes === undefined ? {} : { attributes: attributes as SandboxEntity["attributes"] }),
      };
    }),
  };
}

function event(value: unknown, index: number, issues: string[]): ScenarioEvent {
  const path = `timeline[${index}]`;
  const item = object(value, path, issues);
  return {
    sequence: number(item.sequence, `${path}.sequence`, issues),
    id: string(item.id, `${path}.id`, issues),
    timestamp: rfc3339(item.timestamp, `${path}.timestamp`, issues),
    eventType: string(item.eventType, `${path}.eventType`, issues),
    title: string(item.title, `${path}.title`, issues),
    description: string(item.description, `${path}.description`, issues),
    act: number(item.act, `${path}.act`, issues),
    domain: oneOf(item.domain, domains, `${path}.domain`, issues),
    severity: oneOf(item.severity, risks, `${path}.severity`, issues),
    relatedEntityIds: array(item.relatedEntityIds, `${path}.relatedEntityIds`, issues).map((id, relatedIndex) => string(id, `${path}.relatedEntityIds[${relatedIndex}]`, issues)),
    delta: delta(item.delta, `${path}.delta`, issues),
    kpis: kpis(item.kpis, `${path}.kpis`, issues),
  };
}

function act(value: unknown, index: number, issues: string[]): ScenarioAct {
  const path = `acts[${index}]`;
  const item = object(value, path, issues);
  const range = array(item.eventRange, `${path}.eventRange`, issues);
  if (range.length !== 2) issues.push(`${path}.eventRange must contain two sequence numbers`);
  return { number: number(item.number, `${path}.number`, issues), title: string(item.title, `${path}.title`, issues), eventRange: [number(range[0], `${path}.eventRange[0]`, issues), number(range[1], `${path}.eventRange[1]`, issues)] };
}

function agentOutput(value: unknown, index: number, issues: string[]): AgentOutput {
  const path = `agentOutputs[${index}]`;
  const item = object(value, path, issues);
  const confidence = number(item.confidence, `${path}.confidence`, issues);
  if (confidence < 0 || confidence > 1) issues.push(`${path}.confidence must be between 0 and 1`);
  if (typeof item.requiresHumanConfirmation !== "boolean") issues.push(`${path}.requiresHumanConfirmation must be a boolean`);
  return {
    id: string(item.id, `${path}.id`, issues),
    eventSequence: number(item.eventSequence, `${path}.eventSequence`, issues),
    agent: oneOf(item.agent, agents, `${path}.agent`, issues),
    title: string(item.title, `${path}.title`, issues),
    recommendation: string(item.recommendation, `${path}.recommendation`, issues),
    evidence: array(item.evidence, `${path}.evidence`, issues).map((entry, evidenceIndex) => string(entry, `${path}.evidence[${evidenceIndex}]`, issues)),
    impact: string(item.impact, `${path}.impact`, issues),
    confidence,
    status: oneOf(item.status, ["suggested", "executed"] as const, `${path}.status`, issues),
    requiresHumanConfirmation: item.requiresHumanConfirmation === true,
  };
}

export function parseSandboxScenario(value: unknown): SandboxScenario {
  const issues: string[] = [];
  const root = object(value, "$", issues);
  const layout = object(root.layout, "layout", issues);
  const initialState = object(root.initialState, "initialState", issues);
  const parsed: SandboxScenario = {
    key: string(root.key, "key", issues),
    name: string(root.name, "name", issues),
    version: string(root.version, "version", issues),
    dataSource: oneOf(root.dataSource, ["preview", "iaos"] as const, "dataSource", issues),
    timezone: oneOf(root.timezone, ["Asia/Shanghai"] as const, "timezone", issues),
    startsAt: rfc3339(root.startsAt, "startsAt", issues),
    endsAt: rfc3339(root.endsAt, "endsAt", issues),
    defaultSpeed: oneOf(root.defaultSpeed, [1, 2, 4] as const, "defaultSpeed", issues),
    acts: array(root.acts, "acts", issues).map((entry, index) => act(entry, index, issues)),
    layout: {
      width: number(layout.width, "layout.width", issues),
      height: number(layout.height, "layout.height", issues),
      nodes: array(layout.nodes, "layout.nodes", issues).map((entry, index) => node(entry, `layout.nodes[${index}]`, issues)),
      edges: array(layout.edges, "layout.edges", issues).map((entry, index) => edge(entry, `layout.edges[${index}]`, issues)),
    },
    initialState: {
      entities: array(initialState.entities, "initialState.entities", issues).map((entry, index) => entity(entry, `initialState.entities[${index}]`, issues)),
      kpis: kpis(initialState.kpis, "initialState.kpis", issues),
    },
    timeline: array(root.timeline, "timeline", issues).map((entry, index) => event(entry, index, issues)),
    agentOutputs: array(root.agentOutputs, "agentOutputs", issues).map((entry, index) => agentOutput(entry, index, issues)),
  };

  const nodeIds = new Set(parsed.layout.nodes.map(({ id }) => id));
  const edgeIds = new Set(parsed.layout.edges.map(({ id }) => id));
  const entityIds = new Set(parsed.initialState.entities.map(({ id }) => id));
  if (nodeIds.size !== parsed.layout.nodes.length) issues.push("layout.nodes ids must be unique");
  if (edgeIds.size !== parsed.layout.edges.length) issues.push("layout.edges ids must be unique");
  if (entityIds.size !== parsed.initialState.entities.length) issues.push("initialState.entities ids must be unique");
  if (new Set(parsed.timeline.map(({ id }) => id)).size !== parsed.timeline.length) issues.push("timeline event ids must be unique");
  if (parsed.acts.length !== 7) issues.push(`acts must contain exactly 7 acts (received ${parsed.acts.length})`);
  if (parsed.timeline.length !== 22) issues.push(`timeline must contain exactly 22 events (received ${parsed.timeline.length})`);
  parsed.timeline.forEach((item, index) => {
    if (item.sequence !== index + 1) issues.push(`timeline[${index}].sequence must be ${index + 1}`);
    if (index > 0 && item.timestamp <= parsed.timeline[index - 1].timestamp) issues.push(`timeline[${index}].timestamp must be later than the previous event`);
    if (!parsed.acts.some(({ number: actNumber, eventRange }) => actNumber === item.act && item.sequence >= eventRange[0] && item.sequence <= eventRange[1])) issues.push(`timeline[${index}] is not covered by act ${item.act}`);
    item.relatedEntityIds.forEach((id) => { if (!entityIds.has(id)) issues.push(`timeline[${index}] references unknown entity ${id}`); });
    item.delta.nodeStatuses?.forEach(({ id }) => { if (!nodeIds.has(id)) issues.push(`timeline[${index}] references unknown node ${id}`); });
    item.delta.edgeStatuses?.forEach(({ id }) => { if (!edgeIds.has(id)) issues.push(`timeline[${index}] references unknown edge ${id}`); });
    item.delta.entityUpdates?.forEach(({ id }) => { if (!entityIds.has(id)) issues.push(`timeline[${index}] updates unknown entity ${id}`); });
  });
  parsed.layout.edges.forEach((item, index) => {
    if (!nodeIds.has(item.source) || !nodeIds.has(item.target)) issues.push(`layout.edges[${index}] references an unknown node`);
  });
  parsed.agentOutputs.forEach((output, index) => {
    if (output.eventSequence < 1 || output.eventSequence > parsed.timeline.length) issues.push(`agentOutputs[${index}].eventSequence is outside the timeline`);
  });

  if (issues.length > 0) throw new ScenarioValidationError(issues);
  return parsed;
}
