package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Product — maydonlar kichik harf bilan (ataylab: tez-tez uchraydigan xato).
type Product struct {
	id    int64
	name  string
	price int64 // tiyin
	stock int
}

func main() {
	p := Product{id: 1, name: "olma", price: 1_200_000, stock: 40}
	data, err := json.Marshal(p)
	if err != nil {
		fmt.Fprintln(os.Stderr, "xato:", err)
		os.Exit(1)
	}
	fmt.Println(string(data))
}
