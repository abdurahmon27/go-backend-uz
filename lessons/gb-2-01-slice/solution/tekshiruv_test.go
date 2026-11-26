package main

// Bu fayl vazifani (TASK.md) tekshiradi: go test uni stock_test.go bilan birga ishga tushiradi.
// Uni oʻzgartirmang.

import (
	"slices"
	"testing"
)

// Aksiya narxlari toʻgʻri: 10 % arzonroq, tiyinda.
func TestVazifaNarx(t *testing.T) {
	prices := []int64{1_200_000, 500_000, 2_500_000}
	sale := discounted(prices, 10)
	want := []int64{1_080_000, 450_000, 2_250_000}
	// slices.Equal: uzunlik va elementlar bir xilmi (slice'larni == bilan solishtirib boʻlmaydi).
	if !slices.Equal(sale, want) {
		t.Errorf("discounted(…, 10) = %v; kutilgan %v", sale, want)
	}
}

// Katalog narxlari oʻzgarmaydi: discounted prices massiviga yozmasligi kerak.
func TestVazifaKatalog(t *testing.T) {
	prices := []int64{1_200_000, 500_000, 2_500_000}
	discounted(prices, 10)
	if want := []int64{1_200_000, 500_000, 2_500_000}; !slices.Equal(prices, want) {
		t.Errorf("discounted katalogni oʻzgartirdi: prices endi %v, %v boʻlib qolishi kerak (TASK.md)", prices, want)
	}
}

// Natija alohida massivda: aksiya narxini keyin oʻzgartirish ham katalogga tegmaydi.
func TestVazifaAlohidaMassiv(t *testing.T) {
	prices := []int64{1_200_000, 500_000, 2_500_000}
	sale := discounted(prices, 10)
	if len(sale) == 0 {
		t.Fatalf("discounted boʻsh slice qaytardi, %d ta narx kutilgan", len(prices))
	}
	sale[0] = 1
	if prices[0] != 1_200_000 {
		t.Errorf("sale[0] = 1 yozilganda prices[0] ham %d boʻldi: natija katalog bilan bitta massivda (TASK.md, make)", prices[0])
	}
}
