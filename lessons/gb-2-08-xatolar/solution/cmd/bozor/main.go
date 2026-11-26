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

	err = st.Sell(1, 41)
	if err != nil {
		err = fmt.Errorf("buyurtma 7: %w", err)
	}
	var se *store.StockError
	if errors.As(err, &se) {
		fmt.Printf("%s: faqat %d dona qoldi\n", se.Name, se.Stock)
	} else if err != nil {
		fmt.Println("xato:", err)
	}
}

// message xatoni (err nil emas) foydalanuvchiga koʻrsatiladigan matnga aylantiradi.
// Xatoni matnga aylantirish bitta joyda: keyinroq server ham shunday funksiyadan foydalanadi.
// errors.Is va errors.As butun zanjirni tekshiradi, shuning uchun qoʻshimcha
// oʻralgan xatolar («buyurtma 7: …») ham taniladi.
func message(err error) string {
	if errors.Is(err, store.ErrNotFound) {
		return "bunday mahsulot yoʻq"
	}
	if errors.Is(err, store.ErrBadQuantity) {
		return "dona soni notoʻgʻri"
	}
	var se *store.StockError
	if errors.As(err, &se) {
		return fmt.Sprintf("%s: faqat %d dona qoldi", se.Name, se.Stock)
	}
	return "xato: " + err.Error()
}
