import{readFileSync}from"node:fs";
import{expect,test}from"@playwright/test";

const projection=JSON.parse(readFileSync(new URL("../../world-contracts/fixtures/game-projection.json",import.meta.url),"utf8"));

test("Enterprise Genesis supports three-view game operation and AI identity candidates",async({page},testInfo)=>{
 await page.addInitScript(()=>{localStorage.setItem("iaos_token","e2e-founder-token");localStorage.setItem("aese_iaos_tenant_id","tenant-demo")});
 await page.route("**/api/aese/v1/game/incorporation/**/projection?frame=*",route=>route.fulfill({json:projection}));
 await page.route("**/api/aese/v1/game/creative/intent",route=>route.fulfill({json:{
  schema_version:"1.0",intent_id:"intent-e2e",tenant_id:"tenant-demo",case_code:"INC-GAME-DEMO-001",
  raw_idea:"创建一家服务新能源汽车制造商的工业热管理公司，产品可靠、工程化并坚持长期主义。",
  industry:"热管理",customers:["新能源汽车制造商"],offerings:["电池冷却板"],brand_traits:["可靠","工程"],
  capital_minor:"100000000",risk_appetite:"balanced",assumptions:[],needs_confirmation:[],created_at:"2026-07-27T10:00:00Z"
 }}));
 await page.route("**/api/aese/v1/game/creative/names",route=>route.fulfill({json:{status:"candidate_only",proposals:[{
  proposal_id:"name-01",chinese_name:"澄流热管理有限公司",english_name:"ClearFlow Thermal Systems",
  short_name:"澄流科技",rationale:"可靠工业热管理",slogan:"让热管理稳定流动",keywords:["可靠","工程"],
  primary_color:"#167C80",risk_hints:["现实核名未完成"],status:"candidate"
 }]}}));
 let committedBody:Record<string,unknown>|undefined;
 await page.route("**:8082/api/v1/incorporations/cases",async route=>{committedBody=route.request().postDataJSON();await route.fulfill({status:201,json:{to:"incorporation_case_opened"}})});
 await page.goto("/#enterprise-genesis");
 await expect(page.getByRole("heading",{name:"企业身份工作室"})).toBeVisible();
 await expect(page.getByRole("img",{name:"企业创生等距世界地图"})).toBeVisible();
 for(const tab of["任务","员工","证据"]){
  await page.getByRole("tab",{name:tab}).click();
  await expect(page.getByRole("tab",{name:tab})).toHaveAttribute("aria-selected","true");
 }
 await page.getByRole("button",{name:"生成公司身份候选"}).click();
 await expect(page.getByText("澄流热管理有限公司")).toBeVisible();
 await page.getByRole("button",{name:/澄流热管理有限公司/}).click();
 await expect(page.getByText(/尚未成为 IAOS 正式企业事实/)).toBeVisible();
 await page.getByRole("button",{name:"确认身份并创建企业"}).click();
 await expect(page.getByRole("button",{name:/已通过 incorporation.case.open 创建/})).toBeVisible();
 expect(committedBody?.proposed_company_name).toBe("澄流热管理有限公司");
 await expect(page.getByRole("button",{name:"暂停"})).toHaveAttribute("aria-pressed","true");
 await page.getByRole("button",{name:"2×"}).click();
 await expect(page.getByRole("button",{name:"2×"})).toHaveAttribute("aria-pressed","true");
 await page.screenshot({path:testInfo.outputPath("enterprise-genesis.png"),fullPage:true});
});
