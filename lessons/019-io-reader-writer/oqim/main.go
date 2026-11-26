// oqim: xuddi shu sanoq, lekin fayl 32 KiB'lik "chelak" bilan boʻlakma-boʻlak oʻqiladi.
package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/abdurahmon27/go-backend-uz/lessons/019-io-reader-writer/ozim"
)

func main() {
	f, err := os.Open(os.Args[1]) // f — ochilgan fayl
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	buyurtma, chaqiruv := 0, 0
	buf := make([]byte, 32*1024) // 32 KiB'lik boʻsh chelak
	for {
		n, err := f.Read(buf) // fayl chelakni toʻldiradi
		chaqiruv++
		buyurtma += bytes.Count(buf[:n], []byte{'\n'}) // chelakning toʻlgan qismida qator oxirlarini sanaymiz

		if err == io.EOF { // fayl tugadi
			break // takrorlashni toʻxtatamiz
		}
		if err != nil { // boshqa xato: toʻxtaymiz
			log.Fatal(err)
		}
	}
	fmt.Printf("buyurtmalar: %d · eng koʻp xotira: %.1fMiB\n", buyurtma, ozim.RSS())
	fmt.Printf("Read chaqiruvi: %d\n", chaqiruv)
}
