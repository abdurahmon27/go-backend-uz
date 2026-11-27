package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProduct(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/products/3", nil)
	rec := httptest.NewRecorder()
	productHandler(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("status %d, kutilgan %d: %s", rec.Code, http.StatusOK, rec.Body.String())
	}
}
