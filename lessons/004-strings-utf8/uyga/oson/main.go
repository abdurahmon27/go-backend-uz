package main

import "fmt"

// 🟢 Oson: roʻyxatdan oʻtish formasi.
// Frontend'da <input maxLength={5}>: React "Gʻani" ni oʻtkazadi ("Gʻani".length === 5).
// Backend esa xuddi shu ismni "juda uzun" deb rad etadi.
// 1) Ishga tushirishdan OLDIN yozing: har bir ism uchun nima chiqadi?
// 2) Nega "Gʻani" rad etildi? Tuzating, toki backend frontend bilan bir xil sanasin.
func ismTekshir(ism string) string {
	if len(ism) <= 5 {
		return "qabul"
	}
	return "juda uzun"
}

func main() {
	for _, ism := range []string{"Anvar", "Gʻani", "Oʻktam"} {
		fmt.Println(ism, ismTekshir(ism))
	}
}
