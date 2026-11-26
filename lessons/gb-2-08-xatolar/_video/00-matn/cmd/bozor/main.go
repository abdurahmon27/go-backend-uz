// bozor dasturi: mahsulot sotadi va xatoni turiga qarab tushuntiradi.
package main

import (
	"fmt"

	"github.com/abdurahmon27/bozor/internal/store"
)

func main() {
	var st store.Store
	st.Add("olma", 1_200_000, 40)
	st.Add("non", 400_000, 100)

	err := st.Sell(9, 1)
	fmt.Println(err)
	err = st.Sell(1, 41)
	fmt.Println(err)
}
