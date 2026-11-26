// Package guru — Guru: rune'larni sanaydigan Writer. Boʻlak chegarasi rune'ni kesib oʻtsa ham
// toʻgʻri sanaydi: tugallanmagan 1–3 baytni keyingi Write'gacha saqlab turadi.
// Oqim tugaganda dum'da qolgan baytlarni Close sanaydi (buzuq oxir — har bayt bitta RuneError, utf8.RuneCount kabi).
package guru

import "unicode/utf8"

type RuneSanagich struct {
	Rune int64
	dum  []byte // oldingi boʻlakdan qolgan, hali tugallanmagan rune baytlari (≤ 3)
}

func (s *RuneSanagich) Write(p []byte) (int, error) {
	b := p
	if len(s.dum) > 0 {
		b = append(s.dum, p...) // faqat dum boʻlganda nusxa (koʻpi bilan 3 + len(p))
		s.dum = nil
	}
	for len(b) > 0 {
		if !utf8.FullRune(b) { // rune keyingi boʻlakda davom etadi
			s.dum = append([]byte(nil), b...)
			break
		}
		_, size := utf8.DecodeRune(b)
		s.Rune++
		b = b[size:]
	}
	return len(p), nil
}

// Close oqim oxirida qolgan tugallanmagan baytlarni ham sanaydi.
func (s *RuneSanagich) Close() error {
	s.Rune += int64(utf8.RuneCount(s.dum))
	s.dum = nil
	return nil
}
