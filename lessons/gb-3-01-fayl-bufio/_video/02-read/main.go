// Video, 02-read: f.Read'ni oʻzimiz chaqiramiz — 20 baytlik buferga boʻlak-boʻlak.
package main

import (
	"fmt"
	"io"
	"os"
)

func run() error {
	f, err := os.Open("katalog.txt")
	if err != nil {
		return err
	}
	defer f.Close()
	buf := make([]byte, 20)
	for {
		n, err := f.Read(buf)
		fmt.Printf("%d %q\n", n, buf[:n])
		if err == io.EOF {
			break // fayl tugadi
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "xato:", err)
		os.Exit(1)
	}
}
