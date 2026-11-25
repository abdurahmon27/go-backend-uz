// narx: non narxini chegirma bilan chiqaradi. Foiz birinchi argument: go run . 10
package main

import (
	"fmt"
	"os"
	"strconv"
)

const non = 3000 // so'm

func main() {
	foiz := 0
	if len(os.Args) > 1 {
		f, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Println("foiz butun son bo'lishi kerak:", err)
			os.Exit(2)
		}
		foiz = f
	}
	narx, err := chegirma(non, foiz)
	if err != nil {
		fmt.Println("xato:", err)
		os.Exit(2)
	}
	fmt.Printf("non: %d so'm\n", narx)
}
