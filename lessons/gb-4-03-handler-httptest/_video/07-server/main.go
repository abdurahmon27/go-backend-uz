package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

func main() {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	}))
	defer srv.Close()
	resp, err := http.Get(srv.URL)
	if err != nil {
		fmt.Println("xato:", err)
		return
	}
	defer resp.Body.Close()
	fmt.Println(resp.Status)
}
