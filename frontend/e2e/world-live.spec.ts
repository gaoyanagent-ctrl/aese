import { expect, test } from '@playwright/test';

test('deployed World API renders the first frame without a page error', async ({ page }) => {
  const pageErrors: string[] = [];
  page.on('pageerror', error => pageErrors.push(error.message));

  await page.goto('/#world');

  await expect(page.getByRole('heading', { name: '世界已退化，系统与角色尚未知' })).toBeVisible();
  await expect(page.getByText('角色尚未知，不可读取完整 World State。')).toBeVisible();
  expect(pageErrors).toEqual([]);
});
