import { expect, test } from "@playwright/test";
test("capability build reaches M12 eligibility", async ({ page }) => {
  await page.goto("/#world-plant-build");
  await page.getByRole("button", { name: "生产能力 Campaign" }).click();
  await expect(page).toHaveURL(/#world-capability-build$/);
  await expect(
    page.getByRole("heading", { name: "消费 M10 设施资格" }),
  ).toBeVisible();
  for (let i = 0; i < 9; i++)
    await page.getByRole("button", { name: "单步" }).click();
  await expect(page.getByText("M12 eligible")).toBeVisible();
  await expect(page.getByText("EQ-LEAK-TEST-01")).toBeVisible();
  await expect(page.getByText("10 / 10")).toBeVisible();
});
