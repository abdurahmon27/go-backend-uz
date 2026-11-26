package main

// Bu fayl vazifani tekshiradi: go test buyurtma qiymatlarini (name, price, active,
// stock, sold) oʻzgartirib, main() ni qayta-qayta ishga tushiradi va ekranga chiqqan
// qatorlarni kutilgan natija bilan solishtiradi. Test yozishni GB-1.08 darsida
// ochamiz; hozir uni oʻzgartirmang.

import (
	"fmt"
	"io"
	"os"
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
	kutilgan string
	vazifa   bool // true: TASK.md'dagi yangi qoida
}

// chek: dastur boshidagi ikki qator + berilgan qatorlar.
func chek(qatorlar ...string) string {
	return "bozor ishga tushdi\nversiya: 0.1.0\n" + strings.Join(qatorlar, "\n") + "\n"
}

var holatlar = []buyurtma{
	// videodagi qoidalar (main.go'da bor, buzilmasin)
	{"7 dona: 3- va 6-dona bepul", "choy", 199_900, true, 40, 7, chek(
		"choy, 7 dona        9995.00",
		"QQS                 1199.40",
		"yetkazish          15000.00",
		"jami               26194.40",
		"aksiya: 2 dona bepul"), false},
	{"2 dona: aksiya yoʻq", "choy", 199_900, true, 40, 2, chek(
		"choy, 2 dona        3998.00",
		"QQS                  479.76",
		"yetkazish          15000.00",
		"jami               19477.76"), false},
	{"miqdor 0", "choy", 199_900, true, 40, 0, chek("xato: miqdor kamida 1 dona boʻlishi kerak"), false},
	{"omborda yetmaydi", "choy", 199_900, true, 40, 50, chek("xato: omborda 40 dona bor, 50 soʻraldi"), false},

	// TASK.md: bitta buyurtmada koʻpi bilan 2 dona bepul
	{"12 dona: koʻpi bilan 2 dona bepul", "choy", 199_900, true, 40, 12, chek(
		"choy, 12 dona      19990.00",
		"QQS                 2398.80",
		"yetkazish          15000.00",
		"jami               37388.80",
		"aksiya: 2 dona bepul"), true},
	{"30 dona non: koʻpi bilan 2 dona bepul", "non", 500_000, true, 40, 30, chek(
		"non, 30 dona      140000.00",
		"QQS                16800.00",
		"yetkazish          10000.00",
		"jami              166800.00",
		"aksiya: 2 dona bepul"), true},
	{"24 dona oʻrik: 500 000 soʻmdan oshadi, yetkazish bepul", "oʻrik", 2_500_000, true, 40, 24, chek(
		"oʻrik, 24 dona    550000.00",
		"QQS                66000.00",
		"yetkazish              0.00",
		"jami              616000.00",
		"aksiya: 2 dona bepul"), true},
}

// chiqish buyurtma qiymatlarini oʻrnatib, main() ekranga yozgan matnni qaytaradi.
func chiqish(t *testing.T, b buyurtma) string {
	t.Helper()
	eskiName, eskiPrice, eskiActive, eskiStock, eskiSold := name, price, active, stock, sold
	defer func() { name, price, active, stock, sold = eskiName, eskiPrice, eskiActive, eskiStock, eskiSold }()
	name, price, active, stock, sold = b.name, b.price, b.active, b.stock, b.sold

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	eski := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = eski }()
	natija := make(chan string)
	go func() {
		out, _ := io.ReadAll(r)
		natija <- string(out)
	}()
	main()
	w.Close()
	return <-natija
}

// apostrof ʻ (U+02BB) oʻrniga ' yoki ’ yozilgan boʻlsa ham xabar toʻgʻri hisoblanadi.
var apostrof = strings.NewReplacer("'", "ʻ", "’", "ʻ", "‘", "ʻ")

func TestAksiya(t *testing.T) {
	for _, b := range holatlar {
		hozir := chiqish(t, b)
		if apostrof.Replace(hozir) == b.kutilgan {
			continue
		}
		maslahat := ""
		if b.vazifa {
			maslahat = "  bu TASK.md'dagi qoida: har uchinchi dona bepul, lekin bitta buyurtmada koʻpi bilan 2 dona\n" +
				"  (const maxFree = 2); 2 ta bepul dona berilgach, qolgan donalarning hammasi pullik\n"
		}
		t.Errorf("%s (name=%q price=%d active=%t stock=%d sold=%d):\n  kerak:\n%s  hozir:\n%s%s",
			b.nomi, b.name, b.price, b.active, b.stock, b.sold, qatorlar(b.kutilgan), qatorlar(hozir), maslahat)
	}
}

// qatorlar matnni qatorma-qator, surib va %q bilan chiqaradi: boʻsh joylar ham koʻrinadi.
func qatorlar(s string) string {
	var sb strings.Builder
	for _, q := range strings.Split(strings.TrimRight(s, "\n"), "\n") {
		fmt.Fprintf(&sb, "    %q\n", q)
	}
	return sb.String()
}
