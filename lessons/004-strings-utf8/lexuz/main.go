// lexuz: HTML sahifadagi o/g dan keyingi apostrof turlarini va
// lotin soʻz ichidagi kirill harflarni sanaydi. Videoda ishlatilmaydi (tarmoq kerak):
//
//	curl -sL https://lex.uz/uz/docs/8548754 | go run ./lexuz
package main

import (
	"fmt"
	"html"
	"io"
	"os"
	"regexp"
	"strings"
	"unicode"
)

func main() {
	b, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	// <script> ichidagi JavaScript matn emas: uni olib tashlaymiz
	b = regexp.MustCompile(`(?is)<script.*?</script>`).ReplaceAll(b, nil)
	r := []rune(html.UnescapeString(string(b)))
	sanoq := map[rune]int{}
	for i := 1; i < len(r); i++ {
		if strings.ContainsRune("'`´‘’ʻʼ", r[i]) && strings.ContainsRune("oOgG", r[i-1]) {
			sanoq[r[i]]++
		}
	}
	for _, a := range []rune("ʻ‘'’`´ʼ") {
		if sanoq[a] > 0 {
			fmt.Printf("%U  %d\n", a, sanoq[a])
		}
	}
	// kirill harf lotin harfi bilan yonma-yon (masalan "p" + kirill "e")
	for i := 1; i < len(r); i++ {
		if unicode.Is(unicode.Cyrillic, r[i]) && r[i-1] < 128 && unicode.IsLetter(r[i-1]) {
			fmt.Printf("lotin+kirill: %q (%U)\n", string(r[i-1:i+1]), r[i])
		}
	}
}
