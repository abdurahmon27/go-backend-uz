package main

import (
	"fmt"
	"time"
)

func main() {
	n := 2                      // masalan, config yoki env'dan oʻqilgan soniyalar soni: tur = int
	time.Sleep(n * time.Second) // ataylab xato: int × time.Duration
	fmt.Println("tayyor")
}
