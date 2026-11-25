// log.Fatal ichida os.Exit bor: defer'lar ishlamaydi.
package main

import (
	"fmt"
	"log"
)

func main() {
	defer fmt.Println("defer: buferdagi loglar yozildi")

	if err := sozla(); err != nil {
		log.Fatal("sozlash: ", err)
	}
}

func sozla() error {
	return fmt.Errorf("PORT berilmagan")
}
