#!/usr/bin/env bash
# Videodagi diagramma: fayl hajmi (MiB) va eng koʻp xotira (peak RSS, MiB): butun faylni yuklash va oqim.
# Fayllar /tmp ga yoziladi (gitga emas): eng kattasi 1 GiB, oxirida oʻchiriladi.
set -euo pipefail
cd "$(dirname "$0")"
dir=${TMPDIR:-/tmp}/io19-olch
mkdir -p "$dir/bin"
go build -o "$dir/bin/" ./yarat ./hammasi ./oqim
rss() { "$@" | grep -oE '[0-9.]+MiB' | tr -d 'MiB'; }
for mb in 16 64 256 1024; do
  f=$dir/b$mb.csv
  "$dir/bin/yarat" "$mb" "$f" > /dev/null
  echo "hammasi $mb $(rss "$dir/bin/hammasi" "$f")"
  echo "oqim $mb $(rss "$dir/bin/oqim" "$f")"
  rm -f "$f"
done
rm -rf "$dir"
