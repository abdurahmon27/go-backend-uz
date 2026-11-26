package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /salom", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "salom!")
	})

	for _, method := range []string{"GET", "POST"} {
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, httptest.NewRequest(method, "/salom", nil))
		fmt.Println(method, "/salom →", w.Code)
	}
}
