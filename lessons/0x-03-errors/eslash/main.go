// Takrorlash (GO-X.02): := chap tomonida kamida bitta yangi nom bo'lsa yetarli, qolganlari qiymat oladi.
package main

import (
	"fmt"
	"strconv"
)

func main() {
	a, err := strconv.Atoi("42")
	b, err := strconv.Atoi("x") // b yangi, err esa qayta yoziladi
	fmt.Println(a, b, err)
}
