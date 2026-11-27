// Bu fayl vazifani tekshiradi: products.sql va arzon.sql'ni haqiqiy Postgres'da bajaradi.
// Uni oʻzgartirmang. Postgres uchun: docker compose up -d --wait (shu papkada).
// Tekshiruv alohida vaqtinchalik sxemada ishlaydi va oxirida uni oʻchiradi: public sxemaga tegmaydi.
// Go dasturidan Postgres'ga ulanishni (pgx) GB-5.05 oʻrgatadi; shuning uchun bu yerda bogʻliqlik yoʻq:
// test shunchaki `docker compose exec db psql` buyrugʻini ishga tushiradi.
package lab

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

// psql SQL matnini db konteynerida psql orqali bajaradi (stdin orqali) va chiqishini qaytaradi.
// -At: faqat qatorlar, ustunlar | bilan; ON_ERROR_STOP: birinchi xatoda toʻxtaydi.
func psql(schema, sql string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "docker", "compose", "exec", "-T", "db",
		"psql", "-X", "-q", "-At", "-v", "ON_ERROR_STOP=1")
	cmd.Stdin = strings.NewReader("SET search_path = " + schema + ";\n" + sql + "\n")
	var out bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

// dockerPsql db konteynerida psql'ni berilgan argumentlar bilan 30 soniyagacha ishlatadi.
func dockerPsql(args ...string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	full := append([]string{"compose", "exec", "-T", "db", "psql", "-X", "-q"}, args...)
	return exec.CommandContext(ctx, "docker", full...).CombinedOutput()
}

// newSchema bu test uchun boʻsh sxema yaratadi; test tugagach uni oʻchiradi.
func newSchema(t *testing.T) string {
	t.Helper()
	name := fmt.Sprintf("lab_%d", time.Now().UnixNano())
	out, err := dockerPsql("-v", "ON_ERROR_STOP=1", "-c", "CREATE SCHEMA "+name)
	if err != nil {
		t.Fatalf("Postgres'ga ulanib boʻlmadi: avval shu papkada `docker compose up -d --wait` ishga tushiring.\n%v\n%s", err, out)
	}
	t.Cleanup(func() {
		if out, err := dockerPsql("-c", "DROP SCHEMA "+name+" CASCADE"); err != nil {
			t.Logf("sxema %s oʻchirilmadi (qoʻlda: DROP SCHEMA %s CASCADE): %v\n%s", name, name, err, out)
		}
	})
	return name
}

func readFile(t *testing.T, name string) string {
	t.Helper()
	b, err := os.ReadFile(name)
	if err != nil {
		t.Fatalf("%s oʻqilmadi: %v", name, err)
	}
	return string(b)
}

// products.sql'ni sxemada bajaradi.
func applySchema(t *testing.T) string {
	t.Helper()
	s := newSchema(t)
	if out, err := psql(s, readFile(t, "products.sql")); err != nil {
		t.Fatalf("products.sql bajarilmadi (sintaksis xatosi?): %v\n%s", err, out)
	}
	return s
}

func TestConstraintlar(t *testing.T) {
	s := applySchema(t)
	// Toʻgʻri qator qabul qilinadi.
	if out, err := psql(s, "INSERT INTO products (name, price, stock) VALUES ('olma', 1200000, 40);"); err != nil {
		t.Fatalf("toʻgʻri qator rad etildi: ('olma', 1200000, 40)\n%s", out)
	}
	cases := []struct {
		name, insert, want, hint string
	}{
		{"dublikat nom", "INSERT INTO products (name, price, stock) VALUES ('olma', 1, 1);", "unique", "name UNIQUE boʻlishi kerak"},
		{"NULL nom", "INSERT INTO products (name, price, stock) VALUES (NULL, 1, 1);", "not-null", "name NOT NULL boʻlishi kerak"},
		{"narx nol", "INSERT INTO products (name, price, stock) VALUES ('a', 0, 1);", "check", "price uchun CHECK (price > 0) kerak"},
		{"manfiy qoldiq", "INSERT INTO products (name, price, stock) VALUES ('b', 10, -1);", "check", "stock uchun CHECK (stock >= 0) kerak"},
		{"boʻsh nom", "INSERT INTO products (name, price, stock) VALUES ('', 10, 1);", "check", "products.sql'da name ustuniga CHECK (name <> '') qoʻshing (TASK.md, 1-qadam)"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := psql(s, tc.insert)
			if err == nil {
				t.Errorf("%s: qator qabul qilindi, baza rad etishi kerak edi. %s", tc.name, tc.hint)
				return
			}
			if !strings.Contains(out, tc.want) {
				t.Errorf("%s: %q turidagi rad etish kutilgan edi, chiqdi:\n%s", tc.name, tc.want, out)
			}
		})
	}
}

func TestArzon(t *testing.T) {
	query := strings.TrimSpace(readFile(t, "arzon.sql"))
	if !strings.Contains(strings.ToUpper(stripComments(query)), "SELECT") {
		t.Fatalf("arzon.sql'da SELECT yoʻq: faol mahsulotlardan narxi 3 000 000 tiyindan kam boʻlgan eng arzon ikkitasini tanlang (TASK.md, 2-qadam)")
	}
	cases := []struct{ name, seed, want string }{
		{"katalog", `INSERT INTO products (name, price, stock, active) VALUES
		  ('olma', 1200000, 40, true), ('non', 400000, 3, true), ('choy', 2500000, 8, true),
		  ('oʻrik', 3500000, 2, true), ('sut', 900000, 15, false), ('qand', 1100000, 10, true);`,
			"non|400000\nqand|1100000"},
		{"chegara", `INSERT INTO products (name, price, stock, active) VALUES
		  ('x', 2999999, 1, true), ('y', 3000000, 1, true), ('z', 3000001, 1, true), ('w', 10, 1, false);`,
			"x|2999999"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := applySchema(t)
			if out, err := psql(s, tc.seed); err != nil {
				t.Fatalf("sinov maʼlumotini yozib boʻlmadi (products.sql ustunlarini oʻzgartirmang: id, name, price, stock, active):\n%s", out)
			}
			out, err := psql(s, query)
			if err != nil {
				t.Fatalf("arzon.sql bajarilmadi:\n%s", out)
			}
			if out != tc.want {
				t.Errorf("arzon.sql natijasi notoʻgʻri.\nkutilgan (name, price; narx oʻsish tartibida; faqat faol; narxi 3 000 000 dan KAM, ya'ni < ; eng koʻpi 2 ta):\n%s\nchiqdi:\n%s", tc.want, out)
			}
		})
	}
}

// stripComments "--" bilan boshlangan qatorlarni olib tashlaydi.
func stripComments(sql string) string {
	var keep []string
	for _, line := range strings.Split(sql, "\n") {
		if !strings.HasPrefix(strings.TrimSpace(line), "--") {
			keep = append(keep, line)
		}
	}
	return strings.Join(keep, "\n")
}
