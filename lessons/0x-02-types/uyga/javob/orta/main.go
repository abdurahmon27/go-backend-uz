// Javob 🟡. map[string]int: JSON'da id va next son bo'lib chiqadi ("42" emas, 42).
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func main() {
	http.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, "id son emas", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]int{"id": id, "next": id + 1})
	})
	log.Println("listening on :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
