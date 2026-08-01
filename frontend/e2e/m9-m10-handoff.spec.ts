import { expect, test } from "@playwright/test";

const money = { value: "10000000000", currency: "CNY", scale: 2 };

test("M9 terminal offers a contextual handoff to M10", async ({ page }) => {
  await page.addInitScript(() => {
    localStorage.setItem("iaos_token", "test-founder-token");
    localStorage.setItem("aese_iaos_tenant_id", "tenant-gx-handoff");
    localStorage.setItem("aese_genesis_case_code", "INC-GX-HANDOFF");
    localStorage.setItem("aese_genesis_workspace_id", "gxw-handoff");
  });
  await page.route("**/api/aese/v1/game/incorporation/INC-GX-HANDOFF/projection?frame=0", (route) =>
    route.fulfill({
      status: 200,
      json: {
        schema_version: "1.0",
        projection_id: "projection-handoff",
        tenant_id: "tenant-gx-handoff",
        case_code: "INC-GX-HANDOFF",
        world_run_id: "world-run-handoff",
        chapter: "enterprise_operational_ready",
        sim_time: "2026-08-01T00:00:00Z",
        time_scale: 1,
        paused: false,
        cursor: 23,
        world_scene: { scene_id: "genesis-city", mode: "2.5d", theme: "genesis" },
        lifecycle: { state: "enterprise_operational_ready", current_step: "enterprise_operational_ready", progress: 100 },
        buildings: [],
        actors: [],
        work_items: [{ work_item_id: "wi-23", title: "完成企业开业检查", kind: "capability", status: "completed", owner_type: "human", owner_id: "founder-principal", capability: "enterprise.readiness.evaluate", requires_me: false, evidence_ref: "iaos:wi-23" }],
        resources: { founder_cash: { ...money, value: "0" }, company_cash: money, capital_committed: money, capital_paid: money, budget_authorized: money, risk_level: "normal" },
        finance_opening: { ready: false, roles: [], debit_minor: 0, credit_minor: 0, trial_balance: [], bank_journal: [], general_ledger: [], opening_balance_sheet: { as_of: "2026-08-01", currency: "CNY", assets: [], liabilities: [], equity: [], total_assets_minor: 0, total_liabilities_minor: 0, total_equity_minor: 0, balanced: true } },
        exchanges: [],
        brand: { status: "selected", company_name: "交接测试制造企业" },
        notifications: [],
        evidence_refs: [],
      },
    }),
  );

  await page.goto("/#enterprise-genesis?tenant=tenant-gx-handoff&case=INC-GX-HANDOFF&workspace=gxw-handoff");
  await expect(page.getByRole("heading", { name: "M9 企业创生已完成" })).toBeVisible();
  await page.getByRole("button", { name: "开始 M10 工厂选址与设施规划" }).click();
  await expect(page).toHaveURL(/#world-plant-build\?tenant=tenant-gx-handoff&case=INC-GX-HANDOFF&from=enterprise-genesis&workspace=gxw-handoff$/);
});
