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
  await expect(page.getByRole('button', { name: /进入 Live 沙盘/ })).toBeVisible();
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
