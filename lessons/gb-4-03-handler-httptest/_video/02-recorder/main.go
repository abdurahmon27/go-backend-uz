package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "ok")
}

func main() {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	healthHandler(rec, req)
	fmt.Printf("%d %q\n", rec.Code, rec.Body.String())
}
