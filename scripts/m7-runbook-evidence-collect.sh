#!/usr/bin/env bash
set -euo pipefail

: "${AESE_BASE:?Missing AESE_BASE, e.g. http://127.0.0.1:8090}"
: "${IAOS_BASE:?Missing IAOS_BASE, e.g. http://127.0.0.1:8082}"
: "${TENANT_ID:?Missing TENANT_ID}"

PACK_DIR="${PACK_DIR:-scenario-packs/hctm}"
PACK_KEY="${PACK_KEY:-hctm}"
STORY_KEY="${STORY_KEY:-order-expedite-01}"
RUN_ID_PREFIX="${RUN_ID_PREFIX:-m7-evidence}"
RUN_ID="${RUN_ID:-${RUN_ID_PREFIX}-$(date -u +%Y%m%dT%H%M%SZ)}"
CLI_APPLY="${CLI_APPLY:-1}"
UI_APPLY="${UI_APPLY:-1}"
DRY_RUN="${DRY_RUN:-0}"
OUTPUT_DIR="${OUTPUT_DIR:-artifacts/m7-acceptance/$(date -u +%Y%m%dT%H%M%SZ)}"

if [[ "$DRY_RUN" == "1" ]]; then
  CLI_APPLY=0
  UI_APPLY=0
fi

REQUIRED_TOOLS=(curl jq go tee)
for tool in "${REQUIRED_TOOLS[@]}"; do
  if ! command -v "$tool" >/dev/null 2>&1; then
    echo "missing required tool: $tool" >&2
    exit 1
  fi
 done

mkdir -p "$OUTPUT_DIR"

normalize_url() {
  printf '%s' "${1%/}"
}

AESE_BASE="$(normalize_url "$AESE_BASE")"
IAOS_BASE="$(normalize_url "$IAOS_BASE")"
IAOS_API_ROOT="${IAOS_API_ROOT:-$IAOS_BASE}"
IAOS_API_ROOT="$(normalize_url "$IAOS_API_ROOT")"
IAOS_API_BASE="${IAOS_API_BASE:-$IAOS_API_ROOT}"
IAOS_API_BASE="$(normalize_url "$IAOS_API_BASE")"
IAOS_TOKEN_BASE="${IAOS_TOKEN_BASE:-$IAOS_API_BASE/api/v1}"
IAOS_TOKEN_BASE="$(normalize_url "$IAOS_TOKEN_BASE")"

: "${IAOS_TOKEN:=}"
if [[ -z "$IAOS_TOKEN" ]]; then
  IAOS_TOKEN="$(curl -fsS "$IAOS_TOKEN_BASE/dev/token?tenant_id=${TENANT_ID}&roles=admin" | jq -r '.token')"
fi
if [[ -z "$IAOS_TOKEN" || "$IAOS_TOKEN" == "null" ]]; then
  echo "failed to resolve IAOS_TOKEN" >&2
  exit 1
fi

RUN_VERSION=""
RUN_CURSOR=""
RESET_CONFIRMATION_TOKEN=""

api_call() {
  local method="$1"
  local url="$2"
  local body="${3:-}"
  local out="$4"
  shift 4

  local args=(curl -fsS -H "Authorization: Bearer $IAOS_TOKEN" -H 'Content-Type: application/json' "$@")
  if [[ "$method" == "GET" ]]; then
    "${args[@]}" "$url" | tee "$out"
  else
    "${args[@]}" -X "$method" -d "$body" "$url" | tee "$out"
  fi
}

record_run_context() {
  local file="$1"
  if [[ -f "$file" ]]; then
    RUN_VERSION="$(jq -r '.run.run_version // empty' "$file")"
    RUN_CURSOR="$(jq -r '.run.cursor // empty' "$file")"
  fi
}

run_action() {
  local action="$1"
  local apply="$2"
  local out_file="$3"
  local idempotency="${4:-}"
  local confirmation="${5:-}"

  local payload="{\"plan_hash\":\"${PLAN_HASH}\",\"apply\":${apply}"
  if [[ -n "$RUN_CURSOR" ]]; then
    payload+=" ,\"expected_cursor\":${RUN_CURSOR}"
  fi
  if [[ -n "$RUN_VERSION" && "$RUN_VERSION" != "null" ]]; then
    payload+=" ,\"run_version\":\"${RUN_VERSION}\""
  fi
  if [[ -n "$confirmation" ]]; then
    payload+=" ,\"confirmation_token\":\"${confirmation}\""
  fi
  payload+="}"

  if [[ "$apply" == "true" && "$action" != "preflight" && "$action" != "reset-plan" ]]; then
    local idem="${idempotency:-m7-$action-$(date -u +%Y%m%dT%H%M%SZ)}"
    api_call POST "$AESE_BASE/api/aese/v1/runs/$RUN_ID/$action" "$payload" "$OUTPUT_DIR/$out_file" -H "Idempotency-Key: $idem"
  else
    api_call POST "$AESE_BASE/api/aese/v1/runs/$RUN_ID/$action" "$payload" "$OUTPUT_DIR/$out_file"
  fi

  record_run_context "$OUTPUT_DIR/$out_file"

  if [[ "$action" == "reset-plan" ]]; then
    RESET_CONFIRMATION_TOKEN="$(jq -r '.run.outcome.reset_confirmation_token // .outcome.reset_confirmation_token // empty' "$OUTPUT_DIR/$out_file")"
  fi
}

emit_summary() {
  cat > "$OUTPUT_DIR/summary.txt" <<EOF_SUM
run_id=${RUN_ID}
tenant=${TENANT_ID}
pack=${PACK_KEY}
story=${STORY_KEY}
ui_apply=${UI_APPLY}
cli_apply=${CLI_APPLY}
dry_run=${DRY_RUN}
plan_hash=${PLAN_HASH}
run_version=${RUN_VERSION}
run_cursor=${RUN_CURSOR}
cli_artifacts=$OUTPUT_DIR
EOF_SUM
}

collect_cli_contracts() {
  local mode="apply"
  local verify_ok=0
  [[ "$CLI_APPLY" == "1" ]] || mode="dry-run"

  go run ./cmd/aese apply "$PACK_DIR" \
    --story "$STORY_KEY" \
    --target "$IAOS_API_BASE" \
    --tenant "$TENANT_ID" \
    --run-id "${RUN_ID}-cli-apply" \
    --actor aese-cli \
    $(if [[ "$CLI_APPLY" == "1" ]]; then echo --apply; fi) \
    | tee "$OUTPUT_DIR/05-cli-apply-${mode}.json"

  go run ./cmd/aese replay "$PACK_DIR" \
    --story "$STORY_KEY" \
    --target "$IAOS_API_BASE" \
    --tenant "$TENANT_ID" \
    --run-id "${RUN_ID}-cli-replay" \
    --actor aese-cli \
    $(if [[ "$CLI_APPLY" == "1" ]]; then echo --apply; fi) \
    | tee "$OUTPUT_DIR/06-cli-replay-${mode}.json"

  for attempt in $(seq 1 10); do
    if go run ./cmd/aese verify "$PACK_DIR" \
      --story "$STORY_KEY" \
      --target "$IAOS_API_BASE" \
      --tenant "$TENANT_ID" \
      --run-id "${RUN_ID}-cli-verify" \
      | tee "$OUTPUT_DIR/07-cli-verify-attempt-${attempt}.json"; then
      cp "$OUTPUT_DIR/07-cli-verify-attempt-${attempt}.json" "$OUTPUT_DIR/07-cli-verify.json"
      verify_ok=1
      break
    fi
    sleep 1
  done
  if [[ "$verify_ok" != "1" ]]; then
    echo "CLI verify did not converge after 10 attempts" >&2
    return 1
  fi

  go run ./cmd/aese reset "$PACK_DIR" \
    --story "$STORY_KEY" \
    --target "$IAOS_API_BASE" \
    --tenant "$TENANT_ID" \
    --run-id "${RUN_ID}-cli-reset" \
    --order-id "" \
    $(if [[ "$CLI_APPLY" == "1" ]]; then echo --apply; fi) \
    | tee "$OUTPUT_DIR/08-cli-reset-${mode}.json"
}

main() {
  api_call GET "$AESE_BASE/health" '{}' "$OUTPUT_DIR/health-aese.json"
  api_call GET "$AESE_BASE/ready" '{}' "$OUTPUT_DIR/ready-aese.json"
  api_call GET "$IAOS_API_ROOT/health" '{}' "$OUTPUT_DIR/health-iaos.json"
  if ! api_call GET "$IAOS_API_ROOT/ready" '{}' "$OUTPUT_DIR/ready-iaos.json"; then
    printf '%s\n' '{"status":"unavailable","message":"iaos ready endpoint not available","base":"'"$IAOS_API_ROOT"'"}' > "$OUTPUT_DIR/ready-iaos.json"
  fi

  api_call POST "$AESE_BASE/api/aese/v1/runs/plan" '{"story_key":"'"$STORY_KEY"'","pack_dir":"'"$PACK_DIR"'"}' "$OUTPUT_DIR/00-plan.json"
  PLAN_HASH="$(jq -r '.plan_hash // empty' "$OUTPUT_DIR/00-plan.json")"
  if [[ -z "$PLAN_HASH" || "$PLAN_HASH" == "null" ]]; then
    echo "failed to read plan_hash from 00-plan.json" >&2
    exit 1
  fi

  api_call POST "$AESE_BASE/api/aese/v1/runs" '{"target":"'"$IAOS_API_ROOT"'","tenant":"'"$TENANT_ID"'","story_key":"'"$STORY_KEY"'","plan_hash":"'"$PLAN_HASH"'","run_id":"'"$RUN_ID"'","actor":"aese-cli","token":"'"$IAOS_TOKEN"'","pack_dir":"'"$PACK_DIR"'"}' "$OUTPUT_DIR/01-run-create.json"
  record_run_context "$OUTPUT_DIR/01-run-create.json"

  run_action preflight "$([[ "$UI_APPLY" == "1" ]] && echo true || echo false)" "02-preflight.json"
  run_action initialize "$([[ "$UI_APPLY" == "1" ]] && echo true || echo false)" "03-initialize.json"
  run_action run-to-end "$([[ "$UI_APPLY" == "1" ]] && echo true || echo false)" "04-run-to-end.json"
  run_action analyze "$([[ "$UI_APPLY" == "1" ]] && echo true || echo false)" "05-analyze.json"
  run_action verify "$([[ "$UI_APPLY" == "1" ]] && echo true || echo false)" "06-verify.json"
  run_action reset-plan false "07-reset-plan.json"

  if [[ -n "$RESET_CONFIRMATION_TOKEN" ]]; then
    run_action reset "$([[ "$UI_APPLY" == "1" ]] && echo true || echo false)" "08-reset.json" "" "$RESET_CONFIRMATION_TOKEN"
  else
    echo "reset token missing; skip reset execute" >&2
  fi

  collect_cli_contracts
  emit_summary
  echo "Evidence artifacts saved to: $OUTPUT_DIR"
}

main
