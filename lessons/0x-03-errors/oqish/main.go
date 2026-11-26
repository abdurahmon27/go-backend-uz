// error qiymat: fayl yo'q bo'lsa, data bo'sh, err esa xato matnini olib keladi.
package main

import (
	"fmt"
	"os"
)

func main() {
	data, err := os.ReadFile("yoq.json")
	fmt.Println(len(data), err)
}
