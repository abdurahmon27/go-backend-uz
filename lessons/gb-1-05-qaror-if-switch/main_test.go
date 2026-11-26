// Bu fayl vazifani tekshiradi: `go test` uni ishga tushiradi. U main.go'dagi besh
// qiymatni (name, price, active, stock, sold) har xil buyurtma bilan almashtirib,
// dasturni qayta-qayta ishga tushiradi va ekranga chiqqan qatorlarni kutilgan
// natija bilan solishtiradi. Test fayllarini GB-1.08 darsida toʻliq oʻqiymiz;
// hozir uni oʻzgartirmang.
package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// buyurtma bitta holat: kirish qiymatlari va kutilgan chiqish.
type buyurtma struct {
	nomi     string // holat nomi, xabarda koʻrinadi
	name     string
	price    int64 // tiyinda
	active   bool
	stock    int
	sold     int
	kutilgan string // ekranning 5-qatoridan boshlab (boshidagi toʻrt qator oʻzgarmaydi)
	vazifa   bool   // true: TASK.md'dagi yangi qoida
}

func chek(qatorlar ...string) string { return strings.Join(qatorlar, "\n") + "\n" }

var holatlar = []buyurtma{
	// videodagi qoidalar (main.go'da bor, buzilmasin)
	{"oddiy buyurtma", "choy", 199_900, true, 40, 7, chek(
		"choy 33",
		"choy, 7 dona       13993.00",
		"QQS                 1679.16",
		"yetkazish          15000.00",
		"jami               30672.16"), false},
	{"miqdor 0", "choy", 199_900, true, 40, 0, chek("xato: miqdor kamida 1 dona boʻlishi kerak"), false},
	{"miqdor -3", "choy", 199_900, true, 40, -3, chek("xato: miqdor kamida 1 dona boʻlishi kerak"), false},
	{"omborda yetmaydi (50 > 40)", "choy", 199_900, true, 40, 50, chek("xato: omborda 10 dona yetmaydi"), false},
	{"omborda yetmaydi (41 > 40)", "choy", 199_900, true, 40, 41, chek("xato: omborda 1 dona yetmaydi"), false},
	{"omborning hammasi (40 = 40)", "choy", 199_900, true, 40, 40, chek(
		"choy 0",
		"choy, 40 dona      79960.00",
		"QQS                 9595.20",
		"yetkazish          15000.00",
		"jami              104555.20"), false},
	{"99 999 soʻm: yetkazish 15 000", "non", 9_999_900, true, 40, 1, chek(
		"non 39",
		"non, 1 dona        99999.00",
		"QQS                11999.88",
		"yetkazish          15000.00",
		"jami              126998.88"), false},
	{"roppa-rosa 100 000 soʻm: yetkazish 10 000", "non", 500_000, true, 40, 20, chek(
		"non 20",
		"non, 20 dona      100000.00",
		"QQS                12000.00",
		"yetkazish          10000.00",
		"jami              122000.00"), false},
	{"roppa-rosa 500 000 soʻm: yetkazish bepul", "oʻrik", 2_500_000, true, 40, 20, chek(
		"oʻrik 20",
		"oʻrik, 20 dona    500000.00",
		"QQS                60000.00",
		"yetkazish              0.00",
		"jami              560000.00"), false},

	// TASK.md: nom, sotuv va narx qoidasi
	{"nom boʻsh", "", 199_900, true, 40, 7, chek("xato: mahsulot nomi boʻsh"), true},
	{"sotuvda emas", "choy", 199_900, false, 40, 7, chek("xato: choy sotuvda emas"), true},
	{"narx 0", "choy", 0, true, 40, 7, chek("xato: narx notoʻgʻri: 0 tiyin"), true},
	{"narx manfiy", "choy", -199_900, true, 40, 7, chek("xato: narx notoʻgʻri: -199900 tiyin"), true},
	{"narx 10 000 000 soʻmdan 1 tiyin katta", "choy", 1_000_000_001, true, 40, 1, chek("xato: narx notoʻgʻri: 1000000001 tiyin"), true},
	{"narx 1 tiyin: toʻgʻri", "choy", 1, true, 40, 1, chek(
		"choy 39",
		"choy, 1 dona           0.01",
		"QQS                    0.00",
		"yetkazish          15000.00",
		"jami               15000.01"), true},
	{"narx roppa-rosa 10 000 000 soʻm: toʻgʻri", "choy", 1_000_000_000, true, 40, 1, chek(
		"choy 39",
		"choy, 1 dona    10000000.00",
		"QQS              1200000.00",
		"yetkazish              0.00",
		"jami            11200000.00"), true},
	{"nom ham, narx ham notoʻgʻri: avval nom", "", 0, true, 40, 7, chek("xato: mahsulot nomi boʻsh"), true},
	{"sotuvda emas, narx ham notoʻgʻri: avval sotuv", "choy", 0, false, 40, 7, chek("xato: choy sotuvda emas"), true},
	{"narx ham, miqdor ham notoʻgʻri: avval narx", "choy", 0, true, 40, 0, chek("xato: narx notoʻgʻri: 0 tiyin"), true},
}

// kirishlar: main() ichida qiymati almashtiriladigan qatorlar, masalan `sold := 7`.
var kirishlar = []string{"name", "price", "active", "stock", "sold"}

// qurish main.go'dagi besh qiymatni kirish funksiyalariga almashtirib, dasturni
// vaqtinchalik papkada quradi va binary yoʻlini qaytaradi.
func qurish(t *testing.T) string {
	t.Helper()
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "main.go", nil, parser.ParseComments)
	if err != nil {
		t.Fatalf("main.go oʻqilmadi: %v", err)
	}
	topildi := map[string]bool{}
	ast.Inspect(f, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if !ok || fn.Name.Name != "main" {
			return true
		}
		for _, st := range fn.Body.List {
			as, ok := st.(*ast.AssignStmt)
			if !ok || len(as.Lhs) != 1 || len(as.Rhs) != 1 {
				continue
			}
			id, ok := as.Lhs[0].(*ast.Ident)
			if !ok || topildi[id.Name] {
				continue
			}
			for _, k := range kirishlar {
				if id.Name == k && (k == "sold" || k == "stock" || as.Tok == token.ASSIGN) {
					as.Rhs[0] = &ast.CallExpr{Fun: ast.NewIdent("kirish_" + k)}
					topildi[k] = true
				}
			}
		}
		return false
	})
	for _, k := range kirishlar {
		if !topildi[k] {
			t.Fatalf("main() ichida %s'ga qiymat beradigan qator topilmadi (masalan `sold := 7`): videodagi qatorlarni oʻchirmang", k)
		}
	}
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, f); err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	mod, err := os.ReadFile("go.mod")
	if err != nil {
		mod, err = os.ReadFile(filepath.Join("..", "go.mod")) // solution/ papkasi
	}
	if err != nil {
		t.Fatal(err)
	}
	fayllar := map[string]string{
		"go.mod":  string(mod),
		"main.go": buf.String(),
		"kirish.go": `package main

import (
	"os"
	"strconv"
)

func kirish_name() string { return os.Getenv("BOZOR_NAME") }
func kirish_active() bool { return os.Getenv("BOZOR_ACTIVE") == "true" }
func kirish_price() int64 { v, _ := strconv.ParseInt(os.Getenv("BOZOR_PRICE"), 10, 64); return v }
func kirish_stock() int   { v, _ := strconv.Atoi(os.Getenv("BOZOR_STOCK")); return v }
func kirish_sold() int    { v, _ := strconv.Atoi(os.Getenv("BOZOR_SOLD")); return v }
`,
	}
	for nom, matn := range fayllar {
		if err := os.WriteFile(filepath.Join(dir, nom), []byte(matn), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	bin := filepath.Join(dir, "bozor")
	cmd := exec.Command("go", "build", "-o", bin, ".")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("dastur qurilmadi: %v\n%s", err, out)
	}
	return bin
}

// apostrof ʻ (U+02BB) oʻrniga ' yoki ’ yozilgan boʻlsa ham xabar toʻgʻri hisoblanadi.
var apostrof = strings.NewReplacer("'", "ʻ", "’", "ʻ", "‘", "ʻ")

func TestQoidalar(t *testing.T) {
	bin := qurish(t)
	for _, b := range holatlar {
		cmd := exec.Command(bin)
		cmd.Env = append(os.Environ(),
			"BOZOR_NAME="+b.name,
			"BOZOR_PRICE="+strconv.FormatInt(b.price, 10),
			"BOZOR_ACTIVE="+strconv.FormatBool(b.active),
			"BOZOR_STOCK="+strconv.Itoa(b.stock),
			"BOZOR_SOLD="+strconv.Itoa(b.sold))
		out, err := cmd.Output()
		if err != nil {
			t.Fatalf("%s: dastur xato bilan tugadi: %v", b.nomi, err)
		}
		qatorlar := strings.SplitAfter(string(out), "\n")
		if len(qatorlar) < 5 || qatorlar[0] != "bozor ishga tushdi\n" {
			t.Fatalf("%s: dastur boshidagi toʻrt qator (bozor ishga tushdi, versiya, ikki sinov qatori) joyida qolsin, hozir:\n%s", b.nomi, out)
		}
		hozir := strings.Join(qatorlar[4:], "")
		if apostrof.Replace(hozir) == b.kutilgan {
			continue
		}
		maslahat := ""
		if b.vazifa {
			maslahat = "  bu TASK.md'dagi qoida: sold := 7'dan keyin, miqdor tekshiruvidan oldin: avval nom, keyin active, keyin narx\n"
		}
		t.Errorf("%s (name=%q price=%d active=%t stock=%d sold=%d):\n  kerak:\n%s  hozir:\n%s%s",
			b.nomi, b.name, b.price, b.active, b.stock, b.sold, korsat(b.kutilgan), korsat(hozir), maslahat)
	}
}

// korsat matnni qatorma-qator, surib va %q bilan chiqaradi: boʻsh joylar ham koʻrinadi.
func korsat(s string) string {
	var sb strings.Builder
	for _, q := range strings.Split(strings.TrimRight(s, "\n"), "\n") {
		fmt.Fprintf(&sb, "    %q\n", q)
	}
	return sb.String()
}
