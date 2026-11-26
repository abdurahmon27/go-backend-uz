package main

import (
	"github.com/abdurahmon27/bozor/internal/store"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const tiyinPerSom = 100 // 1 soʻm = 100 tiyin

func main() {
	var st store.Store
	st.Add("olma", 1_200_000, 40)
	st.Add("non", 400_000, 100)

	printer := message.NewPrinter(language.Uzbek)
	for _, p := range st.List() {
		printer.Printf("%d. %s: %d soʻm\n", p.ID, p.Name, p.Price/tiyinPerSom)
	}
}
