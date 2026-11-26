package main

import "fmt"

const tiyinPerSom = 100

func main() {
	var vat int64 = 191_904 // choy, 8 dona: QQS tiyinda
	fmt.Printf("QQS: %d.%d soʻm\n", vat/tiyinPerSom, vat%tiyinPerSom)
}
