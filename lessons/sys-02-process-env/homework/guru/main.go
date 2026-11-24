package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(3 * time.Second)
		fmt.Fprintln(w, "javob: Salom!")
	})
	srv := &http.Server{Addr: ":8077"}
	go srv.ListenAndServe()
	fmt.Println("server: :8077 da tinglayapman")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	<-ctx.Done()
	stop() // signal'ni qaytarib beramiz: ikkinchi Ctrl+C / SIGTERM endi dasturni darhol toʻxtatadi
	fmt.Println("server: signal keldi, kutyapman (yana bir marta bossangiz, darhol chiqaman)")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
	fmt.Println("server: xayr")
}
