import { ArrowRight, Building2, LoaderCircle, LogOut, Plus, RefreshCw } from "lucide-react";
import { useCallback, useEffect, useState } from "react";
import { listGenesisWorkspaces, resumeGenesisWorkspace, signOutGenesisPlayer } from "../../game/api";
import type { GenesisWorkspace, GenesisWorkspaceResult } from "../../game/types";
import { GenesisHome } from "./GenesisHome";
import "./GenesisCompanyLobby.css";

export function GenesisCompanyLobby({username,onCreateEnterprise,onOpenDemo,onOpenWorld,onSignedOut,onResume}:{username:string;onCreateEnterprise:()=>void;onOpenDemo:()=>void;onOpenWorld:()=>void;onSignedOut:()=>void;onResume:(workspace:GenesisWorkspaceResult)=>void}){
 const[items,setItems]=useState<GenesisWorkspace[]>([]),[loading,setLoading]=useState(true),[opening,setOpening]=useState(""),[error,setError]=useState("");
 const handleFailure=useCallback((reason:unknown)=>{
  if(!sessionStorage.getItem("aese_genesis_player_token")){
   signOutGenesisPlayer();onSignedOut();return;
  }
  setError(reason instanceof Error?reason.message:String(reason));
 },[onSignedOut]);
 const load=useCallback(async()=>{setLoading(true);setError("");try{setItems(await listGenesisWorkspaces())}catch(reason){handleFailure(reason)}finally{setLoading(false)}},[handleFailure]);
 useEffect(()=>{void load()},[load]);
 const resume=async(workspace:GenesisWorkspace)=>{setOpening(workspace.workspace_id);setError("");try{onResume(await resumeGenesisWorkspace(workspace as GenesisWorkspaceResult))}catch(reason){handleFailure(reason);setOpening("")}};
 const portfolio=<section className="gx-lobby" aria-labelledby="my-enterprises-title">
   <div className="gx-lobby__heading"><div><p>FOUNDER PORTFOLIO</p><h2 id="my-enterprises-title">我的企业</h2><span>选择一个已创建的独立 IAOS 租户继续经营。</span></div><div><button className="gx-lobby__refresh" onClick={()=>void load()} disabled={loading}><RefreshCw className={loading?"gx-spin":""} aria-hidden="true"/>刷新</button><button className="gx-lobby__create" onClick={onCreateEnterprise}><Plus aria-hidden="true"/>创建新企业</button></div></div>
   {error&&<p className="gx-lobby__error" role="alert">{error}</p>}
   {loading?<div className="gx-lobby__state"><LoaderCircle className="gx-spin" aria-hidden="true"/>正在加载你的企业…</div>:items.length===0?<div className="gx-lobby__empty"><Building2 aria-hidden="true"/><strong>还没有企业</strong><p>创建第一个独立企业空间，开始 M9 企业创生流程。</p><button onClick={onCreateEnterprise}><Plus aria-hidden="true"/>创建我的第一家企业</button></div>:<div className="gx-lobby__grid">{items.map(item=><article key={item.workspace_id}>
    <div className="gx-lobby__card-head"><span>{item.display_name.slice(0,1)}</span><em className={`gx-lobby__status gx-lobby__status--${item.status}`}>{item.status==="active"?"运行中":item.status==="failed"?"需恢复":"创建中"}</em></div>
    <h3>{item.display_name}</h3><p>{item.current_step==="identity_studio_ready"?"企业身份与设立流程":"企业创生空间"}</p>
    <dl><div><dt>企业案件</dt><dd>{item.case_code}</dd></div><div><dt>IAOS 租户</dt><dd>{item.tenant_id}</dd></div></dl>
    <button onClick={()=>void resume(item)} disabled={opening!==""||item.status!=="active"}>{opening===item.workspace_id?<><LoaderCircle className="gx-spin" aria-hidden="true"/>正在进入…</>:<>继续游戏<ArrowRight aria-hidden="true"/></>}</button>
   </article>)}</div>}
   <button className="gx-lobby__signout" onClick={()=>{signOutGenesisPlayer();onSignedOut()}}><LogOut aria-hidden="true"/>退出当前游戏用户</button>
  </section>;
 return <GenesisHome onCreateEnterprise={onCreateEnterprise} onOpenDemo={onOpenDemo} onOpenWorld={onOpenWorld} username={username} onSignOut={()=>{signOutGenesisPlayer();onSignedOut()}}>{portfolio}</GenesisHome>;
}
