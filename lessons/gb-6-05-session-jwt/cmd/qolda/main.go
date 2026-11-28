// qolda — HS256 imzoli JWT'ni kutubxonasiz yasaydi va uch qismini alohida qatorda chiqaradi (GB-6.05).
// Token: header.payload.imzo — uchala qism nuqta bilan ulanadi. Kalit har ishga tushishda yangi.
//
//	go run ./cmd/qolda
package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/abdurahmon27/bozor/internal/token"
)

func main() {
	key := token.NewKey()
	exp := time.Now().Add(15 * time.Minute).Unix() // Unix vaqt: 1970-01-01 UTC'dan beri soniyalar (RFC 7519 §2: NumericDate)
	header := `{"alg":"HS256","typ":"JWT"}`
	payload := fmt.Sprintf(`{"sub":"1","iss":"bozor","aud":"bozor-api","exp":%d}`, exp)
	t := token.SignHS256(header, payload, key)

	parts := strings.Split(t, ".")
	fmt.Println("header: ", parts[0])
	fmt.Println("payload:", parts[1])
	fmt.Println("imzo:   ", parts[2])
}
