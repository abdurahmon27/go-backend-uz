#!/usr/bin/env bash
# takrorla.sh <nom> <tizim.txt|-> <1-xabar.txt> [2-xabar.txt ...]
# Bitta suhbatni Claude bilan N marta (standart 3, N=... bilan oʻzgaradi) boshidan oxirigacha oʻtkazadi.
# Har suhbat yangi, boʻsh papkada; 2-xabardan boshlab oʻsha suhbat davom ettiriladi (--resume).
#   javoblar/<nom>/<i>-<k>.txt   i-suhbatdagi k-javob, oʻzgartirilmagan
#   javoblar/<nom>/<i>-<k>.json  sana, model, tokenlar, vaqt
# tizim.txt: ustoz koʻrsatmasi (Claude.ai'da Project koʻrsatmasi kabi); "-" = boʻsh tizim matni.
set -euo pipefail
nom=$1; tizim=$2; shift 2
n=${N:-3}; model=${MODEL:-claude-sonnet-5-5}
dir=$(realpath javoblar)/$nom; mkdir -p "$dir"
xabarlar=(); for f in "$@"; do xabarlar+=("$(realpath "$f")"); done
tizim_matn=""; [ "$tizim" != "-" ] && tizim_matn=$(cat "$tizim")
# Claude Code ichidan ishga tushirilsa, uning muhit oʻzgaruvchilari bola jarayonga oʻtmasin
for v in $(env | grep -oE '^(CLAUDECODE|CLAUDE_CODE_[A-Z_]+|CLAUDE_PID|CLAUDE_EFFORT|CLAUDE_AGENT_SDK_VERSION|AI_AGENT)=' || true); do
  unset "${v%=}"
done
for i in $(seq 1 "$n"); do
  tmp=$(mktemp -d); cd "$tmp"   # boʻsh papka: CLAUDE.md va loyiha fayllari taʼsir qilmasin
  sessiya=""; k=0
  for x in "${xabarlar[@]}"; do
    k=$((k + 1))
    davom=(); [ -n "$sessiya" ] && davom=(--resume "$sessiya")
    claude -p --model "$model" --effort high --system-prompt "$tizim_matn" --tools "" \
      --setting-sources "" --strict-mcp-config --disable-slash-commands \
      "${davom[@]}" --output-format json < "$x" > raw.json
    sessiya=$(jq -r .session_id raw.json)
    jq -r .result raw.json > "$dir/$i-$k.txt"
    jq --arg sana "$(date -u +%FT%TZ)" --arg model "$model" --arg xabar "$(basename "$x")" \
      --arg tizim "$(basename "$tizim")" --argjson k "$k" \
      '{sana: $sana, model: $model, tizim: $tizim, xabar: $xabar, navbat: $k, effort: "high",
        stop_reason, duration_ms, output_tokens: .usage.output_tokens,
        input_tokens_jami: (.usage.input_tokens + .usage.cache_creation_input_tokens
                       + .usage.cache_read_input_tokens),
        models: (.modelUsage | keys)}' raw.json > "$dir/$i-$k.json"
    echo "$dir/$i-$k.txt: $(wc -w < "$dir/$i-$k.txt") soʻz"
  done
  echo "$tmp" >> "${TAKRORLA_LOG:-/dev/null}"; cd - > /dev/null; rm -rf "$tmp"
done
