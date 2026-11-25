// net/http handler'dagi panic'ni ushlaydi, handler ochgan goroutine'dagisini esa yo'q.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	http.HandleFunc("/narx", func(w http.ResponseWriter, r *http.Request) {
		narxlar := []int{1200, 3400, 560}
		fmt.Fprintln(w, narxlar[5]) // panic: server uni ushlaydi
	})
	http.HandleFunc("/salom", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "server tirik")
	})
	http.HandleFunc("/fon", func(w http.ResponseWriter, r *http.Request) {
		go func() {
			time.Sleep(50 * time.Millisecond) // fondagi ish
			var hisob map[string]int
			hisob["fon"]++ // panic boshqa goroutine'da: hech kim ushlamaydi
		}()
		fmt.Fprintln(w, "fonda ishlayapti")
	})

	srv := &http.Server{Addr: "127.0.0.1:8015", ErrorLog: log.New(os.Stdout, "", 0)}
	go srv.ListenAndServe()
	time.Sleep(100 * time.Millisecond)

	for _, yol := range []string{"/narx", "/salom", "/fon", "/salom"} {
		so(yol)
		time.Sleep(100 * time.Millisecond)
	}
}

func so(yol string) {
	res, err := http.Get("http://127.0.0.1:8015" + yol)
	if err != nil {
		fmt.Println("javob", yol+":", "xato,", err)
		return
	}
	defer res.Body.Close()
	fmt.Println("javob", yol+":", res.Status)
}
