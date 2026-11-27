package main

// Vazifani tekshiradi (go test .). Oʻzgartirmang.
// POST /orders ni serversiz, mux orqali chaqiradi (GB-4.03): createOrder va validateOrder yozilgan boʻlishi kerak.

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func postOrder(body string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	newMux().ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader(body)))
	return rec
}

func TestVazifaBuyurtma(t *testing.T) {
	rec := postOrder(`{"product_id":1,"quantity":3}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("buyurtma: status %d, kutilgan %d (http.StatusCreated): decodeJSON, tekshiruv va writeJSON(w, http.StatusCreated, Order{…}) yozing; body: %q", rec.Code, http.StatusCreated, rec.Body.String())
	}
	if got := rec.Result().Header.Get("Content-Type"); got != "application/json" {
		t.Errorf("buyurtma: Content-Type «%s», kutilgan «application/json»: writeJSON ishlating", got)
	}
	var got Order
	dec := json.NewDecoder(rec.Body)
	if err := dec.Decode(&got); err != nil {
		t.Fatalf("buyurtma: body JSON emas: %v", err)
	}
	want := Order{ProductID: 1, Quantity: 3, Total: 3 * 1_200_000}
	if got != want {
		t.Errorf("buyurtma: %+v, kutilgan %+v (Total = narx × dona, tiyinda)", got, want)
	}
	if dec.More() {
		t.Errorf("buyurtma: body'da bitta JSON boʻlishi kerak: xato javobdan keyin return unutilgan boʻlishi mumkin")
	}
}

func TestVazifaTekshiruv(t *testing.T) {
	cases := []struct {
		name, body string
		fields     []string
	}{
		{"ikkalasi notoʻgʻri", `{"product_id":0,"quantity":0}`, []string{"product_id", "quantity"}},
		{"dona juda koʻp", `{"product_id":1,"quantity":101}`, []string{"quantity"}},
		{"dona manfiy", `{"product_id":2,"quantity":-1}`, []string{"quantity"}},
		{"maydonlar yoʻq", `{}`, []string{"product_id", "quantity"}},
	}
	for _, c := range cases {
		rec := postOrder(c.body)
		if rec.Code != http.StatusUnprocessableEntity {
			t.Errorf("%s: status %d, kutilgan %d (http.StatusUnprocessableEntity): qiymatlar notoʻgʻri, JSON esa toʻgʻri", c.name, rec.Code, http.StatusUnprocessableEntity)
			continue
		}
		e := errorBody(t, rec)
		if e.Code != "validation_failed" {
			t.Errorf("%s: error.code «%s», kutilgan «validation_failed»", c.name, e.Code)
		}
		for _, k := range c.fields {
			if e.Fields[k] == "" {
				t.Errorf("%s: error.fields[%q] boʻsh: validateOrder shu maydon uchun xabar qaytarsin", c.name, k)
			}
		}
	}
}

func TestVazifaTekshiruvJoyida(t *testing.T) {
	for _, in := range []orderInput{{1, 1}, {3, 100}} {
		if f := validateOrder(in); len(f) != 0 {
			t.Errorf("validateOrder(%+v) = %v, kutilgan boʻsh jadval: bu buyurtma toʻgʻri", in, f)
		}
	}
}

func TestVazifaTopilmadi(t *testing.T) {
	rec := postOrder(`{"product_id":999,"quantity":1}`)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("mahsulot katalogda yoʻq: status %d, kutilgan %d (http.StatusNotFound); body: %q", rec.Code, http.StatusNotFound, rec.Body.String())
	}
	if e := errorBody(t, rec); e.Code != "not_found" {
		t.Errorf("mahsulot katalogda yoʻq: error.code «%s», kutilgan «not_found»", e.Code)
	}
}

func TestVazifaBodyQoidalari(t *testing.T) {
	cases := []struct {
		name, body string
		status     int
		code       string
	}{
		{"nomaʼlum kalit", `{"product_id":1,"quantity":1,"discount":50}`, http.StatusBadRequest, "bad_json"},
		{"buzuq JSON", `{"product_id":`, http.StatusBadRequest, "bad_json"},
		{"katta body", `{"product_id":1,"quantity":1,"izoh":"` + strings.Repeat("a", 2000) + `"}`, http.StatusRequestEntityTooLarge, "too_large"},
	}
	for _, c := range cases {
		rec := postOrder(c.body)
		if rec.Code != c.status {
			t.Errorf("%s: status %d, kutilgan %d: createOrder body'ni decodeJSON bilan oʻqisin va xatoda return qilsin (videodagidek)", c.name, rec.Code, c.status)
			continue
		}
		// errorBody body'da bitta JSON borligini ham tekshiradi: decodeJSON false qaytarsa, darhol return kerak.
		if e := errorBody(t, rec); e.Code != c.code {
			t.Errorf("%s: error.code «%s», kutilgan «%s»", c.name, e.Code, c.code)
		}
	}
}
