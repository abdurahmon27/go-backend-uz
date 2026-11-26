// hammasi: butun faylni bir yoʻla xotiraga yuklab, buyurtmalarni (qatorlarni) sanaydi.
package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/abdurahmon27/go-backend-uz/lessons/019-io-reader-writer/ozim"
)

func main() {
	data, err := os.ReadFile(os.Args[1]) // butun fayl — bitta katta slice
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	buyurtma := bytes.Count(data, []byte{'\n'}) // har buyurtma — bitta qator
	fmt.Printf("buyurtmalar: %d · eng koʻp xotira: %.1fMiB\n", buyurtma, ozim.RSS())
}
