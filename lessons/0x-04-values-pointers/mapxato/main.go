// O'rtadagi sinov (GO-X.04): map'dagi struct maydoniga to'g'ridan-to'g'ri yozib bo'lmaydi. Ataylab build bo'lmaydi.
package main

import "fmt"

type User struct {
	ID      string
	Balance int
}

func main() {
	users := map[string]User{"42": {ID: "42", Balance: 100}}
	// JS: users.get("42").balance += 50 — ishlaydi
	users["42"].Balance += 50
	fmt.Println(users["42"].Balance)
}
