import { expect, test } from "@playwright/test";

test.skip(
  !process.env.GENESIS_PLAYER_PASSWORD,
  "requires a real IAOS Genesis Player credential",
);

test("real Genesis Player login restores owner-scoped enterprises", async ({page})=>{
  const browserErrors:string[]=[];
  page.on("console",message=>{
    if(message.type()==="error")browserErrors.push(message.text());
  });
  await page.goto("/");
  await page.getByLabel("用户名").fill(process.env.GENESIS_PLAYER_USERNAME??"founder-principal");
  await page.getByLabel("密码",{exact:true}).fill(process.env.GENESIS_PLAYER_PASSWORD!);
  await page.getByRole("button",{name:"安全登录"}).click();

  await expect(page.getByRole("heading",{name:"我的企业"})).toBeVisible();
  await expect(page.getByRole("heading",{name:"高岩的牛牛公司"})).toBeVisible();
  await expect(page.getByText("founder-principal")).toBeVisible();
  expect(browserErrors).toEqual([]);
});
