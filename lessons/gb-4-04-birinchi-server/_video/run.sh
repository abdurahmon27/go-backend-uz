#!/usr/bin/env bash
# Videodagi bitta qadamni xuddi videodagidek ishga tushiradi: _video/handlers.go (GB-4.03 handler'lari) va
# qadamning main.go'si vaqtinchalik bozor papkasiga shu go.mod bilan koʻchiriladi, u yerda berilgan buyruq bajariladi.
#   ./_video/run.sh 02-server go run .        (server: toʻxtatish — Ctrl+C)
#   ./_video/run.sh final go run .            (final = papkadagi main.go)
# Server ishlab turganda boshqa terminalda: curl http://localhost:8765/health
set -u
export GOTOOLCHAIN=${GOTOOLCHAIN:-go1.27.4} CGO_ENABLED=0
here=$(cd "$(dirname "$0")/.." && pwd)
step=${1:?qadam nomi kerak, masalan 02-server}
shift
t=$(mktemp -d)/bozor && mkdir -p "$t"
cp "$here/go.mod" "$t/"
if [ "$step" = final ]; then
	cp "$here/main.go" "$t/"
else
	cp "$here/_video/handlers.go" "$here/_video/$step/"* "$t/"
fi
cd "$t" || exit 1
"$@"
echo "(chiqish kodi: $?)"
