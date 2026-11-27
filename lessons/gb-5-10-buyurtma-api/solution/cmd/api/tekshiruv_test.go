// VAZIFA tekshiruvi (TASK.md, 2–3-qadam): GET /orders/{id}. Bu faylni oʻzgartirmang.
package main

import (
	"strings"
	"testing"
)

func TestGetOrder(t *testing.T) {
	tests := []struct {
		name, path string
		status     int
		want       string
	}{
		{"bor", "/orders/1", 200, `{"id":1,"customer":"Dilnoza","status":"yangi","items":[{"product_id":1,"quantity":2,"unit_price":1200000},{"product_id":2,"quantity":3,"unit_price":400000}],"total":3600000}`},
		{"qatorsiz", "/orders/3", 200, `{"id":3,"customer":"Malika","status":"yangi","items":[],"total":0}`},
		{"yoq", "/orders/99", 404, `"code":"not_found"`},
		{"son_emas", "/orders/abc", 400, `"code":"bad_id"`},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h, _ := testAPI(t)
			rec := send(h, "GET", tc.path, "")
			if rec.Code != tc.status {
				t.Fatalf("GET %s: status %d, kutilgan %d (TASK.md, 2–3-qadam): %s", tc.path, rec.Code, tc.status, rec.Body)
			}
			if !strings.Contains(rec.Body.String(), tc.want) {
				t.Errorf("GET %s: body'da\n  %s\nyoʻq:\n  %s", tc.path, tc.want, rec.Body)
			}
		})
	}
}

// POST /orders javobidagi Location manzili GET bilan oʻsha buyurtmani qaytarishi kerak.
func TestOrderLocation(t *testing.T) {
	h, _ := testAPI(t)
	rec := send(h, "POST", "/orders", `{"customer":"Aziza","items":[{"product_id":3,"quantity":2}]}`)
	if rec.Code != 201 {
		t.Fatalf("POST /orders: status %d: %s", rec.Code, rec.Body)
	}
	loc := rec.Header().Get("Location")
	got := send(h, "GET", loc, "")
	if got.Code != 200 || got.Body.String() != rec.Body.String() {
		t.Errorf("GET %s: status %d, body %s\nkutilgan: 200 va POST javobi: %s (TASK.md, 2-qadam)", loc, got.Code, got.Body, rec.Body)
	}
}
