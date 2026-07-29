import { expect, test } from "@playwright/test";

test("root URL opens the Enterprise Genesis product home", async ({ page }) => {
  await page.goto("/");

  await expect(page.getByRole("heading",{name:"回到你的企业世界"})).toBeVisible();
  await page.getByLabel("游戏用户名").fill("founder-principal");
  await page.getByRole("button",{name:/进入企业世界/}).click();
  await expect(
    page.getByRole("heading", {
      name: "从一个想法，创建一家真正运行的企业",
    }),
  ).toBeVisible();
  await expect(
    page.getByRole("heading", { name: /客户追加订单下的交付承诺重算/ }),
  ).toHaveCount(0);
  await expect(page.getByRole("button", { name: /创建新企业/ }).first()).toBeVisible();
  await expect(
    page.getByRole("button", { name: /体验华辰样板世界/ }),
  ).toBeVisible();
});

test("sample-world action keeps the former sandbox available", async ({
  page,
}) => {
  await page.goto("/");
  await page.getByLabel("游戏用户名").fill("founder-principal");
  await page.getByRole("button",{name:/进入企业世界/}).click();
  await page.getByRole("button", { name: /体验华辰样板世界/ }).click();

  await expect(page).toHaveURL(/#sandbox$/);
  await expect(
    page.getByRole("heading", { name: /客户追加订单下的交付承诺重算/ }),
  ).toBeVisible();
});

test("new enterprise starts with server-assigned tenant provisioning", async ({
  page,
}) => {
  await page.goto("/");
  await page.getByLabel("游戏用户名").fill("founder-principal");
  await page.getByRole("button",{name:/进入企业世界/}).click();
  await page.getByRole("button", { name: /创建新企业/ }).first().click();

  await expect(
    page.getByRole("heading", { name: "先为你的企业建立独立运行空间" }),
  ).toBeVisible();
  await expect(page.getByText("你不需要填写或选择 tenant ID。")).toBeVisible();
  await expect(
    page.getByRole("button", { name: /创建空间并进入 AI 身份工作室/ }),
  ).toBeVisible();
});

test("signed-in founder can select an existing enterprise and continue", async ({page})=>{
  const workspace={
    workspace_id:"gxw-existing",owner_player_id:"player-local-existing",display_name:"我的热管理企业",
    tenant_id:"tenant-gx-existing",world_run_id:"world-gx-existing",case_code:"INC-GX-EXISTING",
    status:"active",current_step:"identity_studio_ready",created_at:"2026-07-28T03:00:00Z",updated_at:"2026-07-28T03:00:00Z",
  };
  await page.route("**/api/aese/v1/genesis/workspaces",route=>route.fulfill({status:200,json:{items:[workspace]}}));
  await page.route("**/api/aese/v1/genesis/workspaces/gxw-existing/session",route=>route.fulfill({status:200,json:{...workspace,tenant_token:"founder-token"}}));
  await page.goto("/");
  await page.getByLabel("游戏用户名").fill("founder-principal");
  await page.getByRole("button",{name:/进入企业世界/}).click();
  await expect(page.getByRole("heading",{name:"我的企业"})).toBeVisible();
  await expect(page.getByRole("heading",{name:"我的热管理企业"})).toBeVisible();
  await page.getByRole("button",{name:"继续游戏"}).click();
  await expect(page).toHaveURL(/#enterprise-genesis\?tenant=tenant-gx-existing&case=INC-GX-EXISTING&workspace=gxw-existing/);
});
