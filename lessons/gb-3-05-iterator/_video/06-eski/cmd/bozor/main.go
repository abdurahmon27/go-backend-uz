package main

import (
	"fmt"
	"slices"

	"github.com/abdurahmon27/bozor/internal/store"
)

func main() {
	var st store.Store
	st.Add("olma", 1_200_000, 40)
	st.Add("non", 400_000, 3)
	st.Add("choy", 2_500_000, 8)
	st.Add("oʻrik", 3_500_000, 2)

	count := map[string]int{"olma": 3, "non": 2, "choy": 1}
	var names []string
	for name := range count {
		names = append(names, name)
	}
	slices.Sort(names)
	for _, name := range names {
		fmt.Println(name, count[name])
	}
}
