import {
  ArrowLeft,
  BadgeCheck,
  Banknote,
  Bot,
  Building2,
  CircleHelp,
  Factory,
  FilePlus2,
  LoaderCircle,
  Pause,
  Play,
  RotateCcw,
  SearchCheck,
  StepForward,
  TriangleAlert,
  XCircle,
} from "lucide-react";
import { type FormEvent, useEffect, useMemo, useState } from "react";
import {
  generatePlantProposals,
  loadPlantBuild,
  loadPlantFinancialConstraint,
  loadPlantPlanningStatus,
  loadPlantRequirement,
  submitPlantProposalReview,
  type FacilityRequirement,
  type PlantBuildTrace,
  type PlantPlanningProviderStatus,
  type PlantFinancialConstraint,
  type ProposalSet,
  type SiteOptionProposal,
} from "../../world/plantBuild";
import "./WorldPlay.css";

const cny = (value: string) =>
  new Intl.NumberFormat("zh-CN", {
    style: "currency",
    currency: "CNY",
    maximumFractionDigits: 0,
  }).format(Number(value));

const OPTION_TYPES = [
  ["lease_and_retrofit", "租赁并改造"],
  ["greenfield_build", "购地自建"],
  ["build_to_suit", "定制建设/长期租赁"],
  ["existing_plant_purchase", "收购现有厂房"],
] as const;

type RequirementDraft = {
  targetRegion: string;
  facilityPurpose: string;
  minimumAreaM2: string;
  minimumElectricKVA: string;
  targetAvailableAt: string;
  candidateCount: string;
  optionTypes: string[];
  investmentRequest: string;
  minimumCashReserve: string;
  preferences: string;
  revisionReason: string;
};

type ReviewState = { action: string; reason: string };

const blankDraft: RequirementDraft = {
  targetRegion: "",
  facilityPurpose: "",
  minimumAreaM2: "",
  minimumElectricKVA: "",
  targetAvailableAt: "",
  candidateCount: "3",
  optionTypes: [],
  investmentRequest: "",
  minimumCashReserve: "",
  preferences: "",
  revisionReason: "首次设施规划",
};

const localDate = (value: string) =>
  value ? new Date(value).toISOString() : "";

function ProposalList({ title, values }: { title: string; values: string[] }) {
  return (
    <div>
      <strong>{title}</strong>
      {values.length ? (
        <ul>{values.map((value) => <li key={value}>{value}</li>)}</ul>
      ) : (
        <p>尚未填写</p>
      )}
    </div>
  );
}

function ProposalCard({
  proposal,
  review,
  onReview,
  onSubmitReview,
  saved,
  busy,
}: {
  proposal: SiteOptionProposal;
  review?: ReviewState;
  onReview: (value: ReviewState) => void;
  onSubmitReview: () => void;
  saved: boolean;
  busy: boolean;
}) {
  const action = review?.action ?? "";
  const formal = proposal.status !== "manual_draft";
  const label = OPTION_TYPES.find(([code]) => code === proposal.option_type)?.[1] ?? proposal.option_type;
  return (
    <article className="plant-proposal-card">
      <header>
        <div>
          <span>{label}</span>
          <h3>{proposal.display_name}</h3>
        </div>
        <strong>{Math.round(Number(proposal.confidence) * 100)}% 置信度</strong>
      </header>
      <p>{proposal.business_rationale}</p>
      <dl>
        <div><dt>投资估算</dt><dd>{cny(proposal.estimated_amount.minimum.value)} – {cny(proposal.estimated_amount.maximum.value)}<small>可能值 {cny(proposal.estimated_amount.likely.value)}</small></dd></div>
        <div><dt>估算依据</dt><dd>{proposal.estimated_amount.basis}</dd></div>
        <div><dt>预计可用</dt><dd>{new Date(proposal.estimated_schedule.likely).toLocaleDateString("zh-CN")}<small>{new Date(proposal.estimated_schedule.earliest).toLocaleDateString("zh-CN")} – {new Date(proposal.estimated_schedule.latest).toLocaleDateString("zh-CN")}</small></dd></div>
      </dl>
      <div className="plant-proposal-evidence">
        <ProposalList title="假设" values={proposal.assumptions} />
        <ProposalList title="待验证事实" values={proposal.facts_required} />
        <ProposalList title="风险" values={proposal.risks} />
        <ProposalList title="来源" values={proposal.source_refs} />
      </div>
      <fieldset className="plant-review-actions">
        <legend>人工审阅决定</legend>
        <button type="button" className={action === "adopt_for_investigation" ? "selected" : ""} onClick={() => onReview({ action: "adopt_for_investigation", reason: review?.reason ?? "" })}><SearchCheck />采纳调研</button>
        <button type="button" className={action === "request_revision" ? "selected" : ""} onClick={() => onReview({ action: "request_revision", reason: review?.reason ?? "" })}><RotateCcw />退回重生成</button>
        <button type="button" className={action === "discard" ? "selected danger" : ""} onClick={() => onReview({ action: "discard", reason: review?.reason ?? "" })}><XCircle />淘汰</button>
        <label>
          审阅理由
          <textarea minLength={6} value={review?.reason ?? ""} onChange={(event) => onReview({ action, reason: event.target.value })} placeholder="说明采纳、退回或淘汰的业务理由（至少 6 个字符）" />
        </label>
        <button type="button" className="plant-review-submit" disabled={!formal || !action || (review?.reason.trim().length ?? 0) < 6 || busy || saved} onClick={onSubmitReview}>
          {busy ? <LoaderCircle className="gx-spin" /> : formal ? <BadgeCheck /> : <TriangleAlert />}{!formal ? "人工候选仍是本地草稿" : saved ? "已保存审阅" : busy ? "正在保存…" : "提交审阅到 IAOS"}
        </button>
      </fieldset>
    </article>
  );
}

function PlantPlanningWorkspace() {
  const [status, setStatus] = useState<PlantPlanningProviderStatus | null>(null);
  const [financial, setFinancial] = useState<PlantFinancialConstraint | null>(null);
  const [draft, setDraft] = useState<RequirementDraft>(blankDraft);
  const [proposalSet, setProposalSet] = useState<ProposalSet | null>(null);
  const [reviews, setReviews] = useState<Record<string, ReviewState>>({});
  const [savedReviews, setSavedReviews] = useState<Record<string, boolean>>({});
  const [reviewBusy, setReviewBusy] = useState("");
  const [nextRevision, setNextRevision] = useState(1);
  const [manualOpen, setManualOpen] = useState(false);
  const [manualName, setManualName] = useState("");
  const [manualRationale, setManualRationale] = useState("");
  const [manualAmount, setManualAmount] = useState("");
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState("");
  const routeParams = new URLSearchParams(window.location.hash.split("?")[1] ?? "");
  const tenant = routeParams.get("tenant") ?? localStorage.getItem("aese_iaos_tenant_id") ?? localStorage.getItem("iaos_tenant_id") ?? "";
  const caseCode = routeParams.get("case") ?? localStorage.getItem("aese_genesis_case_code") ?? "";
  const requirementID = `facility-requirement-${caseCode || "draft"}`;
  useEffect(() => {
    const controller = new AbortController();
    Promise.all([
      loadPlantPlanningStatus(controller.signal).then(setStatus),
      loadPlantFinancialConstraint(caseCode, controller.signal).then(setFinancial),
      loadPlantRequirement(requirementID, controller.signal).then((existing) => setNextRevision((existing?.revision ?? 0) + 1)),
    ]).catch((reason) => { if (reason.name !== "AbortError") setError(String(reason)); });
    return () => controller.abort();
  }, [caseCode, requirementID]);
  const update = (name: keyof RequirementDraft, value: string | string[]) =>
    setDraft((current) => ({ ...current, [name]: value }));
  const entityCode = financial?.legal_entity_code ?? "";

  const submit = async (event: FormEvent) => {
    event.preventDefault();
    setBusy(true);
    setError("");
    try {
      if (!financial) throw new Error("尚未取得 IAOS 权威资金与预算快照");
      const money = (value: string) => ({ value, currency: "CNY" as const, scale: 2 as const });
      const requirement: FacilityRequirement = {
        schema_version: "1.0",
        requirement_id: requirementID,
        tenant_id: tenant,
        case_code: caseCode,
        legal_entity_code: entityCode,
        target_region: draft.targetRegion.trim(),
        facility_purpose: draft.facilityPurpose.trim(),
        minimum_area_m2: Number(draft.minimumAreaM2),
        minimum_electricity_kva: Number(draft.minimumElectricKVA),
        target_available_at: localDate(draft.targetAvailableAt),
        candidate_count: Number(draft.candidateCount),
        allowed_option_types: draft.optionTypes,
        investment_request: money(draft.investmentRequest),
        minimum_cash_reserve: money(draft.minimumCashReserve),
        financial_constraint: financial.financial_constraint,
        preferences: draft.preferences.split("\n").map((value) => value.trim()).filter(Boolean),
        revision: nextRevision,
        revision_reason: draft.revisionReason.trim(),
      };
      const result = await generatePlantProposals(requirement);
      setProposalSet(result.proposal_set);
      setReviews({});
      setSavedReviews({});
      setNextRevision((current) => current + 1);
    } catch (reason) {
      setError(reason instanceof Error ? reason.message : String(reason));
    } finally {
      setBusy(false);
    }
  };

  const submitReview = async (proposalID: string) => {
    if (!proposalSet) return;
    const review = reviews[proposalID];
    if (!review || !["adopt_for_investigation", "request_revision", "discard"].includes(review.action)) {
      setError("请先选择采纳调研、退回重生成或淘汰");
      return;
    }
    setReviewBusy(proposalID);
    setError("");
    try {
      await submitPlantProposalReview({
        proposal_set_id: proposalSet.proposal_set_id,
        proposal_id: proposalID,
        action: review.action as "adopt_for_investigation" | "request_revision" | "discard",
        reason: review.reason.trim(),
        expected_revision: proposalSet.revision,
      });
      setSavedReviews((current) => ({ ...current, [proposalID]: true }));
    } catch (reason) {
      setError(reason instanceof Error ? reason.message : String(reason));
    } finally {
      setReviewBusy("");
    }
  };

  const addManual = (event: FormEvent) => {
    event.preventDefault();
    const amount = manualAmount || "0.00";
    const now = new Date().toISOString();
    const proposal: SiteOptionProposal = {
      proposal_id: `manual-${Date.now()}`,
      option_type: "manual_option",
      display_name: manualName.trim(),
      business_rationale: manualRationale.trim(),
      estimated_amount: { minimum: { value: amount, currency: "CNY", scale: 2 }, likely: { value: amount, currency: "CNY", scale: 2 }, maximum: { value: amount, currency: "CNY", scale: 2 }, basis: "人工初始估算，待调研校验" },
      estimated_schedule: { earliest: now, likely: now, latest: now },
      assumptions: [], facts_required: ["补充外部调研和权威证据"], risks: [], source_refs: ["manual:user-input"], confidence: "0", status: "manual_draft",
    };
    setProposalSet((current) => current ? { ...current, proposals: [...current.proposals, proposal] } : {
      schema_version: "1.0", proposal_set_id: "local-manual-draft", requirement_id: "local-manual", revision: 1, status: "candidate_only", proposals: [proposal],
      evidence: { provider: "human", model: "none", prompt_version: "manual", input_hash: "local", output_hash: "local", validated_at: now },
    });
    setReviews((current) => ({ ...current, [proposal.proposal_id]: { action: "add_manual_option", reason: "由项目负责人手工补充候选" } }));
    setManualName(""); setManualRationale(""); setManualAmount(""); setManualOpen(false);
  };

  return (
    <section className="plant-planning-workspace" aria-labelledby="plant-planning-title">
      <header className="plant-planning-heading">
        <div><span>M10 INTERACTIVE PLANNING</span><h2 id="plant-planning-title">设施需求与候选方案</h2><p>先由人员定义约束和可调整金额，再让设施规划 Agent 提出候选；人员负责调研、选择和治理决定。</p></div>
        <div className={`plant-model-state ${status?.state ?? "loading"}`}><Bot /> <span>规划模型</span><strong>{status?.state === "connected" ? `${status.provider} · ${status.model}` : status?.state === "not_configured" ? "未启用外部模型" : "正在检查"}</strong></div>
      </header>
      <details className="plant-feature-help">
        <summary><CircleHelp />功能说明：这一步解决什么问题？</summary>
        <div><p>业务目的：把工厂需求、投资边界和工期转成可审阅的场址候选，不让系统用固定答案代替管理判断。</p><p>关系：M9 资格/资金与预算 → 设施需求 → Agent 候选 → 人工调研与选择 → IAOS 投资审批/合同/WBS → World 验收 → M11 能力建设资格。</p><p>Agent 只生成建议和待验证事实，不能批准投资、确认权属或伪造外部证明。可用现金和已批预算来自 IAOS 权威账务，只读且不能在本页面伪造。</p></div>
      </details>
      <form className="plant-requirement-form" onSubmit={submit}>
        <fieldset><legend>1 · 设施业务需求</legend>
          <label>目标区域 <small>Agent 搜索和比较方案的地理范围</small><input required value={draft.targetRegion} onChange={(e) => update("targetRegion", e.target.value)} placeholder="例如：江苏省苏州市及周边" /></label>
          <label>设施用途 <small>说明产品、工艺和运营目标</small><textarea required value={draft.facilityPurpose} onChange={(e) => update("facilityPurpose", e.target.value)} placeholder="例如：建设电池冷却板制造、质检与仓储基地" /></label>
          <label>最小面积（m²）<small>硬约束，可按企业情况修改</small><input required min="1" type="number" value={draft.minimumAreaM2} onChange={(e) => update("minimumAreaM2", e.target.value)} /></label>
          <label>最小电力容量（kVA）<small>需由供电资料或调研进一步验证</small><input required min="1" type="number" value={draft.minimumElectricKVA} onChange={(e) => update("minimumElectricKVA", e.target.value)} /></label>
          <label>目标可用时间<small>期望场地达到可使用状态的日期</small><input required type="datetime-local" value={draft.targetAvailableAt} onChange={(e) => update("targetAvailableAt", e.target.value)} /></label>
          <label>候选数量<small>Agent 生成 2–8 个不同候选</small><input required min="2" max="8" type="number" value={draft.candidateCount} onChange={(e) => update("candidateCount", e.target.value)} /></label>
        </fieldset>
        <fieldset><legend>2 · 投资与资金边界（CNY）</legend>
          <label>本次投资申请金额<small>本次希望申请的投资额度，可修改</small><input required min="0" step="0.01" type="number" value={draft.investmentRequest} onChange={(e) => update("investmentRequest", e.target.value)} /></label>
          <label>最低现金保留额<small>项目后仍需保留的安全现金，可修改</small><input required min="0" step="0.01" type="number" value={draft.minimumCashReserve} onChange={(e) => update("minimumCashReserve", e.target.value)} /></label>
          <label>可用现金快照<small>{financial?.financial_constraint.cash_source_ref || "正在读取 IAOS 总账"}</small><input readOnly aria-readonly="true" value={financial?.financial_constraint.available_cash.value ?? ""} placeholder="由 IAOS 权威账务提供" /></label>
          <label>已批准预算快照<small>{financial?.financial_constraint.budget_source_ref || "正在读取 IAOS 已批预算"}</small><input readOnly aria-readonly="true" value={financial?.financial_constraint.approved_budget.value ?? ""} placeholder="由 IAOS 已批预算提供" /></label>
        </fieldset>
        <fieldset><legend>3 · 方案范围与偏好</legend>
          <div className="plant-option-types"><strong>允许的方案类型</strong>{OPTION_TYPES.map(([code, label]) => <label key={code}><input type="checkbox" checked={draft.optionTypes.includes(code)} onChange={(e) => update("optionTypes", e.target.checked ? [...draft.optionTypes, code] : draft.optionTypes.filter((value) => value !== code))} />{label}</label>)}</div>
          <label>业务偏好<small>每行一项，如交通、人才、扩展性；不是硬约束</small><textarea value={draft.preferences} onChange={(e) => update("preferences", e.target.value)} placeholder={"靠近主要客户\n便于后续扩建"} /></label>
          <label>本次修订原因<small>将保存为需求版本 {nextRevision}，用于解释为什么产生或变化</small><input required value={draft.revisionReason} onChange={(e) => update("revisionReason", e.target.value)} /></label>
        </fieldset>
        <div className="plant-generation-actions"><button disabled={busy || status?.state !== "connected" || !financial || draft.optionTypes.length === 0}>{busy ? <LoaderCircle className="gx-spin" /> : <Bot />}{busy ? "Agent 正在生成…" : "保存需求并让 Agent 生成候选"}</button><button type="button" onClick={() => setManualOpen((value) => !value)}><FilePlus2 />人工新增候选</button><small>需求和候选将通过 IAOS Capability 留下审计证据，但不会自动批准投资。</small></div>
      </form>
      {status?.state === "not_configured" && <p className="plant-inline-warning" role="status"><TriangleAlert />外部模型未启用，不能生成虚拟固定候选；请配置模型，或使用“人工新增候选”。</p>}
      {error && <p className="plant-inline-error" role="alert">{error}</p>}
      {manualOpen && <form className="plant-manual-form" onSubmit={addManual}><h3>人工新增候选草稿</h3><label>候选名称<input required value={manualName} onChange={(e) => setManualName(e.target.value)} /></label><label>业务理由<textarea required value={manualRationale} onChange={(e) => setManualRationale(e.target.value)} /></label><label>初始估算（CNY）<input required min="0" step="0.01" type="number" value={manualAmount} onChange={(e) => setManualAmount(e.target.value)} /></label><button>加入本地审阅列表</button></form>}
      {proposalSet && <section className="plant-proposals" aria-live="polite"><header><div><span>{proposalSet.status}</span><h2>候选方案 · {proposalSet.proposals.length}</h2></div><small>{proposalSet.evidence.provider} / {proposalSet.evidence.model} · {proposalSet.evidence.prompt_version}</small></header><div className="plant-proposal-grid">{proposalSet.proposals.map((proposal) => <ProposalCard key={proposal.proposal_id} proposal={proposal} review={reviews[proposal.proposal_id]} onReview={(review) => { setReviews((current) => ({ ...current, [proposal.proposal_id]: review })); setSavedReviews((current) => ({ ...current, [proposal.proposal_id]: false })); }} onSubmitReview={() => submitReview(proposal.proposal_id)} saved={Boolean(savedReviews[proposal.proposal_id])} busy={reviewBusy === proposal.proposal_id} />)}</div><p className="plant-persistence-note"><BadgeCheck />Agent 候选和相应人工审阅通过 IAOS Capability 保存；人工新增候选在其正式 Capability 完成前明确保留为本地草稿。“采纳调研”不等于投资批准或合同签署。</p></section>}
    </section>
  );
}

export function PlantBuildPlay({ onExit }: { onExit: () => void }) {
  const [trace, setTrace] = useState<PlantBuildTrace | null>(null);
  const [step, setStep] = useState(0);
  const [playing, setPlaying] = useState(false);
  const [error, setError] = useState("");
  useEffect(() => { const controller = new AbortController(); loadPlantBuild(controller.signal).then(setTrace).catch((reason) => { if (reason.name !== "AbortError") setError(String(reason)); }); return () => controller.abort(); }, []);
  useEffect(() => { if (!playing || !trace) return; const id = setInterval(() => setStep((value) => { if (value >= trace.frames.length - 1) { setPlaying(false); return value; } return value + 1; }), 750); return () => clearInterval(id); }, [playing, trace]);
  const frame = useMemo(() => trace?.frames[step], [trace, step]);
  if (error) return <main className="world-error" role="alert">Plant Build 加载失败：{error}</main>;
  if (!trace || !frame) return <main className="world-loading">正在加载设施建设世界…</main>;
  return <div className="world-play">
    <header className="world-toolbar"><button className="world-back" onClick={onExit}><ArrowLeft />企业成立</button><div><span>PROJECT GENESIS · PLANT BUILD</span><h1>M10 工厂选址与设施建设</h1></div><div className="world-clock"><small>虚拟时间 · Asia/Shanghai</small><strong>{new Date(frame.sim_time).toLocaleString("zh-CN", { timeZone: "Asia/Shanghai", hour12: false })}</strong></div></header>
    <main className="world-main"><PlantPlanningWorkspace />
      <details className="plant-reference-replay"><summary>查看已封存的确定性参考回放（fixture-only）</summary><p>以下场址和金额仅用于历史回归测试，不会进入上方需求表单，也不能作为新企业的真实候选或预算。</p>
        <nav className="world-controls"><button onClick={() => setPlaying((value) => !value)}>{playing ? <Pause /> : <Play />}{playing ? "暂停" : "运行"}</button><button onClick={() => setStep((value) => Math.min(value + 1, trace.frames.length - 1))} disabled={step === trace.frames.length - 1}><StepForward />单步</button><button onClick={() => { setPlaying(false); setStep(0); }}><RotateCcw />复位</button><button onClick={() => (window.location.hash = "world-capability-build")}><Factory />生产能力 Campaign</button><span>{step + 1}/{trace.frames.length} · {frame.phase}</span></nav>
        <section className="world-status"><span className={`world-state-badge ${frame.capability_build_eligible ? "closed" : "active"}`}>{frame.capability_build_eligible ? <BadgeCheck /> : <Factory />}{frame.capability_build_eligible ? "M11 eligible" : frame.phase}</span><h2>{frame.title}</h2><p>World {frame.world_progress}% / IAOS plan {frame.iaos_plan_progress}% · {frame.discrepancy} · cursor {frame.iaos_cursor}</p></section>
        <section className="three-state-grid"><article><header><Building2 />场址与项目 <small>World owns reality</small></header><dl><div><dt>选中场址</dt><dd>{frame.selected_site || "尚未批准"}</dd></div>{frame.assessments.map((assessment) => <div key={assessment.site_code}><dt>{assessment.site_code.replace("SITE-SZ-", "")}</dt><dd>{assessment.feasible ? "可行" : "不可行"} · {assessment.weighted_score}<small>{assessment.hard_failures.join(", ") || "硬约束通过"}</small></dd></div>)}</dl></article><article><header><Banknote />治理资金 <small>IAOS owns approvals</small></header><dl><div><dt>可用现金</dt><dd>{cny(frame.cash.value)}</dd></div><div><dt>合同承诺</dt><dd>{cny(frame.committed.value)}</dd></div><div><dt>里程碑应付</dt><dd>{cny(frame.payable.value)}</dd></div><div><dt>已付款</dt><dd>{cny(frame.paid.value)}</dd></div></dl></article><article><header><TriangleAlert />Knowledge / discrepancy <small>actor scoped</small></header><dl><div><dt>项目负责人已知</dt><dd>{frame.knowledge.length} 条</dd></div><div><dt>差异</dt><dd>{frame.discrepancy}</dd></div><div><dt>因果</dt><dd>{frame.causation_id}</dd></div></dl></article></section>
        {frame.zones.length > 0 && <section className="facility-zones">{frame.zones.filter((zone) => !["site", "building"].includes(zone.purpose)).map((zone) => <article key={zone.code}><Factory /><strong>{zone.purpose}</strong><code>{zone.code}</code><small>{zone.area_m2} m² · {zone.status}</small></article>)}</section>}
        <section className="world-timeline">{trace.frames.map((item, index) => <button key={item.step} className={index === step ? "current" : ""} onClick={() => { setPlaying(false); setStep(index); }}><span>{index + 1}</span><strong>{item.phase}</strong><small>{item.title}</small></button>)}</section>
      </details>
    </main>
  </div>;
}
