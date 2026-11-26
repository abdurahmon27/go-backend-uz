// Qiyin: mijozlar boʻyicha jami summa; jami kamayish tartibida, teng boʻlsa ism boʻyicha.
package qiyin

import (
	"cmp"
	"maps"
	"slices"
	"strings"
)

type Buyurtma struct {
	Mijoz string
	Summa int64
}

func Reyting(bb []Buyurtma) []string {
	jami := map[string]int64{}
	for _, b := range bb {
		jami[b.Mijoz] += b.Summa
	}
	return slices.SortedFunc(maps.Keys(jami), func(a, b string) int {
		return cmp.Or(
			cmp.Compare(jami[b], jami[a]), // summa: kattasi oldin
			strings.Compare(a, b),         // teng boʻlsa: ism boʻyicha
		)
	})
}
