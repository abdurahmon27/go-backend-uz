package main

import "fmt"

// chegirma narxdan foiz ulushini ayiradi. foiz 0 dan 100 gacha bo'lishi kerak.
func chegirma(narx, foiz int) (int, error) {
	if foiz < 0 || foiz > 100 {
		return 0, fmt.Errorf("foiz 0..100 oralig'ida bo'lishi kerak, berildi: %d", foiz)
	}
	return narx - narx*foiz/100, nil
}
