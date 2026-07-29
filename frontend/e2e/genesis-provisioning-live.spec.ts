import { expect, test } from "@playwright/test";

test.skip(
  process.env.GENESIS_LIVE_PROVISIONING !== "1",
  "creates a real IAOS tenant; enable explicitly",
);

test("provisions an isolated tenant before opening the AI identity studio", async ({
  page,
}) => {
  await page.goto("/");
  await page.getByLabel("游戏用户名").fill(`founder-live-${Date.now()}`);
  await page.getByRole("button",{name:/进入企业世界/}).click();
  await page.getByRole("button", { name: /创建新企业/ }).first().click();
  await page
    .getByLabel("创业项目名称")
    .fill(`浏览器验收企业-${Date.now()}`);
  await page
    .getByRole("button", { name: /创建空间并进入 AI 身份工作室/ })
    .click();

  await expect(
    page.getByText("创始人办公室", { exact: true }).first(),
  ).toBeVisible({ timeout: 120_000 });
  const tenant = await page.evaluate(() =>
    localStorage.getItem("aese_iaos_tenant_id"),
  );
  expect(tenant).toMatch(/^tenant-gx-[0-9a-f]{16}$/);
  await expect(page).toHaveURL(/#enterprise-genesis\?tenant=tenant-gx-/);
});
