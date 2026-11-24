package tolov

import "testing"

func TestTolash(t *testing.T) {
	if err := Tolash(0); err != ErrRad {
		t.Fatalf("kutilgan ErrRad, keldi %v", err)
	}
	// TODO: musbat summa uchun test
}
