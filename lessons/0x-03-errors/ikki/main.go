// Bashorat: os.ReadFile ikki qiymat qaytaradi, biz bittasini oldik. Ataylab build bo'lmaydi.
package main

import (
	"fmt"
	"os"
)

func main() {
	data := os.ReadFile("config.json")
	fmt.Println(len(data))
}
