package main

// Bu fayl vazifani (TASK.md) tekshiradi: go test uni main.go bilan birga ishga tushiradi.
// Uni oʻzgartirmang.

import "testing"

// Holatlarning raqamlari: ular saqlanadi, shuning uchun oʻzgarmasligi kerak (videoning «Chuqurroq» qismi).
func TestHolatRaqamlari(t *testing.T) {
	tests := []struct {
		name string
		got  OrderStatus
		want OrderStatus
	}{
		{"StatusNew", StatusNew, 0},
		{"StatusPaid", StatusPaid, 1},
		{"StatusShipped", StatusShipped, 2},
		{"StatusCancelled", StatusCancelled, 3},
	}
	for _, tc := range tests {
		if tc.got != tc.want {
			t.Errorf("%s = %d; kutilgan %d (yangi holat faqat roʻyxat oxiriga qoʻshiladi)", tc.name, tc.got, tc.want)
		}
	}
}

// 1-qadam: orderTotal har bir qatorning narxini dona soniga koʻpaytirib qoʻshadi (tiyinda).
func TestVazifaOrderTotal(t *testing.T) {
	tests := []struct {
		name  string
		order Order
		want  int64
	}{
		{"bitta_qator", Order{Items: []OrderItem{{ProductID: 1, Price: 1_200_000, Quantity: 3}}}, 3_600_000},
		{"ikki_qator", Order{Items: []OrderItem{
			{ProductID: 1, Price: 1_200_000, Quantity: 3},
			{ProductID: 2, Price: 400_000, Quantity: 2},
		}}, 4_400_000},
		{"bitta_dona", Order{Items: []OrderItem{{ProductID: 3, Price: 2_500_000, Quantity: 1}}}, 2_500_000},
		{"bosh_buyurtma", Order{}, 0},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := orderTotal(tc.order); got != tc.want {
				t.Errorf("orderTotal(%+v) = %d; kutilgan %d tiyin (TASK.md, 1-qadam)", tc.order, got, tc.want)
			}
		})
	}
}

// 2-qadam: yangi va toʻlangan buyurtmani bekor qilish mumkin, joʻnatilgan va bekor qilinganni — yoʻq.
func TestVazifaCanCancel(t *testing.T) {
	tests := []struct {
		name   string
		status OrderStatus
		want   bool
	}{
		{"yangi", StatusNew, true},
		{"tolangan", StatusPaid, true},
		{"jonatilgan", StatusShipped, false},
		{"bekor_qilingan", StatusCancelled, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			o := Order{ID: 7, Status: tc.status}
			if got := canCancel(o); got != tc.want {
				t.Errorf("canCancel(Status: %d) = %v; kutilgan %v (TASK.md, 2-qadam)", tc.status, got, tc.want)
			}
		})
	}
	// Zero value'li Order — yangi buyurtma, uni ham bekor qilish mumkin.
	var empty Order
	if !canCancel(empty) {
		t.Errorf("canCancel(Order{}) = false; kutilgan true: zero value'li buyurtmaning holati StatusNew (TASK.md, 2-qadam)")
	}
}
