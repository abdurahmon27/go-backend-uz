package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

const tiyinPerSom = 100        // 1 soʻm = 100 tiyin
const maxPriceSom = 10_000_000 // bitta mahsulot narxining yuqori chegarasi, soʻmda (GB-1.07)

// Product — katalogdagi bitta mahsulot.
type Product struct {
	ID    int64
	Name  string
	Price int64 // tiyin
	Stock int
}

// parsePrice soʻmdagi narx matnini tiyinga oʻgiradi (GB-1.07).
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

// parseLine «ID;nom;narx;qoldiq» qatorini Product'ga aylantiradi; narx soʻmda (GB-3.01).
// Masalan "1;olma;12000;40" → Product{ID: 1, Name: "olma", Price: 1200000, Stock: 40}.
func parseLine(line string) (Product, error) {
	parts := strings.Split(line, ";")
	if len(parts) != 4 {
		return Product{}, fmt.Errorf("4 ta maydon kerak, bor %d", len(parts))
	}
	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return Product{}, fmt.Errorf("ID: %w", err)
	}
	name := parts[1]
	if name == "" {
		return Product{}, errors.New("nom boʻsh")
	}
	if !utf8.ValidString(name) {
		return Product{}, errors.New("nom UTF-8 emas")
	}
	price, err := parsePrice(parts[2])
	if err != nil {
		return Product{}, fmt.Errorf("narx: %w", err)
	}
	stock, err := strconv.Atoi(parts[3])
	if err != nil {
		return Product{}, fmt.Errorf("qoldiq: %w", err)
	}
	if stock < 0 {
		return Product{}, fmt.Errorf("qoldiq manfiy boʻlmasin, keldi %d", stock)
	}
	return Product{ID: id, Name: name, Price: price, Stock: stock}, nil
}
