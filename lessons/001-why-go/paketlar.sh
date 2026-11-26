#!/bin/sh
# Go binary uchun "bundle analyzer": mashina kodi (funksiyalar) baytlarini paketlar bo'yicha qo'shadi.
# Ishlatish: ./paketlar.sh <binary> [qator soni, standart 6]
# Eslatma: -ldflags="-s" bilan build qilingan faylda simvollar jadvali yo'q, natija bo'sh bo'ladi.
go tool nm -size "$1" | awk '
  $3 == "T" || $3 == "t" {
    size = $2
    name = $0; sub(/^ *[0-9a-f]+ +[0-9]+ +[A-Za-z] +/, "", name)   # nomda bo`sh joy bo`lishi mumkin
    sub(/^go:/, "", name); sub(/^type:\.(eq|hash)\./, "", name)     # kompilyator yasagan eq/hash funksiyalari
    sub(/\[.*/, "", name)                                            # generic: Foo[...] -> Foo
    slash = 0; for (i = 1; i <= length(name); i++) if (substr(name, i, 1) == "/") slash = i
    rest = substr(name, slash + 1); dot = index(rest, ".")
    pkg = (dot > 0) ? substr(name, 1, slash + dot - 1) : "runtime(asm)"   # paketsiz nom = runtime assembler kodi
    sum[pkg] += size; total += size
  }
  END { for (p in sum) printf "%8d %5.1f%% %s\n", sum[p], 100 * sum[p] / total, p; printf "%8d 100.0%% JAMI\n", total }
' | sort -rn | head -n "${2:-6}"
