// Bashorat 1 (GO-X.04): funksiya parametri ham qiymat berish — topUp user'ning nusxasini oladi.
package main

import "fmt"

type User struct {
	ID      string
	Balance int
}

func topUp(user User, amount int) {
	user.Balance += amount
}

func main() {
	user := User{ID: "42", Balance: 100}
	topUp(user, 50)           // JS: topUp(user) → 150
	fmt.Println(user.Balance) // 100
}
