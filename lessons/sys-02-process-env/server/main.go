package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT") // sozlama muhitdan keladi
	if port == "" {
		port = "8077"
	}
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("server:", err) // port band: yolgʻon gapirmaymiz
		os.Exit(1)
	}
	srv := &http.Server{Handler: http.HandlerFunc(sekin)}
	go srv.Serve(ln)
	fmt.Println("server: :" + port + " da tinglayapman")
	if os.Getenv("KUT") == "0" {
		select {} // signal ushlanmaydi: Go odatdagidek darhol chiqadi
	}
	kut(srv)
}
