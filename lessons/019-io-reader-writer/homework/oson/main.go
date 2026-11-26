// Oson: oʻz `cat`ingiz. Fayl hajmi qancha boʻlsa ham xotira oʻzgarmaydi.
//
//	go run ./homework/oson /tmp/io19/buyurtma.csv > /dev/null
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/abdurahmon27/go-backend-uz/lessons/019-io-reader-writer/ozim"
)

func main() {
	for _, nom := range os.Args[1:] {
		f, err := os.Open(nom)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		_, err = io.Copy(os.Stdout, f)
		f.Close()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
	fmt.Fprintf(os.Stderr, "eng koʻp xotira: %.1fMiB\n", ozim.RSS())
}
