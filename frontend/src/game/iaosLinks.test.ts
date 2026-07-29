import { describe, expect, it } from "vitest";
import { financeWorkspaceUrl } from "./iaosLinks";

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
