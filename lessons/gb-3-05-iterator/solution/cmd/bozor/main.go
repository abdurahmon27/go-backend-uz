// bozor dasturi: katalogni iterator bilan aylanadi va buyurtmadagi donalarni tartib bilan chiqaradi (GB-3.05).
//
//	go run ./cmd/bozor
package main

import (
	"fmt"
	"maps"
	"slices"

	"github.com/abdurahmon27/bozor/internal/store"
)

func main() {
	var st store.Store
	st.Add("olma", 1_200_000, 40)
	st.Add("non", 400_000, 3)
	st.Add("choy", 2_500_000, 8)
	st.Add("oʻrik", 3_500_000, 2)

	// Butun katalog: slice nusxalanmaydi.
	for p := range st.All() {
		fmt.Println(p.Name, p.Stock)
	}

	// Birinchi kam qolgan mahsulot: topilgach break, iterator ham toʻxtaydi.
	for p := range st.All() {
		if p.Stock < 5 {
			fmt.Println("birinchi kam qolgan:", p.Name)
			break
		}
	}

	// Buyurtmadagi donalar: map tartibi kafolatlanmagan (GB-2.02), kalitlarni tartiblaymiz.
	count := map[string]int{"olma": 3, "non": 2, "choy": 1}
	for _, name := range slices.Sorted(maps.Keys(count)) {
		fmt.Println(name, count[name])
	}
}
