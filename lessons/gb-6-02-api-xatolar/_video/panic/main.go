// Recall uchun (GB-6.01): recovery middleware'siz handler panic qilsa nima boʻladi.
package main

import "net/http"

func main() {
	http.HandleFunc("GET /boom", func(w http.ResponseWriter, r *http.Request) {
		var m map[string]int
		m["x"] = 1 // nil map'ga yozish: panic
	})
	http.ListenAndServe(":8766", nil)
}
