// Xotira tugadi: fatal error, recover ishlamaydi.
// Demo uchun: ulimit -v 1000000 (taxminan 1 GB virtual xotira).
package main

import "fmt"

func main() {
	defer func() {
		r := recover()
		fmt.Println("qutqarildi:", r)
	}()

	buf := make([]byte, 2<<30) // 2 GiB
	fmt.Println(len(buf))
}
