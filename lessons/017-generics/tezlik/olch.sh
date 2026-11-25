#!/usr/bin/env bash
# olch.sh: har bir benchmark'ni 7 marta ishga tushirib, medianasini chiqaradi (shovqinni kamaytirish uchun).
# ./tezlik/olch.sh   (lessons/017-generics papkasidan)
set -euo pipefail
cd "$(dirname "$0")"
go test -run '^$' -bench . -count 7 -cpu 1 . |
	awk '/^Benchmark/ { v[$1] = v[$1] " " $3 }
	END {
		for (k in v) {
			n = split(substr(v[k], 2), a, " ")
			for (i = 1; i <= n; i++) for (j = i + 1; j <= n; j++) if (a[j] < a[i]) { t = a[i]; a[i] = a[j]; a[j] = t }
			printf "%-20s %8.0f ns/op (mediana, %d marta)\n", k, a[int((n + 1) / 2)], n
		}
	}' | sort -k2 -n -r
