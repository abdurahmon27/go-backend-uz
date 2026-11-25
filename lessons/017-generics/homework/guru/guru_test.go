package guru

import "testing"

func TestJami(t *testing.T) {
	cases := []struct {
		nom  string
		got  int64
		want int64
	}{
		{"*Buyurtma", JamiGeneric([]*Buyurtma{{100}, {200}}), 300},
		{"*Mijoz", JamiGeneric([]*Mijoz{{50}}), 50},
		{"Narx", JamiGeneric([]Narx{10, 20}), 30},
		{"Chegirma", JamiGeneric([]Chegirma{5}), -5},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("%s: %d, kutilgan %d", c.nom, c.got, c.want)
		}
	}
}
