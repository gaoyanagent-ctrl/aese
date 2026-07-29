import { ArrowLeft, Bot, Check, LoaderCircle, MapPin, MessageSquareText, ShieldCheck, Sparkles } from "lucide-react";
import { useMemo, useState } from "react";
import { analyzeFounderIntent, createIncorporationCase, generateCompanyNames } from "../../game/api";
import type { FounderIntentRequest, NamingProposal } from "../../game/types";
import { FounderOfficeCanvas } from "./FounderOfficeCanvas";
import "./FounderOfficeRPG.css";

const avatars=[
 {name:"林澜",role:"产品型创始人",tone:"冷静、工程导向",skin:"#d4a070",hair:"#25201d"},
 {name:"周启",role:"市场型创始人",tone:"果断、增长导向",skin:"#bd825a",hair:"#171715"},
 {name:"沈知夏",role:"长期主义创始人",tone:"稳健、组织导向",skin:"#ca946b",hair:"#31241f"},
 {name:"陆行舟",role:"产业型创始人",tone:"务实、制造导向",skin:"#a96e50",hair:"#201a18"},
];
const choices={
 industry:["新能源汽车热管理","工业智能装备","绿色能源技术"],
 customer:["新能源汽车制造商","高端装备制造企业","成长型工业企业"],
 offering:["电池冷却板与热管理系统","工业软件与数字化服务","关键零部件研发制造"],
 trait:["可靠 · 工程 · 长期主义","敏捷 · 创新 · 客户成功","稳健 · 品质 · 产业深耕"],
};
const lines=[
 ["欢迎来到你的创始办公室。先选择一个代表你的虚拟形象。","这不是审批表。接下来我会把开办企业所需的信息，逐项变成经营选择。"],
 ["我们先决定要进入哪个产业。产业会影响后续客户、产品和风险事件。"],
 ["很好。第一批客户是谁？这会决定企业的价值主张。"],
 ["你准备靠什么产品或服务赢得市场？"],
 ["最后定下公司的性格。数字员工会据此调整建议风格。"],
 ["我已经整理出创业构想。你可以补充一句自己的想法，然后让 MiniMax M3 形成身份提案。"],
 ["身份提案已经送达。选一个你愿意带领十年的名字。"],
 ["正式提交前，确认企业落在哪里、准备经营什么。之后 IAOS 会创建受治理的设立案。"],
];
const split=(value:string)=>value.split(" · ");

export function FounderOfficeRPG({caseCode,onExit,onCreated}:{caseCode:string;onExit:()=>void;onCreated:()=>void}){
 const[stage,setStage]=useState(0),[avatar,setAvatar]=useState(0),[industry,setIndustry]=useState(""),[customer,setCustomer]=useState(""),[offering,setOffering]=useState(""),[trait,setTrait]=useState("");
 const[idea,setIdea]=useState("我想创建一家真正理解制造现场、重视工程可靠性并坚持长期主义的企业。");
 const[names,setNames]=useState<NamingProposal[]>([]),[selected,setSelected]=useState(""),[address,setAddress]=useState("江苏省苏州市工业园区创生大道1号"),[scope,setScope]=useState("工业产品的研发、制造、销售与技术服务");
 const[busy,setBusy]=useState(false),[error,setError]=useState("");
 const progress=Math.round(stage/7*100);
 const line=lines[Math.min(stage,7)];
 const next=(value:string)=>{
  if(stage===1)setIndustry(value);if(stage===2)setCustomer(value);if(stage===3)setOffering(value);if(stage===4)setTrait(value);
  setStage(current=>current+1);
 };
 const request=useMemo<FounderIntentRequest>(()=>({tenant_id:localStorage.getItem("aese_iaos_tenant_id")??"",case_code:caseCode,raw_idea:idea,industry,customers:[customer],offerings:[offering],brand_traits:split(trait),capital_minor:"100000000",risk_appetite:"balanced"}),[caseCode,idea,industry,customer,offering,trait]);
 const generate=async()=>{setBusy(true);setError("");try{const intent=await analyzeFounderIntent(request);setNames(await generateCompanyNames(intent));setStage(6)}catch(reason){setError(reason instanceof Error?reason.message:String(reason))}finally{setBusy(false)}};
 const commit=async()=>{const choice=names.find(item=>item.proposal_id===selected);if(!choice)return;setBusy(true);setError("");try{localStorage.setItem("aese_founder_avatar",JSON.stringify({...avatars[avatar],index:avatar}));await createIncorporationCase({case_code:caseCode,case_name:`${choice.short_name}企业设立案`,proposed_company_name:choice.chinese_name,registered_address:address,business_scope:scope});onCreated()}catch(reason){setError(reason instanceof Error?reason.message:String(reason))}finally{setBusy(false)}};
 return <main className="gx-rpg">
  <header className="gx-rpg__header"><button onClick={onExit}><ArrowLeft aria-hidden="true"/>返回企业大厅</button><div><small>CHAPTER 01</small><strong>创始人办公室</strong></div><div className="gx-rpg__quest"><span>主线任务 · 定义企业身份</span><progress value={progress} max="100"/><em>{progress}%</em></div></header>
  <section className="gx-rpg__world">
   <FounderOfficeCanvas avatar={avatar} stage={stage}/>
   <div className="gx-rpg__hud"><span><MapPin aria-hidden="true"/>临时创始办公室</span><span><ShieldCheck aria-hidden="true"/>IAOS 尚未提交</span></div>
   <div className="gx-rpg__founder"><i style={{"--skin":avatars[avatar].skin,"--hair":avatars[avatar].hair} as React.CSSProperties}/><div><small>你 · 创始人</small><strong>{avatars[avatar].name}</strong><span>{avatars[avatar].role}</span></div></div>
   <aside className="gx-rpg__quest-list"><p>当前任务</p>{["建立创始人档案","确定产业方向","选择目标客户","确定核心产品","定义品牌性格","生成企业身份","提交设立案"].map((item,index)=><div key={item} className={index<stage?"done":index===stage?"active":""}>{index<stage?<Check/>:<span>{index+1}</span>}<b>{item}</b></div>)}</aside>
  </section>
  <section className="gx-rpg__dialog" aria-live="polite">
   <div className="gx-rpg__advisor"><span><Bot aria-hidden="true"/></span><div><small>数字员工 · 创业顾问</small><strong>纪元</strong></div></div>
   <div className="gx-rpg__speech"><MessageSquareText aria-hidden="true"/><p>{line[0]}</p>{line[1]&&<span>{line[1]}</span>}
    {stage===0&&<div className="gx-rpg__avatars">{avatars.map((item,index)=><button key={item.name} className={avatar===index?"selected":""} onClick={()=>setAvatar(index)}><i style={{"--skin":item.skin,"--hair":item.hair} as React.CSSProperties}/><strong>{item.name}</strong><small>{item.role}</small></button>)}<button className="gx-rpg__continue" onClick={()=>setStage(1)}>使用这个形象 <ArrowLeft/></button></div>}
    {stage>=1&&stage<=4&&<div className="gx-rpg__choices">{Object.values(choices)[stage-1].map(choice=><button key={choice} onClick={()=>next(choice)}>{choice}<span>选择</span></button>)}</div>}
    {stage===5&&<div className="gx-rpg__idea"><label>补充你的创业宣言<textarea value={idea} onChange={event=>setIdea(event.target.value)}/></label><button onClick={()=>void generate()} disabled={busy}>{busy?<LoaderCircle className="gx-spin"/>:<Sparkles/>}{busy?"数字员工正在与 MiniMax M3 协作…":"让 AI 形成 4 组企业身份提案"}</button></div>}
    {stage===6&&<div className="gx-rpg__names">{names.map(item=><button key={item.proposal_id} className={selected===item.proposal_id?"selected":""} style={{"--brand":item.primary_color} as React.CSSProperties} onClick={()=>setSelected(item.proposal_id)}><i>{item.short_name.slice(0,1)}</i><div><strong>{item.chinese_name}</strong><small>{item.english_name}</small><p>{item.slogan}</p></div>{selected===item.proposal_id&&<Check/>}</button>)}{selected&&<button className="gx-rpg__continue" onClick={()=>setStage(7)}>我选择这个名字 <ArrowLeft/></button>}</div>}
    {stage===7&&<div className="gx-rpg__final"><label>注册地址<input value={address} onChange={event=>setAddress(event.target.value)}/></label><label>经营范围<textarea value={scope} onChange={event=>setScope(event.target.value)}/></label><button onClick={()=>void commit()} disabled={busy||!selected}>{busy?<LoaderCircle className="gx-spin"/>:<ShieldCheck/>}{busy?"正在写入 IAOS 企业事实…":"签署创始人指令并启动企业设立"}</button></div>}
    {error&&<p className="gx-rpg__error" role="alert">{error}</p>}
   </div>
  </section>
 </main>;
}
