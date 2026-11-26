#!/usr/bin/env bash
# uslub.sh <Default|Learning> [N]
# Claude Code'ga bitta vazifani (promptlar/5-claude-code.txt) bozor/ nusxasida beradi, N marta (standart 3).
# Claude fayllarni oʻqiy va tahrirlay oladi (Read, Edit, Write), buyruq ishlata olmaydi (Bash yoʻq).
#   javoblar/uslub-<style>/<i>.txt   Claude'ning javobi
#   javoblar/uslub-<style>/<i>.diff  Claude faylga nima yozgani (diff -u)
#   javoblar/uslub-<style>/<i>.json  sana, model, tokenlar, vaqt
set -euo pipefail
style=$1; n=${2:-3}; model=${MODEL:-claude-sonnet-5-5}
dir=$(realpath javoblar)/uslub-$style; mkdir -p "$dir"
prompt=$(realpath promptlar/5-claude-code.txt); asl=$(realpath bozor/main.go)
for v in $(env | grep -oE '^(CLAUDECODE|CLAUDE_CODE_[A-Z_]+|CLAUDE_PID|CLAUDE_EFFORT|CLAUDE_AGENT_SDK_VERSION|AI_AGENT)=' || true); do
  unset "${v%=}"
done
for i in $(seq 1 "$n"); do
  tmp=$(mktemp -d); mkdir "$tmp/bozor"; cp "$asl" "$tmp/bozor/main.go"
  printf 'module bozor\n\ngo 1.26\n' > "$tmp/go.mod"
  # /config → Output style ham shu faylga yozadi (code.claude.com/docs/en/output-styles)
  mkdir "$tmp/.claude"; printf '{"outputStyle": "%s"}\n' "$style" > "$tmp/.claude/settings.local.json"
  cd "$tmp"
  claude -p --model "$model" --effort high \
    --setting-sources local --strict-mcp-config --tools "Read,Edit,Write,Glob,Grep" \
    --permission-mode acceptEdits --no-session-persistence --output-format json < "$prompt" > raw.json
  jq -r .result raw.json > "$dir/$i.txt"
  (cd "$tmp" && diff -u --label a/bozor/main.go --label b/bozor/main.go "$asl" bozor/main.go > "$dir/$i.diff" || true)
  jq --arg sana "$(date -u +%FT%TZ)" --arg model "$model" --arg style "$style" \
    '{sana: $sana, model: $model, outputStyle: $style, effort: "high", stop_reason, num_turns,
      duration_ms, output_tokens: .usage.output_tokens, models: (.modelUsage | keys)}' raw.json > "$dir/$i.json"
  echo "$dir/$i: $(wc -w < "$dir/$i.txt") soʻz, diff $(grep -c '^[+-][^+-]' "$dir/$i.diff" || true) qator"
  echo "$tmp" >> "${TAKRORLA_LOG:-/dev/null}"; cd - > /dev/null; rm -rf "$tmp"
done
