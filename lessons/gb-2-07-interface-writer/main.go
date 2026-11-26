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

// String mahsulotni «olma (12000 soʻm)» koʻrinishida qaytaradi.
func (p Product) String() string {
	return fmt.Sprintf("%s (%d soʻm)", p.Name, p.Price/tiyinPerSom)
}

// writeReport katalog hisobotini w'ga yozadi.
func writeReport(w io.Writer, products []Product) error {
	for _, p := range products {
		_, err := fmt.Fprintf(w, "%-6s %6d soʻm %4d dona\n", p.Name, p.Price/tiyinPerSom, p.Stock)
		if err != nil {
			return err
		}
	}
	return nil
}

// byteCounter oʻziga yozilgan baytlarni sanaydi.
type byteCounter struct {
	n int
}

// Write p'dagi baytlar sonini n'ga qoʻshadi.
func (c *byteCounter) Write(p []byte) (int, error) {
	c.n += len(p)
	return len(p), nil
}

// lineCounter oʻziga yozilgan qatorlarni sanaydi.
// VAZIFA (2-qadam): *lineCounter io.Writer boʻlsin (TASK.md).
type lineCounter struct {
	lines int
}

// writeLowStock qoldigʻi limit'dan kam mahsulotlarni w'ga yozadi:
// har biriga bitta qator, masalan «oʻrik (35000 soʻm) — 5 dona qoldi».
// VAZIFA (1-qadam): hozir u w'ga emas, ekranga yozadi (TASK.md).
func writeLowStock(w io.Writer, products []Product, limit int) error {
	for _, p := range products {
		if p.Stock < limit {
			fmt.Printf("%v — %d dona qoldi\n", p, p.Stock)
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

	f, err := os.Create("hisobot.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "xato:", err)
		os.Exit(1)
	}
	err = writeReport(f, products)
	if err != nil {
		fmt.Fprintln(os.Stderr, "xato:", err)
		os.Exit(1)
	}
	err = f.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, "xato:", err)
		os.Exit(1)
	}

	var c byteCounter
	err = writeReport(&c, products)
	if err != nil {
		fmt.Fprintln(os.Stderr, "xato:", err)
		os.Exit(1)
	}
	fmt.Println("hisobot:", c.n, "bayt")

	fmt.Println(products)
}
