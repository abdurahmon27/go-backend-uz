// Tuzatilgan hook: Unmarshal xatosi tekshiriladi, buzuq config bilan xizmat ishga tushmaydi.
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
	if err := json.Unmarshal(data, &cfg); err != nil {
		log.Fatal(err)
	}

	addr := fmt.Sprintf(":%d", cfg.Port)
	http.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	log.Println("listening on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
