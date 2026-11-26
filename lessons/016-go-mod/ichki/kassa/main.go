// kassa ombor/ dan tashqarida: internal/narx ni import qila olmaydi.
package main

import (
	"fmt"

	"ichki/ombor/internal/narx"
)

func main() {
	fmt.Println(narx.QQS(1000))
}
