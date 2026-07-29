import { Bot, Landmark, MessageSquareText, Play, UserRound } from "lucide-react";
import type { GameWorkItem } from "../../game/types";
import "./MissionBriefing.css";

const stories:Record<string,{place:string;npc:string;role:string;story:string;prompt:string}>={
 "founder.resolution.prepare":{place:"创始人办公室",npc:"纪元",role:"企业设立专员",story:"把刚才的创业构想整理成一份能约束公司未来的创始人决议。",prompt:"你希望决议首先保护什么：长期方向、工程品质，还是现金安全？"},
 "founder.resolution.approve":{place:"临时董事会",npc:"治理顾问",role:"董事会秘书",story:"第一份创始人决议已经摆上桌面。只有你能决定是否让它成为公司的治理起点。",prompt:"审阅数字员工的建议，承担创始人的最终责任。"},
 "capital.commitment.record":{place:"创始人办公室",npc:"财务负责人",role:"资金顾问",story:"公司需要第一笔可信的资本承诺。金额将决定我们能走多快，也决定风险有多大。",prompt:"决定你愿意为这家公司承诺多少启动资本。"},
 "registration.package.validate":{place:"法务作战室",npc:"法务合规专员",role:"登记顾问",story:"公司名称、地址和经营范围需要整理成经得起审查的登记材料。",prompt:"让法务员工检查材料，发现问题会回来找你补正。"},
 "registration.submit":{place:"城市政务大厅",npc:"登记窗口专员",role:"外部办事员",story:"登记材料已经准备好。提交后，企业将第一次进入外部规则世界。",prompt:"确认提交，并等待政务世界给出结果。"},
 "registration.observation.commit":{place:"城市政务大厅",npc:"登记窗口专员",role:"外部办事员",story:"窗口送回了企业登记结果。我们需要把外部结果带回企业世界。",prompt:"查看结果并决定下一步。"},
 "bank.account.opening.submit":{place:"商业银行大厅",npc:"企业客户经理",role:"银行客户经理",story:"没有基本账户，公司无法真正接收资本、支付工资或开展交易。",prompt:"向银行说明公司的业务与资金来源。"},
 "bank.account.observation.commit":{place:"商业银行大厅",npc:"企业客户经理",role:"银行客户经理",story:"银行完成了尽调，开户决定已经送达。",prompt:"接受结果，或根据拒绝原因准备补充材料。"},
 "capital.contribution.verify":{place:"财务控制室",npc:"财务负责人",role:"资金守门人",story:"承诺不是现金。现在要核对真正到账的资金，任何差异都会影响开业。",prompt:"确认实缴金额与资金凭证。"},
 "organization.establish":{place:"组织设计室",npc:"组织顾问",role:"治理组织专员",story:"公司不能只有创始人。我们需要建立最小但清晰的组织骨架。",prompt:"确定岗位如何分工、谁对结果负责。"},
 "executive.appointment.propose":{place:"人才会客室",npc:"CEO 候选人",role:"职业经理人",story:"一位 CEO 候选人来到办公室。他想知道自己将获得多大授权，也想知道你会保留什么权力。",prompt:"提出任命条件与经营期待。"},
 "executive.appointment.acceptance.commit":{place:"人才会客室",npc:"CEO 候选人",role:"职业经理人",story:"候选人已经回复任命邀请。一次任命也是一次双向选择。",prompt:"查看接受、拒绝或附加条件。"},
 "executive.appointment.approve":{place:"临时董事会",npc:"董事会秘书",role:"治理顾问",story:"CEO 任命进入最终表决。你要判断这个人是否配得上公司的第一份经营授权。",prompt:"作出最终任命决定。"},
 "operating.mandate.grant":{place:"董事长办公室",npc:"新任 CEO",role:"经营负责人",story:"职位并不等于权力。现在要明确 CEO 可以决定什么、不能决定什么。",prompt:"授予一份边界清晰、可撤回的经营 Mandate。"},
 "initial.budget.prepare":{place:"经营会议室",npc:"财务负责人",role:"预算规划者",story:"开业前的每一元钱都要有去处。数字员工正在形成第一份经营预算。",prompt:"选择增长速度与现金安全之间的平衡。"},
 "initial.budget.approve":{place:"经营会议室",npc:"新任 CEO",role:"经营负责人",story:"预算代表公司的真实优先级。批准后，团队就会据此行动。",prompt:"批准、缩减或要求重新编制。"},
 "enterprise.readiness.evaluate":{place:"开业指挥中心",npc:"独立审计专员",role:"开业评估官",story:"执照、账户、资本、团队和预算都已准备。最后一次检查决定企业能否正式开业。",prompt:"让审计员工检查所有证据，迎接开业时刻。"},
};

export function MissionBriefing({item,onStart}:{item:GameWorkItem;onStart:()=>void}){
 const scene=stories[item.capability]??{place:"企业指挥中心",npc:"数字员工",role:"任务顾问",story:item.title,prompt:"查看当前情况并作出经营决定。"};
 return <section className="gx-mission" aria-label="当前剧情任务">
  <div className="gx-mission__location"><Landmark aria-hidden="true"/><span><small>当前地点</small><strong>{scene.place}</strong></span></div>
  <div className="gx-mission__npc"><span>{item.requires_me?<UserRound aria-hidden="true"/>:<Bot aria-hidden="true"/>}</span><div><small>{scene.role}</small><strong>{scene.npc}</strong></div></div>
  <div className="gx-mission__story"><MessageSquareText aria-hidden="true"/><div><strong>{scene.story}</strong><p>{scene.prompt}</p></div></div>
  <button onClick={onStart}><Play aria-hidden="true"/>{item.requires_me?"作出决定":"推进剧情"}</button>
 </section>;
}
