// API testlari: haqiqiy handler'lar, haqiqiy mux, haqiqiy Postgres (mock yoʻq). Server ishga tushmaydi:
// httptest soʻrovni toʻgʻridan-toʻgʻri mux'ga beradi (GB-4.03). Har qator oʻz sxemasida (pgtest.New).
package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/abdurahmon27/bozor/internal/pgtest"
	"github.com/abdurahmon27/bozor/internal/store"
	"github.com/jackc/pgx/v5/pgxpool"
)

// send soʻrovni h'ga beradi va javobni (recorder) qaytaradi (GB-4.06, GB-4.08).
func send(h http.Handler, method, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

// testAPI — yangi sxema (seed.sql bilan) va shu bazaga ulangan mux.
func testAPI(t *testing.T) (http.Handler, *pgxpool.Pool) {
	t.Helper()
	pool := pgtest.New(t)
	return newMux(store.New(pool)), pool
}

func TestRoutes(t *testing.T) {
	tests := []struct {
		name, method, path, body string
		status                   int
		want                     string
	}{
		{"katalog", "GET", "/products", "", 200, `"name":"oʻrik"`},
		{"bitta", "GET", "/products/4", "", 200, `"stock":2`},
		{"yoq_id", "GET", "/products/99", "", 404, `"code":"not_found"`},
		{"son_emas", "GET", "/products/abc", "", 400, `"code":"bad_id"`},
		{"yangi_mahsulot", "POST", "/products", `{"name":"anor","price":2000000,"stock":5}`, 201, `"id":5`},
		{"takror_nom", "POST", "/products", `{"name":"olma","price":100,"stock":1}`, 409, `"code":"conflict"`},
		{"buyurtma", "POST", "/orders", `{"customer":"Aziza","items":[{"product_id":1,"quantity":2},{"product_id":4,"quantity":1}]}`, 201, `"total":5900000`},
		{"qoldiq_yetmaydi", "POST", "/orders", `{"customer":"Aziza","items":[{"product_id":1,"quantity":2},{"product_id":4,"quantity":5}]}`, 409, `"message":"oʻrik: qoldiq 2, soʻralgan 5"`},
		{"yoq_mahsulot", "POST", "/orders", `{"customer":"Aziza","items":[{"product_id":99,"quantity":1}]}`, 404, `"code":"not_found"`},
		{"nol_dona", "POST", "/orders", `{"customer":"Aziza","items":[{"product_id":1,"quantity":0}]}`, 422, `"items[0].quantity"`},
		{"boʻsh_savat", "POST", "/orders", `{"customer":"Aziza","items":[]}`, 422, `"items"`},
		{"takror_qator", "POST", "/orders", `{"customer":"Aziza","items":[{"product_id":1,"quantity":1},{"product_id":1,"quantity":2}]}`, 422, `"items[1].product_id"`},
		{"notoʻgʻri_json", "POST", "/orders", `{"customer":`, 400, `"code":"bad_json"`},
		{"oʻchirish", "DELETE", "/products/3", "", 204, ""},
		{"buyurtmadagi", "DELETE", "/products/1", "", 409, `"code":"conflict"`},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h, _ := testAPI(t)
			rec := send(h, tc.method, tc.path, tc.body)
			if rec.Code != tc.status {
				t.Errorf("%s %s: status %d, kutilgan %d: %s", tc.method, tc.path, rec.Code, tc.status, rec.Body)
			}
			if !strings.Contains(rec.Body.String(), tc.want) {
				t.Errorf("%s %s: body'da %s yoʻq: %s", tc.method, tc.path, tc.want, rec.Body)
			}
		})
	}
}

// 409 dan keyin baza oʻzgarmagan: birinchi qatordagi olma qoldigʻi joyida, yangi buyurtma yoʻq.
func TestOrderConflictLeavesNoTrace(t *testing.T) {
	h, pool := testAPI(t)
	rec := send(h, "POST", "/orders", `{"customer":"Aziza","items":[{"product_id":1,"quantity":2},{"product_id":4,"quantity":5}]}`)
	if rec.Code != http.StatusConflict {
		t.Fatalf("status %d, kutilgan 409: %s", rec.Code, rec.Body)
	}
	if rec := send(h, "GET", "/products/1", ""); !strings.Contains(rec.Body.String(), `"stock":40`) {
		t.Errorf("olma: %s; kutilgan stock 40 (transaction bekor boʻlishi kerak)", rec.Body)
	}
	var n int
	if err := pool.QueryRow(t.Context(), "SELECT count(*) FROM orders").Scan(&n); err != nil || n != 3 {
		t.Errorf("orders: %d qator (%v), kutilgan 3", n, err)
	}
}
