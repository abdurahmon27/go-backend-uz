// Javob 🟢 (GO-X.03): uchala xato tekshiriladi.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Config struct {
	Port int `json:"port"`
}

func loadConfig(path string) (Config, error) {
	var cfg Config
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func main() {
	path := os.Getenv("CONFIG")
	if path == "" {
		path = "config.json"
	}
	cfg, err := loadConfig(path)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("port:", cfg.Port)
}
