#!/usr/bin/env bash
# olch.sh: har bir benchmark'ni 7 marta ishga tushirib, ns/op medianasini va allocs/op ni chiqaradi.
# ./tezlik/olch.sh   (lessons/018-iterators papkasidan)
set -euo pipefail
cd "$(dirname "$0")"
go test -run '^$' -bench . -benchmem -count 7 -cpu 1 . |
	awk '/^Benchmark/ { v[$1] = v[$1] " " $3; al[$1] = $(NF-1) }
	END {
		m = split("BenchmarkOddiy BenchmarkValues BenchmarkOz BenchmarkSeq", tartib, " ")
		for (q = 1; q <= m; q++) {
			k = tartib[q]
			n = split(substr(v[k], 2), a, " ")
			for (i = 1; i <= n; i++) for (j = i + 1; j <= n; j++) if (a[j] < a[i]) { t = a[i]; a[i] = a[j]; a[j] = t }
			printf "%-16s %6.0f ns/op  %s allocs/op\n", k, a[int((n + 1) / 2)], al[k]
		}
	}'
