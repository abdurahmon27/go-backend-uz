package main

import (
	"errors"
	"fmt"
)

// Temporal'dagi xatoning soddalashtirilgan shakli (2021).
// Vazifa: saqla() xatoni qaytarishi kerak, lekin nil qaytaryapti. Tuzating.

func kodla(s string) (string, error) { return "[" + s + "]", nil }
func yoz(blob string) error          { return errors.New("baza: yozib boʻlmadi") }

func saqla(yangiShox bool) error {
	var err error
	if yangiShox {
		blob, err := kodla("tarix")
		if err != nil {
			return err
		}
		err = yoz(blob) // bu ichki err'ga yoziladi!
	}
	if err != nil {
		return err
	}
	return nil
}

func main() {
	fmt.Println("saqla:", saqla(true))
}
