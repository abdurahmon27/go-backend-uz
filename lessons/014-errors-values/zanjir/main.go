// zanjir: %w bilan oʻralgan xatolar zanjiri. errors.Unwrap har safar bitta qavatni ochadi.
// go run ./zanjir
package main

import (
	"errors"
	"fmt"
)

var ErrTopilmadi = errors.New("buyurtma topilmadi")

func top(id int) error { return ErrTopilmadi }

func olish(id int) error {
	if err := top(id); err != nil {
		return fmt.Errorf("buyurtma %d: %w", id, err)
	}
	return nil
}

func tolov(id int) error {
	if err := olish(id); err != nil {
		return fmt.Errorf("toʻlov: %w", err)
	}
	return nil
}

func main() {
	err := tolov(42)
	for i, e := 0, err; e != nil; i, e = i+1, errors.Unwrap(e) {
		fmt.Printf("%d  %-20T  %s\n", i, e, e)
	}
	fmt.Println("errors.Is(err, ErrTopilmadi):", errors.Is(err, ErrTopilmadi))
}
