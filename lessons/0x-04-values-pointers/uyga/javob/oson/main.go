// Javob 🟢 (GO-X.04): deactivate — indeks orqali elementning o'zi; rename — pointer (&users[1]).
package main

import "fmt"

type User struct {
	ID     string
	Name   string
	Active bool
}

var users = []User{
	{ID: "42", Name: "Dilnoza", Active: true},
	{ID: "7", Name: "Jasur", Active: true},
}

func deactivate(id string) bool {
	for i := range users {
		if users[i].ID == id {
			users[i].Active = false // nusxa emas, slice'dagi element
			return true
		}
	}
	return false
}

func rename(u *User, name string) {
	u.Name = name // u — pointer: chaqiruvchining qiymatiga yozadi
}

func main() {
	fmt.Println(deactivate("42"))
	rename(&users[1], "Bekzod") // slice elementining manzili
	fmt.Println(users)
}
