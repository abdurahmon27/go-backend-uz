#!/usr/bin/env bash
# Uy vazifasi javoblarini tekshiradi (kurs papkasidan: bash homework/javoblar.sh). 8078 va 9000 portlar boʻsh boʻlsin.
set -u
cd "$(dirname "$0")/.."
[ -f .env ] || { cp .env.example .env; sed -i "s/^DB_PASSWORD=$/DB_PASSWORD=javob-paroli/" .env; }
c() { docker compose -f "homework/$1.yaml" --env-file .env "${@:2}"; }

echo "== Oson: APP_PORT=9000"
APP_PORT=9000 c oson up -d --wait >/dev/null 2>&1; sleep 1
curl -s localhost:9000
c oson down -v >/dev/null 2>&1

for v in orta qiyin; do
  [ $v = orta ] && echo "== Oʻrta: restart: on-failure, healthcheck yoʻq" || echo "== Qiyin: server bazani oʻzi kutadi, healthcheck yoʻq"
  c $v up -d --build >/dev/null 2>&1; sleep 6
  c $v ps -a --format '{{.Service}}: {{.Status}}'
  [ $v = orta ] && docker inspect -f 'app qayta ishga tushdi: {{.RestartCount}} marta' $v-app-1
  c $v logs app --no-log-prefix | tail -3
  curl -s localhost:8078
  c $v down -v >/dev/null 2>&1
done

echo "== Guru: parol secret faylda"
c guru up -d --build --wait >/dev/null 2>&1; sleep 1
curl -s localhost:8078
docker inspect -f '{{range .Config.Env}}{{println .}}{{end}}' guru-db-1 guru-app-1 | grep -E 'PASSWORD|DATABASE_URL'
c guru exec db sh -c 'ls -l /run/secrets'
c guru down -v >/dev/null 2>&1
