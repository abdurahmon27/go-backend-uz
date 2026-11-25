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

// kut signal kelguncha kutadi, keyin serverni toza toʻxtatadi (SYS-0.02).
func kut(srv *http.Server) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer stop()
	<-ctx.Done()
	fmt.Println("server: signal keldi, toʻxtayapman")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Println("server: ulgurmadi:", err)
	}
	fmt.Println("server: xayr")
}
