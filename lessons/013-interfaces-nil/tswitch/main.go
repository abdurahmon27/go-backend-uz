// tswitch: type switch — bir nechta turni bitta switch'da tekshirish.
// go run ./tswitch
package main

import (
	"errors"
	"fmt"
)

type Buyurtma struct {
	Summa   int64
	Tolandi bool
}

func tavsifla(a any) string {
	switch v := a.(type) {
	case nil:
		return "hech narsa"
	case int, int64:
		return fmt.Sprintf("son: %v", v)
	case Buyurtma:
		return fmt.Sprintf("buyurtma, summa: %d", v.Summa)
	case error:
		return "xato: " + v.Error()
	default:
		return fmt.Sprintf("notanish tur: %T", v)
	}
}

func main() {
	for _, a := range []any{nil, 42, Buyurtma{Summa: 150000}, errors.New("disk toʻla"), 3.14} {
		fmt.Println(tavsifla(a))
	}
}
