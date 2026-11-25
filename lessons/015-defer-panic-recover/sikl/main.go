// Ataylab xato: sikl ichidagi defer fayllarni funksiya tugaguncha yopmaydi.
package main

import (
	"fmt"
	"os"
)

func sana(nomlar []string) (int, error) {
	ochildi := 0
	for _, nom := range nomlar {
		f, err := os.Open(nom)
		if err != nil {
			return ochildi, err
		}
		defer f.Close() // XATO: sikl oxirida emas, sana tugaganda yopiladi
		ochildi++
	}
	return ochildi, nil
}

func main() {
	dir, nomlar := loglar(2000)
	defer os.RemoveAll(dir)

	n, err := sana(nomlar)
	fmt.Println("ochildi:", n, "ta fayl")
	if err != nil {
		fmt.Println("xato:", err)
	}
}

// loglar vaqtinchalik papkada n ta bo'sh fayl yaratadi va o'sha papkaga o'tadi.
func loglar(n int) (string, []string) {
	dir, err := os.MkdirTemp("", "loglar")
	if err != nil {
		panic(err)
	}
	if err := os.Chdir(dir); err != nil {
		panic(err)
	}
	nomlar := make([]string, n)
	for i := range nomlar {
		nomlar[i] = fmt.Sprintf("log-%04d.txt", i+1)
		if err := os.WriteFile(nomlar[i], nil, 0o644); err != nil {
			panic(err)
		}
	}
	return dir, nomlar
}
