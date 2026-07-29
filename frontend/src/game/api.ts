import type{FounderIntent,FounderIntentRequest,GameProjection,GenesisWorkspaceResult,NamingProposal}from"./types";
import{resolveIaosLifecycleBase,submitIncorporationObservation}from"../world/incorporation";
export class GameApiError extends Error{
 constructor(public status:number,message:string){super(message);this.name="GameApiError"}
}
export async function loadGameProjection(caseCode:string,frame:number,signal?:AbortSignal){
 const token=localStorage.getItem("iaos_token")??"";
 const tenant=localStorage.getItem("aese_iaos_tenant_id")??localStorage.getItem("iaos_tenant_id")??"";
 const response=await fetch(`/api/aese/v1/game/incorporation/${encodeURIComponent(caseCode)}/projection?frame=${frame}`,{signal,headers:{...(token?{Authorization:`Bearer ${token}`}:{}) ,...(tenant?{"X-IAOS-Tenant-Id":tenant}:{})}});
 if(!response.ok)throw new GameApiError(response.status,`GameProjection API ${response.status}: ${await response.text()}`);
 return response.json() as Promise<GameProjection>;
}
async function postJSON<T>(path:string,body:unknown){
 const token=localStorage.getItem("iaos_token")??"";
 const tenant=localStorage.getItem("aese_iaos_tenant_id")??localStorage.getItem("iaos_tenant_id")??"";
 const workspace=localStorage.getItem("aese_genesis_workspace_id")??"";
 const response=await fetch(path,{method:"POST",headers:{"content-type":"application/json",...(token?{Authorization:`Bearer ${token}`}:{}) ,...(tenant?{"X-IAOS-Tenant-Id":tenant}:{}),...(workspace?{"X-Genesis-Workspace-Id":workspace}:{})},body:JSON.stringify(body)});
 if(!response.ok)throw new Error(`Creative API ${response.status}`);
 return response.json() as Promise<T>;
}
export const analyzeFounderIntent=(request:FounderIntentRequest)=>postJSON<FounderIntent>("/api/aese/v1/game/creative/intent",request);
export type CreativeProviderStatus={state:"connected"|"degraded"|"fallback"|"not_configured";provider:string;model:string;base_url_host?:string;prompt_version:string};
export async function loadCreativeProviderStatus(){
 const response=await fetch("/api/aese/v1/game/creative/status");
 if(!response.ok)throw new Error(`Creative Provider 状态加载失败 ${response.status}`);
 return response.json() as Promise<CreativeProviderStatus>;
}
export async function generateCompanyNames(intent:FounderIntent){
 const result=await postJSON<{status:"candidate_only";proposals:NamingProposal[]}>("/api/aese/v1/game/creative/names",intent);
 return result.proposals;
}
const genesisPlayerID=()=>{
 let player=localStorage.getItem("aese_genesis_player_id");
 if(!player){player=`player-local-${crypto.randomUUID?.()??`${Date.now()}-${Math.random().toString(16).slice(2)}`}`;localStorage.setItem("aese_genesis_player_id",player)}
 return player;
};
const localPlayerMapKey="aese_genesis_local_players";
export const currentGenesisUsername=()=>localStorage.getItem("aese_genesis_username")??"";
export function signInGenesisPlayer(username:string){
 const normalized=username.trim();
 if(normalized.length<2)throw new Error("游戏用户名至少需要 2 个字符");
 let players:Record<string,string>={};
 try{players=JSON.parse(localStorage.getItem(localPlayerMapKey)??"{}") as Record<string,string>}catch{players={}}
 let player=players[normalized];
 if(!player){
  const unclaimed=localStorage.getItem("aese_genesis_player_id");
  const claimed=new Set(Object.values(players));
  player=unclaimed&&!claimed.has(unclaimed)?unclaimed:`player-local-${crypto.randomUUID?.()??`${Date.now()}-${Math.random().toString(16).slice(2)}`}`;
  players[normalized]=player;
  localStorage.setItem(localPlayerMapKey,JSON.stringify(players));
 }
 localStorage.setItem("aese_genesis_username",normalized);
 localStorage.setItem("aese_genesis_player_id",player);
 return{username:normalized,player_id:player};
}
export function signOutGenesisPlayer(){
 localStorage.removeItem("aese_genesis_username");
 localStorage.removeItem("aese_genesis_player_id");
 localStorage.removeItem("iaos_token");
 localStorage.removeItem("aese_iaos_tenant_id");
 localStorage.removeItem("iaos_tenant_id");
 localStorage.removeItem("aese_genesis_workspace_id");
 localStorage.removeItem("aese_genesis_case_code");
 localStorage.removeItem("aese_genesis_player_token");
}
export async function listGenesisWorkspaces(){
 const player=genesisPlayerID();
 const token=localStorage.getItem("aese_genesis_player_token")??localStorage.getItem("iaos_token")??"";
 const response=await fetch("/api/aese/v1/genesis/workspaces",{headers:{"X-Genesis-Player-Id":player,...(token?{Authorization:`Bearer ${token}`}:{})}});
 if(!response.ok)throw new Error(`企业列表加载失败 ${response.status}: ${await response.text()}`);
 const result=await response.json() as{items:GenesisWorkspaceResult[]};
 return result.items;
}
export async function createGenesisWorkspace(input:{display_name:string;idempotency_key:string;template_key:string;region:string;timezone:string;realism_level:"standard"|"strict";data_retention_confirmed:boolean}){
 const player=genesisPlayerID();
 const token=localStorage.getItem("aese_genesis_player_token")??localStorage.getItem("iaos_token")??"";
 const response=await fetch("/api/aese/v1/genesis/workspaces",{method:"POST",headers:{"content-type":"application/json","X-Genesis-Player-Id":player,...(token?{Authorization:`Bearer ${token}`}:{})},body:JSON.stringify({owner_player_id:player,...input})});
 if(!response.ok)throw new Error(`创业空间创建失败 ${response.status}: ${await response.text()}`);
 const result=await response.json() as GenesisWorkspaceResult;
 if(token)localStorage.setItem("aese_genesis_player_token",token);
 localStorage.setItem("iaos_token",result.tenant_token);
 localStorage.setItem("aese_iaos_tenant_id",result.tenant_id);
 localStorage.setItem("iaos_tenant_id",result.tenant_id);
 localStorage.setItem("aese_genesis_workspace_id",result.workspace_id);
 localStorage.setItem("aese_genesis_case_code",result.case_code);
 return result;
}
async function refreshGenesisSession(){
 const player=genesisPlayerID();
 const workspace=localStorage.getItem("aese_genesis_workspace_id");
 if(!workspace)throw new Error("缺少创业空间标识，请从游戏主页重新进入");
 const token=localStorage.getItem("aese_genesis_player_token")??localStorage.getItem("iaos_token")??"";
 const response=await fetch(`/api/aese/v1/genesis/workspaces/${encodeURIComponent(workspace)}/session`,{method:"POST",headers:{"X-Genesis-Player-Id":player,...(token?{Authorization:`Bearer ${token}`}:{})}});
 if(!response.ok)throw new Error(`Founder 会话刷新失败 ${response.status}: ${await response.text()}`);
 const result=await response.json() as GenesisWorkspaceResult;
 localStorage.setItem("iaos_token",result.tenant_token);
 return result.tenant_token;
}
export async function resumeGenesisWorkspace(workspace:GenesisWorkspaceResult){
 const player=genesisPlayerID();
 const token=localStorage.getItem("aese_genesis_player_token")??localStorage.getItem("iaos_token")??"";
 const response=await fetch(`/api/aese/v1/genesis/workspaces/${encodeURIComponent(workspace.workspace_id)}/session`,{method:"POST",headers:{"X-Genesis-Player-Id":player,...(token?{Authorization:`Bearer ${token}`}:{})}});
 if(!response.ok)throw new Error(`进入企业失败 ${response.status}: ${await response.text()}`);
 const result=await response.json() as GenesisWorkspaceResult;
 localStorage.setItem("iaos_token",result.tenant_token);
 localStorage.setItem("aese_iaos_tenant_id",result.tenant_id);
 localStorage.setItem("iaos_tenant_id",result.tenant_id);
 localStorage.setItem("aese_genesis_workspace_id",result.workspace_id);
 localStorage.setItem("aese_genesis_case_code",result.case_code);
 return result;
}
export async function createIncorporationCase(input:{case_code:string;case_name:string;proposed_company_name:string;registered_address:string;business_scope:string}){
 let token=localStorage.getItem("iaos_token");
 const tenant=localStorage.getItem("aese_iaos_tenant_id")??localStorage.getItem("iaos_tenant_id")??"tenant-hctm-genesis";
 if(!token)throw new Error("缺少 IAOS 登录凭据，请从 IAOS 重新进入企业创生");
 const submit=()=>fetch(`${resolveIaosLifecycleBase()}/api/v1/incorporations/cases`,{method:"POST",headers:{Authorization:`Bearer ${token}`,"X-Tenant-ID":tenant,"content-type":"application/json"},body:JSON.stringify(input)});
 let response=await submit();
 if(response.status===422){
  const detail=await response.clone().text();
  if(detail.includes("authenticated principal does not match")){
   token=await refreshGenesisSession();
   response=await submit();
  }
 }
 if(!response.ok)throw new Error(`IAOS incorporation.case.open ${response.status}: ${await response.text()}`);
 return response.json();
}
const iaosHeaders=()=>{
 const token=localStorage.getItem("iaos_token");
 const tenant=localStorage.getItem("aese_iaos_tenant_id")??localStorage.getItem("iaos_tenant_id")??"tenant-hctm-genesis";
 if(!token)throw new Error("缺少 IAOS 登录凭据，请从 IAOS 重新进入企业创生");
 return{Authorization:`Bearer ${token}`,"X-Tenant-ID":tenant,"content-type":"application/json"};
};
async function iaosPost(path:string,body:unknown){
 const response=await fetch(`${resolveIaosLifecycleBase()}${path}`,{method:"POST",headers:iaosHeaders(),body:JSON.stringify(body)});
 if(!response.ok)throw new Error(`IAOS ${response.status}: ${await response.text()}`);
 return response.json();
}
export type WorkItemInput={
 amount_minor?:number;currency?:string;business_note?:string;
 resolution_objective?:string;key_proposals?:string;risk_notes?:string;
};
export async function executeWorkItem(caseCode:string,sequence:number,input:WorkItemInput={}){
 return iaosPost(`/api/v1/incorporations/${encodeURIComponent(caseCode)}/work-items/${sequence}/execute`,input);
}
export async function dispatchAgentWorkItem(caseCode:string,sequence:number,input:WorkItemInput){
 return iaosPost(`/api/v1/incorporations/${encodeURIComponent(caseCode)}/work-items/${sequence}/dispatch-agent`,input);
}
export async function approveAndExecuteWorkItem(caseCode:string,sequence:number,capability:string,gate:string,input:WorkItemInput={}){
 const materialSuffix=input.amount_minor!==undefined?`-${input.amount_minor}`:"";
 const correlation=`corr-game-${caseCode}-${sequence}${materialSuffix}`;
 const submitted=await iaosPost(`/api/v1/incorporations/${encodeURIComponent(caseCode)}/gates/${encodeURIComponent(gate)}/submit`,{
  capability,correlation_id:correlation,intent:{case_code:caseCode,...input}
 }) as{approval?:{id?:string;status?:string}};
 const approvalId=submitted.approval?.id;
 if(!approvalId)throw new Error(`IAOS ${gate} 未返回审批请求`);
 if(submitted.approval?.status!=="approved"){
  await iaosPost(`/api/v1/approvals/${encodeURIComponent(approvalId)}/approve`,{
   note:`Enterprise Genesis 玩家批准 ${capability}`
  });
 }
 const executeInput:WorkItemInput&{correlation_id:string}={correlation_id:correlation};
 if(input.amount_minor!==undefined)executeInput.amount_minor=input.amount_minor;
 if(input.currency!==undefined)executeInput.currency=input.currency;
 return iaosPost(`/api/v1/incorporations/${encodeURIComponent(caseCode)}/work-items/${sequence}/execute`,executeInput);
}
export async function resolveWorldWorkItem(caseCode:string,sequence:number,capability:string,correlation:string,result:string){
 const contracts:Record<string,string>={
  "registration.observation.commit":"registration.decision.observed.v1",
  "bank.account.observation.commit":"bank.account.decision.observed.v1",
  "executive.appointment.acceptance.commit":"appointment.acceptance.observed.v1"
 };
 const payloadType=contracts[capability];
 if(!payloadType)throw new Error(`没有 World Observation 合同: ${capability}`);
 await submitIncorporationObservation(caseCode,payloadType,result,correlation);
 await iaosPost(`/api/v1/incorporations/${encodeURIComponent(caseCode)}/work-items/${sequence}/execute`,{correlation_id:correlation});
}
