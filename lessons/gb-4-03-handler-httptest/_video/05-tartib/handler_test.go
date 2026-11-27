package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNotFound(t *testing.T) {
	rec := httptest.NewRecorder()
	notFoundHandler(rec, httptest.NewRequest(http.MethodGet, "/yoq", nil))
	if rec.Code != http.StatusNotFound {
		t.Errorf("status %d, kutilgan %d", rec.Code, http.StatusNotFound)
	}
}
