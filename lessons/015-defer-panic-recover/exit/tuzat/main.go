// Tuzatish: ish run() ichida, uning defer'lari ishlab bo'lgach main chiqadi.
package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	if err := run(); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func run() error {
	defer fmt.Println("defer: buferdagi loglar yozildi")

	if err := sozla(); err != nil {
		return fmt.Errorf("sozlash: %w", err)
	}
	return nil
}

func sozla() error {
	return fmt.Errorf("PORT berilmagan")
}
