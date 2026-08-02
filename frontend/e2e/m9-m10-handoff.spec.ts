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

test("M10 opens as a playable enterprise scene instead of a standalone form", async ({ page }) => {
  await page.addInitScript(() => {
    localStorage.setItem("iaos_token", "test-founder-token");
    localStorage.setItem("aese_iaos_tenant_id", "tenant-gx-handoff");
    localStorage.setItem("aese_genesis_case_code", "INC-GX-HANDOFF");
  });
  await page.route("**/api/aese/v1/world/plant-build**", (route) => {
    const url = new URL(route.request().url());
    if (url.pathname.endsWith("/planning-status")) return route.fulfill({ json: { state: "not_configured", provider: "none", prompt_version: "plant-planning-v2" } });
    if (url.pathname.endsWith("/financial-constraints")) return route.fulfill({ json: { case_code: "INC-GX-HANDOFF", legal_entity_code: "LE-GX-HANDOFF", financial_constraint: { available_cash: { value: "20000000.00", currency: "CNY", scale: 2 }, approved_budget: { value: "15000000.00", currency: "CNY", scale: 2 }, cash_source_ref: "gl:BOOK:1002", budget_source_ref: "budget:BUD-1", snapshot_hash: "sha256:test" } } });
    if (url.pathname.includes("/requirements/") || url.pathname.endsWith("/proposals")) return route.fulfill({ status: 404, body: "not found" });
    if (url.pathname.endsWith("/investigations") || url.pathname.endsWith("/site-selections")) return route.fulfill({ json: { items: [] } });
    return route.fulfill({ json: { schema_version: "1.0", campaign: "plant-build", world_run_id: "run-1", timezone: "Asia/Shanghai", policy_version: "v1", m9_terminal_hash: "hash", frames: [{ step: 0, phase: "requirement", sim_time: "2026-08-01T08:00:00Z", title: "定义设施需求", causation_id: "cause-1", selected_site: "", assessments: [], zones: [], work_packages: [], utilities: {}, knowledge: [], world_progress: 0, iaos_plan_progress: 0, discrepancy: "open", cash: { value: "0.00", currency: "CNY", scale: 2 }, committed: { value: "0.00", currency: "CNY", scale: 2 }, payable: { value: "0.00", currency: "CNY", scale: 2 }, paid: { value: "0.00", currency: "CNY", scale: 2 }, capability_build_eligible: false, iaos_cursor: 0 }] } });
  });

  await page.goto("/#world-plant-build?tenant=tenant-gx-handoff&case=INC-GX-HANDOFF&from=enterprise-genesis");
  await expect(page.getByRole("heading", { name: "在总部规划室明确工厂需要什么" })).toBeVisible();
  await expect(page.getByRole("navigation", { name: "M10 可进入地点" })).toBeVisible();
  await expect(page.getByRole("button", { name: /总部规划中心/ })).toHaveAttribute("aria-pressed", "true");
  await expect(page.getByAltText("你的创始人角色")).toBeVisible();
  await expect(page.getByText("纪元", { exact: true })).toBeVisible();
  await expect(page.getByRole("dialog", { name: "当前经营任务" })).toHaveCount(0);
  await expect(page.getByLabel("目标区域")).toHaveCount(0);
  await page.getByRole("button", { name: /纪元：与规划 Agent 制定需求/ }).click();
  await expect(page.getByRole("dialog", { name: "当前经营任务" })).toBeVisible();
  await expect(page.getByLabel("目标区域")).toBeVisible();
  await page.getByRole("button", { name: "关闭当前任务" }).click();
  await expect(page.getByRole("dialog", { name: "当前经营任务" })).toHaveCount(0);
  expect(await page.evaluate(() => document.documentElement.scrollWidth <= window.innerWidth)).toBe(true);
});

test("M10 explains a failed site gate with requirement, observation and difference", async ({ page }) => {
  await page.addInitScript(() => {
    localStorage.setItem("iaos_token", "test-founder-token");
    localStorage.setItem("aese_iaos_tenant_id", "tenant-gx-handoff");
    localStorage.setItem("aese_genesis_case_code", "INC-GX-HANDOFF");
  });
  const requirement = {
    schema_version: "1.0", requirement_id: "facility-requirement-INC-GX-HANDOFF", tenant_id: "tenant-gx-handoff", case_code: "INC-GX-HANDOFF", legal_entity_code: "LE-GX-HANDOFF",
    target_region: "苏州", facility_purpose: "制造基地", minimum_area_m2: 8000, minimum_electricity_kva: 2000, target_available_at: "2026-12-01T00:00:00Z",
    candidate_count: 2, allowed_option_types: ["lease_and_retrofit"], investment_request: { value: "10000000.00", currency: "CNY", scale: 2 }, minimum_cash_reserve: { value: "2000000.00", currency: "CNY", scale: 2 },
    financial_constraint: { available_cash: { value: "20000000.00", currency: "CNY", scale: 2 }, approved_budget: { value: "15000000.00", currency: "CNY", scale: 2 }, cash_source_ref: "gl:BOOK:1002", budget_source_ref: "budget:BUD-1", snapshot_hash: "sha256:test" }, preferences: [], revision: 1, revision_reason: "test",
  };
  const proposalSet = {
    schema_version: "1.0", proposal_set_id: "SET-1", requirement_id: requirement.requirement_id, revision: 1, status: "candidate_only",
    proposals: [{ proposal_id: "SITE-1", option_type: "lease_and_retrofit", display_name: "候选场址一", business_rationale: "测试", estimated_amount: { minimum: { value: "7000000.00", currency: "CNY", scale: 2 }, likely: { value: "8000000.00", currency: "CNY", scale: 2 }, maximum: { value: "9000000.00", currency: "CNY", scale: 2 }, basis: "测试" }, estimated_schedule: { earliest: "2026-09-01T00:00:00Z", likely: "2026-10-01T00:00:00Z", latest: "2026-11-01T00:00:00Z" }, assumptions: [], facts_required: [], risks: [], source_refs: [], confidence: "0.8", status: "proposed" }],
    evidence: { provider: "MiniMax", model: "MiniMax-M3", prompt_version: "plant-planning-v2", input_hash: "sha256:in", output_hash: "sha256:out", validated_at: "2026-08-01T00:00:00Z" },
  };
  await page.route("**/api/aese/v1/world/plant-build**", (route) => {
    const url = new URL(route.request().url());
    if (url.pathname.endsWith("/planning-status")) return route.fulfill({ json: { state: "connected", provider: "MiniMax", model: "MiniMax-M3", prompt_version: "plant-planning-v2" } });
    if (url.pathname.endsWith("/financial-constraints")) return route.fulfill({ json: { case_code: "INC-GX-HANDOFF", legal_entity_code: "LE-GX-HANDOFF", financial_constraint: requirement.financial_constraint } });
    if (url.pathname.includes("/requirements/")) return route.fulfill({ json: requirement });
    if (url.pathname.endsWith("/proposals")) return route.fulfill({ json: proposalSet });
    if (url.pathname.endsWith("/investigations")) return route.fulfill({ json: { items: [{ request: { investigation_request_id: "INV-1", proposal_set_id: "SET-1", proposal_id: "SITE-1", expected_revision: 1, scope: [] }, status: "observed", work_item_status: "completed", observation: { schema_version: "1.0", observation_id: "OBS-1", investigation_request_id: "INV-1", proposal_id: "SITE-1", result: "completed", ownership_status: "verified", available_area_m2: 9000, electricity_kva: 1500, quoted_amount: { value: "8000000.00", currency: "CNY", scale: 2 }, available_at: "2026-10-01T00:00:00Z", permit_status: "eligible", evidence_refs: ["world-document:quote-1"], notes: "", external_actor_id: "park", observed_at: "2026-08-01T00:00:00Z" } }] } });
    if (url.pathname.endsWith("/site-selections")) return route.fulfill({ json: { items: [] } });
    return route.fulfill({ status: 404, body: "not found" });
  });

  await page.goto("/#world-plant-build?tenant=tenant-gx-handoff&case=INC-GX-HANDOFF");
  await page.getByRole("button", { name: /周衡：比较事实并提交推荐/ }).click();
  await expect(page.getByText("最低要求 2,000 kVA")).toBeVisible();
  await expect(page.getByText("实测 1,500 kVA")).toBeVisible();
  await expect(page.getByText("短缺 500 kVA")).toBeVisible();
  await expect(page.getByText(/当前权重来源：界面默认比较偏好/)).toBeVisible();
  await expect(page.getByLabel("成本权重")).toBeDisabled();
  await page.getByRole("button", { name: "修订设施需求" }).click();
  await expect(page.getByText("修订设施需求 · 将保存为第 2 版")).toBeVisible();
  await expect(page.getByLabel("目标区域")).toHaveValue("苏州");
  await expect(page.getByLabel("最小电力容量（kVA）")).toHaveValue("2000");
  await expect(page.getByLabel("本次修订原因")).toHaveValue("");
  await page.getByRole("button", { name: "返回外部事实比较" }).click();
  await expect(page.getByText("短缺 500 kVA")).toBeVisible();
  await page.getByText("这些分数具体怎样计算？").click();
  await expect(page.getByText(/目标日期当天可用为 50 分/)).toBeVisible();
});
