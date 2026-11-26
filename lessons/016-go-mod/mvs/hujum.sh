#!/usr/bin/env bash
# Hujum mashqi: zaharli/ proxy'da example.com/log v1.3.0 ning kodi almashtirilgan,
# versiya raqami esa o'sha. dokon'ni toza module cache bilan (yangi kompyuterdagidek)
# build qilamiz: avval zaharli/, topilmasa proxy/ dan yuklanadi.
# go.sum'dagi hash mos kelmaydi va go build to'xtaydi.
set -u
cd "$(dirname "$0")/dokon"
kesh=$(mktemp -d)
GOFLAGS=-modcacherw GOMODCACHE="$kesh" GONOSUMDB=example.com \
	GOPROXY="file://$(cd .. && pwd)/zaharli,file://$(cd .. && pwd)/proxy" \
	go build -o /dev/null .
kod=$?
rm -rf "$kesh"
exit $kod
