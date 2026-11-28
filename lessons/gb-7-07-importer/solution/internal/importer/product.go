package importer

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

const tiyinPerSom = 100        // 1 soʻm = 100 tiyin (GB-1.03)
const maxPriceSom = 10_000_000 // bitta mahsulot narxining yuqori chegarasi, soʻmda (GB-1.07)

// Product — katalogdagi bitta mahsulot: products jadvalining name, price, stock ustunlari.
// ID yoʻq: uni baza beradi (GENERATED ALWAYS AS IDENTITY, GB-5.03); importda kalit — nom (UNIQUE).
type Product struct {
	Name  string
	Price int64 // tiyin
	Stock int
}

// parsePrice soʻmdagi narx matnini tiyinga oʻgiradi (GB-3.01 bilan bir xil).
func parsePrice(s string) (int64, error) {
	som, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	if som <= 0 {
		return 0, fmt.Errorf("narx musbat boʻlishi kerak, keldi %d", som)
	}
	if som > maxPriceSom {
		return 0, fmt.Errorf("narx %d soʻmdan oshmasin, keldi %d", maxPriceSom, som)
	}
	return int64(som) * tiyinPerSom, nil
}

// parseLine «nom;narx;qoldiq» qatorini Product'ga aylantiradi; narx soʻmda.
// GB-3.01'dagi formatning ID'siz varianti: "olma;12000;40" → Product{Name: "olma", Price: 1200000, Stock: 40}.
func parseLine(line string) (Product, error) {
	parts := strings.Split(line, ";")
	if len(parts) != 3 {
		return Product{}, fmt.Errorf("3 ta maydon kerak, bor %d", len(parts))
	}
	name := parts[0]
	if name == "" {
		return Product{}, errors.New("nom boʻsh")
	}
	if !utf8.ValidString(name) {
		return Product{}, errors.New("nom UTF-8 emas") // GB-2.10
	}
	price, err := parsePrice(parts[1])
	if err != nil {
		return Product{}, fmt.Errorf("narx: %w", err)
	}
	stock, err := strconv.Atoi(parts[2])
	if err != nil {
		return Product{}, fmt.Errorf("qoldiq: %w", err)
	}
	if stock < 0 {
		return Product{}, fmt.Errorf("qoldiq manfiy boʻlmasin, keldi %d", stock)
	}
	return Product{Name: name, Price: price, Stock: stock}, nil
}
