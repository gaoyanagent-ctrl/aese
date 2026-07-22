#!/usr/bin/env bash
set -euo pipefail
repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$repo_root"
BASE_URL="${IAOS_BASE_URL:-http://127.0.0.1:8082}" TENANT_ID="${IAOS_TENANT_ID:-tenant-hctm}" TOKEN="${IAOS_TOKEN:-}" COMMIT_SHA="${ATLAS_COMMIT_SHA:-$(git rev-parse HEAD)}"
if [ -z "$TOKEN" ]; then TOKEN="$(curl -fsS "$BASE_URL/api/v1/dev/token?tenant_id=$TENANT_ID&roles=admin" | jq -er '.token')"; fi
count=0
while IFS= read -r file; do
  payload="$(jq --arg commit "$COMMIT_SHA" '{update_key,node_key,update_type,summary,detail,status_after:(.status_after//""),progress_after,current_state:(.current_state//""),source_ref:.evidence_ref,commit_sha:(if (.commit_sha//"")=="" then $commit else .commit_sha end),occurred_at:(.occurred_at//"")}' "$file")"
  curl -fsS -X POST "$BASE_URL/api/v1/system-atlas/updates" -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' --data-binary "$payload" >/dev/null
  count=$((count+1))
done < <(find atlas-updates -maxdepth 1 -name '*.json' ! -name schema.json -print0 | while IFS= read -r -d '' file; do jq -r --arg file "$file" '[(.occurred_at//""),.update_key,$file]|@tsv' "$file"; done | sort -k1,1 -k2,2 | cut -f3-)
echo "system-atlas sync: ok ($count declaration(s))"
