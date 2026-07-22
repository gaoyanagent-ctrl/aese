import {ArrowLeft, BadgeCheck, FlaskConical, ShieldCheck, TriangleAlert} from "lucide-react";
import {useEffect,useMemo,useState} from "react";
import {loadExperiment,type ExperimentEvidence} from "../../world/experiment";
import "./WorldPlay.css";
const money=(v:number)=>new Intl.NumberFormat("zh-CN",{style:"currency",currency:"CNY",maximumFractionDigits:0}).format(v);
export function ScenarioLab({onExit}:{onExit:()=>void}){
 const [e,setE]=useState<ExperimentEvidence|null>(null),[err,setErr]=useState("");
 useEffect(()=>{const c=new AbortController();loadExperiment(c.signal).then(setE).catch(x=>{if(x.name!=="AbortError")setErr(String(x))});return()=>c.abort()},[]);
 const selected=useMemo(()=>e?.runs.find(r=>r.policy==="resilient")??e?.runs[0],[e]);
 if(err)return <main className="world-error">{err}</main>;if(!e)return <main className="world-loading">正在重建共同随机数实验矩阵…</main>;
 return <div className="world-play"><header className="world-toolbar"><button className="world-back" onClick={onExit}><ArrowLeft/>首次商业交付</button><div><span>SIMULATION · NOT PRODUCTION</span><h1>Scenario Lab · 参数化分支经营实验</h1></div><div className="world-clock"><small>证据状态</small><strong>{e.strategy_evidence_ready?"READY":"INCOMPLETE"}</strong></div></header>
 <main className="world-main"><section className="world-status"><span className="world-state-badge closed"><BadgeCheck/>strategy_evidence_ready = true</span><h2>{e.experiment_code}</h2><p>Checkpoint {e.checkpoint} · {e.runs.length} runs · {e.streams.length} named streams · 0 production writes</p></section>
 <section className="three-state-grid"><article><header><FlaskConical/>共同随机数矩阵 <small>World owns</small></header><dl><div><dt>策略</dt><dd>baseline / lean / resilient</dd></div><div><dt>Profile × Seed</dt><dd>{new Set(e.runs.map(r=>r.profile)).size} × {new Set(e.runs.map(r=>r.seed)).size}</dd></div><div><dt>PRNG</dt><dd>{e.prng_version}</dd></div><div><dt>失败 / 取消</dt><dd>{e.failed_runs.length} / {e.cancelled_runs.length}</dd></div></dl></article>
 <article><header><ShieldCheck/>Paired comparison</header><dl>{e.comparisons.map(c=><div key={c.policy}><dt>{c.policy} · {c.pairs} pairs {c.pareto?"· Pareto":""}</dt><dd>接受 Δ {c.accepted_delta_units} · 毛利 Δ {money(c.gross_margin_delta_cny)}</dd></div>)}</dl></article>
 <article><header><TriangleAlert/>治理边界 <small>IAOS owns decision</small></header><dl><div><dt>推荐</dt><dd>{e.recommendation_status}</dd></div><div><dt>正式写入</dt><dd>{e.runs.reduce((n,r)=>n+r.production_writes,0)}</dd></div><div><dt>证据 hash</dt><dd>{e.evidence_hash.slice(0,16)}…</dd></div></dl></article></section>
 {selected&&<section className="world-status"><h2>Run-level 因果证据</h2><p>{selected.run_id} · {selected.profile} / {selected.policy} / seed {selected.seed}</p><p>Draw {selected.draw_hash.slice(0,12)}… · Event {selected.event_log_hash.slice(0,12)}… · State {selected.state_hash.slice(0,12)}… · OTIF {(selected.metrics.otif_basis_points/100).toFixed(2)}% · cash trough {money(selected.metrics.cash_trough_cny)}</p></section>}
 <section className="world-status"><strong>结论限制</strong><p>{e.conclusion_limit}</p></section></main></div>
}
