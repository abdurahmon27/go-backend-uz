package main

import (
	"encoding/json"
	"fmt"
)

type Item struct {
	Price int `json:"price"`
}

func main() {
	body := `{"price": "2"}`
	var item Item
	err := json.Unmarshal([]byte(body), &item)
	fmt.Println(item.Price, err)
}
