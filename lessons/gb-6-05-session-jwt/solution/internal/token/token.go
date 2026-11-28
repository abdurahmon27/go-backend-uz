package token

import (
	"crypto/rand"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// iss va aud: kim bergan va kim uchun (RFC 7519 §4.1.1, §4.1.3). Parse ikkalasini ham talab qiladi:
// boshqa xizmat uchun berilgan token bozor API'siga yaramasin (RFC 8725 §3.9).
const (
	Issuer   = "bozor"
	Audience = "bozor-api"
)

// KeySize — HS256 kaliti kamida 256 bit boʻlsin (RFC 7518 §3.2: hash chiqishi bilan bir xil yoki uzunroq).
const KeySize = 32

// NewKey — 32 tasodifiy bayt. Kim shu kalitni bilsa, istalgan token yasay oladi: u secret (GB-8.01).
func NewKey() []byte {
	key := make([]byte, KeySize)
	rand.Read(key) // Go 1.24+: hech qachon xato qaytarmaydi (go doc crypto/rand.Read)
	return key
}

// Claims — tekshiruvdan oʻtgan tokendan bizga keraklari.
type Claims struct {
	UserID  int64     // sub
	JTI     string    // jti: shu tokenning noyob ID'si (RFC 7519 §4.1.7) — denylist kaliti
	Expires time.Time // exp
}

// Issue foydalanuvchiga ttl muddatli token beradi: sub, iss, aud, iat, exp va tasodifiy jti.
func Issue(userID int64, ttl time.Duration, key []byte) (string, error) {
	if len(key) < KeySize {
		return "", fmt.Errorf("token: kalit %d bayt, kamida %d kerak", len(key), KeySize)
	}
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(userID, 10),
		Issuer:    Issuer,
		Audience:  jwt.ClaimStrings{Audience},
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		ID:        rand.Text(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(key)
}

// Parse tokenni tekshiradi va claim'larini qaytaradi. Qoidalar (RFC 8725 §3.1, §3.8, §3.9):
//   - algoritm faqat HS256: header'dagi alg'ga ishonmaymiz ("none", boshqa alg — rad);
//   - imzo shu kalit bilan toʻgʻri;
//   - iss = bozor, aud = bozor-api;
//   - exp bor va oʻtmagan (kutubxona exp'siz tokenni oʻzi rad etmaydi: WithExpirationRequired).
func Parse(s string, key []byte) (Claims, error) {
	rc := &jwt.RegisteredClaims{}
	_, err := jwt.ParseWithClaims(s, rc, func(*jwt.Token) (any, error) {
		return key, nil
	},
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithIssuer(Issuer),
		jwt.WithAudience(Audience),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return Claims{}, err
	}
	id, err := strconv.ParseInt(rc.Subject, 10, 64)
	if err != nil {
		return Claims{}, fmt.Errorf("token: sub son emas: %w", err)
	}
	c := Claims{UserID: id, JTI: rc.ID}
	if rc.ExpiresAt != nil { // WithExpirationRequired bilan doim bor; nil'ga baribir tayanmaymiz
		c.Expires = rc.ExpiresAt.Time
	}
	return c, nil
}
