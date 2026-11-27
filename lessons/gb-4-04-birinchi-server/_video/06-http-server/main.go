package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/products", catalogHandler)

	srv := &http.Server{
		Addr:              ":8765",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	fmt.Println("bozor API: http://localhost:8765")
	err := srv.ListenAndServe()
	fmt.Fprintln(os.Stderr, "xato:", err)
	os.Exit(1)
}
