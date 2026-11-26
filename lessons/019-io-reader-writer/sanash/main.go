// sanash: oqim/ dagi siklni io.Copy bajaradi: manba — fayl (Reader), manzil — Hisoblagich (Writer).
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/abdurahmon27/go-backend-uz/lessons/019-io-reader-writer/ozim"
)

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	var h ozim.Hisoblagich
	_, err = io.Copy(&h, f) // qayerga ← qayerdan: Hisoblagich (&h — nusxasi emas, oʻzi) ← fayl
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("buyurtmalar: %d · bayt: %d · eng koʻp xotira: %.1fMiB\n", h.Qator, h.Bayt, ozim.RSS())
}
