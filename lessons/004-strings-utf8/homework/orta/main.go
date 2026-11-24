// Oʻrta: oʻzbekcha soʻzdagi harflarni sanash: oʻ, gʻ, sh, ch bitta harf.
// Tutuq belgisi (ʼ) alohida harf emas (CLDR: common/collation/uz.xml izohi).
package main

import (
	"fmt"
	"strings"
)

// harflar soʻzni oʻzbek harflariga ajratadi. Avval apostroflar
// darsdagi tozala() bilan bir xil qilingan deb hisoblaymiz.
func harflar(s string) []string {
	r := []rune(s)
	var out []string
	for i := 0; i < len(r); i++ {
		if r[i] == 'ʼ' {
			continue // tutuq belgisi harf emas
		}
		if i+1 < len(r) {
			pair := strings.ToLower(string(r[i : i+2]))
			switch pair {
			case "oʻ", "gʻ", "sh", "ch":
				out = append(out, string(r[i:i+2]))
				i++
				continue
			}
		}
		out = append(out, string(r[i]))
	}
	return out
}

func main() {
	for _, s := range []string{"oʻzbek", "shahar", "choʻgʻ", "maʼno", "Isʼhoq"} {
		h := harflar(s)
		fmt.Printf("%-7s rune=%d harf=%d %q\n", s, len([]rune(s)), len(h), h)
	}
}
