package kabisa

import "testing"

func TestJadval(t *testing.T) {
	jadval := []struct {
		nom  string
		yil  int
		want bool
	}{
		{"oddiy", 2026, false},
		{"4 ga", 2024, true},
		{"100 ga", 1900, false},
		{"400 ga", 2000, true},
	}
	for _, tc := range jadval {
		t.Run(tc.nom, func(t *testing.T) {
			if got := kabisa(tc.yil); got != tc.want {
				t.Errorf("kabisa(%d) = %v, kutilgan %v", tc.yil, got, tc.want)
			}
		})
	}
}
