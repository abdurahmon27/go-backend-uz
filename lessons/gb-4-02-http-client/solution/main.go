package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const tiyinPerSom = 100 // 1 soʻm = 100 tiyin (GB-1.03)

// rateURL — ochiq kurs API'si: 1 USD necha valyutaga teng (kalitsiz, bepul).
const rateURL = "https://open.er-api.com/v6/latest/USD"

// Product — katalogdagi mahsulot (GB-2.03); narx tiyinda.
type Product struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Price int64  `json:"price"` // tiyin
	Stock int    `json:"stock"`
}

// rateResponse — kurs API'si javobining bizga kerak qismi.
type rateResponse struct {
	Result    string             `json:"result"`
	ErrorType string             `json:"error-type"`
	BaseCode  string             `json:"base_code"`
	Rates     map[string]float64 `json:"rates"`
}

// StatusError — server 200 dan boshqa status bilan javob bergani.
// fetchRate shu turdagi xato qaytaradi: chaqiruvchi errors.As bilan kodini oladi.
type StatusError struct {
	Code int
}

func (e *StatusError) Error() string {
	return fmt.Sprintf("server %d status bilan javob berdi", e.Code)
}

// fetchRate url'dan kurs javobini oladi va code valyutasining kursini qaytaradi
// (1 USD = ? code).
func fetchRate(client *http.Client, url, code string) (float64, error) {
	resp, err := client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, &StatusError{Code: resp.StatusCode}
	}
	var data rateResponse
	// Tashqi server javobining hajmiga ishonmaymiz: ko'pi bilan 1 MB (1<<20 bayt).
	body := io.LimitReader(resp.Body, 1<<20)
	if err := json.NewDecoder(body).Decode(&data); err != nil {
		return 0, fmt.Errorf("kurs javobi: %w", err)
	}
	if data.Result != "success" {
		return 0, fmt.Errorf("kurs API xatosi: %s", data.ErrorType)
	}
	rate, ok := data.Rates[code]
	if !ok {
		return 0, fmt.Errorf("kurs topilmadi: %s", code)
	}
	return rate, nil
}

func run() error {
	client := &http.Client{Timeout: 5 * time.Second}
	rate, err := fetchRate(client, rateURL, "UZS")
	if err != nil {
		return err
	}
	p := Product{ID: 1, Name: "olma", Price: 1_200_000, Stock: 40}
	som := float64(p.Price) / tiyinPerSom
	fmt.Printf("1 USD = %.2f soʻm\n", rate)
	fmt.Printf("%s: %.0f soʻm = %.2f USD\n", p.Name, som, som/rate)
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "xato:", err)
		os.Exit(1)
	}
}
