// defer amalda: xizmat boshqa xizmatdan (users API) foydalanuvchini HTTP orqali oladi.
// resp.Body — tarmoq ulanishidan oqim; uni yopish shart. Demo uchun users API'ni httptest bilan shu dasturda ko'taramiz.
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
	err = json.NewDecoder(resp.Body).Decode(&u)
	return u, err
}

func main() {
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"id": 42, "name": "Dilnoza"}`)
	}))
	defer api.Close()

	u, err := fetchUser(api.URL, "42")
	fmt.Println(u, err)

	_, err = fetchUser("http://localhost:9", "42") // 9-port: hech kim tinglamaydi
	fmt.Println(err)
}
