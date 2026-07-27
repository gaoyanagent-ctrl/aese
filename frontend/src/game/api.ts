import type{FounderIntent,FounderIntentRequest,GameProjection,NamingProposal}from"./types";
import{resolveIaosLifecycleBase,submitIncorporationObservation}from"../world/incorporation";
export async function loadGameProjection(caseCode:string,frame:number,signal?:AbortSignal){
 const token=localStorage.getItem("iaos_token")??"";
 const tenant=localStorage.getItem("aese_iaos_tenant_id")??localStorage.getItem("iaos_tenant_id")??"";
 const response=await fetch(`/api/aese/v1/game/incorporation/${encodeURIComponent(caseCode)}/projection?frame=${frame}`,{signal,headers:{...(token?{Authorization:`Bearer ${token}`}:{}) ,...(tenant?{"X-IAOS-Tenant-Id":tenant}:{})}});
 if(!response.ok)throw new Error(`GameProjection API ${response.status}`);
 return response.json() as Promise<GameProjection>;
}
async function postJSON<T>(path:string,body:unknown){
 const response=await fetch(path,{method:"POST",headers:{"content-type":"application/json"},body:JSON.stringify(body)});
 if(!response.ok)throw new Error(`Creative API ${response.status}`);
 return response.json() as Promise<T>;
}
export const analyzeFounderIntent=(request:FounderIntentRequest)=>postJSON<FounderIntent>("/api/aese/v1/game/creative/intent",request);
export async function generateCompanyNames(intent:FounderIntent){
 const result=await postJSON<{status:"candidate_only";proposals:NamingProposal[]}>("/api/aese/v1/game/creative/names",intent);
 return result.proposals;
}
export async function createIncorporationCase(input:{case_code:string;case_name:string;proposed_company_name:string;registered_address:string;business_scope:string}){
 const token=localStorage.getItem("iaos_token");
 const tenant=localStorage.getItem("aese_iaos_tenant_id")??localStorage.getItem("iaos_tenant_id")??"tenant-hctm-genesis";
 if(!token)throw new Error("缺少 IAOS 登录凭据，请从 IAOS 重新进入企业创生");
 const response=await fetch(`${resolveIaosLifecycleBase()}/api/v1/incorporations/cases`,{method:"POST",headers:{Authorization:`Bearer ${token}`,"X-Tenant-ID":tenant,"content-type":"application/json"},body:JSON.stringify(input)});
 if(!response.ok)throw new Error(`IAOS incorporation.case.open ${response.status}: ${await response.text()}`);
 return response.json();
}
export function iaosWorkItemURL(caseCode:string,sequence:number,capability:string){
 const configured=localStorage.getItem("aese_iaos_ui_url")?.trim();
 const base=configured||`http://${window.location.hostname||"127.0.0.1"}:3000/`;
 const query=new URLSearchParams({tenant:localStorage.getItem("aese_iaos_tenant_id")??"tenant-hctm-genesis",case:caseCode,step:String(sequence),capability});
 return `${base.replace(/\/$/,"")}/#enterprise_lifecycle?${query}`;
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
 const token=localStorage.getItem("iaos_token");
 const tenant=localStorage.getItem("aese_iaos_tenant_id")??"tenant-hctm-genesis";
 if(!token)throw new Error("缺少 IAOS 登录凭据");
 const response=await fetch(`${resolveIaosLifecycleBase()}/api/v1/incorporations/${encodeURIComponent(caseCode)}/work-items/${sequence}/execute`,{method:"POST",headers:{Authorization:`Bearer ${token}`,"X-Tenant-ID":tenant,"content-type":"application/json"},body:JSON.stringify({correlation_id:correlation})});
 if(!response.ok)throw new Error(`World wait commit ${response.status}: ${await response.text()}`);
}
