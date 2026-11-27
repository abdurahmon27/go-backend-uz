// Package store — bozor katalogi: mahsulotlar, qoldiq qoidalari va xatolar (GB-2.06, GB-2.08).
// Faylga yozish va fayldan oʻqish file.go'da. Paket internal/ ichida: uni faqat bozor modulining
// oʻz kodi import qila oladi.
//
// GB-4.08: Store endi HTTP serverda ishlaydi, server esa har soʻrovni alohida goroutine'da bajaradi
// (GB-4.07). Shuning uchun products'ga HAR murojaat — oʻqish ham — mu qulfi ostida, va hech bir metod
// katalog ichiga koʻrsatkichni qulfdan tashqariga bermaydi: tashqariga faqat nusxa chiqadi.
package store

import (
	"errors"
	"fmt"
	"iter"
	"slices"
	"sync"
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

// Yuqori chegaralar (GB-4.08): API'ga istalgan son kelishi mumkin, chegarasiz narx × dona int64'dan
// oshib, manfiy summa berardi (GB-1.03: toʻlib ketish). Chegaralar store'da — CLI ham, API ham shu qoidadan oʻtadi.
const (
	maxPrice = 10_000_000 * 100 // tiyin: bitta mahsulot 10 mln soʻmdan qimmat emas (CLI'dagi maxPriceSom bilan bir xil)
	maxStock = 1_000_000        // bitta mahsulot qoldigʻining yuqori chegarasi, dona
)

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
	case p.Price > maxPrice:
		return fmt.Errorf("%w: narx %d tiyindan oshmasin, keldi %d", ErrBadProduct, maxPrice, p.Price)
	case p.Stock < 0:
		return fmt.Errorf("%w: qoldiq manfiy boʻlmasin, keldi %d", ErrBadProduct, p.Stock)
	case p.Stock > maxStock:
		return fmt.Errorf("%w: qoldiq %d donadan oshmasin, keldi %d", ErrBadProduct, maxStock, p.Stock)
	}
	return nil
}

// Store — katalog. Nol qiymati tayyor (boʻsh katalog, ochiq qulf); fayldan Load bilan yuklanadi.
// Store'ni nusxalamang (qulf ham nusxalanadi, go vet ogohlantiradi): doim *Store bilan ishlang.
type Store struct {
	mu       sync.Mutex // products'ga har murojaat shu qulf ostida (GB-4.07)
	products []Product  // ID tartibida
}

// Add katalogga mahsulot qoʻshadi: ID = eng katta ID + 1. Qoidaga zid mahsulot qoʻshilmaydi.
func (s *Store) Add(name string, price int64, stock int) (Product, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
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
// Nusxa qulf ostida olinadi, yield esa qulfsiz chaqiriladi: sikl tanasi Sell yoki Add'ni chaqirsa
// ham qotib qolmaydi (GB-4.07 Chuqurroq: yield qulf ostida boʻlsa — deadlock).
func (s *Store) All() iter.Seq[Product] {
	return func(yield func(Product) bool) {
		s.mu.Lock()
		products := slices.Clone(s.products)
		s.mu.Unlock()
		for _, p := range products {
			if !yield(p) {
				return
			}
		}
	}
}

// Get id'li mahsulotning nusxasini qaytaradi; yoʻq boʻlsa ErrNotFound.
// Nusxa (Product), koʻrsatkich (*Product) emas: koʻrsatkich qulfdan keyin ham katalog ichiga
// qarab turadi va uni oʻqish boshqa soʻrovning Sell'i bilan poygaga kiradi (GB-4.08 videosidagi bug).
func (s *Store) Get(id int64) (*Product, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p := s.find(id)
	if p == nil {
		return nil, fmt.Errorf("mahsulot %d: %w", id, ErrNotFound)
	}
	return p, nil
}

// find id'li mahsulotning katalog ichidagi manzilini qaytaradi, yoʻq boʻlsa nil (GB-2.04).
// Faqat qulf ostida chaqiriladi va natijasi metoddan tashqariga chiqmaydi.
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
	s.mu.Lock()
	defer s.mu.Unlock()
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

// maxRestock — bitta Restock'da qoʻshiladigan eng koʻp dona (xato yozilgan sondan himoya).
const maxRestock = 10_000

// Restock id'li mahsulotning qoldigʻiga quantity dona qoʻshadi va yangilangan mahsulotni qaytaradi.
// Xato boʻlsa katalog oʻzgarmaydi: ErrBadQuantity (quantity 1..maxRestock emas), ErrNotFound.
func (s *Store) Restock(id int64, quantity int) (Product, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if quantity <= 0 {
		return Product{}, fmt.Errorf("%w, hozir %d", ErrBadQuantity, quantity)
	}
	if quantity > maxRestock {
		return Product{}, fmt.Errorf("%w: bir martada koʻpi bilan %d, hozir %d", ErrBadQuantity, maxRestock, quantity)
	}
	p := s.find(id)
	if p == nil {
		return Product{}, fmt.Errorf("mahsulot %d: %w", id, ErrNotFound)
	}
	if p.Stock+quantity > maxStock {
		return Product{}, fmt.Errorf("%w: qoldiq %d donadan oshmasin, hozir %d", ErrBadQuantity, maxStock, p.Stock)
	}
	p.Stock += quantity
	return *p, nil
}

// Delete id'li mahsulotni katalogdan oʻchiradi; yoʻq boʻlsa ErrNotFound. Qolganlar oʻz tartibida qoladi.
// VAZIFA (TASK.md, 1-qadam): shu metodni yozing.
func (s *Store) Delete(id int64) error {
	return errors.New("Delete hali yozilmagan (TASK.md, 1-qadam)")
}
