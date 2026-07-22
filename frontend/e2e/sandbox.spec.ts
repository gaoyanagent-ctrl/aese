import { expect, test } from '@playwright/test';

test('renders a non-empty factory and plays all 22 events', async ({ page }, testInfo) => {
  await page.goto('/');
  await expect(page.getByRole('heading', { name: /客户追加订单下的交付承诺重算/ })).toBeVisible();
  await expect(page.getByText('事件 0/22')).toBeVisible();

  if (testInfo.project.name === 'mobile-390') {
    await page.getByRole('tab', { name: 'A 线画布' }).click();
  }
  await expect(page.getByLabel('苏州制造基地电池冷却板 A 线工艺与物流画布')).toBeVisible();
  await expect(page.locator('.factory-node')).toHaveCount(14);

  for (let step = 1; step <= 22; step += 1) {
    await page.getByRole('button', { name: '下一个事件' }).click();
  }
  await expect(page.getByText('事件 22/22')).toBeVisible();
  await expect(page.getByRole('button', { name: '下一个事件' })).toBeDisabled();

  await page.screenshot({ path: `test-results/${testInfo.project.name}-completed.png` });
});

test('keeps primary controls within the viewport', async ({ page }) => {
  await page.goto('/');
  await expect(page.getByRole('button', { name: '播放故事' })).toBeVisible();
  const overflow = await page.evaluate(() => document.documentElement.scrollWidth > document.documentElement.clientWidth);
  expect(overflow).toBe(false);
});

test('captures governed run console screenshots for required viewports', async ({ page }, testInfo) => {
  const runBase = {
    run_id: 'run-hctm-evidence-20260721',
    run_version: 'v1',
    pack_key: 'hctm',
    pack_version: '0.1.0',
    scenario_key: 'order-expedite-01',
    plan_hash: 'plan-hctm-evidence-20260721',
    status: 'planned',
    current_act: 0,
    total_acts: 7,
    cursor: 0,
    tenant: 'tenant-hctm',
    target: 'http://127.0.0.1:3000',
    created_at: '2026-07-21T00:00:00Z',
    updated_at: '2026-07-21T00:00:00Z',
    allowed_actions: ['preflight'],
    outcome: {},
  };

  const plan = {
    pack_key: 'hctm',
    pack_version: '0.1.0',
    scenario_key: 'order-expedite-01',
    correlation_id: 'corr-so-202607-0001',
    total_events: 22,
    stages: [
      { stage: 'preflight', event_ids: ['evt-0001'], event_types: ['dryrun.preflight'], event_count: 1, action_hints: ['预检'] },
      { stage: 'initialize', event_ids: ['evt-init'], event_types: ['action.initialize'], event_count: 1, action_hints: ['初始化'] },
      { stage: 'act-1', event_ids: ['act-1-event'], event_types: ['act-1'], event_count: 1, action_hints: ['推进'] },
      { stage: 'act-2', event_ids: ['act-2-event'], event_types: ['act-2'], event_count: 1, action_hints: ['推进'] },
      { stage: 'act-3', event_ids: ['act-3-event'], event_types: ['act-3'], event_count: 1, action_hints: ['推进'] },
      { stage: 'act-4', event_ids: ['act-4-event'], event_types: ['act-4'], event_count: 1, action_hints: ['推进'] },
      { stage: 'act-5', event_ids: ['act-5-event'], event_types: ['act-5'], event_count: 1, action_hints: ['推进'] },
      { stage: 'act-6', event_ids: ['act-6-event'], event_types: ['act-6'], event_count: 1, action_hints: ['推进'] },
      { stage: 'act-7', event_ids: ['act-7-event'], event_types: ['act-7'], event_count: 14, action_hints: ['运行到结束'] },
      { stage: 'analyze', event_ids: ['evt-analyze'], event_types: ['agent.analysis'], event_count: 0, action_hints: ['分析'] },
      { stage: 'verify', event_ids: ['evt-verify'], event_types: ['verify'], event_count: 0, action_hints: ['验证'] },
      { stage: 'reset', event_ids: ['evt-reset-plan'], event_types: ['reset.plan'], event_count: 0, action_hints: ['复位'] },
    ],
    act_count: 7,
    allowable_run_actions: ['preflight', 'initialize', 'advance', 'run-to-end', 'analyze', 'verify', 'reset-plan', 'reset'],
    plan_hash: 'plan-hctm-evidence-20260721',
  };

  let state = { ...runBase };

  const transition = (action: string) => {
    if (action === 'preflight') {
      state = {
        ...state,
        status: 'ready',
        run_version: 'v2',
        allowed_actions: ['initialize', 'advance', 'run-to-end', 'analyze', 'verify', 'reset-plan', 'reset'],
      };
      return;
    }
    if (action === 'initialize') {
      state = {
        ...state,
        status: 'running',
        current_act: 1,
        cursor: 3,
        run_version: 'v3',
        allowed_actions: ['run-to-end', 'analyze', 'verify', 'reset-plan', 'reset'],
      };
    }
  };

  await page.route('**/api/v1/dev/token?**', (route) => route.fulfill({ json: { token: 'governed-run-token' } }));
  await page.route('**/api/v1/profile', (route) => route.fulfill({ json: { tenant_id: 'tenant-hctm', tenant_name: '华辰热管理系统集团' } }));
  await page.route('**/api/v1/scenarios/hctm/order-expedite-01/snapshot', (route) => route.fulfill({ json: { cursor: 22, completeness: 'partial', entities: [], events: [], recommendations: [], gaps: ['cost_actuals'], kpis: {} } }));
  await page.route('**/api/v1/scenarios/hctm/order-expedite-01/events?**', (route) => route.fulfill({ json: { items: [], next_cursor: 22, has_more: false } }));
  const totals: Record<string, number> = { sales_order: 1, work_order: 4, inventory: 6, equipment: 1 };
  await page.route('**/api/v1/entities/*/records?**', (route) => {
    const match = new URL(route.request().url()).pathname.match(/\/entities\/([^/]+)\/records/);
    return route.fulfill({ json: { total: totals[match?.[1] ?? ''] ?? 0, data: [] } });
  });

  await page.route('**/api/aese/v1/runs/plan', (route) => route.fulfill({ json: plan }));
  await page.route('**/api/aese/v1/runs', (route) => {
    if (route.request().method() === 'POST') {
      state = { ...runBase };
      return route.fulfill({ json: state });
    }
    return route.fallback();
  });
  await page.route(/.*\/api\/aese\/v1\/runs\/[^/]+(?:\?.*)?$/, (route) => {
    if (route.request().method() === 'GET') {
      return route.fulfill({ json: state });
    }
    return route.fallback();
  });
  await page.route(/.*\/api\/aese\/v1\/runs\/[^/]+\/preflight$/, (route) => {
    transition('preflight');
    return route.fulfill({ json: { run: state, action: 'preflight' } });
  });
  await page.route(/.*\/api\/aese\/v1\/runs\/[^/]+\/initialize$/, (route) => {
    transition('initialize');
    return route.fulfill({ json: { run: state, action: 'initialize' } });
  });

  await page.goto('/');
  await page.getByRole('button', { name: '打开 AESE 与 IAOS 联动中心' }).click();
  await page.getByRole('button', { name: /一键连接并检查/ }).click();
  await page.getByRole('tab', { name: '运行场景' }).click();
  await page.getByRole('button', { name: '创建并预检一个运行' }).click();
  await expect(page.getByText('预检 执行成功')).toBeVisible();
  const overflow = await page.evaluate(() => document.documentElement.scrollWidth > document.documentElement.clientWidth);
  expect(overflow).toBe(false);
  await page.screenshot({ path: `test-results/governed-console-${testInfo.project.name}.png` });

  await page.getByRole('button', { name: '初始化' }).click();
  await expect(page.getByText('初始化 执行成功')).toBeVisible();
  const overflowAfterInit = await page.evaluate(() => document.documentElement.scrollWidth > document.documentElement.clientWidth);
  expect(overflowAfterInit).toBe(false);
  await page.screenshot({ path: `test-results/governed-console-running-${testInfo.project.name}.png` });
});

test('configures and verifies the IAOS link without technical setup', async ({ page }) => {
  await page.route('**/api/v1/dev/token?**', (route) => route.fulfill({ json: { token: 'visual-hctm-token' } }));
  await page.route('**/api/v1/profile', (route) => route.fulfill({ json: { tenant_id: 'tenant-hctm', tenant_name: '华辰热管理系统集团' } }));
  await page.route('**/api/v1/scenarios/hctm/order-expedite-01/snapshot', (route) => route.fulfill({ json: { cursor: 12, completeness: 'partial', entities: [], events: [], recommendations: [], gaps: ['cost_actuals'], kpis: {} } }));
  await page.route('**/api/v1/scenarios/hctm/order-expedite-01/events?**', (route) => route.fulfill({ json: { items: [], next_cursor: 12, has_more: false } }));
  const totals: Record<string, number> = { sales_order: 1, work_order: 4, inventory: 6, equipment: 1 };
  await page.route('**/api/v1/entities/*/records?**', (route) => {
    const match = new URL(route.request().url()).pathname.match(/\/entities\/([^/]+)\/records/);
    return route.fulfill({ json: { total: totals[match?.[1] ?? ''] ?? 0, data: [] } });
  });

  await page.goto('/');
  await page.getByRole('button', { name: '打开 AESE 与 IAOS 联动中心' }).click();
  await expect(page.getByRole('dialog', { name: '企业沙盘联动中心' })).toBeVisible();
  await page.getByRole('button', { name: '一键连接并检查' }).click();
  await expect(page.getByText('全部可用')).toBeVisible();
  await expect(page.getByText('华辰热管理系统集团').last()).toBeVisible();
  await expect(page.getByText('IAOS · 4 条')).toBeVisible();
  await expect(page.getByText('Snapshot 与事件增量通道可用')).toBeVisible();
  await expect(page.getByRole('button', { name: '进入 AESE 在线 Live 沙盘' })).toBeVisible();
});

test('opens business object details and explains all three Agent roles', async ({ page }, testInfo) => {
  await page.goto('/');
  if (testInfo.project.name === 'mobile-390') {
    await page.getByRole('tab', { name: 'A 线画布' }).click();
  }
  await page.locator('.factory-node').filter({ hasText: '激光焊接' }).click();
  await expect(page.getByRole('dialog', { name: '2 号激光焊接机' })).toBeVisible();
  await page.getByRole('button', { name: '关闭对象详情' }).click();

  await page.getByRole('button', { name: /跳转到第 7 幕/ }).click();
  await page.getByRole('button', { name: '下一个事件' }).click();
  if (testInfo.project.name === 'mobile-390') {
    await page.getByRole('tab', { name: '事件 / Agent' }).click();
  }
  await page.getByRole('tab', { name: /Agent 建议/ }).click();
  await expect(page.getByText('计划 Agent').first()).toBeVisible();
  await expect(page.getByText('质量 Agent')).toBeVisible();
  await expect(page.getByText('经营分析 Agent')).toBeVisible();
  await expect(page.getByText('交付复盘')).toBeVisible();
});

test('renders governed Live facts, gaps, events and Agent evidence', async ({ page }, testInfo) => {
  const observedAt = '2026-07-20T09:00:00Z';
  const events = [
    { cursor: 1, event_id: 'evt-release', event_type: 'proc.production_order.released', occurred_at: observedAt, correlation_id: 'corr-so-202607-0001', business_object_type: 'production_order', business_object_code: 'WO-202607-0001', payload: {} },
    { cursor: 2, event_id: 'evt-ship', event_type: 'o2d.shipment.dispatched', occurred_at: observedAt, correlation_id: 'corr-so-202607-0001', business_object_type: 'shipment', business_object_code: 'SHIP-202607-0002', payload: {} },
  ];
  const snapshot = {
    snapshot_version: '1.0.0', pack_key: 'hctm', scenario_key: 'order-expedite-01', observed_at: observedAt,
    cursor: 2, completeness: 'partial', gaps: ['cost_actuals'],
    kpis: {
      order_demand: { value: 12000, unit: 'pcs' }, cumulative_available: { value: 11700, unit: 'pcs' },
      cumulative_shipped: { value: 11700, unit: 'pcs' }, ending_finished_goods: { value: 0, unit: 'pcs' }, delivery_gap: { value: 300, unit: 'pcs' },
    },
    entities: [
      { id: 'order-live', type: 'sales_order', business_code: 'SO-202607-0001', name: '客户订单', status: 'partially_shipped', attributes: { demand_qty: 12000, shipped_qty: 11700 } },
      { id: 'equipment-live', type: 'equipment', business_code: 'LAS-WLD-02', name: '关键焊接设备', status: 'maintenance', attributes: {} },
      { id: 'shipment-live', type: 'shipment', business_code: 'SHIP-202607-0002', name: '第二次发运', status: 'dispatched', attributes: { quantity: 2700 } },
    ],
    events,
    recommendations: [{ agent_key: 'business_analysis', summary: '在线交付复盘', recommendations: ['保留 300 件缺口并由人工确认后续措施'], object_refs: ['SO-202607-0001'], tool_call_ids: ['00000000-0000-0000-0000-000000000001'], completeness: 'partial', data_gaps: ['cost_actuals'], confidence: '0.92', status: 'suggested', requires_human_confirmation: true }],
  };
  await page.route('**/api/v1/scenarios/hctm/order-expedite-01/snapshot', (route) => route.fulfill({ json: snapshot }));
  await page.route('**/api/v1/scenarios/hctm/order-expedite-01/events?**', (route) => route.fulfill({ json: { items: [], next_cursor: 2, has_more: false } }));
  await page.route('**/api/v1/scenarios/hctm/order-expedite-01/events/stream?**', (route) => route.fulfill({ contentType: 'text/event-stream', body: ': heartbeat\n\n' }));
  await page.goto('/');
  await page.getByRole('button', { name: 'Live' }).click();
  await expect(page.getByRole('heading', { name: '华辰苏州基地 · 在线企业沙盘' })).toBeVisible();
  await expect(page.getByText('12,000 件')).toBeVisible();
  await expect(page.getByText('11,700 件').first()).toBeVisible();
  await expect(page.getByText('300 件')).toBeVisible();
  await expect(page.getByText(/数据缺口：cost_actuals/)).toBeVisible();
  if (testInfo.project.name === 'mobile-390') await page.getByRole('tab', { name: 'A 线画布' }).click();
  await expect(page.locator('.factory-node')).toHaveCount(14);
  await expect(page.locator('.factory-node').first()).toBeVisible();
  await page.screenshot({ path: `test-results/live-${testInfo.project.name}.png` });
});

test('runs governed scenario operations end-to-end and restores on page refresh', async ({ page }) => {
  const runStates: Record<string, Record<string, unknown>> = {
    base: {
      run_id: 'run-hctm-20260721',
      run_version: 'v1',
      pack_key: 'hctm',
      pack_version: '0.1.0',
      scenario_key: 'order-expedite-01',
      plan_hash: 'plan-hctm-m7-20260721',
      status: 'planned',
      current_act: 0,
      total_acts: 7,
      cursor: 0,
      tenant: 'tenant-hctm',
      target: 'http://127.0.0.1:3000',
      created_at: '2026-07-21T00:00:00Z',
      updated_at: '2026-07-21T00:00:00Z',
      allowed_actions: ['preflight'],
      reset_confirmation_required: false,
      outcome: {},
    },
    plan: {
      pack_key: 'hctm',
      pack_version: '0.1.0',
      scenario_key: 'order-expedite-01',
      correlation_id: 'corr-so-202607-0001',
      total_events: 22,
      stages: [
        { stage: 'preflight', event_ids: ['evt-0001'], event_types: ['dryrun.preflight'], event_count: 1, action_hints: ['预检'] },
        { stage: 'initialize', event_ids: ['evt-init'], event_types: ['action.initialize'], event_count: 1, action_hints: ['初始化'] },
        { stage: 'act-1', event_ids: ['act-1-event'], event_types: ['act-1'], event_count: 1, action_hints: ['推进'] },
        { stage: 'act-2', event_ids: ['act-2-event'], event_types: ['act-2'], event_count: 1, action_hints: ['推进'] },
        { stage: 'act-3', event_ids: ['act-3-event'], event_types: ['act-3'], event_count: 1, action_hints: ['推进'] },
        { stage: 'act-4', event_ids: ['act-4-event'], event_types: ['act-4'], event_count: 1, action_hints: ['推进'] },
        { stage: 'act-5', event_ids: ['act-5-event'], event_types: ['act-5'], event_count: 1, action_hints: ['推进'] },
        { stage: 'act-6', event_ids: ['act-6-event'], event_types: ['act-6'], event_count: 1, action_hints: ['推进'] },
        { stage: 'act-7', event_ids: ['act-7-event'], event_types: ['act-7'], event_count: 14, action_hints: ['运行到结束'] },
        { stage: 'analyze', event_ids: ['evt-analyze'], event_types: ['agent.analysis'], event_count: 0, action_hints: ['分析'] },
        { stage: 'verify', event_ids: ['evt-verify'], event_types: ['verify'], event_count: 0, action_hints: ['验证'] },
        { stage: 'reset', event_ids: ['evt-reset-plan'], event_types: ['reset.plan'], event_count: 0, action_hints: ['复位'] },
      ],
      act_count: 7,
      allowable_run_actions: ['preflight', 'initialize', 'advance', 'run-to-end', 'analyze', 'verify', 'reset-plan', 'reset'],
      plan_hash: 'plan-hctm-m7-20260721',
    },
  };
  let state = { ...runStates.base };

  const transition = (action: string) => {
    if (action === 'preflight') {
      state = { ...state, status: 'ready', run_version: 'v2', allowed_actions: ['initialize', 'advance', 'run-to-end', 'analyze', 'verify', 'reset-plan', 'reset'], cursor: 0 };
      return;
    }
    if (action === 'initialize') {
      state = { ...state, status: 'running', current_act: 1, run_version: 'v3', cursor: 3, allowed_actions: ['advance', 'run-to-end', 'analyze', 'verify', 'reset-plan', 'reset'], reset_confirmation_required: false };
      return;
    }
    if (action === 'advance') {
      state = { ...state, current_act: 2, run_version: 'v4', cursor: 8, allowed_actions: ['run-to-end', 'analyze', 'verify', 'reset-plan', 'reset'] };
      return;
    }
    if (action === 'run-to-end') {
      state = { ...state, status: 'awaiting_analysis', current_act: 7, run_version: 'v5', cursor: 20, allowed_actions: ['analyze', 'verify', 'reset-plan', 'reset'] };
      return;
    }
    if (action === 'analyze') {
      state = { ...state, status: 'awaiting_verification', current_act: 7, run_version: 'v6', cursor: 20, allowed_actions: ['verify', 'reset-plan', 'reset'] };
      return;
    }
    if (action === 'verify') {
      state = { ...state, status: 'completed', run_version: 'v7', cursor: 22, allowed_actions: ['reset-plan', 'reset'] };
      return;
    }
    if (action === 'reset-plan') {
      state = {
        ...state,
        run_version: 'v8',
        reset_confirmation_required: true,
        outcome: {
          reset_confirmation_token: 'reset-token-20260721',
          confirmation_expires_at: '2026-07-21T00:05:00Z',
        },
        allowed_actions: ['reset'],
      };
      return;
    }
    if (action === 'reset') {
      state = { ...state, status: 'reset', current_act: 0, run_version: 'v9', allowed_actions: ['preflight'], reset_confirmation_required: false };
    }
  };

  await page.route('**/api/v1/dev/token?**', (route) => route.fulfill({ json: { token: 'governed-run-token' } }));
  await page.route('**/api/v1/profile', (route) => route.fulfill({ json: { tenant_id: 'tenant-hctm', tenant_name: '华辰热管理系统集团' } }));
  await page.route('**/api/v1/scenarios/hctm/order-expedite-01/snapshot', (route) => route.fulfill({ json: { cursor: 22, completeness: 'partial', entities: [], events: [], recommendations: [], gaps: ['cost_actuals'], kpis: {} } }));
  await page.route('**/api/v1/scenarios/hctm/order-expedite-01/events?**', (route) => route.fulfill({ json: { items: [], next_cursor: 22, has_more: false } }));
  const totals: Record<string, number> = { sales_order: 1, work_order: 4, inventory: 6, equipment: 1 };
  await page.route('**/api/v1/entities/*/records?**', (route) => {
    const match = new URL(route.request().url()).pathname.match(/\/entities\/([^/]+)\/records/);
    return route.fulfill({ json: { total: totals[match?.[1] ?? ''] ?? 0, data: [] } });
  });

  await page.route('**/api/aese/v1/runs/plan', (route) => route.fulfill({ json: runStates.plan }));
  await page.route('**/api/aese/v1/runs', (route) => {
    if (route.request().method() === 'POST') {
      state = { ...runStates.base };
      return route.fulfill({ json: state });
    }
    return route.fallback();
  });
  await page.route(/.*\/api\/aese\/v1\/runs\/[^/]+(?:\?.*)?$/, (route) => {
    if (route.request().method() === 'GET') {
      return route.fulfill({ json: state });
    }
    return route.fallback();
  });
  await page.route(/.*\/api\/aese\/v1\/runs\/[^/]+\/(preflight|initialize|advance|run-to-end|analyze|verify|reset-plan|reset)$/, (route) => {
    const action = route.request().url().match(/\/api\/aese\/v1\/runs\/[^/]+\/([^/]+)$/)?.[1] ?? '';
    if (!action || route.request().method() !== 'POST') return route.fallback();
    transition(action);
    return route.fulfill({ json: { run: state, action } });
  });

  await page.goto('/');
  await page.getByRole('button', { name: '打开 AESE 与 IAOS 联动中心' }).click();
  await page.getByRole('button', { name: /一键连接并检查/ }).click();
  await page.getByRole('tab', { name: '运行场景' }).click();
  await page.getByRole('button', { name: '创建并预检一个运行' }).click();
  await expect(page.getByText('预检 执行成功')).toBeVisible();
  await page.getByRole('button', { name: '初始化' }).click();
  await expect(page.getByText('初始化 执行成功')).toBeVisible();
  await page.getByRole('button', { name: '运行到结束' }).click();
  await expect(page.getByText('运行到结束 执行成功')).toBeVisible();
  await page.getByRole('button', { name: 'Agent 分析' }).click();
  await expect(page.getByText('Agent 分析 执行成功')).toBeVisible();
  await page.getByRole('button', { name: '结果验证' }).click();
  await expect(page.getByText('结果验证 执行成功')).toBeVisible();
  await page.getByRole('button', { name: '复位预览' }).click();
  await expect(page.getByText('复位预览 执行成功')).toBeVisible();
  await expect(page.getByText('复位预览已就绪')).toBeVisible();
  await page.getByRole('button', { name: '执行复位' }).click();
  await expect(page.getByText('执行复位 执行成功')).toBeVisible();

  await page.reload();
  await page.getByRole('button', { name: '打开 AESE 与 IAOS 联动中心' }).click();
  await page.getByRole('button', { name: /一键连接并检查/ }).click();
  await page.getByRole('tab', { name: '运行场景' }).click();
  await expect(page.getByText('run-hctm-20260721', { exact: true })).toBeVisible();
  await expect(page.getByText('执行日志')).toBeVisible();
});

test('guards duplicate run action dispatch and surfaces permission errors', async ({ page }) => {
  const runBase = {
    run_id: 'run-hctm-permission',
    run_version: 'v1',
    pack_key: 'hctm',
    pack_version: '0.1.0',
    scenario_key: 'order-expedite-01',
    plan_hash: 'plan-hctm-permission',
    status: 'ready',
    current_act: 0,
    total_acts: 7,
    cursor: 0,
    tenant: 'tenant-hctm',
    target: 'http://127.0.0.1:3000',
    created_at: '2026-07-21T00:00:00Z',
    updated_at: '2026-07-21T00:00:00Z',
    allowed_actions: ['initialize', 'reset'],
    outcome: {},
  };
  let initializeCalls = 0;
  let resetCalls = 0;
  let state = runBase;

  await page.route('**/api/v1/dev/token?**', (route) => route.fulfill({ json: { token: 'governed-run-token' } }));
  await page.route('**/api/v1/profile', (route) => route.fulfill({ json: { tenant_id: 'tenant-hctm', tenant_name: '华辰热管理系统集团' } }));
  await page.route('**/api/v1/scenarios/hctm/order-expedite-01/snapshot', (route) => route.fulfill({ json: { cursor: 22, completeness: 'partial', entities: [], events: [], recommendations: [], gaps: ['cost_actuals'], kpis: {} } }));
  await page.route('**/api/v1/scenarios/hctm/order-expedite-01/events?**', (route) => route.fulfill({ json: { items: [], next_cursor: 22, has_more: false } }));
  const totals: Record<string, number> = { sales_order: 1, work_order: 4, inventory: 6, equipment: 1 };
  await page.route('**/api/v1/entities/*/records?**', (route) => {
    const match = new URL(route.request().url()).pathname.match(/\/entities\/([^/]+)\/records/);
    return route.fulfill({ json: { total: totals[match?.[1] ?? ''] ?? 0, data: [] } });
  });
  await page.route('**/api/aese/v1/runs/plan', (route) => route.fulfill({ json: { pack_key: 'hctm', pack_version: '0.1.0', scenario_key: 'order-expedite-01', correlation_id: 'corr-so-202607-0001', total_events: 22, stages: [{ stage: 'initialize', event_ids: [], event_types: [], event_count: 0, action_hints: [] }, { stage: 'analyze', event_ids: [], event_types: [], event_count: 0, action_hints: [] }, { stage: 'verify', event_ids: [], event_types: [], event_count: 0, action_hints: [] }, { stage: 'reset', event_ids: [], event_types: [], event_count: 0, action_hints: [] }], act_count: 7, allowable_run_actions: ['initialize', 'analyze', 'verify', 'reset'], plan_hash: 'plan-hctm-permission' } }));
  await page.route('**/api/aese/v1/runs', (route) => {
    if (route.request().method() === 'POST') return route.fulfill({ json: runBase });
    return route.fallback();
  });
  await page.route(/.*\/api\/aese\/v1\/runs\/[^/]+(?:\?.*)?$/, (route) => {
    if (route.request().method() === 'GET') return route.fulfill({ json: state });
    return route.fallback();
  });
  await page.route(/.*\/api\/aese\/v1\/runs\/[^/]+\/preflight$/, (route) => {
    state = { ...state, status: 'ready', allowed_actions: ['initialize', 'reset'], run_version: 'v2' };
    return route.fulfill({ json: { run: state, action: 'preflight' } });
  });
  await page.route(/.*\/api\/aese\/v1\/runs\/[^/]+\/initialize$/, (route) => {
    initializeCalls += 1;
    if (initializeCalls > 1) {
      return route.fallback();
    }
    state = { ...state, status: 'running', run_version: 'v2', allowed_actions: ['reset-plan', 'verify', 'analyze'], current_act: 1, cursor: 3 };
    return route.fulfill({ json: { run: state, action: 'initialize' } });
  });
  await page.route(/.*\/api\/aese\/v1\/runs\/[^/]+\/verify$/, (route) => {
    state = { ...state, allowed_actions: ['reset'], current_act: 7, status: 'completed', run_version: 'v3', outcome: { reset_confirmation_token: 'permission-reset' } };
    return route.fallback();
  });
  await page.route(/.*\/api\/aese\/v1\/runs\/[^/]+\/reset-plan$/, (route) => {
    state = { ...state, status: 'resetting', run_version: 'v3', allowed_actions: ['reset'], reset_confirmation_required: true, outcome: { reset_confirmation_token: 'permission-reset' } };
    return route.fulfill({ json: { run: state, action: 'reset-plan' } });
  });
  await page.route(/.*\/api\/aese\/v1\/runs\/[^/]+\/reset$/, (route) => {
    resetCalls += 1;
    if (resetCalls === 1) {
      return route.fulfill({
        status: 403,
        json: {
          error: 'missing permission',
          code: 'forbidden',
          required_permission: 'scenario.run.reset',
          run_id: state.run_id,
          run_version: state.run_version,
        },
      });
    }
    state = { ...state, status: 'reset', reset_confirmation_required: false, current_act: 0, run_version: 'v4', allowed_actions: ['initialize'] };
    return route.fulfill({ json: { run: state, action: 'reset' } });
  });

  await page.goto('/');
  await page.getByRole('button', { name: '打开 AESE 与 IAOS 联动中心' }).click();
  await page.getByRole('button', { name: /一键连接并检查/ }).click();
  await page.getByRole('tab', { name: '运行场景' }).click();
  await page.getByRole('button', { name: '创建并预检一个运行' }).click();
  await expect(page.getByText('run-hctm-permission', { exact: true })).toBeVisible();

  const initButton = page.getByRole('button', { name: '初始化' });
  await initButton.evaluate((button) => {
    button.click();
    button.click();
  });
  await expect.poll(() => initializeCalls).toBe(1);

  await page.getByRole('button', { name: '复位预览' }).click();
  await page.getByRole('button', { name: '执行复位' }).click();
  await expect(page.getByText(/缺少权限: scenario.run.reset/)).toBeVisible();
});
