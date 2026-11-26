package main

// Bu fayl vazifani (TASK.md) tekshiradi: go test uni main.go bilan birga ishga tushiradi.
// Uni oʻzgartirmang.

import "testing"

// Videodagi String metodi: har bir holatning nomi (bu test boshidanoq oʻtadi).
func TestHolatNomlari(t *testing.T) {
	tests := []struct {
		status OrderStatus
		want   string
	}{
		{StatusNew, "yangi"},
		{StatusPaid, "toʻlangan"},
		{StatusShipped, "joʻnatilgan"},
		{StatusCancelled, "bekor qilingan"},
		{OrderStatus(9), "nomaʼlum"},
	}
	for _, tc := range tests {
		if got := tc.status.String(); got != tc.want {
			t.Errorf("OrderStatus(%d).String() = %q; kutilgan %q", int(tc.status), got, tc.want)
		}
	}
}

// 1-qadam: bir mahsulot ikki marta qoʻshilsa, yangi qator emas, oʻsha qatorning Quantity'si oshadi.
func TestVazifaAddItem(t *testing.T) {
	olma := Product{ID: 1, Name: "olma", Price: 1_200_000, Stock: 40}
	non := Product{ID: 2, Name: "non", Price: 400_000, Stock: 100}

	var o Order
	o.AddItem(olma, 3)
	o.AddItem(non, 2)
	o.AddItem(olma, 1)

	if len(o.Items) != 2 {
		t.Fatalf("olma 3, non 2, olma 1 qoʻshildi: len(Items) = %d; kutilgan 2 — olma ikkinchi qator boʻlib qoʻshilmasligi kerak (TASK.md, 1-qadam)", len(o.Items))
	}
	if got := o.Items[0]; got.ProductID != 1 || got.Quantity != 4 {
		t.Errorf("birinchi qator %+v; kutilgan ProductID 1, Quantity 4 (3 + 1) (TASK.md, 1-qadam)", got)
	}
	if got := o.Items[1]; got.ProductID != 2 || got.Quantity != 2 {
		t.Errorf("ikkinchi qator %+v; kutilgan ProductID 2, Quantity 2 (TASK.md, 1-qadam)", got)
	}
	if got := o.Total(); got != 5_600_000 {
		t.Errorf("Total() = %d; kutilgan 5600000 tiyin (olma 4 × 1 200 000 + non 2 × 400 000) (TASK.md, 1-qadam)", got)
	}
}

// 2-qadam: Cancel yangi va toʻlangan buyurtmani bekor qiladi, qolganlarida error qaytaradi.
func TestVazifaCancel(t *testing.T) {
	tests := []struct {
		name    string
		status  OrderStatus
		wantErr bool
		want    OrderStatus
	}{
		{"yangi", StatusNew, false, StatusCancelled},
		{"tolangan", StatusPaid, false, StatusCancelled},
		{"jonatilgan", StatusShipped, true, StatusShipped},
		{"bekor_qilingan", StatusCancelled, true, StatusCancelled},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			o := Order{ID: 7, Status: tc.status}
			err := o.Cancel()
			if tc.wantErr && err == nil {
				t.Errorf("«%v» buyurtmada Cancel() xato qaytarmadi; kutilgan: error (TASK.md, 2-qadam)", tc.status)
			}
			if !tc.wantErr && err != nil {
				t.Errorf("«%v» buyurtmada Cancel() = %v; kutilgan nil (TASK.md, 2-qadam)", tc.status, err)
			}
			if o.Status != tc.want {
				t.Errorf("Cancel()'dan keyin holat «%v»; kutilgan «%v» — receiver oʻzgarishni saqlaydimi? (TASK.md, 2-qadam)", o.Status, tc.want)
			}
		})
	}
}
