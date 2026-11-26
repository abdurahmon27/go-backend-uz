package main

import (
	"fmt"
	"time"
)

func main() {
	// zona ko'rsatilmagan: time.Local (scratch'da jimgina UTC)
	fmt.Println(time.Now().Format("15:04 MST"))
}
