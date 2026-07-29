import{Bot,Building2,CheckCircle2,FileCheck2,KeyRound,Landmark,LoaderCircle,Play,RefreshCw,ShieldCheck,Stamp,UserRound,X}from"lucide-react";
import{useMemo,useState}from"react";
import type{ReactNode}from"react";
import{approveAndExecuteWorkItem,dispatchAgentWorkItem,executeWorkItem,resolveWorldWorkItem,type WorkItemInput}from"../../game/api";
import type{GameProjection,GameWorkItem}from"../../game/types";
import{RPGEventIntro}from"./RPGEventIntro";

const labels:Record<string,{title:string;description:string;action:string}>={
 "founder.resolution.prepare":{title:"让企业设立专员起草创始决议",description:"输入你的经营目标，数字员工将形成可审计草案并提交 IAOS。",action:"派遣数字员工"},
 "founder.resolution.approve":{title:"批准创始人设立决议",description:"这是正式治理动作。批准后企业才能登记资本承诺。",action:"批准并提交"},
 "capital.commitment.record":{title:"登记创始资本承诺",description:"输入计划投入公司的资本，IAOS 将其记录为认缴，不会冒充实际到账。",action:"让财务员工登记"},
 "registration.package.validate":{title:"校验公司登记材料",description:"法务数字员工检查企业名称、地址和经营范围的一致性。",action:"开始材料校验"},
 "registration.submit":{title:"提交公司登记申请",description:"你授权向模拟登记机构发出正式申请；外部结果仍需等待。",action:"批准并提交登记"},
 "registration.observation.commit":{title:"等待登记机构办理",description:"推进外部世界，模拟登记机构审查并把可信结果回传 IAOS。",action:"前往政务中心办理"},
 "bank.account.opening.submit":{title:"申请企业银行账户",description:"你授权向模拟银行提交开户申请。",action:"批准开户申请"},
 "bank.account.observation.commit":{title:"等待银行开户",description:"推进银行办理过程，开户结果会以 World Observation 回传。",action:"前往银行办理"},
 "capital.contribution.verify":{title:"注入并核验实缴资本",description:"输入实际到账金额。IAOS 会校验它与资本承诺是否一致。",action:"批准资金核验"},
 "organization.establish":{title:"建立初始组织",description:"让 IAOS 创建 CEO、CFO 与工厂项目负责人岗位。",action:"建立组织"},
 "executive.appointment.propose":{title:"提名管理团队",description:"治理数字员工准备三项任命，并向候选人发出接受请求。",action:"生成任命方案"},
 "executive.appointment.acceptance.commit":{title:"候选人接受任命",description:"推进人才世界，模拟候选人逐一接受岗位。",action:"组织任命会谈"},
 "executive.appointment.approve":{title:"批准管理层任命",description:"对已接受的 CEO、CFO 与项目负责人任命进行正式审批。",action:"批准任命"},
 "operating.mandate.grant":{title:"授予经营授权",description:"向管理层授予有范围、限额和有效期的经营权限。",action:"批准经营授权"},
 "initial.budget.prepare":{title:"编制首年启动预算",description:"让财务员工编制预算。预算是授权上限，不会直接消耗现金。",action:"让财务员工编制"},
 "initial.budget.approve":{title:"批准首年启动预算",description:"预算金额不得超过已核验公司现金。",action:"批准预算"},
 "enterprise.readiness.evaluate":{title:"执行企业开业检查",description:"审计数字员工核验法人、账户、资本、组织、授权和预算。",action:"开始开业检查"}
};
const needsAmount=(capability:string)=>["capital.commitment.record","capital.contribution.verify","initial.budget.prepare","initial.budget.approve"].includes(capability);
const registrationDocs=["公司设立登记申请书","创始人决议","法定代表人任职文件","注册地址使用证明","经营范围说明"];
const bankDocs=["营业执照","法定代表人身份证明","公司章程与创始人决议","公章、财务章、法人章印鉴","注册地址与经营场所证明"];
const banks=[
 {id:"river",name:"江海商业银行",speed:"1–2 个工作日",focus:"制造业企业服务",risk:"经营范围表述不清时可能要求补充合同"},
 {id:"city",name:"东吴城市银行",speed:"即时预审",focus:"本地创业企业",risk:"注册地址证明不足会拒绝开户"},
 {id:"industry",name:"华业产业银行",speed:"2–3 个工作日",focus:"供应链与跨境结算",risk:"资金来源说明不完整会进入强化尽调"}
];

export function WorkItemActionPanel({item,projection,onClose,onDone}:{item:GameWorkItem;projection:GameProjection;onClose:()=>void;onDone:()=>void}){
 const meta=labels[item.capability]??{title:item.title,description:"通过 IAOS 受治理能力执行当前工作项。",action:"执行当前任务"};
 const suggested=useMemo(()=>projection.resources.capital_committed.value!=="0"?String(Number(projection.resources.capital_committed.value)/100):"1000000",[projection.resources.capital_committed.value]);
 const[amount,setAmount]=useState(suggested),[objective,setObjective]=useState("批准启动企业设立工作，建立面向工业客户的长期经营主体"),[proposals,setProposals]=useState("确认拟设公司身份、注册地址、经营范围、初始资本与治理结构，并授权后续登记办理"),[risk,setRisk]=useState("登记与银行结果必须等待可信外部回传；预算与现金严格分离"),[busy,setBusy]=useState(false),[error,setError]=useState("");
 const bankStorageKey=`gx-bank-${projection.case_code}`;
 const[selectedBank,setSelectedBank]=useState(()=>localStorage.getItem(bankStorageKey)??"river"),[checkedDocs,setCheckedDocs]=useState<string[]>(item.capability.includes("bank")?bankDocs:registrationDocs),[externalStage,setExternalStage]=useState<"prepare"|"rejected"|"approved">("prepare"),[feedback,setFeedback]=useState("");
 const chooseBank=(bank:string)=>{setSelectedBank(bank);localStorage.setItem(bankStorageKey,bank)};
 const toggleDoc=(doc:string)=>setCheckedDocs(current=>current.includes(doc)?current.filter(value=>value!==doc):[...current,doc]);
 const committedYuan=Number(projection.resources.capital_committed.value)/100;
 const enteredYuan=Number(amount);
 const contributionMismatch=item.capability==="capital.contribution.verify"&&Number.isFinite(enteredYuan)&&enteredYuan!==committedYuan;
 const submit=async()=>{
  setBusy(true);setError("");
  const input:WorkItemInput={};
  if(item.kind==="agent_task")input.business_note=`玩家在 Enterprise Genesis 中派遣数字员工执行 ${item.capability}`;
  if(needsAmount(item.capability)){const value=Number(amount);if(!Number.isFinite(value)||value<=0){setError("请输入大于 0 的金额");setBusy(false);return}input.amount_minor=Math.round(value*100);input.currency="CNY"}
  if(item.capability==="capital.contribution.verify"&&input.amount_minor!==Number(projection.resources.capital_committed.value)){
   setError(`本次到账金额必须与认缴资本一致：¥${committedYuan.toLocaleString("zh-CN",{minimumFractionDigits:2})}`);setBusy(false);return;
  }
  if(item.capability==="founder.resolution.prepare"){input.resolution_objective=objective;input.key_proposals=proposals;input.risk_notes=risk}
  if(item.capability==="registration.submit")input.business_note=`提交登记资料：${checkedDocs.join("、")}`;
  if(item.capability==="bank.account.opening.submit")input.business_note=`申请开户银行：${banks.find(bank=>bank.id===selectedBank)?.name}；资料：${checkedDocs.join("、")}`;
  try{
   const sequence=Number(item.work_item_id.replace(/\D/g,""));
   if(item.kind==="world_wait"){
    const required=item.capability.includes("bank")?bankDocs:registrationDocs;
    const missing=required.filter(doc=>!checkedDocs.includes(doc));
    if(missing.length){
     try{await resolveWorldWorkItem(projection.case_code,sequence,item.capability,projection.world_run_id,"rejected")}catch{/* IAOS 正确拒绝推进；rejected observation 已保留为外部反馈证据。 */}
     setExternalStage("rejected");setFeedback(`${item.capability.includes("bank")?"银行尽调未通过":"登记材料不合格"}：缺少${missing.join("、")}。请补齐后重新申请。`);setBusy(false);return;
    }
    const result=item.capability.includes("appointment")?"accepted":item.capability.includes("bank")?"opened":"registered";
    await resolveWorldWorkItem(projection.case_code,sequence,item.capability,projection.world_run_id,result);
    setExternalStage("approved");
    setFeedback(item.capability.includes("bank")?"银行尽调通过，基本账户已开立。请领取账户资料与企业网银 U 盾。":"登记审查通过，企业法律主体成立。请领取营业执照与印章套装。");
    return;
   }else if(item.kind==="agent_task"){
    await dispatchAgentWorkItem(projection.case_code,sequence,input);
   }else if(item.kind==="approval"){
    if(!item.gate)throw new Error("当前审批任务缺少 IAOS Gate 合同");
    await approveAndExecuteWorkItem(projection.case_code,sequence,item.capability,item.gate,input);
   }else await executeWorkItem(projection.case_code,sequence,input);
   onDone();
  }catch(e){setError(e instanceof Error?e.message:String(e))}finally{setBusy(false)}
 };
 return <section className="gx-action-card" aria-labelledby="gx-action-title">
  <div className="gx-action-head"><span>{item.kind==="agent_task"?<Bot/>:item.kind==="world_wait"?<Landmark/>:item.kind==="approval"?<ShieldCheck/>:<UserRound/>}</span><div><small>{item.kind==="agent_task"?"数字员工任务":item.kind==="world_wait"?"外部世界办理":item.kind==="approval"?`治理审批 ${item.gate??""}`:"IAOS 系统任务"}</small><h3 id="gx-action-title">{meta.title}</h3></div><button onClick={onClose} aria-label="关闭任务操作"><X/></button></div>
  <RPGEventIntro item={item}/>
  <p>{meta.description}</p>
  {item.capability==="registration.submit"&&<ApplicationPackage title="本次提交的企业设立登记申请" subtitle="批准后，这一资料包将提交到模拟登记机构" docs={registrationDocs} checked={checkedDocs} onToggle={toggleDoc}>
   <dl><div><dt>申请企业</dt><dd>{projection.brand.company_name}</dd></div><div><dt>申请类型</dt><dd>有限责任公司设立登记</dd></div><div><dt>审查重点</dt><dd>名称、地址、经营范围、治理决议与资本承诺一致性</dd></div></dl>
  </ApplicationPackage>}
  {item.capability==="registration.observation.commit"&&<ExternalReview kind="registration" stage={externalStage} feedback={feedback} docs={registrationDocs} checked={checkedDocs} onToggle={toggleDoc} company={projection.brand.company_name??"拟设企业"}/>}
  {item.capability==="bank.account.opening.submit"&&<div className="gx-application-flow"><h4><Building2/>选择开户银行</h4><div className="gx-bank-grid">{banks.map(bank=><button type="button" className={selectedBank===bank.id?"selected":""} onClick={()=>chooseBank(bank.id)} key={bank.id}><strong>{bank.name}</strong><span>{bank.focus} · {bank.speed}</span><small>可能补件：{bank.risk}</small></button>)}</div><ApplicationPackage title="企业基本账户开户申请" subtitle="选择银行后核对完整资料包" docs={bankDocs} checked={checkedDocs} onToggle={toggleDoc}><dl><div><dt>开户主体</dt><dd>{projection.brand.company_name}</dd></div><div><dt>开户银行</dt><dd>{banks.find(bank=>bank.id===selectedBank)?.name}</dd></div><div><dt>账户用途</dt><dd>资本金、经营收支、工资与税费结算</dd></div></dl></ApplicationPackage></div>}
  {item.capability==="bank.account.observation.commit"&&<ExternalReview kind="bank" stage={externalStage} feedback={feedback} docs={bankDocs} checked={checkedDocs} onToggle={toggleDoc} company={projection.brand.company_name??"企业"} bankName={banks.find(bank=>bank.id===selectedBank)?.name}/>}
  {item.capability==="founder.resolution.prepare"&&<div className="gx-output-explainer">
   <strong>本步骤会形成什么？</strong>
   <span>你输入的三项内容将由企业设立专员整理为《创始人设立决议草案》，持久化到 IAOS 治理决议记录，并送到下一步 G1 供你逐条审阅。</span>
  </div>}
  {item.kind==="approval"&&item.approval_review&&<article className="gx-review-sheet" aria-label="待审批内容">
   <header><div><small>待审议文件 · {item.gate}</small><h4>{item.approval_review.title}</h4></div><span>{item.approval_review.status}</span></header>
   <p>{item.approval_review.summary}</p>
   <dl>
    <div><dt>起草 / 提交</dt><dd>{item.approval_review.prepared_by}</dd></div>
    {item.approval_review.fields.filter(field=>field.value).map(field=><div key={field.label}><dt>{field.label}</dt><dd>{field.value}</dd></div>)}
   </dl>
   {item.approval_review.risks.length>0&&<section><strong>风险与限制</strong>{item.approval_review.risks.map(text=><p key={text}>{text}</p>)}</section>}
   <footer><strong>批准后的效果</strong><span>{item.approval_review.approval_effect}</span><small>证据：{item.approval_review.evidence_ref}</small></footer>
  </article>}
  {item.kind==="approval"&&!item.approval_review&&<p className="gx-action-error" role="alert">IAOS 尚未返回可审阅的审批对象，已禁止盲目批准。请刷新后重试。</p>}
  {item.capability==="founder.resolution.prepare"&&<div className="gx-action-fields">
   <label>决议目标<textarea value={objective} onChange={e=>setObjective(e.target.value)}/></label>
   <label>核心提案<textarea value={proposals} onChange={e=>setProposals(e.target.value)}/></label>
   <label>风险与限制<textarea value={risk} onChange={e=>setRisk(e.target.value)}/></label>
  </div>}
  {needsAmount(item.capability)&&<label className="gx-money-input">金额（人民币元）<div><span>¥</span><input type="number" min="1" step="1000" value={amount} onChange={e=>setAmount(e.target.value)}/></div><small>提交时按分精确写入 IAOS；认缴、实缴、预算分别管理。</small></label>}
  {item.capability==="capital.contribution.verify"&&<section className={`gx-capital-reconcile ${contributionMismatch?"mismatch":"matched"}`} aria-live="polite">
   <header><strong>资金到账核对</strong><span>{contributionMismatch?"存在差额，不能核验":"金额一致，可以核验"}</span></header>
   <dl><div><dt>认缴资本</dt><dd>¥{committedYuan.toLocaleString("zh-CN",{minimumFractionDigits:2})}</dd></div><div><dt>本次到账</dt><dd>¥{(Number.isFinite(enteredYuan)?enteredYuan:0).toLocaleString("zh-CN",{minimumFractionDigits:2})}</dd></div><div><dt>差额</dt><dd>¥{(Number.isFinite(enteredYuan)?enteredYuan-committedYuan:-committedYuan).toLocaleString("zh-CN",{minimumFractionDigits:2})}</dd></div></dl>
   {contributionMismatch&&<button type="button" onClick={()=>setAmount(String(committedYuan))}><RefreshCw/>按认缴金额修正为 ¥{committedYuan.toLocaleString("zh-CN")}</button>}
  </section>}
  <div className="gx-governance-note"><CheckCircle2/><span>操作将写入 IAOS 业务记录、Process Work Item、审计日志与 Outbox，可重复验证。</span></div>
  {error&&<p className="gx-action-error" role="alert">{error}</p>}
  {externalStage==="approved"?<button className="gx-action-submit" onClick={onDone}><CheckCircle2/>领取资产并继续</button>:<button className="gx-action-submit" disabled={busy||contributionMismatch||(item.kind==="approval"&&!item.approval_review)} onClick={()=>void submit()}>{busy?<LoaderCircle className="gx-spin"/>:externalStage==="rejected"?<RefreshCw/>:<Play/>}{busy?"正在通过 IAOS 执行…":externalStage==="rejected"?"已补正，重新申请":item.kind==="approval"?"已审阅，批准并执行":meta.action}</button>}
 </section>
}

function ApplicationPackage({title,subtitle,docs,checked,onToggle,children}:{title:string;subtitle:string;docs:string[];checked:string[];onToggle:(doc:string)=>void;children:ReactNode}){
 return <section className="gx-application-package"><header><FileCheck2/><div><h4>{title}</h4><p>{subtitle}</p></div><span>{checked.length}/{docs.length}</span></header>{children}<fieldset><legend>提交资料清单</legend>{docs.map(doc=><label key={doc}><input type="checkbox" checked={checked.includes(doc)} onChange={()=>onToggle(doc)}/><span>{doc}</span></label>)}</fieldset></section>
}
function ExternalReview({kind,stage,feedback,docs,checked,onToggle,company,bankName}:{kind:"registration"|"bank";stage:"prepare"|"rejected"|"approved";feedback:string;docs:string[];checked:string[];onToggle:(doc:string)=>void;company:string;bankName?:string}){
 const bank=kind==="bank";
 return <section className={`gx-external-review ${stage}`}><header><Landmark/><div><small>{bank?"银行尽职调查":"登记机构形式与实质审查"}</small><h4>{stage==="approved"?"申请已通过":stage==="rejected"?"申请退回补正":"模拟外部机构正在受理"}</h4></div></header>
  {stage!=="approved"&&<><p>{stage==="rejected"?feedback:`窗口将检查 ${docs.length} 项资料。你可以故意取消一项，体验退回、补正和重新申请。`}</p><fieldset>{docs.map(doc=><label key={doc}><input type="checkbox" checked={checked.includes(doc)} onChange={()=>onToggle(doc)}/>{doc}</label>)}</fieldset></>}
  {stage==="approved"&&(bank?<BankAsset company={company} bankName={bankName??"虚拟商业银行"}/>:<RegistrationAssets company={company}/>)}
 </section>
}
function RegistrationAssets({company}:{company:string}){
 return <div className="gx-reward-assets"><article className="gx-license"><img src="/assets/enterprise-genesis/sprites/license-v1.png" alt="虚构营业执照收藏品"/><small>企业创生沙盘 · 虚构证照</small><Stamp/><h5>营业执照</h5><strong>{company}</strong><span>统一代码：GX · {Math.abs(company.length*7919)} · SIM</span><p>状态：存续（沙盘）　登记机关：企业创生虚拟登记中心</p></article><article className="gx-seal-kit"><img src="/assets/enterprise-genesis/sprites/seal-kit-v1.png" alt="虚构企业印章套装"/><small>虚构印章套装</small><div><span><Stamp/>企业公章</span><span><Stamp/>财务专用章</span><span><Stamp/>法定代表人章</span></div></article></div>
}
function BankAsset({company,bankName}:{company:string;bankName:string}){
 return <div className="gx-reward-assets"><article className="gx-account-card"><small>企业创生沙盘 · 虚构账户</small><Building2/><h5>基本存款账户信息</h5><strong>{company}</strong><span>{bankName}</span><span>账号：GX88 · 0001 · {String(company.length*104729).padStart(8,"0")}</span><p>账户状态：正常　网银权限：待首次登录激活</p></article><article className="gx-ukey"><img src="/assets/enterprise-genesis/sprites/ukey-v1.png" alt="虚构企业网银 U 盾"/><KeyRound/><div><strong>企业网银 U 盾</strong><span>设备编号 GX-UKEY-{company.length}A</span><small>虚构数字资产 · 不可用于真实银行系统</small></div></article></div>
}
