import { expect, test } from "@playwright/test";

test("plant build campaign", async ({ page }) => {
  await page.goto("/#world-incorporation");
  await page.getByRole("button", { name: "工厂建设 Campaign" }).click();
  await expect(page).toHaveURL(/#world-plant-build$/);
  await expect(
    page.getByRole("heading", { name: "消费 M9 机器资格" }),
  ).toBeVisible();
  for (let i = 0; i < 9; i += 1) {
    await page.getByRole("button", { name: "单步" }).click();
  }
  await expect(page.getByText("M11 eligible")).toBeVisible();
  await expect(page.getByText("ZONE-HCTM-SZ-PRODUCTION")).toBeVisible();
});
