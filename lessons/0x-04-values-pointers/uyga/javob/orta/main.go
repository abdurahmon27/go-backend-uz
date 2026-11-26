// Javob 🟡 (GO-X.04): a) usul — o'qing, o'zgartiring, qaytarib yozing. b) usul pastda izohda.
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
		u, ok := users[id] // u — map'dagi qiymatning nusxasi
		if !ok {
			http.NotFound(w, r)
			return
		}
		u.Name = body.Name
		users[id] = u // nusxani qaytarib yozamiz
		json.NewEncoder(w).Encode(u)
		// b) usul: var users = map[string]*User{…}; u, ok := users[id]; …; u.Name = body.Name
		// (u pointer bo'lgani uchun qaytarib yozish shart emas; sync.Mutex esa ikkala usulda ham kerak)
	})
	log.Println("listening on :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
