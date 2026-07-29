import { Bot, Landmark, MessageCircleMore, ShieldCheck, UserRound } from "lucide-react";
import { useState } from "react";
import type { GameWorkItem } from "../../game/types";
import "./RPGEventIntro.css";

type Scene = { chapter:string; location:string; npc:string; role:string; opening:string; choices:string[]; consequence:string };

const scenes:Record<string,Scene>={
 "founder.resolution.prepare":{chapter:"创始办公室",location:"经营桌",npc:"纪元",role:"企业设立专员",opening:"你的创业构想需要变成第一份能约束未来管理层的文件。",choices:["优先守住长期方向","优先守住工程品质","优先守住现金安全"],consequence:"选择会作为创始决议的叙事重点，正式内容仍由下方文件确认。"},
 "founder.resolution.approve":{chapter:"临时董事会",location:"治理会议桌",npc:"纪元",role:"董事会秘书",opening:"决议草案已经送达。签字意味着你愿意承担这家企业的最终责任。",choices:["逐条审阅风险","核对授权边界","确认后签署"],consequence:"审批对象与证据会永久进入 IAOS 治理档案。"},
 "capital.commitment.record":{chapter:"资本启程",location:"创始办公室",npc:"周衡",role:"资金顾问",opening:"承诺多少资本，决定公司能跑多快，也决定你的风险暴露。",choices:["稳健起步","平衡投入","快速扩张"],consequence:"经营取向只帮助思考，认缴金额以下方输入为准。"},
 "registration.package.validate":{chapter:"公司登记",location:"法务作战室",npc:"林岚",role:"登记顾问",opening:"名称、地址、经营范围和治理文件必须讲述同一个企业。",choices:["先查主体一致性","先查缺失资料","先查高风险表述"],consequence:"数字员工将形成可审计的材料校验结果。"},
 "registration.submit":{chapter:"公司登记",location:"政务服务中心",npc:"林岚",role:"登记窗口专员",opening:"材料已经装订完毕。提交以后，企业第一次进入外部规则世界。",choices:["核对申请主体","核对五项附件","正式递交"],consequence:"外部机构可能批准，也可能给出补正意见。"},
 "registration.observation.commit":{chapter:"公司登记",location:"登记综合窗口",npc:"林岚",role:"登记窗口专员",opening:"审查结果回来了。先听完反馈，再决定领取证照还是补充材料。",choices:["查看审查意见","补齐缺失材料","领取企业身份"],consequence:"可信 World Observation 才能推进 IAOS 状态。"},
 "bank.account.opening.submit":{chapter:"银行资本",location:"合作银行",npc:"周衡",role:"企业客户经理",opening:"银行想了解企业如何赚钱、资金从哪里来，以及谁能动用账户。",choices:["比较银行方案","说明业务模式","核对开户资料"],consequence:"你的选择形成开户意向，正式申请以下方资料包为准。"},
 "bank.account.observation.commit":{chapter:"银行资本",location:"企业金融柜台",npc:"周衡",role:"企业客户经理",opening:"尽调已经结束。银行可能开户，也可能要求补充经营与地址证据。",choices:["听取尽调结论","提交补充证据","领取账户与 U 盾"],consequence:"开户结果由外部观察回传，不由玩家直接改写。"},
 "capital.contribution.verify":{chapter:"银行资本",location:"资金核验室",npc:"周衡",role:"资金守门人",opening:"承诺不是现金。现在必须逐分核对真正到账的资本。",choices:["核对银行回单","核对认缴金额","确认资金来源"],consequence:"金额不一致会被 IAOS 拒绝，不能用剧情选择绕过。"},
 "organization.establish":{chapter:"人才治理",location:"企业总部",npc:"顾远",role:"组织设计顾问",opening:"公司不能只靠创始人。我们要先建立最小但清晰的责任骨架。",choices:["先建立经营岗位","先明确汇报关系","先冻结职责边界"],consequence:"IAOS 将创建岗位，不会凭空任命人员。"},
 "executive.appointment.propose":{chapter:"人才治理",location:"人才会客室",npc:"顾远",role:"任命顾问",opening:"三位候选人正在等你的条件。职位、授权和责任必须分开谈。",choices:["讨论 CEO 条件","讨论财务守门职责","讨论工厂交付责任"],consequence:"数字员工会形成任命方案并发出邀请。"},
 "executive.appointment.acceptance.commit":{chapter:"人才治理",location:"候选人会谈区",npc:"顾远",role:"人才谈判主持人",opening:"候选人带着自己的条件回来。任命从来是双向选择。",choices:["听取接受条件","回应权责疑问","确认加入意向"],consequence:"候选人的外部反馈必须先被可信观察记录。"},
 "executive.appointment.approve":{chapter:"人才治理",location:"治理会议桌",npc:"顾远",role:"董事会秘书",opening:"候选人已经接受邀请，现在轮到创始人完成最终任命。",choices:["核对岗位匹配","核对利益冲突","表决任命"],consequence:"批准后，任命才成为 IAOS committed 事实。"},
 "operating.mandate.grant":{chapter:"开业准备",location:"董事长办公室",npc:"顾远",role:"治理组织顾问",opening:"职位不等于权力。每一项经营权限都必须有范围、限额和期限。",choices:["保留重大事项","授权日常经营","设置撤回条件"],consequence:"授权书将成为 Agent 和管理层行动的治理边界。"},
 "initial.budget.prepare":{chapter:"开业准备",location:"经营会议室",npc:"周衡",role:"预算规划者",opening:"团队的每一项雄心，都要落到现金和预算约束上。",choices:["优先产品研发","优先市场开拓","优先现金安全"],consequence:"财务员工据此编制方案，金额仍需单独确认。"},
 "initial.budget.approve":{chapter:"开业准备",location:"预算会议桌",npc:"周衡",role:"财务负责人",opening:"预算代表真正的优先级。批准以后，团队会据此申请资源。",choices:["审阅资金上限","审阅部门分配","批准启动预算"],consequence:"预算只是授权上限，不会自动消耗公司现金。"},
 "enterprise.readiness.evaluate":{chapter:"开业准备",location:"开业指挥中心",npc:"纪元",role:"独立开业评估官",opening:"执照、账户、资本、团队、授权和预算都必须经得起最后检查。",choices:["检查法律身份","检查经营能力","发起开业评估"],consequence:"只有完整证据链通过，企业才会进入 operational ready。"},
};

export function RPGEventIntro({item}:{item:GameWorkItem}){
 const scene=scenes[item.capability]??{chapter:"企业事件",location:"企业总部",npc:"数字员工",role:"任务顾问",opening:item.title,choices:["了解情况","查看证据","继续处理"],consequence:"正式结果由 IAOS committed state 决定。"};
 const[selected,setSelected]=useState(0);
 const Actor=item.kind==="agent_task"?Bot:item.kind==="world_wait"?Landmark:item.kind==="approval"?ShieldCheck:UserRound;
 return <section className="gx-rpg-event" aria-label={`${scene.chapter}剧情`}>
  <header><span>{scene.chapter}</span><strong>{scene.location}</strong></header>
  <div className="gx-rpg-event__portrait"><Actor/><small>{scene.role}</small><strong>{scene.npc}</strong></div>
  <div className="gx-rpg-event__dialogue"><MessageCircleMore/><p>{scene.opening}</p><div>{scene.choices.map((choice,index)=><button type="button" aria-pressed={selected===index} onClick={()=>setSelected(index)} key={choice}><span>{index+1}</span>{choice}</button>)}</div><small>{scene.consequence}</small></div>
 </section>;
}
