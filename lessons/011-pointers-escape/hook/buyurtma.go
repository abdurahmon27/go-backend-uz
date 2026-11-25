// hook: ikki funksiya bir xil buyurtma yaratadi. Biri qiymat qaytaradi,
// ikkinchisi pointer. Qaysi biri xotira ajratadi?
package hook

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

//go:noinline
func yaratP(summa int64) *Buyurtma {
	b := Buyurtma{Summa: summa}
	return &b
}
