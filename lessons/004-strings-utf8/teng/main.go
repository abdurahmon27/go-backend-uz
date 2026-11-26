package main

import "fmt"

func main() {
	saytda := "Oʻzbekiston"  // saytdagi roʻyxatda: ʻ (U+02BB)
	qidiruv := "O'zbekiston" // foydalanuvchi klaviaturadan yozdi: ' (U+0027)
	fmt.Println(saytda == qidiruv)
}
