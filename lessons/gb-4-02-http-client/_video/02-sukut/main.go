package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	fmt.Println("soʻrov yuborildi")
	resp, err := http.Get("http://localhost:8089")
	if err != nil {
		fmt.Fprintln(os.Stderr, "xato:", err)
		os.Exit(1)
	}
	fmt.Println(resp.Status)
}
