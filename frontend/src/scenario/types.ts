export type RiskLevel = "normal" | "watch" | "critical";

export type SandboxNodeKind =
  | "supplier"
  | "warehouse"
  | "quality"
  | "process"
  | "shipping"
  | "customer";

export interface SandboxNode {
  id: string;
  businessCode: string;
  label: string;
  kind: SandboxNodeKind;
  position: { x: number; y: number };
  status: RiskLevel;
  entityId?: string;
}

export interface SandboxEdge {
  id: string;
  source: string;
  target: string;
  label?: string;
  status: RiskLevel;
}

export type EntityAttribute = string | number | boolean | null;

export interface SandboxEntity {
  id: string;
  type: string;
  businessCode: string;
  name: string;
  status: string;
  risk: RiskLevel;
  attributes: Record<string, EntityAttribute>;
}

export interface KpiMetric {
  value: number;
  unit: string;
  risk: RiskLevel;
  trend: "up" | "down" | "flat";
}

export interface KpiSnapshot {
  orderDemand: KpiMetric;
  availableFinishedGoods: KpiMetric;
  materialShortageRisk: KpiMetric;
  capacityRisk: KpiMetric;
  deliveryRisk: KpiMetric;
}

export interface EntityDelta {
  id: string;
  status?: string;
  risk?: RiskLevel;
  attributes?: Record<string, EntityAttribute>;
}

export interface VisualDelta {
  nodeStatuses?: Array<{ id: string; status: RiskLevel }>;
  edgeStatuses?: Array<{ id: string; status: RiskLevel }>;
  entityUpdates?: EntityDelta[];
}

export interface ScenarioAct {
  number: number;
  title: string;
  eventRange: [number, number];
}

export interface ScenarioEvent {
  sequence: number;
  id: string;
  timestamp: string;
  eventType: string;
  title: string;
  description: string;
  act: number;
  domain: "order" | "planning" | "supply" | "equipment" | "quality" | "production" | "logistics";
  severity: RiskLevel;
  relatedEntityIds: string[];
  delta: VisualDelta;
  kpis: KpiSnapshot;
}

export type AgentKind = "planning" | "quality" | "business_analysis";

export interface AgentOutput {
  id: string;
  eventSequence: number;
  agent: AgentKind;
  title: string;
  recommendation: string;
  evidence: string[];
  impact: string;
  confidence: number;
  status: "suggested" | "executed";
  requiresHumanConfirmation: boolean;
}

export interface SandboxScenario {
  key: string;
  name: string;
  version: string;
  dataSource: "preview" | "iaos";
  timezone: "Asia/Shanghai";
  startsAt: string;
  endsAt: string;
  defaultSpeed: 1 | 2 | 4;
  acts: ScenarioAct[];
  layout: {
    width: number;
    height: number;
    nodes: SandboxNode[];
    edges: SandboxEdge[];
  };
  initialState: {
    entities: SandboxEntity[];
    kpis: KpiSnapshot;
  };
  timeline: ScenarioEvent[];
  agentOutputs: AgentOutput[];
}

export interface IaosMetric { value: number; unit: string }
export interface IaosScenarioEvent {
  cursor: number;
  event_id: string;
  event_type: string;
  occurred_at: string;
  correlation_id: string;
  business_object_type: string;
  business_object_code: string;
  payload: Record<string, unknown>;
}
export interface IaosScenarioEntity {
  id: string;
  type: string;
  business_code: string;
  name: string;
  status: string;
  attributes: Record<string, EntityAttribute>;
}
export interface IaosRecommendation {
  agent_key: AgentKind;
  summary: string;
  recommendations: string[];
  object_refs: string[];
  tool_call_ids: string[];
  completeness: string;
  data_gaps?: string[];
  confidence: string;
  status: 'suggested';
  requires_human_confirmation: boolean;
  version?: number;
  run_id?: string;
  updated_at?: string;
}
export interface IaosScenarioSnapshot {
  snapshot_version: string;
  pack_key: string;
  scenario_key: string;
  observed_at: string;
  cursor: number;
  completeness: string;
  gaps: string[];
  entities: IaosScenarioEntity[];
  kpis: {
    order_demand: IaosMetric;
    cumulative_available: IaosMetric;
    cumulative_shipped: IaosMetric;
    ending_finished_goods: IaosMetric;
    delivery_gap: IaosMetric;
  };
  events: IaosScenarioEvent[];
  recommendations: IaosRecommendation[];
}
