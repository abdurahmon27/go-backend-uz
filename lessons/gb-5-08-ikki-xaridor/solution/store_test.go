// Bu fayl vazifani tekshiradi: Store metodlarini haqiqiy Postgres'da, bir vaqtda koʻp xaridor bilan ishga tushiradi.
// Uni oʻzgartirmang. Postgres uchun: docker compose up -d --wait (shu papkada).
// Har test alohida vaqtinchalik sxemada ishlaydi (bozor.sql shu yerga qoʻllanadi) va oxirida sxemani oʻchiradi:
// public sxemaga tegmaydi. Boshqa manzil kerak boʻlsa DATABASE_URL qoʻying.
package bozor

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const labURL = "postgres://bozor:bozor-lab@localhost:5458/bozor"

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
	cfg.MaxConns = 40 // 20 xaridor bir vaqtda; videodagi xato versiya (README) har biri uchun ikkinchi ulanish ham soʻraydi. DATABASE_URL da pool_max_conns qolgan boʻlsa ham
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
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	t.Cleanup(cancel)
	return ctx
}

func stockOf(t *testing.T, s *Store, id int64) int {
	t.Helper()
	var n int
	if err := s.pool.QueryRow(testCtx(t), "SELECT stock FROM products WHERE id = $1", id).Scan(&n); err != nil {
		t.Fatalf("mahsulot %d qoldigʻi oʻqilmadi: %v", id, err)
	}
	return n
}

func setStock(t *testing.T, s *Store, id int64, n int) {
	t.Helper()
	if _, err := s.pool.Exec(testCtx(t), "UPDATE products SET stock = $1 WHERE id = $2", n, id); err != nil {
		t.Fatalf("qoldiq oʻrnatilmadi: %v", err)
	}
}

// runBuyers n ta xaridorni bir vaqtda ishga tushiradi (har biri id 4 dan bittadan oladi) va nechtasiga "olindi" deganini qaytaradi.
func runBuyers(t *testing.T, ctx context.Context, buy func(context.Context, int64, int) error, n int) int {
	t.Helper()
	errs := make([]error, n)
	var wg sync.WaitGroup
	for i := range errs {
		wg.Go(func() {
			errs[i] = buy(ctx, 4, 1)
		})
	}
	wg.Wait()
	ok := 0
	for _, err := range errs {
		switch {
		case err == nil:
			ok++
		case errors.Is(err, ErrNotEnoughStock): // rad etilgan xaridor: kutilgan holat
		default:
			t.Fatalf("kutilmagan xato: %v", err)
		}
	}
	return ok
}

// buyers: sotish usullari. Ikkalasi ham bir xil shartni bajarishi kerak.
var buyers = []struct {
	name string
	buy  func(*Store, context.Context, int64, int) error
}{
	{"Buy", (*Store).Buy},
	{"BuyLocked", (*Store).BuyLocked},
}

// Boshlangʻich holat (bozor.sql): olma(1) 40, non(2) 3, choy(3) 8, oʻrik(4) 2.

func TestBuyOneByOne(t *testing.T) {
	for _, b := range buyers {
		t.Run(b.name, func(t *testing.T) {
			s := newStore(t)
			if err := b.buy(s, testCtx(t), 1, 2); err != nil {
				t.Fatalf("%s(olma, 2) xato qaytardi: %v; kutilgan: nil", b.name, err)
			}
			if got := stockOf(t, s, 1); got != 38 {
				t.Errorf("olma qoldigʻi = %d, kutilgan: 38", got)
			}
			err := b.buy(s, testCtx(t), 2, 4) // non: 3 ta bor
			if !errors.Is(err, ErrNotEnoughStock) {
				t.Errorf("%s(non, 4) = %v; kutilgan: ErrNotEnoughStock (non 3 ta)", b.name, err)
			}
			if got := stockOf(t, s, 2); got != 3 {
				t.Errorf("non qoldigʻi = %d, kutilgan: 3 (rad etilgan xarid qoldiqni oʻzgartirmasligi kerak)", got)
			}
			for _, bad := range []int{0, -2} {
				if err := b.buy(s, testCtx(t), 4, bad); !errors.Is(err, ErrInvalidQuantity) {
					t.Errorf("%s(oʻrik, %d) = %v; kutilgan: ErrInvalidQuantity", b.name, bad, err)
				}
			}
			if got := stockOf(t, s, 4); got != 2 {
				t.Errorf("oʻrik qoldigʻi = %d, kutilgan: 2 (notoʻgʻri miqdor qoldiqni oʻzgartirmasligi kerak)", got)
			}
		})
	}
}

// 20 xaridor bir vaqtda 1 donadan oladi, omborda 5 ta: roppa-rosa 5 tasi olishi kerak. Tekshiruv 5 marta takrorlanadi:
// xato usul ba'zan tasodifan oʻtib ketmasin.
func TestConcurrentBuyers(t *testing.T) {
	for _, b := range buyers {
		t.Run(b.name, func(t *testing.T) {
			s := newStore(t)
			ctx := testCtx(t)
			buy := func(ctx context.Context, id int64, qty int) error { return b.buy(s, ctx, id, qty) }
			for round := 1; round <= 5; round++ {
				setStock(t, s, 4, 5)
				ok := runBuyers(t, ctx, buy, 20)
				if ok != 5 {
					t.Fatalf("%d-urinish: 20 xaridor, omborda 5 ta: %d tasiga \"olindi\" dedik; kutilgan: 5", round, ok)
				}
				if got := stockOf(t, s, 4); got != 0 {
					t.Fatalf("%d-urinish: oxirgi qoldiq = %d, kutilgan: 0", round, got)
				}
			}
		})
	}
}
