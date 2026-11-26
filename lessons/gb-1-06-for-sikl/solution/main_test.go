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
	{"7 dona: 3- va 6-dona bepul", "choy", 199_900, true, 40, 7, chek(
		"choy 33",
		"choy, 7 dona        9995.00",
		"QQS                 1199.40",
		"yetkazish          15000.00",
		"jami               26194.40",
		"aksiya: 2 dona bepul"), false},
	{"2 dona: aksiya yoʻq", "choy", 199_900, true, 40, 2, chek(
		"choy 38",
		"choy, 2 dona        3998.00",
		"QQS                  479.76",
		"yetkazish          15000.00",
		"jami               19477.76"), false},
	{"miqdor 0", "choy", 199_900, true, 40, 0, chek("xato: miqdor kamida 1 dona boʻlishi kerak"), false},
	{"omborda yetmaydi (50 > 40)", "choy", 199_900, true, 40, 50, chek("xato: omborda 10 dona yetmaydi"), false},
	{"nom boʻsh", "", 199_900, true, 40, 7, chek("xato: mahsulot nomi boʻsh"), false},
	{"sotuvda emas", "choy", 199_900, false, 40, 7, chek("xato: choy sotuvda emas"), false},
	{"narx 0", "choy", 0, true, 40, 7, chek("xato: narx notoʻgʻri: 0 tiyin"), false},

	// TASK.md: bitta buyurtmada koʻpi bilan 2 dona bepul
	{"12 dona: koʻpi bilan 2 dona bepul", "choy", 199_900, true, 40, 12, chek(
		"choy 28",
		"choy, 12 dona      19990.00",
		"QQS                 2398.80",
		"yetkazish          15000.00",
		"jami               37388.80",
		"aksiya: 2 dona bepul"), true},
	{"30 dona non: koʻpi bilan 2 dona bepul", "non", 500_000, true, 40, 30, chek(
		"non 10",
		"non, 30 dona      140000.00",
		"QQS                16800.00",
		"yetkazish          10000.00",
		"jami              166800.00",
		"aksiya: 2 dona bepul"), true},
	{"24 dona oʻrik: 500 000 soʻmdan oshadi, yetkazish bepul", "oʻrik", 2_500_000, true, 40, 24, chek(
		"oʻrik 16",
		"oʻrik, 24 dona    550000.00",
		"QQS                66000.00",
		"yetkazish              0.00",
		"jami              616000.00",
		"aksiya: 2 dona bepul"), true},
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

func TestAksiya(t *testing.T) {
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
			maslahat = "  bu TASK.md'dagi qoida: har uchinchi dona bepul, lekin bitta buyurtmada koʻpi bilan 2 dona\n" +
				"  (const maxFree = 2); 2 ta bepul dona berilgach, qolgan donalarning hammasi pullik\n"
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
