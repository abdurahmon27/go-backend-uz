// Bu fayl vazifani tekshiradi: `go test` uni ishga tushiradi.
// Test fayllarini GB-1.08 darsida toʻliq oʻqiymiz; hozir uni oʻzgartirmang.
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
	text := strings.TrimSuffix(<-out, "\n")
	if text == "" {
		return nil
	}
	return strings.Split(text, "\n")
}

func TestChek(t *testing.T) {
	got := runMain(t)
	if len(got) == 0 || got[0] != "bozor ishga tushdi" {
		t.Fatalf("1-qator «bozor ishga tushdi» boʻlishi kerak: oʻtgan darslardagi qatorlar qolsin, oʻzgartirish faqat chekda")
	}
	want := []string{
		"choy, 7 dona       13993.00",
		"QQS                 1679.16",
		"yetkazish          15000.00",
		"jami               30672.16",
	}
	if len(got) < len(want)+1 {
		t.Fatalf("dastur %d qator chiqardi, oxirgi %d qator chek boʻlishi kerak.\nkerak:\n%s",
			len(got), len(want), strings.Join(want, "\n"))
	}
	// Vazifadan oldingi holat: chek oxiri QQS, keyin darhol jami.
	last3 := strings.Join(got[len(got)-3:], "\n")
	if last3 == strings.Join([]string{want[0], want[1], want[3]}, "\n") {
		t.Fatalf("chek 4 qator boʻlishi kerak, hozir 3 qator: QQS va jami orasida «yetkazish» qatori yoʻq.\nkerak:\n%s\nhozir:\n%s",
			strings.Join(want, "\n"), last3)
	}
	tail := got[len(got)-len(want):]
	for i := range want {
		if tail[i] != want[i] {
			t.Errorf("chekning %d-qatori:\n  kerak: %q\n  hozir: %q", i+1, want[i], tail[i])
		}
	}
	if tail[2] != want[2] {
		t.Logf("maslahat: deliverySom soʻmda, chekdagi boshqa summalar tiyinda; shablon xuddi boshqa qatorlardagidek: %s", `"%-14s %9d.%02d\n"`)
	}
}

// Yetkazish summasi qoʻlda yozilmasin: u deliverySom'dan hisoblanadi.
// main.go'da "15000…" matni yoki const'dan tashqarida 15000 / 1500000 soni boʻlmasligi kerak.
func TestSummaHisoblanadi(t *testing.T) {
	f, err := parser.ParseFile(token.NewFileSet(), "main.go", nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	ast.Inspect(f, func(n ast.Node) bool {
		if d, ok := n.(*ast.GenDecl); ok && d.Tok == token.CONST {
			return false // const deliverySom = 15_000 shu yerda: unga tegilmaydi
		}
		lit, ok := n.(*ast.BasicLit)
		if !ok {
			return true
		}
		v := strings.ReplaceAll(lit.Value, "_", "")
		switch {
		case lit.Kind == token.STRING && strings.Contains(v, "15000"):
			t.Errorf("main.go'da %s matni bor: yetkazish summasini deliverySom'dan Printf bilan chiqaring, tayyor matn qilib yozmang", lit.Value)
		case lit.Kind == token.INT && (v == "15000" || v == "1500000"):
			t.Errorf("main.go'da %s soni bor: yetkazish summasi deliverySom'dan hisoblansin, sonni qoʻlda yozmang", lit.Value)
		}
		return true
	})
}
