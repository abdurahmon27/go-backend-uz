// if err != nil: xatoni yuqoriga qaytarish (loadConfig) va ishga tushishda to'xtatish (main).
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
	err = json.Unmarshal(data, &cfg)
	return cfg, err
}

func main() {
	cfg, err := loadConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("port:", cfg.Port)
}
