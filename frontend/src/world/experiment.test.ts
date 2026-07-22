import {describe,expect,it,vi} from "vitest"; import {loadExperiment} from "./experiment";
describe("experiment API",()=>{it("loads immutable evidence",async()=>{vi.stubGlobal("fetch",vi.fn().mockResolvedValue({ok:true,json:async()=>({runs:[],comparisons:[],strategy_evidence_ready:true})}));expect((await loadExperiment()).strategy_evidence_ready).toBe(true);vi.unstubAllGlobals()})});
