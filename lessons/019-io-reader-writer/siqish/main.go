// siqish: Writer ichida Writer — gzip baytlarni siqadi va ichidagi faylga yozadi.
package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"

	"github.com/abdurahmon27/go-backend-uz/lessons/019-io-reader-writer/ozim"
)

func main() {
	f, err := os.Open(os.Args[1])
	tek(err)
	defer f.Close()
	out, err := os.Create(os.Args[2])
	tek(err)

	gz, _ := gzip.NewWriterLevel(out, gzip.BestSpeed) // gz — arxivlovchi: yangi faylni (out) oʻrab oladi, eng tez rejimda
	n, err := io.Copy(gz, f)                          // nasos: fayl → arxivlovchi → yangi fayl
	tek(err)
	tek(gz.Close()) // qolgan siqilgan baytlarni yozib tugatadi
	tek(out.Close())

	st, err := os.Stat(os.Args[2])
	tek(err)
	mib := func(b int64) float64 { return float64(b) / (1 << 20) }
	fmt.Printf("oʻqildi: %.1fMiB · yozildi: %.1fMiB · eng koʻp xotira: %.1fMiB\n", mib(n), mib(st.Size()), ozim.RSS())
}

func tek(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
