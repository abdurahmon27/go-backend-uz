// Pointer (GO-X.04): &user — user'ning manzili, *User — "User'ga pointer". topUp nusxa emas, manzil oladi.
package main

import "fmt"

type User struct {
	ID      string
	Balance int
}

func topUp(user *User, amount int) {
	user.Balance += amount // Go o'zi (*user).Balance deb o'qiydi
}

func main() {
	user := User{ID: "42", Balance: 100}
	topUp(&user, 50)          // JS'da havola yashirin, Go'da &
	fmt.Println(user.Balance) // 150
}
