// Tuzatish: har bir fayl alohida funksiyada ochiladi, defer o'sha funksiya tugaganda ishlaydi.
package main

import (
	"fmt"
	"os"
)

func bittasi(nom string) error {
	f, err := os.Open(nom)
	if err != nil {
		return err
	}
	defer f.Close() // bittasi tugashi bilan yopiladi
	return nil
}

func sana(nomlar []string) (int, error) {
	ochildi := 0
	for _, nom := range nomlar {
		if err := bittasi(nom); err != nil {
			return ochildi, err
		}
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
