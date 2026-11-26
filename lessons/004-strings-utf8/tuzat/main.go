package main

import (
	"fmt"
	"strings"
)

func main() {
	saytda := "Oʻzbekiston"
	qidiruv := "O'zbekiston"
	qidiruv = strings.ReplaceAll(qidiruv, "'", "ʻ") // har bir ' (U+0027) ni ʻ (U+02BB) ga almashtir
	fmt.Println(saytda == qidiruv)
	fmt.Println(strings.ReplaceAll("ma'no", "'", "ʻ")) // cheklov: toʻgʻrisi maʼno (ʼ U+02BC, tutuq belgisi)
}
