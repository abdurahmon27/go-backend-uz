// bozor dasturi: mahsulot sotadi va xatoni turiga qarab tushuntiradi.
package main

import (
	"errors"
	"fmt"

	"github.com/abdurahmon27/bozor/internal/store"
)

func main() {
	var st store.Store
	st.Add("olma", 1_200_000, 40)
	st.Add("non", 400_000, 100)

	err := st.Sell(9, 1)
	if errors.Is(err, store.ErrNotFound) {
		fmt.Println("bunday mahsulot yoʻq")
	} else if err != nil {
		fmt.Println("xato:", err)
	}
	fmt.Println(errors.Is(err, errors.New("mahsulot topilmadi")))
}
