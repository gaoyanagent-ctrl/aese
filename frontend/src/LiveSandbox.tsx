import { AlertTriangle, LoaderCircle } from 'lucide-react';
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import type { SandboxScenario } from './scenario/types';
import type { AgentOutput, IaosScenarioEntity, IaosScenarioEvent, IaosScenarioSnapshot, RiskLevel, ScenarioEvent } from './scenario/types';
import { liveDataSourceFromEnvironment } from './scenario';
import { ControlBar } from './components/ControlBar';
import { EnterpriseNav } from './components/EnterpriseNav';
import { EntityDrawer } from './components/EntityDrawer';
import { FactoryCanvas } from './components/FactoryCanvas';
import { InfoPanel } from './components/InfoPanel';

type Connection = 'connecting' | 'connected' | 'reconnecting' | 'disconnected';
type MobileView = 'enterprise' | 'canvas' | 'events';

function risk(status: string): RiskLevel {
  if (['failed', 'down', 'blocked', 'shortage'].includes(status)) return 'critical';
  if (['delayed', 'partially_shipped', 'pending', 'maintenance'].includes(status)) return 'watch';
  return 'normal';
}

const eventLabels: Record<string, { title: string; act: number; domain: ScenarioEvent['domain'] }> = {
  'proc.production_order.released': { title: '生产工单正式下达', act: 6, domain: 'production' },
  'proc.operation.started': { title: 'A 线生产启动', act: 6, domain: 'production' },
  'proc.operation.completed': { title: '生产作业完工', act: 6, domain: 'production' },
  'whs.finished_goods.received': { title: '10,500 件成品入库', act: 6, domain: 'production' },
  'o2d.shipment.dispatched': { title: '客户订单分批发运', act: 7, domain: 'logistics' },
};

function mapEvent(event: IaosScenarioEvent, index: number, snapshot: IaosScenarioSnapshot): ScenarioEvent {
  const meta = eventLabels[event.event_type] ?? { title: event.event_type, act: 7, domain: 'logistics' as const };
  const severity = snapshot.kpis.delivery_gap.value > 0 && meta.act === 7 ? 'watch' : 'normal';
  return {
    sequence: index + 1, id: event.event_id, timestamp: event.occurred_at, eventType: event.event_type,
    title: meta.title, description: `${event.business_object_type} · ${event.business_object_code} · cursor ${event.cursor}`,
    act: meta.act, domain: meta.domain, severity, relatedEntityIds: [event.business_object_code], delta: {},
    kpis: {
      orderDemand: { value: snapshot.kpis.order_demand.value, unit: '件', risk: 'normal', trend: 'flat' },
      availableFinishedGoods: { value: snapshot.kpis.cumulative_available.value, unit: '件', risk: 'normal', trend: 'flat' },
      materialShortageRisk: { value: 0, unit: '%', risk: 'normal', trend: 'flat' },
      capacityRisk: { value: 0, unit: '%', risk: 'normal', trend: 'flat' },
      deliveryRisk: { value: snapshot.kpis.delivery_gap.value, unit: '件', risk: severity, trend: 'flat' },
    },
  };
}

function mapAgent(snapshot: IaosScenarioSnapshot): AgentOutput[] {
  return snapshot.recommendations.map((item, index) => ({
    id: `${item.agent_key}-${index}`, eventSequence: snapshot.events.length, agent: item.agent_key,
    title: item.summary, recommendation: item.recommendations.join('；'),
    evidence: [...item.object_refs, ...(item.version ? [`建议版本 v${item.version}`] : []), ...item.tool_call_ids.map((id) => `Tool Call ${id}`)],
    impact: item.data_gaps?.length ? `数据缺口：${item.data_gaps.join('、')}` : `完整度：${item.completeness}`,
    confidence: Number.parseFloat(item.confidence) || 0, status: 'suggested',
    requiresHumanConfirmation: item.requires_human_confirmation,
  }));
}

function entityForNode(nodeId: string, entities: IaosScenarioEntity[]) {
  const matchers: Record<string, (entity: IaosScenarioEntity) => boolean> = {
    'supplier-alpha': (e) => e.business_code === 'PO-202607-0001',
    'supplier-beta': (e) => e.business_code === 'PO-202607-0002',
    'incoming-quality': (e) => e.business_code === 'IQC-202607-0002',
    welding: (e) => e.business_code === 'LAS-WLD-02',
    'finished-warehouse': (e) => e.type === 'inventory' && e.business_code.includes('WH-SZ-FG'),
    shipping: (e) => e.type === 'shipment' && e.business_code === 'SHIP-202607-0002',
    customer: (e) => e.type === 'sales_order',
    forming: (e) => e.business_code === 'WO-202607-0001',
  };
  return entities.find(matchers[nodeId] ?? (() => false));
}

export function LiveSandbox({ layoutScenario, onModeChange, onOpenIntegration }: { layoutScenario: SandboxScenario; onModeChange: (mode: 'preview' | 'live') => void; onOpenIntegration: () => void }) {
  const source = useMemo(() => liveDataSourceFromEnvironment(), []);
  const [snapshot, setSnapshot] = useState<IaosScenarioSnapshot | null>(null);
  const [connection, setConnection] = useState<Connection>('connecting');
  const [error, setError] = useState<string | null>(null);
  const [generation, setGeneration] = useState(0);
  const [selectedNodeId, setSelectedNodeId] = useState<string | null>(null);
  const [mobileView, setMobileView] = useState<MobileView>('canvas');
  const cursor = useRef(0);

  const refresh = useCallback(async () => {
    setError(null);
    try {
      const next = await source.snapshot('hctm', 'order-expedite-01');
      cursor.current = next.cursor;
      setSnapshot(next);
    } catch (reason) {
      setError(reason instanceof Error ? reason.message : '在线快照加载失败');
      throw reason;
    }
  }, [source]);

  useEffect(() => {
    const abort = new AbortController();
    let retry = 0;
    let timer: ReturnType<typeof setTimeout> | undefined;
    const addEvent = (event: IaosScenarioEvent) => {
      if (event.cursor <= cursor.current) return;
      cursor.current = event.cursor;
      setSnapshot((old) => old ? { ...old, cursor: event.cursor, observed_at: event.occurred_at, events: [...old.events.filter((item) => item.event_id !== event.event_id), event] } : old);
    };
    const connect = async () => {
      try {
        setConnection(retry ? 'reconnecting' : 'connecting');
        await refresh();
        const missed = await source.eventsAfter('hctm', 'order-expedite-01', cursor.current, abort.signal);
        missed.forEach(addEvent);
        setConnection('connected');
        retry = 0;
        await source.stream('hctm', 'order-expedite-01', cursor.current, abort.signal, addEvent);
        if (!abort.signal.aborted) throw new Error('IAOS 事件流已断开');
      } catch (reason) {
        if (abort.signal.aborted) return;
        setConnection('reconnecting');
        setError(reason instanceof Error ? reason.message : '在线事件流中断');
        retry += 1;
        timer = setTimeout(connect, Math.min(1000 * 2 ** retry, 15000));
      }
    };
    void connect();
    return () => { abort.abort(); if (timer) clearTimeout(timer); setConnection('disconnected'); };
  }, [generation, refresh, source]);

  if (!snapshot && error) return <main className="error-state"><AlertTriangle /><strong>Live 模式不可用</strong><p>{error}</p><button onClick={() => setGeneration((v) => v + 1)}>重新连接</button><button onClick={() => onModeChange('preview')}>返回 Preview</button></main>;
  if (!snapshot) return <main className="loading-state"><LoaderCircle className="loading-spinner" /><strong>正在读取 IAOS 在线快照…</strong></main>;

  const events = snapshot.events.map((event, index) => mapEvent(event, index, snapshot));
  const entities = snapshot.entities.map((entity) => ({ ...entity, businessCode: entity.business_code, risk: risk(entity.status) }));
  const nodes = layoutScenario.layout.nodes.map((node) => {
    const entity = entityForNode(node.id, snapshot.entities);
    return { ...node, entityId: entity?.id, status: entity ? risk(entity.status) : node.status };
  });
  const selectedEntity = entities.find((entity) => entity.id === nodes.find((node) => node.id === selectedNodeId)?.entityId) ?? null;
  const actObserved = (act: number) => {
    const has = (type: string, status?: string) => snapshot.entities.some((entity) => entity.type === type && (!status || entity.status === status));
    if (act === 1) return has('sales_order');
    if (act === 2) return has('purchase_order');
    if (act === 3) return has('purchase_order', 'delayed');
    if (act === 4) return snapshot.entities.some((entity) => entity.type === 'equipment' && ['maintenance', 'down'].includes(entity.status));
    if (act === 5) return has('inspection_order', 'failed');
    if (act === 6) return has('finished_goods_receipt', 'posted');
    return has('shipment', 'dispatched');
  };
  const statusLabel = connection === 'connected' ? '已连接' : connection === 'reconnecting' ? '重连中' : '连接中';
  const kpis = [
    ['订单需求', snapshot.kpis.order_demand], ['累计可供', snapshot.kpis.cumulative_available],
    ['累计实发', snapshot.kpis.cumulative_shipped], ['期末成品', snapshot.kpis.ending_finished_goods], ['交付缺口', snapshot.kpis.delivery_gap],
  ] as const;

  return <div className="app-shell live-shell">
    <ControlBar scenarioName="华辰苏州基地 · 在线企业沙盘" currentTime={new Date(snapshot.observed_at).toLocaleString('zh-CN', { hour12: false })}
      step={events.length} totalSteps={events.length} playing={false} speed={1} onTogglePlay={() => {}} onPrevious={() => {}} onNext={() => {}} onReset={() => {}} onSpeedChange={() => {}}
      mode="live" onModeChange={onModeChange} sourceStatus={statusLabel} onRefresh={() => { void refresh().catch(() => undefined); }} onReconnect={() => setGeneration((v) => v + 1)} onOpenIntegration={onOpenIntegration} />
    <div className="live-integrity" role="status"><span>IAOS · tenant-hctm</span><span>游标 {snapshot.cursor}</span><span>最后更新 {new Date(snapshot.observed_at).toLocaleTimeString('zh-CN')}</span><span>完整度 {snapshot.completeness}</span>{snapshot.gaps.length > 0 && <strong>数据缺口：{snapshot.gaps.join('、')}</strong>}</div>
    <div className="mobile-tabs" role="tablist">{([['enterprise', '企业'], ['canvas', 'A 线画布'], ['events', '事件 / Agent']] as const).map(([value, label]) => <button key={value} role="tab" aria-selected={mobileView === value} className={mobileView === value ? 'active' : ''} onClick={() => setMobileView(value)}>{label}</button>)}</div>
    <main className="workspace">
      <div className={mobileView === 'enterprise' ? 'mobile-view-active enterprise-nav-wrapper' : 'enterprise-nav-wrapper'}><EnterpriseNav selectedId={selectedNodeId} onSelect={(id) => setSelectedNodeId(nodes.find((node) => node.businessCode === id)?.id ?? null)} /></div>
      <section className={`canvas-section ${mobileView === 'canvas' ? 'mobile-view-active' : ''}`}><div className="canvas-header"><span className="eyebrow">LIVE · GOVERNED IAOS FACTS</span><h2>电池冷却板 A 线</h2><p>{events.at(-1)?.title ?? '等待在线业务事件'}</p></div><FactoryCanvas nodes={nodes} edges={layoutScenario.layout.edges} selectedNodeId={selectedNodeId} highlightedNodeIds={[]} highlightedEdgeIds={[]} onNodeSelect={(node) => setSelectedNodeId(node.id)} /><div className="act-progress">{layoutScenario.acts.map((act) => <span key={act.number} className={`act-step ${actObserved(act.number) ? 'completed' : ''}`}>{act.number}. {act.title}</span>)}</div></section>
      <div className={mobileView === 'events' ? 'mobile-view-active info-panel-wrapper' : 'info-panel-wrapper'}><InfoPanel events={events} currentStep={events.length} agentOutputs={mapAgent(snapshot)} onSelectEvent={() => {}} /></div>
    </main>
    <section className="kpi-strip" aria-label="在线关键经营指标">{kpis.map(([label, metric]) => <article className={`kpi-card ${label === '交付缺口' && metric.value > 0 ? 'status-watch' : 'status-normal'}`} key={label}><div className="kpi-label"><span>{label}</span></div><strong>{new Intl.NumberFormat('zh-CN').format(metric.value)} 件</strong><span className="kpi-state">IAOS 在线事实</span></article>)}</section>
    <EntityDrawer entity={selectedEntity} onClose={() => setSelectedNodeId(null)} />
  </div>;
}
