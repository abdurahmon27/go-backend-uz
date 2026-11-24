package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	port := 8080
	if s := os.Getenv("PORT"); s != "" {
		port, err := strconv.Atoi(s)
		if err != nil {
			fmt.Println("xato:", err)
			return
		}
		fmt.Println("sozlandi:", port)
	}
	fmt.Println("server ishga tushdi:", port)
}
