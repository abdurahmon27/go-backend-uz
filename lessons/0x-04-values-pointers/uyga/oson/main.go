// Uy vazifasi 🟢 (GO-X.04): ikkala funksiya users'ni o'zgartirishi kerak, lekin hech narsa o'zgarmaydi.
//
//  1. Avval taxmin qiling, keyin dars papkasidan ishga tushiring: go run ./uyga/oson
//  2. deactivate'ni indeks bilan, rename'ni pointer bilan tuzating (main'dagi chaqiruv ham o'zgaradi).
//  3. Kutilgan natija:
//     true
//     [{42 Dilnoza false} {7 Bekzod true}]
//
// Ishora: golangci-lint'da govet'ning unusedwrite tekshiruvini yoqsangiz, ikkala joyni ham ko'rsatadi
// (README'dagi .golangci.yml namunasi).
// Javob: uyga/javob/oson/main.go
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

// deactivate id'si mos foydalanuvchini o'chiradi (Active = false).
func deactivate(id string) bool {
	for _, u := range users {
		if u.ID == id {
			u.Active = false
			return true
		}
	}
	return false
}

// rename foydalanuvchining nomini o'zgartiradi.
func rename(u User, name string) {
	u.Name = name
}

func main() {
	fmt.Println(deactivate("42"))
	rename(users[1], "Bekzod")
	fmt.Println(users)
}
