import { describe, expect, it } from "vitest";
import { financeWorkspaceUrl, m9KnowledgeUrl } from "./iaosLinks";

describe("IAOS finance drill-through", () => {
  it("opens the responsibility and work-item view for the current tenant and case", () => {
    const url = financeWorkspaceUrl(
      { tenant_id: "tenant-gx-demo", case_code: "INC-DEMO-001" },
      "operations",
    );
    expect(url).toContain("tenant=tenant-gx-demo");
    expect(url).toContain("case=INC-DEMO-001");
    expect(url).toContain("view=operations");
    expect(url.endsWith("#finance_workspace")).toBe(true);
  });
});

describe("M9 knowledge drill-through", () => {
  it("preserves tenant, case, article and current work-item context", () => {
    localStorage.setItem("aese_genesis_workspace_id", "gxw-demo");
    const url = m9KnowledgeUrl(
      { tenant_id: "tenant-gx-demo", case_code: "INC-DEMO-001", world_run_id: "world-run-1" },
      {
        work_item_id: "work-item-04",
        capability: "capital.commitment.record",
        kind: "agent_task",
        owner_id: "finance-agent",
        owner_type: "agent",
      },
    );
    expect(url).toContain("article_id=KB-HCTM-M9-INCORPORATION");
    expect(url).toContain("workspace_id=gxw-demo");
    expect(url).toContain("case_code=INC-DEMO-001");
    expect(url).toContain("world_run_id=world-run-1");
    expect(url).toContain("node_id=work-item-04");
    expect(url).toContain("actor_id=finance-agent");
    expect(url).toContain("actor_type=agent");
    expect(url).toContain("task_type=agent_task");
    expect(url).toContain("capability=capital.commitment.record");
    expect(url).toContain("tenant=tenant-gx-demo");
    expect(url.endsWith("#knowledge_center")).toBe(true);
  });
});
