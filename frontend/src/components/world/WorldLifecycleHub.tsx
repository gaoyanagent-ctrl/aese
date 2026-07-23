import {
  ArrowLeft,
  ArrowRight,
  BadgeCheck,
  Building2,
  Factory,
  FlaskConical,
  Network,
  PackageCheck,
  Radar,
  ShieldCheck,
  Sparkles,
  Users,
  Wrench,
} from "lucide-react";
import "./WorldLifecycleHub.css";

type Stage = {
  milestone: string;
  title: string;
  summary: string;
  exchange: string;
  href: string;
  entryLabel: string;
  icon: typeof Building2;
  kind?: "foundation" | "delivery" | "governance";
};

const stages: Stage[] = [
  { milestone: "M8", title: "三态世界与 IAOS Bridge", summary: "验证客观世界、IAOS 管理事实与角色认知如何发现并关闭偏差。", exchange: "Observation → Journal → Intent → Outcome", href: "#world-tristate", entryLabel: "进入 M8 三态架构验证", icon: Network, kind: "foundation" },
  { milestone: "M9", title: "公司成立与治理", summary: "从投资人出资开始，完成法人、治理岗位、组织和初始预算。", exchange: "World 经济事实 ↔ IAOS 治理记录", href: "#world-incorporation", entryLabel: "进入 M9 公司成立", icon: Building2 },
  { milestone: "M10", title: "工厂选址与设施建设", summary: "比较选址方案，推进场地、公用工程、项目付款与设施验收。", exchange: "工程进展 ↔ 项目、预算与审批", href: "#world-plant-build", entryLabel: "进入 M10 工厂建设", icon: Factory },
  { milestone: "M11", title: "生产能力建设", summary: "完成设备、实验室、仓储、人员招聘、培训和岗位资格。", exchange: "设备/人员状态 ↔ 采购、资产与资格", href: "#world-capability-build", entryLabel: "进入 M11 能力建设", icon: Wrench },
  { milestone: "M12", title: "产品工业化", summary: "贯通 RFQ、产品与工艺设计、APQP、试制、PPAP 和量产批准。", exchange: "试制与质量结果 ↔ 工程版本与放行", href: "#world-industrialization", entryLabel: "进入 M12 产品工业化", icon: FlaskConical },
  { milestone: "M13", title: "第一次商业交付", summary: "从正式订单走到采购、生产、质量、三批交付、开票、回款和毛利。", exchange: "物理履约 ↔ O2D、财务与客户接受", href: "#world-first-delivery", entryLabel: "进入 M13 商业交付", icon: PackageCheck, kind: "delivery" },
  { milestone: "M14", title: "经营策略实验", summary: "用共同随机数比较 baseline、lean 与 resilient，形成不可自动投放的证据。", exchange: "隔离 World 分支 → EvidenceBundle", href: "#world-experiments", entryLabel: "进入 M14 策略实验", icon: Radar, kind: "governance" },
  { milestone: "M15", title: "策略发布与试点", summary: "经过独立审议、shadow、有限 pilot、监控、回滚和补偿后决定采纳。", exchange: "Evidence → StrategyRelease → Pilot outcome", href: "#world-strategy-control", entryLabel: "进入 M15 策略治理", icon: ShieldCheck, kind: "governance" },
  { milestone: "M16", title: "持续策略保障", summary: "观察后续经营数据，检查漂移、校准假设并复审已采纳策略。", exchange: "Canonical observations → AssuranceDecision", href: "#world-assurance", entryLabel: "进入 M16 策略保障", icon: BadgeCheck, kind: "governance" },
  { milestone: "M17–M24", title: "企业经营扩展与平台化", summary: "覆盖 IBP、组合、多基地、售后质量、EHS、集团价值、多 Agent 和场景产品化。", exchange: "八个受治理 terminal 与 exact evidence", href: "#world-aese3", entryLabel: "进入 M17 到 M24 经营扩展", icon: Users, kind: "governance" },
];

export function WorldLifecycleHub({ onExit }: { onExit: () => void }) {
  return (
    <div className="lifecycle-hub">
      <header className="lifecycle-hero">
        <button className="lifecycle-back" onClick={onExit}><ArrowLeft aria-hidden="true" />返回企业沙盘</button>
        <div className="lifecycle-hero-copy">
          <span className="lifecycle-eyebrow"><Sparkles aria-hidden="true" />PROJECT GENESIS · WORLD OPERATIONS</span>
          <h1>企业生命周期运营中心</h1>
          <p>从公司成立到产品交付与持续经营</p>
        </div>
        <div className="lifecycle-boundary"><strong>AESE World</strong><span>物理与经济后果</span><i aria-hidden="true">↔</i><strong>IAOS</strong><span>业务事实与治理</span></div>
      </header>
      <main className="lifecycle-content" id="lifecycle-content">
        <section className="lifecycle-overview" aria-labelledby="lifecycle-overview-title">
          <div><span>完整主路径</span><h2 id="lifecycle-overview-title">建立企业 → 建成工厂 → 形成能力 → 产品量产 → 完成交付</h2></div>
          <p>选择任一阶段查看工作步骤、World/IAOS 交换、审批边界和终态证据。M8 是支撑整条路径的架构验证，不再作为默认业务首页。</p>
        </section>
        <ol className="lifecycle-stages" aria-label="M8 到 M24 企业生命周期">
          {stages.map((stage, index) => {
            const Icon = stage.icon;
            return <li className={`lifecycle-stage ${stage.kind ?? "build"}`} key={stage.milestone}>
              <div className="lifecycle-stage-index"><span>{stage.milestone}</span><small>{String(index + 1).padStart(2, "0")}</small></div>
              <div className="lifecycle-stage-icon"><Icon aria-hidden="true" /></div>
              <div className="lifecycle-stage-copy"><h3>{stage.title}</h3><p>{stage.summary}</p><div className="lifecycle-exchange"><span>交换</span>{stage.exchange}</div></div>
              <a className="lifecycle-enter" href={stage.href} aria-label={stage.entryLabel}>查看过程<ArrowRight aria-hidden="true" /></a>
            </li>;
          })}
        </ol>
      </main>
    </div>
  );
}

const journey = [
  ["总览", "world"], ["M8", "world-tristate"], ["M9", "world-incorporation"], ["M10", "world-plant-build"], ["M11", "world-capability-build"], ["M12", "world-industrialization"], ["M13", "world-first-delivery"], ["M14", "world-experiments"], ["M15", "world-strategy-control"], ["M16", "world-assurance"], ["M17–24", "world-aese3"],
] as const;

export function WorldJourneyBar({ current }: { current: string }) {
  return <nav className="world-journey-bar" aria-label="企业生命周期快速导航">
    {journey.map(([label, route]) => <a key={route} href={`#${route}`} aria-current={current === route ? "page" : undefined}>{label}</a>)}
  </nav>;
}
