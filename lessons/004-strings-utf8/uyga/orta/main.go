package main

import "fmt"

// 🟡 Oʻrta: sayt ismni "koʻpi bilan 5 ta belgi" deb len bilan tekshiradi.
// Ekranda "Gʻani" 5 ta belgi (G ʻ a n i), lekin dastur "juda uzun" deydi.
// 1) Nega? Avval oʻylab koʻring, keyin tushuntiring.
// 2) Dasturni tuzating, toki "Gʻani" qabul qilinsin.
func main() {
	ism := "Gʻani"
	if len(ism) <= 5 {
		fmt.Println("qabul")
	} else {
		fmt.Println("juda uzun")
	}
}
