// Package store — bozor katalogi: mahsulotlar, qoldiq qoidalari va xatolar (GB-2.06, GB-2.08).
// Faylga yozish va fayldan oʻqish file.go'da. Paket internal/ ichida: uni faqat bozor modulining
// oʻz kodi import qila oladi.
package store

import (
	"errors"
	"fmt"
	"iter"
	"unicode/utf8"
)

// Product — katalogdagi bitta mahsulot. Tag'lar — bozor.json'dagi kalitlar (GB-3.02).
type Product struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Price int64  `json:"price"` // tiyin
	Stock int    `json:"stock"`
}

// ErrNotFound — katalogda bunday ID'li mahsulot yoʻq.
var ErrNotFound = errors.New("mahsulot topilmadi")

// ErrBadQuantity — soʻralgan dona soni musbat emas.
var ErrBadQuantity = errors.New("dona soni musbat boʻlishi kerak")

// ErrBadProduct — mahsulot katalog qoidalariga zid (nom, narx yoki qoldiq notoʻgʻri).
var ErrBadProduct = errors.New("mahsulot notoʻgʻri")

// StockError — omborda soʻralgancha mahsulot yoʻq.
type StockError struct {
	Name  string
	Stock int // omborda bor
	Want  int // soʻralgan
}

func (e *StockError) Error() string {
	return fmt.Sprintf("%s: qoldiq %d, soʻralgan %d", e.Name, e.Stock, e.Want)
}

// maxNameRunes — mahsulot nomining eng koʻp belgisi (GB-2.10): list ustuni shunga moslangan.
const maxNameRunes = 20

// validate bitta mahsulotni tekshiradi (GB-3.01'dagi parseLine va GB-2.10'dagi checkName qoidalari).
// Add ham, fayldan oʻqilgan katalog ham (Load) shu qoidalardan oʻtadi.
func validate(p Product) error {
	switch {
	case p.ID <= 0:
		return fmt.Errorf("%w: ID musbat boʻlishi kerak, keldi %d", ErrBadProduct, p.ID)
	case p.Name == "":
		return fmt.Errorf("%w: nom boʻsh", ErrBadProduct)
	case !utf8.ValidString(p.Name):
		return fmt.Errorf("%w: nom UTF-8 emas", ErrBadProduct)
	case utf8.RuneCountInString(p.Name) > maxNameRunes:
		return fmt.Errorf("%w: nom %d belgi, eng koʻpi %d", ErrBadProduct, utf8.RuneCountInString(p.Name), maxNameRunes)
	case p.Price <= 0:
		return fmt.Errorf("%w: narx musbat boʻlishi kerak, keldi %d", ErrBadProduct, p.Price)
	case p.Stock < 0:
		return fmt.Errorf("%w: qoldiq manfiy boʻlmasin, keldi %d", ErrBadProduct, p.Stock)
	}
	return nil
}

// Store — katalog. Nol qiymati tayyor (boʻsh katalog); fayldan Load bilan yuklanadi.
type Store struct {
	products []Product
}

// Add katalogga mahsulot qoʻshadi: ID = eng katta ID + 1. Qoidaga zid mahsulot qoʻshilmaydi.
func (s *Store) Add(name string, price int64, stock int) (Product, error) {
	var last int64
	for _, p := range s.products {
		last = max(last, p.ID)
	}
	p := Product{ID: last + 1, Name: name, Price: price, Stock: stock}
	if err := validate(p); err != nil {
		return Product{}, err
	}
	s.products = append(s.products, p)
	return p, nil
}

// All katalogni ID tartibida beradigan iterator (GB-3.05): ichki slice tashqariga chiqmaydi.
func (s *Store) All() iter.Seq[Product] {
	return func(yield func(Product) bool) {
		for _, p := range s.products {
			if !yield(p) {
				return
			}
		}
	}
}

// find id'li mahsulotning katalog ichidagi manzilini qaytaradi, yoʻq boʻlsa nil (GB-2.04).
func (s *Store) find(id int64) *Product {
	for i := range s.products {
		if s.products[i].ID == id {
			return &s.products[i]
		}
	}
	return nil
}

// Sell id'li mahsulotdan quantity dona sotadi va sotuvdan keyingi mahsulotni qaytaradi.
// Xato boʻlsa katalog oʻzgarmaydi. Chaqiruvchi xatoni tekshira oladi (GB-2.08):
// errors.Is(err, ErrBadQuantity), errors.Is(err, ErrNotFound), errors.As(err, &se) (se *StockError).
func (s *Store) Sell(id int64, quantity int) (Product, error) {
	if quantity <= 0 {
		return Product{}, fmt.Errorf("%w, hozir %d", ErrBadQuantity, quantity)
	}
	p := s.find(id)
	if p == nil {
		return Product{}, fmt.Errorf("mahsulot %d: %w", id, ErrNotFound)
	}
	if p.Stock < quantity {
		return Product{}, &StockError{Name: p.Name, Stock: p.Stock, Want: quantity}
	}
	p.Stock -= quantity
	return *p, nil
}

// Restock id'li mahsulotning qoldigʻiga quantity dona qoʻshadi va yangilangan mahsulotni qaytaradi.
// VAZIFA (1-qadam, TASK.md): hozir u hech narsa qilmaydi.
func (s *Store) Restock(id int64, quantity int) (Product, error) {
	return Product{}, errors.New("Restock hali yozilmagan (TASK.md, 1-qadam)")
}
