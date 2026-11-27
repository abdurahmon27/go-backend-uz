package main

import (
	"fmt"
	"net/http"
)

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "topilmadi")
	w.WriteHeader(http.StatusNotFound)
}

func main() {}
