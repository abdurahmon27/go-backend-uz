// VAZIFA tekshiruvi (TASK.md, 1-qadam): store.Order. Bu faylni oʻzgartirmang.
package store_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/abdurahmon27/bozor/internal/pgtest"
	"github.com/abdurahmon27/bozor/internal/store"
)

func TestOrder(t *testing.T) {
	st := store.New(pgtest.New(t))
	o, err := st.Order(testCtx(t), 1)
	if err != nil {
		t.Fatalf("Order(1) = %v; kutilgan: seed.sql'dagi Dilnoza buyurtmasi (TASK.md, 1-qadam)", err)
	}
	want := store.Order{
		ID: 1, Customer: "Dilnoza", Status: "yangi",
		Items: []store.OrderItem{
			{ProductID: 1, Quantity: 2, UnitPrice: 1200000},
			{ProductID: 2, Quantity: 3, UnitPrice: 400000},
		},
		Total: 3600000,
	}
	if !reflect.DeepEqual(o, want) {
		t.Errorf("Order(1) =\n  %+v\nkutilgan:\n  %+v\n(qatorlar order_items.id tartibida, Total = Σ quantity × unit_price)", o, want)
	}
}

// seed.sql'dagi 3-buyurtmada (Malika) qator yoʻq: Items boʻsh slice boʻlsin (JSON'da [], null emas).
func TestOrderWithoutItems(t *testing.T) {
	st := store.New(pgtest.New(t))
	o, err := st.Order(testCtx(t), 3)
	if err != nil {
		t.Fatalf("Order(3) = %v; kutilgan: Malikaning qatorsiz buyurtmasi (TASK.md, 1-qadam)", err)
	}
	if o.Items == nil || len(o.Items) != 0 || o.Total != 0 {
		t.Errorf("Order(3).Items = %#v, Total = %d; kutilgan: []store.OrderItem{} (nil emas) va 0 (TASK.md: JSON'da [] boʻlsin)", o.Items, o.Total)
	}
}

func TestOrderNotFound(t *testing.T) {
	st := store.New(pgtest.New(t))
	if _, err := st.Order(testCtx(t), 99); !errors.Is(err, store.ErrOrderNotFound) {
		t.Errorf("Order(99) = %v; kutilgan: ErrOrderNotFound (TASK.md, 1-qadam)", err)
	}
}

// PlaceOrder yozgan buyurtma Order bilan aynan shunday oʻqilishi kerak.
func TestOrderRoundTrip(t *testing.T) {
	st := store.New(pgtest.New(t))
	placed, err := st.PlaceOrder(testCtx(t), "Aziza", []store.OrderItem{{ProductID: 3, Quantity: 2}})
	if err != nil {
		t.Fatalf("PlaceOrder: %v", err)
	}
	got, err := st.Order(testCtx(t), placed.ID)
	if err != nil {
		t.Fatalf("Order(%d) = %v; kutilgan: hozir yozilgan buyurtma (TASK.md, 1-qadam)", placed.ID, err)
	}
	if !reflect.DeepEqual(got, placed) {
		t.Errorf("Order(%d) =\n  %+v\nPlaceOrder qaytargani:\n  %+v", placed.ID, got, placed)
	}
}
