// Bashorat 2 (GO-X.04): uchta sikl — qaysi biri users'ni o'zgartiradi?
package main

import "fmt"

type User struct {
	ID      string
	Balance int
}

func fresh() []User { return []User{{ID: "42", Balance: 100}} }

func main() {
	users := fresh()
	for _, user := range users { // A: user — elementning nusxasi
		user.Balance += 50
	}
	fmt.Println("A:", users[0].Balance)

	users = fresh()
	for i := range users { // B: users[i] — elementning o'zi
		users[i].Balance += 50
	}
	fmt.Println("B:", users[0].Balance)

	users = fresh()
	for _, user := range users { // C: &user — NUSXAning manzili
		p := &user
		p.Balance += 50
	}
	fmt.Println("C:", users[0].Balance)
}
