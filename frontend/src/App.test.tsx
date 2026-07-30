import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
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
    sessionStorage.clear();
  });
  afterEach(() => vi.restoreAllMocks());

  it('opens the game user login at the root URL', () => {
    render(<App />);

    expect(
      screen.getByRole('heading', { name: '回到你的企业世界' }),
    ).toBeInTheDocument();
    expect(screen.getByLabelText('用户名')).toBeInTheDocument();
    expect(screen.getByLabelText('密码')).toBeInTheDocument();
  });

  it('logs in before starting tenant provisioning', async () => {
    const user = userEvent.setup();
    vi.spyOn(globalThis, "fetch")
      .mockResolvedValueOnce(new Response(JSON.stringify({
        status:"success",
        token:"player-token",
        expires_at:"2026-07-31T00:00:00Z",
        player:{subject_id:"founder-principal",username:"founder-principal",display_name:"Founder"},
      }),{status:200}))
      .mockResolvedValueOnce(new Response(JSON.stringify({items:[]}),{status:200}));
    render(<App />);

    await user.type(screen.getByLabelText('用户名'), 'founder-principal');
    await user.type(screen.getByLabelText('密码'), 'FounderPass123');
    await user.click(screen.getByRole('button', { name: /安全登录/ }));
    await user.click(screen.getAllByRole('button', { name: /创建新企业/ })[0]);

    expect(
      screen.getByRole('heading', { name: '先定义创业项目' }),
    ).toBeInTheDocument();
    expect(
      screen.getByRole('button', { name: '下一步' }),
    ).toBeInTheDocument();
  });

  it('registers a Player identity before showing the enterprise lobby', async () => {
    const user = userEvent.setup();
    const fetchMock=vi.spyOn(globalThis,"fetch")
      .mockResolvedValueOnce(new Response(JSON.stringify({
        status:"success",token:"registered-player-token",expires_at:"2026-07-31T00:00:00Z",
        player:{subject_id:"player-new-founder",username:"new-founder",display_name:"新创始人"},
      }),{status:201}))
      .mockResolvedValueOnce(new Response(JSON.stringify({items:[]}),{status:200}));
    render(<App/>);

    await user.click(screen.getByRole("tab",{name:"注册"}));
    await user.type(screen.getByLabelText("用户名"),"new-founder");
    await user.type(screen.getByLabelText("显示名称 选填"),"新创始人");
    await user.type(screen.getByLabelText("密码",{exact:true}),"NewFounder2026");
    await user.type(screen.getByLabelText("确认密码"),"NewFounder2026");
    await user.click(screen.getByRole("button",{name:/注册并进入/}));

    expect(await screen.findByRole("heading",{name:"我的企业"})).toBeInTheDocument();
    expect(JSON.parse(String(fetchMock.mock.calls[0][1]?.body))).toEqual({
      username:"new-founder",display_name:"新创始人",password:"NewFounder2026",
    });
    expect(sessionStorage.getItem("aese_genesis_player_token")).toBe("registered-player-token");
  });
});
