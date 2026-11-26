// Javob 🟡 (GO-X.03): 200 bo'lmagan javob xato sifatida qaytariladi.
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func fetchUser(baseURL, id string) (User, error) {
	var u User
	resp, err := http.Get(baseURL + "/users/" + id)
	if err != nil {
		return u, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return u, fmt.Errorf("users API: %s", resp.Status)
	}
	err = json.NewDecoder(resp.Body).Decode(&u)
	return u, err
}

// usersAPI — demo uchun soxta users API: faqat 42-foydalanuvchi bor.
func usersAPI() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path != "/users/42" {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, `{"error": "not found"}`)
			return
		}
		fmt.Fprint(w, `{"id": 42, "name": "Dilnoza"}`)
	}))
}

func main() {
	api := usersAPI()
	defer api.Close()

	fmt.Println(fetchUser(api.URL, "42"))
	fmt.Println(fetchUser(api.URL, "7"))
}
