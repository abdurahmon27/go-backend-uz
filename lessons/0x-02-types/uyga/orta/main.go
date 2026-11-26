// Uy vazifasi 🟡: oldingi darsdagi (GO-X.01) API. r.PathValue("id") string qaytaradi.
// Vazifa: id ni int'ga aylantiring va {"id": 42, "next": 43} qaytaring.
// id son bo'lmasa (masalan /users/abc), 400 Bad Request qaytaring.
// Yordam: strconv.Atoi(s) ikki qiymat qaytaradi: son va xato. Xato tekshiruvi keyingi darsda (GO-X.03),
// hozircha shu qolipni ishlating:
//
//	id, err := strconv.Atoi(r.PathValue("id"))
//	if err != nil {
//		http.Error(w, "id son emas", http.StatusBadRequest)
//		return
//	}
//
// Ishga tushirish: go run uyga/orta/main.go, keyin boshqa terminalda:
//
//	curl -i localhost:3000/users/42
//	curl -i localhost:3000/users/abc
package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// TODO: id ni int qiling, "next" = id + 1
		json.NewEncoder(w).Encode(map[string]string{"id": r.PathValue("id")})
	})
	log.Println("listening on :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
