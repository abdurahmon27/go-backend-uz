package main

// Katalog API testlari: har route serversiz, newMux orqali (GB-4.05), har holatga yangi katalog.

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
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
// Body boʻlsa, Content-Type: application/json qoʻyadi (GB-6.01: requireJSON shuni talab qiladi).
func send(h http.Handler, method, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
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
		{"yoq_id", "GET", "/products/99", "", 404, `"title":"Not Found"`},
		{"xato_id", "GET", "/products/abc", "", 400, `"detail":"id butun son boʻlishi kerak"`},
		{"yangi", "POST", "/products", `{"name":"anor","price":1800000,"stock":7}`, 201, `"id":5`},
		{"buzuq_json", "POST", "/products", `{"name":`, 400, `"title":"Bad Request"`},
		{"narx_matn", "POST", "/products", `{"name":"choy","price":"25000"}`, 400, `"detail":"price: qiymat turi notoʻgʻri"`},
		{"ortiqcha_kalit", "POST", "/products", `{"name":"choy","rang":"qora"}`, 400, `"status":400`},
		{"buyurtma", "POST", "/orders", `{"product_id":1,"quantity":3}`, 201, `"total":3600000`},
		{"qoldiq_yetmaydi", "POST", "/orders", `{"product_id":4,"quantity":5}`, 409, `"type":"/problems/insufficient-stock"`},
		{"yoq_mahsulot", "POST", "/orders", `{"product_id":99,"quantity":1}`, 404, `"status":404`},
		{"health", "GET", "/health", "", 200, "ok"},
	}
	captureLog(t) // writeError har xatoni log'ga yozadi: test chiqishida koʻrinmasin
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

// Har xato javobi bitta shaklda: application/problem+json, body'dagi status = HTTP status,
// va ichki matn (Go tur nomlari, store matni) yoʻq.
func TestProblemShape(t *testing.T) {
	captureLog(t)
	h := newHandler(testStore(t))
	for _, c := range []struct{ method, path, body string }{
		{"GET", "/products/99", ""},
		{"GET", "/products/abc", ""},
		{"POST", "/products", `{"name":"choy","price":"25000"}`},
		{"POST", "/products", `[1]`},
		{"POST", "/orders", `{"product_id":4,"quantity":5}`},
		{"GET", "/narsalar", ""},  // mux'ning oʻz 404'i
		{"DELETE", "/health", ""}, // mux'ning oʻz 405'i
	} {
		rec := send(h, c.method, c.path, c.body)
		if ct := rec.Header().Get("Content-Type"); ct != "application/problem+json" {
			t.Errorf("%s %s: Content-Type %q, kutilgan application/problem+json", c.method, c.path, ct)
		}
		var p problem
		if err := json.Unmarshal(rec.Body.Bytes(), &p); err != nil {
			t.Fatalf("%s %s: body problem JSON emas: %v: %s", c.method, c.path, err, rec.Body)
		}
		if p.Status != rec.Code || p.Title == "" || p.Type == "" {
			t.Errorf("%s %s: status %d, problem %+v: status teng, type va title boʻsh boʻlmasin", c.method, c.path, rec.Code, p)
		}
		if p.RequestID == "" || p.RequestID != rec.Header().Get("X-Request-Id") {
			t.Errorf("%s %s: request_id %q, X-Request-Id %q: teng boʻlsin", c.method, c.path, p.RequestID, rec.Header().Get("X-Request-Id"))
		}
		for _, leak := range []string{"Go struct", "Go value", "productInput", "main.", "int64", "strconv", "mahsulot 99"} {
			if strings.Contains(rec.Body.String(), leak) {
				t.Errorf("%s %s: body'da ichki matn %q: %s", c.method, c.path, leak, rec.Body)
			}
		}
	}
	if rec := send(h, "DELETE", "/health", ""); rec.Header().Get("Allow") == "" {
		t.Errorf("405 javobida Allow sarlavhasi yoʻqoldi")
	}
}

// POST /products 201 javobida Location sarlavhasi yangi mahsulotning manzilini beradi.
func TestCreateLocation(t *testing.T) {
	rec := send(newMux(testStore(t)), "POST", "/products", `{"name":"anor","price":1800000,"stock":7}`)
	if got := rec.Header().Get("Location"); got != "/products/5" {
		t.Errorf("Location %q, kutilgan /products/5", got)
	}
}

// Kutilmagan xato: 500, ichki matn (masalan fayl yoʻli) mijozga chiqmaydi, log'ga esa toʻliq yoziladi.
func TestInternalError(t *testing.T) {
	logs := captureLog(t)
	rec := httptest.NewRecorder()
	writeError(rec, errors.New("open /srv/bozor/maxfiy.json: permission denied"))
	if rec.Code != 500 || strings.Contains(rec.Body.String(), "maxfiy") {
		t.Errorf("writeError(kutilmagan xato): %d %s, kutilgan 500 va ichki matnsiz", rec.Code, rec.Body)
	}
	if !strings.Contains(logs.String(), "maxfiy.json: permission denied") {
		t.Errorf("log'da toʻliq xato yoʻq: %q", logs.String())
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
	captureLog(t)
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

// Server sozlamalari: Handler — middleware zanjiri (newHandler), ReadHeaderTimeout qoʻyilgan (GB-4.04).
func TestServer(t *testing.T) {
	captureLog(t)
	srv := newServer(":0", testStore(t))
	if srv.ReadHeaderTimeout <= 0 {
		t.Errorf("ReadHeaderTimeout qoʻyilmagan")
	}
	if rec := send(srv.Handler, "GET", "/health", ""); rec.Code != 200 {
		t.Errorf("srv.Handler: GET /health status %d, kutilgan 200", rec.Code)
	} else if rec.Header().Get("X-Request-Id") == "" {
		t.Errorf("srv.Handler: javobda X-Request-Id yoʻq — Handler'ga newHandler(st) berilsin, newMux emas")
	}
}

// 3-qadam javobi (TASK.md): DELETE uchun oʻz testimiz (GB-1.08 jadval shakli).
func TestDeleteProduct(t *testing.T) {
	captureLog(t)
	tests := []struct {
		name   string
		path   string
		status int
	}{
		{"bor", "/products/3", 204},
		{"yoq", "/products/99", 404},
		{"xato_id", "/products/abc", 400},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h := newMux(testStore(t))
			if rec := send(h, "DELETE", tc.path, ""); rec.Code != tc.status {
				t.Errorf("DELETE %s: status %d, kutilgan %d: %s", tc.path, rec.Code, tc.status, rec.Body)
			}
		})
	}
	// oʻchirilgan mahsulot GET'da ham yoʻq
	h := newMux(testStore(t))
	send(h, "DELETE", "/products/3", "")
	if rec := send(h, "GET", "/products/3", ""); rec.Code != 404 {
		t.Errorf("DELETE'dan keyin GET /products/3: status %d, kutilgan 404", rec.Code)
	}
}
