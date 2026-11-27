package main

import (
	"fmt"

	"github.com/abdurahmon27/bozor/internal/store"
)

func main() {
	var st store.Store
	st.Add("olma", 1_200_000, 40)
	st.Add("non", 400_000, 3)
	st.Add("choy", 2_500_000, 8)
	st.Add("oʻrik", 3_500_000, 2)

	for _, p := range st.List() {
		if p.Stock < 5 {
			fmt.Println("birinchi kam qolgan:", p.Name)
			break
		}
	}
}
