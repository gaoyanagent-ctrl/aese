import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it } from 'vitest';
import App from './App';

describe('AESE sandbox', () => {
  it('loads the HCTM scenario and advances deterministically', async () => {
    const user = userEvent.setup();
    render(<App />);
    expect(await screen.findByRole('heading', { name: /客户追加订单下的交付承诺重算/ })).toBeInTheDocument();
    expect(screen.getByText('事件 0/22')).toBeInTheDocument();
    await user.click(screen.getByRole('button', { name: '下一个事件' }));
    expect(screen.getByText('事件 1/22')).toBeInTheDocument();
    await user.click(screen.getByRole('button', { name: '重置故事' }));
    await waitFor(() => expect(screen.getByText('事件 0/22')).toBeInTheDocument());
  });

  it('filters the event feed by domain', async () => {
    const user = userEvent.setup();
    render(<App />);
    await screen.findByText('事件 0/22');
    await user.selectOptions(screen.getByLabelText('事件领域'), 'equipment');
    expect(screen.getByText('焊接设备停机')).toBeInTheDocument();
    expect(screen.queryByText('收到原订单')).not.toBeInTheDocument();
  });
});
