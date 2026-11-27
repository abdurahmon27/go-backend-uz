package main

// Bu fayl vazifani (TASK.md) tekshiradi: go test uni main.go bilan birga ishga tushiradi.
// Uni oʻzgartirmang.
//   t.Fatal(err) — t.Fatalf kabi, faqat shablonsiz: xatoni chiqarib, testni toʻxtatadi
//   t.Run("tur mos emas", …) — natijada nomdagi probellar _ boʻlib chiqadi (tur_mos_emas)
//   t.TempDir() — test uchun vaqtinchalik papka, test tugagach oʻchiriladi
//   filepath.Join(a, b) — papka va fayl nomidan yoʻl yasaydi
//   os.WriteFile / os.ReadFile, json.Marshal / json.Unmarshal — videodagidek
//   got[i] != want[i] — struct'lar hamma maydonlari teng boʻlsa teng
//   strings.HasPrefix(s, p) — s matni p bilan boshlanadimi
//   *json.SyntaxError — JSON matni buzuq (vergul yoʻq va hokazo),
//   *json.UnmarshalTypeError — kalitdagi qiymat turi maydon turiga mos emas;
//   ikkalasini errors.As (GB-2.08) oʻralgan xato ichidan topadi

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func testProducts() []Product {
	return []Product{
		{ID: 1, Name: "olma", Price: 1_200_000, Stock: 40},
		{ID: 2, Name: "oʻrik", Price: 3_500_000, Stock: 5, Note: "quritilgan"},
	}
}

// Videodagi saveProducts va loadProducts: yozilgan katalog xuddi oʻzidek qaytadi
// (bu test boshidanoq oʻtadi).
func TestVideo(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bozor.json")
	want := testProducts()
	if err := saveProducts(path, want); err != nil {
		t.Fatalf("saveProducts xato qaytardi: %v", err)
	}
	got, err := loadProducts(path)
	if err != nil {
		t.Fatalf("loadProducts xato qaytardi: %v", err)
	}
	if len(got) != len(want) {
		t.Fatalf("loadProducts %d ta mahsulot qaytardi, kutilgan %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("%d-mahsulot: %+v; kutilgan %+v", i, got[i], want[i])
		}
	}
}

// 1-qadam: buyurtma JSON'ining kalitlari.
func TestVazifaOrderJSON(t *testing.T) {
	o := Order{
		ID:     7,
		Items:  []OrderItem{{ProductID: 1, Price: 1_200_000, Quantity: 3}},
		Status: StatusNew,
	}
	data, err := json.Marshal(o)
	if err != nil {
		t.Fatalf("json.Marshal(Order) xato qaytardi: %v", err)
	}
	want := `{"id":7,"items":[{"product_id":1,"price":1200000,"quantity":3}],"status":0}`
	if string(data) != want {
		t.Errorf("json.Marshal(Order):\n  hozir    %s\n  kutilgan %s\nOrder va OrderItem maydonlariga json tag'larini qoʻying; status'da omitempty boʻlmasin (TASK.md, 1-qadam)", data, want)
	}

	in := `{"id":8,"items":[{"product_id":2,"price":3500000,"quantity":1}],"status":1}`
	var back Order
	if err := json.Unmarshal([]byte(in), &back); err != nil {
		t.Fatalf("json.Unmarshal(%s) xato qaytardi: %v", in, err)
	}
	// Unmarshal kalitni katta-kichik harfni farqlamay qidiradi: "id" → ID tag'siz ham tushadi,
	// "product_id" esa ProductID'ga faqat tag bilan tushadi.
	if back.ID != 8 || back.Status != StatusPaid || len(back.Items) != 1 {
		t.Fatalf("json.Unmarshal(%s) = %+v; id, status yoki items maydonga tushmadi (TASK.md, 1-qadam)", in, back)
	}
	if back.Items[0].ProductID != 2 {
		t.Errorf("json.Unmarshal: «product_id» kaliti OrderItem.ProductID'ga tushmadi (hozir %d, kutilgan 2) — tag kerak (TASK.md, 1-qadam)", back.Items[0].ProductID)
	}
	if back.Items[0].Price != 3_500_000 || back.Items[0].Quantity != 1 {
		t.Errorf("json.Unmarshal: price/quantity = %d/%d; kutilgan 3500000/1 (TASK.md, 1-qadam)", back.Items[0].Price, back.Items[0].Quantity)
	}
}

// 2-qadam (a): bozor birinchi marta ishga tushganda bozor.json hali yoʻq.
func TestVazifaFaylYoq(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bozor.json") // bu fayl yaratilmagan
	got, err := loadProducts(path)
	if err != nil {
		t.Errorf("fayl yoʻq: loadProducts xato qaytardi: %v; kutilgan: boʻsh katalog va nil — errors.Is(err, os.ErrNotExist) (TASK.md, 2-qadam)", err)
	}
	if len(got) != 0 {
		t.Errorf("fayl yoʻq: loadProducts %d ta mahsulot qaytardi; kutilgan 0", len(got))
	}
}

// 2-qadam (b): buzuq faylning xatosida fayl nomi boʻlsin, asl JSON xatosi esa yoʻqolmasin.
func TestVazifaBuzuqFayl(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{"tur mos emas", `[{"id":1,"name":"olma","price":"1200000","stock":40}]`},
		{"vergul yoʻq", `[{"id":1,"name":"olma" "price":1200000}]`},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "bozor.json")
			if err := os.WriteFile(path, []byte(tc.content), 0o644); err != nil {
				t.Fatal(err)
			}
			got, err := loadProducts(path)
			if err == nil {
				t.Fatalf("buzuq fayl: loadProducts xato qaytarmadi (%d ta mahsulot)", len(got))
			}
			if got != nil {
				t.Errorf("buzuq fayl: xato bilan birga %d ta mahsulot qaytdi; kutilgan nil", len(got))
			}
			if !strings.HasPrefix(err.Error(), path+": ") {
				t.Errorf("xato «%v»\n  fayl nomi bilan boshlanishi kerak: «%s: …» — fmt.Errorf bilan path'ni qoʻshing (TASK.md, 2-qadam)", err, path)
			}
			var se *json.SyntaxError
			var te *json.UnmarshalTypeError
			if !errors.As(err, &se) && !errors.As(err, &te) {
				t.Errorf("xato «%v» ichida asl JSON xatosi yoʻq: %%w bilan oʻrang, %%v bilan emas (TASK.md, 2-qadam)", err)
			}
		})
	}
}
