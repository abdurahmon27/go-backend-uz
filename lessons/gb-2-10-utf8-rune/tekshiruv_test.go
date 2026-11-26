package main

// Bu fayl vazifani (TASK.md) tekshiradi: go test uni main.go bilan birga ishga tushiradi.
// Uni oʻzgartirmang.

import (
	"testing"
	"unicode/utf8"
)

// 1-qadam: shortLabel rune'lar bilan sanaydi va «…» qoʻshadi.
func TestVazifaShortLabel(t *testing.T) {
	tests := []struct {
		name  string
		width int
		want  string
	}{
		{"olma", 10, "olma"},
		{"oʻrik qoqi", 10, "oʻrik qoqi"},    // 10 rune, 11 bayt: sigʻadi
		{"sariq sabzi", 10, "sariq sab…"},   // 11 rune, hammasi ASCII: 9 rune + «…»
		{"tovuq goʻshti", 10, "tovuq goʻ…"}, // kesish ʻ dan keyin
		{"paxta yogʻi", 9, "paxta yo…"},
		{"oʻrik", 3, "oʻ…"},
		{"choy", 4, "choy"},
	}
	for _, tc := range tests {
		got := shortLabel(tc.name, tc.width)
		if got != tc.want {
			t.Errorf("shortLabel(%q, %d) = %q; kerak %q (TASK.md, 1-qadam)", tc.name, tc.width, got, tc.want)
			continue
		}
		if !utf8.ValidString(got) {
			t.Errorf("shortLabel(%q, %d) = %q: UTF-8 buzilgan, rune oʻrtasidan kesilgan (TASK.md, 1-qadam)", tc.name, tc.width, got)
		}
	}
}

// 2-qadam: checkName nomni rune'lar bilan tekshiradi.
func TestVazifaCheckName(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
		why     string
	}{
		{"oʻrik", false, "toʻgʻri nom"},
		{"oʻgʻil bola koʻylagi", false, "20 rune, 23 bayt: ruxsat"},
		{"", true, "boʻsh nom"},
		{"o\xca", true, "UTF-8 emas: rune oʻrtasidan kesilgan"},
		{"oʻrik qoqi, 500 gramm", true, "21 rune: uzun"},
		{"olma qizil, 500 gramm", true, "21 rune, 21 bayt: uzun"},
	}
	for _, tc := range tests {
		err := checkName(tc.name)
		if tc.wantErr && err == nil {
			t.Errorf("checkName(%q): xato kutilgan (%s), nil qaytdi (TASK.md, 2-qadam)", tc.name, tc.why)
		}
		if !tc.wantErr && err != nil {
			t.Errorf("checkName(%q): xato kutilmagan (%s), keldi: %v (TASK.md, 2-qadam)", tc.name, tc.why, err)
		}
	}
}
