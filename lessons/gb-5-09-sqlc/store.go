package main

import (
	"context"
	"errors"
	"math"

	"github.com/abdurahmon27/bozor/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Product — products jadvalidagi bitta qator (narx tiyinda, GB-1.03).
type Product struct {
	ID    int64
	Name  string
	Price int64
	Stock int
}

// ErrNotFound — jadvalda bunday mahsulot yoʻq (GB-2.08: chaqiruvchi errors.Is bilan tekshiradi).
var ErrNotFound = errors.New("mahsulot topilmadi")

// ErrNotEnoughStock — omborda yetarli qoldiq yoʻq: baza CHECK (stock >= 0) bilan rad etdi (GB-5.07).
var ErrNotEnoughStock = errors.New("qoldiq yetarli emas")

// ErrInvalidQuantity — miqdor noldan katta boʻlishi va integer ustuniga (int32) sigʻishi kerak.
var ErrInvalidQuantity = errors.New("miqdor notoʻgʻri: 1 dan 2147483647 gacha boʻlishi kerak")

// OrderItem — buyurtmadagi bitta qator: qaysi mahsulot va nechta.
type OrderItem struct {
	ProductID int64
	Quantity  int
}

// Store — bozor jadvallariga Go'dan kirish. pool — transaction ochish uchun (Begin),
// q — sqlc query.sql'dan generatsiya qilgan soʻrovlar (db papkasi).
type Store struct {
	pool *pgxpool.Pool
	q    *db.Queries
}

// NewStore pool ustida Store yaratadi: db.New pool'ni sqlc soʻrovlariga ulaydi.
func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool, q: db.New(pool)}
}

// Get id boʻyicha bitta mahsulotni oʻqiydi; yoʻq boʻlsa ErrNotFound.
func (s *Store) Get(ctx context.Context, id int64) (Product, error) {
	r, err := s.q.GetProduct(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return Product{}, ErrNotFound
	}
	if err != nil {
		return Product{}, err
	}
	return Product{ID: r.ID, Name: r.Name, Price: r.Price, Stock: int(r.Stock)}, nil
}

// ListActive faol mahsulotlarni id tartibida oʻqiydi. Mahsulot yoʻq boʻlsa boʻsh (nil emas) slice qaytaradi.
func (s *Store) ListActive(ctx context.Context) ([]Product, error) {
	rows, err := s.q.ListActiveProducts(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]Product, 0, len(rows))
	for _, r := range rows {
		out = append(out, Product{ID: r.ID, Name: r.Name, Price: r.Price, Stock: int(r.Stock)})
	}
	return out, nil
}

// PlaceOrder buyurtmani bitta transaction'da yozadi (GB-5.07): orders qatori, har mahsulot uchun
// products.stock kamayishi va order_items qatori. Soʻrovlar sqlc'niki, WithTx ularni tx ustida ishlatadi.
func (s *Store) PlaceOrder(ctx context.Context, customer string, items []OrderItem) (int64, error) {
	for _, it := range items {
		if it.Quantity <= 0 || it.Quantity > math.MaxInt32 {
			return 0, ErrInvalidQuantity
		}
	}
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)
	qtx := s.q.WithTx(tx)

	orderID, err := qtx.CreateOrder(ctx, customer)
	if err != nil {
		return 0, err
	}
	for _, it := range items {
		var price int64
		price, err = qtx.DecrementStock(ctx, db.DecrementStockParams{
			Quantity: int32(it.Quantity),
			ID:       it.ProductID,
		})
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrNotFound
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23514" { // 23514 = check_violation: CHECK (stock >= 0) buzildi
			return 0, ErrNotEnoughStock
		}
		if err != nil {
			return 0, err
		}
		err = qtx.AddOrderItem(ctx, db.AddOrderItemParams{
			OrderID:   orderID,
			ProductID: it.ProductID,
			Quantity:  int32(it.Quantity),
			UnitPrice: price,
		})
		if err != nil {
			return 0, err
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}
	return orderID, nil
}

// ---------------------------------------------------------------------------------------------
// VAZIFA (TASK.md): quyidagi metodlar hali qoʻlda yozilgan (GB-5.05, GB-5.07): SQL matni va Scan.
// ByName, Cheapest va OrderTotal'ning soʻrovlarini query.sql'ga koʻchiring, sqlc generate qiling va ularni s.q orqali yozing.
// ---------------------------------------------------------------------------------------------

// ByName nom boʻyicha bitta mahsulotni oʻqiydi; yoʻq boʻlsa ErrNotFound.
func (s *Store) ByName(ctx context.Context, name string) (Product, error) {
	var p Product
	err := s.pool.QueryRow(ctx,
		"SELECT id, name, price, stock FROM products WHERE name = $1", name).
		Scan(&p.ID, &p.Name, &p.Price, &p.Stock)
	if errors.Is(err, pgx.ErrNoRows) {
		return Product{}, ErrNotFound
	}
	return p, err
}

// Cheapest eng arzon n ta faol mahsulotni narx boʻyicha oʻsish tartibida oʻqiydi.
func (s *Store) Cheapest(ctx context.Context, n int) ([]Product, error) {
	rows, err := s.pool.Query(ctx,
		"SELECT id, name, price, stock FROM products WHERE active ORDER BY price, id LIMIT $1", n)
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

// OrderTotal bitta buyurtmaning jami summasini (tiyinda) hisoblaydi. Buyurtma yoʻq boʻlsa 0.
func (s *Store) OrderTotal(ctx context.Context, orderID int64) (int64, error) {
	var total int64
	err := s.pool.QueryRow(ctx,
		"SELECT COALESCE(SUM(quantity*unit_price), 0) FROM order_items WHERE order_id = $1", orderID).
		Scan(&total)
	return total, err
}

// Restock tovar qabul qilishni yozadi (GB-5.07, qoʻlda): har mahsulot qoldigʻi Quantity ga oshadi.
// Hammasi bitta transaction'da: bitta mahsulot topilmasa, hech birining qoldigʻi oʻzgarmaydi.
// 🟡 Ixtiyoriy (TASK.md, 4-qadam): soʻrovni query.sql'ga AddStock :execrows qilib koʻchiring.
func (s *Store) Restock(ctx context.Context, items []OrderItem) error {
	for _, it := range items {
		if it.Quantity <= 0 || it.Quantity > math.MaxInt32 {
			return ErrInvalidQuantity
		}
	}
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, it := range items {
		var stock int
		err = tx.QueryRow(ctx,
			"UPDATE products SET stock = stock + $1 WHERE id = $2 RETURNING stock",
			it.Quantity, it.ProductID).Scan(&stock)
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}
