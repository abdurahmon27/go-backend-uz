package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	port := os.Getenv("PORT") // sozlama muhitdan (SYS-0.02)
	if port == "" {
		port = "8077"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	db, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err == nil {
		err = db.Ping(ctx) // pgxpool.New ulanmaydi: ulanishni shu yerda tekshiramiz
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "server: baza:", err)
		os.Exit(1)
	}
	defer db.Close()
	if _, err := db.Exec(ctx, jadval); err != nil {
		fmt.Fprintln(os.Stderr, "server: jadval:", err)
		os.Exit(1)
	}
	fmt.Println("server: bazaga ulandim")

	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Fprintln(os.Stderr, "server:", err)
		os.Exit(1)
	}
	srv := &http.Server{Handler: daftar(db)}
	go srv.Serve(ln)
	fmt.Println("server: :" + port + " da tinglayapman")
	kut(srv)
}
