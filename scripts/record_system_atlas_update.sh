#!/usr/bin/env bash
set -euo pipefail

if [ "$#" -lt 3 ]; then
  echo "usage: $0 NODE_KEY UPDATE_TYPE SUMMARY [STATUS] [PROGRESS] [DETAIL] [SOURCE_REF] [COMMIT_SHA]" >&2
  exit 2
fi
BASE_URL="${IAOS_BASE_URL:-http://127.0.0.1:8082}" TENANT_ID="${IAOS_TENANT_ID:-tenant-hctm}" TOKEN="${IAOS_TOKEN:-}"
if [ -z "$TOKEN" ]; then TOKEN="$(curl -fsS "$BASE_URL/api/v1/dev/token?tenant_id=$TENANT_ID&roles=admin" | jq -er '.token')"; fi
NODE_KEY="$1" UPDATE_TYPE="$2" SUMMARY="$3" STATUS="${4:-}" PROGRESS="${5:-}" DETAIL="${6:-}" SOURCE_REF="${7:-}" COMMIT_SHA="${8:-}"
if [ -n "$PROGRESS" ]; then PROGRESS_JSON="$PROGRESS"; else PROGRESS_JSON="null"; fi
jq -n --arg node_key "$NODE_KEY" --arg update_type "$UPDATE_TYPE" --arg summary "$SUMMARY" \
  --arg status_after "$STATUS" --arg detail "$DETAIL" --arg source_ref "$SOURCE_REF" \
  --arg commit_sha "$COMMIT_SHA" --argjson progress_after "$PROGRESS_JSON" \
  '{node_key:$node_key,update_type:$update_type,summary:$summary,status_after:$status_after,progress_after:$progress_after,detail:$detail,source_ref:$source_ref,commit_sha:$commit_sha}' |
curl -fsS -X POST "$BASE_URL/api/v1/system-atlas/updates" -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' --data-binary @-
echo
