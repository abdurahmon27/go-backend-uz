package main

import "fmt"

// 🟡 Oʻrta: bildirishnoma (push, Telegram) uchun matnning qisqa koʻrinishi.
// qisqa(s, n) matnning birinchi n ta belgisini qaytarishi va kesilgan boʻlsa oxiriga "…" qoʻshishi kerak.
// Hozirgi kod baytlab kesadi: oʻzbekcha matnda belgi yarmidan kesiladi va ilovada "�" chiqadi.
// Tuzating: kesish belgilar (rune'lar) boʻyicha boʻlsin. %q buzilgan baytni \xca kabi koʻrsatadi.
func qisqa(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}

func main() {
	fmt.Printf("%q\n", qisqa("Toʻlov qabul qilindi", 3)) // kutilgan: "Toʻ…"
	fmt.Printf("%q\n", qisqa("Gʻalaba", 7))              // kutilgan: "Gʻalaba"
	fmt.Printf("%q\n", qisqa("Payme orqali toʻlov", 5))  // kutilgan: "Payme…"
}
