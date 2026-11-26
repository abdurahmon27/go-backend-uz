#!/usr/bin/env bash
# poyga.sh [usul] [necha marta] — ikki jarayon (A va B) bitta JSON faylga 1000 tadan buyurtma yozadi.
# Har urinish toza fayldan boshlanadi. Oxirida: fayl necha marta buzildi, butun qolganida nechta buyurtma yoʻqoldi.
set -u
cd "$(dirname "$0")"
usul=${1:-oddiy}; marta=${2:-1}
go build -o bin/fayl ./fayl || exit 1
buzildi=0; butun=0; yoqoldi=0; eng_kam=2000; eng_kop=0; buz_kam=2000; buz_kop=0
for ((k = 1; k <= marta; k++)); do
  rm -f buyurtmalar.json buyurtmalar.json.*
  if ((marta == 1)); then
    ./bin/fayl -usul "$usul" -mijoz A & ./bin/fayl -usul "$usul" -mijoz B & wait
    ./bin/fayl -sana; exit
  fi
  ./bin/fayl -usul "$usul" -mijoz A >/dev/null & ./bin/fayl -usul "$usul" -mijoz B >/dev/null & wait
  if natija=$(./bin/fayl -sana); then
    jami=$(sed -E 's/faylda: ([0-9]+) .*/\1/' <<<"$natija")
    ((butun++, yoqoldi += 2000 - jami))
    ((2000 - jami < eng_kam)) && eng_kam=$((2000 - jami))
    ((2000 - jami > eng_kop)) && eng_kop=$((2000 - jami))
  else
    ((buzildi++))
    b=$(sed -nE 's/^boshida ([0-9]+) ta.*/\1/p' <<<"$natija"); b=${b:-0}
    ((b < buz_kam)) && buz_kam=$b
    ((b > buz_kop)) && buz_kop=$b
  fi
done
echo "$marta urinish ($usul): fayl buzildi $buzildi marta, butun qoldi $butun marta"
((buzildi > 0)) && echo "buzilganda faylning butun boshida $buz_kam–$buz_kop ta buyurtma bor edi (2000 dan)"
((butun > 0)) && echo "butun qolganida yoʻqolgan buyurtma: oʻrtacha $((yoqoldi / butun)), eng kami $eng_kam, eng koʻpi $eng_kop (2000 dan)"
rm -f buyurtmalar.json buyurtmalar.json.*
exit 0
