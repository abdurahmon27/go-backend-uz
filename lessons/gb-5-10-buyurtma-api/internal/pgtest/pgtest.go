// Package pgtest — testlar uchun haqiqiy Postgres (mock emas). Har test oʻz boʻsh sxemasini oladi:
// migratsiyalar (goose) va seed.sql shu sxemaga qoʻllanadi, test tugagach sxema oʻchiriladi.
// Shuning uchun testlar bir-biriga xalaqit bermaydi, `go test ./...` paketlarni parallel
// ishga tushirganda ham (GB-5.10 videosidagi xato), va sizning public sxemangizga tegmaydi.
package pgtest

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/abdurahmon27/bozor"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

// LabURL — shu papkadagi compose.yaml koʻtaradigan Postgres. Boshqa manzil uchun DATABASE_URL qoʻying.
const LabURL = "postgres://bozor:bozor-lab@localhost:5510/bozor"

// New bu test uchun yangi sxema yaratadi, unga migratsiyalar va seed.sql'ni qoʻllaydi va shu sxemaga
// ulangan pool qaytaradi. Postgres ishlamasa test aniq xabar bilan yiqiladi (oʻtkazib yuborilmaydi).
func New(t *testing.T) *pgxpool.Pool {
	t.Helper()
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		url = LabURL
	}
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		t.Fatalf("DATABASE_URL notoʻgʻri: %v", err)
	}
	cfg.MaxConns = max(cfg.MaxConns, 4) // URL'da pool_max_conns=1 qolgan boʻlsa ham parallel test toʻxtab qolmasin (GB-5.07)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	admin, err := pgxpool.NewWithConfig(ctx, cfg)
	if err == nil {
		err = admin.Ping(ctx)
	}
	if err != nil {
		t.Fatalf("Postgres'ga ulanib boʻlmadi (%s:%d): avval shu papkada `docker compose up -d --wait` ishga tushiring.\n%v",
			cfg.ConnConfig.Host, cfg.ConnConfig.Port, err)
	}
	schema := fmt.Sprintf("test_%d", time.Now().UnixNano())
	if _, err := admin.Exec(ctx, "CREATE SCHEMA "+schema); err != nil {
		t.Fatalf("sxema yaratilmadi: %v", err)
	}
	t.Cleanup(func() {
		if _, err := admin.Exec(context.Background(), "DROP SCHEMA "+schema+" CASCADE"); err != nil {
			t.Logf("sxema %s oʻchirilmadi (qoʻlda: DROP SCHEMA %s CASCADE): %v", schema, schema, err)
		}
		admin.Close()
	})

	poolCfg := cfg.Copy()                                    // nusxa: admin pool'ning sozlamasi oʻzgarmasin (u ham shu cfg'dan yasalgan)
	poolCfg.ConnConfig.RuntimeParams["search_path"] = schema // har ulanish jadvallarni shu sxemadan qidiradi
	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		t.Fatalf("pool: %v", err)
	}
	t.Cleanup(pool.Close)

	db := stdlib.OpenDBFromPool(pool) // goose database/sql bilan ishlaydi: pool ustidan *sql.DB
	defer db.Close()
	provider, err := bozor.NewProvider(db)
	if err != nil {
		t.Fatalf("goose: %v", err)
	}
	if _, err := provider.Up(ctx); err != nil {
		t.Fatalf("migratsiya qoʻllanmadi: %v", err)
	}
	// Oddiy protokol: bitta matnda bir nechta SQL buyrugʻi boʻlishi mumkin (GB-5.07 testlari).
	if _, err := pool.Exec(ctx, bozor.SeedSQL, pgx.QueryExecModeSimpleProtocol); err != nil {
		t.Fatalf("seed.sql bajarilmadi: %v", err)
	}
	return pool
}
