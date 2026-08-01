import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeEach, describe, expect, it, vi } from "vitest";
import type { GameProjection } from "../../game/types";
import { EnterpriseGenesisGame } from "./EnterpriseGenesisGame";

const { completeProjection } = vi.hoisted(() => ({ completeProjection: {
  schema_version: "1.0",
  projection_id: "projection-complete",
  tenant_id: "tenant-gx-test",
  case_code: "INC-GX-TEST",
  world_run_id: "world-run-test",
  chapter: "enterprise_operational_ready",
  sim_time: "2026-08-01T00:00:00Z",
  time_scale: 1,
  paused: false,
  cursor: 23,
  world_scene: { scene_id: "genesis-city", mode: "2.5d", theme: "genesis" },
  lifecycle: {
    state: "enterprise_operational_ready",
    current_step: "enterprise_operational_ready",
    progress: 100,
  },
  buildings: [],
  actors: [],
  work_items: [
    {
      work_item_id: "wi-23",
      title: "完成企业开业检查",
      kind: "capability",
      status: "completed",
      owner_type: "human",
      owner_id: "founder-principal",
      capability: "enterprise.readiness.evaluate",
      requires_me: false,
      evidence_ref: "iaos:incorporation:INC-GX-TEST:work-item:23",
    },
  ],
  resources: {
    founder_cash: { value: "0", currency: "CNY", scale: 2 },
    company_cash: { value: "10000000000", currency: "CNY", scale: 2 },
    capital_committed: { value: "10000000000", currency: "CNY", scale: 2 },
    capital_paid: { value: "10000000000", currency: "CNY", scale: 2 },
    budget_authorized: { value: "10000000000", currency: "CNY", scale: 2 },
    risk_level: "normal",
  },
  finance_opening: {
    ready: false,
    roles: [],
    debit_minor: 0,
    credit_minor: 0,
    trial_balance: [],
    bank_journal: [],
    general_ledger: [],
    opening_balance_sheet: {
      as_of: "2026-08-01",
      currency: "CNY",
      assets: [],
      liabilities: [],
      equity: [],
      total_assets_minor: 0,
      total_liabilities_minor: 0,
      total_equity_minor: 0,
      balanced: true,
    },
  },
  exchanges: [],
  brand: { status: "selected", company_name: "测试制造企业" },
  notifications: [],
  evidence_refs: [],
} as GameProjection }));

vi.mock("../../game/api", async () => {
  const actual = await vi.importActual<typeof import("../../game/api")>("../../game/api");
  return { ...actual, loadGameProjection: vi.fn().mockResolvedValue(completeProjection) };
});

describe("EnterpriseGenesisGame milestone handoff", () => {
  beforeEach(() => {
    window.location.hash =
      "enterprise-genesis?tenant=tenant-gx-test&case=INC-GX-TEST&workspace=gxw-test";
  });

  it("offers an explicit M10 handoff after the governed M9 terminal", async () => {
    const onContinueToPlantBuild = vi.fn();
    const user = userEvent.setup();
    render(
      <EnterpriseGenesisGame
        onExit={() => undefined}
        onContinueToPlantBuild={onContinueToPlantBuild}
      />,
    );

    expect(
      await screen.findByRole("heading", { name: "M9 企业创生已完成" }),
    ).toBeInTheDocument();
    await user.click(
      screen.getByRole("button", { name: "开始 M10 工厂选址与设施规划" }),
    );
    expect(onContinueToPlantBuild).toHaveBeenCalledWith({
      tenantId: "tenant-gx-test",
      caseCode: "INC-GX-TEST",
      workspaceId: "gxw-test",
    });
  });
});
