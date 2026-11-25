package qiyin

import (
	"slices"
	"testing"
)

type Buyurtma struct {
	Summa   int64
	Tolandi bool
}

func TestGuruhla(t *testing.T) {
	r := Royxat[Buyurtma]{{150000, true}, {99000, false}, {42000, true}}
	g := r.Guruhla(func(b Buyurtma) bool { return b.Tolandi })
	if len(g[true]) != 2 || len(g[false]) != 1 {
		t.Fatalf("guruhlar: %v", g)
	}
	// kalit turi boshqa: string
	h := r.Guruhla(func(b Buyurtma) string {
		if b.Summa >= 100000 {
			return "katta"
		}
		return "kichik"
	})
	want := Royxat[Buyurtma]{{99000, false}, {42000, true}}
	if !slices.Equal(h["kichik"], want) {
		t.Fatalf("kichik: %v", h["kichik"])
	}
}
