// O'rtadagi sinov: if ichida := bilan e'lon qilingan nom if'dan keyin yo'q. Ataylab build bo'lmaydi.
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
	if cfg, err := loadConfig("config.json"); err != nil {
		log.Fatal(err)
	}
	fmt.Println(cfg.Port)
}
