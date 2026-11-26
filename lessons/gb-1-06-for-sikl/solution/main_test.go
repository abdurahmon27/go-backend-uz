package main

// Bu fayl vazifani tekshiradi: go test buyurtma qiymatlarini (name, price, qty,
// stock) oʻzgartirib, main() ni qayta-qayta ishga tushiradi va ekranga chiqqan
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
	qty      int
	stock    int
	kutilgan string
	vazifa   bool // true: TASK.md'dagi yangi qoida
}

var holatlar = []buyurtma{
	{"7 dona: 3- va 6-dona bepul", "choy", 199_900, 7, 40, "" +
		"choy, 7 dona        9995.00\n" +
		"QQS                 1199.40\n" +
		"yetkazish          15000.00\n" +
		"jami               26194.40\n" +
		"aksiya: 2 dona bepul\n", false},
	{"2 dona: aksiya yoʻq", "choy", 199_900, 2, 40, "" +
		"choy, 2 dona        3998.00\n" +
		"QQS                  479.76\n" +
		"yetkazish          15000.00\n" +
		"jami               19477.76\n", false},
	{"miqdor 0", "choy", 199_900, 0, 40, "xato: miqdor kamida 1 dona boʻlishi kerak\n", false},
	{"omborda yetmaydi", "choy", 199_900, 50, 40, "xato: omborda 40 dona bor, 50 soʻraldi\n", false},
	{"12 dona: koʻpi bilan 2 dona bepul", "choy", 199_900, 12, 40, "" +
		"choy, 12 dona      19990.00\n" +
		"QQS                 2398.80\n" +
		"yetkazish          15000.00\n" +
		"jami               37388.80\n" +
		"aksiya: 2 dona bepul\n", true},
	{"30 dona non: koʻpi bilan 2 dona bepul", "non", 500_000, 30, 40, "" +
		"non, 30 dona      140000.00\n" +
		"QQS                16800.00\n" +
		"yetkazish          10000.00\n" +
		"jami              166800.00\n" +
		"aksiya: 2 dona bepul\n", true},
	{"24 dona oʻrik: 500 000 soʻmdan oshadi, yetkazish bepul", "oʻrik", 2_500_000, 24, 40, "" +
		"oʻrik, 24 dona    550000.00\n" +
		"QQS                66000.00\n" +
		"yetkazish              0.00\n" +
		"jami              616000.00\n" +
		"aksiya: 2 dona bepul\n", true},
}

// chiqish buyurtma qiymatlarini oʻrnatib, main() ekranga yozgan matnni qaytaradi.
func chiqish(t *testing.T, b buyurtma) string {
	t.Helper()
	eskiName, eskiPrice, eskiQty, eskiStock := name, price, qty, stock
	defer func() { name, price, qty, stock = eskiName, eskiPrice, eskiQty, eskiStock }()
	name, price, qty, stock = b.name, b.price, b.qty, b.stock

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	eski := os.Stdout
	os.Stdout = w
	main()
	os.Stdout = eski
	w.Close()
	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	return string(out)
}

func TestAksiya(t *testing.T) {
	for _, b := range holatlar {
		hozir := chiqish(t, b)
		if hozir == b.kutilgan {
			continue
		}
		maslahat := ""
		if b.vazifa {
			maslahat = "  bu TASK.md'dagi qoida: har uchinchi dona bepul, lekin bitta buyurtmada koʻpi bilan 2 dona\n" +
				"  (const maxFree = 2); 2 ta bepul dona berilgach, qolgan donalarning hammasi pullik"
		}
		t.Errorf("%s (name=%q price=%d qty=%d stock=%d):\n  kerak:\n%s  hozir:\n%s%s",
			b.nomi, b.name, b.price, b.qty, b.stock, qatorlar(b.kutilgan), qatorlar(hozir), maslahat)
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
