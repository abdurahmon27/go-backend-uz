package main

// Videodagi importStocks'ning testi: starterda ham oʻtadi.

import (
	"context"
	"testing"
)

func TestImportStocks(t *testing.T) {
	resetStats()
	stocks, err := importStocks(context.Background(), catalogIDs(100))
	if err != nil {
		t.Fatalf("importStocks xato qaytardi: %v", err)
	}
	if stats.peak > workers {
		t.Errorf("bir vaqtda %d ta soʻrov ketdi, koʻpi bilan %d boʻlishi kerak", stats.peak, workers)
	}
	if stocks[6] != 7 { // ID 7 → 7 % 50 = 7 dona
		t.Errorf("stocks[6] = %d, kutilgan 7 (natija ids tartibida boʻlishi kerak)", stocks[6])
	}
}
