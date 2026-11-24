package kabisa

import "testing"

func TestKabisa(t *testing.T) {
	if !kabisa(2024) {
		t.Errorf("kabisa(2024) = false, kutilgan true")
	}
}
