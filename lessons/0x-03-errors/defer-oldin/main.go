// Xato: defer xato tekshiruvidan OLDIN. http.Get xato qaytarsa resp nil bo'ladi, defer qatori esa resp.Body'ni
// darhol o'qiydi — dastur nil pointer bilan yiqiladi (panic va defer tafsilotlari: GO-1.15).
package main

import (
	"fmt"
	"net/http"
)

func main() {
	resp, err := http.Get("http://localhost:9/users/42") // 9-port: hech kim tinglamaydi
	defer resp.Body.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}
