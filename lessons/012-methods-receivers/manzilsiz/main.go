// manzilsiz: pointer metodni qaysi qiymatda chaqirsa boʻladi? Faqat addressable qiymatda.
// Bu papka ATAYLAB kompilyatsiya boʻlmaydi: go build -o /dev/null ./manzilsiz
package main

type Buyurtma struct {
	Summa   int64
	Tolandi bool
}

func (b *Buyurtma) Tola() { b.Tolandi = true }

type Savat struct{ B Buyurtma }

func yarat() Buyurtma { return Buyurtma{Summa: 50_000} }

func main() {
	var b Buyurtma
	b.Tola() // ✅ oʻzgaruvchi
	s := []Buyurtma{{}}
	s[0].Tola() // ✅ slice elementi
	var a [2]Buyurtma
	a[1].Tola() // ✅ addressable massiv elementi
	var sv Savat
	sv.B.Tola() // ✅ addressable struct'ning maydoni
	p := &b
	(*p).Tola() // ✅ pointer orqali qiymat (pointer indirection)

	(&Buyurtma{}).Tola() // ✅ lekin bu yerda receiver allaqachon pointer: &T{} — & uchun istisno

	m := map[string]Buyurtma{}
	m["a1"].Tola()    // ❌ map elementi
	yarat().Tola()    // ❌ funksiya natijasi
	Buyurtma{}.Tola() // ❌ composite literal'ning oʻzi
}
