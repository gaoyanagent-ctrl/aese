import {
  ArrowLeft,
  BadgeCheck,
  Banknote,
  Building2,
  Factory,
  Pause,
  Play,
  RotateCcw,
  StepForward,
  TriangleAlert,
} from "lucide-react";
import { useEffect, useState } from "react";
import { loadPlantBuild, type PlantBuildTrace } from "../../world/plantBuild";
import "./WorldPlay.css";
const cny = (v: string) =>
  new Intl.NumberFormat("zh-CN", {
    style: "currency",
    currency: "CNY",
    maximumFractionDigits: 0,
  }).format(Number(v));
export function PlantBuildPlay({ onExit }: { onExit: () => void }) {
  const [t, setT] = useState<PlantBuildTrace | null>(null);
  const [s, setS] = useState(0);
  const [play, setPlay] = useState(false);
  const [err, setErr] = useState("");
  useEffect(() => {
    const c = new AbortController();
    loadPlantBuild(c.signal)
      .then(setT)
      .catch((e) => {
        if (e.name !== "AbortError") setErr(String(e));
      });
    return () => c.abort();
  }, []);
  useEffect(() => {
    if (!play || !t) return;
    const id = setInterval(
      () =>
        setS((v) => {
          if (v >= t.frames.length - 1) {
            setPlay(false);
            return v;
          }
          return v + 1;
        }),
      750,
    );
    return () => clearInterval(id);
  }, [play, t]);
  if (err)
    return (
      <main className="world-error" role="alert">
        Plant Build 加载失败：{err}
      </main>
    );
  if (!t) return <main className="world-loading">正在加载设施建设世界…</main>;
  const f = t.frames[s];
  return (
    <div className="world-play">
      <header className="world-toolbar">
        <button className="world-back" onClick={onExit}>
          <ArrowLeft />
          企业成立
        </button>
        <div>
          <span>PROJECT GENESIS · PLANT BUILD</span>
          <h1>苏州一期选址与设施建设</h1>
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
        <button onClick={() => setPlay((v) => !v)}>
          {play ? <Pause /> : <Play />}
          {play ? "暂停" : "运行"}
        </button>
        <button
          onClick={() => setS((v) => Math.min(v + 1, t.frames.length - 1))}
          disabled={s === t.frames.length - 1}
        >
          <StepForward />
          单步
        </button>
        <button
          onClick={() => {
            setPlay(false);
            setS(0);
          }}
        >
          <RotateCcw />
          复位
        </button>
        <button
          onClick={() => (window.location.hash = "world-capability-build")}
        >
          <Factory />
          生产能力 Campaign
        </button>
        <span>
          {s + 1}/{t.frames.length} · {f.phase}
        </span>
      </nav>
      <main className="world-main">
        <section className="world-status">
          <span
            className={`world-state-badge ${f.capability_build_eligible ? "closed" : "active"}`}
          >
            {f.capability_build_eligible ? <BadgeCheck /> : <Factory />}
            {f.capability_build_eligible ? "M11 eligible" : f.phase}
          </span>
          <h2>{f.title}</h2>
          <p>
            World {f.world_progress}% / IAOS plan {f.iaos_plan_progress}% ·{" "}
            {f.discrepancy} · cursor {f.iaos_cursor}
          </p>
        </section>
        <section className="three-state-grid">
          <article>
            <header>
              <Building2 />
              场址与项目 <small>World owns reality</small>
            </header>
            <dl>
              <div>
                <dt>选中场址</dt>
                <dd>{f.selected_site || "尚未批准"}</dd>
              </div>
              {f.assessments.map((a) => (
                <div key={a.site_code}>
                  <dt>{a.site_code.replace("SITE-SZ-", "")}</dt>
                  <dd>
                    {a.feasible ? "可行" : "不可行"} · {a.weighted_score}
                    <small>{a.hard_failures.join(", ") || "硬约束通过"}</small>
                  </dd>
                </div>
              ))}
            </dl>
          </article>
          <article>
            <header>
              <Banknote />
              治理资金 <small>IAOS owns approvals</small>
            </header>
            <dl>
              <div>
                <dt>可用现金</dt>
                <dd>{cny(f.cash.value)}</dd>
              </div>
              <div>
                <dt>合同承诺</dt>
                <dd>{cny(f.committed.value)}</dd>
              </div>
              <div>
                <dt>里程碑应付</dt>
                <dd>{cny(f.payable.value)}</dd>
              </div>
              <div>
                <dt>已付款</dt>
                <dd>{cny(f.paid.value)}</dd>
              </div>
            </dl>
          </article>
          <article>
            <header>
              <TriangleAlert />
              Knowledge / discrepancy <small>actor scoped</small>
            </header>
            <dl>
              <div>
                <dt>项目负责人已知</dt>
                <dd>{f.knowledge.length} 条</dd>
              </div>
              <div>
                <dt>差异</dt>
                <dd>{f.discrepancy}</dd>
              </div>
              <div>
                <dt>因果</dt>
                <dd>{f.causation_id}</dd>
              </div>
            </dl>
          </article>
        </section>
        {f.zones.length > 0 && (
          <section className="facility-zones">
            {f.zones
              .filter((z) => !["site", "building"].includes(z.purpose))
              .map((z) => (
                <article key={z.code}>
                  <Factory />
                  <strong>{z.purpose}</strong>
                  <code>{z.code}</code>
                  <small>
                    {z.area_m2} m² · {z.status}
                  </small>
                </article>
              ))}
          </section>
        )}
        <section className="world-timeline">
          {t.frames.map((x, i) => (
            <button
              key={x.step}
              className={i === s ? "current" : ""}
              onClick={() => {
                setPlay(false);
                setS(i);
              }}
            >
              <span>{i + 1}</span>
              <strong>{x.phase}</strong>
              <small>{x.title}</small>
            </button>
          ))}
        </section>
      </main>
    </div>
  );
}
