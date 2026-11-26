// Package orta — Oʻrta: boshqa Reader'ni oʻrab, Windows qator oxiridagi '\r' baytlarini olib tashlaydigan Reader.
package orta

import "io"

type CRsiz struct{ R io.Reader }

func (c CRsiz) Read(p []byte) (int, error) {
	for {
		n, err := c.R.Read(p) // avval ichki Reader'dan oʻqiymiz
		k := 0
		for _, b := range p[:n] { // '\r' dan boshqa baytlarni chelak boshiga suramiz
			if b != '\r' {
				p[k] = b
				k++
			}
		}
		// Chelakda faqat '\r' boʻlgan boʻlsa, (0, nil) qaytarmaymiz — hujjat buni tavsiya qilmaydi, qayta oʻqiymiz.
		if k > 0 || err != nil || n == 0 {
			return k, err
		}
	}
}
