// Tuzatish (GO-X.04): map'dagi struct — o'qing, o'zgartiring, qaytarib yozing; yoki map'da pointer saqlang.
package main

import "fmt"

type User struct {
	ID      string
	Balance int
}

func main() {
	users := map[string]User{"42": {ID: "42", Balance: 100}}
	user := users["42"] // nusxa
	user.Balance += 50
	users["42"] = user               // qaytarib yozamiz
	fmt.Println(users["42"].Balance) // 150

	ptrs := map[string]*User{"42": &User{ID: "42", Balance: 100}}
	ptrs["42"].Balance += 50        // pointer orqali: build bo'ladi
	fmt.Println(ptrs["42"].Balance) // 150
	// Diqqat: kalit yo'q bo'lsa ptrs["x"] nil, ptrs["x"].Balance esa panic beradi.
}
