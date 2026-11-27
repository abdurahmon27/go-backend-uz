// Bu fayl vazifani tekshiradi: query.sql, sqlc generatsiya qilgan db/query.sql.go va store.go'ni oʻqiydi,
// keyin Store metodlarini haqiqiy Postgres'da ishga tushiradi. Uni oʻzgartirmang.
// Postgres uchun: docker compose up -d --wait (shu papkada). Har test alohida vaqtinchalik sxemada ishlaydi:
// migrations/*.sql fayllarining «-- +goose Up» qismi va seed.sql shu sxemaga qoʻllanadi, oxirida sxema oʻchiriladi
// (public sxemaga tegmaydi). Boshqa manzil kerak boʻlsa DATABASE_URL qoʻying.
package main

import (
	"context"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const labURL = "postgres://bozor:bozor-lab@localhost:5409/bozor"

// ---------------------------------------------------------------- fayllar: query.sql, db/, store.go

// TestQueriesMovedToSqlc: vazifadagi uchta soʻrov query.sql'da, sqlc ularni generatsiya qilgan,
// va ByName, Cheapest, OrderTotal metodlari (qaysi .go faylda boʻlmasin) s.q orqali ishlaydi: ularda Scan ham,
// pool.Query/QueryRow ham yoʻq. Kod go/parser bilan oʻqiladi, izohlar hisobga olinmaydi.
func TestQueriesMovedToSqlc(t *testing.T) {
	query := readFile(t, "query.sql")
	gen := readFile(t, filepath.Join("db", "query.sql.go"))
	methods := storeMethods(t)

	want := []struct {
		method, query, kind, signature, hint string
	}{
		{"ByName", "ProductByName", ":one", "ProductByName(ctx context.Context, name string) (ProductByNameRow, error)",
			"SELECT … WHERE name = $1"},
		{"Cheapest", "CheapestProducts", ":many", "CheapestProducts(ctx context.Context, limit int32) ([]CheapestProductsRow, error)",
			"SELECT … WHERE active ORDER BY price, id LIMIT $1"},
		{"OrderTotal", "OrderTotal", ":one", "OrderTotal(ctx context.Context, orderID int64) (int64, error)",
			"TASK.md, 3-maslahat: hisoblangan ustun ::bigint bilan"},
	}
	for _, w := range want {
		annotation := regexp.MustCompile(`(?m)^-- name: ` + w.query + `\s+` + w.kind + `\s*$`)
		switch {
		case !annotation.MatchString(query):
			t.Errorf("query.sql'da «-- name: %s %s» soʻrovi yoʻq (TASK.md, 1-qadam)", w.query, w.kind)
		case !strings.Contains(gen, "func (q *Queries) "+w.query+"("):
			t.Errorf("query.sql'da %s bor, lekin db/query.sql.go'da yoʻq: sqlc generate'ni ishga tushiring (TASK.md, 2-qadam)", w.query)
		case !strings.Contains(gen, "func (q *Queries) "+w.signature):
			t.Errorf("db/query.sql.go'dagi %s kutilgandek emas; kutilgan: func (q *Queries) %s (%s)", w.query, w.signature, w.hint)
		}
		body, ok := methods[w.method]
		if !ok {
			t.Errorf("Store'da %s metodi topilmadi", w.method)
			continue
		}
		if bad := handSQL(body); len(bad) > 0 {
			t.Errorf("%s hali qoʻlda yozilgan (%s): s.q.%s orqali yozing (TASK.md, 3-qadam)", w.method, strings.Join(bad, ", "), w.query)
		} else if !callsQueries(body, w.query) {
			t.Errorf("%s s.q.%s'ni chaqirmaydi (TASK.md, 3-qadam)", w.method, w.query)
		}
	}
}

// TestRestockUsesSqlc — 🟡 ixtiyoriy: Restock hali qoʻlda boʻlsa xabar bilan oʻtkazib yuboriladi (u baribir TestRestock'da tekshiriladi).
func TestRestockUsesSqlc(t *testing.T) {
	body, ok := storeMethods(t)["Restock"]
	if !ok {
		t.Fatal("Store'da Restock metodi topilmadi")
	}
	if !callsQueries(body, "AddStock") {
		t.Skip("🟡 Restock hali qoʻlda (ixtiyoriy, TASK.md 4-qadam: AddStock :execrows)")
	}
	if !regexp.MustCompile(`(?m)^-- name: AddStock\s+:execrows\s*$`).MatchString(readFile(t, "query.sql")) {
		t.Errorf("query.sql'da «-- name: AddStock :execrows» yoʻq")
	}
	if bad := handSQL(body); len(bad) > 0 {
		t.Errorf("Restock'da hali qoʻlda yozilgan SQL bor (%s)", strings.Join(bad, ", "))
	}
}

// storeMethods papkadagi test boʻlmagan .go fayllardan *Store metodlarining tanasini oladi.
func storeMethods(t *testing.T) map[string]*ast.BlockStmt {
	t.Helper()
	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	out := map[string]*ast.BlockStmt{}
	fset := token.NewFileSet()
	for _, f := range files {
		if strings.HasSuffix(f, "_test.go") {
			continue
		}
		file, err := parser.ParseFile(fset, f, nil, 0)
		if err != nil {
			t.Fatalf("%s oʻqilmadi: %v", f, err)
		}
		for _, d := range file.Decls {
			fn, ok := d.(*ast.FuncDecl)
			if !ok || fn.Recv == nil || len(fn.Recv.List) != 1 || fn.Body == nil {
				continue
			}
			if star, ok := fn.Recv.List[0].Type.(*ast.StarExpr); ok {
				if id, ok := star.X.(*ast.Ident); ok && id.Name == "Store" {
					out[fn.Name.Name] = fn.Body
				}
			}
		}
	}
	return out
}

// handSQL tanadagi qoʻlda yozilgan SQL belgilarini qaytaradi: Scan chaqiruvi yoki pool.Query/QueryRow/Exec.
func handSQL(body *ast.BlockStmt) []string {
	var bad []string
	ast.Inspect(body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		switch name := sel.Sel.Name; {
		case name == "Scan":
			bad = append(bad, "Scan")
		case name == "Query" || name == "QueryRow" || name == "Exec":
			if inner, ok := sel.X.(*ast.SelectorExpr); ok && inner.Sel.Name == "pool" {
				bad = append(bad, "pool."+name)
			}
			if id, ok := sel.X.(*ast.Ident); ok && id.Name == "tx" {
				bad = append(bad, "tx."+name)
			}
		}
		return true
	})
	return bad
}

// callsQueries tana sqlc metodini (s.q.Name yoki qtx.Name) chaqiradimi.
func callsQueries(body *ast.BlockStmt, name string) bool {
	found := false
	ast.Inspect(body, func(n ast.Node) bool {
		if sel, ok := n.(*ast.SelectorExpr); ok && sel.Sel.Name == name {
			found = true
		}
		return !found
	})
	return found
}

// ---------------------------------------------------------------- Postgres

func TestGet(t *testing.T) {
	s := newStore(t)
	p, err := s.Get(testCtx(t), 1)
	if err != nil {
		t.Fatalf("Get(1): kutilmagan xato: %v", err)
	}
	if want := (Product{ID: 1, Name: "olma", Price: 1200000, Stock: 40}); p != want {
		t.Errorf("Get(1) = %+v; kutilgan %+v", p, want)
	}
	if _, err := s.Get(testCtx(t), 99); !errors.Is(err, ErrNotFound) {
		t.Errorf("Get(99) xatosi = %v; kutilgan ErrNotFound", err)
	}
}

func TestListActive(t *testing.T) {
	s := newStore(t)
	exec(t, s, "UPDATE products SET active = false WHERE name = 'choy'")
	got, err := s.ListActive(testCtx(t))
	if err != nil {
		t.Fatalf("ListActive: %v", err)
	}
	if names := namesOf(got); names != "olma non oʻrik" {
		t.Errorf("ListActive nomlari = %q; kutilgan \"olma non oʻrik\"", names)
	}
}

func TestPlaceOrder(t *testing.T) {
	s := newStore(t)
	id, err := s.PlaceOrder(testCtx(t), "Aziza", []OrderItem{{ProductID: 1, Quantity: 2}, {ProductID: 2, Quantity: 1}})
	if err != nil {
		t.Fatalf("PlaceOrder: kutilmagan xato: %v", err)
	}
	if id != 4 {
		t.Errorf("PlaceOrder id = %d; kutilgan 4 (seed.sql'da 3 ta buyurtma bor)", id)
	}
	if got := stockOf(t, s, 1); got != 38 {
		t.Errorf("olma qoldigʻi = %d; kutilgan 38", got)
	}
	_, err = s.PlaceOrder(testCtx(t), "Botir", []OrderItem{{ProductID: 1, Quantity: 1}, {ProductID: 2, Quantity: 5}})
	if !errors.Is(err, ErrNotEnoughStock) {
		t.Errorf("PlaceOrder(non x5) xatosi = %v; kutilgan ErrNotEnoughStock", err)
	}
	if got := stockOf(t, s, 1); got != 38 {
		t.Errorf("xato buyurtmadan keyin olma qoldigʻi = %d; kutilgan 38 (transaction bekor boʻlishi kerak)", got)
	}
	if _, err := s.PlaceOrder(testCtx(t), "Botir", []OrderItem{{ProductID: 99, Quantity: 1}}); !errors.Is(err, ErrNotFound) {
		t.Errorf("PlaceOrder(99-mahsulot) xatosi = %v; kutilgan ErrNotFound", err)
	}
	if _, err := s.PlaceOrder(testCtx(t), "Botir", []OrderItem{{ProductID: 1, Quantity: 0}}); !errors.Is(err, ErrInvalidQuantity) {
		t.Errorf("PlaceOrder(miqdor 0) xatosi = %v; kutilgan ErrInvalidQuantity", err)
	}
}

func TestByName(t *testing.T) {
	s := newStore(t)
	p, err := s.ByName(testCtx(t), "oʻrik")
	if err != nil {
		t.Fatalf("ByName(oʻrik): kutilmagan xato: %v", err)
	}
	if want := (Product{ID: 4, Name: "oʻrik", Price: 3500000, Stock: 2}); p != want {
		t.Errorf("ByName(oʻrik) = %+v; kutilgan %+v", p, want)
	}
	if _, err := s.ByName(testCtx(t), "anor"); !errors.Is(err, ErrNotFound) {
		t.Errorf("ByName(anor) xatosi = %v; kutilgan ErrNotFound (pgx.ErrNoRows → ErrNotFound)", err)
	}
}

func TestCheapest(t *testing.T) {
	s := newStore(t)
	exec(t, s, "UPDATE products SET active = false WHERE name = 'non'")
	got, err := s.Cheapest(testCtx(t), 2)
	if err != nil {
		t.Fatalf("Cheapest(2): %v", err)
	}
	if names := namesOf(got); names != "olma choy" {
		t.Errorf("Cheapest(2) nomlari = %q; kutilgan \"olma choy\" (faqat faol, narx boʻyicha)", names)
	}
	exec(t, s, "UPDATE products SET active = false")
	got, err = s.Cheapest(testCtx(t), 2)
	if err != nil {
		t.Fatalf("Cheapest(2), faol mahsulot yoʻq: %v", err)
	}
	if got == nil || len(got) != 0 {
		t.Errorf("Cheapest(2), faol mahsulot yoʻq = %#v; kutilgan boʻsh slice []Product{} (nil emas: JSON'da null emas [] boʻlsin)", got)
	}
}

func TestOrderTotal(t *testing.T) {
	s := newStore(t)
	cases := []struct {
		order int64
		want  int64
	}{
		{1, 3600000}, // 2 olma × 1 200 000 + 3 non × 400 000
		{2, 4700000}, // 1 oʻrik × 3 500 000 + 1 olma × 1 200 000
		{3, 0},       // qatorsiz buyurtma: SUM NULL, COALESCE → 0
	}
	for _, c := range cases {
		got, err := s.OrderTotal(testCtx(t), c.order)
		if err != nil {
			t.Fatalf("OrderTotal(%d): %v", c.order, err)
		}
		if got != c.want {
			t.Errorf("OrderTotal(%d) = %d; kutilgan %d", c.order, got, c.want)
		}
	}
}

// TestRestock: Restock (GB-5.07) qoʻlda ham, sqlc bilan ham bir xil ishlashi kerak.
func TestRestock(t *testing.T) {
	s := newStore(t)
	err := s.Restock(testCtx(t), []OrderItem{{ProductID: 2, Quantity: 5}, {ProductID: 4, Quantity: 1}})
	if err != nil {
		t.Fatalf("Restock: kutilmagan xato: %v", err)
	}
	if a, b := stockOf(t, s, 2), stockOf(t, s, 4); a != 8 || b != 3 {
		t.Errorf("Restock'dan keyin non = %d, oʻrik = %d; kutilgan 8 va 3", a, b)
	}
	err = s.Restock(testCtx(t), []OrderItem{{ProductID: 1, Quantity: 1}, {ProductID: 99, Quantity: 1}})
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("Restock(99-mahsulot) xatosi = %v; kutilgan ErrNotFound (RowsAffected 0)", err)
	}
	if got := stockOf(t, s, 1); got != 40 {
		t.Errorf("Restock(99) dan keyin olma = %d; kutilgan 40 (hammasi bitta transaction'da bekor boʻlishi kerak)", got)
	}
	if err := s.Restock(testCtx(t), []OrderItem{{ProductID: 1, Quantity: -2}}); !errors.Is(err, ErrInvalidQuantity) {
		t.Errorf("Restock(miqdor -2) xatosi = %v; kutilgan ErrInvalidQuantity", err)
	}
}

// ---------------------------------------------------------------- yordamchilar

func readFile(t *testing.T, name string) string {
	t.Helper()
	b, err := os.ReadFile(name)
	if err != nil {
		t.Fatalf("%s oʻqilmadi: %v", name, err)
	}
	return string(b)
}

// gooseUp migratsiya faylining «-- +goose Up» va «-- +goose Down» orasidagi qismini qaytaradi.
func gooseUp(t *testing.T, path string) string {
	t.Helper()
	s := readFile(t, path)
	_, up, ok := strings.Cut(s, "-- +goose Up")
	if !ok {
		t.Fatalf("%s: «-- +goose Up» qatori yoʻq", path)
	}
	up, _, _ = strings.Cut(up, "-- +goose Down")
	return up
}

// newStore bu test uchun boʻsh sxema yaratadi, migratsiyalarni va seed.sql'ni qoʻllaydi, shu sxemaga ulangan Store qaytaradi.
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

	files, err := filepath.Glob(filepath.Join("migrations", "*.sql"))
	if err != nil || len(files) == 0 {
		t.Fatalf("migrations/*.sql topilmadi: %v", err)
	}
	var script strings.Builder
	for _, f := range files { // Glob nomlarni alifbo tartibida beradi: 00001, 00002, …
		script.WriteString(gooseUp(t, f))
	}
	script.WriteString(readFile(t, "seed.sql"))
	// Oddiy protokol: bitta matnda bir nechta SQL buyrugʻi boʻlishi mumkin.
	if _, err := pool.Exec(ctx, script.String(), pgx.QueryExecModeSimpleProtocol); err != nil {
		t.Fatalf("migratsiyalar yoki seed.sql bajarilmadi: %v", err)
	}
	return NewStore(pool)
}

func testCtx(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	t.Cleanup(cancel)
	return ctx
}

func exec(t *testing.T, s *Store, sql string) {
	t.Helper()
	if _, err := s.pool.Exec(testCtx(t), sql); err != nil {
		t.Fatalf("%s: %v", sql, err)
	}
}

func stockOf(t *testing.T, s *Store, id int64) int {
	t.Helper()
	var n int
	if err := s.pool.QueryRow(testCtx(t), "SELECT stock FROM products WHERE id = $1", id).Scan(&n); err != nil {
		t.Fatalf("mahsulot %d qoldigʻi oʻqilmadi: %v", id, err)
	}
	return n
}

func namesOf(ps []Product) string {
	names := make([]string, len(ps))
	for i, p := range ps {
		names[i] = p.Name
	}
	return strings.Join(names, " ")
}
