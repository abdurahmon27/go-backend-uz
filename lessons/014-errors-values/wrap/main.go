// wrap: %v xatoni matnga aylantiradi, %w esa asl xatoni ichida saqlaydi.
// go run ./wrap
package main

import (
	"errors"
	"fmt"
)

var ErrTopilmadi = errors.New("buyurtma topilmadi")

func top(id int) error {
	if id != 7 {
		return ErrTopilmadi
	}
	return nil
}

func olishV(id int) error {
	if err := top(id); err != nil {
		return fmt.Errorf("buyurtma %d: %v", id, err) // faqat matn qoladi
	}
	return nil
}

func olishW(id int) error {
	if err := top(id); err != nil {
		return fmt.Errorf("buyurtma %d: %w", id, err) // asl xato ichida
	}
	return nil
}

func main() {
	for _, x := range []struct {
		nom string
		err error
	}{{"oddiy", top(42)}, {"%v", olishV(42)}, {"%w", olishW(42)}} {
		fmt.Printf("%-5s  %-34q ==: %-5v  errors.Is: %v\n",
			x.nom, x.err.Error(), x.err == ErrTopilmadi, errors.Is(x.err, ErrTopilmadi))
	}
}
