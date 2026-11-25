package guru

import "testing"

type Buyurtma struct{ Summa int64 }

func TestBosh(t *testing.T) {
	var p *Buyurtma
	var m map[string]int
	var s []int
	var f func()
	jadval := []struct {
		nom  string
		a    any
		want bool
	}{
		{"nil", nil, true},
		{"nil pointer", p, true},
		{"nil map", m, true},
		{"nil slice", s, true},
		{"nil func", f, true},
		{"pointer", &Buyurtma{}, false},
		{"0", 0, false},
		{"bosh string", "", false},
	}
	for _, c := range jadval {
		if got := bosh(c.a); got != c.want {
			t.Errorf("%s: bosh = %v, kutilgan %v", c.nom, got, c.want)
		}
		if c.nom == "nil pointer" && c.a == nil {
			t.Error("typed nil a == nil boʻlmasligi kerak")
		}
	}
}
