// inl: hook bilan bir xil, faqat yaratP tepasidagi go:noinline qatori oʻchirilgan.
package inl

type Buyurtma struct {
	Summa      int64
	Tolandi    bool
	Yetkazildi bool
}

// go:noinline: kompilyator funksiyani chaqiruv joyiga koʻchirmasin
// (inline qilmasin), shunda har bir funksiyani alohida oʻlchaymiz.

//go:noinline
func yarat(summa int64) Buyurtma {
	b := Buyurtma{Summa: summa}
	return b
}

func yaratP(summa int64) *Buyurtma {
	b := Buyurtma{Summa: summa}
	return &b
}
