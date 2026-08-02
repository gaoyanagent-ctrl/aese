import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { PlantBuildPlay } from "./PlantBuildPlay";

const trace = {
  schema_version: "1.0",
  campaign: "plant-build",
  world_run_id: "run-1",
  timezone: "Asia/Shanghai",
  policy_version: "v1",
  m9_terminal_hash: "hash",
  frames: [{
    step: 0,
    phase: "requirement",
    sim_time: "2026-08-01T08:00:00Z",
    title: "消费 M9 机器资格",
    causation_id: "cause-1",
    selected_site: "",
    assessments: [], zones: [], work_packages: [], utilities: {}, knowledge: [],
    world_progress: 0, iaos_plan_progress: 0, discrepancy: "open",
    cash: { value: "0.00", currency: "CNY", scale: 2 },
    committed: { value: "0.00", currency: "CNY", scale: 2 },
    payable: { value: "0.00", currency: "CNY", scale: 2 },
    paid: { value: "0.00", currency: "CNY", scale: 2 },
    capability_build_eligible: false,
    iaos_cursor: 0,
  }],
};

describe("PlantBuildPlay interactive planning", () => {
  beforeEach(() => {
    window.location.hash = "";
    localStorage.clear();
    localStorage.setItem("aese_iaos_tenant_id", "tenant-gx-test");
    localStorage.setItem("aese_genesis_case_code", "INC-GX-TEST");
    localStorage.setItem("iaos_token", "test-token");
  });
  afterEach(() => vi.unstubAllGlobals());

  const openNpcTask = async (name: RegExp) => {
    fireEvent.click(await screen.findByRole("button", { name }));
    await screen.findByRole("dialog", { name: "当前经营任务" });
  };

  it("fails closed when the external planning model is unavailable", async () => {
    vi.stubGlobal("fetch", vi.fn(async (input: RequestInfo | URL) => {
      const path = String(input);
      if (path.includes("/requirements/")) return { ok: false, status: 404, text: async () => "not found" };
      if (path.includes("/proposals?")) return { ok: false, status: 404, text: async () => "not found" };
      return {
        ok: true,
        json: async () => path.endsWith("planning-status")
          ? { state: "not_configured", provider: "none", prompt_version: "plant-planning-v1" }
          : path.includes("financial-constraints")
            ? { case_code: "INC-GX-TEST", legal_entity_code: "LE-GX-TEST", financial_constraint: { available_cash: { value: "20000000.00", currency: "CNY", scale: 2 }, approved_budget: { value: "15000000.00", currency: "CNY", scale: 2 }, cash_source_ref: "gl:BOOK:1002", budget_source_ref: "budget:BUD-1", snapshot_hash: "sha256:authority" } }
            : trace,
      };
    }));
    render(<PlantBuildPlay onExit={() => undefined} />);
    await waitFor(() => expect(screen.getByRole("heading", { name: "在总部规划室明确工厂需要什么" })).toBeInTheDocument());
    expect(screen.getByRole("navigation", { name: "M10 可进入地点" })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: /总部规划中心/ })).toHaveAttribute("aria-pressed", "true");
    expect(screen.queryByRole("dialog", { name: "当前经营任务" })).not.toBeInTheDocument();
    expect(screen.queryByLabelText(/目标区域/)).not.toBeInTheDocument();
    await openNpcTask(/纪元：与规划 Agent 制定需求/);
    await waitFor(() => expect(screen.getByText("未启用外部模型")).toBeInTheDocument());
    expect(screen.getByRole("button", { name: /让 Agent 准备需求方案/ })).toBeDisabled();
    expect(screen.getByText(/不能生成虚拟固定候选/)).toBeInTheDocument();
    expect(screen.getByText(/IAOS 可用现金/)).toBeInTheDocument();
  });

  it("allows a project lead to establish the first manual proposal set", async () => {
    const requirement = {
      schema_version: "1.0", requirement_id: "facility-requirement-INC-GX-TEST", tenant_id: "tenant-gx-test", case_code: "INC-GX-TEST", legal_entity_code: "LE-GX-TEST",
      target_region: "苏州", facility_purpose: "制造基地", minimum_area_m2: 8000, minimum_electricity_kva: 2000, target_available_at: "2026-12-01T00:00:00Z",
      candidate_count: 2, allowed_option_types: ["lease_and_retrofit"], investment_request: { value: "10000000.00", currency: "CNY", scale: 2 }, minimum_cash_reserve: { value: "2000000.00", currency: "CNY", scale: 2 },
      financial_constraint: { available_cash: { value: "20000000.00", currency: "CNY", scale: 2 }, approved_budget: { value: "15000000.00", currency: "CNY", scale: 2 }, cash_source_ref: "gl:BOOK:1002", budget_source_ref: "budget:BUD-1", snapshot_hash: "sha256:authority" },
      preferences: [], revision: 1, revision_reason: "首次规划",
    };
    vi.stubGlobal("fetch", vi.fn(async (input: RequestInfo | URL) => {
      const path = String(input);
      if (path.endsWith("planning-status")) return { ok: true, json: async () => ({ state: "not_configured", provider: "none", prompt_version: "plant-planning-v1" }) };
      if (path.includes("financial-constraints")) return { ok: true, json: async () => ({ case_code: "INC-GX-TEST", legal_entity_code: "LE-GX-TEST", financial_constraint: requirement.financial_constraint }) };
      if (path.includes("/requirements/")) return { ok: true, json: async () => requirement };
      if (path.includes("/proposals?")) return { ok: false, status: 404, text: async () => "not found" };
      if (path.includes("/investigations") || path.includes("/site-selections")) return { ok: true, json: async () => ({ items: [] }) };
      return { ok: true, json: async () => trace };
    }));
    render(<PlantBuildPlay onExit={() => undefined} />);
    await openNpcTask(/纪元：继续填写需求并生成候选/);
    await waitFor(() => expect(screen.getByText("未启用外部模型")).toBeInTheDocument());
    fireEvent.click(screen.getByRole("button", { name: "人工新增候选" }));
    expect(screen.getByText(/建立第 1 版人工候选/)).toBeInTheDocument();
    expect(screen.getByRole("button", { name: /提交到 IAOS 候选集/ })).toBeEnabled();
    fireEvent.click(screen.getByRole("button", { name: "关闭当前任务" }));
    fireEvent.click(screen.getByRole("button", { name: /总部规划中心/ }));
    await screen.findByRole("heading", { name: "企业总部规划中心" });
    fireEvent.click(screen.getByRole("tab", { name: "场景档案" }));
    expect(screen.getByText("设施需求 · 第 1 版")).toBeInTheDocument();
  });

  it("submits editable requirements and renders agent evidence for human review", async () => {
    const fetchMock = vi.fn(async (input: RequestInfo | URL, init?: RequestInit) => {
      const path = String(input);
      if (path.endsWith("planning-status")) return { ok: true, json: async () => ({ state: "connected", provider: "MiniMax", model: "MiniMax-M2.5", prompt_version: "plant-planning-v1" }) };
      if (path.includes("financial-constraints")) return { ok: true, json: async () => ({ case_code: "INC-GX-TEST", legal_entity_code: "LE-GX-TEST", financial_constraint: { available_cash: { value: "20000000.00", currency: "CNY", scale: 2 }, approved_budget: { value: "15000000.00", currency: "CNY", scale: 2 }, cash_source_ref: "gl:BOOK:1002", budget_source_ref: "budget:BUD-1", snapshot_hash: "sha256:authority" } }) };
      if (path.endsWith("requirement-options") && init?.method === "POST") return { ok: true, json: async () => ({ schema_version: "1.0", options: [{ option_id: "requirement-option-1", title: "轻资产快速投产", business_rationale: "平衡上市时间和现金占用", target_region: "江苏省苏州市", facility_purpose: "冷却板制造基地", minimum_area_m2: 12000, minimum_electricity_kva: 2200, target_available_at: "2026-11-01T08:00:00Z", candidate_count: 3, allowed_option_types: ["lease_and_retrofit"], investment_request: { value: "15000000.00", currency: "CNY", scale: 2 }, minimum_cash_reserve: { value: "3000000.00", currency: "CNY", scale: 2 }, preferences: ["快速投产"], tradeoffs: ["扩展空间有限"] }], evidence: { provider: "MiniMax", model: "MiniMax-M2.5", prompt_version: "plant-requirement-adviser-v1", input_hash: "sha256:brief-in", output_hash: "sha256:brief-out", validated_at: "2026-08-01T08:00:00Z" } }) };
      if (path.includes("/requirements/")) return { ok: true, status: 200, json: async () => ({ revision: 3 }) };
      if (path.endsWith("investigations")) return init?.method === "POST"
        ? { ok: true, json: async () => ({ status: "waiting_world", investigation_request: {}, work_item: {} }) }
        : { ok: true, json: async () => ({ items: [] }) };
      if (path.endsWith("reviews")) return { ok: true, json: async () => ({ status: "committed", proposal_review: {} }) };
      if (path.endsWith("proposals/manual") && init?.method === "POST") return { ok: true, status: 201, json: async () => ({ status: "committed", manual_proposal_id: "manual-site-1", proposal_set: {
        schema_version: "1.0", proposal_set_id: "set-1", requirement_id: "facility-requirement-INC-GX-TEST", revision: 2, status: "candidate_only", proposals: [],
        evidence: { provider: "human", model: "manual-entry", prompt_version: "manual-candidate-v1", source_type: "human_manual", parent_revision: 1, input_hash: "sha256:in", output_hash: "sha256:out", validated_at: "2026-08-01T09:00:00Z" },
      } }) };
      if (path.includes("/proposals?") && init?.method !== "POST") return { ok: false, status: 404, text: async () => "not found" };
      if (path.endsWith("proposals") && init?.method === "POST") return { ok: true, json: async () => ({ proposal_set: {
        schema_version: "1.0", proposal_set_id: "set-1", requirement_id: "facility-requirement-INC-GX-TEST", revision: 1, status: "candidate_only",
        proposals: [{ proposal_id: "site-1", option_type: "lease_and_retrofit", display_name: "北区租赁改造建议", business_rationale: "平衡上市时间和初始现金占用", estimated_amount: { minimum: { value: "8000000.00", currency: "CNY", scale: 2 }, likely: { value: "10000000.00", currency: "CNY", scale: 2 }, maximum: { value: "12000000.00", currency: "CNY", scale: 2 }, basis: "需求参数与区域历史区间" }, estimated_schedule: { earliest: "2026-10-01T00:00:00Z", likely: "2026-11-01T00:00:00Z", latest: "2026-12-01T00:00:00Z" }, assumptions: ["存在合规标准厂房"], facts_required: ["核验供电容量"], risks: ["消防改造周期"], source_refs: ["requirement:facility-requirement-INC-GX-TEST"], confidence: "0.78", status: "proposed" }],
        evidence: { provider: "MiniMax", model: "MiniMax-M2.5", prompt_version: "plant-planning-v1", input_hash: "sha256:in", output_hash: "sha256:out", validated_at: "2026-08-01T08:00:00Z" },
      }, agent_job: {}, idempotent_replay: false }) };
      return { ok: true, json: async () => trace };
    });
    vi.stubGlobal("fetch", fetchMock);
    render(<PlantBuildPlay onExit={() => undefined} />);
    await openNpcTask(/纪元：与规划 Agent 制定需求/);
    await waitFor(() => expect(screen.getByText(/MiniMax · MiniMax-M2.5/)).toBeInTheDocument());
    fireEvent.click(screen.getByRole("button", { name: "让 Agent 准备需求方案" }));
    fireEvent.click(await screen.findByRole("button", { name: /轻资产快速投产/ }));
    expect(screen.getByLabelText(/本次投资申请金额/)).toHaveValue(15000000);
    expect(screen.getByLabelText(/目标可用时间/)).toHaveValue("2026-11-01T08:00");
    fireEvent.submit(screen.getByRole("button", { name: /确认草案并生成场址候选/ }).closest("form")!);
    await waitFor(() => expect(screen.getByRole("heading", { name: "北区租赁改造建议" })).toBeInTheDocument());
    expect(screen.getByText("核验供电容量")).toBeInTheDocument();
    fireEvent.click(screen.getByRole("button", { name: "采纳调研" }));
    expect(screen.getByRole("button", { name: "采纳调研" })).toHaveClass("selected");
    fireEvent.change(screen.getByLabelText("审阅理由"), { target: { value: "优先核验供电容量和消防改造周期" } });
    fireEvent.click(screen.getByRole("button", { name: "提交审阅到 IAOS" }));
    await waitFor(() => expect(screen.getByRole("button", { name: "已保存审阅" })).toBeInTheDocument());
    fireEvent.click(screen.getByRole("button", { name: "发起外部调研工作项" }));
    await waitFor(() => expect(fetchMock.mock.calls.some(([path, init]) => String(path).endsWith("investigations") && init?.method === "POST")).toBe(true));
    const request = JSON.parse(fetchMock.mock.calls.find(([path]) => String(path).endsWith("proposals"))?.[1]?.body as string);
    expect(request.investment_request.value).toBe("15000000.00");
    expect(request.revision).toBe(4);
    expect(request.allowed_option_types).toEqual(["lease_and_retrofit"]);
    expect(request.financial_constraint.snapshot_hash).toBe("sha256:authority");
    const review = JSON.parse(fetchMock.mock.calls.find(([path]) => String(path).endsWith("reviews"))?.[1]?.body as string);
    expect(review.action).toBe("adopt_for_investigation");
    expect(review.proposal_set_id).toBe("set-1");
    const investigation = JSON.parse(fetchMock.mock.calls.find(([path, init]) => String(path).endsWith("investigations") && init?.method === "POST")?.[1]?.body as string);
    expect(investigation.proposal_id).toBe("site-1");
    expect(investigation.scope).toContain("commercial_quote");
  });

  it("compares only delivered observations and keeps hard constraints visible", async () => {
    const requirement = {
      schema_version: "1.0", requirement_id: "facility-requirement-INC-GX-TEST", tenant_id: "tenant-gx-test", case_code: "INC-GX-TEST", legal_entity_code: "LE-GX-TEST",
      target_region: "苏州", facility_purpose: "制造基地", minimum_area_m2: 8000, minimum_electricity_kva: 2000, target_available_at: "2026-12-01T00:00:00Z",
      candidate_count: 2, allowed_option_types: ["lease_and_retrofit"], investment_request: { value: "10000000.00", currency: "CNY", scale: 2 }, minimum_cash_reserve: { value: "2000000.00", currency: "CNY", scale: 2 },
      financial_constraint: { available_cash: { value: "20000000.00", currency: "CNY", scale: 2 }, approved_budget: { value: "15000000.00", currency: "CNY", scale: 2 }, cash_source_ref: "gl:BOOK:1002", budget_source_ref: "budget:BUD-1", snapshot_hash: "sha256:authority" },
      preferences: [], revision: 1, revision_reason: "test",
    };
    const currentProposalSet = {
      schema_version: "1.0", proposal_set_id: "SET-1", requirement_id: requirement.requirement_id, revision: 1, status: "candidate_only",
      proposals: [{ proposal_id: "SITE-1", option_type: "lease_and_retrofit", display_name: "候选场址一", business_rationale: "测试", estimated_amount: { minimum: { value: "7000000.00", currency: "CNY", scale: 2 }, likely: { value: "8000000.00", currency: "CNY", scale: 2 }, maximum: { value: "9000000.00", currency: "CNY", scale: 2 }, basis: "测试" }, estimated_schedule: { earliest: "2026-09-01T00:00:00Z", likely: "2026-10-01T00:00:00Z", latest: "2026-11-01T00:00:00Z" }, assumptions: [], facts_required: [], risks: [], source_refs: [], confidence: "0.8", status: "proposed" }],
      evidence: { provider: "MiniMax", model: "MiniMax-M2.5", prompt_version: "plant-planning-v2", input_hash: "sha256:in", output_hash: "sha256:out", validated_at: "2026-08-01T00:00:00Z" },
    };
    vi.stubGlobal("fetch", vi.fn(async (input: RequestInfo | URL) => {
      const path = String(input);
      if (path.endsWith("planning-status")) return { ok: true, json: async () => ({ state: "connected", provider: "MiniMax", model: "MiniMax-M2.5", prompt_version: "plant-planning-v1" }) };
      if (path.includes("financial-constraints")) return { ok: true, json: async () => ({ case_code: "INC-GX-TEST", legal_entity_code: "LE-GX-TEST", financial_constraint: requirement.financial_constraint }) };
      if (path.includes("/requirements/")) return { ok: true, json: async () => requirement };
      if (path.includes("/proposals?")) return { ok: true, json: async () => currentProposalSet };
      if (path.includes("/investigations")) return { ok: true, json: async () => ({ items: [{
        request: { schema_version: "1.0", investigation_request_id: "INV-1", case_code: "INC-GX-TEST", proposal_set_id: "SET-1", proposal_id: "SITE-1", expected_revision: 1, world_run_id: "RUN-1", scope: [], requested_by: "owner", requested_at: "2026-08-01T00:00:00Z", status: "waiting_world" },
        status: "observed", work_item_status: "completed", observation: { schema_version: "1.0", observation_id: "OBS-1", investigation_request_id: "INV-1", proposal_id: "SITE-1", result: "completed", ownership_status: "verified", available_area_m2: 9000, electricity_kva: 1500, quoted_amount: { value: "8000000.00", currency: "CNY", scale: 2 }, available_at: "2026-10-01T00:00:00Z", permit_status: "eligible", evidence_refs: ["world-document:quote-1"], notes: "", external_actor_id: "park", observed_at: "2026-08-01T00:00:00Z" },
      }] }) };
      return { ok: true, json: async () => trace };
    }));
    render(<PlantBuildPlay onExit={() => undefined} />);
    await openNpcTask(/周衡：比较事实并提交推荐/);
    await waitFor(() => expect(screen.getByRole("heading", { name: "外部事实比较" })).toBeInTheDocument());
    expect(screen.getByText("硬约束不通过")).toBeInTheDocument();
    expect(screen.getByText("可用电力低于最低要求")).toBeInTheDocument();
    expect(screen.getByText("最低要求 2,000 kVA")).toBeInTheDocument();
    expect(screen.getByText("实测 1,500 kVA")).toBeInTheDocument();
    expect(screen.getByText("短缺 500 kVA")).toBeInTheDocument();
    expect(screen.getByText(/当前权重来源：界面默认比较偏好/)).toBeInTheDocument();
    expect(screen.getByText("world-document:quote-1")).toBeInTheDocument();
    expect(screen.getByLabelText("成本权重")).toHaveValue(35);
    expect(screen.getByLabelText("成本权重")).toBeDisabled();
    expect(screen.getByText(/不是正式推荐、选址批准或投资批准/)).toBeInTheDocument();
    fireEvent.click(screen.getByRole("button", { name: "修订设施需求" }));
    expect(screen.getByText("修订设施需求 · 将保存为第 2 版")).toBeInTheDocument();
    expect(screen.getByLabelText(/目标区域/)).toHaveValue("苏州");
    expect(screen.getByLabelText(/最小电力容量/)).toHaveValue(2000);
    expect(screen.getByLabelText(/本次投资申请金额/)).toHaveValue(10000000);
    expect(screen.getByLabelText(/本次修订原因/)).toHaveValue("");
    expect(screen.getByRole("button", { name: "保存修订并让 Agent 重新生成候选" })).toBeEnabled();
    fireEvent.click(screen.getByRole("button", { name: "返回外部事实比较" }));
    expect(screen.getByRole("heading", { name: "外部事实比较" })).toBeInTheDocument();
  });

  it("offers the site control mission immediately after formal site selection", async () => {
    const recommendation = {
      schema_version: "1.0", recommendation_id: "REC-CONTROL-1", case_code: "INC-GX-TEST",
      proposal_set_id: "SET-1", proposal_set_revision: 1, selected_proposal_id: "SITE-1",
      assessment_policy_version: "site-assessment-v1", weights: { cost: 35, schedule: 25, capacity: 20, control: 20 },
      recommendation_reason: "完成外部调研后推荐该正式场址", alternative_comparison: "相较其他候选更符合经营约束",
      recommended_at: "2026-08-01T12:00:00Z", requirement_id: "REQ-1", input_hash: "sha256:recommendation",
      approval_flow_key: "genesis.site.selection.approval", approval_request_id: "APR-1", eligible_count: 1,
      status: "formalized", recommended_by: "project-owner", assessments: [],
    };
    const decision = {
      selection_id: "SEL-1", recommendation_id: "REC-CONTROL-1", case_code: "INC-GX-TEST",
      selected_proposal_id: "SITE-1", approval_request_id: "APR-1", formalized_by: "project-owner",
      formalized_at: "2026-08-01T13:00:00Z",
    };
    vi.stubGlobal("fetch", vi.fn(async (input: RequestInfo | URL) => {
      const path = String(input);
      if (path.endsWith("planning-status")) return { ok: true, json: async () => ({ state: "connected", provider: "MiniMax", model: "MiniMax-M3", prompt_version: "plant-planning-v2" }) };
      if (path.includes("financial-constraints")) return { ok: true, json: async () => ({ case_code: "INC-GX-TEST", legal_entity_code: "LE-GX-TEST", financial_constraint: { available_cash: { value: "20000000.00", currency: "CNY", scale: 2 }, approved_budget: { value: "15000000.00", currency: "CNY", scale: 2 }, cash_source_ref: "gl:BOOK:1002", budget_source_ref: "budget:BUD-1", snapshot_hash: "sha256:authority" } }) };
      if (path.includes("/requirements/") || path.includes("/proposals?")) return { ok: false, status: 404, text: async () => "not found" };
      if (path.includes("/investigations") || path.includes("/site-controls")) return { ok: true, json: async () => ({ items: [] }) };
      if (path.includes("/site-selections")) return { ok: true, json: async () => ({ items: [{ recommendation, decision, record_status: "formalized", approval_status: "consumed" }] }) };
      if (path.includes("/plant-build/approvals/APR-1")) return { ok: true, json: async () => ({ item: { id: "APR-1", status: "consumed", requester_id: "project-owner" }, detail: { flow_key: "genesis.site.selection.approval", flow_version: 1, flow_name: "工厂场址正式选择审批", subject: { title: "INC-GX-TEST · 工厂场址正式选择", summary: "正式选址已经生效", operation: "site.selection.formalize" }, assignments: [], can_decide: false } }) };
      return { ok: true, json: async () => trace };
    }));
    render(<PlantBuildPlay onExit={() => undefined} />);
    const mission = await screen.findByRole("button", { name: /园区权利方：办理协议与场地交付/ });
    fireEvent.click(mission);
    await screen.findByRole("dialog", { name: "当前经营任务" });
    expect(screen.getByRole("heading", { name: "场地控制与实际交付" })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "发起场地控制交付工作项" })).toBeInTheDocument();
  });

  it("lets the player confirm delivery while the World engine creates all evidence", async () => {
    const recommendation = {
      schema_version: "1.0", recommendation_id: "REC-CONTROL-1", case_code: "INC-GX-TEST", proposal_set_id: "SET-1", proposal_set_revision: 1,
      selected_proposal_id: "SITE-1", assessment_policy_version: "site-assessment-v1", weights: { cost: 35, schedule: 25, capacity: 20, control: 20 },
      recommendation_reason: "完成外部调研后推荐该正式场址", alternative_comparison: "相较其他候选更符合经营约束", recommended_at: "2026-08-01T12:00:00Z",
      requirement_id: "REQ-1", input_hash: "sha256:recommendation", approval_flow_key: "genesis.site.selection.approval", approval_request_id: "APR-1", eligible_count: 1,
      status: "formalized", recommended_by: "project-owner", assessments: [],
    };
    const decision = { selection_id: "SEL-1", recommendation_id: "REC-CONTROL-1", case_code: "INC-GX-TEST", selected_proposal_id: "SITE-1", approval_request_id: "APR-1", formalized_by: "project-owner", formalized_at: "2026-08-01T13:00:00Z" };
    const control = {
      request: { schema_version: "1.0", control_request_id: "CTRL-1", selection_id: "SEL-1", case_code: "INC-GX-TEST", selected_proposal_id: "SITE-1", world_run_id: "WORLD-1", agreement_mode: "lease", requested_handover_at: "2026-10-01T00:00:00Z", required_evidence: ["executed_agreement", "handover_record", "possession_authority"], requested_by: "project-owner", requested_at: "2026-08-02T10:00:00Z", status: "waiting_world" },
      status: "waiting_world",
    };
    const fetchMock = vi.fn(async (input: RequestInfo | URL, init?: RequestInit) => {
      const path = String(input);
      if (path.endsWith("planning-status")) return { ok: true, json: async () => ({ state: "connected", provider: "MiniMax", model: "MiniMax-M3", prompt_version: "plant-planning-v2" }) };
      if (path.includes("financial-constraints")) return { ok: true, json: async () => ({ case_code: "INC-GX-TEST", legal_entity_code: "LE-GX-TEST", financial_constraint: { available_cash: { value: "20000000.00", currency: "CNY", scale: 2 }, approved_budget: { value: "15000000.00", currency: "CNY", scale: 2 }, cash_source_ref: "gl:BOOK:1002", budget_source_ref: "budget:BUD-1", snapshot_hash: "sha256:authority" } }) };
      if (path.includes("/requirements/") || path.includes("/proposals?")) return { ok: false, status: 404, text: async () => "not found" };
      if (path.includes("/investigations")) return { ok: true, json: async () => ({ items: [] }) };
      if (path.includes("/site-selections")) return { ok: true, json: async () => ({ items: [{ recommendation, decision, record_status: "formalized", approval_status: "consumed" }] }) };
      if (path.includes("/plant-build/approvals/APR-1")) return { ok: true, json: async () => ({ item: { id: "APR-1", status: "consumed", requester_id: "project-owner" }, detail: { flow_key: "genesis.site.selection.approval", flow_version: 1, flow_name: "工厂场址正式选择审批", subject: { title: "INC-GX-TEST · 工厂场址正式选择", summary: "正式选址已经生效", operation: "site.selection.formalize" }, assignments: [], can_decide: false } }) };
      if (path.endsWith("/site-controls/observations") && init?.method === "POST") return { ok: true, status: 201, json: async () => ({ status: "committed", idempotent_replay: false, world_message_id: "world-OBS-1", observation: {} }) };
      if (path.includes("/site-controls")) return { ok: true, json: async () => ({ items: [control] }) };
      return { ok: true, json: async () => trace };
    });
    vi.stubGlobal("fetch", fetchMock);
    render(<PlantBuildPlay onExit={() => undefined} />);
    fireEvent.click(await screen.findByRole("button", { name: /园区权利方：办理协议与场地交付/ }));
    expect(await screen.findByText("场地已经准备交付")).toBeInTheDocument();
    expect(screen.getByText(/只需核对本次交付并确认接收/)).toBeInTheDocument();
    expect(screen.queryByLabelText("已签协议引用")).not.toBeInTheDocument();
    expect(screen.queryByLabelText("交付记录引用")).not.toBeInTheDocument();
    expect(screen.queryByLabelText("控制权生效时间")).not.toBeInTheDocument();
    fireEvent.click(screen.getByRole("button", { name: "确认接收场地" }));
    await waitFor(() => expect(fetchMock.mock.calls.some(([path, init]) => String(path).endsWith("/site-controls/observations") && init?.method === "POST")).toBe(true));
    const request = JSON.parse(fetchMock.mock.calls.find(([path, init]) => String(path).endsWith("/site-controls/observations") && init?.method === "POST")?.[1]?.body as string);
    expect(request).toEqual({ schema_version: "1.0", case_code: "INC-GX-TEST", control_request_id: "CTRL-1", action: "accept_delivery" });
    expect(request).not.toHaveProperty("observation");
  });

  it("lets the routed governance authority decide a site approval inside AESE", async () => {
    let approved = false;
    const requirement = {
      schema_version: "1.0", requirement_id: "facility-requirement-INC-GX-TEST", tenant_id: "tenant-gx-test", case_code: "INC-GX-TEST", legal_entity_code: "LE-GX-TEST",
      target_region: "苏州", facility_purpose: "制造基地", minimum_area_m2: 8000, minimum_electricity_kva: 2000, target_available_at: "2026-12-01T00:00:00Z", candidate_count: 2, allowed_option_types: ["lease_and_retrofit"], investment_request: { value: "10000000.00", currency: "CNY", scale: 2 }, minimum_cash_reserve: { value: "2000000.00", currency: "CNY", scale: 2 }, financial_constraint: { available_cash: { value: "20000000.00", currency: "CNY", scale: 2 }, approved_budget: { value: "15000000.00", currency: "CNY", scale: 2 }, cash_source_ref: "gl:BOOK:1002", budget_source_ref: "budget:BUD-1", snapshot_hash: "sha256:authority" }, preferences: [], revision: 1, revision_reason: "test",
    };
    const proposalSet = { schema_version: "1.0", proposal_set_id: "SET-1", requirement_id: requirement.requirement_id, revision: 1, status: "candidate_only", proposals: [{ proposal_id: "SITE-1", option_type: "lease_and_retrofit", display_name: "候选场址一", business_rationale: "测试", estimated_amount: { minimum: { value: "7000000.00", currency: "CNY", scale: 2 }, likely: { value: "8000000.00", currency: "CNY", scale: 2 }, maximum: { value: "9000000.00", currency: "CNY", scale: 2 }, basis: "测试" }, estimated_schedule: { earliest: "2026-09-01T00:00:00Z", likely: "2026-10-01T00:00:00Z", latest: "2026-11-01T00:00:00Z" }, assumptions: [], facts_required: [], risks: [], source_refs: [], confidence: "0.8", status: "proposed" }], evidence: { provider: "MiniMax", model: "MiniMax-M3", prompt_version: "plant-planning-v2", input_hash: "sha256:in", output_hash: "sha256:out", validated_at: "2026-08-01T00:00:00Z" } };
    const recommendation = { schema_version: "1.0", recommendation_id: "REC-1", case_code: "INC-GX-TEST", proposal_set_id: "SET-1", proposal_set_revision: 1, selected_proposal_id: "SITE-1", assessment_policy_version: "site-assessment-v1", weights: { cost: 35, schedule: 25, capacity: 20, control: 20 }, recommendation_reason: "该场址满足当前经营与投产要求", alternative_comparison: "相较其他方案其交付时间和资金占用更优", recommended_at: "2026-08-01T12:00:00Z", requirement_id: requirement.requirement_id, input_hash: "sha256:recommendation", approval_flow_key: "genesis.site.selection.approval", approval_request_id: "APR-1", eligible_count: 1, status: "waiting_approval", recommended_by: "project-owner", assessments: [{ proposal_id: "SITE-1", display_name: "候选场址一", observation_id: "OBS-1", eligible: true, hard_failures: [], total_score: 86.2, evidence_refs: ["world-document:quote-1"] }] };
    const fetchMock = vi.fn(async (input: RequestInfo | URL, init?: RequestInit) => {
      const path = String(input);
      if (path.endsWith("planning-status")) return { ok: true, json: async () => ({ state: "connected", provider: "MiniMax", model: "MiniMax-M3", prompt_version: "plant-planning-v2" }) };
      if (path.includes("financial-constraints")) return { ok: true, json: async () => ({ case_code: "INC-GX-TEST", legal_entity_code: "LE-GX-TEST", financial_constraint: requirement.financial_constraint }) };
      if (path.includes("/requirements/")) return { ok: true, json: async () => requirement };
      if (path.includes("/proposals?")) return { ok: true, json: async () => proposalSet };
      if (path.endsWith("/investigations")) return { ok: true, json: async () => ({ items: [] }) };
      if (path.includes("/site-selections")) return { ok: true, json: async () => ({ items: [{ recommendation, record_status: "waiting_approval", approval_status: approved ? "approved" : "pending" }] }) };
      if (path.includes("/plant-build/approvals/APR-1")) return { ok: true, json: async () => ({ item: { id: "APR-1", status: approved ? "approved" : "pending", requester_id: "project-owner" }, detail: { flow_key: "genesis.site.selection.approval", flow_version: 1, flow_name: "工厂场址正式选择审批", subject: { title: "INC-GX-TEST · 工厂场址正式选择", summary: "审阅调研事实和推荐理由", operation: "site.selection.formalize" }, assignments: [{ id: "ASG-1", stage_code: "governance_authority", stage_name: "企业治理权责审批", mode: "all", display_name: "创始董事长", selector_type: "position", selector_value: "chair", status: approved ? "approved" : "active" }], can_decide: !approved } }) };
      if (path.includes("/commands/iaos/approvals/APR-1/approve") && init?.method === "POST") { approved = true; return { ok: true, json: async () => ({ item: { id: "APR-1", status: "approved" } }) }; }
      return { ok: true, json: async () => trace };
    });
    vi.stubGlobal("fetch", fetchMock);
    render(<PlantBuildPlay onExit={() => undefined} />);
    await openNpcTask(/林岚：查看审批与正式选址/);
    await waitFor(() => expect(screen.getByText("INC-GX-TEST · 工厂场址正式选择")).toBeInTheDocument());
    expect(screen.getByText(/创始董事长 · position:chair/)).toBeInTheDocument();
    expect(screen.getByText("当前身份可决定")).toBeInTheDocument();
    fireEvent.change(screen.getByLabelText(/审批意见/), { target: { value: "同意该场址进入正式选址与后续合同阶段" } });
    fireEvent.click(screen.getByRole("button", { name: "批准选址推荐" }));
    await waitFor(() => expect(screen.getByRole("button", { name: "批准已生效 · 正式选址" })).toBeInTheDocument());
    expect(fetchMock.mock.calls.some(([path]) => String(path).includes("/commands/iaos/approvals/APR-1/approve"))).toBe(true);
  });

  it("starts contractor sourcing from the approved WBS without price or evidence inputs", async () => {
    const project = { project_id: "PROJECT-1", project_name: "苏州制造基地项目", status: "active", wbs_items: [{ wbs_code: "WBS-02", name: "现场施工", phase: "construction", sequence: 2, owner_position: "plant-project-lead", planned_start_at: "2026-10-01T00:00:00Z", planned_finish_at: "2027-02-01T00:00:00Z", budget_share_bps: 3500, acceptance_criteria: "施工验收" }] };
    const fetchMock = vi.fn(async (input: RequestInfo | URL, init?: RequestInit) => {
      const path = String(input);
      if (path.endsWith("planning-status")) return { ok: true, json: async () => ({ state: "connected", provider: "MiniMax", model: "MiniMax-M3", prompt_version: "plant-planning-v2" }) };
      if (path.includes("financial-constraints")) return { ok: true, json: async () => ({ case_code: "INC-GX-TEST", legal_entity_code: "LE-GX-TEST", financial_constraint: { available_cash: { value: "20000000.00", currency: "CNY", scale: 2 }, approved_budget: { value: "15000000.00", currency: "CNY", scale: 2 }, cash_source_ref: "gl:BOOK:1002", budget_source_ref: "budget:BUD-1", snapshot_hash: "sha256:authority" } }) };
      if (path.includes("/requirements/") || path.includes("/proposals?")) return { ok: false, status: 404, text: async () => "not found" };
      if (path.includes("/facility-projects")) return { ok: true, json: async () => ({ items: [{ plan: { plan_id: "PLAN-1", case_code: "INC-GX-TEST", project_name: project.project_name, delivery_strategy: "design_build", budget_ceiling: { value: "12000000.00", currency: "CNY", scale: 2 }, target_start_at: "2026-09-01T00:00:00Z", target_ready_at: "2027-03-01T00:00:00Z", wbs_items: project.wbs_items, status: "approved" }, plan_hash: "sha256:plan", status: "active", approval_request_id: "APR-PROJECT", approval_status: "consumed", project }] }) };
      if (path.includes("/contract-awards")) return { ok: true, json: async () => ({ items: [] }) };
      if (path.endsWith("/contract-rfqs") && init?.method === "POST") return { ok: true, status: 201, json: async () => ({ status: "waiting_world" }) };
      return { ok: true, json: async () => path.includes("/plant-build/") ? { items: [] } : trace };
    });
    vi.stubGlobal("fetch", fetchMock);
    render(<PlantBuildPlay onExit={() => undefined} />);
    await openNpcTask(/顾远：选择采购包并发布 RFQ/);
    expect(screen.queryByLabelText(/合同金额|承包商|证据/)).not.toBeInTheDocument();
    fireEvent.click(screen.getByRole("button", { name: /WBS-02 · 现场施工/ }));
    fireEvent.click(screen.getByRole("button", { name: "确认采购包并发布 RFQ" }));
    await waitFor(() => expect(fetchMock.mock.calls.some(([path, init]) => String(path).endsWith("/contract-rfqs") && init?.method === "POST")).toBe(true));
    const request = JSON.parse(fetchMock.mock.calls.find(([path, init]) => String(path).endsWith("/contract-rfqs") && init?.method === "POST")?.[1]?.body as string);
    expect(request).toEqual({ case_code: "INC-GX-TEST", package_code: "WBS-02", sourcing_strategy: "specialist_packages" });
  });
});
