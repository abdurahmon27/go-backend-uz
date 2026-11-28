// Package password — parollarni argon2id bilan saqlash va tekshirish (GB-6.03).
//
// Bazaga parolning oʻzi ham, tez hash'i (sha256) ham yozilmaydi: faqat Hash qaytargan satr.
// Satrda hamma narsa bor: algoritm, versiya, parametrlar, salt va hosil boʻlgan kalit —
// shuning uchun Verify parametrlarni satrdan oʻqiydi va parametrlar keyin oʻzgarsa ham eski satrlar tekshiriladi.
package password

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// argon2id parametrlari — OWASP Password Storage Cheat Sheet'ning birinchi qatori (vazifa, 1-qadam):
// m=47104 KiB (46 MiB), t=1, p=1. Videodagi m=19456, t=2 bilan bir xil himoya, faqat xotira koʻproq.
// Eski satrlar (m=19456,t=2) Verify'dan oʻtaveradi: parametrlar satrdan oʻqiladi.
const (
	memory  = 46 * 1024 // m: xotira, KiB
	passes  = 1         // t: xotira ustidan necha marta oʻtiladi
	threads = 1         // p: parallel yoʻlaklar soni
	keyLen  = 32        // hosil boʻladigan kalit, bayt
	saltLen = 16        // salt, bayt
)

// b64 — PHC satr formatidagi base64: standart alifbo, oxirida «=» toʻldiruvchisiz.
var b64 = base64.RawStdEncoding

// Hash parolni argon2id bilan hash qiladi va bazaga yoziladigan satrni qaytaradi:
//
//	$argon2id$v=19$m=19456,t=2,p=1$<salt>$<kalit>
//
// Har chaqiruvda yangi tasodifiy salt: bir xil parol ham har safar boshqa satr beradi.
func Hash(pw string) string {
	salt := make([]byte, saltLen)
	rand.Read(salt) // Go 1.24'dan beri xato qaytarmaydi
	key := argon2.IDKey([]byte(pw), salt, passes, memory, threads, keyLen)
	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version,
		memory, passes, threads, b64.EncodeToString(salt), b64.EncodeToString(key))
}

// ErrBadHash — saqlangan satr argon2id formatida emas (buzilgan yoki boshqa algoritm).
var ErrBadHash = errors.New("hash satri notoʻgʻri")

// Verify parol saqlangan satrga mos kelishini tekshiradi. Parametrlar va salt satrdan olinadi,
// kalit qaytadan hisoblanadi va saqlangani bilan doimiy vaqtda solishtiriladi.
// Parol notoʻgʻri boʻlsa — false, nil; satrning oʻzi buzuq boʻlsa — xato (ErrBadHash).
func Verify(pw, encoded string) (bool, error) {
	p, err := decode(encoded)
	if err != nil {
		return false, err
	}
	key := argon2.IDKey([]byte(pw), p.salt, p.passes, p.memory, p.threads, uint32(len(p.key)))
	return subtle.ConstantTimeCompare(key, p.key) == 1, nil
}

// NeedsRehash — saqlangan satr hozirgi parametrlardan boshqasi bilan yasalganmi.
// Login paytida parol toʻgʻri chiqsa va bu true boʻlsa, parolni qaytadan Hash qilib saqlash kerak.
func NeedsRehash(encoded string) (bool, error) {
	p, err := decode(encoded)
	if err != nil {
		return false, err
	}
	return p.memory != memory || p.passes != passes || p.threads != threads || len(p.key) != keyLen, nil
}

// params — satrdan oʻqilgan qismlar.
type params struct {
	memory, passes uint32
	threads        uint8
	salt, key      []byte
}

// decode satrni «$» boʻyicha olti qismga ajratadi: "", "argon2id", "v=19", "m=…,t=…,p=…", salt, kalit.
func decode(encoded string) (params, error) {
	var p params
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 || parts[1] != "argon2id" || parts[2] != "v=19" {
		return p, ErrBadHash
	}
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &p.memory, &p.passes, &p.threads)
	if err != nil || p.passes < 1 || p.threads < 1 {
		return p, fmt.Errorf("%w: parametrlar %q", ErrBadHash, parts[3])
	}
	if p.salt, err = b64.DecodeString(parts[4]); err != nil {
		return p, fmt.Errorf("%w: salt: %w", ErrBadHash, err)
	}
	if p.key, err = b64.DecodeString(parts[5]); err != nil {
		return p, fmt.Errorf("%w: kalit: %w", ErrBadHash, err)
	}
	if len(p.key) < 16 { // boʻsh kalit har qanday parolga «mos» kelardi
		return p, fmt.Errorf("%w: kalit %d bayt", ErrBadHash, len(p.key))
	}
	return p, nil
}
