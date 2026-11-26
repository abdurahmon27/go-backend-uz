// Bu fayl vazifani tekshiradi: `go test` uni ishga tushiradi.
// Test fayllarini GB-1.08 darsida toʻliq oʻqiymiz; hozir uni oʻzgartirish shart emas.
package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"strings"
	"testing"
)

// runMain main() ni ishga tushiradi va u ekranga chiqargan qatorlarni qaytaradi.
func runMain(t *testing.T) []string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	out := make(chan string)
	go func() {
		b, _ := io.ReadAll(r)
		out <- string(b)
	}()
	stdout := os.Stdout
	os.Stdout = w
	func() {
		defer func() {
			os.Stdout = stdout
			w.Close()
		}()
		main()
	}()
	return strings.Split(strings.TrimRight(<-out, "\n"), "\n")
}

func TestBuyurtma(t *testing.T) {
	got := runMain(t)
	if got[0] != "bozor ishga tushdi" {
		t.Errorf("1-qator «bozor ishga tushdi» boʻlishi kerak, hozir «%s»: oʻtgan darsdagi qatorlar qolsin", got[0])
	}
	want := []string{
		"choy 7 dona",
		"QQS: 1679 soʻm 16 tiyin",
		"jami: 30672 soʻm 16 tiyin",
	}
	if len(got) < len(want)+1 {
		t.Fatalf("dastur %d qator chiqardi: oxirgi %d qator buyurtma boʻlishi kerak (TASK.md)", len(got), len(want))
	}
	tail := got[len(got)-len(want):]
	for i, w := range want {
		if tail[i] != w {
			t.Errorf("oxiridan %d-qator «%s» boʻlishi kerak, hozir «%s»", len(want)-i, w, tail[i])
		}
	}
}

// main.go: price int64 (tiyinda), sold := 7 (int) va hisobda int64(sold) bor.
func TestPulTiyinda(t *testing.T) {
	f, err := parser.ParseFile(token.NewFileSet(), "main.go", nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	price, sold, conv := false, false, false
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.ValueSpec:
			if typ, ok := x.Type.(*ast.Ident); ok {
				for _, id := range x.Names {
					if id.Name == "price" && typ.Name == "int64" {
						price = true
					}
				}
			}
		case *ast.AssignStmt:
			if id, ok := x.Lhs[0].(*ast.Ident); ok && id.Name == "sold" && x.Tok == token.DEFINE && len(x.Rhs) == 1 {
				if lit, ok := x.Rhs[0].(*ast.BasicLit); ok && lit.Value == "7" {
					sold = true
				}
			}
		case *ast.CallExpr:
			if fn, ok := x.Fun.(*ast.Ident); ok && fn.Name == "int64" && len(x.Args) == 1 {
				if a, ok := x.Args[0].(*ast.Ident); ok && a.Name == "sold" {
					conv = true
				}
			}
		}
		return true
	})
	if !price {
		t.Error("main.go'da «var price int64» boʻlishi kerak: narx tiyinda, int64")
	}
	if !sold {
		t.Error("main.go'da «sold := 7» qatori boʻlishi kerak (sold turi int boʻlib qolsin)")
	}
	if !conv {
		t.Error("hisobda int64(sold) conversion boʻlishi kerak")
	}
}
