import{ArrowLeft,Bot,Building2,FileCheck2,Home,Landmark,MapPinned,Plus,ShieldCheck,UserRound,WalletCards,Volume2,VolumeX,Trophy}from"lucide-react";
import{useEffect,useRef,useState}from"react";
import{GameApiError,loadGameProjection}from"../../game/api";
import type{GameBuilding,GameProjection,GameWorkItem}from"../../game/types";
import{FounderOfficeRPG}from"./FounderOfficeRPG";
import{WorkItemActionPanel}from"./WorkItemActionPanel";
import{MissionBriefing}from"./MissionBriefing";
import{LocationScene,missionPlace}from"./LocationScene";
import"./EnterpriseGenesisGame.css";
import"./GameSceneAssets.css";
import"./WorkItemActionPanel.css";
import{gameSoundEnabled,playGameTone,setGameSoundEnabled}from"../../game/audio";
import{financeWorkspaceUrl}from"../../game/iaosLinks";

const chapters=["创业构想","企业身份","创始办公室","公司登记","银行资本","人才治理","开业准备","经营世界"];
const missionTitles:Record<string,string>={
 "founder.resolution.prepare":"形成第一份创始人决议","founder.resolution.approve":"召开临时董事会",
 "capital.commitment.record":"确认企业启动资本","registration.package.validate":"准备企业登记材料",
 "registration.submit":"前往政务中心提交登记","registration.observation.commit":"查看登记机构反馈",
 "bank.account.opening.submit":"选择合作银行并申请开户","bank.account.observation.commit":"查看银行尽调结果",
 "capital.contribution.verify":"注入并核验实缴资本","organization.establish":"建立企业初始组织",
 "executive.appointment.propose":"提名首届管理团队","executive.appointment.acceptance.commit":"与管理层候选人会谈",
 "executive.appointment.approve":"批准首届管理层任命","operating.mandate.grant":"授予管理层经营权限",
 "initial.budget.prepare":"编制首年启动预算","initial.budget.approve":"召开首年预算会议",
 "enterprise.readiness.evaluate":"完成企业开业检查"
};
const money=(v:string)=>new Intl.NumberFormat("zh-CN",{style:"currency",currency:"CNY",maximumFractionDigits:0}).format(Number(v)/100);

export function EnterpriseGenesisGame({onExit}:{onExit:()=>void}){
 const params=new URLSearchParams(window.location.hash.split("?")[1]??"");
 const caseCode=params.get("case")??"INC-GAME-DEMO-001",handedOff=params.get("auth_token"),tenantParam=params.get("tenant");
 const[revision,setRevision]=useState(0),[projection,setProjection]=useState<GameProjection|null>(null);
 const[selected,setSelected]=useState<GameBuilding|null>(null),[inside,setInside]=useState(false),[selectedTask,setSelectedTask]=useState<GameWorkItem|null>(null);
 const[travelling,setTravelling]=useState<GameBuilding|null>(null);
 const[currentPlace,setCurrentPlace]=useState("office");
 const travelTimer=useRef<number|undefined>(undefined);
 const[tab,setTab]=useState<"event"|"team"|"archive">("event"),[error,setError]=useState(""),[creating,setCreating]=useState(false);
 const[founderAvatar,setFounderAvatar]=useState(()=>{try{return Number(JSON.parse(localStorage.getItem("aese_founder_avatar")??"{}").index??0)}catch{return 0}});
 const[sound,setSound]=useState(gameSoundEnabled),[reward,setReward]=useState("");
 useEffect(()=>{if(handedOff)localStorage.setItem("iaos_token",handedOff);if(tenantParam)localStorage.setItem("aese_iaos_tenant_id",tenantParam)},[handedOff,tenantParam]);
 useEffect(()=>()=>{if(travelTimer.current)window.clearTimeout(travelTimer.current)},[]);
 useEffect(()=>{
  const controller=new AbortController();setProjection(null);
  loadGameProjection(caseCode,0,controller.signal).then(value=>{
   const mission=value.work_items.find(item=>item.status!=="completed"&&item.status!=="locked");
   setProjection(value);setCreating(false);setSelected(value.buildings.find(building=>building.kind===missionPlace(mission?.capability??""))??value.buildings.find(building=>building.available)??null);setError("");
  }).catch(reason=>{
   if(reason.name==="AbortError")return;
   if(reason instanceof GameApiError&&reason.status===404){setCreating(true);setError("");return}
   setError(String(reason));
  });
  return()=>controller.abort();
 },[caseCode,revision]);
 const startNew=()=>{
  const tenant=localStorage.getItem("aese_iaos_tenant_id")??tenantParam??"tenant-hctm-genesis";
  const fresh=`INC-PLAYER-${Date.now()}`;
  window.location.hash=`enterprise-genesis?tenant=${encodeURIComponent(tenant)}&case=${fresh}`;
  window.location.reload();
 };
 const enterLocation=(building:GameBuilding|null)=>{
  if(!building||!building.available)return;
  playGameTone("travel");
  setSelected(building);setInside(false);setTravelling(building);
  if(travelTimer.current)window.clearTimeout(travelTimer.current);
  travelTimer.current=window.setTimeout(()=>{setCurrentPlace(building.kind);setTravelling(null);setInside(true)},420);
 };
 const finishTravel=()=>{if(travelTimer.current)window.clearTimeout(travelTimer.current);if(travelling)setCurrentPlace(travelling.kind);setTravelling(null);setInside(true)};
 if(error)return <main className="gx-state" role="alert"><strong>企业创生世界无法加载</strong><p>{error}</p><button onClick={onExit}>返回</button></main>;
 if(creating)return <FounderOfficeRPG caseCode={caseCode} onExit={onExit} onCreated={()=>setRevision(value=>value+1)}/>;
 if(!projection)return <main className="gx-state" aria-live="polite">正在读取 IAOS 企业状态…</main>;
 const active=Math.min(7,Math.round(projection.lifecycle.progress*7/100));
 const currentMission=projection.work_items.find(item=>item.status!=="completed"&&item.status!=="locked");
 return <div className="gx-shell">
  <header className="gx-header"><button onClick={onExit}><ArrowLeft/>企业大厅</button><div><span>ENTERPRISE GENESIS · 企业创生</span><h1>{projection.brand.company_name||"尚未命名的企业"}</h1></div><div className="gx-header-actions"><button aria-label={sound?"关闭游戏音效":"开启游戏音效"} title={sound?"关闭游戏音效":"开启游戏音效"} onClick={()=>{const next=!sound;setSound(next);setGameSoundEnabled(next);if(next)playGameTone("inspect")}}>{sound?<Volume2/>:<VolumeX/>}</button><button className="gx-header-founder" title="点击切换本机创始人头像" onClick={()=>{const next=(founderAvatar+1)%4;setFounderAvatar(next);localStorage.setItem("aese_founder_avatar",JSON.stringify({index:next}))}}><i data-avatar={founderAvatar}/><span><small>你 · 创始人</small><strong>{localStorage.getItem("aese_genesis_username")??"Founder"}</strong></span></button><button onClick={startNew}><Plus/>新建企业</button></div></header>
  <nav className="gx-chapters" aria-label="企业创生旅程">{chapters.map((chapter,index)=><div key={chapter} data-state={index<active?"completed":index===active?"current":"locked"} aria-current={index===active?"step":undefined}><span>{index<active?"✓":index+1}</span>{chapter}</div>)}</nav>
  <main className="gx-main"><section className="gx-world"><div className="gx-title"><span>{projection.chapter.replaceAll("_"," ")}</span><h2>{missionTitles[currentMission?.capability??""]??"企业世界正在运行"}</h2><p>前往带有事件标记的地点，与人物和场景互动</p></div>
   <div className="gx-isometric" role="img" aria-label="企业创生城市地图">{currentMission&&missionPlace(currentMission.capability)!=="office"&&<svg className={`gx-city-route route-${missionPlace(currentMission.capability)}`} viewBox="0 0 1000 620" preserveAspectRatio="none" aria-hidden="true"><path d={missionPlace(currentMission.capability)==="government"?"M250 420 C210 350 210 275 240 205":missionPlace(currentMission.capability)==="bank"?"M250 420 C340 330 430 245 540 185":"M250 420 C420 430 610 360 770 250"}/></svg>}<div className={`gx-map-player at-${travelling?.kind??currentPlace} ${travelling?"walking":""}`} aria-label="玩家在城市中的位置"><i data-avatar={founderAvatar}/><span>你</span></div>{projection.buildings.map(building=>{const Icon=building.kind==="government"?Landmark:building.kind==="bank"?WalletCards:building.kind==="office"?Home:Building2;const hasMission=currentMission&&missionPlace(currentMission.capability)===building.kind;return <button key={building.code} data-kind={building.kind} className={`gx-building ${selected?.code===building.code?"selected":""} ${hasMission?"has-mission":""}`} disabled={!building.available} onClick={()=>enterLocation(building)}><span className="gx-building-icon"><Icon/></span><strong>{building.label}</strong><small>{hasMission?"有新事件":building.available?"可进入":"尚未开放"}</small>{hasMission&&<em>!</em>}</button>})}</div>
   {currentMission&&<MissionBriefing item={currentMission} onStart={()=>enterLocation(projection.buildings.find(building=>building.kind===missionPlace(currentMission.capability))??selected)}/>}
   <div className="gx-world-status"><MapPinned/><span><small>企业城市</small><strong>{projection.lifecycle.progress}% 创生进度</strong></span><i><b style={{width:`${projection.lifecycle.progress}%`}}/></i></div>
   {travelling&&<div className={`gx-travel gx-travel-${travelling.kind}`} role="status" aria-live="polite"><MapPinned/><div><small>正在穿过企业城市</small><strong>前往{travelling.label}</strong><span>{currentMission?missionTitles[currentMission.capability]:"查看地点资产"}</span></div><i><b/></i><button onClick={finishTravel}>立即到达</button></div>}
   {inside&&selected&&<LocationScene building={selected} projection={projection} mission={currentMission&&missionPlace(currentMission.capability)===selected.kind?currentMission:undefined} onStart={()=>currentMission&&setSelectedTask(currentMission)} onBack={()=>setInside(false)}/>}
  </section><aside className="gx-desk"><div className="gx-place"><MapPinned/><div><small>当前视角</small><strong>{inside?selected?.label:"企业城市地图"}</strong><code>{inside?selected?.evidence_ref??"IAOS 地点投影":"选择有事件标记的建筑进入场景"}</code></div></div>
   <div className="gx-resources">{[["创始人现金",projection.resources.founder_cash],["公司现金",projection.resources.company_cash],["实缴资本",projection.resources.capital_paid],["预算授权",projection.resources.budget_authorized]].map(([label,value])=><article key={label as string}><small>{label as string}</small><strong>{money((value as {value:string}).value)}</strong></article>)}</div>
   <div className="gx-tabs" role="tablist"><button role="tab" aria-selected={tab==="event"} onClick={()=>setTab("event")}><FileCheck2/>当前事件</button><button role="tab" aria-selected={tab==="team"} onClick={()=>setTab("team")}><Bot/>团队</button><button role="tab" aria-selected={tab==="archive"} onClick={()=>setTab("archive")}><ShieldCheck/>治理档案</button></div>
   <div className="gx-panel">{tab==="event"&&(currentMission?<article className="gx-task-ready"><div><b>{missionTitles[currentMission.capability]??currentMission.title}</b><small>前往{projection.buildings.find(building=>building.kind===missionPlace(currentMission.capability))?.label}处理</small></div><em>{currentMission.requires_me?"等待创始人":"数字员工协作"}</em><button className="gx-task-action" onClick={()=>enterLocation(projection.buildings.find(building=>building.kind===missionPlace(currentMission.capability))??selected)}>进入地点</button></article>:<article><ShieldCheck/><div><b>当前章节已完成</b><small>企业世界没有待处理事件</small></div></article>)}{tab==="team"&&projection.actors.map(actor=><article key={actor.actor_id}>{actor.actor_type==="human"?<UserRound/>:<Bot/>}<div><b>{actor.display_name}</b><small>{actor.position} · {actor.state}</small></div></article>)}{tab==="archive"&&<><article><ShieldCheck/><div><b>企业大事记</b><small>{projection.work_items.filter(item=>item.status==="completed").length} 项 committed 事件 · 只读</small></div></article>{projection.finance_opening?.ready&&<article className="gx-finance-archive"><WalletCards/><div><b>实缴资本开业入账</b><small>{projection.finance_opening.journal_entry_no} · 借贷平衡 {money(String(projection.finance_opening.debit_minor))}</small><code>{projection.finance_opening.book_code} · {projection.finance_opening.accounting_standard} · {projection.finance_opening.period_code}</code><strong>银行日记账</strong>{projection.finance_opening.bank_journal.map(line=><span key={`${line.entry_no}-${line.business_date}`}>{line.business_date} · {line.description} · 收入 {money(String(line.debit_minor))} · 余额 {money(String(line.balance_minor))}</span>)}<strong>总账与试算平衡</strong>{projection.finance_opening.general_ledger.map(line=><span key={line.account_code}>{line.account_code} {line.account_name} · 期末 {money(String(line.closing_balance_minor))}</span>)}<strong>开业资产负债表 · {projection.finance_opening.opening_balance_sheet.as_of}</strong><span>资产 {money(String(projection.finance_opening.opening_balance_sheet.total_assets_minor))}</span><span>负债 {money(String(projection.finance_opening.opening_balance_sheet.total_liabilities_minor))} · 所有者权益 {money(String(projection.finance_opening.opening_balance_sheet.total_equity_minor))}</span><span>{projection.finance_opening.opening_balance_sheet.balanced?"资产 = 负债 + 所有者权益":"报表不平衡，企业就绪已阻断"}</span><nav className="gx-finance-archive-actions" aria-label="IAOS 财务系统入口"><a href={financeWorkspaceUrl(projection,"ledger")} target="_blank" rel="noreferrer">查看系统账务</a><a href={financeWorkspaceUrl(projection,"reports")} target="_blank" rel="noreferrer">查看财务报表</a></nav></div></article>}{projection.work_items.filter(item=>item.status==="completed").map((item,index)=><article className="gx-history-event" key={item.work_item_id}><span>{index+1}</span><div><b>{missionTitles[item.capability]??item.title}</b><small>{item.kind} · 已提交</small><code>{item.evidence_ref}</code></div></article>)}</>}</div>
  </aside></main>
  {reward&&<div className="gx-reward-toast" role="status"><Trophy/><div><strong>企业里程碑达成</strong><span>{reward}</span></div></div>}
  {selectedTask&&<div className="gx-action-overlay" role="dialog" aria-modal="true"><button className="gx-action-backdrop" aria-label="关闭" onClick={()=>setSelectedTask(null)}/><WorkItemActionPanel item={selectedTask} projection={projection} onClose={()=>setSelectedTask(null)} onDone={()=>{const title=missionTitles[selectedTask.capability]??selectedTask.title;setSelectedTask(null);setInside(false);setReward(title);playGameTone("success");window.setTimeout(()=>setReward(""),2600);setRevision(value=>value+1)}}/></div>}
 </div>
}
