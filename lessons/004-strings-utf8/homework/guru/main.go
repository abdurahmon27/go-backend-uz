// Guru: soʻzlarni oʻzbek alifbosi tartibida saralash.
// Hozirgi (1995) alifbo: ... x y z oʻ gʻ sh ch (CLDR: common/collation/uz.xml).
// Oddiy slices.Sort baytlarni solishtiradi, shuning uchun bu tartibni bilmaydi.
package main

import (
	"fmt"
	"slices"
	"strings"
)

var alifbo = strings.Fields("a b d e f g h i j k l m n o p q r s t u v x y z oʻ gʻ sh ch")

// tartib: har bir harfning alifbodagi oʻrni; notanish belgi oxirida, rune qiymati boʻyicha.
func tartib(soz string) []int {
	r := []rune(strings.ToLower(soz))
	var out []int
	for i := 0; i < len(r); i++ {
		if r[i] == 'ʼ' {
			continue
		}
		if i+1 < len(r) {
			if k := slices.Index(alifbo, string(r[i:i+2])); k >= 0 {
				out = append(out, k)
				i++
				continue
			}
		}
		if k := slices.Index(alifbo, string(r[i])); k >= 0 {
			out = append(out, k)
		} else {
			out = append(out, len(alifbo)+int(r[i]))
		}
	}
	return out
}

func main() {
	sozlar := []string{"olma", "pomidor", "oʻrik", "zarang", "gilos", "gʻoz", "shaftoli", "sabzi", "choy", "chumoli"}
	slices.Sort(sozlar)
	fmt.Println("bayt: ", sozlar)
	slices.SortFunc(sozlar, func(a, b string) int { return slices.Compare(tartib(a), tartib(b)) })
	fmt.Println("alifbo:", sozlar)
}
