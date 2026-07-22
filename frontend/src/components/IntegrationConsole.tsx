import {
  Activity,
  ArrowRight,
  Boxes,
  Check,
  ChevronDown,
  CircleAlert,
  Clock3,
  Copy,
  Database,
  ExternalLink,
  Factory,
  LoaderCircle,
  Play,
  PlugZap,
  Radio,
  RefreshCw,
  Settings2,
  ShoppingCart,
  SlidersHorizontal,
  Wrench,
  X,
} from 'lucide-react';
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import {
  OrchestrationError,
  OrchestrationPlan,
  OrchestrationRun,
  OrchestrationRunContext,
  OrchestrationActionResponse,
  connectAndInspectIAOS,
  defaultIntegrationConfiguration,
  createScenarioRun,
  clearStoredActiveRunId,
  clearStoredRunContext,
  executeRunAction,
  getRunPlan,
  getRunStatus,
  getStoredActiveRunId,
  getStoredRunContext,
  setStoredActiveRunId,
  setStoredRunContext,
  toRunContext,
  type IntegrationCheckResult,
  type IntegrationConfiguration,
  iaosBusinessUrl,
} from '../integration/iaosIntegration';

interface IntegrationConsoleProps {
  open: boolean;
  onClose: () => void;
  onConnected: () => void;
  onRunContextChange?: (context: OrchestrationRunContext | null) => void;
}

type IntegrationStatus = 'idle' | 'checking' | 'connected' | 'error';
type ActionBusy = string | null;
type ConsoleMode = 'check' | 'run';
type RunErrorContext = { code?: string; retryable?: boolean; requiredPermission?: string; message: string };

type ActionLog = {
  id: string;
  action: string;
  status: 'success' | 'error';
  message: string;
  detail: string;
  time: string;
  code?: string;
};

type PlanSummaryMode = 'plan' | 'run';

const ACTION_LABELS: Record<string, string> = {
  preflight: '预检',
  initialize: '初始化',
  advance: '推进下一幕',
  'run-to-end': '运行到结束',
  analyze: 'Agent 分析',
  verify: '结果验证',
  'reset-plan': '复位预览',
  reset: '执行复位',
  create: '创建运行',
};

const STAGE_LABELS: Record<string, string> = {
  preflight: '预检',
  initialize: '初始化',
  'act-1': '第1幕',
  'act-2': '第2幕',
  'act-3': '第3幕',
  'act-4': '第4幕',
  'act-5': '第5幕',
  'act-6': '第6幕',
  'act-7': '第7幕',
  analyze: 'Agent 分析',
  verify: '结果验证',
  reset: '复位',
};

const STATUS_LABELS: Record<string, string> = {
  planned: '待预检',
  initializing: '初始化中',
  ready: '可运行',
  running: '运行中',
  awaiting_analysis: '等待分析',
  analyzing: '分析中',
  awaiting_verification: '等待验收',
  completed: '已完成',
  failed: '失败',
  resetting: '复位中',
  reset: '已复位',
};

const OPERATION_ORDER = [
  'preflight',
  'initialize',
  'advance',
  'run-to-end',
  'analyze',
  'verify',
  'reset-plan',
  'reset',
];

const DANGER_ACTIONS = new Set(['reset', 'reset-plan']);

function formatTime(value: Date): string {
  return value.toLocaleTimeString('zh-CN', {
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  });
}

function buildActionIdempotencyKey(run: OrchestrationRun, action: string): string {
  const snapshot = `${run.run_id}|${action}|${run.current_act}|${run.total_acts}|${run.cursor}|${run.run_version}`;
  return btoa(unescape(encodeURIComponent(snapshot)));
}

function toErrorContext(error: unknown): RunErrorContext {
  if (error instanceof OrchestrationError) {
    return {
      code: error.code,
      retryable: error.retryable,
      requiredPermission: error.requiredPermission,
      message: error.message,
    };
  }

  if (error instanceof Error) {
    return { message: error.message };
  }

  return { message: '操作失败，请重试' };
}

function normalizeJsonText(value: unknown): string {
  if (typeof value === 'string') return value;
  if (value === undefined || value === null) return '';
  try {
    return JSON.stringify(value, null, 2);
  } catch {
    return String(value);
  }
}

function summarizeOutcome(action: string, outcome: unknown): string {
  if (!outcome || typeof outcome !== 'object') {
    return `${ACTION_LABELS[action] ?? action} 已提交`;
  }

  const record = outcome as Record<string, unknown>;
  const summary = record.summary;
  const dryRun = record.initialize_dry_run;
  const pieces: string[] = [];

  if (record.event_count !== undefined) {
    pieces.push(`events=${record.event_count}`);
  }
  if (record.warnings !== undefined) {
    const warnings = Array.isArray(record.warnings) ? record.warnings : [];
    pieces.push(`warnings=${warnings.length}`);
  }

  const appendSummary = (source: unknown, title: string) => {
    if (!source || typeof source !== 'object') return;
    const typed = source as Record<string, unknown>;
    const inserted = typeof typed.inserted === 'number' ? typed.inserted : undefined;
    const updated = typeof typed.updated === 'number' ? typed.updated : undefined;
    const deleted = typeof typed.deleted === 'number' ? typed.deleted : undefined;
    const noOp = typed.no_op === true;
    if (inserted !== undefined || updated !== undefined || deleted !== undefined || typeof noOp === 'boolean') {
      pieces.push(`${title}: inserted=${inserted ?? 0} updated=${updated ?? 0} deleted=${deleted ?? 0} no_op=${noOp}`);
    }
  };

  if (action === 'preflight' && dryRun) {
    appendSummary(dryRun, '预检');
  }
  if (summary && typeof summary === 'object') {
    appendSummary(summary, action);
  }
  if (record.reset_confirmation_token) {
    pieces.push('复位预览 token 已返回');
  }
  if (pieces.length === 0) return normalizeJsonText(outcome);
  return pieces.join('；');
}

function describeError(errorContext: RunErrorContext): string {
  const parts = [errorContext.message];
  if (errorContext.code) parts.push(`(${errorContext.code})`);
  if (errorContext.requiredPermission) {
    parts.push(`缺少权限: ${errorContext.requiredPermission}`);
  }
  if (errorContext.retryable) {
    parts.push('可重试');
  }
  return parts.join(' ');
}

function actionOutcomeText(action: string, outcome: unknown): string {
  if (!outcome || typeof outcome !== 'object') {
    return '';
  }
  const record = outcome as Record<string, unknown>;
  const actionRef = typeof record.action === 'string' ? record.action : '';
  const actor = typeof record.actor === 'string' ? record.actor : '';
  const correlation = typeof record.correlation === 'string' ? record.correlation : '';
  const chunks = [actionRef ? `action=${actionRef}` : `action=${ACTION_LABELS[action] ?? action}`];
  if (actor) chunks.push(`actor=${actor}`);
  if (correlation) chunks.push(`correlation=${correlation}`);
  const summary = summarizeOutcome(action, outcome);
  if (summary) chunks.push(summary);
  return chunks.join(' · ');
}

function extractResetToken(response: OrchestrationActionResponse | null): string {
  if (!response) return '';
  const outcome = response.run.outcome;
  if (typeof outcome !== 'object' || outcome === null) return '';
  const record = outcome as Record<string, unknown>;
  const token = record.reset_confirmation_token;
  return typeof token === 'string' ? token : '';
}

function extractResetTokenExpiresAt(response: OrchestrationActionResponse | null): string {
  if (!response) return '';
  const outcome = response.run.outcome;
  if (typeof outcome !== 'object' || outcome === null) return '';
  const record = outcome as Record<string, unknown>;
  const token = record.confirmation_expires_at;
  return typeof token === 'string' ? token : '';
}

function inferRunActionState(
  stage: string,
  run: OrchestrationRun | null,
): 'completed' | 'current' | 'future' {
  if (!run) {
    return stage === 'preflight' ? 'current' : 'future';
  }

  if (stage === 'preflight') {
    return run.status === 'planned' ? 'current' : 'completed';
  }

  if (stage === 'initialize') {
    if (run.status === 'initializing') return 'current';
    return [
      'ready',
      'running',
      'awaiting_analysis',
      'analyzing',
      'awaiting_verification',
      'completed',
      'failed',
      'resetting',
      'reset',
    ].includes(run.status)
      ? 'completed'
      : 'future';
  }

  if (stage.startsWith('act-')) {
    const actIndex = Number(stage.split('-')[1]);
    if (!Number.isNaN(actIndex)) {
      if (run.current_act > actIndex) return 'completed';
      if (
        run.current_act === actIndex &&
        ['ready', 'running', 'awaiting_analysis', 'analyzing', 'awaiting_verification', 'completed', 'failed'].includes(run.status)
      ) {
        return 'current';
      }
    }
    return 'future';
  }

  if (stage === 'analyze') {
    if (run.status === 'analyzing') return 'current';
    if (['awaiting_verification', 'completed', 'failed'].includes(run.status)) return 'completed';
    return 'future';
  }

  if (stage === 'verify') {
    if (run.status === 'awaiting_verification') return 'current';
    if (['completed', 'failed'].includes(run.status)) return 'completed';
    return 'future';
  }

  if (stage === 'reset') {
    if (run.status === 'resetting') return 'current';
    if (run.status === 'reset') return 'completed';
    if (run.status === 'failed' && run.allowed_actions.includes('reset-plan')) return 'current';
    return 'future';
  }

  return 'future';
}

const mappings = [
  { key: 'sales_order', countKey: 'sales_order', label: '销售交付', view: '客户需求与订单', icon: ShoppingCart },
  { key: 'work_order', countKey: 'work_order', label: '计划排产', view: 'A 线工单状态', icon: Factory },
  { key: 'live_inventory', countKey: 'inventory', label: '实物库存', view: '仓库与批次数量', icon: Boxes },
  { key: 'equipment', countKey: 'equipment', label: '设备', view: '焊接设备状态', icon: Wrench },
] as const;

export function IntegrationConsole({ open, onClose, onConnected, onRunContextChange }: IntegrationConsoleProps) {
  const [configuration, setConfiguration] = useState<IntegrationConfiguration>(() => defaultIntegrationConfiguration());
  const [advanced, setAdvanced] = useState(false);
  const [mode, setMode] = useState<ConsoleMode>('check');
  const [status, setStatus] = useState<IntegrationStatus>('idle');
  const [result, setResult] = useState<IntegrationCheckResult | null>(null);
  const [error, setError] = useState('');
  const [plan, setPlan] = useState<OrchestrationPlan | null>(null);
  const [planLoading, setPlanLoading] = useState(false);
  const [planError, setPlanError] = useState('');
  const [run, setRun] = useState<OrchestrationRun | null>(null);
  const [runLoading, setRunLoading] = useState(false);
  const [activeRunId, setActiveRunId] = useState<string>(() => getStoredActiveRunId());
  const [runError, setRunError] = useState('');
  const [actionBusy, setActionBusy] = useState<ActionBusy>(null);
  const [actionLogs, setActionLogs] = useState<ActionLog[]>([]);
  const [resetToken, setResetToken] = useState('');
  const [resetExpiresAt, setResetExpiresAt] = useState('');
  const [summaryMode, setSummaryMode] = useState<PlanSummaryMode>('plan');
  const actionBusyRef = useRef<ActionBusy>(null);
  const closeButton = useRef<HTMLButtonElement>(null);
  const dialog = useRef<HTMLElement>(null);
  const persistRunContext = useCallback((next: OrchestrationRun | null) => {
    if (!next) {
      clearStoredRunContext();
      onRunContextChange?.(null);
      return;
    }
    const context = toRunContext(next);
    setStoredRunContext(context);
    onRunContextChange?.(context);
  }, [onRunContextChange]);

  const addLog = useCallback((action: string, status: 'success' | 'error', context: string, detail = '', code?: string) => {
    setActionLogs((current) => [
      {
        id: `${Date.now()}-${Math.random().toString(36).slice(2, 8)}`,
        action,
        status,
        message: context,
        detail,
        time: formatTime(new Date()),
        code,
      },
      ...current,
    ]);
  }, []);

  const copyText = useCallback(async (text: string) => {
    try {
      await navigator.clipboard.writeText(text);
    } catch {
      // noop: 若浏览器不支持，保留只读显示
    }
  }, []);

  const clearRunState = useCallback(() => {
    setPlan(null);
    setRun(null);
    setActiveRunId('');
    setPlanError('');
    setRunError('');
    setResetToken('');
    setResetExpiresAt('');
    setActionLogs([]);
    clearStoredRunContext();
    onRunContextChange?.(null);
  }, [onRunContextChange]);

  const loadPlan = useCallback(async () => {
    setPlanLoading(true);
    setPlanError('');
    try {
      const next = await getRunPlan(configuration, configuration.scenarioKey);
      setPlan(next);
    } catch (reason) {
      setPlanError(toErrorContext(reason).message);
      addLog('plan', 'error', `获取计划失败：${describeError(toErrorContext(reason))}`);
    } finally {
      setPlanLoading(false);
    }
  }, [addLog, configuration]);

  const restoreActiveRun = useCallback(async (runId: string) => {
    if (!runId) return;
    setRunLoading(true);
    setRunError('');
    try {
      const next = await getRunStatus(configuration, runId);
      setRun(next);
      setActiveRunId(next.run_id);
      setStoredActiveRunId(next.run_id);
      persistRunContext(next);
      if (typeof next.reset_confirmation_required === 'boolean' && next.reset_confirmation_required === false) {
        setResetToken('');
        setResetExpiresAt('');
      }
      addLog('恢复运行', 'success', `已恢复运行 ${next.run_id}（${STATUS_LABELS[next.status] ?? next.status}）`);
    } catch (reason) {
      const context = toErrorContext(reason);
      if (context.code === 'not_found') {
        clearStoredActiveRunId();
        setActiveRunId('');
        setRun(null);
        persistRunContext(null);
        addLog('恢复运行', 'error', '本地运行记录已失效，需重新创建运行');
      } else {
        setRunError(describeError(context));
      }
    } finally {
      setRunLoading(false);
    }
  }, [addLog, configuration, persistRunContext]);

  const refreshRun = useCallback(async () => {
    if (!run?.run_id) return;
    setRunLoading(true);
    try {
      const next = await getRunStatus(configuration, run.run_id);
      setRun(next);
      persistRunContext(next);
      if (next.reset_confirmation_required === false) {
        setResetToken('');
        setResetExpiresAt('');
      }
    } catch (reason) {
      const context = toErrorContext(reason);
      if (context.code === 'not_found') {
        clearStoredActiveRunId();
        setActiveRunId('');
        setRun(null);
        persistRunContext(null);
      }
      setRunError(describeError(context));
    } finally {
      setRunLoading(false);
    }
  }, [configuration, persistRunContext, run?.run_id]);

  const runAction = useCallback(async (action: string, targetRun?: OrchestrationRun) => {
    const currentRun = targetRun ?? run;
    if (!currentRun) {
      setRunError('请先创建或恢复运行');
      return;
    }
    if (actionBusyRef.current) return;

    actionBusyRef.current = action;
    setActionBusy(action);
    setRunError('');
    const actionIdempotencyKey = buildActionIdempotencyKey(currentRun, action);
    const payload: Record<string, unknown> = {
      plan_hash: plan?.plan_hash ?? currentRun.plan_hash,
      run_version: currentRun.run_version,
      expected_cursor: currentRun.cursor,
      apply: ['preflight', 'reset-plan'].includes(action) ? false : true,
      dry_run: action === 'preflight',
      idempotency_key: actionIdempotencyKey,
    };

    if (action === 'reset') {
      payload.confirmation_token = resetToken;
    }

    try {
      const response = await executeRunAction(configuration, currentRun.run_id, action, payload);
      const next = response.run;
      setRun(next);
      persistRunContext(next);
      setStoredActiveRunId(next.run_id);
      setActiveRunId(next.run_id);
      if (response.run.status === 'reset') {
        setResetToken('');
        setResetExpiresAt('');
      } else {
        const token = extractResetToken(response);
        if (token) {
          setResetToken(token);
          setResetExpiresAt(extractResetTokenExpiresAt(response));
        } else if (action !== 'reset-plan') {
          setResetToken('');
          setResetExpiresAt('');
        }
      }
      const detail = actionOutcomeText(action, response.run.outcome);
      addLog(action, 'success', `${ACTION_LABELS[action] || action} 执行成功`, detail);
      } catch (reason) {
        const context = toErrorContext(reason);
        setRunError(describeError(context));
        addLog(action, 'error', `${ACTION_LABELS[action] || action} 执行失败`, context.message, context.code);
        if (context.code === 'plan_hash_mismatch' || context.code === 'cursor_mismatch' || context.code === 'run_version_mismatch') {
          await loadPlan();
          await refreshRun();
        }
        if (context.code === 'not_found') {
          clearStoredActiveRunId();
          setActiveRunId('');
          setRun(null);
          persistRunContext(null);
        }
      } finally {
        actionBusyRef.current = null;
        setActionBusy(null);
      }
  }, [addLog, configuration, loadPlan, plan?.plan_hash, refreshRun, resetToken, run, persistRunContext]);

  const createRun = useCallback(async () => {
    if (!plan) return;
    setActionBusy('create');
    setRunError('');
    try {
      const created = await createScenarioRun(
        configuration,
        configuration.baseUrl,
        {
          planHash: plan.plan_hash,
          actor: 'aese-user',
        },
      );
      setStoredActiveRunId(created.run_id);
      setActiveRunId(created.run_id);
      setRun(created);
      setResetToken('');
      setResetExpiresAt('');
      persistRunContext(created);
      addLog('create', 'success', `运行已创建（${created.run_id}）`, `状态：${STATUS_LABELS[created.status] ?? created.status}`, created.run_id);
      await runAction('preflight', created);
    } catch (reason) {
      const context = toErrorContext(reason);
      setRunError(describeError(context));
      addLog('create', 'error', describeError(context), `action=create failure`, context.code);
    } finally {
      setActionBusy(null);
    }
  }, [addLog, configuration, plan, runAction, persistRunContext]);

  const allowedActions = useMemo(() => {
    if (!run) return [];
    const set = new Set(run.allowed_actions);
    return OPERATION_ORDER.filter((action) => set.has(action));
  }, [run]);

  const stepper = useMemo(() => {
    if (!plan?.stages.length) return [] as { stage: string; label: string; state: 'completed' | 'current' | 'future' }[];
    return plan.stages.map((entry) => ({
      stage: entry.stage,
      label: STAGE_LABELS[entry.stage] ?? entry.stage,
      state: inferRunActionState(entry.stage, run),
    }));
  }, [plan?.stages, run]);

  const stageCountLabel = useMemo(() => {
    if (!plan) return '';
    return `共 ${plan.stages.length} 步（${plan.total_events} 业务事件）`;
  }, [plan]);

  useEffect(() => {
    if (!open) return;
    const previousFocus = document.activeElement instanceof HTMLElement ? document.activeElement : null;
    const previousOverflow = document.body.style.overflow;
    document.body.style.overflow = 'hidden';
    closeButton.current?.focus();
    const handleKey = (event: KeyboardEvent) => {
      if (event.key === 'Escape') onClose();
      if (event.key !== 'Tab') return;
      const focusable = [
        ...(
          dialog.current?.querySelectorAll<HTMLElement>('button:not([disabled]), select:not([disabled]), input:not([disabled]), a[href], [role="button"]') ?? []
        ),
      ];
      if (focusable.length === 0) return;
      const first = focusable[0];
      const last = focusable[focusable.length - 1];
      if (event.shiftKey && document.activeElement === first) {
        event.preventDefault();
        last.focus();
      } else if (!event.shiftKey && document.activeElement === last) {
        event.preventDefault();
        first.focus();
      }
    };

    window.addEventListener('keydown', handleKey);
    return () => {
      window.removeEventListener('keydown', handleKey);
      document.body.style.overflow = previousOverflow;
      previousFocus?.focus();
    };
  }, [onClose, open]);

  useEffect(() => {
    if (!open || status !== 'connected') return;
    const bootstrap = async () => {
      await loadPlan();
      const savedRunId = getStoredActiveRunId();
      setActiveRunId(savedRunId);
      if (savedRunId) {
        await restoreActiveRun(savedRunId);
      }
    };
    bootstrap();
  }, [loadPlan, open, restoreActiveRun, status]);

  useEffect(() => {
    if (!open || status !== 'connected' || !run?.run_id) return;
    const timer = setInterval(() => {
      void refreshRun();
    }, 4500);
    return () => {
      clearInterval(timer);
    };
  }, [open, refreshRun, run?.run_id, status]);

  if (!open) return null;

  const update = (key: keyof IntegrationConfiguration, value: string) => {
    setConfiguration((current) => ({ ...current, [key]: value }));
    setStatus('idle');
    setError('');
    clearRunState();
  };

  const connect = async () => {
    setStatus('checking');
    setError('');
    setPlanError('');
    setRunError('');
    try {
      const next = await connectAndInspectIAOS(configuration);
      setResult(next);
      setStatus('connected');
    } catch (reason) {
      setStatus('error');
      setError(reason instanceof Error ? reason.message : '联动检查失败');
    }
  };

  const primaryMode = mode === 'run' ? 'run' : 'check';
  const currentStatus = run?.status ? STATUS_LABELS[run.status] ?? run.status : '未运行';

  const currentRunId = run?.run_id || activeRunId;
  const canRefreshRun = Boolean(currentRunId);
  const canCreateRun = Boolean(plan) && status === 'connected' && !run;

  return (
    <div className="integration-backdrop" role="presentation" onMouseDown={(event) => event.target === event.currentTarget && onClose()}>
      <section ref={dialog} className="integration-console" role="dialog" aria-modal="true" aria-labelledby="integration-title">
        <header className="integration-header">
          <div className="integration-title-icon"><PlugZap aria-hidden="true" /></div>
          <div>
            <span className="eyebrow">AESE × IAOS</span>
            <h2 id="integration-title">企业沙盘联动中心</h2>
            <p>选择场景，一键检查 IAOS 业务数据并进入运行控制。</p>
          </div>
          <button ref={closeButton} type="button" className="icon-button integration-close" onClick={onClose} aria-label="关闭联动中心"><X aria-hidden="true" /></button>
        </header>

          <div className="integration-mode-switch" role="tablist" aria-label="联动中心视图">
            <button
              type="button"
              role="tab"
              aria-selected={primaryMode === 'check'}
              aria-label="切换到连接检查视图"
              onClick={() => setMode('check')}
            >
            连接检查
          </button>
            <button
              type="button"
              role="tab"
              aria-selected={primaryMode === 'run'}
              disabled={status !== 'connected'}
              aria-label="切换到运行场景视图"
              onClick={() => setMode('run')}
            >
            运行场景
          </button>
        </div>

        <div className="integration-body">
          <section className="integration-config" aria-labelledby="integration-config-title">
            <div className="section-title-row">
              <div>
                <span className="step-number">1</span>
                <h3 id="integration-config-title">联动对象与目标</h3>
              </div>
              <span className="safe-badge">只读检查 + 可恢复运行</span>
            </div>

            <div className="integration-fields">
              <label><span>IAOS 租户</span><select value={configuration.tenantId} onChange={(event) => update('tenantId', event.target.value)}><option value="tenant-hctm">华辰热管理系统集团</option></select></label>
              <label><span>AESE 场景</span><select value={configuration.scenarioKey} onChange={(event) => update('scenarioKey', event.target.value)}><option value="order-expedite-01">客户追加订单 · A 线交付</option></select></label>
            </div>

            <button
              type="button"
              className="advanced-toggle"
              aria-expanded={advanced}
              onClick={() => setAdvanced((value) => !value)}
            >
              <Settings2 aria-hidden="true" />
              <span>高级设置</span>
              <ChevronDown className={advanced ? 'expanded' : ''} aria-hidden="true" />
            </button>

            {advanced && (
              <div className="advanced-fields">
                <label><span>IAOS API 地址</span><input value={configuration.baseUrl} onChange={(event) => update('baseUrl', event.target.value)} placeholder="同源连接（推荐）" /><small>留空时由 AESE 安全代理到 IAOS。</small></label>
                <label><span>IAOS 工作台地址</span><input value={configuration.uiUrl} onChange={(event) => update('uiUrl', event.target.value)} /><small>用于“在 IAOS 查看”。</small></label>
                <label><span>Orchestrator 地址</span><input value={configuration.orchestratorUrl} onChange={(event) => update('orchestratorUrl', event.target.value)} /><small>默认：当前页主机 8090。</small></label>
              </div>
            )}

            <button
              type="button"
              className="integration-primary"
              disabled={status === 'checking'}
              onClick={() => void connect()}
              aria-label={status === 'checking' ? '正在检查 IAOS 与场景数据' : status === 'connected' ? '重新检查联动' : '一键连接并检查'}
            >
              {status === 'checking' ? <LoaderCircle className="loading-spinner" aria-hidden="true" /> : <PlugZap aria-hidden="true" />}
              <span>{status === 'checking' ? '正在检查 IAOS 与场景数据…' : status === 'connected' ? '重新检查联动' : '一键连接并检查'}</span>
            </button>
          </section>

          <section className="integration-results" aria-labelledby="integration-result-title">
            {mode === 'check' && <>
              <div className="section-title-row">
                <div><span className="step-number">2</span><h3 id="integration-result-title">联动状态</h3></div>
                {status === 'connected' && <span className="connected-badge"><Check aria-hidden="true" />全部可用</span>}
              </div>
              {status === 'idle' && <div className="integration-empty"><Activity aria-hidden="true" /><strong>等待联动检查</strong><p>点击按钮后，这里会显示业务对象和在线通道状态。</p></div>}
              {status === 'checking' && <div className="integration-loading" aria-live="polite">{['验证租户身份', '读取场景快照', '检查业务对象', '确认事件通道'].map((label, index) => <div key={label}><LoaderCircle className="loading-spinner" /><span>{label}</span><small>步骤 {index + 1}/4</small></div>)}</div>}
              {result && <>
                <ol className="linkage-flow" aria-label="AESE 与 IAOS 联动流程">
                  {[
                    ['连接 IAOS', '租户身份有效'],
                    ['读取业务数据', `${Object.values(result.counts).reduce((sum, count) => sum + count, 0)} 条对象`],
                    ['同步场景事件', `${result.snapshot.events.length} 个事件`],
                    ['加载 Agent 建议', `${result.snapshot.recommendations.length} 条建议`],
                  ].map(([label, detail]) => (
                    <li key={label}>
                      <span><Check aria-hidden="true" /></span>
                      <div><strong>{label}</strong><small>{detail}</small></div>
                    </li>
                  ))}
                </ol>

                <div className="integration-summary">
                  <article><Database aria-hidden="true" /><span>IAOS 租户</span><strong>{result.tenantName}</strong></article>
                  <article><Radio aria-hidden="true" /><span>场景游标</span><strong>{result.snapshot.cursor}</strong></article>
                  <article><Activity aria-hidden="true" /><span>在线完整度</span><strong>{result.snapshot.completeness}</strong></article>
                </div>

                <div className="mapping-list">
                  {mappings.map(({ key, countKey, label, view, icon: Icon }) => (
                    <article className="mapping-row" key={key}>
                      <div className="mapping-icon"><Icon aria-hidden="true" /></div>
                      <div><strong>{view}</strong><span>AESE 视图</span></div>
                      <ArrowRight aria-hidden="true" className="mapping-arrow" />
                      <div><strong>{label}</strong><span>IAOS · {result.counts[countKey]} 条</span></div>
                      <a href={iaosBusinessUrl(configuration.uiUrl, key)} target="_blank" rel="noreferrer">在 IAOS 查看<ExternalLink aria-hidden="true" /></a>
                    </article>
                  ))}
                </div>

                <div className="channel-status">
                  <Check aria-hidden="true" />
                  <div><strong>Snapshot 与事件增量通道可用</strong><span>进入 Live 后，AESE 自动读取快照并按游标补发。</span></div>
                </div>
                <button
                  type="button"
                  className="enter-live-button"
                  onClick={() => {
                    onRunContextChange?.(getStoredRunContext());
                    onConnected();
                    onClose();
                  }}
                  aria-label="进入 AESE 在线 Live 沙盘"
                >
                  进入 Live 沙盘<ArrowRight aria-hidden="true" />
                </button>
              </>}
              {status === 'error' && <div className="integration-empty error"><CircleAlert aria-hidden="true" /><strong>未建立联动</strong><p>{error || '修正连接后可直接重试，不会改变 IAOS 业务数据。'}</p></div>}
            </>}

            {mode === 'run' && <>
              <div className="section-title-row">
                <div><span className="step-number">3</span><h3 id="integration-result-title">运行场景</h3></div>
                {run && <span className="connected-badge">{currentStatus}</span>}
              </div>

              {planLoading && <div className="integration-loading"><div><LoaderCircle className="loading-spinner" /><span>正在加载运行计划…</span></div></div>}
              {planError && <div className="integration-empty error"><CircleAlert aria-hidden="true" /><strong>计划加载失败</strong><p>{planError}</p></div>}

              {plan && <div className="integration-run-summary">
                <div className="integration-summary">
                  <article><SlidersHorizontal aria-hidden="true" /><span>Pack</span><strong>{plan.pack_key}</strong></article>
                  <article><Clock3 aria-hidden="true" /><span>Plan Hash</span><strong>{plan.plan_hash.slice(0, 18)}…</strong></article>
                  <article><Check aria-hidden="true" /><span>{stageCountLabel}</span><strong>Act {plan.act_count}</strong></article>
                </div>
                <div className="integration-run-actions">
                  <button
                    type="button"
                    className="summary-switch"
                    disabled={status !== 'connected'}
                    aria-label="查看计划摘要"
                    onClick={() => setSummaryMode('plan')}
                  >
                    计划摘要
                  </button>
                  <button
                    type="button"
                    className="summary-switch"
                    disabled={!run}
                    aria-label="查看运行视图"
                    onClick={() => setSummaryMode('run')}
                  >
                    运行视图
                  </button>
                </div>
                <p className="integration-summary-text">
                  {summaryMode === 'plan' ? `plan_hash=${plan.plan_hash}` : run ? `run_id=${run.run_id}` : '尚未创建运行'}
                </p>
              </div>}

              {!run && (
                <div className="run-start-panel">
                  <p>先创建/恢复一个场景运行，再执行预检、初始化和逐幕操作。</p>
                  <button
                    type="button"
                    aria-label="创建并预检一个运行"
                    className="integration-primary"
                    disabled={planLoading || canCreateRun === false}
                    onClick={() => void createRun()}
                  >
                    {actionBusy === 'create' ? <LoaderCircle className="loading-spinner" aria-hidden="true" /> : <Check aria-hidden="true" />}
                    <span>创建并预检运行</span>
                  </button>
                  {activeRunId && (
                    <button
                      type="button"
                      className="integration-primary"
                      aria-label={`恢复本地运行 ${activeRunId}`}
                      disabled={runLoading}
                      onClick={() => void restoreActiveRun(activeRunId)}
                    >
                      {runLoading ? <LoaderCircle className="loading-spinner" aria-hidden="true" /> : <Clock3 aria-hidden="true" />}
                      <span>恢复本地运行（{activeRunId.slice(0, 10)}…）</span>
                    </button>
                  )}
                </div>
              )}

              {run && <>
                <div className="integration-run-status">
                  <article><strong>运行 ID</strong><small>{run.run_id}</small></article>
                  <article><strong>版本</strong><small>{run.run_version}</small></article>
                  <article><strong>Cursor</strong><small>{run.cursor}</small></article>
                  <article><strong>当前幕</strong><small>{run.current_act}/{run.total_acts}</small></article>
                </div>

                <div className="stepper-track" aria-label="七幕阶段状态">
                  {stepper.map(({ stage, label, state }) => (
                    <span key={stage} className={`step-chip ${state}`}>{label}</span>
                  ))}
                </div>

                <div className="integration-actions">
                  {canRefreshRun && <button
                    type="button"
                    className="integration-run-control"
                    aria-label="刷新运行状态"
                    onClick={() => void refreshRun()}
                    disabled={runLoading}
                  >
                    <RefreshCw aria-hidden="true" />
                    刷新运行状态
                  </button>}

                  {allowedActions.map((action) => {
                    const isDanger = DANGER_ACTIONS.has(action);
                    return (
                      <button
                        type="button"
                        key={action}
                        className={`integration-run-control ${isDanger ? 'run-danger' : ''}`}
                        disabled={Boolean(actionBusy) || runLoading}
                        aria-label={ACTION_LABELS[action] ?? action}
                        onClick={() => void runAction(action)}
                      >
                        {actionBusy === action ? <LoaderCircle className="loading-spinner" aria-hidden="true" /> : <Play aria-hidden="true" />}
                        <span>{ACTION_LABELS[action] ?? action}</span>
                        {action === 'reset' && run.reset_confirmation_required && <small>需确认 token</small>}
                      </button>
                    );
                  })}
                </div>

                {resetToken ? (
                  <div className="integration-reset-card">
                    <strong>复位预览已就绪</strong>
                    <span>token: {resetToken.slice(0, 12)}…</span>
                    {resetExpiresAt && <span>有效期: {resetExpiresAt}</span>}
                  </div>
                ) : <></>}

                {runError && <div className="integration-error" role="alert"><CircleAlert aria-hidden="true" /><div><strong>最近失败</strong><p>{runError}</p></div></div>}
                {run.last_error && <div className="run-last-error" aria-live="polite">{run.last_error}</div>}
              </>}

              <div className="integration-log">
                <div className="section-title-row"><div><span className="step-number">4</span><h3>执行日志</h3></div><Clock3 aria-hidden="true" /></div>
                {actionLogs.length === 0 && <div className="integration-empty"><Activity aria-hidden="true" /><strong>暂未有运行日志</strong><p>点击预检/初始化/推进/分析/验证会写入时间戳化日志。</p></div>}
                {actionLogs.length > 0 && <div className="integration-log-list">
                  {actionLogs.map((entry) => (
                    <article key={entry.id} className={`action-log ${entry.status}`}>
                      <header>
                        <span>{entry.status === 'success' ? <Check aria-hidden="true" /> : <CircleAlert aria-hidden="true" />}</span>
                        <strong>{entry.action}</strong>
                        <small>{entry.time}</small>
                        {entry.code && <small className="badge">code: {entry.code}</small>}
                      </header>
                      <p>{entry.message}</p>
                      {entry.detail && <pre className="integration-log-code">{normalizeJsonText(entry.detail)}</pre>}
                      <button
                        type="button"
                        onClick={() => void copyText(`${entry.message}\n${entry.detail}`)}
                        className="integration-copy"
                        aria-label={`复制日志：${entry.action}`}
                      >
                        <Copy aria-hidden="true" />复制日志
                      </button>
                    </article>
                  ))}
                </div>}
              </div>
            </>}
          </section>
        </div>
      </section>
    </div>
  );
}
