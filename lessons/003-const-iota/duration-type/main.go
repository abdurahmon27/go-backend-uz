package main

import (
	"fmt"
	"time"
)

func main() {
	// time paketida: type Duration int64 — asosi int64, lekin alohida tur
	var n int = 2
	var d time.Duration = n // ataylab xato: int → time.Duration yashirin oʻtmaydi
	fmt.Println(d)
}
