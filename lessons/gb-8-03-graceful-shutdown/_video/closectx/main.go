// closectx — research.md F4: Shutdown muddati tugagach srv.Close() faol handler'ning r.Context()'ini bekor qiladimi?
// go run ./_video/closectx  →  "shutdown: context deadline exceeded", "ctx: context canceled" (go1.27.4)
package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func main() {
	done := make(chan string, 1)
	srv := &http.Server{Addr: "127.0.0.1:8793", Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
			done <- "ctx: " + r.Context().Err().Error()
		case <-time.After(10 * time.Second):
			done <- "timer"
		}
	})}
	go srv.ListenAndServe()
	time.Sleep(200 * time.Millisecond)
	go http.Get("http://127.0.0.1:8793/")
	time.Sleep(200 * time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	fmt.Println("shutdown:", srv.Shutdown(ctx))
	srv.Close()
	select {
	case s := <-done:
		fmt.Println(s)
	case <-time.After(3 * time.Second):
		fmt.Println("handler ctx not canceled within 3s")
	}
}
