// GO-X.03 hook: GO-X.01 xizmati config.json'ni o'qiydi. Unmarshal'ning xatosi tekshirilmagan.
// Ishga tushirish (dars papkasidan): go run ./hook   (config.json'da ortiqcha vergul bor)
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Config struct {
	Port int `json:"port"`
}

func main() {
	data, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}
	var cfg Config
	json.Unmarshal(data, &cfg)

	addr := fmt.Sprintf(":%d", cfg.Port)
	http.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	log.Println("listening on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
