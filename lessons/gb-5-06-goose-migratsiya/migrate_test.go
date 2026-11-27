// Bu fayl vazifani tekshiradi: migrations/ papkasidagi fayllarni haqiqiy Postgres'da goose bilan qoʻllaydi.
// Uni oʻzgartirmang. Postgres uchun: docker compose up -d --wait (shu papkada).
// Har test alohida vaqtinchalik sxemada ishlaydi va oxirida sxemani oʻchiradi: public sxemaga va
// sizning goose_db_version jadvalingizga tegmaydi. Boshqa manzil kerak boʻlsa DATABASE_URL qoʻying.
package bozor

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

const labURL = "postgres://bozor:bozor-lab@localhost:5486/bozor"

// newDB bu test uchun boʻsh sxema yaratadi va shu sxemaga ulangan *sql.DB qaytaradi.
func newDB(t *testing.T) (*sql.DB, string) {
	t.Helper()
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		url = labURL
	}
	cfg, err := pgx.ParseConfig(url)
	if err != nil {
		t.Fatalf("DATABASE_URL notoʻgʻri: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	admin := stdlib.OpenDB(*cfg)
	if err := admin.PingContext(ctx); err != nil {
		admin.Close()
		t.Fatalf("Postgres'ga ulanib boʻlmadi (%s:%d): avval shu papkada `docker compose up -d --wait` ishga tushiring.\n%v", cfg.Host, cfg.Port, err)
	}
	schema := fmt.Sprintf("lab_%d", time.Now().UnixNano())
	if _, err := admin.ExecContext(ctx, "CREATE SCHEMA "+schema); err != nil {
		admin.Close()
		t.Fatalf("sxema yaratilmadi: %v", err)
	}
	t.Cleanup(func() {
		if _, err := admin.ExecContext(context.Background(), "DROP SCHEMA "+schema+" CASCADE"); err != nil {
			t.Logf("sxema %s oʻchirilmadi (qoʻlda: DROP SCHEMA %s CASCADE): %v", schema, schema, err)
		}
		admin.Close()
	})

	labCfg := cfg.Copy()
	labCfg.RuntimeParams["search_path"] = schema // jadvallar, goose_db_version ham, shu sxemada yaratiladi
	db := stdlib.OpenDB(*labCfg)
	t.Cleanup(func() { db.Close() })
	return db, schema
}

func testCtx(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	t.Cleanup(cancel)
	return ctx
}

func provider(t *testing.T, db *sql.DB) *goose.Provider {
	t.Helper()
	p, err := NewProvider(db)
	if err != nil {
		t.Fatalf("goose provider: %v", err)
	}
	return p
}

func count(t *testing.T, db *sql.DB, query string, args ...any) int {
	t.Helper()
	var n int
	if err := db.QueryRowContext(testCtx(t), query, args...).Scan(&n); err != nil {
		t.Fatalf("%s: %v", query, err)
	}
	return n
}

// Videoda yozilgan migratsiyalar (1-3): starterda ham oʻtadi.

// Qoʻllangan migratsiya faylini oʻzgartirmang: yangi raqamli fayl yarating (videodagi qoida).
func TestAppliedMigrationsUnchanged(t *testing.T) {
	want := map[string]string{
		"migrations/00001_products.sql":      "ff05415d55b55a5ad7959d09a947c37b3b7dbcb5eabda663f9902ea9a3c6cbb7",
		"migrations/00002_orders.sql":        "8472674e19f3b1c35f8ef86804a2c65d76c1fe7c6a8060d605594005874268f0",
		"migrations/00003_orders_status.sql": "c359804e6a6296192be8aeba59549102708b39688ebc38c04435f181a6feb1f4",
	}
	for name, sum := range want {
		b, err := fs.ReadFile(files, name)
		if err != nil {
			t.Fatalf("%s oʻqilmadi: %v (videodagi migratsiya fayllarini oʻchirmang)", name, err)
		}
		h := sha256.Sum256([]byte(strings.ReplaceAll(string(b), "\r\n", "\n")))
		if got := hex.EncodeToString(h[:]); got != sum {
			t.Fatalf("%s oʻzgartirilgan. Qoʻllangan migratsiyani oʻzgartirmang: yangi 00004_ fayl yarating (TASK.md)", name)
		}
	}
}

func TestUpDownUp(t *testing.T) {
	db, schema := newDB(t)
	p := provider(t, db)
	ctx := testCtx(t)

	if _, err := p.UpTo(ctx, 3); err != nil {
		t.Fatalf("1-3 migratsiyalar toza bazaga qoʻllanmadi: %v", err)
	}
	for _, table := range []string{"products", "orders", "order_items"} {
		if n := count(t, db, "SELECT count(*) FROM information_schema.tables WHERE table_schema = $1 AND table_name = $2", schema, table); n != 1 {
			t.Fatalf("up'dan keyin %s jadvali yoʻq", table)
		}
	}
	if _, err := p.DownTo(ctx, 0); err != nil {
		t.Fatalf("hammasini qaytarish (down) xato berdi: %v", err)
	}
	if n := count(t, db, "SELECT count(*) FROM information_schema.tables WHERE table_schema = $1 AND table_name <> 'goose_db_version'", schema); n != 0 {
		t.Fatalf("hamma down'dan keyin %d ta jadval qoldi; kutilgan 0 (Down boʻlimi Up'ni toʻliq qaytarishi kerak)", n)
	}
	if _, err := p.Up(ctx); err != nil {
		t.Fatalf("down'dan keyin qayta up xato berdi: %v", err)
	}
}

func TestOrdersStatusMigration(t *testing.T) {
	db, _ := newDB(t)
	p := provider(t, db)
	ctx := testCtx(t)
	if _, err := p.UpTo(ctx, 2); err != nil {
		t.Fatal(err)
	}
	if _, err := db.ExecContext(ctx, "INSERT INTO orders (customer) VALUES ('Dilnoza')"); err != nil {
		t.Fatal(err)
	}
	if _, err := p.UpTo(ctx, 3); err != nil {
		t.Fatalf("3-migratsiya (orders.status) qatorlari bor jadvalga qoʻllanmadi: %v", err)
	}
	if n := count(t, db, "SELECT count(*) FROM orders WHERE status = 'yangi'"); n != 1 {
		t.Fatalf("mavjud buyurtma status = 'yangi' boʻlishi kerak edi (DEFAULT), topildi: %d", n)
	}
	if _, err := db.ExecContext(ctx, "INSERT INTO orders (customer, status) VALUES ('Jasur', '')"); err == nil {
		t.Fatalf("status = '' qabul qilindi; kutilgan: CHECK (status <> '') rad etsin")
	}
}

// VAZIFA: migrations/00004_products_category.sql.

func TestCategoryMigration(t *testing.T) {
	db, schema := newDB(t)
	p := provider(t, db)
	ctx := testCtx(t)

	has4 := false
	for _, s := range p.ListSources() {
		if s.Version == 4 {
			has4 = true
		}
	}
	if !has4 {
		t.Fatalf("4-migratsiya yoʻq: migrations/ papkasida 00004_ bilan boshlanadigan .sql fayl yarating (TASK.md, 1-qadam)")
	}

	if _, err := p.UpTo(ctx, 3); err != nil {
		t.Fatal(err)
	}
	if n := count(t, db, "SELECT count(*) FROM information_schema.columns WHERE table_schema = $1 AND table_name = 'products' AND column_name = 'category'", schema); n != 0 {
		t.Fatalf("products.category 3-versiyada ham bor: ustunni faqat 00004 migratsiya qoʻshsin (qoʻllangan fayllarni oʻzgartirmang)")
	}
	if _, err := db.ExecContext(ctx, "INSERT INTO products (name, price, stock) VALUES ('olma', 1200000, 40), ('non', 400000, 3)"); err != nil {
		t.Fatal(err)
	}
	if _, err := p.UpTo(ctx, 4); err != nil {
		hint := "TASK.md, 2-qadam"
		if strings.Contains(err.Error(), "23502") {
			hint = "qatorlari bor jadvalga DEFAULT'siz NOT NULL ustun qoʻshilmaydi: DEFAULT 'boshqa' yozing (TASK.md, 2-qadam)"
		}
		t.Fatalf("4-migratsiya qatorlari bor products jadvaliga qoʻllanmadi: %v\n(%s)", err, hint)
	}

	var typ, nullable string
	err := db.QueryRowContext(ctx, "SELECT data_type, is_nullable FROM information_schema.columns WHERE table_schema = $1 AND table_name = 'products' AND column_name = 'category'", schema).Scan(&typ, &nullable)
	if err != nil {
		t.Fatalf("products.category ustuni topilmadi (%v); kutilgan: text, NOT NULL (TASK.md, 2-qadam)", err)
	}
	if typ != "text" || nullable != "NO" {
		t.Fatalf("products.category: tur %q, NULL mumkin: %q; kutilgan: text, NOT NULL", typ, nullable)
	}
	if n := count(t, db, "SELECT count(*) FROM products WHERE category = 'boshqa'"); n != 2 {
		t.Fatalf("mavjud 2 ta mahsulotdan %d tasida category = 'boshqa'; kutilgan 2 (DEFAULT 'boshqa')", n)
	}
	if _, err := db.ExecContext(ctx, "INSERT INTO products (name, price, category) VALUES ('choy', 2500000, 'ichimlik')"); err != nil {
		t.Fatalf("category bilan mahsulot qoʻshilmadi: %v", err)
	}
	_, err = db.ExecContext(ctx, "INSERT INTO products (name, price, category) VALUES ('oʻrik', 3500000, '')")
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) || pgErr.Code != "23514" {
		t.Fatalf("category = '' uchun CHECK buzilishi (23514) kutilgan edi, keldi: %v (CHECK (category <> ''), TASK.md, 2-qadam)", err)
	}

	if _, err := p.Down(ctx); err != nil {
		t.Fatalf("4-migratsiya Down'i xato berdi: %v (TASK.md, 3-qadam)", err)
	}
	if n := count(t, db, "SELECT count(*) FROM information_schema.columns WHERE table_schema = $1 AND table_name = 'products' AND column_name = 'category'", schema); n != 0 {
		t.Fatalf("Down'dan keyin products.category ustuni qoldi; kutilgan: ustun oʻchirilsin (TASK.md, 3-qadam)")
	}
	if v, err := p.GetDBVersion(ctx); err != nil || v != 3 {
		t.Fatalf("Down'dan keyin versiya %d (%v); kutilgan 3", v, err)
	}
	if _, err := p.Up(ctx); err != nil {
		t.Fatalf("Down'dan keyin qayta Up xato berdi: %v", err)
	}
}
