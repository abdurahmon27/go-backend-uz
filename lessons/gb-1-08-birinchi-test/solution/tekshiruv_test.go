package main

// Bu fayl vazifani (TASK.md) tekshiradi: go test uni price_test.go va sizning
// sold_test.go faylingiz bilan birga ishga tushiradi. Uni oʻzgartirmang.

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

// 1-qadam: sold_test.go bor va unda TestParseSold jadvali t.Run bilan yozilgan.
func TestVazifaJadvalBor(t *testing.T) {
	// os.ReadFile faylni oʻqib, uning baytlarini va xatoni qaytaradi.
	kod, err := os.ReadFile("sold_test.go")
	if err != nil {
		t.Fatalf("sold_test.go oʻqilmadi (%v): avval parseSold uchun jadvalli test yozing (TASK.md, 1-qadam)", err)
	}
	// strings.Contains matn ichida boshqa matn bor-yoʻqligini aytadi (true/false).
	for _, kerak := range []string{"func TestParseSold(", "t.Run("} {
		if !strings.Contains(string(kod), kerak) {
			t.Errorf("sold_test.go ichida %q topilmadi: test TestParseSold deb nomlansin va har qatorni t.Run bilan ishga tushirsin (TASK.md, 1-qadam)", kerak)
		}
	}
}

// 2-qadam: parseSold 1 dan 1000 gacha qabul qiladi, qolganini xato bilan rad etadi.
func TestVazifaParseSold(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    int
		wantErr bool
	}{
		{"oddiy", "5", 5, false},
		{"eng_kam", "1", 1, false},
		{"chegara", "1000", 1000, false},
		{"nol", "0", 0, true},
		{"ortiq", "1001", 0, true},
		{"million", "1000000", 0, true},
		{"harf", "besh", 0, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseSold(tc.in)
			if (err != nil) != tc.wantErr {
				t.Fatalf("parseSold(%q) = %d, %v; xato kutilgan: %v (TASK.md, 2-qadam)", tc.in, got, err, tc.wantErr)
			}
			if got != tc.want {
				t.Errorf("parseSold(%q) = %d; kutilgan %d", tc.in, got, tc.want)
			}
		})
	}
}

// Xato matni ham aniq boʻlsin: foydalanuvchi uni terminalda oʻqiydi.
func TestVazifaXatoMatni(t *testing.T) {
	_, err := parseSold("1001")
	if err == nil {
		t.Fatalf("parseSold(%q) xato qaytarishi kerak, hozir xato yoʻq (TASK.md, 2-qadam)", "1001")
	}
	want := "miqdor 1000 donadan oshmasin, keldi 1001"
	if got := fmt.Sprintf("%v", err); got != want {
		t.Errorf("parseSold(%q) xatosi %q boʻlishi kerak, hozir %q", "1001", want, got)
	}
}
