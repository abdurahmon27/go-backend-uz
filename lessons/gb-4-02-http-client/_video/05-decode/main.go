package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type rateResponse struct {
	Result   string             `json:"result"`
	BaseCode string             `json:"base_code"`
	Rates    map[string]float64 `json:"rates"`
}

func main() {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get("https://open.er-api.com/v6/latest/USD")
	if err != nil {
		fmt.Fprintln(os.Stderr, "xato:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	var data rateResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		fmt.Fprintln(os.Stderr, "xato:", err)
		os.Exit(1)
	}
	fmt.Println(data.BaseCode, data.Rates["UZS"])
}
