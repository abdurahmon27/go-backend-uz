package main

import (
	"fmt"
	"slices"
	"strings"
	"unicode"
)

func main() {
	baza := "Oʻzbekiston poytaxti Toshkent" // bazada: U+02BB
	soz := "O'zbekiston"                    // klaviaturadan: U+0027
	fmt.Println(strings.Contains(baza, soz))

	harfEmas := func(r rune) bool { return !unicode.IsLetter(r) }
	fmt.Println(strings.FieldsFunc("oʻzbek o'zbek o’zbek", harfEmas))

	sozlar := []string{"olma", "pomidor", "oʻrik", "zarang", "gilos", "gʻoz"}
	slices.Sort(sozlar)
	fmt.Println(sozlar)
}
