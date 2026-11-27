package main

// Katalog API testlari: har route serversiz, newMux orqali (GB-4.05), har holatga yangi katalog.

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/abdurahmon27/bozor/internal/store"
)

// testStore videodagi toʻrt mahsulotli katalog (bozor.json bilan bir xil, narx tiyinda).
func testStore(t *testing.T) *store.Store {
	t.Helper()
	st := &store.Store{}
	for _, p := range []store.Product{
		{Name: "olma", Price: 1_200_000, Stock: 40},
		{Name: "non", Price: 400_000, Stock: 3},
		{Name: "choy", Price: 2_500_000, Stock: 8},
		{Name: "oʻrik", Price: 3_500_000, Stock: 2},
	} {
		if _, err := st.Add(p.Name, p.Price, p.Stock); err != nil {
			t.Fatal(err)
		}
	}
	return st
}

// send soʻrovni serversiz h'ga beradi va recorder'ni qaytaradi (GB-4.03, GB-4.06).
func send(h http.Handler, method, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestRoutes(t *testing.T) {
	tests := []struct {
		name, method, path, body string
		status                   int
		want                     string // body'da boʻlishi kerak boʻlgan boʻlak
	}{
		{"katalog", "GET", "/products", "", 200, `"name":"choy"`},
		{"bitta", "GET", "/products/4", "", 200, `"name":"oʻrik"`},
		{"yoq_id", "GET", "/products/99", "", 404, `"code":"not_found"`},
		{"xato_id", "GET", "/products/abc", "", 400, `"code":"bad_id"`},
		{"yangi", "POST", "/products", `{"name":"anor","price":1800000,"stock":7}`, 201, `"id":5`},
		{"bosh_nom", "POST", "/products", `{"name":"","price":1800000}`, 422, `"name":"boʻsh boʻlmasin"`},
		{"uzun_nom", "POST", "/products", `{"name":"` + strings.Repeat("a", 21) + `","price":1}`, 422, `"validation_failed"`},
		{"buzuq_json", "POST", "/products", `{"name":`, 400, `"code":"bad_json"`},
		{"katta_body", "POST", "/products", `{"name":"` + strings.Repeat("a", 2000) + `"}`, 413, `"code":"too_large"`},
		{"juda_qimmat", "POST", "/products", `{"name":"oltin","price":1000000001}`, 422, `"validation_failed"`},
		{"buyurtma", "POST", "/orders", `{"product_id":1,"quantity":3}`, 201, `"total":3600000`},
		{"qoldiq_yetmaydi", "POST", "/orders", `{"product_id":4,"quantity":5}`, 409, `"code":"insufficient_stock"`},
		{"yoq_mahsulot", "POST", "/orders", `{"product_id":99,"quantity":1}`, 404, `"code":"not_found"`},
		{"nol_dona", "POST", "/orders", `{"product_id":1,"quantity":0}`, 422, `"quantity"`},
		{"health", "GET", "/health", "", 200, "ok"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := send(newMux(testStore(t)), tc.method, tc.path, tc.body)
			if rec.Code != tc.status {
				t.Errorf("%s %s: status %d, kutilgan %d: %s", tc.method, tc.path, rec.Code, tc.status, rec.Body)
			}
			if !strings.Contains(rec.Body.String(), tc.want) {
				t.Errorf("%s %s: body %q, unda %q yoʻq", tc.method, tc.path, rec.Body, tc.want)
			}
		})
	}
}

// POST /products 201 javobida Location sarlavhasi yangi mahsulotning manzilini beradi.
func TestCreateLocation(t *testing.T) {
	rec := send(newMux(testStore(t)), "POST", "/products", `{"name":"anor","price":1800000,"stock":7}`)
	if got := rec.Header().Get("Location"); got != "/products/5" {
		t.Errorf("Location %q, kutilgan /products/5", got)
	}
}

// Kutilmagan xato: 500, ichki matn (masalan fayl yoʻli) mijozga chiqmaydi.
func TestInternalError(t *testing.T) {
	log.SetOutput(io.Discard) // writeStoreError log'ga yozadi: test chiqishida koʻrinmasin
	t.Cleanup(func() { log.SetOutput(os.Stderr) })
	rec := httptest.NewRecorder()
	writeStoreError(rec, errors.New("open /srv/bozor/maxfiy.json: permission denied"))
	if rec.Code != 500 || strings.Contains(rec.Body.String(), "maxfiy") {
		t.Errorf("writeStoreError(kutilmagan xato): %d %s, kutilgan 500 va ichki matnsiz", rec.Code, rec.Body)
	}
}

// Buyurtma katalogni oʻzgartiradi: keyingi GET yangi qoldiqni koʻradi.
func TestOrderChangesStock(t *testing.T) {
	h := newMux(testStore(t))
	send(h, "POST", "/orders", `{"product_id":1,"quantity":3}`)
	if rec := send(h, "GET", "/products/1", ""); !strings.Contains(rec.Body.String(), `"stock":37`) {
		t.Errorf("buyurtmadan keyin GET /products/1: %s, kutilgan stock 37", rec.Body)
	}
}

// Bir vaqtda koʻp soʻrov: oddiy go test oʻtadi, data race'ni faqat go test -race koʻradi (GB-4.07).
func TestParallel(t *testing.T) {
	st := testStore(t)
	h := newMux(st)
	var wg sync.WaitGroup
	for range 30 {
		wg.Go(func() { send(h, "POST", "/orders", `{"product_id":1,"quantity":1}`) })
		wg.Go(func() { send(h, "GET", "/products/1", "") })
	}
	wg.Wait()
	if p, _ := st.Get(1); p.Stock != 10 {
		t.Errorf("30 ta buyurtmadan keyin olma qoldigʻi %d, kutilgan 10", p.Stock)
	}
}

// Server sozlamalari: Handler — newMux, ReadHeaderTimeout qoʻyilgan (GB-4.04).
func TestServer(t *testing.T) {
	srv := newServer(":0", testStore(t))
	if srv.ReadHeaderTimeout <= 0 {
		t.Errorf("ReadHeaderTimeout qoʻyilmagan")
	}
	if rec := send(srv.Handler, "GET", "/health", ""); rec.Code != 200 {
		t.Errorf("srv.Handler: GET /health status %d, kutilgan 200", rec.Code)
	}
}
