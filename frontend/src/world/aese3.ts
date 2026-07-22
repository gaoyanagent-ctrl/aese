export type AESE3Evidence={code:string;value:number;unit:string;source_ref:string};
export type AESE3Milestone={code:string;title:string;design:string;terminal:string;terminal_ready:boolean;world_owner:string;business_owner:string;evidence:AESE3Evidence[];automatic_business_writes:number;evidence_hash:string};
export type AESE3Program={schema_version:string;code:string;tenant:string;timezone:string;parent_terminal:string;milestones:AESE3Milestone[];industry_simulation_platform_ready:boolean;program_hash:string;limitations:string};
const base=(import.meta.env.VITE_AESE_API_BASE_URL as string|undefined)?.replace(/\/$/,"")??"";
export async function loadAESE3(signal?:AbortSignal):Promise<AESE3Program>{const r=await fetch(`${base}/api/aese/v1/world/aese3`,{signal});if(!r.ok)throw new Error(`AESE 3 API ${r.status}`);return r.json() as Promise<AESE3Program>}
