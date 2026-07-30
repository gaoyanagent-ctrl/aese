import { ArrowLeft, ArrowRight, Building2, LoaderCircle, ShieldCheck, Sparkles } from "lucide-react";
import { useEffect, useRef, useState } from "react";
import { createGenesisWorkspace } from "../../game/api";
import type { GenesisWorkspaceResult } from "../../game/types";
import "./GenesisOnboarding.css";

type Draft={
  display_name:string;
  template_key:"manufacturing-enterprise";
  region:"CN-JS";
  timezone:"Asia/Shanghai";
  realism_level:"standard"|"strict";
  data_retention_confirmed:boolean;
};
const draftKey="aese_genesis_onboarding_draft";
const idempotencyKey="aese_genesis_onboarding_idempotency_key";
const initialDraft:Draft={display_name:"我的新企业",template_key:"manufacturing-enterprise",region:"CN-JS",timezone:"Asia/Shanghai",realism_level:"standard",data_retention_confirmed:false};
const steps=["创业项目","行业模板","区域与时区","真实性与数据","确认创建"];

function restoreDraft():Draft{
  try{return{...initialDraft,...JSON.parse(localStorage.getItem(draftKey)??"{}") as Partial<Draft>}}catch{return initialDraft}
}

export function restoreGenesisOnboardingIdempotencyKey(){
  const existing=localStorage.getItem(idempotencyKey)?.trim();
  if(existing)return existing;
  const created=`genesis-create-${Date.now()}-${crypto.randomUUID?.()??Math.random().toString(16).slice(2)}`;
  localStorage.setItem(idempotencyKey,created);
  return created;
}

export function clearGenesisOnboardingIdempotencyKey(){
  localStorage.removeItem(idempotencyKey);
}

export function GenesisOnboarding({onBack,onReady}:{onBack:()=>void;onReady:(workspace:GenesisWorkspaceResult)=>void}){
  const[draft,setDraft]=useState<Draft>(restoreDraft);
  const[step,setStep]=useState(0),[busy,setBusy]=useState(false),[error,setError]=useState("");
  const[provisioned,setProvisioned]=useState<GenesisWorkspaceResult|null>(null);
  const key=useRef(restoreGenesisOnboardingIdempotencyKey());
  useEffect(()=>{localStorage.setItem(draftKey,JSON.stringify(draft))},[draft]);
  const valid=[
    draft.display_name.trim().length>=2,
    draft.template_key==="manufacturing-enterprise",
    draft.region==="CN-JS"&&draft.timezone==="Asia/Shanghai",
    draft.data_retention_confirmed,
    true,
  ][step];
  const create=async()=>{
    setBusy(true);setError("");
    try{
      const workspace=await createGenesisWorkspace({...draft,idempotency_key:key.current});
      localStorage.removeItem(draftKey);clearGenesisOnboardingIdempotencyKey();setProvisioned(workspace);
    }catch(reason){setError(reason instanceof Error?reason.message:String(reason))}
    finally{setBusy(false)}
  };
  if(provisioned)return <main className="gx-onboarding">
    <header><button onClick={onBack}><ArrowLeft aria-hidden="true"/>返回主页</button><span>ENTERPRISE GENESIS · READY</span></header>
    <section className="gx-onboarding__card">
      <div className="gx-onboarding__icon"><ShieldCheck aria-hidden="true"/></div>
      <p className="gx-onboarding__eyebrow">独立企业空间已激活</p>
      <h1>八个可恢复检查点均已留下证据</h1>
      <div className="gx-onboarding__choice">
        <strong>{provisioned.display_name}</strong>
        <p>Workspace {provisioned.workspace_id} · Tenant {provisioned.tenant_id}</p>
        <ol>{(provisioned.steps??[]).map(item=><li key={item.step_key} className={item.status==="completed"?"active":""}><span>{item.status==="completed"?"✓":"…"}</span><div><strong>{item.step_key}</strong><small>{item.evidence_ref}</small></div></li>)}</ol>
      </div>
      <button className="gx-onboarding__submit" onClick={()=>onReady(provisioned)}><Sparkles/>进入 AI 身份工作室</button>
      <p className="gx-onboarding__security"><ShieldCheck/>当前浏览器仅取得 tenant-scoped genesis_owner 会话。</p>
    </section>
  </main>;
  return <main className="gx-onboarding">
    <header><button onClick={onBack}><ArrowLeft aria-hidden="true"/>返回主页</button><span>ENTERPRISE GENESIS · ZERO START</span></header>
    <section className="gx-onboarding__card">
      <div className="gx-onboarding__icon"><Building2 aria-hidden="true"/></div>
      <p className="gx-onboarding__eyebrow">第 {step+1} 步 · {steps[step]}</p>
      <h1>{step===0?"先定义创业项目":step===1?"选择制造企业运行模板":step===2?"确认业务区域与时间口径":step===3?"选择仿真真实性与数据保留":"确认并创建独立企业空间"}</h1>
      <div className="gx-onboarding__progress" aria-label="创建进度">{steps.map((label,index)=><button key={label} className={index===step?"active":index<step?"done":""} onClick={()=>index<=step&&setStep(index)}><span>{index+1}</span>{label}</button>)}</div>
      {step===0&&<label>创业项目名称<input value={draft.display_name} minLength={2} maxLength={80} onChange={event=>setDraft({...draft,display_name:event.target.value})}/><small>这是创建阶段名称；正式公司名称稍后由身份工作室产生。</small></label>}
      {step===1&&<div className="gx-onboarding__choice"><strong>制造企业 · M9 创生模板</strong><p>安装成立治理语义、25 个 Capability、23 个工作项、审批、Agent 与开业财务。</p><span>模板版本：m9/1.8.0</span></div>}
      {step===2&&<div className="gx-onboarding__choice"><strong>中国 · 江苏</strong><p>业务时区固定为 Asia/Shanghai；金额使用 CNY minor unit，不使用浮点舍入。</p><span>区域代码：CN-JS</span></div>}
      {step===3&&<><label>真实性等级<select value={draft.realism_level} onChange={event=>setDraft({...draft,realism_level:event.target.value as Draft["realism_level"]})}><option value="standard">标准经营仿真</option><option value="strict">严格治理与失败注入</option></select></label><label className="gx-onboarding__consent"><input type="checkbox" checked={draft.data_retention_confirmed} onChange={event=>setDraft({...draft,data_retention_confirmed:event.target.checked})}/><span>确认保存 Workspace、World Run、流程、审计和 CreativeJob 证据；不保存平台管理员凭据。</span></label></>}
      {step===4&&<div className="gx-onboarding__choice"><strong>{draft.display_name}</strong><p>制造企业 · CN-JS · Asia/Shanghai · {draft.realism_level==="strict"?"严格治理":"标准仿真"}</p><span>平台将分配不可预测的 Workspace、Tenant、World Run 与案件编码。</span></div>}
      {error&&<p className="gx-onboarding__error" role="alert">{error}</p>}
      <div className="gx-onboarding__wizard-actions">
        <button disabled={busy||step===0} onClick={()=>setStep(step-1)}><ArrowLeft/>上一步</button>
        {step<4?<button disabled={!valid} onClick={()=>setStep(step+1)}>下一步<ArrowRight/></button>:<button className="gx-onboarding__submit" disabled={busy} onClick={()=>void create()}>{busy?<LoaderCircle className="gx-spin"/>:<Sparkles/>}{busy?"正在执行八个可恢复检查点…":"创建空间并进入身份工作室"}</button>}
      </div>
      <p className="gx-onboarding__security"><ShieldCheck/>草稿自动保存；Tenant 由 IAOS 服务端分配，浏览器不能指定。</p>
    </section>
  </main>
}
