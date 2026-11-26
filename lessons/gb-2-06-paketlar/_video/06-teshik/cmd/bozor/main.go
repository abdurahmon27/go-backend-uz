package main

import (
	"fmt"

	"github.com/abdurahmon27/bozor/internal/store"
)

func main() {
	var st store.Store
	st.Add("olma", 1_200_000, 40)
	st.Add("non", 400_000, 100)
	st.List()[0].Stock = -5
	fmt.Println(st.List()[0].Stock)
}
