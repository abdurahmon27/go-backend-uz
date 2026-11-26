package main

import (
	"encoding/json"
	"fmt"
)

// Vazifa: AuthorID ham frontend'da (JavaScript) buzilmasin.
// Tekshirish: go run ./uyga/oson | node uyga/check.js
type Post struct {
	ID       int64 `json:"id,string"`
	AuthorID int64 `json:"author_id"`
}

func main() {
	b, _ := json.Marshal(Post{ID: 1234567890123456789, AuthorID: 9007199254740993})
	fmt.Println(string(b))
}
