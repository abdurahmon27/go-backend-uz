// katalog — sinov uchun katta katalog faylini yozadi (GB-7.07):
//
//	go run ./cmd/katalog 200000 > katta.txt
package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/abdurahmon27/bozor/internal/importer"
)

func main() {
	n := 200_000
	if len(os.Args) > 1 {
		v, err := strconv.Atoi(os.Args[1])
		if err != nil || v < 0 {
			fmt.Fprintln(os.Stderr, "qatorlar soni musbat butun son boʻlishi kerak:", os.Args[1])
			os.Exit(1)
		}
		n = v
	}
	if err := importer.Generate(os.Stdout, n); err != nil {
		fmt.Fprintln(os.Stderr, "xato:", err)
		os.Exit(1)
	}
}
