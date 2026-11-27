// Bu fayl vazifani tekshiradi: sotilgan.sql va izohsiz.sql'ni haqiqiy Postgres'da bajaradi.
// Uni oʻzgartirmang. Postgres uchun: docker compose up -d --wait (shu papkada).
// Tekshiruv products.sql, orders.sql va seed.sql'ni alohida vaqtinchalik sxemada bajaradi va oxirida uni oʻchiradi:
// public sxemaga tegmaydi. Go dasturidan Postgres'ga ulanishni (pgx) GB-5.05 oʻrgatadi; shuning uchun bu yerda
// bogʻliqlik yoʻq: test shunchaki `docker compose exec db psql` buyrugʻini ishga tushiradi.
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

// newLab products.sql, orders.sql va seed.sql'ni yangi sxemada bajaradi; more boʻlsa, qoʻshimcha maʼlumot yozadi.
func newLab(t *testing.T, more string) string {
	t.Helper()
	s := newSchema(t)
	for _, f := range []string{"products.sql", "orders.sql", "seed.sql"} {
		if out, err := psql(s, readFile(t, f)); err != nil {
			t.Fatalf("%s bajarilmadi (oʻzgartirmagan boʻlsangiz, bu laboratoriya xatosi):\n%s", f, out)
		}
	}
	if more != "" {
		if out, err := psql(s, more); err != nil {
			t.Fatalf("sinov maʼlumotini yozib boʻlmadi:\n%s", out)
		}
	}
	return s
}

// moreData: ikkinchi sinov maʼlumoti. Mahsulotlar: sut (id 5), qand (id 6, hech sotilmagan).
// Buyurtmalar: Bobur (id 4, izohi boʻsh matn, NULL emas), Zarina (id 5), Sardor (id 6, izohsiz, choy olgan).
// UPDATE: 1-buyurtmaning qatori jadval oxiriga koʻchadi, shuning uchun ORDER BY'siz soʻrov boshqa tartibda qaytaradi.
const moreData = `
UPDATE orders SET note = note WHERE id = 1;
INSERT INTO products (name, price, stock) VALUES ('sut', 900000, 5), ('qand', 1100000, 7);
INSERT INTO orders (customer, note) VALUES ('Bobur', ''), ('Zarina', 'qoʻngʻiroq qiling'), ('Sardor', NULL);
INSERT INTO order_items (order_id, product_id, quantity, unit_price)
VALUES (4, 5, 6, 900000), (5, 5, 1, 900000), (5, 2, 2, 400000), (6, 3, 3, 2500000);
`

// needSelect fayl hali boʻsh (faqat izoh) boʻlsa, testni vazifa xabari bilan toʻxtatadi.
func needSelect(t *testing.T, file, hint string) {
	t.Helper()
	if !strings.Contains(strings.ToUpper(stripComments(readFile(t, file))), "SELECT") {
		t.Fatalf("%s'da SELECT yoʻq: %s (TASK.md)", file, hint)
	}
}

// run soʻrovni ikki rejalashtirish usulida (oddiy va enable_hashagg = off) bajaradi: natija soʻrovning oʻzidan
// (ORDER BY), rejadan tasodifan emas, kelishi kerak.
func run(t *testing.T, file, more, want, spec string) {
	t.Helper()
	query := strings.TrimSpace(readFile(t, file))
	s := newLab(t, more)
	for _, setup := range []string{"", "SET enable_hashagg = off;\n"} {
		out, err := psql(s, setup+query)
		if err != nil {
			t.Fatalf("%s bajarilmadi:\n%s", file, out)
		}
		if out != want {
			t.Errorf("%s natijasi notoʻgʻri.\nkutilgan (%s):\n%s\nchiqdi:\n%s", file, spec, want, out)
			return
		}
	}
}

func TestSotilgan(t *testing.T) {
	needSelect(t, "sotilgan.sql", "har mahsulot uchun name va sotilgan (order_items.quantity yigʻindisi) ni tanlang; hech sotilmaganda 0 chiqsin")
	const spec = "name | sotilgan; sotilgan kamayish, teng boʻlsa name oʻsish tartibida; hech sotilmagan mahsulotda 0"
	t.Run("seed", func(t *testing.T) {
		run(t, "sotilgan.sql", "", "non|3\nolma|3\noʻrik|1\nchoy|0", spec)
	})
	t.Run("koʻproq maʼlumot", func(t *testing.T) {
		run(t, "sotilgan.sql", moreData, "sut|7\nnon|5\nchoy|3\nolma|3\noʻrik|1\nqand|0", spec)
	})
}

func TestIzohsiz(t *testing.T) {
	needSelect(t, "izohsiz.sql", "izohi yoʻq (NULL) buyurtmalarning id va customer ustunlarini tanlang")
	const spec = "id | customer; faqat note IS NULL (boʻsh matn '' emas); id oʻsish tartibida"
	t.Run("seed", func(t *testing.T) {
		run(t, "izohsiz.sql", "", "1|Dilnoza\n3|Malika", spec)
	})
	t.Run("koʻproq maʼlumot", func(t *testing.T) {
		run(t, "izohsiz.sql", moreData, "1|Dilnoza\n3|Malika\n6|Sardor", spec)
	})
}
