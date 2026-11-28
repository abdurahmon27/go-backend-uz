// Package token — bozor'ning JWT'lari (GB-6.05).
//
// qolda.go: HS256 imzoli JWT'ni kutubxonasiz, qoʻlda yasash — har bayti koʻrinsin deb (RFC 7519, RFC 7515).
// token.go: ishdagi yoʻl — golang-jwt/jwt/v5 bilan yasash va qatʼiy tekshirish (RFC 8725).
package token

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

// b64 — base64url, "="siz (RFC 7515 §2: Base64url Encoding, padding yoʻq).
// Alifbo: A–Z a–z 0–9 - _ ; oddiy base64'dagi + va / oʻrniga - va _, shuning uchun URL va sarlavhada xavfsiz.
func b64(b []byte) string {
	return base64.RawURLEncoding.EncodeToString(b)
}

// SignHS256 — JWT'ni qoʻlda yasaydi: base64url(header) . base64url(payload) . base64url(imzo).
// Imzo — nuqtagacha boʻlgan matnning HMAC-SHA256'i, kalit bilan (RFC 7518 §3.2).
// header va payload — tayyor JSON matnlar; ular shifrlanmaydi, faqat imzolanadi.
func SignHS256(header, payload string, key []byte) string {
	input := b64([]byte(header)) + "." + b64([]byte(payload))
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(input))
	return input + "." + b64(mac.Sum(nil))
}
