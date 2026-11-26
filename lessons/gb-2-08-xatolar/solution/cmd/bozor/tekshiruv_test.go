package main

// Bu fayl vazifani (TASK.md, 2-qadam) tekshiradi: go test ./... uni main.go bilan birga ishga tushiradi.
// Uni oʻzgartirmang.

import (
	"errors"
	"fmt"
	"testing"

	"github.com/abdurahmon27/bozor/internal/store"
)

// 2-qadam: message xatoni turiga qarab foydalanuvchi matniga aylantiradi.
// Baʼzi xatolar yana bir qavat oʻralgan («buyurtma 7: …»): == va oddiy type assertion
// ularni topmaydi, errors.Is va errors.As topadi.
func TestVazifaMessage(t *testing.T) {
	var st store.Store
	st.Add("olma", 1_200_000, 40)

	tests := []struct {
		name string
		err  error
		want string
		hint string
	}{
		{"topilmadi", st.Sell(9, 1), "bunday mahsulot yoʻq",
			"errors.Is(err, store.ErrNotFound)"},
		{"topilmadi_ikki_qavat", fmt.Errorf("buyurtma 7: %w", st.Sell(9, 1)), "bunday mahsulot yoʻq",
			"errors.Is butun zanjirni tekshiradi"},
		{"qoldiq", st.Sell(1, 41), "olma: faqat 40 dona qoldi",
			"errors.As(err, &se), se — *store.StockError; matn se.Name va se.Stock'dan"},
		{"qoldiq_ikki_qavat", fmt.Errorf("buyurtma 7: %w", st.Sell(1, 41)), "olma: faqat 40 dona qoldi",
			"errors.As ham butun zanjirni tekshiradi"},
		{"dona_soni", fmt.Errorf("buyurtma 7: %w", store.ErrBadQuantity), "dona soni notoʻgʻri",
			"errors.Is(err, store.ErrBadQuantity)"},
		{"boshqa", errors.New("disk toʻldi"), "xato: disk toʻldi",
			"qolgan xatolar oʻzgarmaydi: «xato: » + err.Error()"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := message(tc.err)
			if got != tc.want {
				t.Errorf("message(%q) = %q; kutilgan %q — %s (TASK.md, 2-qadam)", tc.err.Error(), got, tc.want, tc.hint)
			}
		})
	}
}
