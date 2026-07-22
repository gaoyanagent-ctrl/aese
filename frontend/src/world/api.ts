import type { GenesisTrace } from './types';
export async function loadGenesisTrace(signal?: AbortSignal): Promise<GenesisTrace> {
  const response = await fetch('/api/aese/v1/world/genesis', { signal });
  if (!response.ok) throw new Error(`World API ${response.status}`);
  const trace = await response.json() as GenesisTrace;
  trace.frames = (trace.frames ?? []).map(frame => ({ ...frame, knowledge: frame.knowledge ?? [] }));
  return trace;
}
