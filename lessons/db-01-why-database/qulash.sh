#!/usr/bin/env bash
# qulash.sh [usul] [necha marta] — bitta jarayon toʻxtovsiz yozadi, biz uni tasodifiy paytda kill -9 qilamiz
# (crash: jarayon hech narsani tugatishga ulgurmaydi). Keyin fayl butunmi, shuni tekshiramiz.
set -u
cd "$(dirname "$0")"
usul=${1:-oddiy}; marta=${2:-100}
go build -o bin/fayl ./fayl || exit 1
bosh=0; yarim=0; butun=0
for ((k = 1; k <= marta; k++)); do
  rm -f buyurtmalar.json buyurtmalar.json.*
  ./bin/fayl -usul "$usul" -n 0 >/dev/null &
  sleep "0.$((RANDOM % 9 + 1))"                     # 0.1–0.9 s yozsin
  kill -9 $! ; wait $! 2>/dev/null
  if [[ ! -e buyurtmalar.json ]]; then ((k--)); continue; fi   # birinchi yozuvgacha oʻldi: hisobga olinmaydi
  if ./bin/fayl -sana >/dev/null; then ((butun++))
  elif [[ ! -s buyurtmalar.json ]]; then ((bosh++))  # 0 bayt: hamma buyurtma yoʻq
  else ((yarim++)); fi
done
echo "$marta marta kill -9 ($usul): fayl butun $butun, boʻm-boʻsh $bosh, yarim yozilgan $yarim"
rm -f buyurtmalar.json buyurtmalar.json.*
