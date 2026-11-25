#!/usr/bin/env bash
# tuzoq.yaml (healthcheck yoʻq) ni koʻtaradi va 3 soniyadan keyin holatni koʻrsatadi.
#   ./tuzoq.sh       toza volume bilan (baza endi yaratiladi)
#   ./tuzoq.sh bor   volume'dagi tayyor baza bilan
c() { docker compose -f tuzoq.yaml "$@"; }
if [ "${1:-}" = bor ]; then c down >/dev/null 2>&1; else c down -v >/dev/null 2>&1; fi
c up -d >/dev/null 2>&1
sleep 3
c ps -a --format '{{.Service}}: {{.Status}}'
c logs app --no-log-prefix
