package main

import (
	"fmt"
	"time"
	_ "time/tzdata"
)

func main() {
	loc, err := time.LoadLocation("Asia/Tashkent")
	if err != nil {
		fmt.Println("xato:", err)
		return
	}
	fmt.Println("Toshkent soati:", time.Now().In(loc).Hour())
}
