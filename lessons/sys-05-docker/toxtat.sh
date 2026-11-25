#!/usr/bin/env bash
# ./toxtat.sh IMAGE [buyruq]
# Container'ni ishga tushiradi, docker stop qancha kutganini oʻlchaydi va chiqish kodini aytadi.
exec 2>&1
docker rm -f toxtat >/dev/null 2>&1
docker run -d --name toxtat "$@" >/dev/null || exit 1
sleep 0.5
bosh=$(date +%s%N)
docker stop toxtat >/dev/null              # SIGTERM, 10 soniyadan keyin SIGKILL
tugash=$(date +%s%N)
docker logs toxtat
ms=$(( (tugash - bosh) / 1000000 ))
kod=$(docker inspect -f '{{.State.ExitCode}}' toxtat)
printf 'docker stop: %d.%d s, chiqish kodi: %s\n' $((ms / 1000)) $((ms % 1000 / 100)) "$kod"
docker rm toxtat >/dev/null
