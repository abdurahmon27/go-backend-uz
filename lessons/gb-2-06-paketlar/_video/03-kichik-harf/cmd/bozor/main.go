package main

import (
	"fmt"

	"github.com/abdurahmon27/bozor/internal/store"
)

func main() {
	var st store.Store
	st.Add("olma", 1_200_000, 40)
	st.add("non", 400_000, 100)
	for _, p := range st.List() {
		fmt.Println(p.ID, p.Name, p.Price)
	}
}
