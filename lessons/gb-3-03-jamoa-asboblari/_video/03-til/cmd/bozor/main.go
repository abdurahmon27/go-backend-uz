// bozor dasturi: katalogni chiqaradi, buyurtmani sotadi va katalogni JSON'da yozadi.
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/abdurahmon27/bozor/internal/store"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const tiyinPerSom = 100 // 1 soʻm = 100 tiyin

func main() {
	var st store.Store
	st.Add("olma", 1_200_000, 40)
	st.Add("non", 400_000, 100)

	til := "uz"
	if len(os.Args) > 1 {
		til = os.Args[1]
	}
	matcher := language.NewMatcher([]language.Tag{language.Uzbek, language.Russian})
	tag, _ := language.MatchStrings(matcher, til)
	printer := message.NewPrinter(tag)
	for _, p := range st.List() {
		printer.Printf("%d. %s: %d soʻm\n", p.ID, p.Name, p.Price/tiyinPerSom)
	}

	if err := st.Sell(1, 50); err != nil {
		fmt.Fprintln(os.Stderr, "xato:", err)
		os.Exit(1)
	}
	if err := st.Sell(2, 3); err != nil {
		fmt.Fprintln(os.Stderr, "xato:", err)
		os.Exit(1)
	}
	fmt.Println("buyurtma qabul qilindi")

	if err := json.NewEncoder(os.Stdout).Encode(st.List()); err != nil {
		fmt.Fprintln(os.Stderr, "xato:", err)
		os.Exit(1)
	}
}
