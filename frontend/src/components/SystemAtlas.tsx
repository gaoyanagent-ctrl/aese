import { useCallback, useEffect, useMemo, useState } from 'react';
import { Background, Controls, Handle, MiniMap, Position, ReactFlow, type Edge, type Node, type NodeProps } from '@xyflow/react';
import { Activity, ArrowLeft, CheckCircle2, Clock3, ExternalLink, Filter, RefreshCw, Search, X } from 'lucide-react';
import './SystemAtlas.css';

type Status='planned'|'designed'|'building'|'validating'|'completed'|'blocked'|'deferred';
type AtlasNode={id:string;key:string;project_key:string;parent_key?:string;title:string;subtitle:string;kind:string;layer:string;summary:string;target_state:string;current_state:string;status:Status;progress:number;health:string;owner:string;sort_order:number;design_refs:string[];code_refs:string[];evidence_refs:string[];updated_at:string};
type AtlasUpdate={id:string;node_key:string;summary:string;detail:string;occurred_at:string;created_by:string};
type AtlasData={summary:{node_count:number;overall_progress:number;by_status:Record<string,number>;last_updated_at?:string};nodes:AtlasNode[];edges:{id:string;source:string;target:string;relation:string}[];recent_updates:AtlasUpdate[]};
const labels:Record<Status,string>={planned:'待规划',designed:'已设计',building:'建设中',validating:'验证中',completed:'已完成',blocked:'受阻',deferred:'后置'};
const colors:Record<Status,string>={planned:'#64748b',designed:'#8b5cf6',building:'#168aad',validating:'#f59e0b',completed:'#059669',blocked:'#dc2626',deferred:'#78716c'};

function Card({data,selected}:NodeProps<Node<AtlasNode>>){return <div className={`aese-atlas-node ${selected?'selected':''}`} style={{borderLeftColor:colors[data.status]}}><Handle type="target" position={Position.Left}/><div><span>{data.layer}</span><i style={{background:colors[data.status]}}/></div><strong>{data.title}</strong><small>{data.subtitle}</small><div className="aese-atlas-meter"><i style={{width:`${data.progress}%`,background:colors[data.status]}}/></div><footer><b>{data.progress}%</b><span>{labels[data.status]}</span></footer><Handle type="source" position={Position.Right}/></div>}
const nodeTypes={atlas:Card};

async function atlasToken(){
  const existing=localStorage.getItem('iaos_token');if(existing)return existing;
  const tenant=localStorage.getItem('aese_iaos_tenant_id')??'tenant-hctm';
  const response=await fetch(`/api/v1/dev/token?tenant_id=${encodeURIComponent(tenant)}&roles=admin`);if(!response.ok)throw new Error('请先通过联动中心建立 IAOS 身份');
  const body=await response.json() as {token?:string};if(!body.token)throw new Error('IAOS 未返回可用身份');localStorage.setItem('iaos_token',body.token);return body.token;
}

function layout(nodes:AtlasNode[]):Node<AtlasNode>[] {
  const roots:Record<string,{x:number;y:number}>={ecosystem:{x:40,y:40},aese:{x:330,y:40},'shared.contract':{x:620,y:40}};const counters:Record<string,number>={model:0,integration:0,intelligence:0,experience:0,operations:0,simulation:0,scale:0,governance:0};
  const lane:Record<string,number>={model:60,integration:60,intelligence:350,experience:350,operations:640,simulation:640,scale:930,governance:930};
  return nodes.map(n=>{let p=roots[n.key];if(!p){const i=counters[n.layer]??0;counters[n.layer]=i+1;p={x:lane[n.layer]??930,y:230+i*145}}return{id:n.key,type:'atlas',position:p,data:n}});
}

export function SystemAtlas({onExit}:{onExit:()=>void}){
  const [data,setData]=useState<AtlasData|null>(null),[selected,setSelected]=useState<AtlasNode|null>(null),[error,setError]=useState(''),[loading,setLoading]=useState(true),[query,setQuery]=useState(''),[status,setStatus]=useState<Status|'all'>('all');
  const load=useCallback(async()=>{setLoading(true);setError('');try{const token=await atlasToken();const response=await fetch('/api/v1/system-atlas?view=aese',{headers:{Authorization:`Bearer ${token}`}});if(!response.ok){const body=await response.json().catch(()=>({}));throw new Error(body.error??`HTTP ${response.status}`)}const next=await response.json() as AtlasData;next.nodes=next.nodes.map(node=>({...node,design_refs:node.design_refs??[],code_refs:node.code_refs??[],evidence_refs:node.evidence_refs??[]}));setData(next)}catch(reason){setError(reason instanceof Error?reason.message:'全景数据载入失败')}finally{setLoading(false)}},[]);
  useEffect(()=>{void load()},[load]);
  const visible=useMemo(()=>new Set((data?.nodes??[]).filter(n=>(status==='all'||n.status===status)&&(!query||`${n.title}${n.subtitle}${n.summary}`.toLowerCase().includes(query.toLowerCase()))).map(n=>n.key)),[data,status,query]);
  const nodes=useMemo(()=>layout((data?.nodes??[]).filter(n=>visible.has(n.key))),[data,visible]);
  const edges=useMemo<Edge[]>(()=>(data?.edges??[]).filter(e=>visible.has(e.source)&&visible.has(e.target)).map(e=>({id:e.id,source:e.source,target:e.target,label:e.relation==='depends_on'?'依赖':'',animated:e.relation==='depends_on',style:{stroke:e.relation==='depends_on'?'#f59e0b':'#94a3b8'}})),[data,visible]);
  const updates=(data?.recent_updates??[]).filter(u=>!selected||u.node_key===selected.key);
  return <main className="aese-atlas">
    <header className="aese-atlas-header"><button className="aese-atlas-back" onClick={onExit} title="返回企业沙盘"><ArrowLeft/></button><div><span>AESE · ENTERPRISE SIMULATION ANATOMY</span><h1>智能企业仿真完成体</h1></div><div className="aese-atlas-kpis"><div><small>总体完成</small><b>{data?.summary.overall_progress??'-'}%</b></div><div><small>目标构件</small><b>{data?.summary.node_count??'-'}</b></div><div><small>已完成</small><b>{data?.summary.by_status.completed??0}</b></div><div><small>建设 / 验证</small><b>{(data?.summary.by_status.building??0)+(data?.summary.by_status.validating??0)}</b></div></div><button className="aese-atlas-refresh" onClick={()=>void load()} title="刷新"><RefreshCw className={loading?'spin':''}/></button></header>
    <div className="aese-atlas-tools"><div><Search/><input placeholder="搜索模型、场景、Agent 或实验" value={query} onChange={e=>setQuery(e.target.value)}/>{query&&<button onClick={()=>setQuery('')} title="清空"><X/></button>}</div><label><Filter/><select value={status} onChange={e=>setStatus(e.target.value as Status|'all')}><option value="all">全部状态</option>{Object.entries(labels).map(([key,value])=><option key={key} value={key}>{value}</option>)}</select></label><span>实时数据 · {data?.summary.last_updated_at?new Date(data.summary.last_updated_at).toLocaleString('zh-CN'):'等待连接'}</span></div>
    {error?<section className="aese-atlas-error"><Activity/><b>无法读取 IAOS System Atlas</b><p>{error}</p><button onClick={()=>void load()}>重新连接</button></section>:<div className="aese-atlas-body"><section className="aese-atlas-canvas"><div className="aese-atlas-legend">{Object.entries(labels).map(([key,value])=><span key={key}><i style={{background:colors[key as Status]}}/>{value}</span>)}</div><ReactFlow nodes={nodes} edges={edges} nodeTypes={nodeTypes} onNodeClick={(_,n)=>setSelected(n.data as AtlasNode)} fitView minZoom={.35} maxZoom={1.7} proOptions={{hideAttribution:true}}><Background color="#cbd5e1" gap={22}/><Controls position="bottom-left"/><MiniMap pannable zoomable nodeColor={n=>colors[(n.data as AtlasNode).status]}/></ReactFlow></section><aside className={`aese-atlas-detail ${selected?'open':''}`}>{selected?<><div className="aese-atlas-detail-head"><div><span>{selected.layer} · {selected.owner}</span><h2>{selected.title}</h2><p>{selected.subtitle}</p></div><button onClick={()=>setSelected(null)} title="关闭"><X/></button></div><div className="aese-atlas-progress"><b>{selected.progress}%</b><div><i style={{width:`${selected.progress}%`,background:colors[selected.status]}}/></div><em style={{color:colors[selected.status]}}>{labels[selected.status]}</em></div><section><h3>仿真职责</h3><p>{selected.summary}</p></section><section><h3>完成体目标</h3><p>{selected.target_state}</p></section><section><h3>当前状态</h3><p>{selected.current_state}</p></section><section><h3>设计与实现依据</h3>{[...selected.design_refs,...selected.code_refs,...selected.evidence_refs].map(ref=><div className="aese-atlas-ref" key={ref}><ExternalLink/><code>{ref}</code></div>)}</section></>:<div className="aese-atlas-empty"><Activity/><b>点击构件查看详情</b><p>查看最终目标、当前完成度、设计文件和可验证证据。</p></div>}<section className="aese-atlas-timeline"><h3><Clock3/>进展记录</h3>{updates.slice(0,8).map(u=><article key={u.id}><i/><div><b>{u.summary}</b><p>{u.detail}</p><span>{new Date(u.occurred_at).toLocaleString('zh-CN')} · {u.created_by}</span></div></article>)}</section></aside></div>}
    <footer className="aese-atlas-footer"><CheckCircle2/>最终蓝图与当前状态来自 IAOS 数据库，不由页面静态估算。</footer>
  </main>
}
