import type { GameProjection, GameWorkItem } from "./types";

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

export function m9KnowledgeUrl(
  projection: Pick<GameProjection, "tenant_id" | "case_code" | "world_run_id">,
  item: Pick<GameWorkItem, "work_item_id" | "capability" | "kind" | "owner_id" | "owner_type">,
): string {
  const workspace =
    typeof window === "undefined"
      ? ""
      : window.localStorage.getItem("aese_genesis_workspace_id") ?? "";
  const query = new URLSearchParams({
    tenant: projection.tenant_id,
    case: projection.case_code,
    article_id: "KB-HCTM-M9-INCORPORATION",
    workspace_id: workspace,
    case_code: projection.case_code,
    world_run_id: projection.world_run_id,
    node_id: item.work_item_id,
    actor_id: item.owner_id,
    actor_type: item.owner_type,
    task_type: item.kind,
    capability: item.capability,
  });
  return `${iaosBaseUrl()}/?${query.toString()}#knowledge_center`;
}
