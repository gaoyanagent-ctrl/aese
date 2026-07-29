import { Building2, LoaderCircle, LogIn, ShieldCheck } from "lucide-react";
import { FormEvent, useState } from "react";
import { signInGenesisPlayer } from "../../game/api";
import "./GenesisLogin.css";

export function GenesisLogin({onSignedIn}:{onSignedIn:(username:string)=>void}){
  const[username,setUsername]=useState("");
  const[busy,setBusy]=useState(false);
  const[error,setError]=useState("");
  const submit=(event:FormEvent)=>{
    event.preventDefault();
    setBusy(true);setError("");
    try{onSignedIn(signInGenesisPlayer(username).username)}
    catch(reason){setError(reason instanceof Error?reason.message:String(reason));setBusy(false)}
  };
  return <main className="gx-login">
    <section className="gx-login__visual" aria-label="Enterprise Genesis 企业世界">
      <div className="gx-login__brand"><span>G</span><div><strong>Enterprise Genesis</strong><small>IAOS × AESE 企业创生世界</small></div></div>
      <div className="gx-login__message">
        <p>FOUNDER ACCESS</p>
        <h1>回到你的企业世界</h1>
        <span>创建、治理并持续经营属于你的数字企业。每个企业拥有独立 IAOS 租户、流程和证据。</span>
      </div>
      <div className="gx-login__city" aria-hidden="true"><i/><i/><i/><i/><i/></div>
    </section>
    <section className="gx-login__panel">
      <form onSubmit={submit}>
        <div className="gx-login__icon"><Building2 aria-hidden="true"/></div>
        <p className="gx-login__eyebrow">游戏用户登录</p>
        <h2>选择你的创始人身份</h2>
        <p className="gx-login__intro">登录后可以选择已经创建的企业，或建立一个新的独立企业空间。</p>
        <label htmlFor="genesis-username">游戏用户名</label>
        <input id="genesis-username" autoFocus autoComplete="username" minLength={2} maxLength={40} value={username} onChange={event=>setUsername(event.target.value)} placeholder="例如：founder-principal"/>
        <small>首次登录会关联此浏览器中已经创建的企业。</small>
        {error&&<p className="gx-login__error" role="alert">{error}</p>}
        <button disabled={busy||username.trim().length<2}>{busy?<LoaderCircle className="gx-spin" aria-hidden="true"/>:<LogIn aria-hidden="true"/>}进入企业世界</button>
        <p className="gx-login__security"><ShieldCheck aria-hidden="true"/>当前为私人开发环境的本机游戏身份；正式多人版将接入 IAOS 账号认证。</p>
      </form>
    </section>
  </main>
}
