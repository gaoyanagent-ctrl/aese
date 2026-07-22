import {
  ArrowLeft,
  BadgeCheck,
  Banknote,
  PackageCheck,
  Pause,
  Play,
  RotateCcw,
  StepForward,
  Truck,
} from "lucide-react";
import { useEffect, useState } from "react";
import {
  loadFirstDelivery,
  type FirstDeliveryTrace,
} from "../../world/firstDelivery";
import "./WorldPlay.css";
const c = (v: string) =>
  new Intl.NumberFormat("zh-CN", {
    style: "currency",
    currency: "CNY",
    maximumFractionDigits: 0,
  }).format(Number(v));
export function FirstDeliveryPlay({ onExit }: { onExit: () => void }) {
  const [t, setT] = useState<FirstDeliveryTrace | null>(null),
    [s, setS] = useState(0),
    [play, setPlay] = useState(false),
    [err, setErr] = useState("");
  useEffect(() => {
    const x = new AbortController();
    loadFirstDelivery(x.signal)
      .then(setT)
      .catch((e) => {
        if (e.name !== "AbortError") setErr(String(e));
      });
    return () => x.abort();
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
      600,
    );
    return () => clearInterval(id);
  }, [play, t]);
  if (err) return <main className="world-error">{err}</main>;
  if (!t) return <main className="world-loading">正在加载首次商业交付…</main>;
  const f = t.frames[s];
  return (
    <div className="world-play">
      <header className="world-toolbar">
        <button className="world-back" onClick={onExit}>
          <ArrowLeft />
          产品工业化
        </button>
        <div>
          <span>PROJECT GENESIS · FIRST DELIVERY</span>
          <h1>首张正式订单商业闭环</h1>
        </div>
        <div className="world-clock">
          <small>虚拟时间</small>
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
        {f.first_commercial_cycle_closed && (
          <button onClick={() => { window.location.hash = "world-experiments"; }}>
            Scenario Lab
          </button>
        )}
      </nav>
      <main className="world-main">
        <section className="world-status">
          <span
            className={`world-state-badge ${f.first_commercial_cycle_closed ? "closed" : "active"}`}
          >
            {f.first_commercial_cycle_closed ? <BadgeCheck /> : <Truck />}
            {f.first_commercial_cycle_closed ? "Genesis cycle closed" : f.phase}
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
              <PackageCheck />
              数量与库存 <small>World owns</small>
            </header>
            <dl>
              <div>
                <dt>正式需求</dt>
                <dd>{f.demand}</dd>
              </div>
              <div>
                <dt>发运 / 接受</dt>
                <dd>
                  {f.shipped} / {f.accepted}
                </dd>
              </div>
              <div>
                <dt>可销售库存</dt>
                <dd>{f.inventory}</dd>
              </div>
              {f.shipments.map((x) => (
                <div key={x.code}>
                  <dt>{x.code}</dt>
                  <dd>
                    {x.quantity} · accepted {x.accepted}
                  </dd>
                </div>
              ))}
            </dl>
          </article>
          <article>
            <header>
              <Banknote />
              发票、应收与现金
            </header>
            <dl>
              <div>
                <dt>含税发票</dt>
                <dd>{c(f.invoice_gross.value)}</dd>
              </div>
              <div>
                <dt>应收</dt>
                <dd>{c(f.ar.value)}</dd>
              </div>
              <div>
                <dt>已核销回款</dt>
                <dd>{c(f.collected.value)}</dd>
              </div>
              <div>
                <dt>现金</dt>
                <dd>{c(f.cash.value)}</dd>
              </div>
            </dl>
          </article>
          <article>
            <header>
              <BadgeCheck />
              实际成本与毛利
            </header>
            <dl>
              <div>
                <dt>收入</dt>
                <dd>{c(f.revenue.value)}</dd>
              </div>
              <div>
                <dt>实际成本</dt>
                <dd>{c(f.actual_cost.value)}</dd>
              </div>
              <div>
                <dt>项目毛利</dt>
                <dd>{c(f.gross_margin.value)}</dd>
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
