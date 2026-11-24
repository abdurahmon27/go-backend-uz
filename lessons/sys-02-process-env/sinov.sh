#!/usr/bin/env bash
# ./sinov.sh TERM | KILL        (KUT=0 ./sinov.sh TERM: signal ushlanmaydi)
# Serverga 1 soniyalik soʻrov yuboramiz va soʻrov oʻrtasida signal joʻnatamiz.
exec 2>&1                                  # xatolar ham shu oqimda, tartib saqlansin
CGO_ENABLED=0 go build -o srv ./server     # C kompilyatorisiz, sof Go
./srv &
srv=$!
sleep 0.3
curl -sS localhost:8077 &
sleep 0.3
kill -s "$1" $srv
wait $srv 2>/dev/null                      # bash'ning "Killed" xabarisiz
echo "server chiqish kodi: $?"
wait
