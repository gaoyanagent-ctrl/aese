import { ArrowLeft, Building2, LoaderCircle, ShieldCheck, Sparkles } from "lucide-react";
import { useRef, useState } from "react";
import { createGenesisWorkspace } from "../../game/api";
import type { GenesisWorkspaceResult } from "../../game/types";
import "./GenesisOnboarding.css";

export function GenesisOnboarding({
  onBack,
  onReady,
}: {
  onBack: () => void;
  onReady: (workspace: GenesisWorkspaceResult) => void;
}) {
  const [name, setName] = useState("我的新企业");
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState("");
  const key = useRef(`genesis-create-${Date.now()}-${Math.random().toString(16).slice(2)}`);
  const create = async () => {
    setBusy(true);
    setError("");
    try {
      const workspace = await createGenesisWorkspace(name, key.current);
      onReady(workspace);
    } catch (reason) {
      setError(reason instanceof Error ? reason.message : String(reason));
    } finally {
      setBusy(false);
    }
  };
  return (
    <main className="gx-onboarding">
      <header>
        <button onClick={onBack}><ArrowLeft aria-hidden="true" />返回主页</button>
        <span>ENTERPRISE GENESIS · ZERO START</span>
      </header>
      <section className="gx-onboarding__card">
        <div className="gx-onboarding__icon"><Building2 aria-hidden="true" /></div>
        <p className="gx-onboarding__eyebrow">第一步 · 创建隔离创业空间</p>
        <h1>先为你的企业建立独立运行空间</h1>
        <p className="gx-onboarding__lead">
          平台会自动分配 IAOS 租户、创始人身份、M9 Runtime 和 AESE World Run。
          你不需要填写或选择 tenant ID。
        </p>
        <label>
          创业项目名称
          <input value={name} minLength={2} maxLength={80} onChange={(event) => setName(event.target.value)} disabled={busy} />
          <small>这是创建阶段的项目名称，正式公司名称稍后由 MiniMax M3 生成。</small>
        </label>
        <ol>
          <li className={busy ? "active" : ""}><span>1</span><div><strong>分配独立 IAOS Tenant</strong><small>服务端生成不可预测的租户标识</small></div></li>
          <li><span>2</span><div><strong>建立 Founder 运行身份</strong><small>初始化租户内创始治理账户</small></div></li>
          <li><span>3</span><div><strong>安装 M9 Runtime</strong><small>Process、Capability、Agent 与权限</small></div></li>
          <li><span>4</span><div><strong>创建 World Run</strong><small>绑定 AESE 企业世界与案件编号</small></div></li>
        </ol>
        {error && <p className="gx-onboarding__error" role="alert">{error}</p>}
        <button className="gx-onboarding__submit" disabled={busy || name.trim().length < 2} onClick={() => void create()}>
          {busy ? <LoaderCircle className="gx-spin" aria-hidden="true" /> : <Sparkles aria-hidden="true" />}
          {busy ? "正在创建独立企业空间…" : "创建空间并进入 AI 身份工作室"}
        </button>
        <p className="gx-onboarding__security"><ShieldCheck aria-hidden="true" />租户由平台分配；企业之间通过 IAOS RLS 隔离。</p>
      </section>
    </main>
  );
}
