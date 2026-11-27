package main

import (
	"encoding/json"
	jsonv2 "encoding/json/v2"
	"fmt"
)

// Product — katalogdagi bitta mahsulot.
type Product struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Price int64  `json:"price"` // tiyin
	Stock int    `json:"stock"`
}

func main() {
	in := []byte(`{"id":1,"name":"olma","name":"non"}`)
	var p Product
	err := json.Unmarshal(in, &p)
	fmt.Println("v1:", err, p.Name)
	err = jsonv2.Unmarshal(in, &p)
	fmt.Println("v2:", err)
}
