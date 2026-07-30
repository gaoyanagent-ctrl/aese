export type GameMoney = { value: string; currency: "CNY"; scale: 2 };
export type GameBuilding = {
  code: string;
  kind: string;
  label: string;
  state: string;
  x: number;
  y: number;
  evidence_ref?: string;
  available: boolean;
};
export type GameActor = {
  actor_id: string;
  actor_type: "human" | "agent" | "runtime" | "world";
  display_name: string;
  position: string;
  state: string;
  work_item_id?: string;
};
export type ApprovalReview = {
  document_type: string;
  title: string;
  summary: string;
  prepared_by: string;
  status: string;
  fields: Array<{ label: string; value: string }>;
  risks: string[];
  approval_effect: string;
  evidence_ref: string;
};
export type GameWorkItem = {
  work_item_id: string;
  title: string;
  kind: "human_task" | "agent_task" | "approval" | "world_wait" | "capability";
  status: string;
  owner_type: string;
  owner_id: string;
  capability: string;
  requires_me: boolean;
  evidence_ref: string;
  gate?: string;
  approval_review?: ApprovalReview;
};
export type GameProjection = {
  schema_version: "1.0";
  projection_id: string;
  tenant_id: string;
  case_code: string;
  world_run_id: string;
  chapter: string;
  sim_time: string;
  time_scale: 0 | 1 | 2 | 4;
  paused: boolean;
  cursor: number;
  world_scene: { scene_id: string; mode: "2d" | "2.5d"; theme: string };
  lifecycle: {
    state: string;
    current_step: string;
    progress: number;
    blocked_by?: string;
  };
  buildings: GameBuilding[];
  actors: GameActor[];
  work_items: GameWorkItem[];
  resources: {
    founder_cash: GameMoney;
    company_cash: GameMoney;
    capital_committed: GameMoney;
    capital_paid: GameMoney;
    budget_authorized: GameMoney;
    risk_level: string;
  };
  finance_opening: {
    ready: boolean;
    organization_code?: string;
    organization_status?: string;
    roles: string[];
    book_code?: string;
    book_name?: string;
    fiscal_year?: number;
    accounting_periods?: Array<{
      period_code: string;
      starts_on: string;
      ends_on: string;
      status: string;
    }>;
    accounting_standard?: string;
    functional_currency?: string;
    period_code?: string;
    period_status?: string;
    journal_entry_no?: string;
    journal_status?: string;
    debit_minor: number;
    credit_minor: number;
    trial_balance: Array<{
      account_code: string;
      account_name: string;
      account_class: string;
      normal_balance: string;
      debit_minor: number;
      credit_minor: number;
      balance_minor: number;
    }>;
    bank_journal: Array<{
      entry_no: string;
      business_date: string;
      description: string;
      debit_minor: number;
      credit_minor: number;
      balance_minor: number;
      source_type: string;
      source_ref: string;
      evidence_ref: string;
    }>;
    general_ledger: Array<{
      account_code: string;
      account_name: string;
      opening_balance_minor: number;
      debit_minor: number;
      credit_minor: number;
      closing_balance_minor: number;
    }>;
    opening_balance_sheet: {
      as_of: string;
      currency: string;
      assets: Array<{
        account_code: string;
        account_name: string;
        amount_minor: number;
      }>;
      liabilities: Array<{
        account_code: string;
        account_name: string;
        amount_minor: number;
      }>;
      equity: Array<{
        account_code: string;
        account_name: string;
        amount_minor: number;
      }>;
      total_assets_minor: number;
      total_liabilities_minor: number;
      total_equity_minor: number;
      balanced: boolean;
    };
    evidence_ref?: string;
  };
  exchanges: Array<{
    exchange_id: string;
    kind: string;
    status: string;
    correlation_id: string;
    evidence_ref: string;
    occurred_at: string;
  }>;
  brand: {
    status: string;
    company_name?: string;
    logo_asset_id?: string;
    primary_color?: string;
  };
  notifications: Array<{
    notification_id: string;
    severity: string;
    message: string;
    action_ref?: string;
  }>;
  evidence_refs: Array<{ ref: string; kind: string }>;
};
export type FounderIntentRequest = {
  tenant_id: string;
  case_code: string;
  raw_idea: string;
  industry: string;
  customers: string[];
  offerings: string[];
  brand_traits: string[];
  capital_minor: string;
  risk_appetite: "conservative" | "balanced" | "aggressive";
};
export type FounderIntent = FounderIntentRequest & {
  schema_version: "1.0";
  intent_id: string;
  assumptions: string[];
  needs_confirmation: string[];
  created_at: string;
};
export type NamingProposal = {
  proposal_id: string;
  chinese_name: string;
  english_name: string;
  short_name: string;
  rationale: string;
  slogan: string;
  keywords: string[];
  primary_color: string;
  risk_hints: string[];
  status: "candidate";
};
export type GenesisWorkspace = {
  workspace_id: string;
  owner_player_id: string;
  display_name: string;
  tenant_id: string;
  world_run_id: string;
  case_code: string;
  status: "provisioning" | "awaiting_world" | "active" | "failed";
  current_step: string;
  last_error?: string;
  created_at: string;
  updated_at: string;
  correlation_id?: string;
  attempt?: number;
  evidence_refs?: Record<string, string>;
  steps?: Array<{
    step_key: string;
    status: string;
    attempt: number;
    evidence_ref?: string;
    last_error?: string;
    updated_at: string;
  }>;
};
export type GenesisWorkspaceResult = GenesisWorkspace & {
  tenant_token: string;
};
