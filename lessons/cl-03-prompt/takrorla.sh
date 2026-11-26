#!/usr/bin/env bash
# takrorla.sh promptlar/2-kontekst.txt [N]
# Promptni Claude'ga N marta (standart 3) beradi, javoblarni javoblar/<prompt>/ ga yozadi:
#   <i>.txt   javob matni, oʻzgartirilmagan
#   <i>.json  sana, model, tokenlar, vaqt
# Javob har safar boshqacha boʻlishi mumkin: shuning uchun uchta javob saqlaymiz.
set -euo pipefail
prompt=$(realpath "$1"); n=${2:-3}; model=${MODEL:-claude-sonnet-5-5}
dir=$(realpath javoblar)/$(basename "$prompt" .txt); mkdir -p "$dir"
cd "$(mktemp -d)"   # boʻsh papka: CLAUDE.md va loyiha fayllari taʼsir qilmasin
for i in $(seq 1 "$n"); do
  # boʻsh system prompt, asboblarsiz: Claude faqat bizning promptni oladi
  claude -p --model "$model" --effort high --system-prompt "" --tools "" \
    --setting-sources "" --strict-mcp-config --disable-slash-commands \
    --no-session-persistence --output-format json < "$prompt" > raw.json
  jq -r .result raw.json > "$dir/$i.txt"
  jq --arg sana "$(date -u +%FT%TZ)" --arg model "$model" \
    --arg prompt "$(basename "$prompt")" \
    '{sana: $sana, model: $model, prompt: $prompt, effort: "high",
      stop_reason, duration_ms, output_tokens: .usage.output_tokens,
      input_tokens_jami: (.usage.input_tokens + .usage.cache_creation_input_tokens
                     + .usage.cache_read_input_tokens),
      models: (.modelUsage | keys)}' raw.json > "$dir/$i.json"
  echo "$dir/$i.txt: $(wc -w < "$dir/$i.txt") soʻz"
done
