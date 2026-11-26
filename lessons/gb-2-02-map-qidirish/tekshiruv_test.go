package main

// Bu fayl vazifani (TASK.md) tekshiradi: go test uni main.go bilan birga ishga tushiradi.
// Uni oʻzgartirmang.

import (
	"fmt"
	"slices"
	"testing"
)

// ombor har safar yangi map qaytaradi: testlar bir-birining qoldigʻini buzmasin.
func ombor() map[string]int {
	return map[string]int{"olma": 5, "non": 2, "choy": 1}
}

// tekshirQoldiq stock'dagi har bir mahsulot qoldigʻini want bilan solishtiradi.
func tekshirQoldiq(t *testing.T, stock, want map[string]int, qadam string) {
	t.Helper() // xato qatori shu yerni emas, chaqirgan joyni koʻrsatsin
	var names []string
	for name := range want {
		names = append(names, name)
	}
	slices.Sort(names)
	for _, name := range names {
		got, ok := stock[name]
		if !ok {
			t.Errorf("reserve'dan keyin %s ombordan yoʻqolib qoldi: kalitni oʻchirmang (%s)", name, qadam)
			continue
		}
		if got != want[name] {
			t.Errorf("reserve'dan keyin %s qoldigʻi %d boʻlishi kerak, hozir %d (%s)", name, want[name], got, qadam)
		}
	}
	if len(stock) != len(want) {
		t.Errorf("reserve'dan keyin omborda %d xil mahsulot boʻlishi kerak, hozir %d (%s)", len(want), len(stock), qadam)
	}
}

// 1-qadam: hammasi yetarli boʻlsa, xato yoʻq va har bir mahsulot sotilgan donaga kamayadi.
func TestVazifaAyiradi(t *testing.T) {
	stock := ombor()
	err := reserve(stock, []string{"olma", "non", "olma"})
	if err != nil {
		t.Fatalf("reserve(olma non olma) xatosiz boʻlishi kerak, hozir xato: %v (TASK.md, 1-qadam)", err)
	}
	tekshirQoldiq(t, stock, map[string]int{"olma": 3, "non": 1, "choy": 1}, "TASK.md, 1-qadam")
}

// 2-qadam: omborda yoʻq mahsulot xato beradi, ombor oʻzgarmaydi.
// 3-qadam: yetmaydigan mahsulot xato beradi, ombor oʻzgarmaydi.
func TestVazifaRadEtadi(t *testing.T) {
	tests := []struct {
		name    string
		order   []string
		wantErr string
		qadam   string
	}{
		{"yoq", []string{"olma", "anor"}, "anor omborda yoʻq", "TASK.md, 2-qadam"},
		{"yetmaydi", []string{"choy", "choy"}, "choy: omborda 1 dona, buyurtmada 2", "TASK.md, 3-qadam"},
		{"yarmi", []string{"olma", "non", "non", "non"}, "non: omborda 2 dona, buyurtmada 3", "TASK.md, 3-qadam"},
		{"birinchisi", []string{"anor", "choy", "choy"}, "anor omborda yoʻq", "TASK.md, 3-qadam: birinchi muammoli nom"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			stock := ombor()
			err := reserve(stock, tc.order)
			if err == nil {
				t.Fatalf("reserve(%v) xato qaytarishi kerak: %q; hozir xato yoʻq (%s)", tc.order, tc.wantErr, tc.qadam)
			}
			if got := fmt.Sprintf("%v", err); got != tc.wantErr {
				t.Errorf("reserve(%v) xatosi %q boʻlishi kerak, hozir %q (%s)", tc.order, tc.wantErr, got, tc.qadam)
			}
			tekshirQoldiq(t, stock, ombor(), tc.qadam+": xato boʻlsa ombor oʻzgarmaydi")
		})
	}
}
