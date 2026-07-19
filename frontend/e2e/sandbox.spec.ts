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
