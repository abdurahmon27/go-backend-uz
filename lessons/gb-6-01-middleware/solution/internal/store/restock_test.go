package store

// 2-qadam javobi (TASK.md): Restock uchun oʻz jadvalli testimiz (GB-1.08 shakli).

import (
	"errors"
	"testing"
)

func TestRestock(t *testing.T) {
	tests := []struct {
		name      string
		id        int64
		quantity  int
		wantStock int
		wantErr   error
	}{
		{"oddiy", 2, 10, 13, nil},
		{"yoq_id", 9, 1, 0, ErrNotFound},
		{"nol_dona", 2, 0, 3, ErrBadQuantity},
		{"juda_kop", 2, 10_001, 3, ErrBadQuantity},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			st := katalog()
			_, err := st.Restock(tc.id, tc.quantity)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("Restock(%d, %d) xatosi %v, kutilgan %v", tc.id, tc.quantity, err, tc.wantErr)
			}
			if p := st.find(tc.id); p != nil && p.Stock != tc.wantStock {
				t.Fatalf("qoldiq %d, kutilgan %d", p.Stock, tc.wantStock)
			}
		})
	}
}
