// yarat: katta test fayli. Manba — oʻz Generator'imiz, manzil — fayl, oʻrtada io.Copy.
package main

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/abdurahmon27/go-backend-uz/lessons/019-io-reader-writer/ozim"
)

func main() {
	mb, _ := strconv.ParseInt(os.Args[1], 10, 64)
	f, err := os.Create(os.Args[2])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	n, err := io.Copy(f, ozim.NewGenerator(mb<<20))
	if err == nil {
		err = f.Close()
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("yozildi: %d bayt · eng koʻp xotira: %.1fMiB\n", n, ozim.RSS())
}
