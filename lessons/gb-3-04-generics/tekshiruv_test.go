package main

// Bu fayl vazifani (TASK.md) tekshiradi: go test uni main.go bilan birga ishga tushiradi.
// Uni oʻzgartirmang. SumBy bu yerda uch xil tur bilan chaqiriladi: Product, Order va string.
// Funksiya generic boʻlmasa, test kompilyatsiya ham boʻlmaydi.

import "testing"

var testProducts = []Product{
	{ID: 1, Name: "olma", Price: 1_200_000, Stock: 40},
	{ID: 2, Name: "non", Price: 400_000, Stock: 3},
	{ID: 3, Name: "choy", Price: 2_500_000, Stock: 0},
}

func testOrders() []Order {
	return []Order{
		{ID: 7, Status: StatusPaid, Items: []OrderItem{{ProductID: 1, Price: 1_200_000, Quantity: 3}}},
		{ID: 8, Status: StatusNew, Items: []OrderItem{{ProductID: 2, Price: 400_000, Quantity: 5}}},
		{ID: 10, Status: StatusPaid, Items: []OrderItem{
			{ProductID: 2, Price: 400_000, Quantity: 2},
			{ProductID: 3, Price: 2_500_000, Quantity: 1},
		}},
	}
}

// Videodagi Filter va Sum (bu test boshidanoq oʻtadi).
func TestVideo(t *testing.T) {
	low := Filter(testProducts, func(p Product) bool { return p.Stock < 5 })
	if len(low) != 2 || low[0].Name != "non" || low[1].Name != "choy" {
		t.Errorf("Filter(qoldiq < 5) = %v; kutilgan non va choy", low)
	}
	paid := Filter(testOrders(), func(o Order) bool { return o.Status == StatusPaid })
	if len(paid) != 2 || paid[0].ID != 7 || paid[1].ID != 10 {
		t.Errorf("Filter(toʻlangan) = %v; kutilgan 7 va 10-buyurtmalar", paid)
	}
	if got := Filter([]string{"olma", "oʻrik", "non"}, func(s string) bool { return s != "non" }); len(got) != 2 {
		t.Errorf("Filter(string) = %v; kutilgan [olma oʻrik]", got)
	}
	if got := Sum([]int{40, 3, 0}); got != 43 {
		t.Errorf("Sum([]int{40, 3, 0}) = %d; kutilgan 43", got)
	}
	if got := Sum([]int64{3_600_000, 800_000}); got != 4_400_000 {
		t.Errorf("Sum([]int64{3_600_000, 800_000}) = %d; kutilgan 4400000", got)
	}
	if !Contains([]string{"olma", "non"}, "non") || Contains([]int64{1, 2}, 3) {
		t.Errorf("Contains notoʻgʻri ishlayapti")
	}
}

// 1-qadam: SumBy har xil tur va har xil son turi bilan.
func TestVazifaSumBy(t *testing.T) {
	const step = "TASK.md, 1-qadam"
	stock := SumBy(testProducts, func(p Product) int { return p.Stock })
	if stock != 43 {
		t.Errorf("SumBy(mahsulotlar, Stock) = %d; kutilgan 43 (40 + 3 + 0). Har bir x uchun value(x)'ni total'ga qoʻshing (%s)", stock, step)
	}
	value := SumBy(testProducts, func(p Product) int64 { return p.Price * int64(p.Stock) })
	if value != 49_200_000 {
		t.Errorf("SumBy(mahsulotlar, Price × Stock) = %d; kutilgan 49200000 tiyin (%s)", value, step)
	}
	letters := SumBy([]string{"olma", "non"}, func(s string) int { return len(s) })
	if letters != 7 {
		t.Errorf("SumBy([olma non], len) = %d; kutilgan 7 (%s)", letters, step)
	}
	if empty := SumBy([]Order(nil), func(o Order) int64 { return o.Total() }); empty != 0 {
		t.Errorf("boʻsh slice uchun SumBy = %d; kutilgan 0 (%s)", empty, step)
	}
}

// 2-qadam: paidTotal faqat toʻlangan buyurtmalarni qoʻshadi.
func TestVazifaPaidTotal(t *testing.T) {
	const step = "TASK.md, 2-qadam"
	orders := testOrders()
	// Testning oʻz maʼlumoti (main'dagidan boshqa): 7: 3 × 1 200 000 = 3 600 000;
	// 10: 2 × 400 000 + 1 × 2 500 000 = 3 300 000; 8 toʻlanmagan.
	if got := paidTotal(orders); got != 6_900_000 {
		t.Errorf("paidTotal = %d; kutilgan 6900000 tiyin (7 va 10-buyurtmalar, 8 toʻlanmagan) (%s)", got, step)
	}
	if got := paidTotal(orders[1:2]); got != 0 {
		t.Errorf("faqat toʻlanmagan buyurtma: paidTotal = %d; kutilgan 0 (%s)", got, step)
	}
}
