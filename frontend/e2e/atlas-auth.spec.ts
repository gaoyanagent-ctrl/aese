import { expect, test } from '@playwright/test';

test('Atlas refreshes a stale IAOS token and loads the AESE graph', async ({ page }) => {
  await page.addInitScript(() => localStorage.setItem('iaos_token', 'expired-demo-token'));
  const atlasResponses: number[] = [];
  page.on('response', response => {
    if (response.url().includes('/api/v1/system-atlas?view=aese')) atlasResponses.push(response.status());
  });

  await page.goto('/#atlas');

  await expect(page.getByText('AESE World', { exact: true }).first()).toBeVisible();
  await expect.poll(() => atlasResponses).toEqual([401, 200]);
  expect(await page.evaluate(() => localStorage.getItem('iaos_token'))).not.toBe('expired-demo-token');
});
