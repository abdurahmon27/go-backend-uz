package store

// Bu fayl vazifaning 1- va 2-qadamini (TASK.md) tekshiradi: go test ./... uni store.go bilan birga
// ishga tushiradi. Uni oʻzgartirmang.

import (
	"errors"
	"os"
	"regexp"
	"strings"
	"testing"
)

// jadvalQatori — bitta qatorda yozilgan jadval qatori: {"oddiy", 2, 10, 13, nil}, yoki {name: "oddiy", …},
var jadvalQatori = regexp.MustCompile(`(?m)^\s*\{.+\},\s*(//.*)?$`)

// 1-qadam: Restock(id, quantity) qoldiqqa quantity dona qoʻshadi; xato boʻlsa katalog oʻzgarmaydi.
func TestVazifaRestock(t *testing.T) {
	tests := []struct {
		name      string
		id        int64
		quantity  int
		wantStock int   // shu ID'li mahsulotning katalogdagi qoldigʻi Restock'dan keyin
		wantErr   error // errors.Is bilan tekshiriladi; nil — xato kutilmaydi
	}{
		{"oddiy", 2, 10, 13, nil},
		{"bitta", 4, 1, 3, nil},
		{"chegara", 1, 10_000, 10_040, nil},
		{"yoq_mahsulot", 9, 5, 0, ErrNotFound},
		{"nol_dona", 2, 0, 3, ErrBadQuantity},
		{"manfiy_dona", 2, -5, 3, ErrBadQuantity},
		{"juda_kop", 2, 10_001, 3, ErrBadQuantity},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			st := katalog()
			p, err := st.Restock(tc.id, tc.quantity)
			if !errors.Is(err, tc.wantErr) {
				want := "xatosiz"
				switch tc.wantErr {
				case ErrNotFound:
					want = "ErrNotFound (%w bilan oʻralgan)"
				case ErrBadQuantity:
					want = "ErrBadQuantity (%w bilan oʻralgan; dona 1..10 000)"
				}
				t.Fatalf("Restock(%d, %d) xatosi «%v», kutilgan: %s (TASK.md, 1-qadam)", tc.id, tc.quantity, err, want)
			}
			if got := st.find(tc.id); got != nil && got.Stock != tc.wantStock {
				t.Fatalf("Restock(%d, %d) dan keyin katalogda qoldiq %d, kutilgan %d (TASK.md, 1-qadam)", tc.id, tc.quantity, got.Stock, tc.wantStock)
			}
			if tc.wantErr == nil && (p.ID != tc.id || p.Stock != tc.wantStock) {
				t.Fatalf("Restock(%d, %d) %+v qaytardi, kutilgan ID %d, qoldiq %d: yangilangan mahsulotni qaytaring (TASK.md, 1-qadam)", tc.id, tc.quantity, p, tc.id, tc.wantStock)
			}
		})
	}
}

// 2-qadam: restock_test.go'da Restock uchun oʻz jadvalli testingiz bor (GB-1.08 shakli).
func TestVazifaOzTestingiz(t *testing.T) {
	// os.ReadFile go test ishlayotgan papkadagi (internal/store) faylni oʻqiydi.
	kod, err := os.ReadFile("restock_test.go")
	if err != nil {
		t.Fatalf("internal/store/restock_test.go oʻqilmadi (%v): Restock uchun jadvalli test yozing (TASK.md, 2-qadam)", err)
	}
	for _, kerak := range []string{"func TestRestock(", "t.Run(", "Restock("} {
		if !strings.Contains(string(kod), kerak) {
			t.Errorf("restock_test.go ichida %q topilmadi: test TestRestock deb nomlansin, jadval qatorlarini t.Run bilan ishga tushirsin va Restock'ni chaqirsin (TASK.md, 2-qadam)", kerak)
		}
	}
	// Jadval qatori: bitta qatorda yozilgan {…}, (maydon nomlari bilan yoki ularsiz).
	if n := len(jadvalQatori.FindAllString(string(kod), -1)); n < 3 {
		t.Errorf("restock_test.go jadvalida %d ta qator topildi, kamida 3 ta kerak, har biri bitta qatorda {…}, (oddiy, yoʻq ID, notoʻgʻri dona; TASK.md, 2-qadam)", n)
	}
}
