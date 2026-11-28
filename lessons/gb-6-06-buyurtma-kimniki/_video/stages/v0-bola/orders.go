// _video bosqichi v0-bola: XATO versiya — Order egasini tekshirmaydi (BOLA). Faqat videodagi namoyish uchun.
package store

// Mahsulotlar (GB-5.10; GB-6.06: oʻzgartirish faqat admin'ga) va buyurtmalar (GB-5.10), egasi bilan (GB-6.06).
//
// Qoida: buyurtmani ID boʻyicha oʻqiydigan yoki oʻzgartiradigan har metod userID'ni ham oladi va
// uni SQL shartiga qoʻyadi. ID mijozdan keladi (yoʻldagi {id}), userID esa session'dan: mijoz uni
// tanlay olmaydi. Begona buyurtma uchun ham, yoʻq buyurtma uchun ham natija bir xil — ErrOrderNotFound.

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"slices"
	"unicode/utf8"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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

// Order — orders qatori va uning order_items qatorlari. UserID — buyurtma egasi (users.id).
// Note — yetkazish izohi (manzil), boʻlmasa JSON'ga chiqmaydi. Total = Σ quantity × unit_price (tiyin).
type Order struct {
	ID     int64       `json:"id"`
	UserID int64       `json:"user_id"`
	Status string      `json:"status"`
	Note   string      `json:"note,omitempty"`
	Items  []OrderItem `json:"items"`
	Total  int64       `json:"total"`
}

// ErrNotFound — bunday ID'li mahsulot yoʻq.
var ErrNotFound = errors.New("mahsulot topilmadi")

// ErrOrderNotFound — bunday ID'li buyurtma yoʻq yoki u boshqa foydalanuvchiniki (ikkalasi bir xil).
var ErrOrderNotFound = errors.New("buyurtma topilmadi")

// ErrBadQuantity — soʻralgan dona soni musbat emas yoki savat boʻsh.
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

// PlaceOrder userID foydalanuvchisi uchun buyurtmani bitta transaction'da yozadi (GB-5.07, GB-5.10):
// orders qatori, har mahsulot uchun qoldiq kamayishi (shartli UPDATE, GB-5.08) va order_items qatori.
// Biror qadam xato bersa, hech narsa yozilmaydi.
func (s *Store) PlaceOrder(ctx context.Context, userID int64, note string, items []OrderItem) (Order, error) {
	if len(items) == 0 {
		return Order{}, fmt.Errorf("%w: savat boʻsh", ErrBadQuantity)
	}
	for _, it := range items {
		if it.Quantity <= 0 {
			return Order{}, fmt.Errorf("%w, hozir %d", ErrBadQuantity, it.Quantity)
		}
	}
	// Qatorlar doim ID tartibida qulflanadi: ikki savat bir-birini kutib qotmaydi (GB-5.08).
	items = slices.Clone(items)
	slices.SortFunc(items, func(a, b OrderItem) int { return cmp.Compare(a.ProductID, b.ProductID) })

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return Order{}, err
	}
	defer tx.Rollback(ctx)

	o := Order{UserID: userID, Note: note, Items: []OrderItem{}}
	err = tx.QueryRow(ctx,
		"INSERT INTO orders (user_id, note) VALUES ($1, NULLIF($2, '')) RETURNING id, status",
		userID, note).Scan(&o.ID, &o.Status)
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

// whyNoRow shartli UPDATE nol qator qaytarganda sababini aniqlaydi (GB-5.08): mahsulot yoʻq → ErrNotFound,
// bor → *StockError. Oʻsha transaction ichida oʻqiydi.
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

// Order id'li buyurtmani qatorlari bilan oʻqiydi — faqat u userID'niki boʻlsa (GB-6.06).
// Yoʻq buyurtma ham, begona buyurtma ham: ErrOrderNotFound. Ikkinchi SELECT'ga egasi kerak emas:
// unga faqat birinchi SELECT topgan (demak shu foydalanuvchiniki boʻlgan) buyurtma ID'si boradi.
func (s *Store) Order(ctx context.Context, id int64) (Order, error) {
	o := Order{ID: id, Items: []OrderItem{}}
	var note *string // orders.note NULL boʻlishi mumkin: NULL → nil
	err := s.pool.QueryRow(ctx, `
		SELECT user_id, status, note FROM orders
		WHERE id = $1`,
		id).Scan(&o.UserID, &o.Status, &note)
	if errors.Is(err, pgx.ErrNoRows) {
		return Order{}, fmt.Errorf("buyurtma %d: %w", id, ErrOrderNotFound)
	}
	if err != nil {
		return Order{}, err
	}
	if note != nil {
		o.Note = *note
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

// UpdateNote buyurtmaning izohini (yetkazish manzilini) oʻzgartiradi; boʻsh satr — izohni boʻshatish.
// VAZIFA (TASK.md, 1-qadam): bu metodda videodagi xato bor — kimning buyurtmasi ekani tekshirilmaydi,
// istalgan foydalanuvchi istalgan buyurtmaning manzilini oʻzgartira oladi. Order metodidagidek tuzating.
func (s *Store) UpdateNote(ctx context.Context, id int64, note string) error {
	tag, err := s.pool.Exec(ctx,
		"UPDATE orders SET note = $2 WHERE id = $1",
		id, note)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("buyurtma %d: %w", id, ErrOrderNotFound)
	}
	return nil
}
