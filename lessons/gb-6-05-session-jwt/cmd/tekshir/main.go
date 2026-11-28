// tekshir — qoʻlda yasalgan JWT'ni va uning buzilgan nusxalarini token.Parse'dan oʻtkazadi (GB-6.05).
// Har qator: holat va Parse natijasi (ok yoki kutubxonaning xato matni).
//
//	go run ./cmd/tekshir
package main

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/abdurahmon27/bozor/internal/token"
)

// signHS512 — SignHS256'ning HS512 nusxasi: header'da alg HS512, imzo HMAC-SHA512.
func signHS512(payload string, key []byte) string {
	b64 := base64.RawURLEncoding.EncodeToString
	input := b64([]byte(`{"alg":"HS512","typ":"JWT"}`)) + "." + b64([]byte(payload))
	mac := hmac.New(sha512.New, key)
	mac.Write([]byte(input))
	return input + "." + b64(mac.Sum(nil))
}

func main() {
	key := token.NewKey()
	exp := time.Now().Add(15 * time.Minute).Unix()
	hs256 := `{"alg":"HS256","typ":"JWT"}`
	claims := func(sub, aud string, exp int64) string {
		if exp == 0 {
			return fmt.Sprintf(`{"sub":%q,"iss":"bozor","aud":%q}`, sub, aud)
		}
		return fmt.Sprintf(`{"sub":%q,"iss":"bozor","aud":%q,"exp":%d}`, sub, aud, exp)
	}
	b64 := base64.RawURLEncoding.EncodeToString

	good := token.SignHS256(hs256, claims("1", "bozor-api", exp), key)
	parts := strings.Split(good, ".")

	cases := []struct{ name, tok string }{
		{"qoʻlda yasalgan", good},
		// payload'da sub 1 → 2, imzo eski: header va imzo oʻsha, faqat oʻrtadagi qism yangi.
		{"sub 1 → 2, imzo eski", parts[0] + "." + b64([]byte(claims("2", "bozor-api", exp))) + "." + parts[2]},
		// alg "none": imzo boʻsh (RFC 7519 §6: unsecured JWT).
		{"alg none, imzo yoʻq", b64([]byte(`{"alg":"none","typ":"JWT"}`)) + "." + parts[1] + "."},
		{"exp yoʻq", token.SignHS256(hs256, claims("1", "bozor-api", 0), key)},
		{"aud boshqa xizmat", token.SignHS256(hs256, claims("1", "hisobot-api", exp), key)},
		{"muddati oʻtgan", token.SignHS256(hs256, claims("1", "bozor-api", time.Now().Add(-time.Minute).Unix()), key)},
		// alg HS512, imzo shu kalit bilan TOʻGʻRI hisoblangan: algoritmni token emas, server tanlaydi.
		{"alg HS512, imzo toʻgʻri", signHS512(claims("1", "bozor-api", exp), key)},
	}
	for _, c := range cases {
		cl, err := token.Parse(c.tok, key)
		if err != nil {
			fmt.Printf("%-24s %v\n", c.name+":", err)
			continue
		}
		fmt.Printf("%-24s ok, sub=%d\n", c.name+":", cl.UserID)
	}
}
