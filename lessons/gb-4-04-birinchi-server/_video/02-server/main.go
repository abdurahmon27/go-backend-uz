package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/products", catalogHandler)

	fmt.Println("bozor API: http://localhost:8765")
	err := http.ListenAndServe(":8765", mux)
	fmt.Fprintln(os.Stderr, "xato:", err)
	os.Exit(1)
}
