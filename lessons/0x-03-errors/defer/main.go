// defer: chaqiruv funksiya tugaguncha kechiktiriladi. Lekin u faqat defer qatori bajarilganda ro'yxatga yoziladi:
// undan oldingi return'da hech narsa yopilmaydi (finally'dan farqi).
package main

import (
	"errors"
	"fmt"
)

func load(fail bool) error {
	fmt.Println("ochildi")
	if fail {
		return errors.New("buzuq JSON")
	}
	defer fmt.Println("yopildi")
	fmt.Println("o'qildi")
	return nil
}

func main() {
	if err := load(false); err != nil {
		fmt.Println("xato:", err)
	}
	if err := load(true); err != nil {
		fmt.Println("xato:", err)
	}
}
