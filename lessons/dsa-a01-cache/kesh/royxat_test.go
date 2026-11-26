package kesh

import (
	"testing"
	"unsafe"
)

// new bilan ketma-ket yaratilgan tugunlar xotirada yonma-yon turadimi?
// go test -run Yonma -v ./kesh
func TestYonmaYon(t *testing.T) {
	n := 1 << 20
	yonma := 0
	for p := RoyxatYangi(n); p.Keyingi != nil; p = p.Keyingi {
		if uintptr(unsafe.Pointer(p.Keyingi))-uintptr(unsafe.Pointer(p)) == unsafe.Sizeof(Tugun{}) {
			yonma++
		}
	}
	foiz := 100 * float64(yonma) / float64(n-1)
	t.Logf("%d tugun: %.1f%% da keyingisi roppa-rosa +%d bayt", n, foiz, unsafe.Sizeof(Tugun{}))
	if foiz < 90 {
		t.Errorf("kutilganidan kam: %.2f%%", foiz)
	}
}
