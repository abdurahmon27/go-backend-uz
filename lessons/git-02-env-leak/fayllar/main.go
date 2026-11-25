package main

import (
	"fmt"
	"os"
)

// Kalit koddan emas, muhitdan o'qiladi (SYS-0.02). Lekin .env fayli ham git'ga tushmasligi kerak.
func main() {
	kalit := os.Getenv("TOLOV_KALIT")
	if kalit == "" {
		fmt.Println("TOLOV_KALIT berilmagan")
		os.Exit(1)
	}
	fmt.Println("to'lov tizimiga ulandik, kalit uzunligi:", len(kalit))
}
