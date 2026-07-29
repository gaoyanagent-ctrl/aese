import { BriefcaseBusiness, Building2, FileBadge2, Landmark, MonitorCog, ShieldCheck, Stamp, UsersRound, WalletCards } from "lucide-react";
import { useState } from "react";
import type { GameBuilding, GameProjection, GameWorkItem } from "../../game/types";
import { MissionBriefing } from "./MissionBriefing";
import { playGameTone } from "../../game/audio";
import { financeWorkspaceUrl } from "../../game/iaosLinks";
import "./LocationScene.css";

const placeCopy:Record<string,{eyebrow:string;title:string;description:string}>={
 office:{eyebrow:"创业街区 · 01",title:"创始办公室",description:"企业的第一处基地。决议、证照、印章和数字员工都会在这里留下永久痕迹。"},
 government:{eyebrow:"行政服务区 · 02",title:"政务服务中心",description:"在登记窗口递交资料、接收补正意见，并领取企业法律身份。"},
 bank:{eyebrow:"金融商务区 · 03",title:"合作银行",description:"比较银行方案、接受尽调、开立基本账户并完成资本到账。"},
 headquarters:{eyebrow:"总部商务区 · 04",title:"企业总部",description:"组织、管理层、经营授权和预算在这里逐步成为真实运营能力。"},
};
const npcByPlace:Record<string,{name:string;role:string;line:string}>={
 office:{name:"纪元",role:"企业设立专员",line:"我已经把今天需要决定的事项放在你的经营桌上。"},
 government:{name:"林岚",role:"企业登记窗口专员",line:"请先核对材料。缺少任何一项，我都会给出明确补正意见。"},
 bank:{name:"周衡",role:"企业客户经理",line:"开户不是形式审查，我们会核对业务、资金来源和实际经营地址。"},
 headquarters:{name:"顾远",role:"治理组织顾问",line:"职位、任命和授权是三件不同的事，我们会逐项建立。"},
};

export function missionPlace(capability:string){
 if(capability.startsWith("registration."))return"government";
 if(capability==="capital.contribution.post")return"headquarters";
 if(capability.startsWith("bank.")||capability.startsWith("capital."))return capability==="capital.commitment.record"?"office":"bank";
 if(capability.startsWith("finance.")||capability.startsWith("accounting.")||capability.startsWith("chart."))return"headquarters";
 if(capability.startsWith("organization.")||capability.startsWith("executive.")||capability.startsWith("operating.")||capability.startsWith("initial.")||capability.startsWith("enterprise."))return"headquarters";
 return"office";
}

export function LocationScene({building,projection,mission,onStart,onBack}:{building:GameBuilding;projection:GameProjection;mission?:GameWorkItem;onStart:()=>void;onBack:()=>void}){
 const copy=placeCopy[building.kind]??{eyebrow:"企业城市",title:building.label,description:"查看地点中的企业事件。"};
 const state=projection.lifecycle.state;
 const registered=Boolean(projection.brand.company_name)&&!["incorporation_opened","incorporation_case_opened","founder_resolution_approved","capital_commitments_confirmed","registration_submitted"].includes(state);
 const banked=Boolean(projection.resources.company_cash.value!=="0"||["bank_account_opened","capital_contribution_verified","organization_established","executive_appointment_proposed","executive_appointments_accepted","operating_mandates_activated","initial_budget_approved","enterprise_operational_ready"].includes(state));
 const organized=["organization_established","executive_appointment_proposed","executive_appointments_accepted","operating_mandates_activated","initial_budget_approved","enterprise_operational_ready"].includes(state);
 const npc=npcByPlace[building.kind]??npcByPlace.office;
 const npcState=mission?.status==="waiting_world"?"等待外部回复":mission?.status==="waiting_approval"?"等待你的决定":mission?.kind==="agent_task"?"正在准备材料":"可以交谈";
 type Inspectable={title:string;status:string;detail:string;source:string;actions?:Array<{label:string;href:string}>};
 const[inspected,setInspected]=useState<Inspectable|null>(null);
 const[playerTarget,setPlayerTarget]=useState(0);
 const financeWorkItems=projection.work_items.filter(item=>["finance.organization.configure","accounting.book.activate","chart.of.accounts.activate","capital.contribution.post","finance.opening.readiness.evaluate"].includes(item.capability));
 const financeCompleted=financeWorkItems.filter(item=>item.status==="completed").length;
 const financeCurrent=financeWorkItems.find(item=>item.status!=="completed"&&item.status!=="locked");
 const items=building.kind==="office"?[
  {title:"创始人经营桌",status:"可使用",detail:"查看当前经营事件并与数字员工协作。",source:"IAOS Work Item"},
  {title:"证照与印章柜",status:registered?"已解锁":"等待登记",detail:registered?"营业执照和三枚企业印章已经归档。":"登记成功后由 committed legal entity 解锁。",source:"IAOS Legal Entity"},
  {title:"企业资产柜",status:banked?"账户已归档":"持续建设中",detail:"集中展示执照、印章、账户和管理团队。",source:"GameProjection committed facts"},
 ]:building.kind==="government"?[
  {title:"登记综合窗口",status:mission?"正在受理":"空闲",detail:"提交登记资料并接收补正或批准意见。",source:"World Registration Observation"},
  {title:"证照领取台",status:registered?"可以领取":"等待审查",detail:"登记通过后领取虚构营业执照和印章。",source:"IAOS Legal Entity"},
 ]:building.kind==="bank"?[
  {title:"企业金融柜台",status:mission?"尽调进行中":"空闲",detail:"选择开户银行、提交资料并处理补件。",source:"World Bank Observation"},
  {title:"账户资料柜",status:banked?"已解锁":"等待开户",detail:"保存基本账户信息和企业网银 U 盾。",source:"IAOS Bank Account"},
 ]:[
  {title:"组织结构墙",status:organized?"已建立":"等待组织建立",detail:"展示 CEO、CFO 和工厂项目负责人岗位。",source:"IAOS Organization"},
  {title:"治理会议桌",status:mission?"会议待召开":"空闲",detail:"处理任命、经营授权和启动预算。",source:"IAOS Approval + Mandate"},
  {title:"开业财务中心",status:projection.finance_opening?.ready?"账套已启用":financeCurrent?`当前：${financeCurrent.title}`:`财务建设 ${financeCompleted}/5`,detail:projection.finance_opening?.ready?`${projection.finance_opening.book_code} · ${projection.finance_opening.accounting_standard} · ${projection.finance_opening.period_code}。这里是账套、凭证和报表的穿透入口，不在 AESE 复制财务系统。`:`五个 IAOS 财务节点已完成 ${financeCompleted}/5。岗位责任、Mandate、职责分离和升级对象在 IAOS“组织与待办”中查看；当前任务仍从本场景的经营事件推进。`,source:projection.finance_opening?.evidence_ref??"IAOS Finance Work Items",actions:[
   {label:"查看财务组织与待办",href:financeWorkspaceUrl(projection,"operations")},
   ...(projection.finance_opening?.ready?[
    {label:"查看系统账务",href:financeWorkspaceUrl(projection,"ledger")},
    {label:"查看财务报表",href:financeWorkspaceUrl(projection,"reports")},
   ]:[]),
  ]},
 ];
 const inspect=(item:Inspectable,index:number)=>{setPlayerTarget(index+1);playGameTone("inspect");window.setTimeout(()=>setInspected(item),220)};
 const financeItem=items.find(item=>item.title==="开业财务中心");
 return <section className={`gx-location-scene gx-location-${building.kind}`} aria-label={`${copy.title}室内场景`}>
  <header><button onClick={onBack}>返回城市地图</button><div><small>{copy.eyebrow}</small><h2>{copy.title}</h2><p>{copy.description}</p></div><span>{building.state==="active"?"地点已开放":"等待解锁"}</span></header>
  <div className="gx-room">
   {building.kind==="office"&&<><div className="gx-room-window"/><div className="gx-room-desk"><MonitorCog/><span>创始人经营桌</span></div><div className={`gx-room-asset ${registered?"unlocked":""}`}><FileBadge2/><span>{registered?"营业执照已上墙":"证照墙 · 等待登记"}</span></div><div className={`gx-room-asset seal ${registered?"unlocked":""}`}><Stamp/><span>{registered?"企业印章已入柜":"印章柜 · 尚未领取"}</span></div><div className="gx-trophy-shelf"><small>企业资产柜</small><span data-on={registered}><FileBadge2/>执照</span><span data-on={registered}><Stamp/>印章</span><span data-on={banked}><WalletCards/>账户</span><span data-on={organized}><UsersRound/>团队</span></div></>}
   {building.kind==="government"&&<><div className="gx-counter"><Landmark/><strong>企业登记综合窗口</strong><span>受理 · 审查 · 补正 · 发照</span></div><div className="gx-queue"><i/><i/><i/><span>办件等候区</span></div><div className={`gx-certificate-desk ${registered?"unlocked":""}`}><FileBadge2/><strong>{registered?"证照领取台":"材料审查台"}</strong></div></>}
   {building.kind==="bank"&&<><div className="gx-counter bank"><WalletCards/><strong>企业金融服务柜台</strong><span>开户 · 尽调 · 网银 · 资本金</span></div><div className="gx-bank-room"><BriefcaseBusiness/><span>客户经理会谈区</span></div><div className={`gx-ukey-desk ${banked?"unlocked":""}`}><ShieldCheck/><strong>{banked?"账户与 U 盾已就绪":"账户资料领取柜"}</strong></div></>}
   {building.kind==="headquarters"&&<><div className="gx-hq-board"><Building2/><strong>{projection.brand.company_name}</strong><span>{organized?"首届管理团队组建中":"总部空间等待组织建立"}</span></div><div className={`gx-exec-seats ${organized?"unlocked":""}`}><i/><i/><i/><span>CEO · CFO · 工厂项目负责人</span></div><div className="gx-board-table"><UsersRound/><span>治理与经营会议桌</span></div>{financeItem&&<button type="button" className="gx-finance-center" aria-label="打开开业财务中心" onClick={()=>inspect(financeItem,2)}><WalletCards/><strong>开业财务中心</strong><span>{projection.finance_opening?.ready?"账套 · 凭证 · 报表入口":`组织与待办 · ${financeCompleted}/5`}</span></button>}</>}
   <div className={`gx-scene-npc ${mission?"active":""}`} aria-label={`${npc.name} NPC`}><i><span/></i><div><small>{npc.role} · {npcState}</small><strong>{npc.name}</strong><p>{mission?npc.line:"今天这里没有新的事项，我会继续关注企业状态。"}</p></div></div>
   <div className={`gx-room-player target-${playerTarget}`} aria-label="创始人在室内的位置"><i/><span>你</span></div>
   <div className="gx-object-actions" aria-label="可检查场景物件">{items.map((item,index)=><button key={item.title} onClick={()=>inspect(item,index)}><ShieldCheck/>{item.title}</button>)}</div>
   {inspected&&<aside className="gx-object-detail" aria-label={`${inspected.title}详情`}><button aria-label="关闭物件详情" onClick={()=>setInspected(null)}>×</button><small>场景物件</small><h3>{inspected.title}</h3><strong>{inspected.status}</strong><p>{inspected.detail}</p>{inspected.actions&&<nav className="gx-object-detail-actions" aria-label={`${inspected.title}系统入口`}>{inspected.actions.map(action=><a key={action.label} href={action.href} target="_blank" rel="noreferrer">{action.label}</a>)}</nav>}<code>来源：{inspected.source}</code></aside>}
  </div>
  {mission?<MissionBriefing item={mission} onStart={onStart}/>:<div className="gx-location-empty"><ShieldCheck/><div><strong>这里暂时没有需要处理的事件</strong><p>返回城市地图，带任务标记的地点会指引下一步。</p></div></div>}
 </section>;
}
