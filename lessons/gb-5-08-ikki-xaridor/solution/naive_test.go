// Videodagi birinchi (xato) versiyani koʻrsatadi: poyga haqiqatan boʻlishini isbotlaydi. Uni oʻzgartirmang.
package bozor

import "testing"

// buyNaive bilan 20 xaridor, omborda 5 ta: ortiqcha sotuv ("olindi" 5 tadan koʻp) kuzatilsa test oʻtadi.
// Tasodifan poyga boʻlmay qolsa 5 marta qayta uriniladi. Natijani koʻrish: go test -run Naive -v
func TestNaiveOversells(t *testing.T) {
	s := newStore(t)
	ctx := testCtx(t)
	best := 0
	for round := 1; round <= 5 && best <= 5; round++ {
		setStock(t, s, 4, 5)
		ok := runBuyers(t, ctx, s.buyNaive, 20)
		t.Logf("%d-urinish: 20 xaridor, omborda 5 ta: %d tasiga \"olindi\" dedik", round, ok)
		best = max(best, ok)
	}
	if best <= 5 {
		t.Errorf("poyga kuzatilmadi (5 urinishda ham 5 tadan koʻp sotilmadi): boshqa kompyuterda kamdan-kam; testni qayta ishga tushiring")
	}
}
