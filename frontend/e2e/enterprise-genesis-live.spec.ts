import{expect,test}from"@playwright/test";

test.skip(!process.env.GX_LIVE,"set GX_LIVE=1 for the deployed IAOS + AESE acceptance");

test("live Enterprise Genesis restores the governed M9 case",async({page,request},testInfo)=>{
 const caseCode=process.env.M9_CASE_CODE??"HCTM-TEST001";
 const login=await request.post("http://127.0.0.1:8082/api/v1/auth/login",{data:{username:"founder-principal",password:"Founder-Lifecycle-2026!",tenant_id:"tenant-hctm-genesis"}});
 expect(login.ok()).toBeTruthy();
 const session=await login.json();
 await page.addInitScript(value=>{localStorage.setItem("iaos_token",value.token);localStorage.setItem("aese_iaos_tenant_id","tenant-hctm-genesis")},session);
 await page.goto(`/#enterprise-genesis?tenant=tenant-hctm-genesis&case=${encodeURIComponent(caseCode)}`);
 if(process.env.M9_EXPECT_READY){
  await expect(page.getByText("enterprise_operational_ready").first()).toBeVisible();
  await expect(page.getByRole("heading",{name:"企业身份工作室"})).toHaveCount(0);
 }else{
  await expect(page.getByRole("heading",{name:"企业身份工作室"})).toBeVisible();
  await expect(page.getByText("founder.resolution.approve").first()).toBeVisible();
 }
 await page.getByRole("tab",{name:"员工"}).click();
 for(const agent of["企业设立专员","治理组织专员","法务合规专员","财务负责人","独立审计专员"]){
  await expect(page.getByText(agent,{exact:true}).first()).toBeVisible();
 }
 await page.reload();
 await expect(page.getByText(process.env.M9_EXPECT_READY?"enterprise_operational_ready":"founder.resolution.approve").first()).toBeVisible();
 await page.screenshot({path:testInfo.outputPath("enterprise-genesis-live.png"),fullPage:true});
});
