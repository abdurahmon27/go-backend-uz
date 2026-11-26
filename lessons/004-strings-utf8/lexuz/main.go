// lexuz: HTML sahifada koʻrinadigan matndagi o/g dan keyingi apostrof turlarini va
// lotin soʻz ichidagi kirill harflarni sanaydi (tarmoq kerak):
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
	// Faqat sahifada koʻrinadigan matn: <head> (meta, title) tashlanadi, <script>/<style> ichidagi
	// kod ham, teglar ham olib tashlanadi. Menyular va sarlavhalar sanaladi (ular ham koʻrinadi).
	if i := strings.Index(strings.ToLower(string(b)), "<body"); i >= 0 {
		b = b[i:]
	}
	b = regexp.MustCompile(`(?is)<(script|style|noscript)\b.*?</(script|style|noscript)>`).ReplaceAll(b, nil)
	b = regexp.MustCompile(`(?s)<!--.*?-->`).ReplaceAll(b, nil)
	b = regexp.MustCompile(`<[^>]+>`).ReplaceAll(b, []byte("\n"))
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
