package store

// Bu fayl vazifani (TASK.md) tekshiradi: go test ./... uni store.go bilan birga ishga tushiradi.
// Uni oʻzgartirmang. Test store paketining ichida, shuning uchun products maydonini ham koʻradi.

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

// katalog ikki mahsulotli yangi katalog qaytaradi: olma (ID 1, qoldiq 40) va non (ID 2, qoldiq 100).
func katalog() *Store {
	var st Store
	st.Add("olma", 1_200_000, 40)
	st.Add("non", 400_000, 100)
	return &st
}

// Videodagi xatolar (bu test boshidanoq oʻtadi).
func TestVideodagiXatolar(t *testing.T) {
	st := katalog()

	err := st.Sell(9, 1)
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("Sell(9, 1) = %v; errors.Is(err, ErrNotFound) true boʻlishi kerak", err)
	}
	if err != nil && err.Error() != "mahsulot 9: mahsulot topilmadi" {
		t.Errorf("Sell(9, 1) matni %q; kutilgan %q", err.Error(), "mahsulot 9: mahsulot topilmadi")
	}

	err = st.Sell(1, 41)
	var se *StockError
	if !errors.As(err, &se) {
		t.Fatalf("Sell(1, 41) = %v; *StockError boʻlishi kerak", err)
	}
	if se.Name != "olma" || se.Stock != 40 || se.Want != 41 {
		t.Errorf("StockError = %+v; kutilgan {Name:olma Stock:40 Want:41}", *se)
	}

	if err := st.Sell(2, 3); err != nil {
		t.Fatalf("Sell(2, 3) = %v; kutilgan nil", err)
	}
	if st.products[0].Stock != 40 || st.products[1].Stock != 97 {
		t.Errorf("qoldiqlar %d va %d; kutilgan 40 va 97", st.products[0].Stock, st.products[1].Stock)
	}
}

// 1-qadam: musbat boʻlmagan dona soni — ErrBadQuantity, %w bilan oʻralgan.
func TestVazifaDonaSoni(t *testing.T) {
	for _, quantity := range []int{0, -2} {
		st := katalog()
		err := st.Sell(1, quantity)
		if err == nil {
			t.Fatalf("Sell(1, %d) = nil; kutilgan error (TASK.md, 1-qadam)", quantity)
		}
		// Shablondagi %%w — matnda foiz belgisining oʻzi, ya'ni «%w» chiqadi.
		if !errors.Is(err, ErrBadQuantity) {
			t.Errorf("Sell(1, %d) xatosi %q, lekin errors.Is(err, ErrBadQuantity) = false: ErrBadQuantity'ni fmt.Errorf'da %%w bilan oʻrab qaytaring (TASK.md, 1-qadam)",
				quantity, err.Error())
		}
		// fmt.Sprint(x) — x'ni matnga aylantiradi; strings.Contains(a, b) — a matnida b bormi.
		if !strings.Contains(err.Error(), fmt.Sprint(quantity)) {
			t.Errorf("Sell(1, %d) xatosi %q: matnda soʻralgan son (%d) boʻlsin (TASK.md, 1-qadam)", quantity, err.Error(), quantity)
		}
		if st.products[0].Stock != 40 {
			t.Errorf("xato boʻlgan Sell(1, %d)'dan keyin olma qoldigʻi %d; kutilgan 40 (TASK.md, 1-qadam)", quantity, st.products[0].Stock)
		}
	}
}
