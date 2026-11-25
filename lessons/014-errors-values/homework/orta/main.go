// Oʻrta: bitta Errorf ichida ikkita %w (Go 1.20+).
// errors.Is ikkalasini ham topadi, errors.Unwrap esa nil qaytaradi:
// bunday xatoda Unwrap() error emas, Unwrap() []error metodi bor.
// go run ./homework/orta
package main

import (
	"errors"
	"fmt"
)

var (
	ErrDisk = errors.New("disk toʻla")
	ErrLog  = errors.New("log yozilmadi")
)

func main() {
	err := fmt.Errorf("saqlash: %w; %w", ErrDisk, ErrLog)
	fmt.Println(err)
	fmt.Println("Is(ErrDisk):", errors.Is(err, ErrDisk), " Is(ErrLog):", errors.Is(err, ErrLog))
	fmt.Println("errors.Unwrap:", errors.Unwrap(err))
	if m, ok := err.(interface{ Unwrap() []error }); ok {
		fmt.Println("Unwrap() []error:", len(m.Unwrap()), "ta")
	}
}
