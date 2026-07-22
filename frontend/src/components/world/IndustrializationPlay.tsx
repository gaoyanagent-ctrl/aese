import {
  ArrowLeft,
  BadgeCheck,
  ClipboardCheck,
  Pause,
  Play,
  RotateCcw,
  StepForward,
  Wrench,
} from "lucide-react";
import { useEffect, useState } from "react";
import {
  loadIndustrialization,
  type IndustrializationTrace,
} from "../../world/industrialization";
import "./WorldPlay.css";
export function IndustrializationPlay({ onExit }: { onExit: () => void }) {
  const [t, setT] = useState<IndustrializationTrace | null>(null),
    [s, setS] = useState(0),
    [play, setPlay] = useState(false),
    [err, setErr] = useState("");
  useEffect(() => {
    const c = new AbortController();
    loadIndustrialization(c.signal)
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
      650,
    );
    return () => clearInterval(id);
  }, [play, t]);
  if (err)
    return (
      <main className="world-error" role="alert">
        {err}
      </main>
    );
  if (!t) return <main className="world-loading">正在加载产品工业化世界…</main>;
  const f = t.frames[s];
  return (
    <div className="world-play">
      <header className="world-toolbar">
        <button className="world-back" onClick={onExit}>
          <ArrowLeft />
          能力建设
        </button>
        <div>
          <span>PROJECT GENESIS · INDUSTRIALIZATION</span>
          <h1>HCTM-BCP-A01 产品工业化</h1>
        </div>
        <div className="world-clock">
          <small>虚拟时间 · Asia/Shanghai</small>
          <strong>{f.sim_time}</strong>
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
        <button onClick={() => setS(0)}>
          <RotateCcw />
          复位
        </button>
        <span>
          {s + 1}/{t.frames.length} · {f.phase}
        </span>
      </nav>
      <main className="world-main">
        <section className="world-status">
          <span
            className={`world-state-badge ${f.serial_production_eligible ? "closed" : "active"}`}
          >
            {f.serial_production_eligible ? <BadgeCheck /> : <Wrench />}
            {f.serial_production_eligible ? "M13 eligible" : f.phase}
          </span>
          <h2>{f.title}</h2>
          <p>
            World {f.world_progress}% / IAOS {f.iaos_progress}% ·{" "}
            {f.discrepancy}
          </p>
        </section>
        <section className="three-state-grid">
          <article>
            <header>
              <ClipboardCheck />
              APQP / PPAP
            </header>
            <dl>
              {Object.entries(f.apqp_gates).map(([k, v]) => (
                <div key={k}>
                  <dt>{k}</dt>
                  <dd>{v ? "通过" : "阻塞"}</dd>
                </div>
              ))}
              <div>
                <dt>PPAP</dt>
                <dd>{f.ppap_status}</dd>
              </div>
            </dl>
          </article>
          <article>
            <header>
              <Wrench />
              试制与质量 <small>World owns</small>
            </header>
            <dl>
              {f.trials.map((x) => (
                <div key={x.code}>
                  <dt>{x.code}</dt>
                  <dd>
                    rev {x.revision} · Cpk {x.Cpk}
                    <small>
                      良率 {x.Yield}% · 泄漏 {x.leak_failures}
                    </small>
                  </dd>
                </div>
              ))}
            </dl>
          </article>
          <article>
            <header>
              <BadgeCheck />
              发布与兼容
            </header>
            <dl>
              {f.releases.map((x) => (
                <div key={x.code}>
                  <dt>{x.code}</dt>
                  <dd>
                    rev {x.revision}
                    <small>{x.hash.slice(0, 12)}</small>
                  </dd>
                </div>
              ))}
              <div>
                <dt>旧 HCTM</dt>
                <dd>{f.compatibility}</dd>
              </div>
            </dl>
          </article>
        </section>
        <section className="world-timeline">
          {t.frames.map((x, i) => (
            <button
              key={x.step}
              className={i === s ? "current" : ""}
              onClick={() => setS(i)}
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
