package main

// Vazifani tekshiradi (go test .). Oʻzgartirmang.
// Soʻrovlar serversiz, newMux().ServeHTTP orqali yuboriladi: route roʻyxatdan oʻtmagan boʻlsa, mux 404 qaytaradi.

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func callStock(method, path string) (int, string, http.Header) {
	rec := httptest.NewRecorder()
	newMux().ServeHTTP(rec, httptest.NewRequest(method, path, nil))
	return rec.Code, rec.Body.String(), rec.Header()
}

func TestVazifaStock(t *testing.T) {
	for _, tc := range []struct {
		path, want string
	}{
		{"/products/1/stock", "40"},
		{"/products/3/stock", "5"},
	} {
		code, body, _ := callStock(http.MethodGet, tc.path)
		if code == http.StatusNotFound && strings.Contains(body, "404 page not found") {
			t.Fatalf("GET %s: mux «404 page not found» qaytardi: newMux'da \"GET /products/{id}/stock\" route'ini stockHandler bilan roʻyxatdan oʻtkazing (TASK.md, 1-qadam)", tc.path)
		}
		if code != http.StatusOK {
			t.Errorf("GET %s: status %d, kutilgan 200 (body: %q)", tc.path, code, body)
			continue
		}
		if got := strings.TrimSpace(body); got != tc.want {
			t.Errorf("GET %s: body «%s», kutilgan «%s»: qoldiqni fmt.Fprintln(w, p.Stock) bilan yozing (TASK.md, 2-qadam)", tc.path, got, tc.want)
		}
	}
}

func TestVazifaStockXato(t *testing.T) {
	code, body, _ := callStock(http.MethodGet, "/products/abc/stock")
	if code != http.StatusBadRequest {
		t.Errorf("GET /products/abc/stock: status %d, kutilgan 400: id son boʻlmasa http.Error(…, http.StatusBadRequest) va return (TASK.md, 2-qadam); body: %q", code, body)
	} else if body != "id butun son boʻlishi kerak\n" {
		t.Errorf("GET /products/abc/stock: body %q, kutilgan faqat %q: http.Error handler'ni toʻxtatmaydi, undan keyin return yozing (TASK.md, 2-qadam)", body, "id butun son boʻlishi kerak\n")
	}
	code, body, _ = callStock(http.MethodGet, "/products/99/stock")
	if code != http.StatusNotFound || !strings.Contains(body, "mahsulot topilmadi") {
		t.Errorf("GET /products/99/stock: %d %q, kutilgan 404 va «mahsulot topilmadi»: findProduct nil qaytarsa http.Error(…, http.StatusNotFound) (TASK.md, 2-qadam)", code, body)
	} else if body != "mahsulot topilmadi\n" {
		t.Errorf("GET /products/99/stock: body %q, kutilgan faqat %q: http.Error'dan keyin return yozing (TASK.md, 2-qadam)", body, "mahsulot topilmadi\n")
	}
}

func TestVazifaStockMetod(t *testing.T) {
	code, _, h := callStock(http.MethodPost, "/products/1/stock")
	if code != http.StatusMethodNotAllowed {
		t.Errorf("POST /products/1/stock: status %d, kutilgan 405: pattern metod bilan boshlansin, \"GET /products/{id}/stock\" (TASK.md, 3-qadam)", code)
		return
	}
	if allow := h.Get("Allow"); allow != "GET, HEAD" {
		t.Errorf("POST /products/1/stock: Allow «%s», kutilgan «GET, HEAD»: 405'ni handler'da qoʻlda yozmang, pattern'ga GET yozing, mux oʻzi javob beradi (TASK.md, 3-qadam)", allow)
	}
}
