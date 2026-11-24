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

// kut signal kelguncha kutadi, keyin serverni toza toʻxtatadi.
func kut(srv *http.Server) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer stop()
	<-ctx.Done()
	fmt.Println("server: signal keldi, boshlangan soʻrovlarni tugatyapman")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Println("server: ulgurmadi:", err) // 5 soniya tugadi: yolgʻon gapirmaymiz
	}
	fmt.Println("server: xayr")
}
