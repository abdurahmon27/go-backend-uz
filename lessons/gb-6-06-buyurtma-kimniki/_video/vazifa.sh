#!/usr/bin/env bash
# TASK.md'dagi namoyish: starter'da Jasur Azizaning 41-buyurtmasini oʻzgartira oladi; solution'da — 404.
# Ishga tushirish (shu papkada): bash _video/vazifa.sh > _video/vazifa.txt   (Docker, curl, goose; 8765, 5516 boʻsh)
set -u
export GOTOOLCHAIN=go1.27.4 CGO_ENABLED=0
lab=$(pwd)
show() { echo "\$ $*"; eval "$@" 2>&1; }
docker compose down -v >/dev/null 2>&1
docker compose up -d --wait >/dev/null 2>&1
export GOOSE_DRIVER=postgres GOOSE_MIGRATION_DIR=migrations GOOSE_DBSTRING=postgres://bozor:bozor-lab@localhost:5516/bozor
goose up >/dev/null 2>&1
export DATABASE_URL=$GOOSE_DBSTRING
for dir in . solution; do
	echo; echo "### $dir"
	docker compose exec -T db psql -q -f /lab/seed.sql >/dev/null
	work=$(mktemp -d)
	(cd "$dir" && go build -o "$work/api" ./cmd/api) || exit 1
	if ss -ltn | grep -q ':8765 '; then echo "8765 band"; exit 1; fi
	"$work/api" >"$work/server.log" 2>&1 &
	pid=$!
	for _ in $(seq 50); do curl -s localhost:8765/health >/dev/null && break; sleep 0.1; done
	cd "$work"
	curl -s -c jasur.txt localhost:8765/login -H "Content-Type: application/json" -d '{"email":"jasur@bozor.test","password":"choy-non-qand-27"}' >/dev/null
	show "curl -s -i -X PATCH -b jasur.txt localhost:8765/orders/41 -H \"Content-Type: application/json\" -d '{\"note\":\"Jasurning manzili\"}' | head -1"
	cd "$lab"
	show "docker compose exec -T db psql -c \"SELECT id, user_id, note FROM orders WHERE id = 41\""
	kill $pid; wait $pid 2>/dev/null
done
docker compose down -v >/dev/null 2>&1 && echo "docker compose down -v: ok"
