import {
  ArrowRight,
  Bot,
  Building2,
  CircleCheck,
  FlaskConical,
  Landmark,
  Play,
  ShieldCheck,
  Sparkles,
  Users,
  UserRound,
} from "lucide-react";
import "./GenesisHome.css";
import type { ReactNode } from "react";

interface GenesisHomeProps {
  onCreateEnterprise: () => void;
  onOpenDemo: () => void;
  onOpenWorld: () => void;
  username?: string;
  onSignOut?: () => void;
  children?: ReactNode;
}

const creationSteps = [
  { icon: Sparkles, title: "定义企业身份", copy: "AI 辅助命名、品牌与经营方向" },
  { icon: Building2, title: "办理企业设立", copy: "执照、治理结构与法定事项" },
  { icon: Users, title: "组建核心团队", copy: "聘任 CEO 与数字员工 Agent" },
  { icon: Landmark, title: "启动企业经营", copy: "银行开户、注资与预算启用" },
];

export function GenesisHome({
  onCreateEnterprise,
  onOpenDemo,
  onOpenWorld,
  username,
  onSignOut,
  children,
}: GenesisHomeProps) {
  return (
    <main className="genesis-home">
      <header className="genesis-home__nav">
        <a className="genesis-home__brand" href="/" aria-label="Enterprise Genesis 首页">
          <span className="genesis-home__brand-mark">G</span>
          <span>
            <strong>Enterprise Genesis</strong>
            <small>IAOS × AESE 企业创生世界</small>
          </span>
        </a>
        <div className="genesis-home__nav-actions">
          <span className="genesis-home__status">
            <i aria-hidden="true" />
            本地开发环境
          </span>
          <button className="genesis-home__text-button" onClick={onOpenWorld}>
            世界地图
          </button>
          {username&&<span className="genesis-home__user"><UserRound aria-hidden="true"/>{username}</span>}
          {onSignOut&&<button className="genesis-home__text-button" onClick={onSignOut}>退出</button>}
        </div>
      </header>

      <section className="genesis-home__hero" aria-labelledby="genesis-home-title">
        <div className="genesis-home__hero-copy">
          <p className="genesis-home__kicker">
            <span>ENTERPRISE CREATION SIMULATION</span>
            <span>从零开始</span>
          </p>
          <h1 id="genesis-home-title">从一个想法，创建一家真正运行的企业</h1>
          <p className="genesis-home__lead">
            你是企业创建者。提出构想、作出决策、雇佣数字员工，并让每一步经营动作
            通过 IAOS 成为可管理、可审计的企业事实。
          </p>
          <div className="genesis-home__hero-actions">
            <button
              className="genesis-home__primary"
              onClick={onCreateEnterprise}
            >
              创建新企业
              <ArrowRight aria-hidden="true" />
            </button>
            <button className="genesis-home__secondary" onClick={onOpenDemo}>
              <Play aria-hidden="true" />
              体验华辰样板世界
            </button>
          </div>
          <div className="genesis-home__trust">
            <span><ShieldCheck aria-hidden="true" /> IAOS 受治理业务运行时</span>
            <span><Bot aria-hidden="true" /> AI 创意与数字员工</span>
            <span><CircleCheck aria-hidden="true" /> 全过程证据链</span>
          </div>
        </div>

        <div className="genesis-home__scene" aria-label="企业创生过程概览">
          <div className="genesis-home__orb genesis-home__orb--one" />
          <div className="genesis-home__orb genesis-home__orb--two" />
          <article className="genesis-home__command-card">
            <div className="genesis-home__command-head">
              <span>创始人指挥中心</span>
              <span className="genesis-home__live">READY</span>
            </div>
            <div className="genesis-home__company">
              <span className="genesis-home__company-logo">新</span>
              <div>
                <small>你的下一家公司</small>
                <strong>等待被创造</strong>
              </div>
            </div>
            <ol className="genesis-home__mini-flow">
              <li className="active"><span>01</span><div><strong>企业构想</strong><small>由你发起</small></div></li>
              <li><span>02</span><div><strong>AI 身份工作室</strong><small>名称与品牌</small></div></li>
              <li><span>03</span><div><strong>IAOS 企业设立</strong><small>流程与数据</small></div></li>
              <li><span>04</span><div><strong>开始经营</strong><small>人、财、事协同</small></div></li>
            </ol>
          </article>
          <div className="genesis-home__agent-card">
            <Bot aria-hidden="true" />
            <span><strong>AI 创意官</strong><small>准备接收你的企业构想</small></span>
          </div>
        </div>
      </section>

      {children}

      <section className="genesis-home__journey" aria-labelledby="creation-journey-title">
        <div className="genesis-home__section-heading">
          <div>
            <p>YOUR FIRST MISSION</p>
            <h2 id="creation-journey-title">M9 · 企业从 0 到 1</h2>
          </div>
          <span>4 个阶段 · 由玩家决策驱动</span>
        </div>
        <div className="genesis-home__steps">
          {creationSteps.map(({ icon: Icon, title, copy }, index) => (
            <article key={title}>
              <span className="genesis-home__step-number">0{index + 1}</span>
              <Icon aria-hidden="true" />
              <h3>{title}</h3>
              <p>{copy}</p>
            </article>
          ))}
        </div>
      </section>

      <section className="genesis-home__lower">
        <article className="genesis-home__feature">
          <FlaskConical aria-hidden="true" />
          <div>
            <p>样板世界</p>
            <h2>先看看一家企业如何运转</h2>
            <span>进入华辰热管理样板场景，体验订单、产线、Agent 与 IAOS 的联动。</span>
          </div>
          <button onClick={onOpenDemo}>进入演示 <ArrowRight aria-hidden="true" /></button>
        </article>
        <article className="genesis-home__build-note">
          <span>运行状态</span>
          <strong>独立租户创建通道已启用</strong>
          <p>每家新企业获得独立 IAOS 租户；可从“我的企业”返回并继续原有创生进度。</p>
        </article>
      </section>
    </main>
  );
}
