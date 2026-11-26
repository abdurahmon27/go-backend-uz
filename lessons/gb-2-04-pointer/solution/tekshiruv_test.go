package main

// Bu fayl vazifani (TASK.md) tekshiradi: go test uni main.go bilan birga ishga tushiradi.
// Uni oʻzgartirmang.

import "testing"

func newCatalog() []Product {
	return []Product{
		{ID: 1, Name: "olma", Price: 1_200_000, Stock: 40},
		{ID: 2, Name: "non", Price: 400_000, Stock: 100},
		{ID: 3, Name: "choy", Price: 2_500_000, Stock: 7},
	}
}

// 1-qadam: sellByID katalogdagi elementning oʻzini oʻzgartiradi; xato boʻlsa qoldiqqa tegmaydi.
func TestVazifaSellByID(t *testing.T) {
	tests := []struct {
		name      string
		id        int64
		quantity  int
		wantErr   bool
		wantStock []int // sotuvdan keyingi qoldiqlar: olma, non, choy
	}{
		{"olma_3_dona", 1, 3, false, []int{37, 100, 7}},
		{"choy_hammasi", 3, 7, false, []int{40, 100, 0}},
		{"yoq_mahsulot", 9, 1, true, []int{40, 100, 7}},
		{"nol_dona", 2, 0, true, []int{40, 100, 7}},
		{"manfiy_dona", 2, -5, true, []int{40, 100, 7}},
		{"qoldiq_yetmaydi", 3, 8, true, []int{40, 100, 7}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			catalog := newCatalog()
			err := sellByID(catalog, tc.id, tc.quantity)
			if tc.wantErr && err == nil {
				t.Errorf("sellByID(id %d, %d dona): xato kutilgan, nil qaytdi (TASK.md, 1-qadam)", tc.id, tc.quantity)
			}
			if !tc.wantErr && err != nil {
				t.Errorf("sellByID(id %d, %d dona): xato kutilmagan, keldi: %v (TASK.md, 1-qadam)", tc.id, tc.quantity, err)
			}
			for i, want := range tc.wantStock {
				if got := catalog[i].Stock; got != want {
					t.Errorf("sellByID(id %d, %d dona) dan keyin %s qoldigʻi %d; kutilgan %d (TASK.md, 1-qadam)",
						tc.id, tc.quantity, catalog[i].Name, got, want)
				}
			}
		})
	}
}

// 2-qadam: cheapest katalogdagi elementning manzilini qaytaradi (nusxaning emas).
func TestVazifaCheapest(t *testing.T) {
	t.Run("bosh_katalog", func(t *testing.T) {
		if got := cheapest(nil); got != nil {
			t.Errorf("cheapest(boʻsh katalog) = %+v; kutilgan nil (TASK.md, 2-qadam)", *got)
		}
	})
	t.Run("uch_mahsulot", func(t *testing.T) {
		catalog := newCatalog()
		got := cheapest(catalog)
		if got == nil {
			t.Fatalf("cheapest(olma, non, choy) = nil; kutilgan non'ning manzili (TASK.md, 2-qadam)")
		}
		if got.Name != "non" {
			t.Fatalf("cheapest(olma, non, choy) → %s; kutilgan non, eng arzoni 400000 tiyin (TASK.md, 2-qadam)", got.Name)
		}
		got.Stock = 0
		if catalog[1].Stock != 0 {
			t.Errorf("cheapest qaytargan pointer orqali yozildi, lekin katalogdagi non qoldigʻi %d qoldi: "+
				"pointer elementning nusxasini koʻrsatyapti, &products[i] kerak (TASK.md, 2-qadam)", catalog[1].Stock)
		}
	})
}
