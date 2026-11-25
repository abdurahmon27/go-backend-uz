// escape: kompilyator qaysi oʻzgaruvchini heap'ga chiqaradi?
// go build -gcflags=-m ./escape
package escape

import "fmt"

type Buyurtma struct {
	Summa   int64
	Tolandi bool
}

// 1. Lokal oʻzgaruvchining manzili qaytariladi:
// funksiya tugagach ham yashashi kerak.
func yangi(summa int64) *Buyurtma {
	b := Buyurtma{Summa: summa}
	return &b
}

// 2. Pointer pastga uzatiladi, lekin hech qayerda saqlanmaydi.
func summa(p *Buyurtma) int64 {
	return p.Summa
}

// 3. Pointer global oʻzgaruvchiga yoziladi.
var oxirgi *Buyurtma

func eslab(p *Buyurtma) {
	oxirgi = p
}

// 4. Qiymat interface'ga (any) aylanadi: fmt.Println ...any oladi.
func chop(b Buyurtma) {
	fmt.Println(b.Summa)
}

// 5. Lokal slice: oʻlchami kichik va oʻzgarmas — stack'da;
// juda katta yoki oʻlchami noma'lum — kompilyator qaror qiladi.
func hajm(n int) int {
	kichik := make([]byte, 64)
	katta := make([]byte, 100_000)
	nomalum := make([]byte, n)
	return len(kichik) + len(katta) + len(nomalum)
}
