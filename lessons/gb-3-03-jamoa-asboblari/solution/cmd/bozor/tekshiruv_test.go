package main

// Bu fayl vazifaning 3-qadamini tekshiradi (TASK.md): dastur qaysi golang.org/x/text versiyasi bilan
// yigʻilganini oʻqiydi. Uni oʻzgartirmang.

import (
	"runtime/debug"
	"strconv"
	"strings"
	"testing"
)

// minClean — 2026-11-27 holatida govulncheck x/text uchun hech qanday zaiflik koʻrsatmaydigan eng past versiya
// (GO-2022-1059 v0.3.8'da, GO-2026-5970 v0.39.0'da, GO-2026-6629 v0.41.0'da tuzatilgan).
const minClean = "v0.41.0"

func TestVazifaXText(t *testing.T) {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		t.Fatal("build maʼlumoti oʻqilmadi")
	}
	var version string
	for _, m := range info.Deps {
		if m.Path == "golang.org/x/text" {
			version = m.Version
			if m.Replace != nil {
				version = m.Replace.Version
			}
		}
	}
	if version == "" {
		t.Fatal("dastur golang.org/x/text'siz yigʻilgan; go.mod'ni tekshiring")
	}
	if less(version, minClean) {
		t.Errorf("go.mod: golang.org/x/text %s — govulncheck bu versiyada zaiflik topadi (main language.MatchStrings'ni chaqiradi); kamida %s kerak. govulncheck ./... nima deydi? (TASK.md, 3-qadam)", version, minClean)
	}
}

// less a < b ekanini aytadi; ikkalasi vX.Y.Z koʻrinishida. vX.Y.Z-… (pre-release yoki pseudo-version)
// vX.Y.Z'dan kichik hisoblanadi.
func less(a, b string) bool {
	pa, preA := parts(a)
	pb, preB := parts(b)
	for i := range 3 {
		if pa[i] != pb[i] {
			return pa[i] < pb[i]
		}
	}
	return preA && !preB
}

func parts(v string) (out [3]int, pre bool) {
	v = strings.TrimPrefix(v, "v")
	if i := strings.IndexByte(v, '+'); i >= 0 {
		v = v[:i]
	}
	if i := strings.IndexByte(v, '-'); i >= 0 {
		v, pre = v[:i], true
	}
	for i, s := range strings.SplitN(v, ".", 3) {
		n, err := strconv.Atoi(s)
		if err != nil {
			return out, pre
		}
		out[i] = n
	}
	return out, pre
}
