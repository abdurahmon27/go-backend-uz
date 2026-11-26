package main

import (
	"encoding/json"
	"fmt"
)

// Javob: AuthorID ga ham ",string" qo'shamiz.
type Post struct {
	ID       int64 `json:"id,string"`
	AuthorID int64 `json:"author_id,string"`
}

func main() {
	b, _ := json.Marshal(Post{ID: 1234567890123456789, AuthorID: 9007199254740993})
	fmt.Println(string(b))
}
