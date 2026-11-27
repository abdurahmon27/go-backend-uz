package main

import (
	"fmt"
	"net/http"
)

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintln(w, "topilmadi")
}

func main() {}
