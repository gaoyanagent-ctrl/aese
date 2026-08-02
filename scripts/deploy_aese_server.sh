#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
env_file="${AESE_ENV_FILE:-$repo_root/.env}"
listen="${AESE_LISTEN:-:8090}"
iaos_base_url="${IAOS_BASE_URL:-http://127.0.0.1:8082}"
request_timeout="${AESE_REQUEST_TIMEOUT:-90s}"
bin_dir="$repo_root/.aese-data/bin"
bin_file="$bin_dir/aese-server"
pid_file="$repo_root/.aese-data/aese-server.pid"
log_file="$repo_root/.aese-data/aese-server.log"
mode="${1:-deploy}"

load_env_file() {
  local file="$1" line name value mode_bits
  if [[ ! -f "$file" ]]; then
    return 0
  fi
  mode_bits="$(stat -c '%a' "$file")"
  if [[ "$mode_bits" != *00 ]]; then
    echo "refusing readable secret file $file (mode=$mode_bits, require *00)" >&2
    return 1
  fi
  while IFS= read -r line || [[ -n "$line" ]]; do
    line="${line%$'\r'}"
    [[ -z "$line" || "$line" == \#* ]] && continue
    if [[ ! "$line" =~ ^([A-Za-z_][A-Za-z0-9_]*)=(.*)$ ]]; then
      echo "invalid environment assignment in $file" >&2
      return 1
    fi
    name="${BASH_REMATCH[1]}"
    value="${BASH_REMATCH[2]}"
    if [[ "$value" =~ ^\"(.*)\"$ || "$value" =~ ^\'(.*)\'$ ]]; then
      value="${BASH_REMATCH[1]}"
    fi
    export "$name=$value"
  done < "$file"
}

validate_creative_config() {
  local configured=0
  [[ -n "${MINMAX_API_KEY:-}" ]] && configured=$((configured + 1))
  [[ -n "${MINMAX_API_BASE:-}" ]] && configured=$((configured + 1))
  [[ -n "${MINMAX_MODEL:-}" ]] && configured=$((configured + 1))
  if [[ "$configured" -ne 0 && "$configured" -ne 3 ]]; then
    echo "MiniMax configuration must set MINMAX_API_KEY, MINMAX_API_BASE and MINMAX_MODEL together" >&2
    return 1
  fi
  if [[ "$configured" -eq 3 ]]; then
    echo "creative provider: MiniMax (${MINMAX_MODEL}); API key present"
  else
    echo "creative provider: deterministic fallback (MiniMax not configured)"
  fi
}

stop_existing_server() {
  local candidates="" pid cmd
  if [[ -f "$pid_file" ]]; then
    candidates="$(cat "$pid_file")"
  fi
  candidates="$candidates $(ss -ltnp 2>/dev/null | sed -n -E "/:${listen#:}[[:space:]]/s/.*pid=([0-9]+).*/\\1/p")"
  for pid in $candidates; do
    [[ "$pid" =~ ^[0-9]+$ && -r "/proc/$pid/cmdline" ]] || continue
    if ! cmd="$(tr '\0' ' ' < "/proc/$pid/cmdline" 2>/dev/null)"; then
      continue
    fi
    if [[ "$cmd" != *aese-server* ]]; then
      echo "refusing to stop non-AESE process $pid on $listen" >&2
      return 1
    fi
    kill "$pid" 2>/dev/null || true
  done
  for _ in $(seq 1 20); do
    if ! curl -fsS "http://127.0.0.1:${listen#:}/health" >/dev/null 2>&1; then
      return 0
    fi
    sleep 0.1
  done
  echo "existing AESE server did not stop cleanly on $listen" >&2
  return 1
}

load_env_file "$env_file"
validate_creative_config
if [[ "$mode" == "--check-config" ]]; then
  exit 0
fi
if [[ "$mode" != "deploy" ]]; then
  echo "usage: $0 [deploy|--check-config]" >&2
  exit 2
fi

mkdir -p "$bin_dir"
stop_existing_server

(
  cd "$repo_root"
  GOPROXY="${GOPROXY:-https://goproxy.cn,direct}" GOSUMDB="${GOSUMDB:-off}" \
    go build -buildvcs=false -o "$bin_file" ./cmd/aese-server
)

setsid "$bin_file" \
  --listen "$listen" \
  --pack-dir "$repo_root/scenario-packs/hctm" \
  --iaos-base-url "$iaos_base_url" \
  --request-timeout "$request_timeout" \
  >"$log_file" 2>&1 < /dev/null &
pid="$!"
printf '%s\n' "$pid" > "$pid_file"

for _ in $(seq 1 50); do
  if ! kill -0 "$pid" 2>/dev/null; then
    echo "AESE server exited before becoming healthy" >&2
    tail -80 "$log_file" >&2 || true
    exit 1
  fi
  if curl -fsS "http://127.0.0.1:${listen#:}/health" >/dev/null 2>&1; then
    status="$(curl -fsS "http://127.0.0.1:${listen#:}/api/aese/v1/game/creative/status")"
    planning_status="$(curl -fsS "http://127.0.0.1:${listen#:}/api/aese/v1/world/plant-build/planning-status")"
    if [[ -n "${MINMAX_API_KEY:-}" ]]; then
      if ! jq -e '.state=="connected" and .provider=="MiniMax"' <<<"$status" >/dev/null; then
        echo "AESE started but M9 creative provider is not connected" >&2
        jq '{state,provider,model,prompt_version}' <<<"$status" >&2
        exit 1
      fi
      if ! jq -e '.state=="connected" and .provider=="MiniMax"' <<<"$planning_status" >/dev/null; then
        echo "AESE started but M10 plant-planning provider is not connected" >&2
        jq '{state,provider,model,prompt_version}' <<<"$planning_status" >&2
        exit 1
      fi
    fi
    echo "AESE server deployed (pid=$pid listen=$listen log=$log_file)"
    jq -n --argjson creative "$status" --argjson planning "$planning_status" \
      '{creative:($creative|{state,provider,model,base_url_host,prompt_version}),plant_planning:($planning|{state,provider,model,prompt_version})}'
    exit 0
  fi
  sleep 0.2
done

echo "AESE server failed health check" >&2
tail -80 "$log_file" >&2 || true
exit 1
