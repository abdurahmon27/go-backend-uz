// Uy vazifasi 🟡 (GO-X.04): PATCH /users/{id} foydalanuvchi nomini o'zgartirishi kerak. Fayl build bo'lmaydi.
//
//  1. go run ./uyga/orta → cannot assign to struct field users[id].Name in map
//  2. Ikki usulda tuzating va qaysi biri sizga yoqishini README'ga yozing:
//     a) o'qing → o'zgartiring → map'ga qaytarib yozing (u := users[id]; …; users[id] = u);
//     b) map[string]*User saqlang (unda GET'dagi qiymat ham pointer bo'ladi).
//  3. Tekshiring (boshqa terminalda):
//     curl -X PATCH -d '{"name":"Dilnoza A."}' localhost:3000/users/42  → {"id":"42","name":"Dilnoza A.","balance":100}
//     curl localhost:3000/users/42                                       → o'sha
//     curl -X PATCH -d '{"name":"X"}' localhost:3000/users/1             → 404 page not found
//
// ⚠️ net/http har so'rovni alohida goroutine'da bajaradi: haqiqiy xizmatda map'ga bir vaqtda yozish uchun
// sync.Mutex kerak (concurrent map writes — GO-1.09). Bu vazifada bitta curl bilan tekshiramiz.
// Javob: uyga/javob/orta/main.go
package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type User struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Balance int    `json:"balance"`
}

var users = map[string]User{"42": {ID: "42", Name: "Dilnoza", Balance: 100}}

func main() {
	http.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		u, ok := users[r.PathValue("id")] // ok = false: kalit yo'q
		if !ok {
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(u)
	})
	http.HandleFunc("PATCH /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		var body struct { // nomsiz struct: faqat shu handler uchun
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		id := r.PathValue("id")
		if _, ok := users[id]; !ok {
			http.NotFound(w, r)
			return
		}
		users[id].Name = body.Name
		json.NewEncoder(w).Encode(users[id])
	})
	log.Println("listening on :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
