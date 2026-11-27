package main

import (
	"fmt"
	"net/http"
	"os"
)

// Koʻp misollardagi shakl: http.HandleFunc va nil — DefaultServeMux (umumiy, global mux).
func main() {
	http.HandleFunc("/health", healthHandler)

	fmt.Println("bozor API: http://localhost:8765")
	err := http.ListenAndServe(":8765", nil)
	fmt.Fprintln(os.Stderr, "xato:", err)
	os.Exit(1)
}
