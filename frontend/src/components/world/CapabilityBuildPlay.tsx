import {
  ArrowLeft,
  BadgeCheck,
  Factory,
  Pause,
  Play,
  RotateCcw,
  StepForward,
  Users,
} from "lucide-react";
import { useEffect, useState } from "react";
import {
  loadCapabilityBuild,
  type CapabilityTrace,
} from "../../world/capabilityBuild";
import "./WorldPlay.css";
export function CapabilityBuildPlay({ onExit }: { onExit: () => void }) {
  const [t, setT] = useState<CapabilityTrace | null>(null),
    [s, setS] = useState(0),
    [play, setPlay] = useState(false),
    [err, setErr] = useState("");
  useEffect(() => {
    const c = new AbortController();
    loadCapabilityBuild(c.signal)
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
      700,
    );
    return () => clearInterval(id);
  }, [play, t]);
  if (err)
    return (
      <main className="world-error" role="alert">
        {err}
      </main>
    );
  if (!t) return <main className="world-loading">正在加载生产能力世界…</main>;
  const f = t.frames[s];
  return (
    <div className="world-play">
      <header className="world-toolbar">
        <button className="world-back" onClick={onExit}>
          <ArrowLeft />
          工厂建设
        </button>
        <div>
          <span>PROJECT GENESIS · CAPABILITY BUILD</span>
          <h1>A 线设备、团队与资格建设</h1>
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
        <button
          onClick={() => {
            setPlay(false);
            setS(0);
          }}
        >
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
            className={`world-state-badge ${f.industrialization_eligible ? "closed" : "active"}`}
          >
            {f.industrialization_eligible ? <BadgeCheck /> : <Factory />}
            {f.industrialization_eligible ? "M12 eligible" : f.phase}
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
              <Factory />
              设备与实验室 <small>World owns ability</small>
            </header>
            <dl>
              {f.equipment.length ? (
                f.equipment.map((x) => (
                  <div key={x.code}>
                    <dt>{x.code}</dt>
                    <dd>
                      {x.status}
                      <small>{x.zone}</small>
                    </dd>
                  </div>
                ))
              ) : (
                <div>
                  <dt>设备</dt>
                  <dd>尚未形成实际能力</dd>
                </div>
              )}
            </dl>
          </article>
          <article>
            <header>
              <Users />
              人员与技能 <small>privacy scoped</small>
            </header>
            <dl>
              <div>
                <dt>实际到岗资格</dt>
                <dd>{f.workers.length} / 10</dd>
              </div>
              <div>
                <dt>角色认知</dt>
                <dd>{f.knowledge.length} 条</dd>
              </div>
            </dl>
          </article>
          <article>
            <header>
              <BadgeCheck />
              联合资格门 <small>IAOS governance</small>
            </header>
            <dl>
              {Object.entries(f.gate).map(([k, v]) => (
                <div key={k}>
                  <dt>{k}</dt>
                  <dd>{v ? "通过" : "阻塞"}</dd>
                </div>
              ))}
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
