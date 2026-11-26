#!/bin/sh
# "oʻzbek" ning uzunligi toʻrt tilda. Kerak: node, python3, php, go (yoʻgʻi oʻtkazib yuboriladi).
# Natija (2026-11-26: Node v24.10.0, Python 3.13.15, PHP 8.4.26, go1.26.8):
#   js .length 6 · python len 6 · php strlen 7 · go len 7
cd "$(dirname "$0")/.." || exit 1
command -v node    >/dev/null && echo "js     .length  $(node -p '"oʻzbek".length')"
command -v python3 >/dev/null && echo "python len      $(python3 -c 'print(len("oʻzbek"))')"
command -v php     >/dev/null && echo "php    strlen   $(php -r 'echo strlen("oʻzbek");')"
command -v go      >/dev/null && echo "go     len      $(go run hook/main.go)"
