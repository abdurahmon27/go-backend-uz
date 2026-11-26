package main

import (
	"encoding/json"
	"fmt"
)

// ",string" teg int64 ni JSON'da string qilib yozadi: JS uni yaxlitlamaydi.
type Post struct {
	ID int64 `json:"id,string"`
}

func main() {
	b, _ := json.Marshal(Post{ID: 1234567890123456789})
	fmt.Println(string(b))
}
