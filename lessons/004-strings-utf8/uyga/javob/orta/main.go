package main

import "fmt"

// Javob: s[:n] n ta BAYT oladi: "Toʻlov" ning 3-bayti ʻ ning birinchi bayti (ca), shuning uchun "To\xca…".
// []rune(s) matnni rune'larga ajratadi: kesish endi butun belgilar boʻyicha.
// Eslatma: rune ham har doim "harf" emas (oʻ = o + ʻ ikki rune), lekin UTF-8 endi buzilmaydi.
func qisqa(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n]) + "…"
}

func main() {
	fmt.Printf("%q\n", qisqa("Toʻlov qabul qilindi", 3))
	fmt.Printf("%q\n", qisqa("Gʻalaba", 7))
	fmt.Printf("%q\n", qisqa("Payme orqali toʻlov", 5))
}
