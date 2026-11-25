#!/usr/bin/env bash
# bash homework/javoblar.sh   (dars papkasidan; avval ./yigish.sh)
# Uy vazifasi javoblarini ishga tushiradi. 9000 va 8078 portlari boʻsh boʻlishi kerak.
set -u
cd "$(dirname "$0")/.." || exit 1
exec 2>&1
tozala() { docker rm -f h1 h2 h3 h4 >/dev/null 2>&1; }
tozala

echo "== Oson: PORT=9000 muhitdan, -p 9000:9000"
docker run -d --name h1 -e PORT=9000 -p 9000:9000 mehmon:kichik >/dev/null
sleep 0.5
docker logs h1
curl -s localhost:9000

echo "== Oʻrta: volume oʻrniga papka (bind mount), fayl egasi siz"
rm -rf homework/papka && mkdir homework/papka
docker run -d --name h2 -p 8078:8077 --user "$(id -u):$(id -g)" \
  -v "$PWD/homework/papka:/data" mehmon:kichik >/dev/null
sleep 0.5
curl -s localhost:8078
curl -s localhost:8078
docker rm -f h2 >/dev/null
echo "xostdagi fayl: homework/papka/mehmonlar.txt, qatorlar: $(wc -l < homework/papka/mehmonlar.txt)"
rm -rf homework/papka

echo "== Qiyin: exec'siz image + --init"
./toxtat.sh --init mehmon-tuzoq

echo "== Guru: USER 65532 (root emas) + volume"
docker volume rm -f guru-daftar >/dev/null
echo "-- oddiy image'ga --user 65532 bersak:"
docker run -d --name h3 -p 8077:8077 --user 65532:65532 -v guru-daftar:/data mehmon:kichik >/dev/null
sleep 0.5
curl -s localhost:8077
docker rm -f h3 >/dev/null
docker volume rm -f guru-daftar >/dev/null
echo "-- homework/guru/Dockerfile bilan:"
docker build -q -t mehmon-boshqa:guru -f homework/guru/Dockerfile . >/dev/null
for i in 1 2; do
  docker run -d --name h4 -p 8077:8077 -v guru-daftar:/data mehmon-boshqa:guru >/dev/null
  sleep 0.5
  curl -s localhost:8077
  docker rm -f h4 >/dev/null
done
docker run -d --name h4 mehmon-boshqa:guru >/dev/null
sleep 0.3
docker top h4 -o user,pid,comm
docker rm -f h4 >/dev/null
docker volume rm guru-daftar >/dev/null
tozala
