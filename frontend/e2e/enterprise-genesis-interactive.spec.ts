import{expect,test}from"@playwright/test";

test("founder creates a new enterprise and operates the first IAOS work item in game",async({page,request})=>{
 test.setTimeout(60_000);
 const login=await request.post("http://127.0.0.1:8082/api/v1/auth/login",{data:{
  username:"founder-principal",password:"Founder-Lifecycle-2026!",tenant_id:"tenant-hctm-genesis"
 }});
 expect(login.ok()).toBeTruthy();
 const session=await login.json();
 const caseCode=`INC-UI-${Date.now()}`;
 await page.addInitScript(value=>{
  localStorage.setItem("iaos_token",value.token);
  localStorage.setItem("aese_iaos_tenant_id","tenant-hctm-genesis");
 },session);
 await page.goto(`/#enterprise-genesis?tenant=tenant-hctm-genesis&case=${caseCode}`);
 await expect(page.getByRole("heading",{name:"创建你的第一家企业"})).toBeVisible();
 await page.getByRole("button",{name:"生成公司身份候选"}).click();
 const proposal=page.locator(".gx-candidates>button").first();
 await expect(proposal).toBeVisible();
 await proposal.click();
 await page.getByRole("button",{name:"确认身份并创建企业"}).click();
 await expect(page.getByText("founder.resolution.prepare").first()).toBeVisible({timeout:15_000});
 await page.getByRole("button",{name:"进入操作"}).click();
 await expect(page.getByRole("heading",{name:"让企业设立专员起草创始决议"})).toBeVisible();
 await expect(page.getByLabel("创始办公室剧情")).toBeVisible();
 await page.getByRole("button",{name:"派遣数字员工"}).click();
 const remaining=[
  "founder.resolution.approve","capital.commitment.record","registration.package.validate",
  "registration.submit","registration.observation.commit","bank.account.opening.submit",
  "bank.account.observation.commit","capital.contribution.verify","organization.establish",
  "executive.appointment.propose","executive.appointment.acceptance.commit",
  "executive.appointment.approve","operating.mandate.grant","initial.budget.prepare",
  "initial.budget.approve","enterprise.readiness.evaluate"
 ];
 for(const capability of remaining){
  await expect(page.getByText(capability).first()).toBeVisible({timeout:15_000});
  const ready=page.locator(".gx-task-ready",{hasText:capability});
  await expect(ready).toBeVisible();
  await ready.getByRole("button",{name:"进入操作"}).click();
  await page.locator(".gx-action-submit").click();
  await expect(page.locator(".gx-action-card")).toBeHidden({timeout:15_000});
 }
 await expect(page.getByText("enterprise_operational_ready").first()).toBeVisible({timeout:15_000});
 const evidence=await request.get(`http://127.0.0.1:8082/api/v1/incorporations/${caseCode}/evidence`,{
  headers:{Authorization:`Bearer ${session.token}`,"X-Tenant-ID":"tenant-hctm-genesis"}
 });
 expect(evidence.ok()).toBeTruthy();
 expect(await evidence.text()).toContain("incorporation_case_opened");
});
