import {
  ArrowLeft,
  BadgeCheck,
  Bot,
  CircleHelp,
  FilePlus2,
  Gavel,
  Handshake,
  LoaderCircle,
  RotateCcw,
  SearchCheck,
  TriangleAlert,
  X,
  XCircle,
} from "lucide-react";
import { type FormEvent, useEffect, useMemo, useState } from "react";
import {
  generatePlantProposals,
  generatePlantRequirementOptions,
  loadPlantFinancialConstraint,
  loadPlantApprovalDetail,
  loadPlantPlanningStatus,
  loadPlantProposalSet,
  loadPlantProposalReviews,
  loadPlantRequirement,
  loadSiteInvestigations,
  loadSiteSelections,
  loadSiteControls,
  requestSiteInvestigation,
  requestSiteControl,
  finalizeSiteSelection,
  confirmSiteInvestigationReport,
  submitSiteSelectionRecommendation,
  confirmSiteControlDelivery,
  generateFacilityProjectOptions,
  submitFacilityProjectOption,
  loadFacilityProjects,
  activateFacilityProject,
  submitFacilityProjectBaseline,
  loadContractAwards,
  issueContractRFQ,
  confirmContractBids,
  generateContractRecommendation,
  submitContractRecommendation,
  awardContract,
  decidePlantApproval,
  submitPlantProposalReview,
  submitManualPlantProposal,
  type FacilityRequirement,
  type PlantPlanningProviderStatus,
  type PlantFinancialConstraint,
  type PlantRequirementOption,
  type PlantRequirementOptionSet,
  type PlantApprovalDetail,
  type ProposalSet,
  type SiteOptionProposal,
  type SiteInvestigationItem,
  type SiteSelectionItem,
  type SiteControlItem,
  type SiteControlRequest,
  type ProjectPlanOption,
  type ProjectPlanOptionSet,
  type FacilityProjectItem,
  type ContractAwardItem,
  type ContractRecommendationAdvice,
  type ContractRFQ,
} from "../../world/plantBuild";
import {
  DEFAULT_SITE_ASSESSMENT_WEIGHTS,
  assessObservedSites,
  currentProposalInvestigations,
  type SiteAssessment,
  type SiteAssessmentWeights,
} from "../../world/siteAssessment";
import { createClientRequestId } from "../../lib/clientRequestId";
import { derivePlantGameStage } from "../../world/plantGame";
import { PlantBuildGameScene, type PlantSceneArchiveEntry } from "./PlantBuildGameScene";
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

const localDateTimeInput = (value: string) =>
  value ? new Date(value).toISOString().slice(0, 16) : "";

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
  investigation,
  investigationBusy,
  onRequestInvestigation,
}: {
  proposal: SiteOptionProposal;
  review?: ReviewState;
  onReview: (value: ReviewState) => void;
  onSubmitReview: () => void;
  saved: boolean;
  busy: boolean;
  investigation?: SiteInvestigationItem;
  investigationBusy: boolean;
  onRequestInvestigation: () => void;
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
        <button type="button" disabled={saved} className={action === "adopt_for_investigation" ? "selected" : ""} onClick={() => onReview({ action: "adopt_for_investigation", reason: review?.reason ?? "" })}><SearchCheck />采纳调研</button>
        <button type="button" disabled={saved} className={action === "request_revision" ? "selected" : ""} onClick={() => onReview({ action: "request_revision", reason: review?.reason ?? "" })}><RotateCcw />退回重生成</button>
        <button type="button" disabled={saved} className={action === "discard" ? "selected danger" : ""} onClick={() => onReview({ action: "discard", reason: review?.reason ?? "" })}><XCircle />淘汰</button>
        <label>
          审阅理由
          <textarea disabled={saved} minLength={6} value={review?.reason ?? ""} onChange={(event) => onReview({ action, reason: event.target.value })} placeholder="说明采纳、退回或淘汰的业务理由（至少 6 个字符）" />
        </label>
        <button type="button" className="plant-review-submit" disabled={!formal || !action || (review?.reason.trim().length ?? 0) < 6 || busy || saved} onClick={onSubmitReview}>
          {busy ? <LoaderCircle className="gx-spin" /> : formal ? <BadgeCheck /> : <TriangleAlert />}{!formal ? "人工候选仍是本地草稿" : saved ? "已保存审阅" : busy ? "正在保存…" : "提交审阅到 IAOS"}
        </button>
        {saved && action === "adopt_for_investigation" && <button type="button" className="plant-investigation-request" disabled={investigationBusy || Boolean(investigation)} onClick={onRequestInvestigation}>
          {investigationBusy ? <LoaderCircle className="gx-spin" /> : <SearchCheck />}{investigation ? investigation.status === "observed" ? "调研事实已提交" : "正在等待外部调研" : investigationBusy ? "正在建立工作项…" : "发起外部调研工作项"}
        </button>}
      </fieldset>
    </article>
  );
}

function InvestigationPanel({ item, caseCode, requirementID, onCommitted }: { item: SiteInvestigationItem; caseCode: string; requirementID: string; onCommitted: () => void }) {
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState("");
  const acceptReport = async () => {
    setBusy(true); setError("");
    try {
      await confirmSiteInvestigationReport(caseCode, requirementID, item.request.investigation_request_id);
      onCommitted();
    } catch (reason) { setError(reason instanceof Error ? reason.message : String(reason)); } finally { setBusy(false); }
  };
  return <article className={`plant-investigation-card ${item.status}`}>
    <header><div><span>World wait · {item.work_item_status}</span><h3>{item.request.investigation_request_id}</h3></div><strong>{item.status === "observed" ? "可信事实已提交" : "等待园区运营方"}</strong></header>
    <p>候选 {item.request.proposal_id} · 调研范围：{item.request.scope.join("、")}</p>
    {item.observation ? <dl className="plant-investigation-facts"><div><dt>权属</dt><dd>{item.observation.ownership_status}</dd></div><div><dt>可用面积</dt><dd>{item.observation.available_area_m2} m²</dd></div><div><dt>电力</dt><dd>{item.observation.electricity_kva} kVA</dd></div><div><dt>正式报价</dt><dd>{cny(item.observation.quoted_amount.value)}</dd></div><div><dt>许可状态</dt><dd>{item.observation.permit_status}</dd></div></dl> :
      <section className="plant-control-confirmation plant-investigation-confirmation">
        <div className="plant-control-world-badge"><SearchCheck /><span><small>WORLD ENGINE · 园区调研团队</small><strong>现场调研报告已经准备完成</strong></span></div>
        <p>园区调研团队会依据权威调研工作项生成权属、面积、电力、正式报价、可用日期、许可与证据。你只需确认接收报告，不需要代替外部人员填写数据。</p>
        <dl><div><dt>调研候选</dt><dd>{item.request.proposal_id}</dd></div><div><dt>核验范围</dt><dd>权属 · 报价 · 面积 · 电力 · 日期 · 许可</dd></div><div><dt>归档位置</dt><dd>World Journal → IAOS Observation</dd></div></dl>
        {error && <p role="alert" className="plant-inline-error">{error}</p>}
        <button type="button" onClick={acceptReport} disabled={busy}>{busy ? <LoaderCircle className="gx-spin" /> : <BadgeCheck />}{busy ? "World 正在归档报告…" : "确认接收调研报告"}</button>
        <small>确认后由服务端生成可重放事实，先进入 World Journal，再由受治理 Capability 核验并关闭工作项。</small>
      </section>}
  </article>;
}

const SCORE_DIMENSIONS: Array<[keyof SiteAssessmentWeights, string, string]> = [
  ["cost", "成本", "正式报价相对投资申请额"],
  ["schedule", "工期", "正式可用日期相对目标日期"],
  ["capacity", "容量", "实际面积与电力容量"],
  ["control", "权属与许可", "权属与许可核验结论"],
];

function AssessmentCard({ assessment }: { assessment: SiteAssessment }) {
  return <article className={`plant-assessment-card ${assessment.eligible ? "eligible" : "ineligible"}`}>
    <header>
      <div><span>{assessment.observation_id}</span><h3>{assessment.display_name}</h3></div>
      <strong>{assessment.eligible ? `${assessment.total_score?.toFixed(1)} 分` : "硬约束不通过"}</strong>
    </header>
    {assessment.hard_failures.length > 0
      ? <div className="plant-hard-failures"><strong><XCircle />先解决 {assessment.hard_failures.length} 项硬约束</strong><ul>{assessment.hard_failures.map((failure) => <li key={failure}>{failure}</li>)}</ul></div>
      : <p className="plant-hard-pass"><BadgeCheck />投资、面积、电力、日期、权属和许可硬约束通过</p>}
    <div className="plant-criteria-comparison" aria-label={`${assessment.display_name} 硬约束逐项比较`}>
      {assessment.criteria.map((criterion) => <article key={criterion.key} className={criterion.passed ? "passed" : "failed"}>
        <span className="plant-criterion-state">{criterion.passed ? <BadgeCheck /> : <XCircle />}<b>{criterion.label}</b><em>{criterion.passed ? "通过" : "不通过"}</em></span>
        <div><small>需求门槛 · Facility Requirement</small><strong>{criterion.required}</strong></div>
        <div><small>现场事实 · World Observation</small><strong>{criterion.observed}</strong></div>
        <div className="plant-criterion-difference"><small>比较结果</small><strong>{criterion.difference}</strong></div>
      </article>)}
    </div>
    <div className="plant-estimate-observation">
      <section><span>Agent 估算 · 非正式事实</span><strong>{assessment.estimated ? cny(assessment.estimated.amount) : "当前会话未加载"}</strong><small>{assessment.estimated ? new Date(assessment.estimated.available_at).toLocaleDateString("zh-CN") : "以 Observation 为评分依据"}</small></section>
      <section><span>World Observation · 评分事实</span><strong>{cny(assessment.observed.quote)}</strong><small>{new Date(assessment.observed.available_at).toLocaleDateString("zh-CN")} · {assessment.observed.area_m2} m² · {assessment.observed.electricity_kva} kVA</small></section>
    </div>
    {assessment.eligible && <dl className="plant-score-components">{SCORE_DIMENSIONS.map(([key, label]) => <div key={key}><dt>{label}</dt><dd>{assessment.component_scores[key].toFixed(1)}</dd></div>)}</dl>}
    <details className="plant-assessment-evidence"><summary>查看事实与证据引用</summary><dl><div><dt>权属</dt><dd>{assessment.observed.ownership}</dd></div><div><dt>许可</dt><dd>{assessment.observed.permit}</dd></div></dl><ul>{assessment.evidence_refs.map((reference) => <li key={reference}>{reference}</li>)}</ul></details>
  </article>;
}

function PlantApprovalPanel({
  selection,
  detail,
  loadError,
  note,
  busy,
  approvalURL,
  onNote,
  onDecide,
  onFormalize,
}: {
  selection: SiteSelectionItem;
  detail: PlantApprovalDetail | null;
  loadError: string;
  note: string;
  busy: boolean;
  approvalURL: string;
  onNote: (value: string) => void;
  onDecide: (decision: "approve" | "reject") => void;
  onFormalize: () => void;
}) {
  const recommendation = selection.recommendation;
  const pending = selection.approval_status === "pending";
  const selectedAssessment = recommendation.assessments.find((item) => item.proposal_id === recommendation.selected_proposal_id);
  return <article id="plant-task-governance" className="plant-selection-status plant-game-approval" aria-live="polite">
    <header><div><span>{recommendation.approval_flow_key}</span><h3>治理会议：审阅场址正式选择</h3></div><strong>{selection.decision ? "已正式选址" : selection.approval_status}</strong></header>
    <section className="plant-approval-subject">
      <div><small>待审批事项</small><h4>{detail?.detail.subject.title ?? `${recommendation.case_code} · 工厂场址正式选择`}</h4><p>{detail?.detail.subject.summary ?? "审阅可信调研事实、硬约束、评分和人员推荐，决定是否允许正式选址。"}</p></div>
      <dl><div><dt>推荐候选</dt><dd>{recommendation.selected_proposal_id}</dd></div><div><dt>推荐理由</dt><dd>{recommendation.recommendation_reason}</dd></div><div><dt>替代方案比较</dt><dd>{recommendation.alternative_comparison}</dd></div>{recommendation.single_source_exception_reason && <div><dt>单一来源例外</dt><dd>{recommendation.single_source_exception_reason}</dd></div>}<div><dt>权威评分</dt><dd>{selectedAssessment?.total_score ?? "未返回"}</dd></div><div><dt>Observation</dt><dd>{selectedAssessment?.observation_id ?? "未返回"}</dd></div><div><dt>输入哈希</dt><dd>{recommendation.input_hash}</dd></div></dl>
    </section>
    {detail && <section className="plant-approval-routing"><header><div><small>IAOS 审批流</small><strong>{detail.detail.flow_name} · v{detail.detail.flow_version}</strong></div><span>{detail.detail.can_decide ? "当前身份可决定" : "当前身份仅可查看"}</span></header><div>{detail.detail.assignments.map((assignment) => <article key={assignment.id}><Gavel /><span><strong>{assignment.stage_name}</strong><small>{assignment.display_name} · {assignment.selector_type}:{assignment.selector_value}</small></span><em>{assignment.status}</em></article>)}</div></section>}
    {loadError && <p className="plant-inline-error" role="alert">{loadError}</p>}
    {pending && detail && <section className="plant-approval-decision"><label>审批意见 <small>批准或驳回都将由 IAOS 记录决定人、时间和意见</small><textarea required minLength={6} value={note} onChange={(event) => onNote(event.target.value)} placeholder="请说明同意的条件，或明确驳回原因" /></label>{detail.detail.can_decide ? <div><button type="button" disabled={busy || note.trim().length < 6} onClick={() => onDecide("approve")}><BadgeCheck />批准选址推荐</button><button type="button" className="danger" disabled={busy || note.trim().length < 6} onClick={() => onDecide("reject")}><XCircle />驳回并退回修订</button></div> : <p><TriangleAlert />当前登录身份不是该阶段的受派审批人。请由上方显示的治理主体进入本企业后，在本会议室作出决定。</p>}</section>}
    <div className="plant-selection-actions"><a href={approvalURL} target="_blank" rel="noreferrer">在 IAOS 查看审计详情</a>{selection.approval_status === "approved" && !selection.decision && <button type="button" disabled={busy} onClick={onFormalize}>批准已生效 · 正式选址</button>}{selection.decision && <span><BadgeCheck />正式选择已落地</span>}{selection.approval_status === "rejected" && <span><XCircle />推荐已驳回，不能正式选址</span>}</div>
  </article>;
}

function PlantPlanningWorkspace({ onExit }: { onExit: () => void }) {
  const [status, setStatus] = useState<PlantPlanningProviderStatus | null>(null);
  const [financial, setFinancial] = useState<PlantFinancialConstraint | null>(null);
  const [activeRequirement, setActiveRequirement] = useState<FacilityRequirement | null>(null);
  const [draft, setDraft] = useState<RequirementDraft>(blankDraft);
  const [requirementOptions, setRequirementOptions] = useState<PlantRequirementOptionSet | null>(null);
  const [selectedRequirementOption, setSelectedRequirementOption] = useState<PlantRequirementOption | null>(null);
  const [requirementAdviserBusy, setRequirementAdviserBusy] = useState(false);
  const [proposalSet, setProposalSet] = useState<ProposalSet | null>(null);
  const [reviews, setReviews] = useState<Record<string, ReviewState>>({});
  const [savedReviews, setSavedReviews] = useState<Record<string, boolean>>({});
  const [reviewBusy, setReviewBusy] = useState("");
  const [investigations, setInvestigations] = useState<SiteInvestigationItem[]>([]);
  const [investigationBusy, setInvestigationBusy] = useState("");
  const [assessmentWeights, setAssessmentWeights] = useState<SiteAssessmentWeights>(DEFAULT_SITE_ASSESSMENT_WEIGHTS);
  const [siteSelections, setSiteSelections] = useState<SiteSelectionItem[]>([]);
  const [siteControls, setSiteControls] = useState<SiteControlItem[]>([]);
  const [projectOptions, setProjectOptions] = useState<ProjectPlanOptionSet | null>(null);
  const [selectedProjectOption, setSelectedProjectOption] = useState<ProjectPlanOption | null>(null);
  const [facilityProjects, setFacilityProjects] = useState<FacilityProjectItem[]>([]);
  const [projectApprovalDetail, setProjectApprovalDetail] = useState<PlantApprovalDetail | null>(null);
  const [projectApprovalNote, setProjectApprovalNote] = useState("");
  const [projectBusy, setProjectBusy] = useState(false);
  const [contractAwards, setContractAwards] = useState<ContractAwardItem[]>([]);
  const [contractAdvice, setContractAdvice] = useState<ContractRecommendationAdvice | null>(null);
  const [contractApprovalDetail, setContractApprovalDetail] = useState<PlantApprovalDetail | null>(null);
  const [selectedPackageCode, setSelectedPackageCode] = useState("");
  const [sourcingStrategy, setSourcingStrategy] = useState<ContractRFQ["sourcing_strategy"]>("specialist_packages");
  const [contractBusy, setContractBusy] = useState(false);
  const [controlMode, setControlMode] = useState<SiteControlRequest["agreement_mode"]>("lease");
  const [controlHandoverAt, setControlHandoverAt] = useState("");
  const [controlBusy, setControlBusy] = useState(false);
  const [selectedProposalID, setSelectedProposalID] = useState("");
  const [recommendationReason, setRecommendationReason] = useState("");
  const [alternativeComparison, setAlternativeComparison] = useState("");
  const [singleSourceReason, setSingleSourceReason] = useState("");
  const [selectionBusy, setSelectionBusy] = useState(false);
  const [approvalDetail, setApprovalDetail] = useState<PlantApprovalDetail | null>(null);
  const [approvalNote, setApprovalNote] = useState("");
  const [approvalLoadError, setApprovalLoadError] = useState("");
  const [nextRevision, setNextRevision] = useState(1);
  const [manualOpen, setManualOpen] = useState(false);
  const [manualName, setManualName] = useState("");
  const [manualRationale, setManualRationale] = useState("");
  const [manualOptionType, setManualOptionType] = useState("");
  const [manualMinimum, setManualMinimum] = useState("");
  const [manualLikely, setManualLikely] = useState("");
  const [manualMaximum, setManualMaximum] = useState("");
  const [manualBasis, setManualBasis] = useState("");
  const [manualAvailableAt, setManualAvailableAt] = useState("");
  const [manualAssumptions, setManualAssumptions] = useState("");
  const [manualFacts, setManualFacts] = useState("");
  const [manualRisks, setManualRisks] = useState("");
  const [manualBusy, setManualBusy] = useState(false);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState("");
  const [taskOpen, setTaskOpen] = useState(false);
  const [taskView, setTaskView] = useState<"current" | "requirement-revision">("current");
  const routeParams = new URLSearchParams(window.location.hash.split("?")[1] ?? "");
  const tenant = routeParams.get("tenant") ?? localStorage.getItem("aese_iaos_tenant_id") ?? localStorage.getItem("iaos_tenant_id") ?? "";
  const caseCode = routeParams.get("case") ?? localStorage.getItem("aese_genesis_case_code") ?? "";
  const requirementID = `facility-requirement-${caseCode || "draft"}`;
  const refreshInvestigations = async () => setInvestigations((await loadSiteInvestigations(caseCode)).items ?? []);
  const refreshSiteSelections = async () => setSiteSelections((await loadSiteSelections(caseCode)).items ?? []);
  const refreshSiteControls = async () => setSiteControls((await loadSiteControls(caseCode)).items ?? []);
  const refreshFacilityProjects = async () => setFacilityProjects((await loadFacilityProjects(caseCode)).items ?? []);
  const refreshContractAwards = async () => setContractAwards((await loadContractAwards(caseCode)).items ?? []);
  useEffect(() => {
    const controller = new AbortController();
    Promise.all([
      loadPlantPlanningStatus(controller.signal).then(setStatus),
      loadPlantFinancialConstraint(caseCode, controller.signal).then(setFinancial),
      loadPlantRequirement(requirementID, controller.signal).then((existing) => {
        setActiveRequirement(existing?.investment_request && existing.target_available_at ? existing : null);
        setNextRevision((existing?.revision ?? 0) + 1);
      }),
      loadPlantProposalSet(requirementID, controller.signal).then(async (set) => {
        setProposalSet(set);
        if (!set?.proposal_set_id) return;
        const persisted = await loadPlantProposalReviews(set.proposal_set_id, controller.signal);
        const reviewState: Record<string, ReviewState> = {};
        const savedState: Record<string, boolean> = {};
        for (const review of persisted.items ?? []) {
          reviewState[review.proposal_id] = { action: review.action, reason: review.reason };
          savedState[review.proposal_id] = true;
        }
        setReviews(reviewState);
        setSavedReviews(savedState);
      }),
      loadSiteInvestigations(caseCode, controller.signal).then((result) => setInvestigations(result.items ?? [])),
      loadSiteSelections(caseCode, controller.signal).then((result) => setSiteSelections(result.items ?? [])),
      loadSiteControls(caseCode, controller.signal).then((result) => setSiteControls(result.items ?? [])),
      loadFacilityProjects(caseCode, controller.signal).then((result) => setFacilityProjects(result.items ?? [])),
      loadContractAwards(caseCode, controller.signal).then((result) => setContractAwards(result.items ?? [])),
    ]).catch((reason) => { if (reason.name !== "AbortError") setError(String(reason)); });
    return () => controller.abort();
  }, [caseCode, requirementID]);
  const update = (name: keyof RequirementDraft, value: string | string[]) =>
    setDraft((current) => ({ ...current, [name]: value }));
  const prepareRequirementOptions = async () => {
    setRequirementAdviserBusy(true); setError("");
    try {
      const result = await generatePlantRequirementOptions(caseCode);
      setRequirementOptions(result);
      setSelectedRequirementOption(null);
    } catch (reason) {
      setError(reason instanceof Error ? reason.message : String(reason));
    } finally {
      setRequirementAdviserBusy(false);
    }
  };
  const chooseRequirementOption = (option: PlantRequirementOption) => {
    setSelectedRequirementOption(option);
    setDraft({
      targetRegion: option.target_region,
      facilityPurpose: option.facility_purpose,
      minimumAreaM2: String(option.minimum_area_m2),
      minimumElectricKVA: String(option.minimum_electricity_kva),
      targetAvailableAt: localDateTimeInput(option.target_available_at),
      candidateCount: String(option.candidate_count),
      optionTypes: [...option.allowed_option_types],
      investmentRequest: option.investment_request.value,
      minimumCashReserve: option.minimum_cash_reserve.value,
      preferences: option.preferences.join("\n"),
      revisionReason: `采纳 Agent 草案：${option.title}`,
    });
  };
  const openRequirementRevision = () => {
    if (!activeRequirement) return;
    setDraft({
      targetRegion: activeRequirement.target_region,
      facilityPurpose: activeRequirement.facility_purpose,
      minimumAreaM2: String(activeRequirement.minimum_area_m2),
      minimumElectricKVA: String(activeRequirement.minimum_electricity_kva),
      targetAvailableAt: localDateTimeInput(activeRequirement.target_available_at),
      candidateCount: String(activeRequirement.candidate_count),
      optionTypes: [...activeRequirement.allowed_option_types],
      investmentRequest: activeRequirement.investment_request.value,
      minimumCashReserve: activeRequirement.minimum_cash_reserve.value,
      preferences: activeRequirement.preferences.join("\n"),
      revisionReason: "",
    });
    setTaskView("requirement-revision");
  };
  const entityCode = financial?.legal_entity_code ?? "";
  const currentInvestigations = useMemo(() => currentProposalInvestigations(proposalSet, investigations), [investigations, proposalSet]);
  const historicalInvestigationCount = investigations.length - currentInvestigations.length;
  const assessments = useMemo(
    () => activeRequirement ? assessObservedSites(activeRequirement, proposalSet, currentInvestigations, assessmentWeights) : [],
    [activeRequirement, assessmentWeights, currentInvestigations, proposalSet],
  );
  const normalizedWeightTotal = Object.values(assessmentWeights).reduce((sum, value) => sum + Math.max(0, value), 0);
  const eligibleAssessments = assessments.filter((item) => item.eligible);
  const latestSelection = siteSelections[0];
  const latestSiteControl = siteControls[0];
  const latestFacilityProject = facilityProjects.at(-1);
  const latestContractAward = contractAwards.at(-1);
  const projectCanRevise = latestFacilityProject?.approval_status === "rejected";
  useEffect(() => {
    const approvalID = latestSelection?.recommendation.approval_request_id;
    if (!approvalID) { setApprovalDetail(null); setApprovalLoadError(""); return; }
    const controller = new AbortController();
    setApprovalLoadError("");
    loadPlantApprovalDetail(approvalID, controller.signal).then(setApprovalDetail).catch((reason) => {
      if (reason.name !== "AbortError") setApprovalLoadError(reason instanceof Error ? reason.message : String(reason));
    });
    return () => controller.abort();
  }, [latestSelection?.recommendation.approval_request_id, latestSelection?.approval_status]);
  useEffect(() => {
    const approvalID = latestFacilityProject?.approval_request_id;
    if (!approvalID) { setProjectApprovalDetail(null); return; }
    const controller = new AbortController();
    loadPlantApprovalDetail(approvalID, controller.signal).then(setProjectApprovalDetail).catch((reason) => {
      if (reason.name !== "AbortError") setError(reason instanceof Error ? reason.message : String(reason));
    });
    return () => controller.abort();
  }, [latestFacilityProject?.approval_request_id, latestFacilityProject?.approval_status]);
  useEffect(() => {
    const approvalID = latestContractAward?.approval_request_id;
    if (!approvalID) { setContractApprovalDetail(null); return; }
    const controller = new AbortController();
    loadPlantApprovalDetail(approvalID, controller.signal).then(setContractApprovalDetail).catch((reason) => {
      if (reason.name !== "AbortError") setError(reason instanceof Error ? reason.message : String(reason));
    });
    return () => controller.abort();
  }, [latestContractAward?.approval_request_id, latestContractAward?.approval_status]);
  const observedCount = currentInvestigations.filter((item) => item.status === "observed").length;
  const adoptedReviewCount = Object.values(reviews).filter((review) => review.action === "adopt_for_investigation").length;
  const derivedGameStage = derivePlantGameStage({
    hasRequirement: Boolean(activeRequirement),
    proposalCount: proposalSet?.proposals.length ?? 0,
    adoptedReviewCount,
    investigationCount: investigations.length,
    observedCount,
    hasRecommendation: Boolean(latestSelection?.recommendation),
    hasDecision: Boolean(latestSelection?.decision),
    hasSiteControlRequest: Boolean(latestSiteControl?.request),
    hasSiteControl: latestSiteControl?.status === "controlled",
    hasProjectPlan: Boolean(latestFacilityProject?.plan?.plan_id),
    projectApprovalStatus: latestFacilityProject?.approval_status ?? "",
    hasActiveProject: Boolean(latestFacilityProject?.project?.project_id),
    hasContractRFQ: Boolean(latestContractAward?.rfq?.rfq_id),
    hasContractBids: Boolean(latestContractAward?.observation?.observation_id),
    hasContractRecommendation: Boolean(latestContractAward?.recommendation?.recommendation_id),
    contractApprovalStatus: latestContractAward?.approval_status ?? "",
    hasAwardedContract: Boolean(latestContractAward?.contract?.contract_id),
  });
  const gameStage = derivedGameStage.key === "proposal" && !proposalSet
    ? { ...derivedGameStage, anchor: "plant-task-requirement", actionLabel: "继续填写需求并生成候选" }
    : derivedGameStage;
  const archives: PlantSceneArchiveEntry[] = [
    ...(activeRequirement ? [{
      id: activeRequirement.requirement_id,
      location: "headquarters" as const,
      title: `设施需求 · 第 ${activeRequirement.revision} 版`,
      state: "已确认",
      summary: `${activeRequirement.target_region} · ${activeRequirement.facility_purpose} · 投资申请 ${cny(activeRequirement.investment_request.value)}`,
      evidence: `${activeRequirement.requirement_id} · ${activeRequirement.financial_constraint.snapshot_hash}`,
    }] : []),
    ...(proposalSet ? [{
      id: proposalSet.proposal_set_id,
      location: "city" as const,
      title: `场址候选集 · 第 ${proposalSet.revision} 版`,
      state: proposalSet.status,
      summary: `${proposalSet.proposals.length} 个候选 · ${proposalSet.evidence.provider} / ${proposalSet.evidence.model}`,
      evidence: `${proposalSet.proposal_set_id} · ${proposalSet.evidence.output_hash}`,
    }] : []),
    ...Object.entries(savedReviews).filter(([, saved]) => saved).map(([proposalID]) => ({
      id: `review-${proposalID}`,
      location: "city" as const,
      title: `候选评审 · ${proposalID}`,
      state: reviews[proposalID]?.action ?? "已评审",
      summary: reviews[proposalID]?.reason ?? "人员评审已经提交 IAOS",
      evidence: `site.proposal.review:${proposalID}`,
    })),
    ...investigations.map((item) => {
      const belongsToCurrentSet = currentInvestigations.some((current) => current.request.investigation_request_id === item.request.investigation_request_id);
      return {
        id: item.request.investigation_request_id,
        location: "site" as const,
        title: `场址调研 · ${item.request.proposal_id}`,
        state: belongsToCurrentSet ? (item.status === "observed" ? "可信事实已提交" : "等待外部参与者") : "历史版本 · 不参与当前比较",
        summary: `${belongsToCurrentSet ? "" : "旧候选集证据 · "}${item.observation ? `面积 ${item.observation.available_area_m2} m² · 电力 ${item.observation.electricity_kva} kVA · 报价 ${cny(item.observation.quoted_amount.value)}` : `调研范围：${item.request.scope.join("、")}`}`,
        evidence: item.observation?.observation_id ?? item.request.investigation_request_id,
      };
    }),
    ...(latestSelection ? [{
      id: latestSelection.recommendation.recommendation_id,
      location: "boardroom" as const,
      title: latestSelection.decision ? "正式场址决定" : "场址推荐与审批",
      state: latestSelection.decision ? "已正式选址" : latestSelection.approval_status,
      summary: `推荐候选 ${latestSelection.recommendation.selected_proposal_id} · 审批请求 ${latestSelection.recommendation.approval_request_id}`,
      evidence: `${latestSelection.recommendation.recommendation_id} · ${latestSelection.recommendation.input_hash}`,
    }, ...(latestSelection.decision ? [{
      id: `formal-${latestSelection.recommendation.recommendation_id}`,
      location: "site" as const,
      title: "正式场址已经落地",
      state: "Committed",
      summary: `已批准候选 ${latestSelection.recommendation.selected_proposal_id}，等待项目与 WBS 建设阶段。`,
      evidence: `site.selection.formalize:${latestSelection.recommendation.recommendation_id}`,
    }] : [])] : []),
    ...(latestSiteControl ? [{
      id: latestSiteControl.request.control_request_id,
      location: "site" as const,
      title: latestSiteControl.status === "controlled" ? "场地控制权已交付" : "场地控制交付工作项",
      state: latestSiteControl.status,
      summary: latestSiteControl.observation
        ? `${latestSiteControl.observation.result} · ${latestSiteControl.observation.agreement_ref || "尚无协议引用"}`
        : `${latestSiteControl.request.agreement_mode} · 等待园区权利方提交交付事实`,
      evidence: latestSiteControl.observation?.observation_id ?? latestSiteControl.request.control_request_id,
    }] : []),
    ...(latestFacilityProject ? [{
      id: latestFacilityProject.plan.plan_id,
      location: "headquarters" as const,
      title: latestFacilityProject.project?.project_id ? "设施项目与 WBS 基线" : "设施项目基线草案",
      state: latestFacilityProject.project?.project_id ? "已激活" : latestFacilityProject.approval_status || latestFacilityProject.status,
      summary: `${latestFacilityProject.plan.project_name} · ${latestFacilityProject.plan.wbs_items?.length ?? 0} 个工作包`,
      evidence: latestFacilityProject.project?.project_id ?? latestFacilityProject.plan_hash,
    }] : []),
    ...(latestContractAward?.rfq?.rfq_id ? [{
      id: latestContractAward.rfq.rfq_id,
      location: "contractor" as const,
      title: `工程采购邀请 · ${latestContractAward.rfq.package_name}`,
      state: latestContractAward.status,
      summary: `${latestContractAward.rfq.sourcing_strategy} · 上限 ${cny(latestContractAward.rfq.contract_ceiling.value)}`,
      evidence: `contractor.rfq.issue:${latestContractAward.rfq.rfq_id}`,
    }] : []),
    ...(latestContractAward?.observation?.observation_id ? [{
      id: latestContractAward.observation.observation_id,
      location: "contractor" as const,
      title: "承包商可信投标",
      state: "World Observation",
      summary: `${latestContractAward.observation.bids.length} 份密封报价与资质事实`,
      evidence: latestContractAward.observation.observation_id,
    }] : []),
    ...(latestContractAward?.recommendation?.recommendation_id ? [{
      id: latestContractAward.recommendation.recommendation_id,
      location: "boardroom" as const,
      title: "合同授予推荐与审批",
      state: latestContractAward.approval_status,
      summary: `推荐投标 ${latestContractAward.recommendation.selected_bid_id}`,
      evidence: latestContractAward.approval_request_id,
    }] : []),
    ...(latestContractAward?.contract?.contract_id ? [{
      id: latestContractAward.contract.contract_id,
      location: "headquarters" as const,
      title: "正式工程合同",
      state: latestContractAward.contract.status ?? "active",
      summary: `${latestContractAward.contract.contractor_name} · ${latestContractAward.contract.committed_amount ? cny(latestContractAward.contract.committed_amount.value) : "已承诺"}`,
      evidence: latestContractAward.contract.contract_id,
    }] : []),
  ];
  useEffect(() => {
    if (!eligibleAssessments.some((item) => item.proposal_id === selectedProposalID)) {
      setSelectedProposalID(eligibleAssessments[0]?.proposal_id ?? "");
    }
  }, [eligibleAssessments, selectedProposalID]);

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
      setActiveRequirement(requirement);
      setProposalSet(result.proposal_set);
      setReviews({});
      setSavedReviews({});
      setNextRevision((current) => current + 1);
      setTaskView("current");
    } catch (reason) {
      setError(reason instanceof Error ? reason.message : String(reason));
	  try {
		const committed = await loadPlantRequirement(requirementID);
		if (committed?.investment_request && committed.target_available_at) {
		  setActiveRequirement(committed);
		  setNextRevision(committed.revision + 1);
		}
	  } catch {
		// The original error remains authoritative; this read only recovers a requirement committed before model failure.
	  }
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

  const addManual = async (event: FormEvent) => {
    event.preventDefault();
    if (!activeRequirement) {
      setError("请先保存设施需求；系统随后可由人工建立首版候选集");
      return;
    }
    setManualBusy(true); setError("");
    try {
      const money = (value: string) => ({ value, currency: "CNY" as const, scale: 2 as const });
      const available = localDate(manualAvailableAt);
      const proposal: Omit<SiteOptionProposal, "proposal_id" | "status"> = {
      option_type: manualOptionType,
      display_name: manualName.trim(),
      business_rationale: manualRationale.trim(),
      estimated_amount: { minimum: money(manualMinimum), likely: money(manualLikely), maximum: money(manualMaximum), basis: manualBasis.trim() },
      estimated_schedule: { earliest: available, likely: available, latest: available },
      assumptions: manualAssumptions.split("\n").map((value) => value.trim()).filter(Boolean),
      facts_required: manualFacts.split("\n").map((value) => value.trim()).filter(Boolean),
      risks: manualRisks.split("\n").map((value) => value.trim()).filter(Boolean),
      source_refs: ["manual:user-input"], confidence: "0.40",
    };
      const result = await submitManualPlantProposal({
        requirement_id: activeRequirement.requirement_id,
        proposal_set_id: proposalSet?.proposal_set_id ?? "",
        expected_revision: proposalSet?.revision ?? 0,
        proposal,
      });
      setProposalSet(result.proposal_set);
      setReviews({}); setSavedReviews({});
      setManualName(""); setManualRationale(""); setManualOptionType("");
      setManualMinimum(""); setManualLikely(""); setManualMaximum(""); setManualBasis("");
      setManualAvailableAt(""); setManualAssumptions(""); setManualFacts(""); setManualRisks(""); setManualOpen(false);
    } catch (reason) {
      setError(reason instanceof Error ? reason.message : String(reason));
    } finally {
      setManualBusy(false);
    }
  };

  const startInvestigation = async (proposalID: string) => {
    if (!proposalSet) return;
    setInvestigationBusy(proposalID); setError("");
    try {
      const now = new Date().toISOString();
      await requestSiteInvestigation({
        schema_version: "1.0", investigation_request_id: createClientRequestId("site-investigation"),
        case_code: caseCode, proposal_set_id: proposalSet.proposal_set_id, proposal_id: proposalID,
        expected_revision: proposalSet.revision, world_run_id: `plant-build-${caseCode}`,
        scope: ["ownership", "commercial_quote", "available_area", "electricity_capacity", "available_date", "permit"],
        requested_by: "resolved-by-aese", requested_at: now, status: "waiting_world",
      });
      await refreshInvestigations();
    } catch (reason) { setError(reason instanceof Error ? reason.message : String(reason)); } finally { setInvestigationBusy(""); }
  };

  const submitRecommendation = async (event: FormEvent) => {
    event.preventDefault();
    if (!proposalSet || !selectedProposalID) return;
    setSelectionBusy(true); setError("");
    try {
      await submitSiteSelectionRecommendation({
        schema_version: "1.0", recommendation_id: createClientRequestId("site-recommendation"),
        case_code: caseCode, proposal_set_id: proposalSet.proposal_set_id, proposal_set_revision: proposalSet.revision,
        selected_proposal_id: selectedProposalID, assessment_policy_version: "site-assessment-v1",
        weights: assessmentWeights, recommendation_reason: recommendationReason.trim(),
        alternative_comparison: alternativeComparison.trim(),
        ...(eligibleAssessments.length < 2 ? { single_source_exception_reason: singleSourceReason.trim() } : {}),
        recommended_at: new Date().toISOString(),
      });
      await refreshSiteSelections();
      setRecommendationReason(""); setAlternativeComparison(""); setSingleSourceReason("");
    } catch (reason) { setError(reason instanceof Error ? reason.message : String(reason)); } finally { setSelectionBusy(false); }
  };

  const formalizeSelection = async (item: SiteSelectionItem) => {
    setSelectionBusy(true); setError("");
    try {
      await finalizeSiteSelection(caseCode, item.recommendation.recommendation_id, item.recommendation.approval_request_id);
      await refreshSiteSelections();
    } catch (reason) { setError(reason instanceof Error ? reason.message : String(reason)); } finally { setSelectionBusy(false); }
  };

  const startSiteControl = async (event: FormEvent) => {
    event.preventDefault();
    if (!latestSelection?.decision || !controlHandoverAt) return;
    setControlBusy(true); setError("");
    try {
      const requestID = createClientRequestId("site-control");
      await requestSiteControl({
        schema_version: "1.0",
        control_request_id: requestID,
        selection_id: latestSelection.decision.selection_id,
        case_code: caseCode,
        selected_proposal_id: latestSelection.decision.selected_proposal_id,
        world_run_id: `world-run-m10-${caseCode}`,
        agreement_mode: controlMode,
        requested_handover_at: localDate(controlHandoverAt),
        required_evidence: ["executed_agreement", "handover_record", "possession_authority"],
        requested_by: "resolved-by-server",
        requested_at: new Date().toISOString(),
        status: "waiting_world",
      });
      await refreshSiteControls();
    } catch (reason) { setError(reason instanceof Error ? reason.message : String(reason)); } finally { setControlBusy(false); }
  };

  const confirmSiteDelivery = async () => {
    if (!latestSiteControl) return;
    setControlBusy(true); setError("");
    try {
      await confirmSiteControlDelivery(caseCode, latestSiteControl.request.control_request_id);
      await refreshSiteControls();
    } catch (reason) { setError(reason instanceof Error ? reason.message : String(reason)); } finally { setControlBusy(false); }
  };

  const prepareProjectOptions = async () => {
    setProjectBusy(true); setError("");
    try {
      const options = await generateFacilityProjectOptions(caseCode);
      setProjectOptions(options); setSelectedProjectOption(null);
    } catch (reason) { setError(reason instanceof Error ? reason.message : String(reason)); }
    finally { setProjectBusy(false); }
  };

  const submitProjectOption = async () => {
    if (!selectedProjectOption || !projectOptions) return;
    setProjectBusy(true); setError("");
    try {
      await submitFacilityProjectOption(caseCode, selectedProjectOption, projectOptions.evidence);
      await refreshFacilityProjects();
      setProjectOptions(null); setSelectedProjectOption(null);
    } catch (reason) { setError(reason instanceof Error ? reason.message : String(reason)); }
    finally { setProjectBusy(false); }
  };

  const decideProjectApproval = async (decision: "approve" | "reject") => {
    if (!latestFacilityProject?.approval_request_id || projectApprovalNote.trim().length < 6) return;
    setProjectBusy(true); setError("");
    try {
      await decidePlantApproval(latestFacilityProject.approval_request_id, decision, projectApprovalNote.trim());
      await refreshFacilityProjects();
      setProjectApprovalDetail(await loadPlantApprovalDetail(latestFacilityProject.approval_request_id));
      setProjectApprovalNote("");
    } catch (reason) { setError(reason instanceof Error ? reason.message : String(reason)); }
    finally { setProjectBusy(false); }
  };

  const activateProjectBaseline = async () => {
    if (!latestFacilityProject?.approval_request_id) return;
    setProjectBusy(true); setError("");
    try {
      await activateFacilityProject(caseCode, latestFacilityProject.plan.plan_id, latestFacilityProject.approval_request_id);
      await refreshFacilityProjects();
    } catch (reason) { setError(reason instanceof Error ? reason.message : String(reason)); }
    finally { setProjectBusy(false); }
  };

  const resumeProjectSubmission = async () => {
    if (!latestFacilityProject?.plan.plan_id || !latestFacilityProject.plan_hash) return;
    setProjectBusy(true); setError("");
    try {
      await submitFacilityProjectBaseline(caseCode, latestFacilityProject.plan.plan_id, latestFacilityProject.plan_hash);
      await refreshFacilityProjects();
    } catch (reason) { setError(reason instanceof Error ? reason.message : String(reason)); }
    finally { setProjectBusy(false); }
  };

  const issueRFQ = async () => {
    if (!selectedPackageCode) return;
    setContractBusy(true); setError("");
    try {
      await issueContractRFQ(caseCode, selectedPackageCode, sourcingStrategy);
      await refreshContractAwards();
    } catch (reason) { setError(reason instanceof Error ? reason.message : String(reason)); }
    finally { setContractBusy(false); }
  };

  const receiveContractBids = async () => {
    if (!latestContractAward?.rfq.rfq_id) return;
    setContractBusy(true); setError("");
    try {
      await confirmContractBids(caseCode, latestContractAward.rfq.rfq_id);
      await refreshContractAwards();
    } catch (reason) { setError(reason instanceof Error ? reason.message : String(reason)); }
    finally { setContractBusy(false); }
  };

  const askContractAgent = async () => {
    if (!latestContractAward?.rfq.rfq_id) return;
    setContractBusy(true); setError("");
    try { setContractAdvice(await generateContractRecommendation(caseCode, latestContractAward.rfq.rfq_id)); }
    catch (reason) { setError(reason instanceof Error ? reason.message : String(reason)); }
    finally { setContractBusy(false); }
  };

  const confirmContractRecommendation = async () => {
    if (!latestContractAward?.rfq.rfq_id || !contractAdvice) return;
    setContractBusy(true); setError("");
    try {
      await submitContractRecommendation(caseCode, latestContractAward.rfq.rfq_id, contractAdvice);
      setContractAdvice(null); await refreshContractAwards();
    } catch (reason) { setError(reason instanceof Error ? reason.message : String(reason)); }
    finally { setContractBusy(false); }
  };

  const decideContractApproval = async (decision: "approve" | "reject") => {
    if (!latestContractAward?.approval_request_id) return;
    setContractBusy(true); setError("");
    try {
      const note = decision === "approve" ? "已核对可信投标、Agent 比选依据和合同承诺金额，同意授予。" : "当前投标或合同条件不满足项目治理要求，退回重新寻源。";
      await decidePlantApproval(latestContractAward.approval_request_id, decision, note);
      await refreshContractAwards();
      setContractApprovalDetail(await loadPlantApprovalDetail(latestContractAward.approval_request_id));
    } catch (reason) { setError(reason instanceof Error ? reason.message : String(reason)); }
    finally { setContractBusy(false); }
  };

  const finalizeContractAward = async () => {
    const recommendationID = latestContractAward?.recommendation?.recommendation_id;
    if (!recommendationID || !latestContractAward?.approval_request_id) return;
    setContractBusy(true); setError("");
    try {
      await awardContract(caseCode, recommendationID, latestContractAward.approval_request_id);
      await refreshContractAwards();
    } catch (reason) { setError(reason instanceof Error ? reason.message : String(reason)); }
    finally { setContractBusy(false); }
  };

  const decideApproval = async (decision: "approve" | "reject") => {
    if (!latestSelection || approvalNote.trim().length < 6) return;
    setSelectionBusy(true); setError("");
    try {
      await decidePlantApproval(latestSelection.recommendation.approval_request_id, decision, approvalNote.trim());
      await refreshSiteSelections();
      setApprovalDetail(await loadPlantApprovalDetail(latestSelection.recommendation.approval_request_id));
      setApprovalNote("");
    } catch (reason) { setError(reason instanceof Error ? reason.message : String(reason)); } finally { setSelectionBusy(false); }
  };

  const approvalURL = latestSelection && typeof window !== "undefined"
    ? `http://${window.location.hostname}:3000/#approvals?request=${encodeURIComponent(latestSelection.recommendation.approval_request_id)}` : "#";

  return (
    <section className="plant-planning-workspace plant-planning-game-mode">
      <PlantBuildGameScene
        stage={gameStage}
        caseCode={caseCode}
        companyCode={entityCode}
        proposalNames={proposalSet?.proposals.map((proposal) => proposal.display_name) ?? []}
        observedCount={observedCount}
        archives={archives}
        onExit={onExit}
        onOpenTask={() => { setTaskView("current"); setTaskOpen(true); }}
      />
      {taskOpen && <div className="plant-task-scrim" role="presentation" onMouseDown={() => setTaskOpen(false)}>
        <section className="plant-task-dialog" role="dialog" aria-modal="true" aria-label="当前经营任务" data-stage={taskView === "requirement-revision" ? "requirement-revision" : gameStage.key} data-has-proposals={Boolean(proposalSet)} onMouseDown={(event) => event.stopPropagation()}>
          <header className="plant-task-dialog-header"><div><span>{taskView === "requirement-revision" ? "需求版本治理" : `${gameStage.npcRole} · ${gameStage.npc}`}</span><strong>{taskView === "requirement-revision" ? `修订设施需求 · 将保存为第 ${nextRevision} 版` : gameStage.actionLabel}</strong></div><nav>{taskView === "requirement-revision" && <button type="button" aria-label="返回外部事实比较" onClick={() => setTaskView("current")}><ArrowLeft /></button>}<button type="button" aria-label="关闭当前任务" onClick={() => setTaskOpen(false)}><X /></button></nav></header>
          <div className="plant-task-dialog-body">
      <header className="plant-planning-heading">
        <div><span>M10 INTERACTIVE PLANNING</span><h2 id="plant-planning-title">设施需求与候选方案</h2><p>先由人员定义约束和可调整金额，再让设施规划 Agent 提出候选；人员负责调研、选择和治理决定。</p></div>
        <div className={`plant-model-state ${status?.state ?? "loading"}`}><Bot /> <span>规划模型</span><strong>{status?.state === "connected" ? `${status.provider} · ${status.model}` : status?.state === "not_configured" ? "未启用外部模型" : "正在检查"}</strong></div>
      </header>
      <details className="plant-feature-help">
        <summary><CircleHelp />功能说明：这一步解决什么问题？</summary>
        <div><p>业务目的：把工厂需求、投资边界和工期转成可审阅的场址候选，不让系统用固定答案代替管理判断。</p><p>关系：M9 资格/资金与预算 → 设施需求 → Agent 候选 → 人工调研与选择 → IAOS 投资审批/合同/WBS → World 验收 → M11 能力建设资格。</p><p>Agent 只生成建议和待验证事实，不能批准投资、确认权属或伪造外部证明。可用现金和已批预算来自 IAOS 权威账务，只读且不能在本页面伪造。</p></div>
      </details>
      {!activeRequirement && taskView === "current" && <section id="plant-task-requirement" className="plant-requirement-adviser">
        <header><Bot /><div><strong>先让设施规划 Agent 准备需求草案</strong><p>Agent 会读取本企业的已批预算和可用现金，提出不同经营取向；不会伪造场址、报价或许可。</p></div></header>
        {!requirementOptions && <button type="button" onClick={prepareRequirementOptions} disabled={requirementAdviserBusy || status?.state !== "connected" || !financial}>{requirementAdviserBusy ? <LoaderCircle className="gx-spin" /> : <Bot />}{requirementAdviserBusy ? "Agent 正在分析企业边界…" : "让 Agent 准备需求方案"}</button>}
        {financial && <div className="plant-authority-snapshot"><span>IAOS 可用现金 {cny(financial.financial_constraint.available_cash.value)}</span><span>已批预算 {cny(financial.financial_constraint.approved_budget.value)}</span></div>}
        {requirementOptions && <div className="plant-requirement-option-grid">{requirementOptions.options.map((option) => <button type="button" key={option.option_id} className={selectedRequirementOption?.option_id === option.option_id ? "selected" : ""} onClick={() => chooseRequirementOption(option)}><span>{option.title}</span><strong>{option.target_region}</strong><p>{option.business_rationale}</p><small>{cny(option.investment_request.value)} · {new Date(option.target_available_at).toLocaleDateString("zh-CN")}</small><em>{option.tradeoffs.join("；")}</em></button>)}</div>}
        {requirementOptions && <small className="plant-agent-evidence">{requirementOptions.evidence.provider} / {requirementOptions.evidence.model} · {requirementOptions.evidence.prompt_version}</small>}
      </section>}
      {(taskView === "requirement-revision" || selectedRequirementOption) && <form id="plant-task-requirement-form" className="plant-requirement-form plant-requirement-confirm" onSubmit={submit}>
        {selectedRequirementOption && taskView === "current" && <section className="plant-selected-requirement"><span>已选择 Agent 草案</span><strong>{selectedRequirementOption.title}</strong><p>{draft.targetRegion} · {draft.facilityPurpose}</p></section>}
        <fieldset className="plant-management-boundaries"><legend>{taskView === "requirement-revision" ? "确认修订的经营边界" : "只需确认两项经营边界"}</legend>
          <label>本次投资申请金额上限（CNY）<small>Agent 建议值，可按本企业情况调整</small><input required min="0" step="0.01" type="number" value={draft.investmentRequest} onChange={(e) => update("investmentRequest", e.target.value)} /></label>
          <label>目标可用时间<small>Agent 建议日期，可由项目负责人调整</small><input required type="datetime-local" value={draft.targetAvailableAt} onChange={(e) => update("targetAvailableAt", e.target.value)} /></label>
        </fieldset>
        <details className="plant-professional-parameters" open={taskView === "requirement-revision"}><summary>查看或调整专业参数</summary><div>
          <label>目标区域<input required value={draft.targetRegion} onChange={(e) => update("targetRegion", e.target.value)} /></label>
          <label>设施用途<textarea required value={draft.facilityPurpose} onChange={(e) => update("facilityPurpose", e.target.value)} /></label>
          <label>最小面积（m²）<input required min="1" type="number" value={draft.minimumAreaM2} onChange={(e) => update("minimumAreaM2", e.target.value)} /></label>
          <label>最小电力容量（kVA）<input required min="1" type="number" value={draft.minimumElectricKVA} onChange={(e) => update("minimumElectricKVA", e.target.value)} /></label>
          <label>候选数量<input required min="2" max="8" type="number" value={draft.candidateCount} onChange={(e) => update("candidateCount", e.target.value)} /></label>
          <label>最低现金保留额（CNY）<input required min="0" step="0.01" type="number" value={draft.minimumCashReserve} onChange={(e) => update("minimumCashReserve", e.target.value)} /></label>
          <div className="plant-option-types"><strong>允许的方案类型</strong>{OPTION_TYPES.map(([code, label]) => <label key={code}><input type="checkbox" checked={draft.optionTypes.includes(code)} onChange={(e) => update("optionTypes", e.target.checked ? [...draft.optionTypes, code] : draft.optionTypes.filter((value) => value !== code))} />{label}</label>)}</div>
          <label>业务偏好<textarea value={draft.preferences} onChange={(e) => update("preferences", e.target.value)} /></label>
          <label>本次修订原因<input required value={draft.revisionReason} onChange={(e) => update("revisionReason", e.target.value)} /></label>
          <div className="plant-authority-snapshot"><span>IAOS 可用现金 {financial ? cny(financial.financial_constraint.available_cash.value) : "读取中"}</span><span>已批预算 {financial ? cny(financial.financial_constraint.approved_budget.value) : "读取中"}</span></div>
        </div></details>
        <div className="plant-generation-actions"><button disabled={busy || status?.state !== "connected" || !financial || draft.optionTypes.length === 0}>{busy ? <LoaderCircle className="gx-spin" /> : <BadgeCheck />}{busy ? "Agent 正在生成场址候选…" : taskView === "requirement-revision" ? "保存修订并让 Agent 重新生成候选" : "确认草案并生成场址候选"}</button><small>确认后通过 IAOS Capability 保存需求；Agent 只提出候选，不会自动批准投资。</small></div>
      </form>}
      {status?.state === "not_configured" && <p className="plant-inline-warning" role="status"><TriangleAlert />外部模型未启用，不能生成虚拟固定候选；请配置模型，或使用“人工新增候选”。</p>}
      {error && <p className="plant-inline-error" role="alert">{error}</p>}
      {activeRequirement && !proposalSet && <button type="button" className="plant-manual-takeover" onClick={() => setManualOpen((value) => !value)}><FilePlus2 />人工新增候选</button>}
      {manualOpen && <form className="plant-manual-form" onSubmit={addManual}><h3>人工新增权威候选</h3><p>{proposalSet ? `将在候选集第 ${proposalSet.revision} 版后追加一个不可覆盖的候选。` : "当前没有候选集；将由项目负责人建立第 1 版人工候选。"}</p><label>候选名称<input required value={manualName} onChange={(e) => setManualName(e.target.value)} /></label><label>方案类型<select required value={manualOptionType} onChange={(e) => setManualOptionType(e.target.value)}><option value="">请选择</option>{(activeRequirement?.allowed_option_types ?? draft.optionTypes).map((value) => <option key={value} value={value}>{OPTION_TYPES.find(([code]) => code === value)?.[1] ?? value}</option>)}</select></label><label>业务理由<textarea required minLength={6} value={manualRationale} onChange={(e) => setManualRationale(e.target.value)} /></label><label>最小估算（CNY）<input required min="0" step="0.01" type="number" value={manualMinimum} onChange={(e) => setManualMinimum(e.target.value)} /></label><label>最可能估算（CNY）<input required min="0" step="0.01" type="number" value={manualLikely} onChange={(e) => setManualLikely(e.target.value)} /></label><label>最大估算（CNY）<input required min="0" step="0.01" type="number" value={manualMaximum} onChange={(e) => setManualMaximum(e.target.value)} /></label><label>估算依据<textarea required value={manualBasis} onChange={(e) => setManualBasis(e.target.value)} /></label><label>预计可用日期<input required type="date" value={manualAvailableAt} onChange={(e) => setManualAvailableAt(e.target.value)} /></label><label>假设 <small>每行一项</small><textarea required value={manualAssumptions} onChange={(e) => setManualAssumptions(e.target.value)} /></label><label>待核验事实 <small>每行一项</small><textarea required value={manualFacts} onChange={(e) => setManualFacts(e.target.value)} /></label><label>主要风险 <small>每行一项</small><textarea required value={manualRisks} onChange={(e) => setManualRisks(e.target.value)} /></label><button disabled={manualBusy || !activeRequirement}>{manualBusy ? <LoaderCircle className="gx-spin" /> : <FilePlus2 />}{manualBusy ? "正在提交…" : "提交到 IAOS 候选集"}</button></form>}
      {proposalSet && <section id="plant-task-proposals" className="plant-proposals" aria-live="polite"><header><div><span>{proposalSet.status}</span><h2>候选方案 · {proposalSet.proposals.length}</h2></div><small>{proposalSet.evidence.provider} / {proposalSet.evidence.model} · {proposalSet.evidence.prompt_version}</small></header>{historicalInvestigationCount > 0 && <p className="plant-inline-warning" role="status"><TriangleAlert />需求或候选集已修订，{historicalInvestigationCount} 条旧调研事实只保留在“候选场址 → 场景档案”，不会参与当前比较或正式推荐。请审阅当前候选并重新发起调研。</p>}<div className="plant-proposal-grid">{proposalSet.proposals.map((proposal) => <ProposalCard key={proposal.proposal_id} proposal={proposal} review={reviews[proposal.proposal_id]} onReview={(review) => { setReviews((current) => ({ ...current, [proposal.proposal_id]: review })); setSavedReviews((current) => ({ ...current, [proposal.proposal_id]: false })); }} onSubmitReview={() => submitReview(proposal.proposal_id)} saved={Boolean(savedReviews[proposal.proposal_id])} busy={reviewBusy === proposal.proposal_id} investigation={currentInvestigations.find((item) => item.request.proposal_id === proposal.proposal_id)} investigationBusy={investigationBusy === proposal.proposal_id} onRequestInvestigation={() => startInvestigation(proposal.proposal_id)} />)}</div><p className="plant-persistence-note"><BadgeCheck />Agent 候选、人工新增候选和相应审阅都通过 IAOS Capability 保存并形成版本；“采纳调研”不等于投资批准或合同签署。</p></section>}
      {currentInvestigations.length > 0 && activeRequirement && <section id="plant-task-investigation" className="plant-investigations" aria-live="polite"><header><div><span>facility.site.investigation.v1</span><h2>场址外部调研工作项</h2></div><small>{currentInvestigations.filter((item) => item.status === "waiting_world").length} 条等待 World</small></header>{currentInvestigations.map((item) => <InvestigationPanel key={item.request.investigation_request_id} item={item} caseCode={caseCode} requirementID={activeRequirement.requirement_id} onCommitted={refreshInvestigations} />)}</section>}
      {assessments.length > 0 && <section id="plant-task-comparison" className="plant-assessments" aria-labelledby="plant-assessment-title">
        <header><div><span>OBSERVATION-ONLY COMPARISON</span><h2 id="plant-assessment-title">外部事实比较</h2><p>第一步逐项比较需求门槛和现场事实；全部硬约束通过后，第二步才按经营偏好排序。Agent 估算不参与硬约束判定。</p></div><small>{eligibleAssessments.length}/{assessments.length} 个候选通过硬约束</small></header>
        <div className={`plant-gate-summary ${eligibleAssessments.length ? "passed" : "blocked"}`}><span>{eligibleAssessments.length ? <BadgeCheck /> : <XCircle />}</span><div><strong>{eligibleAssessments.length ? "已有候选进入排序" : "当前没有候选可以进入排序"}</strong><p>{eligibleAssessments.length ? "下面的权重只用于比较这些合格候选。" : "请先根据红色差额补充场址或修订需求；无论怎样调整权重，都不能抵消硬约束失败。"}</p>{!eligibleAssessments.length && activeRequirement && <button type="button" className="plant-revise-requirement" onClick={openRequirementRevision}>修订设施需求</button>}</div></div>
        <div className="plant-assessment-grid">{assessments.map((assessment) => <AssessmentCard key={assessment.observation_id} assessment={assessment} />)}</div>
        <fieldset className="plant-assessment-weights" disabled={eligibleAssessments.length === 0}><legend>第二步 · 调整合格候选的排序偏好</legend><header><div><strong>当前权重来源：界面默认比较偏好</strong><small>默认成本 35%、工期 25%、容量 20%、权属与许可 20%；不是 Agent 生成，也不是现场事实。</small></div><button type="button" onClick={() => setAssessmentWeights(DEFAULT_SITE_ASSESSMENT_WEIGHTS)}><RotateCcw />恢复默认</button></header>{SCORE_DIMENSIONS.map(([key, label, help]) => { const effective = normalizedWeightTotal ? assessmentWeights[key] / normalizedWeightTotal * 100 : 0; return <label key={key}><span>{label}<b>{effective.toFixed(0)}%</b></span><small>{help}</small><input aria-label={`${label}权重`} type="number" min="0" max="100" value={assessmentWeights[key]} onChange={(event) => setAssessmentWeights((current) => ({ ...current, [key]: Math.max(0, Number(event.target.value) || 0) }))} /></label>; })}<p>{eligibleAssessments.length ? `输入值合计 ${normalizedWeightTotal}，系统会自动换算为百分比。权重只影响合格候选的排序预览，提交推荐时 IAOS 按 site-assessment-v1 重新计算。` : "当前没有合格候选，因此排序偏好暂不生效。"}</p></fieldset>
        <details className="plant-score-method"><summary><CircleHelp />这些分数具体怎样计算？</summary><div><p><strong>成本：</strong>投资申请额中尚未被正式报价占用的比例，越节省分数越高。</p><p><strong>工期：</strong>在目标日期当天可用为 50 分；每提前一天加分，提前 180 天达到 100 分。</p><p><strong>容量：</strong>分别计算面积和电力相对最低门槛的裕量，再取平均。</p><p><strong>权属与许可：</strong>已核验/满足为 100 分，有条件结论为 60 分。</p><p>综合分 = 四项分数 × 自动归一化后的权重。该公式来自版本化策略 <code>site-assessment-v1</code>，提交推荐时 IAOS 会用权威数据重新计算。</p></div></details>
        <p className="plant-assessment-governance"><TriangleAlert />页面分数只是预览，不是正式推荐、选址批准或投资批准。提交推荐时 IAOS 会从权威需求、候选版本和 Observation 重新计算；浏览器不能指定审批人或篡改批准结果。</p>
        {!latestSelection?.decision && <form className="plant-selection-form" onSubmit={submitRecommendation}>
          <fieldset><legend>提交人工场址推荐</legend>
            <div className="plant-selection-options" role="radiogroup" aria-label="选择合格候选">{eligibleAssessments.map((assessment) => <label key={assessment.proposal_id} className={selectedProposalID === assessment.proposal_id ? "selected" : ""}><input type="radio" name="selected-site" value={assessment.proposal_id} checked={selectedProposalID === assessment.proposal_id} onChange={() => setSelectedProposalID(assessment.proposal_id)} /><span>{assessment.display_name}</span><strong>{assessment.total_score?.toFixed(1)} 分</strong></label>)}</div>
            <label>推荐理由 <small>说明为什么该候选最符合经营目标；至少 12 个字符</small><textarea required minLength={12} value={recommendationReason} onChange={(event) => setRecommendationReason(event.target.value)} /></label>
            <label>替代方案比较 <small>明确与其他候选的差异和未选原因；至少 12 个字符</small><textarea required minLength={12} value={alternativeComparison} onChange={(event) => setAlternativeComparison(event.target.value)} /></label>
            {eligibleAssessments.length < 2 && <label>单一来源例外说明 <small>当前少于两个合格候选，必须解释为何仍可进入审批；至少 20 个字符</small><textarea required minLength={20} value={singleSourceReason} onChange={(event) => setSingleSourceReason(event.target.value)} /></label>}
            <button disabled={selectionBusy || !selectedProposalID}>{selectionBusy ? <LoaderCircle className="gx-spin" /> : <BadgeCheck />}{selectionBusy ? "正在提交…" : "提交 IAOS 选址审批"}</button>
          </fieldset>
        </form>}
      </section>}
      {latestSelection && <PlantApprovalPanel selection={latestSelection} detail={approvalDetail} loadError={approvalLoadError} note={approvalNote} busy={selectionBusy} approvalURL={approvalURL} onNote={setApprovalNote} onDecide={decideApproval} onFormalize={() => formalizeSelection(latestSelection)} />}
      {latestSelection?.decision && <section id="plant-task-site-control" className="plant-site-control">
        <header><div><span>facility.plant.delivery.v1</span><h2>场地控制与实际交付</h2><p>正式选址只证明企业批准了哪个场址；必须取得已签协议、交付记录和占有权限，才能建立可执行设施项目。</p></div><strong>{latestSiteControl?.status ?? "尚未发起"}</strong></header>
        {!latestSiteControl && <form onSubmit={startSiteControl}>
          <fieldset><legend>由工厂项目负责人发起交付请求</legend>
            <label>场地取得方式<select value={controlMode} onChange={(event) => setControlMode(event.target.value as SiteControlRequest["agreement_mode"])}><option value="lease">租赁</option><option value="purchase">购买</option><option value="build_to_suit">定制代建</option><option value="use_agreement">场地使用协议</option></select></label>
            <label>期望交付时间<input required type="datetime-local" value={controlHandoverAt} onChange={(event) => setControlHandoverAt(event.target.value)} /></label>
            <div className="plant-control-evidence"><strong>系统强制核验</strong><span>已签协议</span><span>交付记录</span><span>占有与使用权限</span></div>
            <button disabled={controlBusy || !controlHandoverAt}>{controlBusy ? <LoaderCircle className="gx-spin" /> : <SearchCheck />}{controlBusy ? "正在发起…" : "发起场地控制交付工作项"}</button>
          </fieldset>
        </form>}
        {latestSiteControl?.status === "waiting_world" && <section className="plant-control-confirmation" aria-labelledby="plant-control-confirmation-title">
          <div className="plant-control-world-badge"><BadgeCheck /><span><small>WORLD ENGINE · 园区权利方</small><strong id="plant-control-confirmation-title">场地已经准备交付</strong></span></div>
          <p>协议、交接单、生效时间和占有权限由 World 引擎依据权威交付请求生成。你只需核对本次交付并确认接收，不需要填写系统编号或 JSON。</p>
          <dl>
            <div><dt>取得方式</dt><dd>{{ lease: "租赁", purchase: "购买", build_to_suit: "定制代建", use_agreement: "场地使用协议" }[latestSiteControl.request.agreement_mode]}</dd></div>
            <div><dt>计划交付</dt><dd>{new Date(latestSiteControl.request.requested_handover_at).toLocaleString("zh-CN", { timeZone: "Asia/Shanghai", hour12: false })}</dd></div>
            <div><dt>引擎将归档</dt><dd>已签协议 · 交接记录 · 占有与使用权限</dd></div>
          </dl>
          <button type="button" onClick={confirmSiteDelivery} disabled={controlBusy}>{controlBusy ? <LoaderCircle className="gx-spin" /> : <BadgeCheck />}{controlBusy ? "World 正在生成交付事实…" : "确认接收场地"}</button>
          <small>确认后：World Journal 记录外部事实 → IAOS Capability 核验 → 关闭等待工作项。玩家不能修改这些权威证据。</small>
        </section>}
        {latestSiteControl?.status === "controlled" && <div className="plant-control-complete"><BadgeCheck /><div><strong>场址取得章节已完成</strong><p>{latestSiteControl.observation?.agreement_ref} · {latestSiteControl.observation?.handover_ref}</p><small>协议、交接和占有权限已归档。请按剧情引导返回总部项目办公室，让设施项目 Agent 准备下一章。</small></div></div>}
      </section>}
      {latestSiteControl?.status === "controlled" && <section id="plant-task-project" className="plant-project-baseline">
        <header><div><span>facility.project.baseline.v1</span><h2>设施项目与 WBS 基线</h2><p>Agent 负责专业分解，玩家只选择经营方案，组织审批后再激活权威基线。</p></div><strong>{latestFacilityProject?.project?.project_id ? "已激活" : latestFacilityProject?.approval_status || "等待 Agent"}</strong></header>
        {(!latestFacilityProject || projectCanRevise) && !projectOptions && <section className="plant-project-agent-start"><Bot /><div><strong>{projectCanRevise ? "审批意见已退回项目办公室" : "项目办公室已准备就绪"}</strong><p>纪元会读取设施需求、投资上限和场地交付事实，生成三套可执行项目方案。你不需要填写 WBS、岗位、证据编号或 JSON。</p></div><button type="button" disabled={projectBusy || status?.state !== "connected"} onClick={prepareProjectOptions}>{projectBusy ? <LoaderCircle className="gx-spin" /> : <Bot />}{projectBusy ? "Agent 正在编制项目方案…" : projectCanRevise ? "让 Agent 根据审批意见修订方案" : "让 Agent 准备项目方案"}</button></section>}
        {(!latestFacilityProject || projectCanRevise) && projectOptions && <section className="plant-project-option-stage">
          <header><div><strong>选择一套项目经营方案</strong><p>这里只决定交付策略、预算上限和投产目标；专业 WBS 可展开查看。</p></div><small>{projectOptions.evidence.provider} / {projectOptions.evidence.model}</small></header>
          <div className="plant-project-option-grid">{projectOptions.options.map((option) => <article key={option.option_id} className={selectedProjectOption?.option_id === option.option_id ? "selected" : ""}><button type="button" onClick={() => setSelectedProjectOption(option)}><span>{option.title}</span><strong>{option.project_name}</strong><p>{option.business_rationale}</p><dl><div><dt>策略</dt><dd>{{ design_bid_build: "设计-招标-施工", design_build: "设计施工总承包", epcm: "EPCM" }[option.delivery_strategy]}</dd></div><div><dt>预算上限</dt><dd>{cny(option.budget_ceiling.value)}</dd></div><div><dt>目标投产</dt><dd>{new Date(option.target_ready_at).toLocaleDateString("zh-CN")}</dd></div></dl><em>{option.tradeoffs.join("；")}</em></button><details><summary>查看 {option.wbs_items.length} 个专业工作包</summary><ol>{option.wbs_items.map((item) => <li key={item.wbs_code}><b>{item.wbs_code} · {item.name}</b><small>{item.phase} · {item.owner_position} · {(item.budget_share_bps / 100).toFixed(0)}%</small></li>)}</ol></details></article>)}</div>
          <button type="button" className="plant-project-confirm" disabled={!selectedProjectOption || projectBusy} onClick={submitProjectOption}>{projectBusy ? <LoaderCircle className="gx-spin" /> : <BadgeCheck />}{projectBusy ? "正在冻结草案并发起审批…" : "确认方案并提交项目基线审批"}</button>
        </section>}
        {latestFacilityProject && !latestFacilityProject.project?.project_id && <section className="plant-project-governance">
          <div className="plant-project-summary"><Gavel /><div><small>{projectApprovalDetail?.detail.flow_name ?? "genesis.facility.project.approval"}</small><strong>{latestFacilityProject.plan.project_name}</strong><p>{latestFacilityProject.plan.delivery_strategy} · {cny(latestFacilityProject.plan.budget_ceiling.value)} · {latestFacilityProject.plan.wbs_items.length} 个 WBS</p></div><em>{latestFacilityProject.approval_status}</em></div>
          {projectApprovalDetail && <div className="plant-approval-routing">{projectApprovalDetail.detail.assignments.map((assignment) => <article key={assignment.id}><Gavel /><span><strong>{assignment.stage_name}</strong><small>{assignment.display_name}</small></span><em>{assignment.status}</em></article>)}</div>}
          {latestFacilityProject.status === "draft" && !latestFacilityProject.approval_request_id && <button type="button" className="plant-project-confirm" disabled={projectBusy} onClick={resumeProjectSubmission}>{projectBusy ? <LoaderCircle className="gx-spin" /> : <Gavel />}{projectBusy ? "正在恢复审批提交…" : "草案已保存 · 继续提交项目审批"}</button>}
          {latestFacilityProject.approval_status === "pending" && projectApprovalDetail?.detail.can_decide && <div className="plant-project-decision"><label>审批意见<textarea minLength={6} value={projectApprovalNote} onChange={(event) => setProjectApprovalNote(event.target.value)} placeholder="说明批准条件或驳回原因（至少 6 个字符）" /></label><div><button type="button" disabled={projectBusy || projectApprovalNote.trim().length < 6} onClick={() => decideProjectApproval("approve")}><BadgeCheck />批准项目基线</button><button type="button" className="danger" disabled={projectBusy || projectApprovalNote.trim().length < 6} onClick={() => decideProjectApproval("reject")}><XCircle />驳回草案</button></div></div>}
          {latestFacilityProject.approval_status === "approved" && <button type="button" className="plant-project-confirm" disabled={projectBusy} onClick={activateProjectBaseline}>{projectBusy ? <LoaderCircle className="gx-spin" /> : <BadgeCheck />}{projectBusy ? "正在激活项目与 WBS…" : "批准已生效 · 激活项目基线"}</button>}
          {latestFacilityProject.approval_status === "rejected" && <p className="plant-inline-error"><XCircle />项目基线已驳回；本版保留审计，请返回项目办公室让 Agent 根据审批意见形成新草案。</p>}
        </section>}
        {latestFacilityProject?.project?.project_id && <div className="plant-control-complete"><BadgeCheck /><div><strong>设施项目与 WBS 基线已激活</strong><p>{latestFacilityProject.project.project_id} · {latestFacilityProject.project.wbs_items?.length ?? latestFacilityProject.plan.wbs_items.length} 个工作包</p><small>项目办公室档案已写入 IAOS。合同、施工、工程财务和验收将引用该稳定基线继续推进。</small></div></div>}
      </section>}
      {latestFacilityProject?.project?.project_id && <section id="plant-task-contract" className="plant-contract-award">
        <header><div><span>facility.contract.award.v1</span><h2>工程承包商寻源与合同授予</h2><p>WBS 决定采购范围；World 提供可信报价；Agent 负责比较；审批主体决定；项目负责人归档合同。</p></div><strong>{latestContractAward?.contract?.contract_id ? "已授予" : latestContractAward?.approval_status || latestContractAward?.status || "准备寻源"}</strong></header>
        {!latestContractAward && <div className="plant-contract-source"><Handshake /><div><strong>选择一个已批准 WBS 采购包</strong><p>合同上限由项目预算 × WBS 预算份额自动计算，交付日期来自 WBS；玩家不能手填金额、承包商或证据。</p></div><div className="plant-contract-package-grid">{(latestFacilityProject.project.wbs_items ?? latestFacilityProject.plan.wbs_items).map((item) => <button type="button" key={item.wbs_code} className={selectedPackageCode === item.wbs_code ? "selected" : ""} onClick={() => setSelectedPackageCode(item.wbs_code)}><span>{item.phase}</span><strong>{item.wbs_code} · {item.name}</strong><small>预算份额 {(item.budget_share_bps / 100).toFixed(0)}% · {new Date(item.planned_finish_at).toLocaleDateString("zh-CN")}</small></button>)}</div><div className="plant-contract-strategy"><button type="button" className={sourcingStrategy === "general_contract" ? "selected" : ""} onClick={() => setSourcingStrategy("general_contract")}>总承包寻源</button><button type="button" className={sourcingStrategy === "specialist_packages" ? "selected" : ""} onClick={() => setSourcingStrategy("specialist_packages")}>专业分包寻源</button><button type="button" className={sourcingStrategy === "epcm_managed" ? "selected" : ""} onClick={() => setSourcingStrategy("epcm_managed")}>EPCM 管理寻源</button></div><button type="button" className="plant-project-confirm" disabled={!selectedPackageCode || contractBusy} onClick={issueRFQ}>{contractBusy ? <LoaderCircle className="gx-spin" /> : <Handshake />}{contractBusy ? "正在发布采购邀请…" : "确认采购包并发布 RFQ"}</button></div>}
        {latestContractAward?.rfq && !latestContractAward.observation?.observation_id && <div className="plant-control-confirmation"><div className="plant-control-world-badge"><Handshake /><span><small>WORLD ENGINE · 承包商市场</small><strong>三家虚构承包商已准备密封投标</strong></span></div><p>系统会依据 RFQ 上限和交付要求生成报价、资质、质保、里程碑与证据。你只需确认收取，不填写任何外部事实。</p><dl><div><dt>采购包</dt><dd>{latestContractAward.rfq.package_name}</dd></div><div><dt>合同上限</dt><dd>{cny(latestContractAward.rfq.contract_ceiling.value)}</dd></div><div><dt>要求完成</dt><dd>{new Date(latestContractAward.rfq.required_ready_at).toLocaleDateString("zh-CN")}</dd></div></dl><button type="button" disabled={contractBusy} onClick={receiveContractBids}>{contractBusy ? <LoaderCircle className="gx-spin" /> : <Handshake />}{contractBusy ? "World 正在封存投标…" : "确认收取正式投标"}</button></div>}
        {latestContractAward?.observation?.observation_id && <div className="plant-contract-bids"><header><div><strong>可信投标已送达</strong><p>以下价格和条件来自 World Observation，不是 Agent 估算。</p></div><small>{latestContractAward.observation.bids.length} 份合格投标</small></header><div className="plant-contract-bid-grid">{latestContractAward.observation.bids.map((bid) => <article key={bid.bid_id} className={(contractAdvice?.selected_bid_id ?? latestContractAward.recommendation?.selected_bid_id) === bid.bid_id ? "selected" : ""}><span>{bid.qualification === "eligible" ? "资质通过" : bid.qualification}</span><strong>{bid.contractor_name}</strong><b>{cny(bid.quoted_amount.value)}</b><p>承诺 {new Date(bid.promised_ready_at).toLocaleDateString("zh-CN")} · 质保 {bid.warranty_months} 月 · {bid.milestone_count} 个里程碑</p></article>)}</div>{!latestContractAward.recommendation?.recommendation_id && !contractAdvice && <button type="button" className="plant-project-confirm" disabled={contractBusy || status?.state !== "connected"} onClick={askContractAgent}>{contractBusy ? <LoaderCircle className="gx-spin" /> : <Bot />}{contractBusy ? "Agent 正在比选…" : "让 Agent 评审正式投标"}</button>}{contractAdvice && <div className="plant-contract-advice"><Bot /><div><small>{contractAdvice.evidence.provider} · {contractAdvice.evidence.model}</small><strong>Agent 推荐 {latestContractAward.observation.bids.find((bid) => bid.bid_id === contractAdvice.selected_bid_id)?.contractor_name}</strong><p>{contractAdvice.recommendation_reason}</p><em>{contractAdvice.alternative_comparison}</em></div><button type="button" disabled={contractBusy} onClick={confirmContractRecommendation}><BadgeCheck />确认推荐并提交审批</button></div>}</div>}
        {latestContractAward?.recommendation?.recommendation_id && !latestContractAward.contract?.contract_id && <div className="plant-project-governance"><div className="plant-project-summary"><Gavel /><div><small>{contractApprovalDetail?.detail.flow_name ?? "genesis.facility.contract.award.approval"}</small><strong>{contractApprovalDetail?.detail.subject.title ?? latestContractAward.rfq.package_name}</strong><p>{contractApprovalDetail?.detail.subject.summary ?? latestContractAward.recommendation.recommendation_reason}</p></div><em>{latestContractAward.approval_status}</em></div>{contractApprovalDetail && <div className="plant-approval-routing">{contractApprovalDetail.detail.assignments.map((assignment) => <article key={assignment.id}><Gavel /><span><strong>{assignment.stage_name}</strong><small>{assignment.display_name}</small></span><em>{assignment.status}</em></article>)}</div>}{latestContractAward.approval_status === "pending" && contractApprovalDetail?.detail.can_decide && <div className="plant-project-decision"><p>点击即提交预置、可审计的审批意见；无需输入 JSON 或系统编号。</p><div><button type="button" disabled={contractBusy} onClick={() => decideContractApproval("approve")}><BadgeCheck />批准合同授予</button><button type="button" className="danger" disabled={contractBusy} onClick={() => decideContractApproval("reject")}><XCircle />驳回并重新寻源</button></div></div>}{latestContractAward.approval_status === "approved" && <button type="button" className="plant-project-confirm" disabled={contractBusy} onClick={finalizeContractAward}><BadgeCheck />批准已生效 · 归档正式合同</button>}{latestContractAward.approval_status === "rejected" && <p className="plant-inline-error"><XCircle />合同授予已驳回，本轮 RFQ 与投标保留审计；下一版本将重新寻源。</p>}</div>}
        {latestContractAward?.contract?.contract_id && <div className="plant-control-complete"><BadgeCheck /><div><strong>正式工程合同已归档</strong><p>{latestContractAward.contract.contractor_name} · {latestContractAward.contract.committed_amount ? cny(latestContractAward.contract.committed_amount.value) : "金额已承诺"}</p><small>{latestContractAward.contract.contract_id} · 当前只形成合同承诺，不会自动生成发票、应付或付款。</small></div></div>}
      </section>}
          </div>
        </section>
      </div>}
    </section>
  );
}

export function PlantBuildPlay({ onExit }: { onExit: () => void }) {
  return <PlantPlanningWorkspace onExit={onExit} />;
}
