// O'rtadagi sinov (GO-X.04): &users[i] — slice elementining o'z manzili (C variantdagi &user esa nusxaning manzili).
package main

import "fmt"

type User struct {
	ID      string
	Balance int
}

func main() {
	users := []User{{ID: "42", Balance: 100}}
	for i := range users {
		p := &users[i]
		p.Balance += 50
	}
	fmt.Println(users[0].Balance) // 150
}
