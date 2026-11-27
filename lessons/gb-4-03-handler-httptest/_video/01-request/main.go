package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

func main() {
	req := httptest.NewRequest(http.MethodGet, "/products/1", nil)
	fmt.Println(req.Method, req.URL.Path, req.Proto)
}
