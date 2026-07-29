import { readFileSync } from "node:fs";
import { expect, test } from "@playwright/test";

const projection=JSON.parse(readFileSync(new URL("../../world-contracts/fixtures/game-projection.json",import.meta.url),"utf8"));

test("founder plays the guided office conversation before IAOS incorporation",async({page})=>{
 let created=false;
 await page.addInitScript(()=>{localStorage.setItem("iaos_token","founder-token");localStorage.setItem("aese_iaos_tenant_id","tenant-game")});
 await page.route("**/api/aese/v1/game/incorporation/**/projection?frame=*",route=>created?route.fulfill({json:projection}):route.fulfill({status:404,json:{code:"incorporation_case_not_found"}}));
 await page.route("**/api/aese/v1/game/creative/intent",route=>route.fulfill({json:{...route.request().postDataJSON(),schema_version:"1.0",intent_id:"intent-rpg",assumptions:[],needs_confirmation:[],created_at:"2026-07-28T04:00:00Z"}}));
 await page.route("**/api/aese/v1/game/creative/names",route=>route.fulfill({json:{status:"candidate_only",proposals:[{proposal_id:"p1",chinese_name:"拓流热管理有限公司",english_name:"FlowForge Thermal",short_name:"拓流",rationale:"工程可靠",slogan:"让制造稳定前行",keywords:["可靠"],primary_color:"#2E8174",risk_hints:[],status:"candidate"}]}}));
 await page.route("**:8082/api/v1/incorporations/cases",async route=>{created=true;await route.fulfill({status:201,json:{to:"incorporation_case_opened"}})});
 await page.goto("/#enterprise-genesis?tenant=tenant-game&case=INC-RPG-001");
 await expect(page.getByText("创始人办公室",{exact:true}).first()).toBeVisible();
 await expect(page.getByRole("img",{name:/PixiJS 渲染的创始人办公室/})).toBeVisible();
 await page.getByRole("button",{name:/林澜/}).click();
 await page.getByRole("button",{name:/使用这个形象/}).click();
 for(const choice of["新能源汽车热管理","新能源汽车制造商","电池冷却板与热管理系统","可靠 · 工程 · 长期主义"]){
  await page.getByRole("button",{name:new RegExp(choice)}).click();
 }
 await page.getByRole("button",{name:/让 AI 形成 4 组企业身份提案/}).click();
 await page.getByRole("button",{name:/拓流热管理有限公司/}).click();
 await page.getByRole("button",{name:/我选择这个名字/}).click();
 await page.getByRole("button",{name:/签署创始人指令并启动企业设立/}).click();
 await expect(page.getByText("founder.resolution.prepare").first()).toBeVisible();
});
