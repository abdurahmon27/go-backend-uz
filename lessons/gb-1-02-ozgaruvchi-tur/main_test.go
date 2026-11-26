package main

// Bu fayl vazifani tekshiradi: `go test` uni ishga tushiradi.
// Ichidagi har bir qatorni GB-1.08 (birinchi test) darsida ochamiz;
// hozir uni oʻzgartirmang, faqat main.go'ni oʻzgartiring.

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"os"
	"strings"
	"testing"
)

// kutilgan: dastur ekranga chiqarishi kerak boʻlgan qatorlar, tartibi bilan.
var kutilgan = []string{
	"bozor ishga tushdi",
	"versiya: 0.1.0",
	" 0 false",
	"olma 12000 true 40 1.5",
	"olma 37",
}

// ekran main() funksiyasini ishga tushiradi va u ekranga yozgan matnni qaytaradi.
func ekran(t *testing.T) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	eski := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = eski }()

	natija := make(chan string)
	go func() {
		var b bytes.Buffer
		io.Copy(&b, r)
		natija <- b.String()
	}()
	main()
	w.Close()
	return <-natija
}

func TestQatorlar(t *testing.T) {
	matn := ekran(t)
	qatorlar := strings.Split(strings.TrimSuffix(matn, "\n"), "\n")
	if matn == "" {
		qatorlar = nil
	}
	for i, k := range kutilgan {
		if i >= len(qatorlar) {
			t.Errorf("%d-qator «%s» boʻlishi kerak, hozir dastur faqat %d qator chiqaryapti (TASK.md)",
				i+1, k, len(qatorlar))
			continue
		}
		if qatorlar[i] != k {
			t.Errorf("%d-qator «%s» boʻlishi kerak, hozir «%s»", i+1, k, qatorlar[i])
		}
	}
	if len(qatorlar) > len(kutilgan) {
		t.Errorf("dastur %d qator chiqarishi kerak, hozir %d qator: ortiqchasi «%s»",
			len(kutilgan), len(qatorlar), qatorlar[len(kutilgan)])
	}
}

// main.go ichida `sold := …` bor va stock'ning yangi qiymati sold'dan hisoblanadi.
func TestSold(t *testing.T) {
	f, err := parser.ParseFile(token.NewFileSet(), "main.go", nil, 0)
	if err != nil {
		t.Fatalf("main.go'da yozuv xatosi bor: %v", err)
	}
	soldBor, stockSolddan := false, false
	ast.Inspect(f, func(n ast.Node) bool {
		a, ok := n.(*ast.AssignStmt)
		if !ok || len(a.Lhs) != 1 {
			return true
		}
		nom, ok := a.Lhs[0].(*ast.Ident)
		if !ok {
			return true
		}
		if nom.Name == "sold" && a.Tok == token.DEFINE {
			soldBor = true
		}
		if nom.Name == "stock" && (a.Tok == token.ASSIGN || a.Tok == token.SUB_ASSIGN) {
			for _, e := range a.Rhs {
				ast.Inspect(e, func(m ast.Node) bool {
					if x, ok := m.(*ast.Ident); ok && x.Name == "sold" {
						stockSolddan = true
					}
					return true
				})
			}
		}
		return true
	})
	if !soldBor {
		t.Error("main.go'da «sold := 3» kabi eʼlon yoʻq: sotilgan qadoqlar soni uchun sold degan oʻzgaruvchi kerak")
	}
	if !stockSolddan {
		t.Error("stock'ning yangi qiymati sold'dan hisoblanib, bor oʻzgaruvchiga yozilishi kerak (:= emas: stock allaqachon bor)")
	}
}

func TestGoFmt(t *testing.T) {
	kod, err := os.ReadFile("main.go")
	if err != nil {
		t.Fatalf("main.go oʻqilmadi: %v", err)
	}
	toza, err := format.Source(kod)
	if err != nil {
		t.Fatalf("main.go'da yozuv xatosi bor: %v", err)
	}
	if !bytes.Equal(kod, toza) {
		t.Errorf("main.go go fmt formatida emas: shu papkada `go fmt ./...` ni ishga tushiring")
	}
}
