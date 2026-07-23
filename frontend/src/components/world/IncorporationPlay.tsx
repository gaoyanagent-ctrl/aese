import {
  ArrowLeft,
  Activity,
  BadgeCheck,
  Banknote,
  Building2,
  Database,
  ExternalLink,
  GitBranch,
  RefreshCw,
  RotateCcw,
  ShieldCheck,
  StepForward,
  Users,
} from "lucide-react";
import { useEffect, useState } from "react";
import {
  loadIncorporation,
  type IncorporationTrace,
} from "../../world/incorporation";
import {
  buildIncorporationStepTrace,
  unlockedIncorporationFrame,
} from "../../world/incorporationStepTrace";
import "./WorldPlay.css";
const money = (value: string) =>
  new Intl.NumberFormat("zh-CN", {
    style: "currency",
    currency: "CNY",
    maximumFractionDigits: 2,
  }).format(Number(value));
export function IncorporationPlay({ onExit }: { onExit: () => void }) {
  const [trace, setTrace] = useState<IncorporationTrace | null>(null);
  const [step, setStep] = useState(0);
  const [error, setError] = useState("");
  const [refreshing, setRefreshing] = useState(false);
  const refresh = () => {
    setRefreshing(true);
    loadIncorporation()
      .then((next) => {
        setTrace(next);
        setStep((current) =>
          Math.min(current, unlockedIncorporationFrame(next.iaos_lifecycle)),
        );
        setError("");
      })
      .catch((e) => setError(String(e)))
      .finally(() => setRefreshing(false));
  };
  useEffect(() => {
    const c = new AbortController();
    loadIncorporation(c.signal)
      .then(setTrace)
      .catch((e) => {
        if (e.name !== "AbortError") setError(String(e));
      });
    return () => c.abort();
  }, []);
  if (error)
    return (
      <main className="world-error" role="alert">
        企业成立 campaign 加载失败：{error}
        <button onClick={onExit}>返回</button>
      </main>
    );
  if (!trace)
    return <main className="world-loading">正在加载企业成立世界…</main>;
  const f = trace.frames[step];
  const lifecycle = trace.iaos_lifecycle;
  const unlockedFrame = unlockedIncorporationFrame(lifecycle);
  const params = new URLSearchParams(window.location.hash.split("?")[1] ?? "");
  const tenant = params.get("tenant") ?? "tenant-hctm-genesis";
  const caseCode = lifecycle?.case_code ?? params.get("case") ?? "";
  const processRun = params.get("process_run") ?? String(lifecycle?.process_runs?.[0]?.id ?? "");
  const correlation = params.get("correlation") ?? String(lifecycle?.world_exchanges?.[0]?.correlation_id ?? "");
  const stepTrace = buildIncorporationStepTrace(f, lifecycle);
  const primaryCapability = stepTrace?.definition.capabilities[0] ?? "";
  const iaosOrigin = `http://${window.location.hostname || "127.0.0.1"}:3000`;
  const iaosLink = `${iaosOrigin}/#enterprise_lifecycle?tenant=${encodeURIComponent(tenant)}&case=${encodeURIComponent(caseCode)}&process_run=${encodeURIComponent(processRun)}&world_run=${encodeURIComponent(trace.world_run_id)}&correlation=${encodeURIComponent(correlation)}&step=${f.step}&capability=${encodeURIComponent(primaryCapability)}`;
  const roles = [
    f.governance.ceo,
    f.governance.cfo,
    f.governance.project_director,
  ];
  return (
    <div className="world-play">
      <header className="world-toolbar">
        <button className="world-back" onClick={onExit}>
          <ArrowLeft />
          生命周期总览
        </button>
        <div>
          <span>PROJECT GENESIS · INCORPORATION</span>
          <h1>华辰苏州制造公司成立与治理</h1>
        </div>
        <div className="world-clock">
          <small>虚拟时间 · Asia/Shanghai</small>
          <strong>
            {new Date(f.sim_time).toLocaleString("zh-CN", {
              timeZone: "Asia/Shanghai",
              hour12: false,
            })}
          </strong>
        </div>
      </header>
      <nav className="world-controls">
        <button onClick={refresh} disabled={refreshing}>
          <RefreshCw />
          {refreshing ? "同步中" : "同步 IAOS"}
        </button>
        <button
          onClick={() => setStep((v) => Math.min(v + 1, unlockedFrame))}
          disabled={step >= unlockedFrame}
          title={step >= unlockedFrame ? "请先在 IAOS 完成当前工作项" : ""}
        >
          <StepForward />
          查看下一已完成阶段
        </button>
        <button
          onClick={() => setStep(0)}
        >
          <RotateCcw />
          复位
        </button>
        <button onClick={() => (window.location.hash = "world-plant-build")}>
          <Building2 />
          工厂建设 Campaign
        </button>
        <a className="world-back" href={iaosLink} target="_blank" rel="noreferrer">
          打开 IAOS 设立案
        </a>
        <span>
          已解锁 {unlockedFrame + 1}/{trace.frames.length} · 当前查看{" "}
          {step + 1} · {f.phase}
        </span>
      </nav>
      <main className="world-main">
        <section className="world-status">
          <span
            className={`world-state-badge ${f.plant_project_eligible ? "closed" : "active"}`}
          >
            {f.plant_project_eligible ? <BadgeCheck /> : <Building2 />}
            {f.plant_project_eligible ? "M10 eligible" : f.phase}
          </span>
          <h2>{f.title}</h2>
          <p>
            因果引用：<code>{f.causation_id}</code> · IAOS cursor{" "}
            {f.iaos_cursor}
          </p>
        </section>
        <section className="three-state-grid">
          <article>
            <header>
              <Building2 />
              法人和外部世界 <small>World owns</small>
            </header>
            <dl>
              <div>
                <dt>法人</dt>
                <dd>HCTM-SZ-MFG · {f.legal_entity_status}</dd>
              </div>
              <div>
                <dt>登记</dt>
                <dd>{f.registration_status}</dd>
              </div>
              <div>
                <dt>公司账户</dt>
                <dd>
                  {f.company.status} · {money(f.company.balance.value)}
                </dd>
              </div>
              <div>
                <dt>M10 资格</dt>
                <dd>{f.plant_project_eligible ? "是" : "否"}</dd>
              </div>
            </dl>
          </article>
          <article>
            <header>
              <Banknote />
              资金与预算 <small>来源：虚构基线</small>
            </header>
            <dl>
              <div>
                <dt>投资人现金</dt>
                <dd>{money(f.investor.balance.value)}</dd>
              </div>
              <div>
                <dt>认缴资本</dt>
                <dd>{money(f.capital_committed.value)}</dd>
              </div>
              <div>
                <dt>实缴资本</dt>
                <dd>{money(f.capital_paid.value)}</dd>
              </div>
              <div>
                <dt>预算授权</dt>
                <dd>
                  {money(f.budget.amount.value)} · {f.budget.status}
                  <small>预算不是现金或实际支出</small>
                </dd>
              </div>
            </dl>
          </article>
          <article>
            <header>
              <Users />
              治理与角色认知 <small>actor scoped</small>
            </header>
            <dl>
              {roles.map((role) => (
                <div key={role.position}>
                  <dt>{role.position.replace("POS-HCTM-", "")}</dt>
                  <dd>
                    {role.status}
                    <small>{role.assignee}</small>
                  </dd>
                </div>
              ))}
              <div>
                <dt>Mandate</dt>
                <dd>{f.governance.mandate_active ? "active" : "inactive"}</dd>
              </div>
              <div>
                <dt>已获知</dt>
                <dd>{f.knowledge.length} 条角色事实</dd>
              </div>
            </dl>
          </article>
        </section>
        <section className="world-timeline incorporation-timeline">
          {trace.frames.map((item, index) => (
            <button
              key={item.step}
              className={index === step ? "current" : ""}
              disabled={index > unlockedFrame}
              title={
                index > unlockedFrame
                  ? "该阶段尚未由 IAOS 的能力、审批与 World Observation 提交"
                  : ""
              }
              onClick={() => setStep(index)}
            >
              <span>{index + 1}</span>
              <strong>{item.phase}</strong>
              <small>{item.title}</small>
            </button>
          ))}
        </section>
        {stepTrace && (
          <section
            className="step-trace"
            data-testid="incorporation-step-trace"
            aria-labelledby="step-trace-heading"
          >
            <header className="step-trace-header">
              <div>
                <span>当前步骤的 IAOS 执行证据</span>
                <h2 id="step-trace-heading">
                  步骤 {f.step + 1} · {stepTrace.definition.process}
                </h2>
                <p>
                  这里只显示当前 World 步骤关联的 IAOS 记录，不混入整案其他步骤。
                </p>
              </div>
              <a href={iaosLink} target="_blank" rel="noreferrer">
                在 IAOS 中定位
                <ExternalLink aria-hidden="true" />
              </a>
            </header>
            <div className="step-trace-summary">
              <article>
                <GitBranch aria-hidden="true" />
                <span>Process</span>
                <strong>{stepTrace.definition.process}</strong>
                <small>
                  run {String(stepTrace.processRun?.id ?? processRun ?? "—")}
                </small>
              </article>
              <article>
                <Activity aria-hidden="true" />
                <span>Capabilities</span>
                <strong>{stepTrace.definition.capabilities.length} 项</strong>
                <small>{stepTrace.definition.states.join(" → ")}</small>
              </article>
              <article>
                <Database aria-hidden="true" />
                <span>Entity 影响</span>
                <strong>{stepTrace.definition.entities.length} 类</strong>
                <small>{stepTrace.definition.entities.join(" · ")}</small>
              </article>
              <article>
                <ShieldCheck aria-hidden="true" />
                <span>治理与事务证据</span>
                <strong>
                  {stepTrace.decisions.length + stepTrace.approvals.length} 条
                </strong>
                <small>
                  Journal {stepTrace.journal.length} · Outbox{" "}
                  {stepTrace.outbox.length}
                </small>
              </article>
            </div>
            <div className="step-capability-list">
              {stepTrace.definition.capabilities.map((capability) => {
                const transition = stepTrace.transitions.find(
                  (item) => item.capability === capability,
                );
                return (
                  <article key={capability}>
                    <span className={transition ? "committed" : "unmatched"}>
                      {transition ? "已提交" : "未匹配"}
                    </span>
                    <div>
                      <strong>{capability}</strong>
                      <small>
                        {String(transition?.actor_id ?? "—")} ·{" "}
                        {String(transition?.idempotency_key ?? "—")}
                      </small>
                    </div>
                  </article>
                );
              })}
            </div>
            <nav className="step-studio-links" aria-label="当前步骤平台资产入口">
              <a
                href={`${iaosOrigin}/#capability_studio?capability=${encodeURIComponent(primaryCapability)}`}
                target="_blank"
                rel="noreferrer"
              >
                Capability Studio
              </a>
              <a
                href={`${iaosOrigin}/#process_studio?process=${encodeURIComponent(stepTrace.definition.process)}`}
                target="_blank"
                rel="noreferrer"
              >
                Process Studio
              </a>
              {stepTrace.definition.entities.map((entity) => (
                <a
                  key={entity}
                  href={`${iaosOrigin}/#entity_explorer?entity=${encodeURIComponent(entity)}`}
                  target="_blank"
                  rel="noreferrer"
                >
                  Entity · {entity}
                </a>
              ))}
            </nav>
            <div className="step-evidence-groups">
              {[
                ["World Bridge", stepTrace.worldExchanges],
                [
                  "Decision / Approval",
                  [...stepTrace.decisions, ...stepTrace.approvals],
                ],
                [
                  "Journal / Outbox",
                  [...stepTrace.journal, ...stepTrace.outbox],
                ],
              ].map(([label, evidence]) => (
                <details key={String(label)}>
                  <summary>
                    {String(label)}
                    <span>{(evidence as unknown[]).length} 条</span>
                  </summary>
                  {(evidence as unknown[]).length ? (
                    <pre>{JSON.stringify(evidence, null, 2)}</pre>
                  ) : (
                    <p>当前步骤没有这类证据。</p>
                  )}
                </details>
              ))}
            </div>
          </section>
        )}
        {lifecycle && <section className="world-status" data-testid="iaos-lifecycle-projection">
          <span>IAOS Effective Runtime Projection</span>
          <h2>{lifecycle.case_code} · {String((lifecycle.state as {state?:string})?.state ?? "")}</h2>
          <p>Process {processRun || "—"} · correlation {correlation || "—"}</p>
          <div className="three-state-grid">
            <article><header>Intent / Observation / CommittedOutcome</header><pre>{JSON.stringify(lifecycle.world_exchanges ?? [],null,2)}</pre></article>
            <article><header>Discrepancy / Decision</header><pre>{JSON.stringify({discrepancies:(lifecycle.state as {discrepancies?:unknown[]})?.discrepancies,decisions:lifecycle.decisions},null,2)}</pre></article>
            <article><header>Runtime lineage</header><pre>{JSON.stringify(lifecycle.lineage ?? {},null,2)}</pre></article>
          </div>
        </section>}
      </main>
    </div>
  );
}
