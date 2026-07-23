import { expect, test } from "@playwright/test";

test("World home exposes the complete enterprise lifecycle", async ({ page }) => {
  await page.goto("/#world");
  await expect(page.getByRole("heading", { name: "企业生命周期运营中心" })).toBeVisible();
  await expect(page.getByText("从公司成立到产品交付与持续经营")).toBeVisible();
  for (const milestone of ["M8", "M9", "M10", "M11", "M12", "M13", "M14", "M15", "M16", "M17–M24"]) {
    await expect(page.getByText(milestone, { exact: true })).toBeVisible();
  }
  await page.getByRole("link", { name: /进入 M9 公司成立/ }).click();
  await expect(page).toHaveURL(/#world-incorporation$/);
  await expect(page.getByRole("heading", { name: "华辰苏州制造公司成立与治理" })).toBeVisible();
});

test("M8 tracer is a secondary architecture validation entry", async ({ page }) => {
  await page.goto("/#world");
  await page.getByRole("link", { name: /进入 M8 三态架构验证/ }).click();
  await expect(page).toHaveURL(/#world-tristate$/);
  await expect(page.getByRole("heading", { name: "LAS-WLD-02 三态偏差闭环" })).toBeVisible();
});
