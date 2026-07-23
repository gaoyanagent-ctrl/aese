import {
  ArrowLeft,
  BadgeCheck,
  Banknote,
  Building2,
  Pause,
  Play,
  RotateCcw,
  StepForward,
  Users,
} from "lucide-react";
import { useEffect, useState } from "react";
import {
  loadIncorporation,
  type IncorporationTrace,
} from "../../world/incorporation";
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
  const [playing, setPlaying] = useState(false);
  const [error, setError] = useState("");
  useEffect(() => {
    const c = new AbortController();
    loadIncorporation(c.signal)
      .then(setTrace)
      .catch((e) => {
        if (e.name !== "AbortError") setError(String(e));
      });
    return () => c.abort();
  }, []);
  useEffect(() => {
    if (!playing || !trace) return;
    const id = setInterval(
      () =>
        setStep((v) => {
          if (v >= trace.frames.length - 1) {
            setPlaying(false);
            return v;
          }
          return v + 1;
        }),
      850,
    );
    return () => clearInterval(id);
  }, [playing, trace]);
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
  const params = new URLSearchParams(window.location.hash.split("?")[1] ?? "");
  const tenant = params.get("tenant") ?? "tenant-hctm-genesis";
  const caseCode = lifecycle?.case_code ?? params.get("case") ?? "";
  const processRun = params.get("process_run") ?? String(lifecycle?.process_runs?.[0]?.id ?? "");
  const correlation = params.get("correlation") ?? String(lifecycle?.world_exchanges?.[0]?.correlation_id ?? "");
  const iaosLink = `http://${window.location.hostname || "127.0.0.1"}:3000/#enterprise_lifecycle?tenant=${encodeURIComponent(tenant)}&case=${encodeURIComponent(caseCode)}&process_run=${encodeURIComponent(processRun)}&world_run=${encodeURIComponent(trace.world_run_id)}&correlation=${encodeURIComponent(correlation)}`;
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
        <button onClick={() => setPlaying((v) => !v)}>
          {playing ? <Pause /> : <Play />}
          {playing ? "暂停" : "运行"}
        </button>
        <button
          onClick={() =>
            setStep((v) => Math.min(v + 1, trace.frames.length - 1))
          }
          disabled={step === trace.frames.length - 1}
        >
          <StepForward />
          单步
        </button>
        <button
          onClick={() => {
            setPlaying(false);
            setStep(0);
          }}
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
          阶段 {step + 1}/{trace.frames.length} · {f.phase}
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
              onClick={() => {
                setPlaying(false);
                setStep(index);
              }}
            >
              <span>{index + 1}</span>
              <strong>{item.phase}</strong>
              <small>{item.title}</small>
            </button>
          ))}
        </section>
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
