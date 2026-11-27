// Video, 01-open: katalog.txt'ni oʻqish uchun ochadi (hali hech narsa oʻqimaydi).
package main

import (
	"fmt"
	"os"
)

func run() error {
	f, err := os.Open("katalog.txt")
	if err != nil {
		return err
	}
	defer f.Close()
	// … f'dan oʻqiymiz
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "xato:", err)
		os.Exit(1)
	}
}
