import { Bot, CheckCircle2, Filter, ListFilter, ShieldAlert } from 'lucide-react';
import { useMemo, useState } from 'react';
import type { AgentOutput, ScenarioEvent } from '../scenario/types';

type Tab = 'events' | 'agents';
type SeverityFilter = 'all' | 'watch' | 'critical';
type DomainFilter = 'all' | ScenarioEvent['domain'];

interface InfoPanelProps {
  events: ScenarioEvent[];
  currentStep: number;
  agentOutputs: AgentOutput[];
  onSelectEvent: (step: number) => void;
}

const agentLabels = { planning: '计划 Agent', quality: '质量 Agent', business_analysis: '经营分析 Agent' } as const;
const domainLabels: Record<ScenarioEvent['domain'], string> = {
  order: '订单', planning: '计划', supply: '供应', equipment: '设备', quality: '质量', production: '生产', logistics: '物流',
};

export function InfoPanel({ events, currentStep, agentOutputs, onSelectEvent }: InfoPanelProps) {
  const [tab, setTab] = useState<Tab>('events');
  const [filter, setFilter] = useState<SeverityFilter>('all');
  const [domain, setDomain] = useState<DomainFilter>('all');
  const visibleEvents = useMemo(() => events.filter((event) =>
    (filter === 'all' || event.severity === filter) && (domain === 'all' || event.domain === domain),
  ), [domain, events, filter]);
  const visibleAgents = agentOutputs.filter((output) => output.eventSequence <= currentStep);
  const currentEvent = currentStep > 0 ? events[currentStep - 1] : undefined;

  return (
    <aside className="info-panel" aria-label="事件和 Agent 信息">
      <div className="info-tabs" role="tablist" aria-label="信息面板">
        <button role="tab" aria-selected={tab === 'events'} onClick={() => setTab('events')}><ListFilter aria-hidden="true" />事件流</button>
        <button role="tab" aria-selected={tab === 'agents'} onClick={() => setTab('agents')}><Bot aria-hidden="true" />Agent 建议<span className="count-badge">{visibleAgents.length}</span></button>
      </div>
      {tab === 'events' ? (
        <div className="info-content" role="tabpanel">
          {currentEvent && (
            <section className={`current-event-card status-${currentEvent.severity}`} aria-label="当前事件详情">
              <span>当前事件 · {domainLabels[currentEvent.domain]}</span>
              <strong>{currentEvent.title}</strong>
              <p>{currentEvent.description}</p>
            </section>
          )}
          <div className="filter-row">
            <Filter aria-hidden="true" />
            {(['all', 'watch', 'critical'] as const).map((value) => (
              <button key={value} className={filter === value ? 'active' : ''} onClick={() => setFilter(value)}>
                {value === 'all' ? '全部' : value === 'watch' ? '关注' : '严重'}
              </button>
            ))}
            <label className="domain-filter">
              <span className="sr-only">事件领域</span>
              <select aria-label="事件领域" value={domain} onChange={(event) => setDomain(event.target.value as DomainFilter)}>
                <option value="all">全部领域</option>
                {Object.entries(domainLabels).map(([value, label]) => <option key={value} value={value}>{label}</option>)}
              </select>
            </label>
          </div>
          <ol className="event-list">
            {visibleEvents.map((event) => (
              <li key={event.id}>
                <button className={`event-item status-${event.severity} ${event.sequence === currentStep ? 'current' : ''} ${event.sequence > currentStep ? 'future' : ''}`} onClick={() => onSelectEvent(event.sequence)}>
                  <span className="event-index">{String(event.sequence).padStart(2, '0')}</span>
                  <span className="event-copy"><strong>{event.title}</strong><small>{new Date(event.timestamp).toLocaleString('zh-CN', { month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit', hour12: false })} · 第 {event.act} 幕</small></span>
                  <span className="severity-label">{event.severity === 'critical' ? '严重' : event.severity === 'watch' ? '关注' : '正常'}</span>
                </button>
              </li>
            ))}
          </ol>
        </div>
      ) : (
        <div className="info-content agent-list" role="tabpanel">
          {visibleAgents.length === 0 ? (
            <div className="empty-state"><Bot aria-hidden="true" /><strong>等待业务事件</strong><span>推进故事后，Agent 会基于确定性场景给出建议。</span></div>
          ) : visibleAgents.map((output) => (
            <article className={`agent-card status-${output.status}`} key={output.id}>
              <div className="agent-card-top"><span>{agentLabels[output.agent]}</span><span className="confidence">置信度 {Math.round(output.confidence * 100)}%</span></div>
              <h3>{output.title}</h3>
              <p>{output.recommendation}</p>
              <dl><dt>依据</dt><dd>{output.evidence.join('；')}</dd><dt>影响</dt><dd>{output.impact}</dd></dl>
              <div className="agent-status">{output.status === 'executed' ? <CheckCircle2 aria-hidden="true" /> : <ShieldAlert aria-hidden="true" />}{output.status === 'executed' ? '已执行' : output.requiresHumanConfirmation ? '建议 · 待人工确认' : '建议'}</div>
            </article>
          ))}
        </div>
      )}
    </aside>
  );
}
