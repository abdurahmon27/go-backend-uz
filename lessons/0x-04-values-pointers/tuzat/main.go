// Tuzatilgan hook (GO-X.04): range indeks bilan, element users[i] orqali o'zgartiriladi.
// ⚠️ Demo: users global, so'rovlar esa parallel goroutine'larda — haqiqiy xizmatda sync.Mutex kerak (GO-1.09).
// Ishga tushirish: go run ./tuzat   (boshqa terminalda: curl -X POST localhost:3000/users/42/topup; curl localhost:3000/users/42)
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
		for i := range users {
			if users[i].ID == r.PathValue("id") {
				users[i].Balance += 50
				json.NewEncoder(w).Encode(users[i])
				return
			}
		}
		http.NotFound(w, r)
	})
	log.Println("listening on :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
