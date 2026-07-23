import {
  ArrowLeft,
  CheckCircle2,
  CircleAlert,
  Database,
  Eye,
  Pause,
  Play,
  RotateCcw,
  StepForward,
} from "lucide-react";
import { useEffect, useState } from "react";
import { loadGenesisTrace } from "../../world/api";
import type { GenesisTrace } from "../../world/types";
import "./WorldPlay.css";
export function WorldPlay({ onExit }: { onExit: () => void }) {
  const [trace, setTrace] = useState<GenesisTrace | null>(null);
  const [step, setStep] = useState(0);
  const [playing, setPlaying] = useState(false);
  const [error, setError] = useState("");
  useEffect(() => {
    const c = new AbortController();
    loadGenesisTrace(c.signal)
      .then(setTrace)
      .catch((e) => {
        if (e.name !== "AbortError") setError(String(e));
      });
    return () => c.abort();
  }, []);
  useEffect(() => {
    if (!playing || !trace) return;
    const id = window.setInterval(
      () =>
        setStep((v) => {
          if (v >= trace.frames.length - 1) {
            setPlaying(false);
            return v;
          }
          return v + 1;
        }),
      900,
    );
    return () => clearInterval(id);
  }, [playing, trace]);
  if (error)
    return (
      <main className="world-error" role="alert">
        <CircleAlert />
        <strong>World API 不可用</strong>
        <p>{error}</p>
        <button onClick={onExit}>返回沙盘</button>
      </main>
    );
  if (!trace)
    return (
      <main className="world-loading" aria-live="polite">
        正在加载三态世界…
      </main>
    );
  const frame = trace.frames[step];
  const closed = frame.discrepancy.status === "closed";
  return (
    <div className="world-play">
      <header className="world-toolbar">
        <button className="world-back" onClick={onExit}>
          <ArrowLeft />
          生命周期总览
        </button>
        <div>
          <span>PROJECT GENESIS · WORLD PLAY</span>
          <h1>LAS-WLD-02 三态偏差闭环</h1>
        </div>
        <div className="world-clock">
          <small>虚拟时间 · Asia/Shanghai</small>
          <strong>
            {new Date(frame.sim_time).toLocaleString("zh-CN", {
              timeZone: "Asia/Shanghai",
              hour12: false,
            })}
          </strong>
        </div>
      </header>
      <nav className="world-controls" aria-label="世界时间控制">
        <button onClick={() => setPlaying((v) => !v)} aria-pressed={playing}>
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
        <button onClick={() => { window.location.hash = "world-incorporation"; }}>
          企业成立 Campaign
        </button>
        <span>
          步骤 {step + 1}/{trace.frames.length}
        </span>
      </nav>
      <main className="world-main">
        <section className="world-status" aria-labelledby="world-frame-title">
          <span className={`world-state-badge ${closed ? "closed" : "active"}`}>
            {closed ? <CheckCircle2 /> : <CircleAlert />}
            {frame.discrepancy.status}
          </span>
          <h2 id="world-frame-title">{frame.title}</h2>
          <p>
            因果引用：<code>{frame.causation_id}</code>
          </p>
        </section>
        <section className="three-state-grid">
          <article>
            <header>
              <Eye />
              客观世界 <small>AESE owns</small>
            </header>
            <dl>
              <div>
                <dt>设备</dt>
                <dd>{frame.world.code}</dd>
              </div>
              <div>
                <dt>实际状态</dt>
                <dd>{frame.world.condition}</dd>
              </div>
              <div>
                <dt>振动</dt>
                <dd>
                  {frame.world.vibration} {frame.world.unit}
                </dd>
              </div>
              <div>
                <dt>可用产能</dt>
                <dd>{frame.world.available_capacity}</dd>
              </div>
            </dl>
          </article>
          <article>
            <header>
              <Database />
              管理事实 <small>IAOS owns</small>
            </header>
            <dl>
              <div>
                <dt>登记状态</dt>
                <dd>{frame.iaos.registered_status}</dd>
              </div>
              <div>
                <dt>Journal cursor</dt>
                <dd>{frame.iaos.cursor}</dd>
              </div>
              <div>
                <dt>维护工单</dt>
                <dd>{frame.iaos.maintenance_order || "尚未登记"}</dd>
              </div>
            </dl>
          </article>
          <article>
            <header>
              <Eye />
              角色认知 <small>Actor scoped</small>
            </header>
            {frame.knowledge.length ? (
              <dl>
                {frame.knowledge.map((k) => (
                  <div key={k.knowledge_id}>
                    <dt>{k.actor_ref.code}</dt>
                    <dd>
                      已观察 · 置信度 {k.confidence}
                      <small>{k.source_ref}</small>
                    </dd>
                  </div>
                ))}
              </dl>
            ) : (
              <p className="world-empty">
                角色尚未知，不可读取完整 World State。
              </p>
            )}
          </article>
        </section>
        <section className="world-timeline" aria-label="偏差时间线">
          {trace.frames.map((item, index) => (
            <button
              key={item.step}
              className={index === step ? "current" : ""}
              onClick={() => {
                setPlaying(false);
                setStep(index);
              }}
              aria-current={index === step ? "step" : undefined}
            >
              <span>{index + 1}</span>
              <strong>{item.discrepancy.status}</strong>
              <small>{item.title}</small>
            </button>
          ))}
        </section>
      </main>
    </div>
  );
}
