import {
  ArrowLeft,
  Bot,
  Building2,
  ChevronRight,
  Factory,
  FileCheck2,
  HardHat,
  Handshake,
  Landmark,
  MapPinned,
  ShieldCheck,
  Trees,
  UserRound,
  Users,
} from "lucide-react";
import { type CSSProperties, useEffect, useMemo, useState } from "react";
import { PLANT_GAME_STAGES, type PlantGameStageDefinition } from "../../world/plantGame";
import "./PlantBuildGameScene.css";

export type PlantLocation = PlantGameStageDefinition["location"];
export type PlantSceneArchiveEntry = {
  id: string;
  location: PlantLocation;
  title: string;
  state: string;
  summary: string;
  evidence: string;
};

const locations: Array<{ key: PlantLocation; label: string; shortLabel: string; description: string; kind: "内部" | "外部"; icon: typeof Building2 }> = [
  { key: "headquarters", label: "企业总部规划中心", shortLabel: "总部规划中心", description: "需求、资金边界与事实比较", kind: "内部", icon: Building2 },
  { key: "city", label: "苏州产业园区地图", shortLabel: "产业园区", description: "候选场址与调研路线", kind: "外部", icon: MapPinned },
  { key: "site", label: "候选场址与园区服务中心", shortLabel: "候选场址", description: "现场调研与外部事实", kind: "外部", icon: Factory },
  { key: "boardroom", label: "企业治理决策会议室", shortLabel: "治理会议室", description: "推荐、审批与正式决定", kind: "内部", icon: Landmark },
  { key: "contractor", label: "工程承包商市场", shortLabel: "承包商市场", description: "正式邀标、密封报价与外部资质事实", kind: "外部", icon: Handshake },
  { key: "construction", label: "工厂建设现场", shortLabel: "施工现场", description: "工程进度、质量、安全与里程碑检查", kind: "外部", icon: HardHat },
];

export function PlantBuildGameScene({ stage, caseCode, companyCode, proposalNames, observedCount, archives, onExit, onOpenTask }: {
  stage: PlantGameStageDefinition;
  caseCode: string;
  companyCode: string;
  proposalNames: string[];
  observedCount: number;
  archives: PlantSceneArchiveEntry[];
  onExit: () => void;
  onOpenTask: () => void;
}) {
  const [location, setLocation] = useState<PlantLocation>(stage.location);
  const [travelling, setTravelling] = useState(false);
  const [tab, setTab] = useState<"event" | "people" | "archive">("event");
  useEffect(() => { setLocation(stage.location); setTravelling(false); setTab("event"); }, [stage.key, stage.location]);
  const currentLocation = useMemo(() => locations.find((item) => item.key === location) ?? locations[0], [location]);
  const locationArchives = archives.filter((item) => item.location === location);
  const NPCIcon = stage.key === "investigation" || stage.key === "contract_world" || stage.key === "construction_world" ? UserRound : stage.key === "selected" || stage.key === "contract_sourcing" || stage.key === "construction_start" || stage.key === "milestone_acceptance" ? HardHat : Bot;
  const LocationIcon = currentLocation.icon;
  const enter = (next: PlantLocation) => {
    if (next === location) return;
    setTravelling(true);
    window.setTimeout(() => { setLocation(next); setTravelling(false); setTab("event"); }, 280);
  };
  const missionHere = stage.location === location;
  const openMission = () => {
    if (stage.key === "selected") {
      setLocation("site");
      setTab("archive");
      return;
    }
    onOpenTask();
  };

  return <div className="plant-game-shell">
    <header className="plant-game-header">
      <button type="button" onClick={onExit}><ArrowLeft />返回企业创生世界</button>
      <div><span>PROJECT GENESIS · M10</span><h1>{companyCode || "新制造企业"} · 工厂选址与设施建设</h1></div>
      <div className="plant-game-founder"><img src="/assets/enterprise-genesis/sprites/founder-v1.png" alt="" /><span><small>你 · 创始人</small><strong>{caseCode || "未绑定案件"}</strong></span></div>
    </header>
    <nav className="plant-game-journey" aria-label="M10 工厂规划旅程">
      {PLANT_GAME_STAGES.map((item, index) => {
        const activeIndex = PLANT_GAME_STAGES.findIndex((value) => value.key === stage.key);
        return <div key={item.key} data-state={index < activeIndex ? "completed" : index === activeIndex ? "current" : "locked"} aria-current={index === activeIndex ? "step" : undefined}><span>{index < activeIndex ? "✓" : index + 1}</span><strong>{item.label}</strong></div>;
      })}
    </nav>
    <main className="plant-game-main">
      <section className="plant-game-world" aria-label="M10 企业与外部世界">
        <div className={`plant-game-scene plant-game-scene--${location}`} role="img" aria-label={`${currentLocation.label} 2.5D 场景`}>
          <div className="plant-game-scene-shade" />
          <header className="plant-game-scene-title"><span>{currentLocation.kind}组织场景</span><h2>{currentLocation.label}</h2><p>{currentLocation.description}</p></header>
          {location === "city" && <div className="plant-game-candidate-pins" aria-label={`${proposalNames.length} 个候选场址`}>{proposalNames.slice(0, 4).map((name, index) => <button key={name} type="button" title={name} style={{ "--pin-index": index, top: index % 2 === 0 ? "30%" : "49%" } as CSSProperties}><MapPinned /><span>{name}</span></button>)}</div>}
          {location === "site" && <div className="plant-game-site-assets" aria-hidden="true"><Trees /><Factory /><HardHat /></div>}
          {location === "boardroom" && <div className="plant-game-board" aria-hidden="true"><ShieldCheck /><span>INVESTMENT & SITE GOVERNANCE</span></div>}
          {location === "contractor" && <div className="plant-game-site-assets" aria-hidden="true"><Handshake /><Building2 /><HardHat /></div>}
          {location === "construction" && <div className="plant-game-site-assets" aria-hidden="true"><HardHat /><Factory /><ShieldCheck /></div>}
          <div className={`plant-game-player ${travelling ? "travelling" : ""}`}><img src="/assets/enterprise-genesis/sprites/founder-v1.png" alt="你的创始人角色" /><span>你</span></div>
          {missionHere && <button type="button" className="plant-game-npc" onClick={openMission} aria-label={`${stage.npc}：${stage.actionLabel}`}><span className="plant-game-npc-avatar"><NPCIcon /></span><span className="plant-game-npc-bubble"><small>{stage.npcRole}</small><strong>{stage.npc}</strong><em>{stage.dialogue}</em><b>{stage.actionLabel}<ChevronRight /></b></span></button>}
          {!missionHere && <div className="plant-game-scene-quiet"><LocationIcon /><div><strong>这里目前没有待处理事件</strong><p>可检查本地点已经确认的事实档案，或前往带有事件标记的地点。</p></div></div>}
          <nav className="plant-game-location-switcher" aria-label="M10 可进入地点">{locations.map((item) => { const Icon = item.icon; const hasMission = item.key === stage.location; return <button key={item.key} type="button" className={location === item.key ? "selected" : ""} aria-pressed={location === item.key} onClick={() => enter(item.key)}><Icon /><span><strong>{item.shortLabel}</strong><small>{item.kind}</small></span>{hasMission && <i aria-label="有新事件">!</i>}</button>; })}</nav>
          {travelling && <div className="plant-game-travel" role="status" aria-live="polite"><MapPinned /><strong>正在切换组织场景…</strong></div>}
        </div>
      </section>
      <aside className="plant-game-desk">
        <div className="plant-game-place"><LocationIcon /><div><small>当前地点</small><strong>{currentLocation.label}</strong><code>{currentLocation.kind} · {locationArchives.length} 条已确认档案</code></div></div>
        <div className="plant-game-tabs" role="tablist" aria-label="M10 场景信息">
          <button role="tab" aria-selected={tab === "event"} onClick={() => setTab("event")}><FileCheck2 />当前事件</button>
          <button role="tab" aria-selected={tab === "people"} onClick={() => setTab("people")}><Users />人物</button>
          <button role="tab" aria-selected={tab === "archive"} onClick={() => setTab("archive")}><ShieldCheck />场景档案</button>
        </div>
        <div className="plant-game-panel">
          {tab === "event" && (missionHere ? <article className="plant-game-event"><span>{stage.label}</span><h2>{stage.mission}</h2><p>{stage.dialogue}</p><button type="button" onClick={openMission}>{stage.actionLabel}<ChevronRight /></button><small>正式结果只来自 IAOS Capability、审批或受信 World Observation。</small></article> : <div className="plant-game-empty"><LocationIcon /><strong>当前地点没有待处理事件</strong><p>前往带有黄色事件标记的地点，或者查看这里已经形成的档案。</p></div>)}
          {tab === "people" && <div className="plant-game-people"><article><Bot /><span><strong>纪元</strong><small>设施与工程采购 Agent · 生成建议，不批准</small></span></article><article><UserRound /><span><strong>园区顾问</strong><small>外部参与者 · 提供场址 Observation</small></span></article><article><Handshake /><span><strong>承包商市场</strong><small>外部参与者 · 提供可信投标 Observation</small></span></article><article><HardHat /><span><strong>顾远</strong><small>项目负责人 · 选择、确认与归档</small></span></article><article><ShieldCheck /><span><strong>治理审批主体</strong><small>由 IAOS 审批流确定，不由玩家指定</small></span></article></div>}
          {tab === "archive" && <div className="plant-game-archives">{locationArchives.length ? locationArchives.map((entry) => <article key={entry.id}><header><span>{entry.state}</span><strong>{entry.title}</strong></header><p>{entry.summary}</p><code>{entry.evidence}</code></article>) : <div className="plant-game-empty"><ShieldCheck /><strong>本场景尚无已确认档案</strong><p>完成受治理任务后，需求、候选、外部事实、审批和决定会按地点归档。</p></div>}</div>}
        </div>
        <footer className="plant-game-facts"><span><Factory />{proposalNames.length} 个候选</span><span><MapPinned />{observedCount} 条外部事实</span><span><ShieldCheck />{archives.length} 条档案</span></footer>
      </aside>
    </main>
  </div>;
}
