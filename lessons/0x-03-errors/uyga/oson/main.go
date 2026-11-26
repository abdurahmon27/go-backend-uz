// Uy vazifasi 🟢 (GO-X.03): bu config yuklovchi xatolarni uch joyda tashlab yuboradi.
//
//  1. Avval taxmin qiling, keyin dars papkasidan ishga tushiring:
//     go run ./uyga/oson                                  # buzuq config.json (ortiqcha vergul)
//     CONFIG=uyga/oson/config.json go run ./uyga/oson     # to'g'ri config
//  2. Uchala joyni ham if err != nil bilan tuzating:
//     loadConfig xatoni chaqiruvchiga qaytarsin, main esa log.Fatal bilan sababini yozib to'xtasin.
//  3. Tekshiring: buzuq config → "invalid character '}' …", CONFIG=yoq.json → "open yoq.json: no such file
//     or directory", to'g'ri config → "port: 8080".
//     Qiziq: golangci-lint run ./uyga/oson hozir faqat Unmarshal'ni ko'rsatadi — _ bilan tashlangan xatolarni
//     errcheck sukut bo'yicha sanamaydi (check-blank: false). Ko'z bilan topasiz.
//
// Javob: uyga/javob/oson/main.go
package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Port int `json:"port"`
}

func loadConfig(path string) (Config, error) {
	var cfg Config
	data, _ := os.ReadFile(path)
	json.Unmarshal(data, &cfg)
	return cfg, nil
}

func main() {
	path := os.Getenv("CONFIG")
	if path == "" {
		path = "config.json"
	}
	cfg, _ := loadConfig(path)
	fmt.Println("port:", cfg.Port)
}
