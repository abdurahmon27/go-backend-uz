package store

// store paketining testlari (GB-1.08 jadval shakli). Test store paketining ichida,
// shuning uchun kichik harfli products maydonini ham koʻradi.

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// katalog videodagi toʻrt mahsulotli katalogni qaytaradi (narx tiyinda).
func katalog() *Store {
	return &Store{products: []Product{
		{ID: 1, Name: "olma", Price: 1_200_000, Stock: 40},
		{ID: 2, Name: "non", Price: 400_000, Stock: 3},
		{ID: 3, Name: "choy", Price: 2_500_000, Stock: 8},
		{ID: 4, Name: "oʻrik", Price: 3_500_000, Stock: 2},
	}}
}

func TestSell(t *testing.T) {
	tests := []struct {
		name      string
		id        int64
		quantity  int
		wantStock int   // shu ID'li mahsulotning katalogdagi qoldigʻi Sell'dan keyin
		wantErr   error // errors.Is bilan tekshiriladi; nil — xato kutilmaydi
	}{
		{"oddiy", 1, 3, 37, nil},
		{"hammasi", 4, 2, 0, nil},
		{"yoq_mahsulot", 9, 1, 0, ErrNotFound},
		{"nol_dona", 1, 0, 40, ErrBadQuantity},
		{"manfiy_dona", 1, -2, 40, ErrBadQuantity},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			st := katalog()
			_, err := st.Sell(tc.id, tc.quantity)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("Sell(%d, %d) xatosi %v, kutilgan %v", tc.id, tc.quantity, err, tc.wantErr)
			}
			if p := st.find(tc.id); p != nil && p.Stock != tc.wantStock {
				t.Fatalf("Sell(%d, %d) dan keyin katalogda qoldiq %d, kutilgan %d", tc.id, tc.quantity, p.Stock, tc.wantStock)
			}
		})
	}
}

// Qoldiq yetmasa — *StockError (GB-2.08), katalog oʻzgarmaydi.
func TestSellStockError(t *testing.T) {
	st := katalog()
	_, err := st.Sell(4, 5)
	var se *StockError
	if !errors.As(err, &se) || se.Stock != 2 || se.Want != 5 {
		t.Fatalf("Sell(4, 5) xatosi %v, kutilgan *StockError{oʻrik, 2, 5}", err)
	}
	if p := st.find(4); p.Stock != 2 {
		t.Fatalf("rad etilgan Sell qoldiqni %d qildi, kutilgan 2", p.Stock)
	}
}

func TestAdd(t *testing.T) {
	tests := []struct {
		name  string
		pname string
		price int64
		stock int
		ok    bool
	}{
		{"oddiy", "anor", 3_000_000, 1, true},
		{"nol_qoldiq", "anor", 3_000_000, 0, true},
		{"bosh_nom", "", 3_000_000, 1, false},
		{"utf8_emas", "o\xca", 3_000_000, 1, false},
		{"uzun_nom", "oʻgʻil bola koʻylagi!", 3_000_000, 1, false}, // 21 belgi
		{"nol_narx", "anor", 0, 1, false},
		{"manfiy_qoldiq", "anor", 3_000_000, -1, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			st := katalog()
			p, err := st.Add(tc.pname, tc.price, tc.stock)
			if (err == nil) != tc.ok {
				t.Fatalf("Add(%q, %d, %d) xatosi %v, qabul qilinishi kerak: %v", tc.pname, tc.price, tc.stock, err, tc.ok)
			}
			if tc.ok && p.ID != 5 {
				t.Fatalf("yangi mahsulot ID'si %d, kutilgan 5", p.ID)
			}
			if !tc.ok && (len(st.products) != 4 || !errors.Is(err, ErrBadProduct)) {
				t.Fatalf("rad etilgan mahsulot: katalogda %d ta, xato %v", len(st.products), err)
			}
		})
	}
}

// Save → Load: katalog fayldan xuddi oʻzidek qaytadi, vaqtinchalik fayl qolmaydi.
func TestSaveLoad(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bozor.json")
	if err := katalog().Save(path); err != nil {
		t.Fatal(err)
	}
	st, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(st.products, katalog().products) {
		t.Fatalf("Load qaytardi %v, kutilgan %v", st.products, katalog().products)
	}
	entries, err := os.ReadDir(filepath.Dir(path))
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("papkada %d ta fayl, kutilgan faqat bozor.json (vaqtinchalik fayl qolmasin)", len(entries))
	}
}

// Save muvaffaqiyatsiz boʻlsa, eski bozor.json butun qoladi va vaqtinchalik fayl qolmaydi.
func TestSaveKeepsOld(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bozor.json")
	if err := katalog().Save(path); err != nil {
		t.Fatal(err)
	}
	old, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(dir, 0o500); err != nil { // papka faqat oʻqish uchun: yangi fayl ochib boʻlmaydi
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chmod(dir, 0o700) })
	var empty Store
	if err := empty.Save(path); err == nil {
		t.Skip("papka yozishdan himoyalanmadi (root?): bu tekshiruvni oʻtkazib boʻlmaydi")
	}
	if now, _ := os.ReadFile(path); string(now) != string(old) {
		t.Fatalf("Save xato berdi, lekin bozor.json oʻzgardi:\n%s", now)
	}
	if entries, _ := os.ReadDir(dir); len(entries) != 1 {
		t.Fatalf("papkada %d ta fayl, kutilgan faqat bozor.json", len(entries))
	}
}

// Boʻsh katalog [] boʻlib yoziladi, null emas.
func TestSaveEmpty(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bozor.json")
	var st Store
	if err := st.Save(path); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(path)
	if got := strings.TrimSpace(string(data)); got != "[]" {
		t.Fatalf("boʻsh katalog %q boʻlib yozildi, kutilgan []", got)
	}
}

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		data    string // "-" — fayl umuman yoʻq
		wantN   int
		wantErr string
	}{
		{"fayl_yoq", "-", 0, ""},
		{"bosh", "[]", 0, ""},
		{"nol_bayt", "", 0, "unexpected end of JSON input"}, // os.WriteFile boʻshatib, yoza olmay qolgan fayl
		{"oddiy", `[{"id":1,"name":"olma","price":1200000,"stock":40}]`, 1, ""},
		{"yarim_fayl", `[{"id":1,"name":"ol`, 0, "unexpected end of JSON input"},
		{"manfiy_qoldiq", `[{"id":1,"name":"olma","price":1200000,"stock":-5}]`, 0, "qoldiq manfiy"},
		{"takror_id", `[{"id":1,"name":"olma","price":1,"stock":1},{"id":1,"name":"non","price":1,"stock":1}]`, 0, "takrorlangan"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "bozor.json")
			if tc.data != "-" {
				if err := os.WriteFile(path, []byte(tc.data), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			st, err := Load(path)
			if tc.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("Load xatosi %v, kutilgan %q", err, tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if len(st.products) != tc.wantN {
				t.Fatalf("Load %d ta mahsulot, kutilgan %d", len(st.products), tc.wantN)
			}
		})
	}
}
