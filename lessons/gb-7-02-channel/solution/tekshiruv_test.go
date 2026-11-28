package main

// Vazifani tekshiradi (TASK.md). Oʻzgartirmang.
//
// Har tekshiruv vaqt chegarasi bilan ishlaydi: bosqich out'ni yopmasa, `for p := range out` abadiy kutadi,
// shuning uchun test 2 soniyadan keyin oʻzbekcha xabar bilan yiqiladi (select va time.After: GB-7.03).

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
	"time"
)

// runUnique in'ga products'ni yuborib yopadi, unique(in, out)ni ishga tushiradi va out'dan hammasini yigʻadi.
func runUnique(products []Product) ([]Product, bool) {
	in := make(chan Product)
	out := make(chan Product)
	go func() {
		for _, p := range products {
			in <- p
		}
		close(in)
	}()
	go unique(in, out)
	done := make(chan []Product)
	go func() {
		var got []Product
		for p := range out {
			got = append(got, p)
		}
		done <- got
	}()
	select {
	case got := <-done:
		return got, true
	case <-time.After(2 * time.Second):
		return nil, false
	}
}

func TestVazifaUnique(t *testing.T) {
	input := []Product{
		{ID: 1, Name: "olma", Price: 1_200_000, Stock: 40},
		{ID: 3, Name: "non", Price: 400_000, Stock: 100},
		{ID: 4, Name: "choy", Price: 2_500_000, Stock: 8},
		{ID: 3, Name: "non", Price: 420_000, Stock: 60},
		{ID: 1, Name: "olma", Price: 1_250_000, Stock: 7},
		{ID: 5, Name: "anor", Price: 1_800_000, Stock: 25},
	}
	got, ok := runUnique(input)
	if !ok {
		t.Fatalf("unique 2 soniyada tugamadi: out yopilmadi, shuning uchun `for p := range out` abadiy kutyapti.\n" +
			"\tunique in'dagi hamma qiymatni oʻqishi (for p := range in) va oxirida out'ni yopishi kerak (defer close(out)).")
	}
	want := []Product{input[0], input[1], input[2], input[5]}
	if len(got) != len(want) {
		t.Fatalf("unique %d ta mahsulot oʻtkazdi, kutilgan %d ta (ID 1, 3, 4, 5; har ID bir marta). Keldi: %v",
			len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("%d-mahsulot %v, kutilgan %v: unique har ID'ning birinchisini oʻtkazadi, tartibni oʻzgartirmaydi",
				i+1, got[i], want[i])
		}
	}
	if got, ok := runUnique(nil); !ok || len(got) != 0 {
		t.Errorf("boʻsh in: unique hech narsa yubormasdan out'ni yopishi kerak (ok=%v, keldi %v)", ok, got)
	}
}

// importAllCallsUnique main.go'dagi importAll ichida unique chaqirilganini tekshiradi
// (natijani oxirgi siklda map bilan tozalash ham testdan oʻtardi, lekin vazifa — bosqich).
func importAllCallsUnique(t *testing.T) bool {
	t.Helper()
	f, err := parser.ParseFile(token.NewFileSet(), "main.go", nil, 0)
	if err != nil {
		t.Fatalf("main.go oʻqilmadi: %v", err)
	}
	found := false
	for _, d := range f.Decls {
		fn, ok := d.(*ast.FuncDecl)
		if !ok || fn.Name.Name != "importAll" {
			continue
		}
		ast.Inspect(fn, func(n ast.Node) bool {
			if call, ok := n.(*ast.CallExpr); ok {
				if id, ok := call.Fun.(*ast.Ident); ok && id.Name == "unique" {
					found = true
				}
			}
			return true
		})
	}
	return found
}

func TestVazifaImportAll(t *testing.T) {
	if !importAllCallsUnique(t) {
		t.Errorf("importAll ichida unique bosqichi yoʻq: TASK.md 2-qadam, inStock → yangi channel → go unique(…, ready)")
	}
	// Goroutine'lar har safar boshqa tartibda ishlaydi, shuning uchun bir necha marta tekshiramiz.
	for run := 1; run <= 50; run++ {
		done := make(chan []Product)
		go func() { done <- importAll(chorsu, oloy) }()
		var got []Product
		select {
		case got = <-done:
		case <-time.After(2 * time.Second):
			t.Fatalf("importAll 2 soniyada tugamadi (%d-urinish): qaysidir channel yopilmadi.\n"+
				"\tHar bosqich oʻz out'ini yopadi; unique'ga yangi channel bering: inStock → stocked → unique → ready.", run)
		}
		count := map[int64]int{}
		for _, p := range got {
			count[p.ID]++
		}
		for _, id := range []int64{1, 3, 4, 5} {
			if count[id] != 1 {
				t.Fatalf("importAll natijasida ID %d %d marta (%d-urinish), kutilgan 1 marta. Natija: %v\n"+
					"\tTASK.md 2-qadam: inStock va natija orasiga unique bosqichini qoʻying.", id, count[id], run, got)
			}
		}
		if len(got) != 4 {
			t.Fatalf("importAll %d ta mahsulot qaytardi, kutilgan 4 ta (olma, non, choy, anor; oʻrikning qoldigʻi 0). Natija: %v",
				len(got), got)
		}
	}
}
