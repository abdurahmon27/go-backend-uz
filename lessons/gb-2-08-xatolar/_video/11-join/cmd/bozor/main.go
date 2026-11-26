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

	err := errors.Join(st.Sell(9, 1), st.Sell(1, 41), st.Sell(2, 3))
	fmt.Println(err)
	fmt.Println(errors.Is(err, store.ErrNotFound))
}
