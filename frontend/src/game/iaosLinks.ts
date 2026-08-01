import type { GameProjection } from "./types";

function iaosBaseUrl(): string {
  const configured = import.meta.env.VITE_IAOS_UI_BASE_URL?.trim();
  if (configured) return configured.replace(/\/$/, "");
  if (typeof window === "undefined") return "http://127.0.0.1:3000";
  return `${window.location.protocol}//${window.location.hostname}:3000`;
}

export function financeWorkspaceUrl(
  projection: Pick<GameProjection, "tenant_id" | "case_code">,
  view: "ledger" | "reports" | "operations",
): string {
  const query = new URLSearchParams({
    tenant: projection.tenant_id,
    case: projection.case_code,
    view,
  });
  return `${iaosBaseUrl()}/?${query.toString()}#finance_workspace`;
}

export function m9KnowledgeUrl(projection: Pick<GameProjection,"tenant_id"|"case_code">, capability:string):string {
  const query=new URLSearchParams({tenant:projection.tenant_id,case:projection.case_code,article_id:"KB-HCTM-M9-INCORPORATION",capability});
  return `${iaosBaseUrl()}/?${query.toString()}#knowledge_center`;
}
