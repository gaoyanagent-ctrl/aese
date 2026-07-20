#!/usr/bin/env bash
set -euo pipefail

base_ref="${1:-}"
head_ref="${2:-HEAD}"
if [ -n "$base_ref" ]; then
  mapfile -t changed < <(git diff --name-only "$base_ref"..."$head_ref" --)
else
  mapfile -t changed < <({ git diff --name-only HEAD --; git ls-files --others --exclude-standard; } | sort -u)
fi
substantive=false
progress_changed=false
manifests=()
for path in "${changed[@]}"; do
  case "$path" in
    cmd/*|internal/*|frontend/src/*|scenario-packs/*|docs/designs/*|docs/decisions/*|docs/plans/*)
      substantive=true ;;
  esac
  [ "$path" = "docs/progress-log.md" ] && progress_changed=true
  case "$path" in atlas-updates/*.json) [ "$path" = "atlas-updates/schema.json" ] || manifests+=("$path");; esac
done
if [ "$substantive" != true ]; then echo "system-atlas tracking: no substantive change"; exit 0; fi
[ "$progress_changed" = true ] || { echo "system-atlas tracking: substantive AESE change requires docs/progress-log.md" >&2; exit 1; }
[ "${#manifests[@]}" -gt 0 ] || { echo "system-atlas tracking: substantive change requires an atlas-updates/*.json declaration" >&2; exit 1; }
for file in "${manifests[@]}"; do
  jq -e '.schema_version==1 and (.update_key|type=="string" and test("^[a-z0-9][a-z0-9._-]{7,239}$")) and (.node_key|type=="string" and length>=2 and length<=160) and (.update_type|IN("design","implementation","test","release","decision","risk","status")) and (.summary|type=="string" and length>=4 and length<=300) and (.detail|type=="string" and length>=4) and (.evidence_ref|type=="string" and length>=2 and length<=500) and ((.status_after//"planned")|IN("planned","designed","building","validating","completed","blocked","deferred")) and ((has("progress_after")|not) or (.progress_after|type=="number" and floor==. and .>=0 and .<=100))' "$file" >/dev/null || { echo "system-atlas tracking: invalid declaration $file" >&2; exit 1; }
  evidence="$(jq -r '.evidence_ref' "$file")"; [ -e "$evidence" ] || { echo "system-atlas tracking: evidence_ref does not exist: $evidence" >&2; exit 1; }
done
duplicates="$(find atlas-updates -maxdepth 1 -name '*.json' ! -name schema.json -print0 | xargs -0 -r jq -r '.update_key' | sort | uniq -d)"
[ -z "$duplicates" ] || { echo "system-atlas tracking: duplicate update_key: $duplicates" >&2; exit 1; }
echo "system-atlas tracking: ok (${#manifests[@]} declaration(s))"
