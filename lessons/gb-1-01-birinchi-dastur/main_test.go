package main

// Bu fayl vazifani tekshiradi: `go test` uni ishga tushiradi.
// Test yozishni GB-1.08 (birinchi test) darsida oʻrganamiz;
// hozir uni oʻzgartirmang, faqat main.go'ni oʻzgartiring.

import (
	"bytes"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// kutilgan: dastur ekranga chiqarishi kerak boʻlgan qatorlar, tartibi bilan.
var kutilgan = []string{
	"bozor ishga tushdi",
	"versiya: 0.1.0",
}

// ekran main() funksiyasini ishga tushiradi va u standart chiqishga yozgan matnni qaytaradi:
// os.Stdout vaqtincha faylga almashtiriladi, keyin fayl oʻqiladi.
func ekran(t *testing.T) string {
	t.Helper()
	f, err := os.Create(filepath.Join(t.TempDir(), "ekran.txt"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	eski := os.Stdout
	os.Stdout = f
	main()
	os.Stdout = eski
	matn, err := os.ReadFile(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	return string(matn)
}

func TestIkkiQator(t *testing.T) {
	matn := ekran(t)
	if matn != "" && !strings.HasSuffix(matn, "\n") {
		t.Errorf("oxirgi qator yangi qator bilan tugashi kerak: fmt.Println shuni qiladi")
	}
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

// TestModul: papka (yoki uning ota papkasi) bozor moduli boʻlishi kerak.
func TestModul(t *testing.T) {
	papka, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		kod, err := os.ReadFile(filepath.Join(papka, "go.mod"))
		if err == nil {
			tekshirGoMod(t, string(kod))
			return
		}
		ota := filepath.Dir(papka)
		if ota == papka {
			t.Fatalf("go.mod topilmadi: papkada `go mod init github.com/<siz>/bozor` ni ishga tushiring")
		}
		papka = ota
	}
}

func tekshirGoMod(t *testing.T, kod string) {
	t.Helper()
	var modul, versiya string
	for _, qator := range strings.Split(kod, "\n") {
		if s, ok := strings.CutPrefix(qator, "module "); ok {
			modul = strings.TrimSpace(s)
		}
		if s, ok := strings.CutPrefix(qator, "go "); ok {
			versiya = strings.TrimSpace(s)
		}
	}
	if modul != "bozor" && !strings.HasSuffix(modul, "/bozor") {
		t.Errorf("go.mod'dagi modul yoʻli «…/bozor» bilan tugashi kerak, hozir «%s»", modul)
	}
	if versiya == "" {
		t.Errorf("go.mod'da `go` qatori yoʻq: uni `go mod init` yozadi")
	}
}
