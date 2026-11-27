// Bu fayl vazifani tekshiradi: Store metodlarini haqiqiy Postgres'da ishga tushiradi.
// Uni oʻzgartirmang. Postgres uchun: docker compose up -d --wait (shu papkada).
// Har test alohida vaqtinchalik sxemada ishlaydi (products.sql + seed.sql shu yerga qoʻllanadi)
// va oxirida sxemani oʻchiradi: public sxemaga tegmaydi. Boshqa manzil kerak boʻlsa DATABASE_URL qoʻying.
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"slices"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const labURL = "postgres://bozor:bozor-lab@localhost:5475/bozor"

// newStore bu test uchun boʻsh sxema yaratadi, jadvalni toʻldiradi va shu sxemaga ulangan Store qaytaradi.
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
	cfg.MaxConns = 4 // darsdagi pool_max_conns=1 tajribasidan qolgan DATABASE_URL testni toʻxtatib qoʻymasin
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
	for _, f := range []string{"products.sql", "seed.sql", "orders.sql", "orders_seed.sql"} {
		b, err := os.ReadFile(f)
		if err != nil {
			t.Fatalf("%s oʻqilmadi: %v", f, err)
		}
		// Oddiy protokol: bitta matnda bir nechta SQL buyrugʻi boʻlishi mumkin.
		if _, err := pool.Exec(ctx, string(b), pgx.QueryExecModeSimpleProtocol); err != nil {
			t.Fatalf("%s bajarilmadi: %v", f, err)
		}
	}
	return &Store{pool: pool}
}

func testCtx(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	t.Cleanup(cancel)
	return ctx
}

func names(ps []Product) []string {
	out := make([]string, 0, len(ps))
	for _, p := range ps {
		out = append(out, p.Name)
	}
	return out
}

// Videoda yozilgan metodlar: starterda ham oʻtadi.

func TestOrderTotal(t *testing.T) {
	s := newStore(t)
	if got, err := s.OrderTotal(testCtx(t), 1); err != nil || got != 3_600_000 {
		t.Fatalf("OrderTotal(1) = %d, %v; kutilgan: 3600000, nil", got, err)
	}
	if got, err := s.OrderTotal(testCtx(t), 99); err != nil || got != 0 {
		t.Fatalf("OrderTotal(99) = %d, %v; kutilgan: 0, nil (SUM NULL beradi: COALESCE bilan nolga aylantiring)", got, err)
	}
}

func TestGet(t *testing.T) {
	s := newStore(t)
	p, err := s.Get(testCtx(t), 1)
	if err != nil || p != (Product{ID: 1, Name: "olma", Price: 1_200_000, Stock: 40}) {
		t.Fatalf("Get(1) = %+v, %v; kutilgan: {1 olma 1200000 40}, nil", p, err)
	}
	if _, err := s.Get(testCtx(t), 99); !errors.Is(err, ErrNotFound) {
		t.Fatalf("Get(99): xato %v; kutilgan ErrNotFound (pgx.ErrNoRows'ni ErrNotFound'ga aylantiring)", err)
	}
}

func TestListActive(t *testing.T) {
	s := newStore(t)
	if _, err := s.pool.Exec(testCtx(t), "UPDATE products SET active = false WHERE name = 'non'"); err != nil {
		t.Fatal(err)
	}
	ps, err := s.ListActive(testCtx(t))
	if err != nil || !slices.Equal(names(ps), []string{"olma", "choy", "oʻrik"}) {
		t.Fatalf("ListActive() = %v, %v; kutilgan: [olma choy oʻrik] (non faol emas), nil", names(ps), err)
	}
}

// VAZIFA, 1-qadam: ByName.

func TestByName(t *testing.T) {
	s := newStore(t)
	p, err := s.ByName(testCtx(t), "choy")
	if err != nil || p != (Product{ID: 3, Name: "choy", Price: 2_500_000, Stock: 8}) {
		t.Fatalf("ByName(\"choy\") = %+v, %v; kutilgan: {3 choy 2500000 8}, nil (TASK.md, 1-qadam)", p, err)
	}
	for _, name := range []string{"yoʻq", "", "x' OR 'a'='a", "olma'; DROP TABLE products; --"} {
		if _, err := s.ByName(testCtx(t), name); !errors.Is(err, ErrNotFound) {
			t.Fatalf("ByName(%q): xato %v; kutilgan ErrNotFound. Nomni SQL matniga qoʻshmang: $1 parametrini ishlating (TASK.md, 1-qadam)", name, err)
		}
	}
	var n int
	if err := s.pool.QueryRow(testCtx(t), "SELECT count(*) FROM products").Scan(&n); err != nil || n != 4 {
		t.Fatalf("products jadvalida %d qator, kutilgan 4 (err: %v): nom SQL buyrugʻiga aylanib ketgan", n, err)
	}
	done, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := s.ByName(done, "olma"); err == nil || errors.Is(err, ErrNotFound) {
		t.Fatalf("bekor qilingan ctx bilan ByName: xato %v; kutilgan context xatosi (ctx'ni QueryRow'ga bering)", err)
	}
}

// VAZIFA, 2-qadam: Cheapest.

func TestCheapest(t *testing.T) {
	s := newStore(t)
	ps, err := s.Cheapest(testCtx(t), 2)
	if err != nil || !slices.Equal(names(ps), []string{"non", "olma"}) {
		t.Fatalf("Cheapest(2) = %v, %v; kutilgan: [non olma], nil (TASK.md, 2-qadam)", names(ps), err)
	}
	if ps[0].Price != 400_000 || ps[0].Stock != 3 {
		t.Fatalf("Cheapest(2)[0] = %+v; kutilgan narx 400000, qoldiq 3", ps[0])
	}
	if _, err := s.pool.Exec(testCtx(t), "UPDATE products SET active = false WHERE name = 'non'"); err != nil {
		t.Fatal(err)
	}
	ps, err = s.Cheapest(testCtx(t), 3)
	if err != nil || !slices.Equal(names(ps), []string{"olma", "choy", "oʻrik"}) {
		t.Fatalf("non faol emas: Cheapest(3) = %v, %v; kutilgan: [olma choy oʻrik] (faqat faol, narx oʻsish tartibida)", names(ps), err)
	}
	ps, err = s.Cheapest(testCtx(t), 10)
	if err != nil || len(ps) != 3 {
		t.Fatalf("Cheapest(10): %d ta mahsulot, %v; kutilgan 3 (faol mahsulotlarning hammasi)", len(ps), err)
	}
	if ps, err = s.Cheapest(testCtx(t), 0); err != nil || len(ps) != 0 || ps == nil {
		t.Fatalf("Cheapest(0) = %#v, %v; kutilgan: boʻsh, nil boʻlmagan slice ([]Product{}), xatosiz", ps, err)
	}
}
