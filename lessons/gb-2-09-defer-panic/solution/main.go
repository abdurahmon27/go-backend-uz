// bozor: hisobotni faylga yozadi; defer, panic va recover darsi (GB-2.09).
package main

import (
	"errors"
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

// ErrNotFound — katalogda bunday ID'li mahsulot yoʻq (GB-2.08).
var ErrNotFound = errors.New("mahsulot topilmadi")

// ErrBadQuantity — soʻralgan dona soni musbat emas (GB-2.08).
var ErrBadQuantity = errors.New("dona soni musbat boʻlishi kerak")

// StockError — omborda soʻralgancha mahsulot yoʻq (GB-2.08).
type StockError struct {
	Name  string
	Stock int // omborda bor
	Want  int // soʻralgan
}

func (e *StockError) Error() string {
	return fmt.Sprintf("%s: qoldiq %d, soʻralgan %d", e.Name, e.Stock, e.Want)
}

// OrderItem — buyurtmadagi bitta qator (GB-2.03).
type OrderItem struct {
	ProductID int64
	Price     int64 // tiyin, buyurtma paytidagi narx
	Quantity  int
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

// saveReport katalog hisobotini path fayliga yozadi.
// Fayl har qanday return'da yopiladi: defer f.Close().
// Oxirgi f.Close() esa yopishdagi xatoni qaytaradi; undan keyin
// defer'dagi Close yopiq faylda xato beradi, u tashlab yuboriladi
// (*os.File uchun ikkinchi Close xavfsiz, hujjatida yozilgan).
func saveReport(path string, products []Product) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	err = writeReport(f, products)
	if err != nil {
		return err
	}
	return f.Close()
}

// findProduct katalogdan ID boʻyicha mahsulotni topadi: topilsa elementning manzili, topilmasa nil (GB-2.04).
func findProduct(products []Product, id int64) *Product {
	for i := range products {
		if products[i].ID == id {
			return &products[i]
		}
	}
	return nil
}

// safely fn'ni chaqiradi; fn panic qilsa, panic'ni error qilib qaytaradi.
func safely(fn func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()
	fn()
	return nil
}

// ErrPanic — buyruq panic bilan tugadi. Chaqiruvchi errors.Is bilan tekshiradi.
var ErrPanic = errors.New("panic")

// handle bozor buyrugʻi fn'ni bajaradi: fn'ning xatosini oʻzgartirmay qaytaradi;
// fn panic qilsa, dastur yiqilmaydi: handle ErrPanic'ni oʻragan xato qaytaradi,
// masalan «sotish: panic: kutilmagan holat».
func handle(name string, fn func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s: %w: %v", name, ErrPanic, r)
		}
	}()
	return fn()
}

// reserve buyurtma qatorlari uchun qoldiqni kamaytiradi: hammasi yoki hech biri.
// Biror qator xato bersa, oldingi qatorlarda kamaytirilgan qoldiqlar joyiga qaytishi kerak.
// Xatolar GB-2.08'dagidek tekshiriladi: errors.Is(err, ErrNotFound),
// errors.Is(err, ErrBadQuantity), errors.As(err, &se) (se *StockError).
func reserve(products []Product, items []OrderItem) (err error) {
	var saved []int // boshidagi qoldiqlar
	for _, p := range products {
		saved = append(saved, p.Stock)
	}
	ok := false
	defer func() {
		if !ok { // xato bilan qaytdik yoki panic boʻldi: qoldiqni tiklaymiz
			for i := range saved {
				products[i].Stock = saved[i]
			}
		}
	}()
	for _, it := range items {
		if it.Quantity <= 0 {
			return fmt.Errorf("%w, hozir %d", ErrBadQuantity, it.Quantity)
		}
		p := findProduct(products, it.ProductID)
		if p == nil {
			return fmt.Errorf("mahsulot %d: %w", it.ProductID, ErrNotFound)
		}
		if p.Stock < it.Quantity {
			return &StockError{Name: p.Name, Stock: p.Stock, Want: it.Quantity}
		}
		p.Stock -= it.Quantity
	}
	ok = true
	return nil
}

// run dasturning asosiy ishi: defer'lari main'dagi os.Exit'dan oldin ishlab boʻladi.
// Videoda hisobot /dev/full'ga yozilgan (doim «disk toʻla»); bu yerda oddiy fayl.
func run() error {
	defer fmt.Println("defer: run tugadi")
	products := []Product{
		{ID: 1, Name: "olma", Price: 1_200_000, Stock: 40},
		{ID: 2, Name: "non", Price: 400_000, Stock: 100},
		{ID: 3, Name: "oʻrik", Price: 3_500_000, Stock: 5},
	}
	return saveReport("hisobot.txt", products)
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "xato:", err)
		os.Exit(1)
	}
}
