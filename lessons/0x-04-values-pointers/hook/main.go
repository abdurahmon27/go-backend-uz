// Hook (GO-X.04): GO-X.01 xizmati + xotiradagi users ro'yxati. POST /users/{id}/topup balansni 50 ga oshirishi
// kerak. Javob 150 deydi, keyingi GET esa 100: range bergan user — elementning NUSXASI.
// Ishga tushirish: go run ./hook   (boshqa terminalda: curl -X POST localhost:3000/users/42/topup; curl localhost:3000/users/42)
package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type User struct {
	ID      string `json:"id"`
	Balance int    `json:"balance"`
}

var users = []User{{ID: "42", Balance: 100}}

func main() {
	http.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		for _, user := range users {
			if user.ID == r.PathValue("id") {
				json.NewEncoder(w).Encode(user)
				return
			}
		}
		http.NotFound(w, r)
	})
	http.HandleFunc("POST /users/{id}/topup", func(w http.ResponseWriter, r *http.Request) {
		for _, user := range users {
			if user.ID == r.PathValue("id") {
				user.Balance += 50
				json.NewEncoder(w).Encode(user)
				return
			}
		}
		http.NotFound(w, r)
	})
	log.Println("listening on :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
