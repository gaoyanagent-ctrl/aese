import { AlertTriangle, LoaderCircle } from 'lucide-react';
import { useEffect, useMemo, useState } from 'react';
import previewJson from '../../scenario-packs/hctm/stories/order-expedite-01/preview.json';
import { ControlBar } from './components/ControlBar';
import { EnterpriseNav } from './components/EnterpriseNav';
import { EntityDrawer } from './components/EntityDrawer';
import { FactoryCanvas } from './components/FactoryCanvas';
import { InfoPanel } from './components/InfoPanel';
import { KpiStrip } from './components/KpiStrip';
import { IntegrationConsole } from './components/IntegrationConsole';
import { getStoredRunContext, type OrchestrationRunContext } from './integration/iaosIntegration';
import { usePlayback } from './playback';
import { StaticScenarioDataSource } from './scenario';
import type { SandboxScenario } from './scenario/types';
import { LiveSandbox } from './LiveSandbox';
import { SystemAtlas, type AtlasEntryRef } from './components/SystemAtlas';
import { WorldPlay } from './components/world/WorldPlay';
import { IncorporationPlay } from './components/world/IncorporationPlay';

const SCENARIO_KEY = 'order-expedite-01';
const scenarioSource = new StaticScenarioDataSource({ [SCENARIO_KEY]: previewJson });
type MobileView = 'enterprise' | 'canvas' | 'events';

function Sandbox({ scenario, onModeChange, onOpenIntegration, onOpenAtlas, onOpenWorld }: { scenario: SandboxScenario; onModeChange: (mode: 'preview' | 'live') => void; onOpenIntegration: () => void; onOpenAtlas: () => void; onOpenWorld: () => void }) {
  const playback = usePlayback(scenario);
  const [selectedEntityId, setSelectedEntityId] = useState<string | null>(null);
  const [selectedNodeId, setSelectedNodeId] = useState<string | null>(null);
  const [mobileView, setMobileView] = useState<MobileView>('canvas');
  const currentEvent = playback.currentEvent;

  const nodes = useMemo(() => scenario.layout.nodes.map((node) => ({
    ...node,
    status: playback.viewState.nodeStatuses[node.id] ?? node.status,
  })), [playback.viewState.nodeStatuses, scenario.layout.nodes]);
  const edges = useMemo(() => scenario.layout.edges.map((edge) => ({
    ...edge,
    status: playback.viewState.edgeStatuses[edge.id] ?? edge.status,
  })), [playback.viewState.edgeStatuses, scenario.layout.edges]);
  const selectedEntity = playback.viewState.entities.find((entity) => entity.id === selectedEntityId) ?? null;
  const highlightedNodeIds = nodes
    .filter((node) => currentEvent?.relatedEntityIds.includes(node.entityId ?? node.id))
    .map((node) => node.id);
  const highlightedEdgeIds = edges
    .filter((edge) => highlightedNodeIds.includes(edge.source) || highlightedNodeIds.includes(edge.target))
    .map((edge) => edge.id);
  const currentTime = currentEvent?.timestamp ?? scenario.startsAt;
  const currentAct = currentEvent?.act ?? 0;

  const handleNodeSelect = (node: (typeof nodes)[number]) => {
    setSelectedNodeId(node.id);
    setSelectedEntityId(node.entityId ?? null);
  };

  const handleEnterpriseSelect = (id: string) => {
    const node = nodes.find((candidate) => candidate.businessCode === id || candidate.entityId === id);
    if (node) handleNodeSelect(node);
  };

  return (
    <div className="app-shell">
      <ControlBar
        scenarioName={scenario.name}
        currentTime={new Date(currentTime).toLocaleString('zh-CN', { month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit', hour12: false, timeZone: scenario.timezone })}
        step={playback.currentStep}
        totalSteps={playback.totalSteps}
        playing={playback.status === 'playing'}
        speed={playback.speed}
        onTogglePlay={playback.controls.toggle}
        onPrevious={playback.controls.previous}
        onNext={playback.controls.next}
        onReset={playback.controls.reset}
        onSpeedChange={playback.controls.setSpeed}
        mode="preview"
        onModeChange={onModeChange}
        onOpenIntegration={onOpenIntegration}
        onOpenAtlas={onOpenAtlas}
        onOpenWorld={onOpenWorld}
      />
      <div className="mobile-tabs" role="tablist" aria-label="移动端沙盘区域">
        {([['enterprise', '企业'], ['canvas', 'A 线画布'], ['events', '事件 / Agent']] as const).map(([value, label]) => (
          <button key={value} className={mobileView === value ? 'active' : ''} role="tab" aria-selected={mobileView === value} onClick={() => setMobileView(value)}>{label}</button>
        ))}
      </div>
      <main id="sandbox-main" className="workspace">
        <div className={mobileView === 'enterprise' ? 'mobile-view-active enterprise-nav-wrapper' : 'enterprise-nav-wrapper'}>
          <EnterpriseNav selectedId={selectedNodeId} onSelect={handleEnterpriseSelect} />
        </div>
        <section className={`canvas-section ${mobileView === 'canvas' ? 'mobile-view-active' : ''}`} aria-labelledby="canvas-title">
          <div className="canvas-header">
            <span className="eyebrow">SUZHOU · BATTERY COOLING PLATE</span>
            <h2 id="canvas-title">电池冷却板 A 线</h2>
            <p>{currentEvent ? `当前：${currentEvent.title}` : '初始状态 · 等待订单故事开始'}</p>
          </div>
          <FactoryCanvas
            nodes={nodes}
            edges={edges}
            selectedNodeId={selectedNodeId}
            highlightedNodeIds={highlightedNodeIds}
            highlightedEdgeIds={highlightedEdgeIds}
            onNodeSelect={handleNodeSelect}
          />
          <div className="act-progress" aria-label={`七幕故事进度，当前第 ${currentAct || 0} 幕`}>
            {scenario.acts.map((act) => (
              <button
                key={act.number}
                className={`act-step ${currentAct > act.number ? 'completed' : ''} ${currentAct === act.number ? 'current' : ''}`}
                onClick={() => playback.controls.seek(act.eventRange[0])}
                aria-label={`跳转到第 ${act.number} 幕：${act.title}`}
              >
                {act.number}. {act.title}
              </button>
            ))}
          </div>
        </section>
        <div className={mobileView === 'events' ? 'mobile-view-active info-panel-wrapper' : 'info-panel-wrapper'}>
          <InfoPanel events={scenario.timeline} currentStep={playback.currentStep} agentOutputs={scenario.agentOutputs} onSelectEvent={playback.controls.seek} />
        </div>
      </main>
      <KpiStrip kpis={playback.viewState.kpis} />
      <EntityDrawer entity={selectedEntity} onClose={() => { setSelectedEntityId(null); setSelectedNodeId(null); }} />
    </div>
  );
}

export default function App() {
  const [scenario, setScenario] = useState<SandboxScenario | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [mode, setMode] = useState<'preview' | 'live' | 'atlas' | 'world' | 'world-incorporation'>('preview');
  const [integrationOpen, setIntegrationOpen] = useState(false);
  const [connectionVersion, setConnectionVersion] = useState(0);
  const [runContext, setRunContext] = useState<OrchestrationRunContext | null>(() => getStoredRunContext());

  const navigate = (target: 'preview' | 'live' | 'atlas' | 'world' | 'world-incorporation') => {
    setMode(target);
    window.location.hash = target === 'preview' ? 'sandbox' : target;
  };

  useEffect(() => {
    const applyHash = () => {
      const target = window.location.hash.replace(/^#/, '');
      if (target === 'atlas') setMode('atlas');
      if (target === 'live') setMode('live');
      if (target === 'world') setMode('world');
      if (target === 'world-incorporation') setMode('world-incorporation');
      if (target === 'sandbox') setMode('preview');
      if (target === 'integration') { setMode('preview'); setIntegrationOpen(true); }
    };
    applyHash();
    window.addEventListener('hashchange', applyHash);
    return () => window.removeEventListener('hashchange', applyHash);
  }, []);

  useEffect(() => {
    let active = true;
    scenarioSource.loadScenario(SCENARIO_KEY)
      .then((loaded) => active && setScenario(loaded))
      .catch((reason: unknown) => active && setError(reason instanceof Error ? reason.message : '未知数据错误'));
    return () => { active = false; };
  }, []);

  if (error) return <main className="error-state"><AlertTriangle aria-hidden="true" /><strong>场景无法加载</strong><p>{error}</p><p>请检查 preview.json 是否存在且包含七幕和 22 个事件。</p></main>;
  if (!scenario) return <main className="loading-state"><LoaderCircle aria-hidden="true" className="loading-spinner" /><strong>正在装载苏州基地场景…</strong></main>;
  const openIntegration = () => setIntegrationOpen(true);
  const navigateAtlasEntry = (entry: AtlasEntryRef) => {
    const target = entry.path.replace(/^#/, '');
    if (target === 'integration') { navigate('preview'); setIntegrationOpen(true); return; }
    navigate(target === 'live' ? 'live' : target === 'atlas' ? 'atlas' : 'preview');
  };
  if (mode === 'atlas') return <SystemAtlas onExit={() => navigate('preview')} onNavigate={navigateAtlasEntry} />;
  if (mode === 'world') return <WorldPlay onExit={() => navigate('preview')} />;
  if (mode === 'world-incorporation') return <IncorporationPlay onExit={() => navigate('world')} />;
  return <>
    {mode === 'live'
      ? <LiveSandbox
          key={connectionVersion}
          layoutScenario={scenario}
          runContext={runContext}
          onModeChange={navigate}
          onOpenIntegration={openIntegration}
          onOpenAtlas={() => navigate('atlas')}
        />
      : <Sandbox scenario={scenario} onModeChange={navigate} onOpenIntegration={openIntegration} onOpenAtlas={() => navigate('atlas')} onOpenWorld={() => navigate('world')} />}
    <IntegrationConsole
      open={integrationOpen}
      onClose={() => setIntegrationOpen(false)}
      onConnected={() => { setConnectionVersion((version) => version + 1); setMode('live'); }}
      onRunContextChange={setRunContext}
    />
  </>;
}
