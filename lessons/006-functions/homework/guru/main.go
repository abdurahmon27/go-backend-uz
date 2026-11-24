package main

import (
	"errors"
	"fmt"
)

// Tuzatilgan versiya: ichki blokda := ishlatmaymiz, tashqi err'ga = bilan yozamiz.

func kodla(s string) (string, error) { return "[" + s + "]", nil }
func yoz(blob string) error          { return errors.New("baza: yozib boʻlmadi") }

func saqla(yangiShox bool) error {
	var err error
	if yangiShox {
		var blob string
		blob, err = kodla("tarix")
		if err != nil {
			return err
		}
		err = yoz(blob)
	}
	return err
}

func main() {
	fmt.Println("saqla:", saqla(true))
}
