// Oson: ErrTolangan sentinel xatosi, %w bilan oʻrash va errors.Is bilan tekshirish.
// go run ./homework/oson
package main

import (
	"errors"
	"fmt"
)

var ErrTolangan = errors.New("buyurtma allaqachon toʻlangan")

func tola(id int, tolangan bool) error {
	if tolangan {
		return fmt.Errorf("toʻlov %d: %w", id, ErrTolangan)
	}
	return nil
}

func main() {
	err := tola(7, true)
	fmt.Println(err)
	fmt.Println("==:", err == ErrTolangan, " errors.Is:", errors.Is(err, ErrTolangan))
}
