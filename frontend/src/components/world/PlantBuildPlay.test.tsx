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

  it("fails closed when the external planning model is unavailable", async () => {
    vi.stubGlobal("fetch", vi.fn(async (input: RequestInfo | URL) => {
      const path = String(input);
      if (path.includes("/requirements/")) return { ok: false, status: 404, text: async () => "not found" };
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
    await waitFor(() => expect(screen.getByText("未启用外部模型")).toBeInTheDocument());
    expect(screen.getByRole("button", { name: /让 Agent 生成候选/ })).toBeDisabled();
    expect(screen.getByText(/不能生成虚拟固定候选/)).toBeInTheDocument();
    expect(screen.getByText(/来自 IAOS 权威账务/)).toBeInTheDocument();
  });

  it("submits editable requirements and renders agent evidence for human review", async () => {
    const fetchMock = vi.fn(async (input: RequestInfo | URL, init?: RequestInit) => {
      void init;
      const path = String(input);
      if (path.endsWith("planning-status")) return { ok: true, json: async () => ({ state: "connected", provider: "MiniMax", model: "MiniMax-M2.5", prompt_version: "plant-planning-v1" }) };
      if (path.includes("financial-constraints")) return { ok: true, json: async () => ({ case_code: "INC-GX-TEST", legal_entity_code: "LE-GX-TEST", financial_constraint: { available_cash: { value: "20000000.00", currency: "CNY", scale: 2 }, approved_budget: { value: "15000000.00", currency: "CNY", scale: 2 }, cash_source_ref: "gl:BOOK:1002", budget_source_ref: "budget:BUD-1", snapshot_hash: "sha256:authority" } }) };
      if (path.includes("/requirements/")) return { ok: true, status: 200, json: async () => ({ revision: 3 }) };
      if (path.endsWith("investigations")) return init?.method === "POST"
        ? { ok: true, json: async () => ({ status: "waiting_world", investigation_request: {}, work_item: {} }) }
        : { ok: true, json: async () => ({ items: [] }) };
      if (path.endsWith("reviews")) return { ok: true, json: async () => ({ status: "committed", proposal_review: {} }) };
      if (path.endsWith("proposals")) return { ok: true, json: async () => ({ proposal_set: {
        schema_version: "1.0", proposal_set_id: "set-1", requirement_id: "facility-requirement-INC-GX-TEST", revision: 1, status: "candidate_only",
        proposals: [{ proposal_id: "site-1", option_type: "lease_and_retrofit", display_name: "北区租赁改造建议", business_rationale: "平衡上市时间和初始现金占用", estimated_amount: { minimum: { value: "8000000.00", currency: "CNY", scale: 2 }, likely: { value: "10000000.00", currency: "CNY", scale: 2 }, maximum: { value: "12000000.00", currency: "CNY", scale: 2 }, basis: "需求参数与区域历史区间" }, estimated_schedule: { earliest: "2026-10-01T00:00:00Z", likely: "2026-11-01T00:00:00Z", latest: "2026-12-01T00:00:00Z" }, assumptions: ["存在合规标准厂房"], facts_required: ["核验供电容量"], risks: ["消防改造周期"], source_refs: ["requirement:facility-requirement-INC-GX-TEST"], confidence: "0.78", status: "proposed" }],
        evidence: { provider: "MiniMax", model: "MiniMax-M2.5", prompt_version: "plant-planning-v1", input_hash: "sha256:in", output_hash: "sha256:out", validated_at: "2026-08-01T08:00:00Z" },
      }, agent_job: {}, idempotent_replay: false }) };
      return { ok: true, json: async () => trace };
    });
    vi.stubGlobal("fetch", fetchMock);
    render(<PlantBuildPlay onExit={() => undefined} />);
    await waitFor(() => expect(screen.getByText(/MiniMax · MiniMax-M2.5/)).toBeInTheDocument());
    fireEvent.change(screen.getByLabelText(/目标区域/), { target: { value: "江苏省苏州市" } });
    fireEvent.change(screen.getByLabelText(/设施用途/), { target: { value: "冷却板制造基地" } });
    fireEvent.change(screen.getByLabelText(/最小面积/), { target: { value: "12000" } });
    fireEvent.change(screen.getByLabelText(/最小电力容量/), { target: { value: "2200" } });
    fireEvent.change(screen.getByLabelText(/目标可用时间/), { target: { value: "2026-11-01T08:00" } });
    fireEvent.change(screen.getByLabelText(/本次投资申请金额/), { target: { value: "15000000" } });
    fireEvent.change(screen.getByLabelText(/最低现金保留额/), { target: { value: "3000000" } });
    expect(screen.getByLabelText(/可用现金快照/)).toHaveValue("20000000.00");
    expect(screen.getByLabelText(/已批准预算快照/)).toHaveValue("15000000.00");
    fireEvent.click(screen.getByLabelText("租赁并改造"));
    fireEvent.submit(screen.getByRole("button", { name: /让 Agent 生成候选/ }).closest("form")!);
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
    expect(request.investment_request.value).toBe("15000000");
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
});
