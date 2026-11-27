// Package store — bozor maʼlumotlari Postgres'da (GB-5.05..5.08): mahsulotlar, buyurtmalar va qoldiq qoidasi.
//
// GB-4.08'dagi katalog xotirada turardi va server toʻxtasa yoʻqolardi. Endi har metod SQL yuboradi.
// Qulf (sync.Mutex) kerak emas: Store'da oʻzgaruvchan maʼlumot yoʻq, pool'ni bir vaqtda koʻp goroutine
// ishlatishi mumkin, toʻgʻrilikni esa Postgres taʼminlaydi (transaction, qator qulfi, shartli UPDATE).
// Xatolar GB-4.08'dagi (ErrNotFound, *StockError, ErrBadQuantity, ErrBadProduct) va uchta yangisi
// (ErrDuplicate, ErrInUse, ErrOrderNotFound); cmd/api'dagi writeStoreError ularni bitta joyda statusga aylantiradi.
// Metodlar: List (GB-4.08'dagi All oʻrniga), Get, Add, Delete, PlaceOrder (Sell oʻrniga, savat bilan), Order.
package store

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"slices"
	"unicode/utf8"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Product — products jadvalidagi bitta qator. Tag'lar — API javobidagi kalitlar (GB-3.02).
type Product struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Price int64  `json:"price"` // tiyin
	Stock int    `json:"stock"`
}

// OrderItem — buyurtmadagi bitta qator: qaysi mahsulot, necha dona, qaysi narxda (sotilgan paytdagi narx, tiyin).
type OrderItem struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
	UnitPrice int64 `json:"unit_price"`
}

// Order — orders qatori va uning order_items qatorlari. Total = Σ quantity × unit_price (tiyin).
type Order struct {
	ID       int64       `json:"id"`
	Customer string      `json:"customer"`
	Status   string      `json:"status"`
	Items    []OrderItem `json:"items"`
	Total    int64       `json:"total"`
}

// ErrNotFound — bunday ID'li mahsulot yoʻq.
var ErrNotFound = errors.New("mahsulot topilmadi")

// ErrOrderNotFound — bunday ID'li buyurtma yoʻq.
var ErrOrderNotFound = errors.New("buyurtma topilmadi")

// ErrBadQuantity — soʻralgan dona soni musbat emas.
var ErrBadQuantity = errors.New("dona soni musbat boʻlishi kerak")

// ErrBadProduct — mahsulot qoidalarga zid (nom, narx yoki qoldiq notoʻgʻri).
var ErrBadProduct = errors.New("mahsulot notoʻgʻri")

// ErrDuplicate — shu nomli mahsulot allaqachon bor (products.name UNIQUE).
var ErrDuplicate = errors.New("bu nomli mahsulot bor")

// ErrInUse — mahsulot buyurtmalarda bor, oʻchirib boʻlmaydi (order_items.product_id → products.id).
var ErrInUse = errors.New("mahsulot buyurtmalarda bor")

// StockError — omborda soʻralgancha mahsulot yoʻq.
type StockError struct {
	Name  string
	Stock int // omborda bor
	Want  int // soʻralgan
}

func (e *StockError) Error() string {
	return fmt.Sprintf("%s: qoldiq %d, soʻralgan %d", e.Name, e.Stock, e.Want)
}

// Postgres xato kodlari (SQLSTATE): https://www.postgresql.org/docs/18/errcodes-appendix.html
const (
	codeForeignKey = "23503" // foreign_key_violation
	codeUnique     = "23505" // unique_violation
)

// pgCode xato Postgres'dan kelgan boʻlsa uning SQLSTATE kodini, aks holda "" qaytaradi (GB-5.07).
func pgCode(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code
	}
	return ""
}

// Chegaralar (GB-4.08): narx × dona int64'dan oshmasin.
const (
	maxNameRunes = 20
	maxPrice     = 10_000_000 * 100 // tiyin: 10 mln soʻm
	maxStock     = 1_000_000
)

// validate yangi mahsulotni tekshiradi (GB-4.08'dagi qoidalar). Baza ham CHECK bilan tekshiradi; bu yerda
// mijozga tushunarli matn chiqadi.
func validate(name string, price int64, stock int) error {
	switch {
	case name == "":
		return fmt.Errorf("%w: nom boʻsh", ErrBadProduct)
	case !utf8.ValidString(name):
		return fmt.Errorf("%w: nom UTF-8 emas", ErrBadProduct)
	case utf8.RuneCountInString(name) > maxNameRunes:
		return fmt.Errorf("%w: nom %d belgidan oshmasin", ErrBadProduct, maxNameRunes)
	case price <= 0 || price > maxPrice:
		return fmt.Errorf("%w: narx 1..%d tiyin boʻlsin, keldi %d", ErrBadProduct, maxPrice, price)
	case stock < 0 || stock > maxStock:
		return fmt.Errorf("%w: qoldiq 0..%d boʻlsin, keldi %d", ErrBadProduct, maxStock, stock)
	}
	return nil
}

// Store — Postgres'dagi bozor jadvallari. pool'ni bir vaqtda koʻp goroutine ishlata oladi (GB-5.05).
type Store struct {
	pool *pgxpool.Pool
}

// New pool ustida Store yasaydi.
func New(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool}
}

// List hamma mahsulotlarni ID tartibida qaytaradi (boʻsh boʻlsa boʻsh slice, nil emas).
func (s *Store) List(ctx context.Context) ([]Product, error) {
	rows, err := s.pool.Query(ctx, "SELECT id, name, price, stock FROM products ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Product{}
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// Get id'li mahsulotni oʻqiydi; yoʻq boʻlsa ErrNotFound.
func (s *Store) Get(ctx context.Context, id int64) (Product, error) {
	var p Product
	err := s.pool.QueryRow(ctx,
		"SELECT id, name, price, stock FROM products WHERE id = $1", id).
		Scan(&p.ID, &p.Name, &p.Price, &p.Stock)
	if errors.Is(err, pgx.ErrNoRows) {
		return Product{}, fmt.Errorf("mahsulot %d: %w", id, ErrNotFound)
	}
	return p, err
}

// Add yangi mahsulot yozadi; ID'ni Postgres beradi (GENERATED ALWAYS AS IDENTITY).
func (s *Store) Add(ctx context.Context, name string, price int64, stock int) (Product, error) {
	if err := validate(name, price, stock); err != nil {
		return Product{}, err
	}
	p := Product{Name: name, Price: price, Stock: stock}
	err := s.pool.QueryRow(ctx,
		"INSERT INTO products (name, price, stock) VALUES ($1, $2, $3) RETURNING id",
		name, price, stock).Scan(&p.ID)
	if pgCode(err) == codeUnique {
		return Product{}, fmt.Errorf("%s: %w", name, ErrDuplicate)
	}
	if err != nil {
		return Product{}, err
	}
	return p, nil
}

// Delete id'li mahsulotni oʻchiradi: yoʻq boʻlsa ErrNotFound, buyurtmalarda boʻlsa ErrInUse
// (order_items uni REFERENCES bilan koʻrsatib turibdi, Postgres oʻchirishni rad etadi).
func (s *Store) Delete(ctx context.Context, id int64) error {
	tag, err := s.pool.Exec(ctx, "DELETE FROM products WHERE id = $1", id)
	if pgCode(err) == codeForeignKey {
		return fmt.Errorf("mahsulot %d: %w", id, ErrInUse)
	}
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("mahsulot %d: %w", id, ErrNotFound)
	}
	return nil
}

// PlaceOrder buyurtmani bitta transaction'da yozadi (GB-5.07): orders qatori, har mahsulot uchun qoldiq
// kamayishi va order_items qatori. Qoldiq qoidasi bazada, shartli UPDATE'da (GB-5.08): bir vaqtdagi
// buyurtmalar qoldiqni hech qachon manfiy qilmaydi. Biror qadam xato bersa, hech narsa yozilmaydi.
func (s *Store) PlaceOrder(ctx context.Context, customer string, items []OrderItem) (Order, error) {
	if len(items) == 0 {
		return Order{}, fmt.Errorf("%w: savat boʻsh", ErrBadQuantity)
	}
	for _, it := range items {
		if it.Quantity <= 0 {
			return Order{}, fmt.Errorf("%w, hozir %d", ErrBadQuantity, it.Quantity)
		}
	}
	// Qatorlar doim bir xil tartibda (ID boʻyicha) qulflanadi: ikki savat bir-birini kutib qotmaydi
	// (deadlock, GB-5.08 Chuqurroq). Chaqiruvchining slice'i oʻzgarmasin: nusxa saralanadi.
	items = slices.Clone(items)
	slices.SortFunc(items, func(a, b OrderItem) int { return cmp.Compare(a.ProductID, b.ProductID) })

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return Order{}, err
	}
	defer tx.Rollback(ctx)

	o := Order{Customer: customer}
	err = tx.QueryRow(ctx,
		"INSERT INTO orders (customer) VALUES ($1) RETURNING id, status", customer).
		Scan(&o.ID, &o.Status)
	if err != nil {
		return Order{}, err
	}
	for _, it := range items {
		err = tx.QueryRow(ctx,
			"UPDATE products SET stock = stock - $1 WHERE id = $2 AND stock >= $1 RETURNING price",
			it.Quantity, it.ProductID).Scan(&it.UnitPrice)
		if errors.Is(err, pgx.ErrNoRows) {
			return Order{}, whyNoRow(ctx, tx, it)
		}
		if err != nil {
			return Order{}, err
		}
		_, err = tx.Exec(ctx,
			"INSERT INTO order_items (order_id, product_id, quantity, unit_price) VALUES ($1, $2, $3, $4)",
			o.ID, it.ProductID, it.Quantity, it.UnitPrice)
		if err != nil {
			return Order{}, err
		}
		o.Items = append(o.Items, it)
		o.Total += it.UnitPrice * int64(it.Quantity)
	}
	if err := tx.Commit(ctx); err != nil {
		return Order{}, err
	}
	return o, nil
}

// whyNoRow shartli UPDATE nol qator qaytarganda sababini aniqlaydi (GB-5.08: «nol qator: sabab noaniq»):
// mahsulot yoʻq → ErrNotFound, bor → *StockError. Oʻsha transaction ichida oʻqiydi.
func whyNoRow(ctx context.Context, tx pgx.Tx, it OrderItem) error {
	var p Product
	err := tx.QueryRow(ctx,
		"SELECT name, stock FROM products WHERE id = $1", it.ProductID).Scan(&p.Name, &p.Stock)
	if errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("mahsulot %d: %w", it.ProductID, ErrNotFound)
	}
	if err != nil {
		return err
	}
	return &StockError{Name: p.Name, Stock: p.Stock, Want: it.Quantity}
}

// Order id'li buyurtmani qatorlari bilan oʻqiydi; yoʻq boʻlsa ErrOrderNotFound.
// Qatorlar order_items.id tartibida (yozilgan tartib), Total Go'da hisoblanadi.
func (s *Store) Order(ctx context.Context, id int64) (Order, error) {
	o := Order{ID: id, Items: []OrderItem{}}
	err := s.pool.QueryRow(ctx,
		"SELECT customer, status FROM orders WHERE id = $1", id).Scan(&o.Customer, &o.Status)
	if errors.Is(err, pgx.ErrNoRows) {
		return Order{}, fmt.Errorf("buyurtma %d: %w", id, ErrOrderNotFound)
	}
	if err != nil {
		return Order{}, err
	}
	rows, err := s.pool.Query(ctx,
		"SELECT product_id, quantity, unit_price FROM order_items WHERE order_id = $1 ORDER BY id", id)
	if err != nil {
		return Order{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var it OrderItem
		if err := rows.Scan(&it.ProductID, &it.Quantity, &it.UnitPrice); err != nil {
			return Order{}, err
		}
		o.Items = append(o.Items, it)
		o.Total += it.UnitPrice * int64(it.Quantity)
	}
	if err := rows.Err(); err != nil {
		return Order{}, err
	}
	return o, nil
}
