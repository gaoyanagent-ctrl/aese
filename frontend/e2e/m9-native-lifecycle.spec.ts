import { expect, test } from "@playwright/test";

const caseCode = process.env.M9_CASE_CODE ?? "INC-HCTM-001";

test("AESE consumes persisted IAOS M9 lifecycle projection after refresh", async ({ page, request }) => {
  const login = await request.post("http://127.0.0.1:8082/api/v1/auth/login", {
    data: {
      username: "founder-principal",
      password: "Founder-Lifecycle-2026!",
      tenant_id: "tenant-hctm-genesis",
    },
  });
  expect(login.ok()).toBeTruthy();
  const session = await login.json();
  await page.addInitScript(({ token }) => {
    localStorage.setItem("iaos_token", token);
    localStorage.setItem("aese_iaos_tenant_id", "tenant-hctm-genesis");
  }, session);
  const target = `/#world-incorporation?tenant=tenant-hctm-genesis&case=${encodeURIComponent(caseCode)}&process_run=&world_run=&correlation=`;
  await page.goto(target);
  await expect(page.getByTestId("iaos-lifecycle-projection")).toBeVisible();
  await expect(page.getByText("Intent / Observation / CommittedOutcome")).toBeVisible();
  const escapedHost = new URL(page.url()).hostname.replaceAll(".", "\\.");
  await expect(page.getByRole("link", { name: "打开 IAOS 设立案" })).toHaveAttribute("href", new RegExp(`^http://${escapedHost}:3000/.*tenant=.*case=.*process_run=.*world_run=.*correlation=`));
  await page.getByRole("button", { name: "复位" }).click();
  const persisted = await request.get(`http://127.0.0.1:8082/api/v1/incorporations/${encodeURIComponent(caseCode)}/trace`, {
    headers: { Authorization: `Bearer ${session.token}`, "X-Tenant-ID": "tenant-hctm-genesis" },
  });
  expect(persisted.ok()).toBeTruthy();
  expect((await persisted.json()).state.state).toBe("enterprise_operational_ready");
  await page.reload();
  await expect(page.getByTestId("iaos-lifecycle-projection")).toContainText(caseCode);
});
