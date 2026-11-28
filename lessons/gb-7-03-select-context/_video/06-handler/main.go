package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

func checkStock(ctx context.Context, id int64) (int, error) {
	select {
	case n := <-supplier(id):
		return n, nil
	case <-ctx.Done():
		return 0, ctx.Err()
	}
}

func stockHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	n, err := checkStock(r.Context(), 3)
	if err != nil {
		fmt.Println("toʻxtadi:", err, time.Since(start).Round(time.Second))
		return
	}
	fmt.Println("javob yozildi", time.Since(start).Round(time.Second))
	fmt.Fprintln(w, "qoldiq:", n)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /stock", stockHandler)
	log.Fatal(http.ListenAndServe("localhost:8093", mux))
}
