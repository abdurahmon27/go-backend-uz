package main

// Bu fayl vazifani (TASK.md) tekshiradi: go test uni main.go bilan birga ishga tushiradi.
// Uni oʻzgartirmang. Bu yerda ham defer va recover bor: handle panic'ni ushlamasa,
// test yiqilmaydi, aniq xabar chiqaradi.
//   t.TempDir() — test uchun vaqtinchalik papka, test tugagach oʻchiriladi
//   filepath.Join(a, b) — papka va fayl nomidan yoʻl yasaydi
//   os.ReadFile(path) — faylning butun ichini bayt slice'i ([]byte) qilib oʻqiydi;
//   string(b) uni matnga aylantiradi
//   strings.Contains(a, b) — a matnida b bormi
//   %%w — xabar matnida foiz belgisining oʻzi (%w)

import (
	"errors"
	"os"
	"path/filepath"
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

// Videodagi saveReport va safely (bu testlar boshidanoq oʻtadi).
func TestVideo(t *testing.T) {
	path := filepath.Join(t.TempDir(), "hisobot.txt")
	if err := saveReport(path, testProducts()); err != nil {
		t.Fatalf("saveReport xato qaytardi: %v", err)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("hisobot oʻqilmadi: %v", err)
	}
	want := "olma    12000 soʻm   40 dona\n" +
		"non      4000 soʻm  100 dona\n" +
		"oʻrik   35000 soʻm    5 dona\n"
	if string(b) != want {
		t.Errorf("hisobot:\n%s\nkutilgan:\n%s", b, want)
	}
	if err := saveReport(filepath.Join(t.TempDir(), "yoʻq", "hisobot.txt"), testProducts()); err == nil {
		t.Errorf("mavjud boʻlmagan papkaga yozishda saveReport xato qaytarmadi")
	}

	err = safely(func() {
		var p *Product // nil
		p.Stock -= 1
	})
	if err == nil || !strings.Contains(err.Error(), "nil pointer") {
		t.Errorf("safely(nil pointer) = %v; kutilgan «panic: runtime error: … nil pointer dereference»", err)
	}
	if err := safely(func() {}); err != nil {
		t.Errorf("safely(panic'siz fn) = %v; kutilgan <nil>", err)
	}
}

// callHandle handle'ni chaqiradi; handle panic'ni oʻtkazib yuborsa, uni ushlab qaytaradi.
func callHandle(name string, fn func() error) (panicked any, err error) {
	defer func() {
		panicked = recover()
	}()
	err = handle(name, fn)
	return nil, err
}

// 1-qadam: handle panic'ni ErrPanic'ni oʻragan xatoga aylantiradi, fn'ning xatosini oʻzgartirmaydi.
func TestVazifaHandle(t *testing.T) {
	panics := []struct {
		name string
		fn   func() error
		text string // xato matnida boʻlishi kerak
	}{
		// panic'dan keyin return shart emas: funksiya bu qatordan keyin davom etmaydi.
		{"sotish", func() error { panic("kutilmagan holat") }, "kutilmagan holat"},
		{"qoldiq", func() error {
			p := findProduct(testProducts(), 9) // nil
			p.Stock -= 1
			return nil
		}, "nil pointer dereference"},
	}
	for _, tc := range panics {
		t.Run(tc.name, func(t *testing.T) {
			panicked, err := callHandle(tc.name, tc.fn)
			if panicked != nil {
				t.Fatalf("handle(%q) panic'ni ushlamadi («%v»): handle'da defer va recover kerak (TASK.md, 1-qadam)", tc.name, panicked)
			}
			if !errors.Is(err, ErrPanic) {
				t.Fatalf("handle(%q) = %v; errors.Is(err, ErrPanic) false: ErrPanic'ni fmt.Errorf'da %%w bilan oʻrang (TASK.md, 1-qadam)", tc.name, err)
			}
			if !strings.Contains(err.Error(), tc.name+": ") || !strings.Contains(err.Error(), tc.text) {
				t.Errorf("handle(%q) xatosi «%v»: unda buyruq nomi «%s: » va panic qiymati «%s» boʻlsin (TASK.md, 1-qadam)", tc.name, err, tc.name, tc.text)
			}
		})
	}

	panicked, err := callHandle("sotish", func() error { return ErrBadQuantity })
	if panicked != nil || err != ErrBadQuantity {
		t.Errorf("fn oʻz xatosini qaytardi, handle = %v; kutilgan aynan «%v» (oʻzgartirmay qaytaring, TASK.md, 1-qadam)", err, ErrBadQuantity)
	}
	panicked, err = callHandle("royxat", func() error { return nil })
	if panicked != nil || err != nil {
		t.Errorf("fn xatosiz tugadi, handle = %v; kutilgan <nil> (TASK.md, 1-qadam)", err)
	}
}

// 2-qadam: reserve — hammasi yoki hech biri.
func TestVazifaReserve(t *testing.T) {
	products := testProducts()
	err := reserve(products, []OrderItem{{ProductID: 1, Quantity: 3}, {ProductID: 2, Quantity: 10}})
	if err != nil {
		t.Fatalf("reserve(olma 3, non 10) xato qaytardi: %v", err)
	}
	if products[0].Stock != 37 || products[1].Stock != 90 || products[2].Stock != 5 {
		t.Fatalf("muvaffaqiyatli buyurtmadan keyin qoldiq: olma %d, non %d, oʻrik %d; kutilgan 37, 90, 5",
			products[0].Stock, products[1].Stock, products[2].Stock)
	}

	fails := []struct {
		name  string
		items []OrderItem
		want  error // errors.Is bilan tekshiriladi; nil — *StockError kutiladi
	}{
		{"qoldiq_yetmaydi", []OrderItem{{ProductID: 1, Quantity: 3}, {ProductID: 2, Quantity: 500}}, nil},
		{"mahsulot_yoq", []OrderItem{{ProductID: 1, Quantity: 3}, {ProductID: 2, Quantity: 10}, {ProductID: 9, Quantity: 1}}, ErrNotFound},
		{"dona_soni", []OrderItem{{ProductID: 3, Quantity: 2}, {ProductID: 1, Quantity: 0}}, ErrBadQuantity},
		{"bir_mahsulot_ikki_qator", []OrderItem{{ProductID: 1, Quantity: 30}, {ProductID: 1, Quantity: 30}}, nil},
	}
	for _, tc := range fails {
		t.Run(tc.name, func(t *testing.T) {
			products := testProducts()
			err := reserve(products, tc.items)
			if err == nil {
				t.Fatalf("reserve(%v) xato qaytarmadi", tc.items)
			}
			var se *StockError
			if tc.want != nil && !errors.Is(err, tc.want) || tc.want == nil && !errors.As(err, &se) {
				t.Errorf("reserve xatosi «%v» tekshirib boʻlmaydigan boʻlib qoldi: xato matnlarini oʻzgartirmang (TASK.md)", err)
			}
			if products[0].Stock != 40 || products[1].Stock != 100 || products[2].Stock != 5 {
				t.Errorf("reserve xato qaytardi («%v»), lekin qoldiq: olma %d, non %d, oʻrik %d; kutilgan 40, 100, 5: xato boʻlsa, oldin kamaytirilgan qoldiqlarni defer bilan qaytaring (TASK.md, 2-qadam)",
					err, products[0].Stock, products[1].Stock, products[2].Stock)
			}
		})
	}
}
