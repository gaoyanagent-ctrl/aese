import { expect, test } from "@playwright/test";
test("industrialization reaches M13 eligibility", async ({ page }) => {
  await page.goto("/#world-capability-build");
  await page.getByRole("button", { name: "产品工业化 Campaign" }).click();
  await expect(page).toHaveURL(/#world-industrialization$/);
  await expect(
    page.getByRole("heading", { name: "消费 M11 工业化资格" }),
  ).toBeVisible();
  for (let i = 0; i < 10; i++)
    await page.getByRole("button", { name: "单步" }).click();
  await expect(page.getByText("M13 eligible")).toBeVisible();
  await expect(page.getByText("BOM-HCTM-BCP-A01-V1")).toBeVisible();
  await expect(page.getByText("hctm-stable-codes-compatible")).toBeVisible();
});
