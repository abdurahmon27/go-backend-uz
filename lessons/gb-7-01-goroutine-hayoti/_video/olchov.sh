#!/usr/bin/env bash
# Videodagi ikki oʻlchov (research.md F3, F7):
#   ./_video/olchov.sh go      — 02-go (kutmasdan go) 200 marta: nechta «boshlandi» va «tayyor» qatori chiqdi
#   ./_video/olchov.sh append  — 05-append (-race'siz) 200 marta: natijada nechta son qoldi (3 = hech biri yoʻqolmadi)
# Natija har safar biroz farq qiladi: goroutine'lar tartibi kafolatlanmaydi.
set -u
export GOTOOLCHAIN=${GOTOOLCHAIN:-go1.27.4} CGO_ENABLED=0
here=$(cd "$(dirname "$0")/.." && pwd)
case ${1:?go yoki append} in
go) step=02-go ;;
append) step=05-append ;;
*) echo "go yoki append" >&2; exit 2 ;;
esac
t=$(mktemp -d)
cp "$here/go.mod" "$here/_video/$step/main.go" "$t/"
cd "$t" && go build -o prog . || exit 1
for _ in $(seq 200); do
	./prog >out.txt
	if [ "$1" = go ]; then
		echo "boshlandi=$(grep -c boshlandi out.txt) tayyor=$(grep -c tayyor out.txt)"
	else
		grep mahsulotlar out.txt | awk -F'[][]' '{print "sonlar=" split($2, a, " ")}'
	fi
done | sort | uniq -c
