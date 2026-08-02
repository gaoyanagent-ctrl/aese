import { expect, test } from "@playwright/test";

test("plant build opens as an enterprise scene with on-demand missions", async ({ page }) => {
  await page.goto("/#world-incorporation");
  await page.getByRole("button", { name: "工厂建设 Campaign" }).click();
  await expect(page).toHaveURL(/#world-plant-build$/);
  await expect(page.getByRole("heading", { name: "新制造企业 · 工厂选址与设施建设" })).toBeVisible();
  await expect(page.getByRole("region", { name: "M10 企业与外部世界" })).toBeVisible();
  await expect(page.getByRole("dialog", { name: "当前经营任务" })).toBeHidden();
  await page.getByRole("button", { name: /纪元：与规划 Agent 制定需求/ }).click();
  await expect(page.getByRole("dialog", { name: "当前经营任务" })).toBeVisible();
  await expect(page.getByRole("heading", { name: "设施需求与候选方案" })).toBeVisible();
});
