package main

// Bu fayl vazifani (TASK.md) tekshiradi: go test uni main.go bilan birga ishga tushiradi.
// Uni oʻzgartirmang.

import (
	"fmt"
	"io"
	"strings"
	"testing"
)

func testProducts() []Product {
	return []Product{
		{ID: 1, Name: "olma", Price: 1_200_000, Stock: 40},
		{ID: 2, Name: "non", Price: 400_000, Stock: 100},
		{ID: 3, Name: "oʻrik", Price: 3_500_000, Stock: 5},
	}
}

// buzuqWriter — har doim xato qaytaradigan io.Writer (masalan, disk toʻlgan fayl).
type buzuqWriter struct{}

func (buzuqWriter) Write(p []byte) (int, error) {
	return 0, fmt.Errorf("disk toʻldi")
}

// Videodagi writeReport va String (bu test boshidanoq oʻtadi).
func TestHisobot(t *testing.T) {
	products := testProducts()
	want := "olma    12000 soʻm   40 dona\n" +
		"non      4000 soʻm  100 dona\n" +
		"oʻrik   35000 soʻm    5 dona\n"

	var sb strings.Builder
	err := writeReport(&sb, products)
	if err != nil {
		t.Fatalf("writeReport xato qaytardi: %v", err)
	}
	got := sb.String()
	if got != want {
		t.Errorf("hisobot:\n%s\nkutilgan:\n%s", got, want)
	}
	if got := fmt.Sprint(testProducts()[0]); got != "olma (12000 soʻm)" {
		t.Errorf("fmt.Sprint(olma) = %q; kutilgan %q", got, "olma (12000 soʻm)")
	}
}

// 1-qadam: writeLowStock hamma narsani w'ga yozadi va yozish xatosini qaytaradi.
func TestVazifaKamQoldi(t *testing.T) {
	products := append(testProducts(), Product{ID: 4, Name: "choy", Price: 2_500_000, Stock: 8})
	var sb strings.Builder
	err := writeLowStock(&sb, products, 10)
	if err != nil {
		t.Fatalf("writeLowStock xato qaytardi: %v", err)
	}
	want := "oʻrik (35000 soʻm) — 5 dona qoldi\n" +
		"choy (25000 soʻm) — 8 dona qoldi\n"
	got := sb.String()
	if got == "" {
		t.Fatalf("writeLowStock w'ga hech narsa yozmadi: fmt.Printf ekranga yozadi, w'ga fmt.Fprintf(w, …) yozadi (TASK.md, 1-qadam)")
	}
	if got != want {
		t.Errorf("writeLowStock(…, 10) yozdi:\n%s\nkutilgan:\n%s(TASK.md, 1-qadam)", got, want)
	}

	if err := writeLowStock(buzuqWriter{}, products, 10); err == nil {
		t.Errorf("w xato qaytardi («disk toʻldi»), writeLowStock esa nil qaytardi; Fprintf'ning xatosini qaytaring (TASK.md, 1-qadam)")
	}
}

// 2-qadam: *lineCounter — io.Writer, u yozilgan qatorlarni sanaydi.
func TestVazifaQatorSanagich(t *testing.T) {
	c := &lineCounter{}
	var x any = c
	w, ok := x.(io.Writer)
	if !ok {
		t.Fatalf("*lineCounter io.Writer emas: unga Write(p []byte) (int, error) metodini yozing (TASK.md, 2-qadam)")
	}

	n, err := w.Write([]byte("bir\nikki\n"))
	if n != 9 || err != nil {
		t.Errorf("Write(\"bir\\nikki\\n\") = %d, %v; kutilgan 9, <nil>: Write len(p) va nil qaytarsin (TASK.md, 2-qadam)", n, err)
	}
	if c.lines != 2 {
		t.Fatalf("ikki qatorli matn yozildi, lines = %d; kutilgan 2 — receiver pointer'mi? (TASK.md, 2-qadam)", c.lines)
	}

	c.lines = 0
	if err := writeReport(w, testProducts()); err != nil {
		t.Fatalf("writeReport xato qaytardi: %v", err)
	}
	if c.lines != 3 {
		t.Errorf("uch mahsulotli hisobot: lines = %d; kutilgan 3 (TASK.md, 2-qadam)", c.lines)
	}
}
