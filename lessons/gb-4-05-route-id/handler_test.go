package main

// Videoda yozilgan testlar: soʻrov serversiz, lekin haqiqiy route'lar orqali (newMux().ServeHTTP) oʻtadi.
// Starterda ham oʻtadi.

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// serve soʻrovni mux'ga beradi: mux pattern'ni tanlaydi, PathValue'ni toʻldiradi va handler'ni chaqiradi.
func serve(method, path string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	rec := httptest.NewRecorder()
	newMux().ServeHTTP(rec, req)
	return rec
}

func TestProduct(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/products/3", nil)
	rec := httptest.NewRecorder()
	newMux().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("status %d, kutilgan %d: %s", rec.Code, http.StatusOK, rec.Body.String())
	}
}

func TestProductXato(t *testing.T) {
	if rec := serve(http.MethodGet, "/products/abc"); rec.Code != http.StatusBadRequest {
		t.Errorf("/products/abc: status %d, kutilgan %d", rec.Code, http.StatusBadRequest)
	}
	rec := serve(http.MethodGet, "/products/99")
	if rec.Code != http.StatusNotFound || !strings.Contains(rec.Body.String(), "mahsulot topilmadi") {
		t.Errorf("/products/99: %d %q, kutilgan 404 va «mahsulot topilmadi»", rec.Code, rec.Body.String())
	}
}

func TestMetod(t *testing.T) {
	rec := serve(http.MethodPost, "/products/3")
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("POST /products/3: status %d, kutilgan %d", rec.Code, http.StatusMethodNotAllowed)
	}
	if got := rec.Header().Get("Allow"); got != "GET, HEAD" {
		t.Errorf("POST /products/3: Allow %q, kutilgan %q", got, "GET, HEAD")
	}
}

func TestHome(t *testing.T) {
	if rec := serve(http.MethodGet, "/"); rec.Body.String() != "bozor API\n" {
		t.Errorf("/: body %q, kutilgan %q", rec.Body.String(), "bozor API\n")
	}
	if rec := serve(http.MethodGet, "/produkts/3"); rec.Code != http.StatusNotFound {
		t.Errorf("/produkts/3: status %d, kutilgan %d (\"/\" pattern'i hamma yoʻlni yutmasin: \"GET /{$}\")", rec.Code, http.StatusNotFound)
	}
}
