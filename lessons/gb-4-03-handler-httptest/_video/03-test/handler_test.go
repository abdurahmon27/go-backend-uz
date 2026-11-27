package main

import (
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
