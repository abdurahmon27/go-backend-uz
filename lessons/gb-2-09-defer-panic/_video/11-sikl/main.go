package main

import (
	"fmt"
	"io"
	"os"
)

const tiyinPerSom = 100 // 1 soʻm = 100 tiyin

// Product — katalogdagi bitta mahsulot.
type Product struct {
	ID    int64
	Name  string
	Price int64 // tiyin
	Stock int
}

// writeReport katalog hisobotini w'ga yozadi (GB-2.07).
func writeReport(w io.Writer, products []Product) error {
	for _, p := range products {
		_, err := fmt.Fprintf(w, "%-6s %6d soʻm %4d dona\n", p.Name, p.Price/tiyinPerSom, p.Stock)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	products := []Product{
		{ID: 1, Name: "olma", Price: 1_200_000, Stock: 40},
		{ID: 2, Name: "non", Price: 400_000, Stock: 100},
		{ID: 3, Name: "oʻrik", Price: 3_500_000, Stock: 5},
	}
	for i := 1; i <= 2000; i++ {
		f, err := os.Create("hisobot.txt")
		if err != nil {
			fmt.Printf("%d: %v\n", i, err)
			return
		}
		defer f.Close() // sikl ichida: main tugaguncha yopilmaydi
		err = writeReport(f, products)
		if err != nil {
			fmt.Printf("%d: %v\n", i, err)
			return
		}
	}
	fmt.Println("2000 ta hisobot yozildi")
}
