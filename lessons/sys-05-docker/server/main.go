package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT") // sozlama muhitdan (SYS-0.02)
	if port == "" {
		port = "8077"
	}
	dir := os.Getenv("DATA") // mehmonlar daftari shu papkada
	if dir == "" {
		dir = "data"
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		fmt.Println("server:", err)
		os.Exit(1)
	}
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("server:", err)
		os.Exit(1)
	}
	srv := &http.Server{Handler: daftar(dir)}
	go srv.Serve(ln)
	fmt.Println("server: :" + port + " da tinglayapman")
	kut(srv)
}
