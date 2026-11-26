// Uy vazifasi 🟡 (GO-X.03): users API noma'lum id uchun 404 qaytaradi, fetchUser esa buni sezmaydi.
//
// Ishga tushiring: go run ./uyga/orta
// Hozir:
//
//	{42 Dilnoza} <nil>
//	{0 } <nil>          ← 404 bo'ldi, lekin xato yo'q va bo'sh foydalanuvchi
//
// Sabab: Go'da (JS'dagi fetch kabi) 404 yoki 500 javobi xato emas: http.Get faqat tarmoq yoki protokol xatosida err qaytaradi
// (go doc net/http.Get: "A non-2xx response doesn't cause an error."). 404 javobining tanasi {"error": …}
// esa User'ga jim tushadi va maydonlar zero value'da qoladi (GO-X.02).
//
// Vazifa: defer'dan keyin resp.StatusCode'ni tekshiring. 200 bo'lmasa, xato qaytaring:
//
//	return u, fmt.Errorf("users API: %s", resp.Status)
//
// Kutilgan natija:
//
//	{42 Dilnoza} <nil>
//	{0 } users API: 404 Not Found
//
// ⭐ Qo'shimcha: GO-X.01 dagi GET /users/{id} handler'ini fetchUser bilan yozing: xato bo'lsa
// http.Error(w, err.Error(), http.StatusBadGateway) va return.
//
// Javob: uyga/javob/orta/main.go
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
	// TODO: resp.StatusCode 200 bo'lmasa, xato qaytaring
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
