import { Building2, Eye, EyeOff, LoaderCircle, LogIn, ShieldCheck, UserPlus } from "lucide-react";
import { FormEvent, useState } from "react";
import { loginGenesisPlayer, registerGenesisPlayer } from "../../game/api";
import "./GenesisLogin.css";

type Mode="login"|"register";

export function GenesisLogin({onSignedIn}:{onSignedIn:(username:string)=>void}){
  const[mode,setMode]=useState<Mode>("login");
  const[username,setUsername]=useState("");
  const[displayName,setDisplayName]=useState("");
  const[password,setPassword]=useState("");
  const[confirmation,setConfirmation]=useState("");
  const[showPassword,setShowPassword]=useState(false);
  const[busy,setBusy]=useState(false);
  const[error,setError]=useState("");
  const isRegistration=mode==="register";
  const valid=username.trim().length>=3&&password.length>=10&&(!isRegistration||confirmation===password);

  const switchMode=(next:Mode)=>{
    setMode(next);setError("");setPassword("");setConfirmation("");
  };
  const submit=async(event:FormEvent)=>{
    event.preventDefault();
    if(isRegistration&&password!==confirmation){
      setError("两次输入的密码不一致");
      return;
    }
    setBusy(true);setError("");
    try{
      const session=isRegistration
        ?await registerGenesisPlayer({username:username.trim(),display_name:displayName.trim(),password})
        :await loginGenesisPlayer({username:username.trim(),password});
      onSignedIn(session.player.username);
    }catch(reason){
      setError(reason instanceof Error?reason.message:String(reason));
    }finally{
      setBusy(false);
    }
  };

  return <main className="gx-login">
    <section className="gx-login__visual" aria-label="Enterprise Genesis 企业世界">
      <div className="gx-login__brand"><span>G</span><div><strong>Enterprise Genesis</strong><small>IAOS × AESE 企业创生世界</small></div></div>
      <div className="gx-login__message">
        <p>FOUNDER ACCESS</p>
        <h1>{isRegistration?"建立创始人身份":"回到你的企业世界"}</h1>
        <span>创始人账号由 IAOS 统一认证。每个企业拥有独立租户、治理流程和可审计证据，账号本身不会自动创建企业。</span>
      </div>
      <div className="gx-login__city" aria-hidden="true"><i/><i/><i/><i/><i/></div>
    </section>
    <section className="gx-login__panel">
      <form onSubmit={submit}>
        <div className="gx-login__icon"><Building2 aria-hidden="true"/></div>
        <p className="gx-login__eyebrow">IAOS PLAYER ACCOUNT</p>
        <h2>{isRegistration?"注册创始人账号":"创始人登录"}</h2>
        <p className="gx-login__intro">
          {isRegistration?"先创建个人身份，登录后再建立或加入企业空间。":"可使用已注册的 Player 账号；现有 IAOS 用户首次登录会安全关联原有身份与企业。"}
        </p>
        <div className="gx-login__tabs" role="tablist" aria-label="认证方式">
          <button type="button" role="tab" aria-selected={!isRegistration} onClick={()=>switchMode("login")}>登录</button>
          <button type="button" role="tab" aria-selected={isRegistration} onClick={()=>switchMode("register")}>注册</button>
        </div>
        <label htmlFor="genesis-username">用户名</label>
        <input id="genesis-username" autoFocus autoComplete="username" minLength={3} maxLength={40} pattern={"[A-Za-z0-9._\\-]{3,40}"} required value={username} onChange={event=>setUsername(event.target.value)} placeholder="例如：founder-principal"/>
        <small>3–40 位，仅支持字母、数字、点、下划线和连字符。</small>
        {isRegistration&&<>
          <label htmlFor="genesis-display-name">显示名称 <span>选填</span></label>
          <input id="genesis-display-name" autoComplete="name" maxLength={80} value={displayName} onChange={event=>setDisplayName(event.target.value)} placeholder="例如：企业创始人"/>
        </>}
        <label htmlFor="genesis-password">密码</label>
        <div className="gx-login__password">
          <input id="genesis-password" type={showPassword?"text":"password"} autoComplete={isRegistration?"new-password":"current-password"} minLength={10} maxLength={128} required value={password} onChange={event=>setPassword(event.target.value)} placeholder="输入密码"/>
          <button type="button" aria-label={showPassword?"隐藏密码":"显示密码"} onClick={()=>setShowPassword(value=>!value)}>{showPassword?<EyeOff aria-hidden="true"/>:<Eye aria-hidden="true"/>}</button>
        </div>
        {isRegistration&&<>
          <small>至少 10 位，并同时包含字母和数字。</small>
          <label htmlFor="genesis-confirmation">确认密码</label>
          <input id="genesis-confirmation" type={showPassword?"text":"password"} autoComplete="new-password" minLength={10} maxLength={128} required value={confirmation} onChange={event=>setConfirmation(event.target.value)} placeholder="再次输入密码"/>
        </>}
        {error&&<p className="gx-login__error" role="alert" aria-live="polite">{error}</p>}
        <button className="gx-login__submit" disabled={busy||!valid}>
          {busy?<LoaderCircle className="gx-spin" aria-hidden="true"/>:isRegistration?<UserPlus aria-hidden="true"/>:<LogIn aria-hidden="true"/>}
          {busy?"正在验证…":isRegistration?"注册并进入":"安全登录"}
        </button>
        <p className="gx-login__security"><ShieldCheck aria-hidden="true"/>密码只提交给 IAOS 认证服务，AESE 不保存密码；连续失败会触发临时锁定。</p>
      </form>
    </section>
  </main>
}
