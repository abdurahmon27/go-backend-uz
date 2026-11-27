package main

import (
	"context"
	"errors"
	"fmt"
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

// ByName nom boʻyicha bitta mahsulotni oʻqiydi; yoʻq boʻlsa ErrNotFound.
func (s *Store) ByName(ctx context.Context, name string) (Product, error) {
	r, err := s.q.ProductByName(ctx, name)
	if errors.Is(err, pgx.ErrNoRows) {
		return Product{}, ErrNotFound
	}
	if err != nil {
		return Product{}, err
	}
	return Product{ID: r.ID, Name: r.Name, Price: r.Price, Stock: int(r.Stock)}, nil
}

// Cheapest eng arzon n ta faol mahsulotni narx boʻyicha oʻsish tartibida oʻqiydi.
// Mahsulot yoʻq boʻlsa boʻsh (nil emas) slice: sqlc'ning slice'i bu holda nil boʻladi.
func (s *Store) Cheapest(ctx context.Context, n int) ([]Product, error) {
	if n < 0 || n > math.MaxInt32 { // LIMIT manfiy boʻlolmaydi; int32 parametrga sigʻishi kerak
		return nil, fmt.Errorf("Cheapest: n = %d, 0 dan %d gacha boʻlishi kerak", n, math.MaxInt32)
	}
	rows, err := s.q.CheapestProducts(ctx, int32(n))
	if err != nil {
		return nil, err
	}
	out := make([]Product, 0, len(rows))
	for _, r := range rows {
		out = append(out, Product{ID: r.ID, Name: r.Name, Price: r.Price, Stock: int(r.Stock)})
	}
	return out, nil
}

// OrderTotal bitta buyurtmaning jami summasini (tiyinda) hisoblaydi. Buyurtma yoʻq boʻlsa 0.
// query.sql'dagi ::bigint boʻlmasa, SUM(bigint) Postgres'da numeric beradi va sqlc interface{} qaytaradi.
func (s *Store) OrderTotal(ctx context.Context, orderID int64) (int64, error) {
	return s.q.OrderTotal(ctx, orderID)
}

// Restock tovar qabul qilishni yozadi: har mahsulot qoldigʻi Quantity ga oshadi, hammasi bitta transaction'da.
// AddStock :execrows yangilangan qatorlar sonini qaytaradi: 0 boʻlsa mahsulot yoʻq, hammasi bekor (ErrNotFound).
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
	qtx := s.q.WithTx(tx)

	for _, it := range items {
		n, err := qtx.AddStock(ctx, db.AddStockParams{
			Quantity: int32(it.Quantity),
			ID:       it.ProductID,
		})
		if err != nil {
			return err
		}
		if n == 0 {
			return ErrNotFound
		}
	}
	return tx.Commit(ctx)
}
