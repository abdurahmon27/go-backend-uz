// Qiymat nusxasi (GO-X.04): c := user butun struct'ni (hamma maydonlarini) yangi joyga nusxalaydi.
// JS'da const c = user — o'sha obyektga ikkinchi havola bo'lardi.
package main

import "fmt"

type User struct {
	ID      string
	Balance int
}

func main() {
	user := User{ID: "42", Balance: 100}
	c := user
	c.Balance += 50
	fmt.Println(user.Balance, c.Balance) // 100 150
}
