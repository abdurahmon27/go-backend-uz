package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	healthHandler(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("status %d, kutilgan %d", rec.Code, http.StatusOK)
	}
	if got := rec.Body.String(); got != "ok\n" {
		t.Errorf("body %q, kutilgan %q", got, "ok\n")
	}
}

func TestProduct(t *testing.T) {
	rec := httptest.NewRecorder()
	productHandler(rec, httptest.NewRequest(http.MethodGet, "/products/1", nil))
	var got Product
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("body JSON emas: %v", err)
	}
	if got != olma {
		t.Errorf("mahsulot %v, kutilgan %v", got, olma)
	}
}

func TestProductMetod(t *testing.T) {
	rec := httptest.NewRecorder()
	productHandler(rec, httptest.NewRequest(http.MethodPost, "/products/1", nil))
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("status %d, kutilgan %d", rec.Code, http.StatusMethodNotAllowed)
	}
	if got := rec.Result().Header.Get("Allow"); got != "GET" {
		t.Errorf("Allow %q, kutilgan %q", got, "GET")
	}
}
