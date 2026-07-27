import{Check,LoaderCircle,Sparkles,TriangleAlert}from"lucide-react";
import{useState}from"react";
import{analyzeFounderIntent,createIncorporationCase,generateCompanyNames}from"../../game/api";
import type{FounderIntentRequest,NamingProposal}from"../../game/types";
import"./BrandStudio.css";

const initial:FounderIntentRequest={
 tenant_id:"tenant-demo",case_code:"INC-GAME-DEMO-001",
 raw_idea:"创建一家服务新能源汽车制造商的工业热管理公司，产品可靠、工程化并坚持长期主义。",
 industry:"热管理",customers:["新能源汽车制造商"],offerings:["电池冷却板"],
 brand_traits:["可靠","工程","长期主义"],capital_minor:"100000000",risk_appetite:"balanced"
};
const words=(value:string)=>value.split(/[，,\s]+/).map(v=>v.trim()).filter(Boolean);

export function BrandStudio({caseCode}:{caseCode:string}){
 const[form,setForm]=useState({...initial,case_code:caseCode}),[names,setNames]=useState<NamingProposal[]>([]),[selected,setSelected]=useState(""),[address,setAddress]=useState("江苏省苏州市工业园区创生大道1号"),[scope,setScope]=useState("工业热管理系统与电池冷却板的研发、制造及技术服务"),[busy,setBusy]=useState(false),[committing,setCommitting]=useState(false),[committed,setCommitted]=useState(false),[error,setError]=useState("");
 const create=async()=>{setBusy(true);setError("");try{const intent=await analyzeFounderIntent(form);setNames(await generateCompanyNames(intent));setSelected("")}catch(e){setError(e instanceof Error?e.message:String(e))}finally{setBusy(false)}};
 return <section className="gx-brand" aria-labelledby="brand-studio-title">
  <div className="gx-brand-heading"><div><span>AI CREATIVE · 候选区</span><h2 id="brand-studio-title">企业身份工作室</h2><p>AI 只生成候选；选择后仍需 IAOS Capability 提交，当前不会改变公司事实。</p></div><Sparkles/></div>
  <div className="gx-brand-grid"><form onSubmit={e=>{e.preventDefault();void create()}}>
   <label>你的创业想法<textarea value={form.raw_idea} minLength={12} onChange={e=>setForm({...form,raw_idea:e.target.value})}/></label>
   <div><label>行业<input value={form.industry} onChange={e=>setForm({...form,industry:e.target.value})}/></label><label>品牌特质<input value={form.brand_traits.join("，")} onChange={e=>setForm({...form,brand_traits:words(e.target.value)})}/></label></div>
   <div><label>目标客户<input value={form.customers.join("，")} onChange={e=>setForm({...form,customers:words(e.target.value)})}/></label><label>产品/服务<input value={form.offerings.join("，")} onChange={e=>setForm({...form,offerings:words(e.target.value)})}/></label></div>
   <button className="gx-generate" disabled={busy}>{busy?<LoaderCircle className="gx-spin"/>:<Sparkles/>}{busy?"正在生成候选…":"生成公司身份候选"}</button>
   {error&&<p className="gx-error" role="alert">{error}</p>}
  </form><div className="gx-candidates" aria-live="polite">
   {!names.length&&!busy&&<div className="gx-empty"><Sparkles/><strong>等待创业构想</strong><p>将生成 4 组名称、英文名、口号、色彩与风险提示。</p></div>}
   {names.map(n=><button key={n.proposal_id} type="button" className={selected===n.proposal_id?"selected":""} onClick={()=>setSelected(n.proposal_id)} style={{"--brand-color":n.primary_color} as React.CSSProperties}>
    <i>{n.short_name.slice(0,1)}</i><div><strong>{n.chinese_name}</strong><small>{n.english_name}</small><p>{n.slogan}</p><em><TriangleAlert/>现实工商核名/商标检索未完成</em></div>{selected===n.proposal_id&&<Check/>}
   </button>)}
   {selected&&<div className="gx-commit"><div className="gx-candidate-note"><Check/><span>已选为草稿候选，尚未成为 IAOS 正式企业事实。</span><button type="button" onClick={()=>setSelected("")}>撤销选择</button></div>
    <div className="gx-commit-fields"><label>注册地址<input value={address} onChange={e=>setAddress(e.target.value)}/></label><label>经营范围<input value={scope} onChange={e=>setScope(e.target.value)}/></label></div>
    <button className="gx-commit-button" type="button" disabled={committing||committed} onClick={async()=>{const choice=names.find(n=>n.proposal_id===selected);if(!choice)return;setCommitting(true);setError("");try{await createIncorporationCase({case_code:form.case_code,case_name:`${choice.short_name}企业设立案`,proposed_company_name:choice.chinese_name,registered_address:address,business_scope:scope});setCommitted(true)}catch(e){setError(e instanceof Error?e.message:String(e))}finally{setCommitting(false)}}}>{committed?<><Check/>已通过 incorporation.case.open 创建</>:committing?<><LoaderCircle className="gx-spin"/>正在提交 IAOS…</>:<><Check/>确认身份并创建企业</>}</button>
   </div>}
  </div></div>
 </section>
}
