#!/usr/bin/env bash
# qulash-baza.sh [necha marta] — A va B toʻxtovsiz INSERT qiladi, biz Postgres'ni kill -9 qilamiz
# (docker compose kill = SIGKILL). Keyin bazani qayta koʻtarib, tasdiq olgan har bir id joyidami, tekshiramiz.
set -u
cd "$(dirname "$0")"
set -a; . ./.env; set +a
marta=${1:-1}
CGO_ENABLED=0 go build -o bin/baza ./baza || exit 1
docker compose up -d --wait >/dev/null 2>&1
butun=0; ortiqcha=0; yoqolgan=0
for ((k = 1; k <= marta; k++)); do
  ./bin/baza -tayyorla || exit 1
  ./bin/baza -mijoz A -n 0 -tasdiq bin/A.ids >bin/A.log & ./bin/baza -mijoz B -n 0 -tasdiq bin/B.ids >bin/B.log &
  sleep "1.$((RANDOM % 9))"
  docker compose kill db >/dev/null 2>&1
  wait
  docker compose up -d --wait >/dev/null 2>&1
  natija=$(./bin/baza -tekshir bin/A.ids bin/B.ids) || { echo "$natija"; exit 1; }
  read -r tasdiq bazada yoq < <(sed -E 's/[^0-9]+/ /g' <<<"$natija")
  if ((marta == 1)); then
    cut -c1-90 bin/A.log bin/B.log
    docker compose logs db --no-log-prefix | grep -E 'not properly shut down|redo (starts|done)|ready to accept' | tail -4 | sed -E 's/^.*LOG:  //; s/ at [0-9A-F]+\/[0-9A-F]+.*//'
    echo "$natija"
    exit
  fi
  ((yoq == 0)) && ((butun++))
  ((yoqolgan += yoq, ortiqcha += bazada - tasdiq + yoq))
  echo "$k: $natija"
done
echo "$marta marta kill -9: $butun marta tasdiqlangan har bir id joyida (yoʻqolgan jami $yoqolgan), tasdiqsiz ortiqcha qator jami $ortiqcha"
