import { expect, test } from "@playwright/test";
test("first commercial cycle closes", async ({ page }) => {
  await page.goto("/#world-industrialization");
  await page.getByRole("button", { name: "首次商业交付 Campaign" }).click();
  await expect(page).toHaveURL(/#world-first-delivery$/);
  await expect(
    page.getByRole("heading", { name: "消费 M12 量产资格并完成财务结转" }),
  ).toBeVisible();
  for (let i = 0; i < 12; i++)
    await page.getByRole("button", { name: "单步" }).click();
  await expect(page.getByText("Genesis cycle closed")).toBeVisible();
  await expect(page.getByText("SHIP-GENESIS-0001-C")).toBeVisible();
  await expect(page.getByText("¥4,200,000")).toBeVisible();
});
