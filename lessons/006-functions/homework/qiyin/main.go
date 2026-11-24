package main

import "fmt"

type ishlovchi func(string) string

// logla: funksiya oladi va uni oʻraydigan yangi funksiya qaytaradi (middleware).
// Qaytgan closure `nom`, `next` va `n` oʻzgaruvchilarini ushlab qoladi.
func logla(nom string, next ishlovchi) ishlovchi {
	n := 0
	return func(s string) string {
		n++
		natija := next(s)
		fmt.Printf("[%s #%d] %q → %q\n", nom, n, s, natija)
		return natija
	}
}

func salom(ism string) string { return "Salom, " + ism }

func main() {
	h := logla("salom", salom)
	h("Aziz")
	h("Nodira")

	// ikki qavat: middleware'ni middleware bilan oʻrash ham mumkin
	h2 := logla("tashqi", logla("ichki", salom))
	h2("Gʻayrat")
}
