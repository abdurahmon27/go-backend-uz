// bozor: katalogni fayldan yoki stdin'dan qatorma-qator oʻqiydi (GB-3.01), vazifa bajarilgan.
//
//	go run . katalog.txt          fayldan
//	cat katalog.txt | go run .    stdin'dan
package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
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

// parseLine «ID;nom;narx;qoldiq» qatorini Product'ga aylantiradi; narx soʻmda.
// Masalan "1;olma;12000;40" → Product{ID: 1, Name: "olma", Price: 1200000, Stock: 40}.
// Xato matni qaysi maydon notoʻgʻri ekanini aytadi, asl xato %w bilan ichida qoladi.
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
		return Product{}, errors.New("nom UTF-8 emas") // GB-2.10
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

// loadCatalog katalogni r'dan qatorma-qator oʻqiydi: har qator bitta mahsulot.
// r — istalgan io.Reader: fayl, os.Stdin yoki testda strings.NewReader.
// Boʻsh qatorlar va # bilan boshlangan izohlar tashlab oʻtiladi; xatoda qator
// raqami aytiladi, asl xato %w bilan ichida qoladi.
func loadCatalog(r io.Reader) ([]Product, error) {
	var products []Product
	sc := bufio.NewScanner(r)
	n := 0 // oʻqilgan qatorlar soni, izoh va boʻsh qatorlar ham
	for sc.Scan() {
		n++
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		p, err := parseLine(line)
		if err != nil {
			return nil, fmt.Errorf("%d-qator: %w", n, err)
		}
		products = append(products, p)
	}
	if err := sc.Err(); err != nil {
		// Scanner n-qatordan keyingisini oʻqiy olmadi.
		return nil, fmt.Errorf("%d-qator: %w", n+1, err)
	}
	return products, nil
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

// run katalogni argumentdagi fayldan, argument boʻlmasa stdin'dan oʻqib, hisobotni chiqaradi.
func run() error {
	var r io.Reader = os.Stdin
	if len(os.Args) > 1 {
		f, err := os.Open(os.Args[1])
		if err != nil {
			return err
		}
		defer f.Close()
		r = f
	}
	products, err := loadCatalog(r)
	if err != nil {
		return err
	}
	return writeReport(os.Stdout, products)
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "xato:", err)
		os.Exit(1)
	}
}
