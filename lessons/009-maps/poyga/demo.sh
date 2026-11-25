#!/usr/bin/env bash
# Record a real failure with the message used in the lesson. Detection is best effort.
set -eu
cd "$(dirname "$0")/.."
work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT
# Build once so each retry only runs the demonstration, without compilation overhead.
timeout 10s go build -o "$work/poyga" ./poyga
for ((attempt = 0; attempt < 20; attempt++)); do
    status=0
    timeout 0.5s "$work/poyga" >"$work/output" 2>&1 || status=$?
    if [[ "$status" == 2 ]] && grep -q '^fatal error: concurrent map writes$' "$work/output"; then
        grep '^fatal error: concurrent map writes$' "$work/output" | head -n 1
        printf 'exit status %s\n' "$status"
        exit "$status"
    fi
done
printf '%s\n' 'Could not capture concurrent map writes; retry the demonstration.' >&2
exit 1
