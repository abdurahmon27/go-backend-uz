package store

import "testing"

// 10 000 mahsulotli katalogda birinchi kam qolgan mahsulotni (2-oʻrinda) topish:
// List bilan va All bilan. go test -bench . -benchmem ./internal/store
func katta() *Store {
	var st Store
	for i := range 10_000 {
		stock := 40
		if i == 1 {
			stock = 3
		}
		st.Add("olma", 1_200_000, stock)
	}
	return &st
}

func BenchmarkList(b *testing.B) {
	st := katta()
	for b.Loop() {
		for _, p := range st.List() {
			if p.Stock < 5 {
				break
			}
		}
	}
}

func BenchmarkAll(b *testing.B) {
	st := katta()
	for b.Loop() {
		for p := range st.All() {
			if p.Stock < 5 {
				break
			}
		}
	}
}
