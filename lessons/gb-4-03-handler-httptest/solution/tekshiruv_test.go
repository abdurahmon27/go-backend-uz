package main

// Vazifani tekshiradi (go test .). Oʻzgartirmang.
// catalogHandler'ni serversiz, httptest.NewRecorder bilan chaqiradi (GB-4.03).
// Sarlavhalar rec.Result().Header'dan olinadi: bu nusxa birinchi Write/WriteHeader paytida olinadi, shuning uchun
// Write'dan KEYIN qoʻyilgan sarlavha unda koʻrinmaydi (rec.Header() esa jonli jadval, uni oʻzgartirsangiz koʻrinib qoladi).

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"
)

func callCatalog(method string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	catalogHandler(rec, httptest.NewRequest(method, "/products", nil))
	return rec
}

func TestVazifaGET(t *testing.T) {
	rec := callCatalog(http.MethodGet)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET: status %d, kutilgan %d", rec.Code, http.StatusOK)
	}
	if got := rec.Result().Header.Get("Content-Type"); got != "application/json" {
		t.Errorf("GET: Content-Type «%s», kutilgan «application/json»: sarlavhani Encode'dan OLDIN qoʻying", got)
	}
	var got []Product
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("GET: body JSON massiv emas (%v): catalog'ni json.NewEncoder(w).Encode bilan yozing, hozirgi body: %q", err, rec.Body.String())
	}
	if !slices.Equal(got, catalog) {
		t.Errorf("GET: javobdagi mahsulotlar %v, kutilgan catalog: %v", got, catalog)
	}
}

func TestVazifaMetod(t *testing.T) {
	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodDelete} {
		t.Run(method, func(t *testing.T) {
			rec := callCatalog(method)
			if rec.Code != http.StatusMethodNotAllowed {
				t.Errorf("%s: status %d, kutilgan 405 (http.StatusMethodNotAllowed): WriteHeader'ni chaqiring va return qiling", method, rec.Code)
			}
			if got := rec.Result().Header.Get("Allow"); got != "GET" {
				t.Errorf("%s: Allow «%s», kutilgan «GET»: sarlavhani WriteHeader'dan OLDIN qoʻying", method, got)
			}
			if rec.Body.Len() != 0 {
				t.Errorf("%s: 405 javobda body boʻlmasligi kerak, hozir %q", method, rec.Body.String())
			}
		})
	}
}

func TestVazifaBosh(t *testing.T) {
	old := catalog
	t.Cleanup(func() { catalog = old })
	catalog = nil
	rec := callCatalog(http.MethodGet)
	if got := strings.TrimSpace(rec.Body.String()); got != "[]" {
		t.Errorf("catalog boʻsh: body «%s», kutilgan «[]» (nil slice JSON'da null boʻladi: boʻsh, lekin nil boʻlmagan slice kerak)", got)
	}
}
