package main

import (
	"fmt"
	"net/http"
)

func main() {
	resp, err := http.Get("https://go.dev")
	if err != nil {
		fmt.Println("xato:", err)
		return
	}
	defer resp.Body.Close()
	fmt.Println(resp.Status)
}
