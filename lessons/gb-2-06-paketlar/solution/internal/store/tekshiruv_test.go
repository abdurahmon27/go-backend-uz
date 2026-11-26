package store

// Bu fayl vazifani (TASK.md) tekshiradi: go test ./... uni store.go bilan birga ishga tushiradi.
// Uni oʻzgartirmang. Test store paketining ichida, shuning uchun products maydonini ham koʻradi.

import "testing"

// katalog ikki mahsulotli yangi katalog qaytaradi: olma (ID 1, qoldiq 40) va non (ID 2, qoldiq 100).
func katalog() *Store {
	var st Store
	st.Add("olma", 1_200_000, 40)
	st.Add("non", 400_000, 100)
	return &st
}

// Videodagi Add: ID'lar 1 dan boshlab beriladi (bu test boshidanoq oʻtadi).
func TestAdd(t *testing.T) {
	st := katalog()
	if len(st.products) != 2 || st.products[0].ID != 1 || st.products[1].ID != 2 {
		t.Fatalf("Add'dan keyin katalog %+v; kutilgan ikki mahsulot, ID 1 va 2", st.products)
	}
}

// 1-qadam: Sell qoidasi.
func TestVazifaSell(t *testing.T) {
	st := katalog()
	if err := st.Sell(2, 3); err != nil {
		t.Fatalf("Sell(2, 3) = %v; kutilgan nil: nonning qoldigʻi 100 (TASK.md, 1-qadam)", err)
	}
	if got := st.products[1].Stock; got != 97 {
		t.Errorf("Sell(2, 3)'dan keyin non qoldigʻi %d; kutilgan 97 (TASK.md, 1-qadam)", got)
	}

	tests := []struct {
		name     string
		id       int64
		quantity int
	}{
		{"qoldiqdan_kop", 1, 41},
		{"nol_dona", 1, 0},
		{"manfiy_dona", 1, -2},
		{"yoq_mahsulot", 9, 1},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			st := katalog()
			err := st.Sell(tc.id, tc.quantity)
			if err == nil {
				t.Errorf("Sell(%d, %d) = nil; kutilgan error (TASK.md, 1-qadam)", tc.id, tc.quantity)
			}
			if st.products[0].Stock != 40 || st.products[1].Stock != 100 {
				t.Errorf("xato boʻlgan Sell(%d, %d)'dan keyin qoldiqlar %d va %d; kutilgan 40 va 100: xato boʻlsa qoldiq oʻzgarmasin (TASK.md, 1-qadam)",
					tc.id, tc.quantity, st.products[0].Stock, st.products[1].Stock)
			}
		})
	}
}

// 2-qadam: List natijasiga yozish katalogni oʻzgartirmasin.
func TestVazifaList(t *testing.T) {
	st := katalog()
	list := st.List()
	if len(list) != 2 {
		t.Fatalf("len(List()) = %d; kutilgan 2", len(list))
	}
	list[0].Stock = -5
	if got := st.products[0].Stock; got != 40 {
		t.Errorf("List() natijasida olma qoldigʻi -5 qilindi va katalogdagi qoldiq ham %d boʻldi; kutilgan 40: List nusxa qaytarsin (TASK.md, 2-qadam)", got)
	}
}
