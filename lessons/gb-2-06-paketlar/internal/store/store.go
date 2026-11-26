// Package store bozor katalogini xotirada saqlaydi.
// U internal/ ichida: uni faqat bozor modulining oʻz kodi import qila oladi.
package store

// Product — katalogdagi bitta mahsulot.
type Product struct {
	ID    int64
	Name  string
	Price int64 // tiyin
	Stock int
}

// Store — xotiradagi katalog. Nol qiymati tayyor: var st store.Store.
// products kichik harf bilan: paketdan tashqaridagi kod unga yoza olmaydi.
type Store struct {
	products []Product
}

// Add katalogga mahsulot qoʻshadi va unga ID beradi (1, 2, 3 …).
func (s *Store) Add(name string, price int64, stock int) Product {
	p := Product{ID: int64(len(s.products)) + 1, Name: name, Price: price, Stock: stock}
	s.products = append(s.products, p)
	return p
}

// List katalogdagi mahsulotlarni qaytaradi.
// TASK.md, 2-qadam: hozir u ichki slice'ning oʻzini qaytaradi, shuning uchun
// tashqaridagi kod natijaga yozib, katalogni oʻzgartira oladi.
func (s *Store) List() []Product {
	return s.products
}

// Sell id'li mahsulotdan quantity dona sotadi va qoldiqni kamaytiradi.
// TASK.md, 1-qadam: hozir u hech narsa qilmaydi va har doim nil qaytaradi.
func (s *Store) Sell(id int64, quantity int) error {
	return nil
}
