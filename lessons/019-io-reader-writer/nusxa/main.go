// nusxa: fayldan faylga io.Copy. Baytlar dasturimiz orqali oʻtadimi?
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/abdurahmon27/go-backend-uz/lessons/019-io-reader-writer/ozim"
)

func main() {
	src, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer src.Close()
	dst, err := os.Create(os.Args[2])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	n, err := io.Copy(dst, src) // dst — *os.File, unda ReadFrom bor
	if err == nil {
		err = dst.Close()
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("nusxa: %d bayt · eng koʻp xotira: %.1fMiB\n", n, ozim.RSS())
}
