// defer tartibi: oxirgi qo'yilgan defer birinchi ishlaydi (LIFO).
package main

import "fmt"

func main() {
	fmt.Println("baza ochildi")
	defer fmt.Println("baza yopildi")

	fmt.Println("fayl ochildi")
	defer fmt.Println("fayl yopildi")

	fmt.Println("lock olindi")
	defer fmt.Println("lock qaytarildi")

	fmt.Println("ish bajarildi")
}
