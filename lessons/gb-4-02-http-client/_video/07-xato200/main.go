package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type rateResponse struct {
	Result   string             `json:"result"`
	BaseCode string             `json:"base_code"`
	Rates    map[string]float64 `json:"rates"`
}

func fetchRate(client *http.Client, url, code string) (float64, error) {
	resp, err := client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("kurs soʻrovi: javob %s", resp.Status)
	}
	var data rateResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, fmt.Errorf("kurs javobi: %w", err)
	}
	return data.Rates[code], nil
}

func main() {
	client := &http.Client{Timeout: 5 * time.Second}
	rate, err := fetchRate(client, "https://open.er-api.com/v6/latest/XXX", "UZS")
	fmt.Println(rate, err)
}
