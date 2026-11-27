package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/products", catalogHandler)

	fmt.Println("bozor API: http://localhost:8765")
	http.ListenAndServe(":8765", mux) // XATO: qaytgan error tashlab yuborildi
}
