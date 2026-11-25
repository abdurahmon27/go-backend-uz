// Qiyin: fondagi goroutine'dagi panic butun serverni yiqitmasin.
package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

var logger = log.New(os.Stdout, "", 0)

// xavfsizGo f'ni yangi goroutine'da ishga tushiradi va uning panic'ini o'sha goroutine ichida ushlaydi.
func xavfsizGo(nom string, f func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Printf("%s: panic ushlandi: %v", nom, r)
			}
		}()
		f()
	}()
}

func main() {
	xavfsizGo("fon", func() {
		var hisob map[string]int
		hisob["fon"]++
	})
	time.Sleep(100 * time.Millisecond)
	fmt.Println("server tirik")
}
