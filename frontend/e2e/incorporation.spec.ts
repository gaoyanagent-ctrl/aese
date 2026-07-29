import { expect, test } from "@playwright/test";

test("Genesis incorporation opens from the lifecycle hub", async ({ page }) => {
  const errors: string[] = [];
  page.on("pageerror", (error) => errors.push(error.message));
  await page.goto("/#world");
  await page.getByRole("link", { name: "进入 M9 公司成立" }).click();
  await expect(page).toHaveURL(/#world-incorporation$/);
  await expect(page.getByRole("heading", { name: "投资人形成设立意图" })).toBeVisible();
  await expect(page.getByText("预算不是现金或实际支出")).toBeVisible();
  await expect(page.getByText("已解锁 1/8")).toBeVisible();
  expect(errors).toEqual([]);
});
