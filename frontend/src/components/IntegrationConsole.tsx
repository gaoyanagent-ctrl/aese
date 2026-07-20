import {
  Activity,
  ArrowRight,
  Boxes,
  Check,
  ChevronDown,
  CircleAlert,
  Database,
  ExternalLink,
  Factory,
  LoaderCircle,
  PlugZap,
  Radio,
  Settings2,
  ShoppingCart,
  Wrench,
  X,
} from 'lucide-react';
import { useEffect, useRef, useState } from 'react';
import {
  connectAndInspectIAOS,
  defaultIntegrationConfiguration,
  iaosBusinessUrl,
  type IntegrationCheckResult,
  type IntegrationConfiguration,
} from '../integration/iaosIntegration';

interface IntegrationConsoleProps {
  open: boolean;
  onClose: () => void;
  onConnected: () => void;
}

const mappings = [
  { key: 'sales_order', countKey: 'sales_order', label: '销售交付', view: '客户需求与订单', icon: ShoppingCart },
  { key: 'work_order', countKey: 'work_order', label: '计划排产', view: 'A 线工单状态', icon: Factory },
  { key: 'live_inventory', countKey: 'inventory', label: '实物库存', view: '仓库与批次数量', icon: Boxes },
  { key: 'equipment', countKey: 'equipment', label: '设备', view: '焊接设备状态', icon: Wrench },
] as const;

export function IntegrationConsole({ open, onClose, onConnected }: IntegrationConsoleProps) {
  const [configuration, setConfiguration] = useState<IntegrationConfiguration>(() => defaultIntegrationConfiguration());
  const [advanced, setAdvanced] = useState(false);
  const [status, setStatus] = useState<'idle' | 'checking' | 'connected' | 'error'>('idle');
  const [result, setResult] = useState<IntegrationCheckResult | null>(null);
  const [error, setError] = useState('');
  const closeButton = useRef<HTMLButtonElement>(null);
  const dialog = useRef<HTMLElement>(null);

  useEffect(() => {
    if (!open) return;
    const previousFocus = document.activeElement instanceof HTMLElement ? document.activeElement : null;
    const previousOverflow = document.body.style.overflow;
    document.body.style.overflow = 'hidden';
    closeButton.current?.focus();
    const handleKey = (event: KeyboardEvent) => {
      if (event.key === 'Escape') onClose();
      if (event.key !== 'Tab') return;
      const focusable = [...(dialog.current?.querySelectorAll<HTMLElement>('button:not([disabled]), select:not([disabled]), input:not([disabled]), a[href]') ?? [])];
      if (focusable.length === 0) return;
      const first = focusable[0];
      const last = focusable[focusable.length - 1];
      if (event.shiftKey && document.activeElement === first) { event.preventDefault(); last.focus(); }
      if (!event.shiftKey && document.activeElement === last) { event.preventDefault(); first.focus(); }
    };
    window.addEventListener('keydown', handleKey);
    return () => {
      window.removeEventListener('keydown', handleKey);
      document.body.style.overflow = previousOverflow;
      previousFocus?.focus();
    };
  }, [onClose, open]);

  if (!open) return null;

  const update = (key: keyof IntegrationConfiguration, value: string) => {
    setConfiguration((current) => ({ ...current, [key]: value }));
    setStatus('idle');
    setResult(null);
  };
  const connect = async () => {
    setStatus('checking');
    setError('');
    try {
      const next = await connectAndInspectIAOS(configuration);
      setResult(next);
      setStatus('connected');
    } catch (reason) {
      setStatus('error');
      setError(reason instanceof Error ? reason.message : '联动检查失败');
    }
  };

  return <div className="integration-backdrop" role="presentation" onMouseDown={(event) => event.target === event.currentTarget && onClose()}>
    <section ref={dialog} className="integration-console" role="dialog" aria-modal="true" aria-labelledby="integration-title">
      <header className="integration-header">
        <div className="integration-title-icon"><PlugZap aria-hidden="true" /></div>
        <div>
          <span className="eyebrow">AESE × IAOS</span>
          <h2 id="integration-title">企业沙盘联动中心</h2>
          <p>选择场景，一键检查 IAOS 业务数据并进入 Live。</p>
        </div>
        <button ref={closeButton} className="icon-button integration-close" onClick={onClose} aria-label="关闭联动中心"><X aria-hidden="true" /></button>
      </header>

      <div className="integration-body">
        <section className="integration-config" aria-labelledby="connection-config-title">
          <div className="section-title-row">
            <div><span className="step-number">1</span><h3 id="connection-config-title">选择联动对象</h3></div>
            <span className="safe-badge">只读检查</span>
          </div>
          <div className="integration-fields">
            <label><span>IAOS 租户</span><select value={configuration.tenantId} onChange={(event) => update('tenantId', event.target.value)}><option value="tenant-hctm">华辰热管理系统集团</option></select></label>
            <label><span>AESE 场景</span><select value={configuration.scenarioKey} onChange={(event) => update('scenarioKey', event.target.value)}><option value="order-expedite-01">客户追加订单 · A 线交付</option></select></label>
          </div>
          <button className="advanced-toggle" aria-expanded={advanced} onClick={() => setAdvanced((value) => !value)}><Settings2 aria-hidden="true" /><span>高级连接设置</span><ChevronDown className={advanced ? 'expanded' : ''} aria-hidden="true" /></button>
          {advanced && <div className="advanced-fields">
            <label><span>IAOS API 地址</span><input value={configuration.baseUrl} onChange={(event) => update('baseUrl', event.target.value)} placeholder="同源连接（推荐）" /><small>留空时由 AESE 安全代理到 IAOS。</small></label>
            <label><span>IAOS 工作台地址</span><input value={configuration.uiUrl} onChange={(event) => update('uiUrl', event.target.value)} /><small>用于“在 IAOS 查看”跳转。</small></label>
          </div>}
          <button className="integration-primary" disabled={status === 'checking'} onClick={() => void connect()}>
            {status === 'checking' ? <LoaderCircle className="loading-spinner" aria-hidden="true" /> : <PlugZap aria-hidden="true" />}
            <span>{status === 'checking' ? '正在检查 IAOS 与场景数据…' : status === 'connected' ? '重新检查联动' : '一键连接并检查'}</span>
          </button>
          {status === 'error' && <div className="integration-error" role="alert"><CircleAlert aria-hidden="true" /><div><strong>联动失败</strong><p>{error}</p><span>请确认 IAOS 服务已启动，或展开高级设置检查地址。</span></div></div>}
        </section>

        <section className="integration-results" aria-labelledby="integration-result-title">
          <div className="section-title-row"><div><span className="step-number">2</span><h3 id="integration-result-title">联动状态</h3></div>{status === 'connected' && <span className="connected-badge"><Check aria-hidden="true" />全部可用</span>}</div>
          {status === 'idle' && <div className="integration-empty"><Activity aria-hidden="true" /><strong>等待联动检查</strong><p>点击左侧按钮后，这里会显示业务对象和在线通道状态。</p></div>}
          {status === 'checking' && <div className="integration-loading" aria-live="polite">{['验证租户身份', '读取场景快照', '检查业务对象', '确认事件通道'].map((label, index) => <div key={label}><LoaderCircle className="loading-spinner" /><span>{label}</span><small>步骤 {index + 1}/4</small></div>)}</div>}
          {result && <>
            <div className="integration-summary">
              <article><Database aria-hidden="true" /><span>IAOS 租户</span><strong>{result.tenantName}</strong></article>
              <article><Radio aria-hidden="true" /><span>场景游标</span><strong>{result.snapshot.cursor}</strong></article>
              <article><Activity aria-hidden="true" /><span>在线完整度</span><strong>{result.snapshot.completeness}</strong></article>
            </div>
            <div className="mapping-list">
              {mappings.map(({ key, countKey, label, view, icon: Icon }) => <article className="mapping-row" key={key}>
                <div className="mapping-icon"><Icon aria-hidden="true" /></div>
                <div><strong>{view}</strong><span>AESE 视图</span></div>
                <ArrowRight aria-hidden="true" className="mapping-arrow" />
                <div><strong>{label}</strong><span>IAOS · {result.counts[countKey]} 条</span></div>
                <a href={iaosBusinessUrl(configuration.uiUrl, key)} target="_blank" rel="noreferrer">在 IAOS 查看<ExternalLink aria-hidden="true" /></a>
              </article>)}
            </div>
            <div className="channel-status"><Check aria-hidden="true" /><div><strong>Snapshot 与事件增量通道可用</strong><span>进入 Live 后，AESE 会自动读取快照、按游标补发并保持连接。</span></div></div>
            <button className="enter-live-button" onClick={() => { onConnected(); onClose(); }}>进入 Live 沙盘<ArrowRight aria-hidden="true" /></button>
          </>}
          {status === 'error' && <div className="integration-empty error"><CircleAlert aria-hidden="true" /><strong>未建立联动</strong><p>修正连接后可以直接重试，不会改变 IAOS 业务数据。</p></div>}
        </section>
      </div>
    </section>
  </div>;
}
