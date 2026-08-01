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

describe("M9 knowledge drill-through",()=>{
  it("preserves tenant, case, article and current capability",()=>{
    const url=m9KnowledgeUrl({tenant_id:"tenant-gx-demo",case_code:"INC-DEMO-001"},"capital.commitment.record");
    expect(url).toContain("article_id=KB-HCTM-M9-INCORPORATION");expect(url).toContain("capability=capital.commitment.record");expect(url).toContain("tenant=tenant-gx-demo");expect(url.endsWith("#knowledge_center")).toBe(true);
  });
});
