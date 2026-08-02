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
  loadPlantProposalSet,
  loadPlantRequirement,
  loadSiteInvestigations,
  loadSiteSelections,
  requestSiteInvestigation,
  finalizeSiteSelection,
  submitSiteInvestigationObservation,
  submitSiteSelectionRecommendation,
  submitPlantProposalReview,
  submitManualPlantProposal,
  type FacilityRequirement,
  type PlantBuildTrace,
  type PlantPlanningProviderStatus,
  type PlantFinancialConstraint,
  type ProposalSet,
  type SiteOptionProposal,
  type SiteInvestigationItem,
  type SiteSelectionItem,
} from "../../world/plantBuild";
import {
  assessObservedSites,
  type SiteAssessment,
  type SiteAssessmentWeights,
} from "../../world/siteAssessment";
import { createClientRequestId } from "../../lib/clientRequestId";
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
        {saved && action === "adopt_for_investigation" && <button type="button" className="plant-investigation-request" disabled={investigationBusy || Boolean(investigation)} onClick={onRequestInvestigation}>
          {investigationBusy ? <LoaderCircle className="gx-spin" /> : <SearchCheck />}{investigation ? investigation.status === "observed" ? "调研事实已提交" : "正在等待外部调研" : investigationBusy ? "正在建立工作项…" : "发起外部调研工作项"}
        </button>}
      </fieldset>
    </article>
  );
}

function InvestigationPanel({ item, caseCode, onCommitted }: { item: SiteInvestigationItem; caseCode: string; onCommitted: () => void }) {
  const [ownership, setOwnership] = useState("verified");
  const [area, setArea] = useState("");
  const [electricity, setElectricity] = useState("");
  const [quote, setQuote] = useState("");
  const [availableAt, setAvailableAt] = useState("");
  const [permit, setPermit] = useState("eligible");
  const [externalActor, setExternalActor] = useState("virtual-park-operator");
  const [evidence, setEvidence] = useState("");
  const [notes, setNotes] = useState("");
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState("");
  const submitObservation = async (event: FormEvent) => {
    event.preventDefault(); setBusy(true); setError("");
    try {
      const observedAt = new Date().toISOString();
      await submitSiteInvestigationObservation(caseCode, item.request.world_run_id, {
        schema_version: "1.0", observation_id: createClientRequestId("site-observation"),
        investigation_request_id: item.request.investigation_request_id, proposal_id: item.request.proposal_id,
        result: "completed", ownership_status: ownership, available_area_m2: Number(area), electricity_kva: Number(electricity),
        quoted_amount: { value: quote, currency: "CNY", scale: 2 }, available_at: localDate(availableAt), permit_status: permit,
        evidence_refs: evidence.split("\n").map((value) => value.trim()).filter(Boolean), notes: notes.trim(),
        external_actor_id: externalActor.trim(), observed_at: observedAt,
      });
      onCommitted();
    } catch (reason) { setError(reason instanceof Error ? reason.message : String(reason)); } finally { setBusy(false); }
  };
  return <article className={`plant-investigation-card ${item.status}`}>
    <header><div><span>World wait · {item.work_item_status}</span><h3>{item.request.investigation_request_id}</h3></div><strong>{item.status === "observed" ? "可信事实已提交" : "等待园区运营方"}</strong></header>
    <p>候选 {item.request.proposal_id} · 调研范围：{item.request.scope.join("、")}</p>
    {item.observation ? <dl className="plant-investigation-facts"><div><dt>权属</dt><dd>{item.observation.ownership_status}</dd></div><div><dt>可用面积</dt><dd>{item.observation.available_area_m2} m²</dd></div><div><dt>电力</dt><dd>{item.observation.electricity_kva} kVA</dd></div><div><dt>正式报价</dt><dd>{cny(item.observation.quoted_amount.value)}</dd></div><div><dt>许可状态</dt><dd>{item.observation.permit_status}</dd></div></dl> :
      <form className="plant-observation-form" onSubmit={submitObservation}>
        <fieldset><legend>外部参与者回传实地调研事实</legend>
          <label>外部参与者标识<input required value={externalActor} onChange={(e) => setExternalActor(e.target.value)} placeholder="park-operator-001" /></label>
          <label>权属核验<select value={ownership} onChange={(e) => setOwnership(e.target.value)}><option value="verified">已核验</option><option value="conditional">有条件有效</option></select></label>
          <label>实际可用面积（m²）<input required min="1" type="number" value={area} onChange={(e) => setArea(e.target.value)} /></label>
          <label>可用电力（kVA）<input required min="1" type="number" value={electricity} onChange={(e) => setElectricity(e.target.value)} /></label>
          <label>正式报价（CNY）<input required min="0.01" step="0.01" type="number" value={quote} onChange={(e) => setQuote(e.target.value)} /></label>
          <label>最早可用时间<input required type="datetime-local" value={availableAt} onChange={(e) => setAvailableAt(e.target.value)} /></label>
          <label>许可条件<select value={permit} onChange={(e) => setPermit(e.target.value)}><option value="eligible">满足</option><option value="conditional">有条件满足</option></select></label>
          <label>证据引用（每行一项）<textarea required value={evidence} onChange={(e) => setEvidence(e.target.value)} placeholder="world-document:园区报价函-001" /></label>
          <label>补充说明<textarea value={notes} onChange={(e) => setNotes(e.target.value)} /></label>
        </fieldset>
        {error && <p role="alert" className="plant-inline-error">{error}</p>}
        <button disabled={busy}>{busy ? <LoaderCircle className="gx-spin" /> : <BadgeCheck />}{busy ? "正在校验并提交…" : "园区运营方确认并提交 Observation"}</button>
        <small>提交后先进入 IAOS World Journal，再由 `site.investigation.observation.commit` 校验并完成工作项；不能直接写业务表。</small>
      </form>}
  </article>;
}

const SCORE_DIMENSIONS: Array<[keyof SiteAssessmentWeights, string, string]> = [
  ["cost", "成本", "正式报价相对投资申请额"],
  ["schedule", "工期", "正式可用日期相对目标日期"],
  ["capacity", "容量", "实际面积与电力容量"],
  ["control", "控制", "权属与许可核验结论"],
];

function AssessmentCard({ assessment }: { assessment: SiteAssessment }) {
  return <article className={`plant-assessment-card ${assessment.eligible ? "eligible" : "ineligible"}`}>
    <header>
      <div><span>{assessment.observation_id}</span><h3>{assessment.display_name}</h3></div>
      <strong>{assessment.eligible ? `${assessment.total_score?.toFixed(1)} 分` : "硬约束不通过"}</strong>
    </header>
    {assessment.hard_failures.length > 0
      ? <ul className="plant-hard-failures">{assessment.hard_failures.map((failure) => <li key={failure}>{failure}</li>)}</ul>
      : <p className="plant-hard-pass"><BadgeCheck />投资、面积、电力、日期、权属和许可硬约束通过</p>}
    <div className="plant-estimate-observation">
      <section><span>Agent 估算 · 非正式事实</span><strong>{assessment.estimated ? cny(assessment.estimated.amount) : "当前会话未加载"}</strong><small>{assessment.estimated ? new Date(assessment.estimated.available_at).toLocaleDateString("zh-CN") : "以 Observation 为评分依据"}</small></section>
      <section><span>World Observation · 评分事实</span><strong>{cny(assessment.observed.quote)}</strong><small>{new Date(assessment.observed.available_at).toLocaleDateString("zh-CN")} · {assessment.observed.area_m2} m² · {assessment.observed.electricity_kva} kVA</small></section>
    </div>
    <dl className="plant-score-components">{SCORE_DIMENSIONS.map(([key, label]) => <div key={key}><dt>{label}</dt><dd>{assessment.component_scores[key].toFixed(1)}</dd></div>)}</dl>
    <details className="plant-assessment-evidence"><summary>查看事实与证据引用</summary><dl><div><dt>权属</dt><dd>{assessment.observed.ownership}</dd></div><div><dt>许可</dt><dd>{assessment.observed.permit}</dd></div></dl><ul>{assessment.evidence_refs.map((reference) => <li key={reference}>{reference}</li>)}</ul></details>
  </article>;
}

function PlantPlanningWorkspace() {
  const [status, setStatus] = useState<PlantPlanningProviderStatus | null>(null);
  const [financial, setFinancial] = useState<PlantFinancialConstraint | null>(null);
  const [activeRequirement, setActiveRequirement] = useState<FacilityRequirement | null>(null);
  const [draft, setDraft] = useState<RequirementDraft>(blankDraft);
  const [proposalSet, setProposalSet] = useState<ProposalSet | null>(null);
  const [reviews, setReviews] = useState<Record<string, ReviewState>>({});
  const [savedReviews, setSavedReviews] = useState<Record<string, boolean>>({});
  const [reviewBusy, setReviewBusy] = useState("");
  const [investigations, setInvestigations] = useState<SiteInvestigationItem[]>([]);
  const [investigationBusy, setInvestigationBusy] = useState("");
  const [assessmentWeights, setAssessmentWeights] = useState<SiteAssessmentWeights>({ cost: 35, schedule: 25, capacity: 20, control: 20 });
  const [siteSelections, setSiteSelections] = useState<SiteSelectionItem[]>([]);
  const [selectedProposalID, setSelectedProposalID] = useState("");
  const [recommendationReason, setRecommendationReason] = useState("");
  const [alternativeComparison, setAlternativeComparison] = useState("");
  const [singleSourceReason, setSingleSourceReason] = useState("");
  const [selectionBusy, setSelectionBusy] = useState(false);
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
  const routeParams = new URLSearchParams(window.location.hash.split("?")[1] ?? "");
  const tenant = routeParams.get("tenant") ?? localStorage.getItem("aese_iaos_tenant_id") ?? localStorage.getItem("iaos_tenant_id") ?? "";
  const caseCode = routeParams.get("case") ?? localStorage.getItem("aese_genesis_case_code") ?? "";
  const requirementID = `facility-requirement-${caseCode || "draft"}`;
  const refreshInvestigations = async () => setInvestigations((await loadSiteInvestigations(caseCode)).items ?? []);
  const refreshSiteSelections = async () => setSiteSelections((await loadSiteSelections(caseCode)).items ?? []);
  useEffect(() => {
    const controller = new AbortController();
    Promise.all([
      loadPlantPlanningStatus(controller.signal).then(setStatus),
      loadPlantFinancialConstraint(caseCode, controller.signal).then(setFinancial),
      loadPlantRequirement(requirementID, controller.signal).then((existing) => {
        setActiveRequirement(existing?.investment_request && existing.target_available_at ? existing : null);
        setNextRevision((existing?.revision ?? 0) + 1);
      }),
      loadPlantProposalSet(requirementID, controller.signal).then(setProposalSet),
      loadSiteInvestigations(caseCode, controller.signal).then((result) => setInvestigations(result.items ?? [])),
      loadSiteSelections(caseCode, controller.signal).then((result) => setSiteSelections(result.items ?? [])),
    ]).catch((reason) => { if (reason.name !== "AbortError") setError(String(reason)); });
    return () => controller.abort();
  }, [caseCode, requirementID]);
  const update = (name: keyof RequirementDraft, value: string | string[]) =>
    setDraft((current) => ({ ...current, [name]: value }));
  const entityCode = financial?.legal_entity_code ?? "";
  const assessments = useMemo(
    () => activeRequirement ? assessObservedSites(activeRequirement, proposalSet, investigations, assessmentWeights) : [],
    [activeRequirement, assessmentWeights, investigations, proposalSet],
  );
  const normalizedWeightTotal = Object.values(assessmentWeights).reduce((sum, value) => sum + Math.max(0, value), 0);
  const eligibleAssessments = assessments.filter((item) => item.eligible);
  const latestSelection = siteSelections[0];
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

  const approvalURL = latestSelection && typeof window !== "undefined"
    ? `http://${window.location.hostname}:3000/#approvals?request=${encodeURIComponent(latestSelection.recommendation.approval_request_id)}` : "#";

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
      {manualOpen && <form className="plant-manual-form" onSubmit={addManual}><h3>人工新增权威候选</h3><p>{proposalSet ? `将在候选集第 ${proposalSet.revision} 版后追加一个不可覆盖的候选。` : "当前没有候选集；将由项目负责人建立第 1 版人工候选。"}</p><label>候选名称<input required value={manualName} onChange={(e) => setManualName(e.target.value)} /></label><label>方案类型<select required value={manualOptionType} onChange={(e) => setManualOptionType(e.target.value)}><option value="">请选择</option>{(activeRequirement?.allowed_option_types ?? draft.optionTypes).map((value) => <option key={value} value={value}>{OPTION_TYPES.find(([code]) => code === value)?.[1] ?? value}</option>)}</select></label><label>业务理由<textarea required minLength={6} value={manualRationale} onChange={(e) => setManualRationale(e.target.value)} /></label><label>最小估算（CNY）<input required min="0" step="0.01" type="number" value={manualMinimum} onChange={(e) => setManualMinimum(e.target.value)} /></label><label>最可能估算（CNY）<input required min="0" step="0.01" type="number" value={manualLikely} onChange={(e) => setManualLikely(e.target.value)} /></label><label>最大估算（CNY）<input required min="0" step="0.01" type="number" value={manualMaximum} onChange={(e) => setManualMaximum(e.target.value)} /></label><label>估算依据<textarea required value={manualBasis} onChange={(e) => setManualBasis(e.target.value)} /></label><label>预计可用日期<input required type="date" value={manualAvailableAt} onChange={(e) => setManualAvailableAt(e.target.value)} /></label><label>假设 <small>每行一项</small><textarea required value={manualAssumptions} onChange={(e) => setManualAssumptions(e.target.value)} /></label><label>待核验事实 <small>每行一项</small><textarea required value={manualFacts} onChange={(e) => setManualFacts(e.target.value)} /></label><label>主要风险 <small>每行一项</small><textarea required value={manualRisks} onChange={(e) => setManualRisks(e.target.value)} /></label><button disabled={manualBusy || !activeRequirement}>{manualBusy ? <LoaderCircle className="gx-spin" /> : <FilePlus2 />}{manualBusy ? "正在提交…" : "提交到 IAOS 候选集"}</button></form>}
      {proposalSet && <section className="plant-proposals" aria-live="polite"><header><div><span>{proposalSet.status}</span><h2>候选方案 · {proposalSet.proposals.length}</h2></div><small>{proposalSet.evidence.provider} / {proposalSet.evidence.model} · {proposalSet.evidence.prompt_version}</small></header><div className="plant-proposal-grid">{proposalSet.proposals.map((proposal) => <ProposalCard key={proposal.proposal_id} proposal={proposal} review={reviews[proposal.proposal_id]} onReview={(review) => { setReviews((current) => ({ ...current, [proposal.proposal_id]: review })); setSavedReviews((current) => ({ ...current, [proposal.proposal_id]: false })); }} onSubmitReview={() => submitReview(proposal.proposal_id)} saved={Boolean(savedReviews[proposal.proposal_id])} busy={reviewBusy === proposal.proposal_id} investigation={investigations.find((item) => item.request.proposal_id === proposal.proposal_id)} investigationBusy={investigationBusy === proposal.proposal_id} onRequestInvestigation={() => startInvestigation(proposal.proposal_id)} />)}</div><p className="plant-persistence-note"><BadgeCheck />Agent 候选、人工新增候选和相应审阅都通过 IAOS Capability 保存并形成版本；“采纳调研”不等于投资批准或合同签署。</p></section>}
      {investigations.length > 0 && <section className="plant-investigations" aria-live="polite"><header><div><span>facility.site.investigation.v1</span><h2>场址外部调研工作项</h2></div><small>{investigations.filter((item) => item.status === "waiting_world").length} 条等待 World</small></header>{investigations.map((item) => <InvestigationPanel key={item.request.investigation_request_id} item={item} caseCode={caseCode} onCommitted={refreshInvestigations} />)}</section>}
      {assessments.length > 0 && <section className="plant-assessments" aria-labelledby="plant-assessment-title">
        <header><div><span>OBSERVATION-ONLY COMPARISON</span><h2 id="plant-assessment-title">外部事实比较</h2><p>只有已送达的 World Observation 参与评分；Agent 估算只用于对照，不会覆盖正式报价或现场事实。</p></div><small>{assessments.filter((item) => item.eligible).length}/{assessments.length} 个候选通过硬约束</small></header>
        <fieldset className="plant-assessment-weights"><legend>本次比较权重（自动归一化）</legend>{SCORE_DIMENSIONS.map(([key, label, help]) => <label key={key}>{label}<small>{help}</small><input aria-label={`${label}权重`} type="number" min="0" max="100" value={assessmentWeights[key]} onChange={(event) => setAssessmentWeights((current) => ({ ...current, [key]: Math.max(0, Number(event.target.value) || 0) }))} /></label>)}<p>当前权重合计 {normalizedWeightTotal}。权重只改变合格候选的比较视图；任何硬约束失败都不会被高分抵消。</p></fieldset>
        <div className="plant-assessment-grid">{assessments.map((assessment) => <AssessmentCard key={assessment.observation_id} assessment={assessment} />)}</div>
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
        {latestSelection && <article className="plant-selection-status" aria-live="polite"><header><div><span>{latestSelection.recommendation.approval_flow_key}</span><h3>推荐与审批状态</h3></div><strong>{latestSelection.decision ? "已正式选址" : latestSelection.approval_status}</strong></header><dl><div><dt>推荐候选</dt><dd>{latestSelection.recommendation.selected_proposal_id}</dd></div><div><dt>推荐编号</dt><dd>{latestSelection.recommendation.recommendation_id}</dd></div><div><dt>审批请求</dt><dd>{latestSelection.recommendation.approval_request_id}</dd></div><div><dt>权威输入哈希</dt><dd>{latestSelection.recommendation.input_hash}</dd></div></dl><div className="plant-selection-actions"><a href={approvalURL} target="_blank" rel="noreferrer">打开 IAOS 审批中心</a><button type="button" disabled={selectionBusy || latestSelection.approval_status !== "approved" || Boolean(latestSelection.decision)} onClick={() => formalizeSelection(latestSelection)}>{latestSelection.decision ? "正式选择已落地" : latestSelection.approval_status === "approved" ? "同步批准并正式选址" : "等待审批决定"}</button></div></article>}
      </section>}
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
