import{ArrowLeft,Bot,Building2,Clock3,ExternalLink,FileCheck2,Globe2,ShieldCheck,UserRound}from"lucide-react";
import{useEffect,useState}from"react";
import{iaosWorkItemURL,loadGameProjection,resolveWorldWorkItem}from"../../game/api";
import type{GameBuilding,GameProjection}from"../../game/types";
import{BrandStudio}from"./BrandStudio";
import{IsometricCanvas}from"./IsometricCanvas";
import"./EnterpriseGenesisGame.css";
import"./GameSceneAssets.css";
const chapters=["创业构想","企业身份","创始办公室","公司登记","银行资本","人才治理","开业准备","经营世界"];
const money=(v:string)=>new Intl.NumberFormat("zh-CN",{style:"currency",currency:"CNY",maximumFractionDigits:0}).format(Number(v)/100);
export function EnterpriseGenesisGame({onExit}:{onExit:()=>void}){
 const params=new URLSearchParams(window.location.hash.split("?")[1]??""),caseCode=params.get("case")??"INC-GAME-DEMO-001",handedOff=params.get("auth_token"),tenantParam=params.get("tenant");
 const[frame,setFrame]=useState(0),[revision,setRevision]=useState(0),[projection,setProjection]=useState<GameProjection|null>(null),[selected,setSelected]=useState<GameBuilding|null>(null),[tab,setTab]=useState<"tasks"|"agents"|"evidence">("tasks"),[speed,setSpeed]=useState<0|1|2|4>(0),[error,setError]=useState("");
 useEffect(()=>{if(handedOff)localStorage.setItem("iaos_token",handedOff);if(tenantParam)localStorage.setItem("aese_iaos_tenant_id",tenantParam)},[handedOff,tenantParam]);
 useEffect(()=>{const c=new AbortController();loadGameProjection(caseCode,frame,c.signal).then(p=>{setProjection(p);setSelected((p.chapter==="operating_world"?p.buildings.find(b=>b.kind==="headquarters"):undefined)??p.buildings.find(b=>b.available)??null);setError("")}).catch(e=>{if(e.name!=="AbortError")setError(String(e))});return()=>c.abort()},[caseCode,frame,revision]);
 if(error)return <main className="gx-state" role="alert"><strong>企业创生世界无法加载</strong><p>{error}</p><button onClick={onExit}>返回</button></main>;
 if(!projection)return <main className="gx-state" aria-live="polite">正在生成企业世界投影…</main>;
 const active=Math.min(7,Math.round(projection.lifecycle.progress*7/100));
 return <div className="gx-shell">
  <header className="gx-header"><button onClick={onExit}><ArrowLeft/>生命周期</button><div><span>ENTERPRISE GENESIS · 企业创生</span><h1>{projection.brand.company_name||"尚未命名的企业"}</h1></div><div className="gx-clock"><Clock3/><small>虚拟时间</small><strong>{new Date(projection.sim_time).toLocaleString("zh-CN",{hour12:false})}</strong></div></header>
  <nav className="gx-chapters" aria-label="企业创生章节">{chapters.map((c,i)=><button key={c} disabled={i>active} aria-current={i===active?"step":undefined} onClick={()=>setFrame(i)}><span>{i+1}</span>{c}</button>)}</nav>
  {(projection.chapter==="founder_intent"||projection.chapter==="founder_office")&&<BrandStudio caseCode={caseCode}/>}
  <main className="gx-main"><section className="gx-world"><div className="gx-title"><span>{projection.chapter.replaceAll("_"," ")}</span><h2>{projection.lifecycle.current_step}</h2><p>场景只展示 IAOS / World 已提交证据</p></div>
   <div className="gx-isometric" role="img" aria-label="企业创生等距世界地图"><IsometricCanvas buildings={projection.buildings}/>{projection.buildings.map(b=><button key={b.code} className={`gx-building ${selected?.code===b.code?"selected":""}`} style={{"--gx-x":b.x,"--gx-y":b.y} as React.CSSProperties} disabled={!b.available} onClick={()=>setSelected(b)}><i/><Building2/><strong>{b.label}</strong><small>{b.state}</small></button>)}</div>
   <div className="gx-controls"><button disabled={!frame} onClick={()=>setFrame(v=>v-1)}>上一步</button><div className="gx-speed" aria-label="虚拟时间倍率">{([0,1,2,4]as const).map(v=><button key={v} aria-pressed={speed===v} onClick={()=>setSpeed(v)}>{v?v+"×":"暂停"}</button>)}</div><strong>{projection.lifecycle.progress}%</strong><button disabled={frame===7} onClick={()=>{setSpeed(0);setFrame(v=>v+1)}}>推进已提交状态</button></div>
  </section><aside className="gx-desk"><div className="gx-place"><Building2/><div><small>当前地点</small><strong>{selected?.label??"企业世界"}</strong><code>{selected?.evidence_ref??"持久投影"}</code></div></div>
   <div className="gx-resources">{[["创始人现金",projection.resources.founder_cash],["公司现金",projection.resources.company_cash],["实缴资本",projection.resources.capital_paid],["预算授权",projection.resources.budget_authorized]].map(([k,v])=><article key={k as string}><small>{k as string}</small><strong>{money((v as {value:string}).value)}</strong></article>)}</div>
   <div className="gx-tabs" role="tablist"><button role="tab" aria-selected={tab==="tasks"} onClick={()=>setTab("tasks")}><FileCheck2/>任务</button><button role="tab" aria-selected={tab==="agents"} onClick={()=>setTab("agents")}><Bot/>员工</button><button role="tab" aria-selected={tab==="evidence"} onClick={()=>setTab("evidence")}><ShieldCheck/>证据</button></div>
   <div className="gx-panel">{tab==="tasks"&&projection.work_items.map(i=>{const sequence=Number(i.work_item_id.replace(/\D/g,""));const actionable=i.status!=="completed"&&i.status!=="locked";return <article key={i.work_item_id}><div><b>{i.title}</b><code>{i.capability}</code></div><em>{i.requires_me?"需要我":i.status==="completed"?"已完成":i.status==="locked"?"未解锁":i.status==="waiting_world"?"等待世界":i.status==="waiting_approval"?"等待审批":"执行中"}</em>{actionable&&(i.kind==="world_wait"?<button className="gx-task-action" onClick={async()=>{setError("");try{const result=i.capability.includes("appointment")?"accepted":i.capability.includes("bank")?"opened":"registered";await resolveWorldWorkItem(caseCode,sequence,i.capability,projection.world_run_id,result);setRevision(v=>v+1)}catch(e){setError(e instanceof Error?e.message:String(e))}}}><Globe2/>模拟世界同意</button>:<a className="gx-task-action" href={iaosWorkItemURL(caseCode,sequence,i.capability)} target="_blank" rel="noreferrer">在 IAOS 处理<ExternalLink/></a>)}</article>})}{tab==="agents"&&projection.actors.map(a=><article key={a.actor_id}>{a.actor_type==="human"?<UserRound/>:<Bot/>}<div><b>{a.display_name}</b><small>{a.position} · {a.state}</small></div></article>)}{tab==="evidence"&&projection.evidence_refs.map(e=><article key={e.ref}><ShieldCheck/><div><b>{e.kind}</b><code>{e.ref}</code></div></article>)}</div>
  </aside></main>
 </div>
}
