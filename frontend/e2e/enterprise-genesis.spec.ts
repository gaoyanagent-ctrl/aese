import{readFileSync}from"node:fs";
import{expect,test}from"@playwright/test";

const projection=JSON.parse(readFileSync(new URL("../../world-contracts/fixtures/game-projection.json",import.meta.url),"utf8"));

test("Enterprise Genesis supports three-view world-first city and location operation",async({page},testInfo)=>{
 await page.addInitScript(()=>{localStorage.setItem("iaos_token","e2e-founder-token");localStorage.setItem("aese_iaos_tenant_id","tenant-demo")});
 await page.route("**/api/aese/v1/game/incorporation/**/projection?frame=*",route=>route.fulfill({json:projection}));
 await page.goto("/#enterprise-genesis");
 await expect(page.getByRole("img",{name:"企业创生城市地图"})).toBeVisible();
 for(const tab of["当前事件","团队","治理档案"]){
  await page.getByRole("tab",{name:tab}).click();
  await expect(page.getByRole("tab",{name:tab})).toHaveAttribute("aria-selected","true");
 }
 await expect(page.getByRole("button",{name:/创始办公室 有新事件/})).toBeVisible();
 await expect(page.getByLabel("玩家在城市中的位置")).toBeVisible();
 await page.getByRole("button",{name:/创始办公室 有新事件/}).click();
 await expect(page.getByRole("status")).toContainText("前往创始办公室");
 await expect(page.getByRole("region",{name:"创始办公室室内场景"})).toBeVisible();
 await expect(page.getByRole("button",{name:"关闭游戏音效"})).toBeVisible();
 const roomBounds=await page.getByRole("region",{name:"创始办公室室内场景"}).boundingBox();
 const briefingBounds=await page.getByRole("region",{name:"创始办公室室内场景"}).getByLabel("当前剧情任务").boundingBox();
 expect(roomBounds).not.toBeNull();
 expect(briefingBounds).not.toBeNull();
 expect(briefingBounds!.x+briefingBounds!.width).toBeLessThanOrEqual(roomBounds!.x+roomBounds!.width);
 const overflowingIntroductions=await page.getByRole("region",{name:"创始办公室室内场景"}).locator(".gx-mission__story > div, .gx-scene-npc > div").evaluateAll(elements=>elements.filter(element=>element.scrollWidth>element.clientWidth+1).length);
 expect(overflowingIntroductions).toBe(0);
 await expect(page.getByRole("button",{name:"创始人经营桌"})).toBeVisible();
 await expect(page.getByLabel("纪元 NPC")).toBeVisible();
 await expect(page.getByLabel("创始人在室内的位置")).toBeVisible();
 expect(await page.getByLabel("纪元 NPC").locator("> i").evaluate(element=>getComputedStyle(element).backgroundImage)).toContain("ji-yuan-v1.png");
 const npcSize=await page.getByLabel("纪元 NPC").locator("> i").boundingBox();
 const playerSize=await page.getByLabel("创始人在室内的位置").locator("i").boundingBox();
 expect(npcSize).not.toBeNull();
 expect(playerSize).not.toBeNull();
 expect(npcSize!.height).toBeGreaterThan(170);
 expect(playerSize!.height).toBeGreaterThan(145);
 expect(npcSize!.height/playerSize!.height).toBeLessThan(1.3);
 await page.getByRole("button",{name:"创始人经营桌"}).click();
 await expect(page.getByLabel("创始人在室内的位置")).toHaveClass(/target-1/);
 await expect(page.getByLabel("创始人经营桌详情")).toContainText("IAOS Work Item");
 await page.getByRole("button",{name:"关闭物件详情"}).click();
 await page.screenshot({path:testInfo.outputPath("enterprise-genesis-room.png"),fullPage:true});
 await page.getByRole("button",{name:"返回城市地图"}).click();
 await expect(page.getByRole("img",{name:"企业创生城市地图"})).toBeVisible();
 await expect(page.getByRole("button",{name:"暂停"})).toHaveCount(0);
 await expect(page.getByRole("button",{name:"查看后续已提交状态"})).toHaveCount(0);
 await page.screenshot({path:testInfo.outputPath("enterprise-genesis.png"),fullPage:true});
});

test("organization mission opens its governed action from the story button",async({page})=>{
 const organizationProjection=structuredClone(projection);
 organizationProjection.chapter="talent_governance";
 organizationProjection.lifecycle={state:"capital_contribution_verified",current_step:"organization.establish",progress:59,blocked_by:"governance-agent"};
 organizationProjection.buildings=organizationProjection.buildings.map((building:{kind:string;state:string;available:boolean})=>({
  ...building,
  state:"active",
  available:true,
 }));
 organizationProjection.work_items=[{
  work_item_id:"WI-011",
  title:"建立企业初始组织",
  kind:"system_task",
  status:"ready",
  owner_type:"service",
  owner_id:"iaos-runtime",
  capability:"organization.establish",
  requires_me:false,
  evidence_ref:"iaos:work-item:INC-DEMO-001:11",
 }];
 await page.addInitScript(()=>{localStorage.setItem("iaos_token","e2e-founder-token");localStorage.setItem("aese_iaos_tenant_id","tenant-demo")});
 await page.route("**/api/aese/v1/game/incorporation/**/projection?frame=*",route=>route.fulfill({json:organizationProjection}));
 await page.goto("/#enterprise-genesis?tenant=tenant-demo&case=INC-DEMO-001");
 await page.getByRole("button",{name:"推进剧情"}).click();
 await page.getByRole("button",{name:"立即到达"}).click();
 await expect(page.getByRole("region",{name:"企业总部室内场景"})).toBeVisible();
 await page.getByRole("region",{name:"企业总部室内场景"}).getByRole("button",{name:"推进剧情"}).click();
 await expect(page.getByRole("heading",{name:"建立初始组织"})).toBeVisible();
});

test("headquarters finance center does not overlap governance table and drills into IAOS",async({page})=>{
 const financeProjection=structuredClone(projection);
 financeProjection.tenant_id="tenant-demo";
 financeProjection.case_code="INC-DEMO-001";
 financeProjection.chapter="operating_world";
 financeProjection.lifecycle={state:"enterprise_operational_ready",current_step:"",progress:100};
 financeProjection.buildings=financeProjection.buildings.map((building:{kind:string;state:string;available:boolean})=>({...building,state:"active",available:true}));
 financeProjection.work_items=[];
 financeProjection.finance_opening={
  ready:true,organization_code:"FIN-INC-DEMO-001",organization_status:"active",roles:["finance_lead"],
  book_code:"BOOK-INC-DEMO-001",accounting_standard:"CAS",functional_currency:"CNY",
  period_code:"2026-07",period_status:"open",journal_entry_no:"OPEN-INC-DEMO-001",journal_status:"posted",
  debit_minor:100000000,credit_minor:100000000,trial_balance:[],bank_journal:[],general_ledger:[],
  opening_balance_sheet:{as_of:"2026-07-29",currency:"CNY",assets:[],liabilities:[],equity:[],total_assets_minor:100000000,total_liabilities_minor:0,total_equity_minor:100000000,balanced:true},
  evidence_ref:"iaos:finance:INC-DEMO-001:OPEN-INC-DEMO-001",
 };
 await page.addInitScript(()=>{localStorage.setItem("iaos_token","e2e-founder-token");localStorage.setItem("aese_iaos_tenant_id","tenant-demo")});
 await page.route("**/api/aese/v1/game/incorporation/**/projection?frame=*",route=>route.fulfill({json:financeProjection}));
 await page.goto("/#enterprise-genesis?tenant=tenant-demo&case=INC-DEMO-001");
 await page.getByRole("button",{name:/企业总部/}).click();
 await page.getByRole("button",{name:"立即到达"}).click();
 const room=page.getByRole("region",{name:"企业总部室内场景"});
 const finance=room.getByRole("button",{name:"打开开业财务中心"});
 const table=room.locator(".gx-board-table");
 await expect(finance).toBeVisible();
 const financeBox=await finance.boundingBox();
 const tableBox=await table.boundingBox();
 expect(financeBox).not.toBeNull();
 expect(tableBox).not.toBeNull();
 const overlaps=financeBox!.x<tableBox!.x+tableBox!.width&&financeBox!.x+financeBox!.width>tableBox!.x&&financeBox!.y<tableBox!.y+tableBox!.height&&financeBox!.y+financeBox!.height>tableBox!.y;
 expect(overlaps).toBe(false);
 await finance.click();
 const detail=page.getByLabel("开业财务中心详情");
 await expect(detail.getByRole("link",{name:"查看系统账务"})).toHaveAttribute("href",/\?tenant=tenant-demo&case=INC-DEMO-001&view=ledger#finance_workspace$/);
 await expect(detail.getByRole("link",{name:"查看财务报表"})).toHaveAttribute("href",/view=reports#finance_workspace$/);
});
