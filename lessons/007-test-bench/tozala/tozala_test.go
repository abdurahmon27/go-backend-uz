package tozala

import "testing"

func TestTozala(t *testing.T) {
	jadval := []struct {
		nom, kirish, want string
	}{
		{"klaviatura", "O'zbekiston", "Oʻzbekiston"},
		{"jingalak", "tog’", "togʻ"},
		{"tutuq", "ma'no", "maʼno"},
		{"qoshtirnoq", "U ‘Salom’ dedi", "U ‘Salom’ dedi"},
	}
	for _, tc := range jadval {
		t.Run(tc.nom, func(t *testing.T) {
			if got := tozala(tc.kirish); got != tc.want {
				t.Errorf("tozala(%+q) = %+q, kutilgan %+q", tc.kirish, got, tc.want)
			}
		})
	}
}
