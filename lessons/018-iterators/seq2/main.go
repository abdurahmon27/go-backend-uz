// Seq2: har qadamda ikki qiymat (indeks va qiymat, kalit va qiymat).
package main

import (
	"fmt"
	"iter"
	"slices"
	"sync"
)

func main() {
	summalar := []int64{150000, 42000}
	for i, s := range slices.All(summalar) { // iter.Seq2[int, int64]
		fmt.Println(i, s)
	}

	// sync.Map'ning Range metodi Go 1.9 dan beri bor va aynan Seq2 shaklida.
	var m sync.Map
	m.Store("Aziz", 150000)
	var _ iter.Seq2[any, any] = m.Range
	for k, v := range m.Range {
		fmt.Println(k, v)
	}
}
