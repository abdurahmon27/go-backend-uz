package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get("https://open.er-api.com/v6/latest/USD")
	defer resp.Body.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, "xato:", err)
		os.Exit(1)
	}
	fmt.Println(resp.Status)
}
