// 🟢 Oson: PORT muhit oʻzgaruvchisidan oʻqilsin (boʻsh boʻlsa 3000).
// Keyin: GOOS=linux GOARCH=arm64 go build -o api-arm ./uyga/oson && go version -m api-arm | grep GOARCH
// Ishora: os.Getenv("PORT") — Node’dagi process.env.PORT. := va if keyingi darsda (X.02), hozircha shunday yozing:
//
//	port := os.Getenv("PORT")
//	if port == "" { port = "3000" }
//
// Javob: ../javob/oson/main.go
package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})
	http.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": r.PathValue("id")})
	})
	// TODO: port := ...
	log.Println("listening on :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
