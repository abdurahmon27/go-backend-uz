// Hamma goroutine uxlab qoldi: fatal error, recover ishlamaydi.
package main

import "fmt"

func main() {
	defer func() {
		r := recover()
		fmt.Println("qutqarildi:", r)
	}()

	natija := make(chan int)
	<-natija // hech kim yozmaydi
}
