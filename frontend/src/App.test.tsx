import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it } from 'vitest';
import App from './App';

describe('AESE sandbox', () => {
  beforeEach(() => {
    window.location.hash = 'sandbox';
  });

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

  it('places the system atlas command in the control bar', async () => {
    render(<App />);
    await screen.findByText('事件 0/22');
    const atlasButton = screen.getByRole('button', { name: '打开系统全景' });
    expect(atlasButton.closest('.playback-controls')).not.toBeNull();
    expect(document.querySelector('.aese-atlas-launch')).toBeNull();
  });
});

describe('AESE game home', () => {
  beforeEach(() => {
    window.history.replaceState(null, '', '/');
    localStorage.clear();
  });

  it('opens the game user login at the root URL', () => {
    render(<App />);

    expect(
      screen.getByRole('heading', { name: '回到你的企业世界' }),
    ).toBeInTheDocument();
    expect(screen.getByLabelText('游戏用户名')).toBeInTheDocument();
  });

  it('logs in before starting tenant provisioning', async () => {
    const user = userEvent.setup();
    render(<App />);

    await user.type(screen.getByLabelText('游戏用户名'), 'founder-principal');
    await user.click(screen.getByRole('button', { name: /进入企业世界/ }));
    await user.click(screen.getAllByRole('button', { name: /创建新企业/ })[0]);

    expect(
      screen.getByRole('heading', { name: '先定义创业项目' }),
    ).toBeInTheDocument();
    expect(
      screen.getByRole('button', { name: '下一步' }),
    ).toBeInTheDocument();
  });
});
