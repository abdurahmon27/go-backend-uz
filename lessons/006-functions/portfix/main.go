package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	port := 8080
	if s := os.Getenv("PORT"); s != "" {
		p, err := strconv.Atoi(s) // yangi nom
		if err != nil {
			fmt.Println("xato:", err)
			return
		}
		port = p // = : tashqi oʻzgaruvchiga yozamiz
		fmt.Println("sozlandi:", port)
	}
	fmt.Println("server ishga tushdi:", port)
}
