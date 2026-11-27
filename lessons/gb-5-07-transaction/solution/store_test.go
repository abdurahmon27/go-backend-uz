// Bu fayl vazifani tekshiradi: Store metodlarini haqiqiy Postgres'da ishga tushiradi.
// Uni oʻzgartirmang. Postgres uchun: docker compose up -d --wait (shu papkada).
// Har test alohida vaqtinchalik sxemada ishlaydi (bozor.sql shu yerga qoʻllanadi)
// va oxirida sxemani oʻchiradi: public sxemaga tegmaydi. Boshqa manzil kerak boʻlsa DATABASE_URL qoʻying.
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const labURL = "postgres://bozor:bozor-lab@localhost:5497/bozor"

// newStore bu test uchun boʻsh sxema yaratadi, bozor.sql ni qoʻllaydi va shu sxemaga ulangan Store qaytaradi.
func newStore(t *testing.T) *Store {
	t.Helper()
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		url = labURL
	}
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		t.Fatalf("DATABASE_URL notoʻgʻri: %v", err)
	}
	cfg.MaxConns = 4 // DATABASE_URL da pool_max_conns=1 qolgan boʻlsa ham test toʻxtab qolmasin
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	admin, err := pgxpool.NewWithConfig(ctx, cfg)
	if err == nil {
		err = admin.Ping(ctx)
	}
	if err != nil {
		t.Fatalf("Postgres'ga ulanib boʻlmadi (%s:%d): avval shu papkada `docker compose up -d --wait` ishga tushiring.\n%v", cfg.ConnConfig.Host, cfg.ConnConfig.Port, err)
	}
	schema := fmt.Sprintf("lab_%d", time.Now().UnixNano())
	if _, err := admin.Exec(ctx, "CREATE SCHEMA "+schema); err != nil {
		t.Fatalf("sxema yaratilmadi: %v", err)
	}
	t.Cleanup(func() {
		if _, err := admin.Exec(context.Background(), "DROP SCHEMA "+schema+" CASCADE"); err != nil {
			t.Logf("sxema %s oʻchirilmadi (qoʻlda: DROP SCHEMA %s CASCADE): %v", schema, schema, err)
		}
		admin.Close()
	})

	poolCfg := cfg.Copy()
	poolCfg.ConnConfig.RuntimeParams["search_path"] = schema
	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		t.Fatalf("pool: %v", err)
	}
	t.Cleanup(pool.Close)
	b, err := os.ReadFile("bozor.sql")
	if err != nil {
		t.Fatalf("bozor.sql oʻqilmadi: %v", err)
	}
	// Oddiy protokol: bitta matnda bir nechta SQL buyrugʻi boʻlishi mumkin.
	if _, err := pool.Exec(ctx, string(b), pgx.QueryExecModeSimpleProtocol); err != nil {
		t.Fatalf("bozor.sql bajarilmadi: %v", err)
	}
	return &Store{pool: pool}
}

func testCtx(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	t.Cleanup(cancel)
	return ctx
}

// count jadvaldagi qatorlar sonini qaytaradi.
func count(t *testing.T, s *Store, table string) int {
	t.Helper()
	var n int
	if err := s.pool.QueryRow(testCtx(t), "SELECT count(*) FROM "+table).Scan(&n); err != nil {
		t.Fatalf("%s qatorlari sanalmadi: %v", table, err)
	}
	return n
}

func stockOf(t *testing.T, s *Store, id int64) int {
	t.Helper()
	var n int
	if err := s.pool.QueryRow(testCtx(t), "SELECT stock FROM products WHERE id = $1", id).Scan(&n); err != nil {
		t.Fatalf("mahsulot %d qoldigʻi oʻqilmadi: %v", id, err)
	}
	return n
}

// Boshlangʻich holat (bozor.sql): olma(1) 40, non(2) 3, choy(3) 8, oʻrik(4) 2; 3 ta buyurtma, 4 ta buyurtma qatori.

func TestPlaceOrder(t *testing.T) {
	s := newStore(t)
	id, err := s.PlaceOrder(testCtx(t), "Aziza", []OrderItem{{ProductID: 1, Quantity: 2}, {ProductID: 2, Quantity: 1}})
	if err != nil {
		t.Fatalf("PlaceOrder xato qaytardi: %v; kutilgan: nil", err)
	}
	if id != 4 {
		t.Errorf("buyurtma id = %d, kutilgan: 4 (oldin 3 ta buyurtma bor edi)", id)
	}
	if got := count(t, s, "order_items WHERE order_id = 4"); got != 2 {
		t.Errorf("4-buyurtmada %d ta qator bor, kutilgan: 2", got)
	}
	if got := stockOf(t, s, 1); got != 38 {
		t.Errorf("olma qoldigʻi = %d, kutilgan: 38", got)
	}
	if got := stockOf(t, s, 2); got != 2 {
		t.Errorf("non qoldigʻi = %d, kutilgan: 2", got)
	}
}

// Ikkinchi mahsulotda qoldiq yetmaydi: birinchi qadamlar ham bekor boʻlishi kerak (transaction).
func TestPlaceOrderRollsBack(t *testing.T) {
	s := newStore(t)
	_, err := s.PlaceOrder(testCtx(t), "Aziza", []OrderItem{{ProductID: 1, Quantity: 2}, {ProductID: 2, Quantity: 5}})
	if err == nil {
		t.Fatalf("PlaceOrder xato qaytarmadi: non qoldigʻi 3 ta, 5 ta soʻralgan")
	}
	if got := count(t, s, "orders"); got != 3 {
		t.Errorf("orders'da %d qator, kutilgan: 3: xato boʻlsa buyurtma yozilib qolmasligi kerak", got)
	}
	if got := count(t, s, "order_items"); got != 4 {
		t.Errorf("order_items'da %d qator, kutilgan: 4: yarim yozilgan buyurtma qoldi", got)
	}
	if got := stockOf(t, s, 1); got != 40 {
		t.Errorf("olma qoldigʻi = %d, kutilgan: 40: birinchi qadam bekor boʻlmadi", got)
	}
}

func TestPlaceOrderUnknownProduct(t *testing.T) {
	s := newStore(t)
	_, err := s.PlaceOrder(testCtx(t), "Aziza", []OrderItem{{ProductID: 1, Quantity: 1}, {ProductID: 99, Quantity: 1}})
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("PlaceOrder(99-mahsulot) = %v; kutilgan: ErrNotFound", err)
	}
	if got := count(t, s, "orders"); got != 3 {
		t.Errorf("orders'da %d qator, kutilgan: 3", got)
	}
}

// VAZIFA, 1-qadam: CHECK (stock >= 0) buzilganini chaqiruvchi ErrNotEnoughStock bilan taniydi.
func TestPlaceOrderNotEnoughStock(t *testing.T) {
	s := newStore(t)
	_, err := s.PlaceOrder(testCtx(t), "Aziza", []OrderItem{{ProductID: 2, Quantity: 5}})
	if !errors.Is(err, ErrNotEnoughStock) {
		t.Fatalf("PlaceOrder(non x5) = %v; kutilgan: ErrNotEnoughStock (TASK.md, 1-qadam: errors.As bilan *pgconn.PgError, Code 23514)", err)
	}
	if got := stockOf(t, s, 2); got != 3 {
		t.Errorf("non qoldigʻi = %d, kutilgan: 3", got)
	}
}

// VAZIFA, 2-qadam.
func TestRestock(t *testing.T) {
	s := newStore(t)
	if err := s.Restock(testCtx(t), []OrderItem{{ProductID: 2, Quantity: 10}, {ProductID: 4, Quantity: 5}}); err != nil {
		t.Fatalf("Restock = %v; kutilgan: nil (TASK.md, 2-qadam)", err)
	}
	if got := stockOf(t, s, 2); got != 13 {
		t.Errorf("non qoldigʻi = %d, kutilgan: 13", got)
	}
	if got := stockOf(t, s, 4); got != 7 {
		t.Errorf("oʻrik qoldigʻi = %d, kutilgan: 7", got)
	}
}

func TestRestockAllOrNothing(t *testing.T) {
	s := newStore(t)
	err := s.Restock(testCtx(t), []OrderItem{{ProductID: 1, Quantity: 10}, {ProductID: 99, Quantity: 5}})
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("Restock(99-mahsulot) = %v; kutilgan: ErrNotFound (TASK.md, 2-qadam)", err)
	}
	if got := stockOf(t, s, 1); got != 40 {
		t.Errorf("olma qoldigʻi = %d, kutilgan: 40: 99-mahsulot topilmasa birinchi qadam ham bekor boʻlishi kerak (bitta transaction)", got)
	}
}

// VAZIFA, 3-qadam: miqdor noldan katta boʻlishi kerak; xato Begin'dan oldin tekshiriladi, hech narsa yozilmaydi.
func TestInvalidQuantity(t *testing.T) {
	s := newStore(t)
	if _, err := s.PlaceOrder(testCtx(t), "Aziza", []OrderItem{{ProductID: 1, Quantity: 1}, {ProductID: 2, Quantity: 0}}); !errors.Is(err, ErrInvalidQuantity) {
		t.Errorf("PlaceOrder(miqdor 0) = %v; kutilgan: ErrInvalidQuantity (TASK.md, 3-qadam)", err)
	}
	if err := s.Restock(testCtx(t), []OrderItem{{ProductID: 2, Quantity: -2}}); !errors.Is(err, ErrInvalidQuantity) {
		t.Errorf("Restock(miqdor -2) = %v; kutilgan: ErrInvalidQuantity (TASK.md, 3-qadam): manfiy miqdor qoldiqni kamaytirib yuboradi", err)
	}
	if got := count(t, s, "orders"); got != 3 {
		t.Errorf("orders'da %d qator, kutilgan: 3", got)
	}
	if got := stockOf(t, s, 1); got != 40 {
		t.Errorf("olma qoldigʻi = %d, kutilgan: 40", got)
	}
	if got := stockOf(t, s, 2); got != 3 {
		t.Errorf("non qoldigʻi = %d, kutilgan: 3", got)
	}
}
