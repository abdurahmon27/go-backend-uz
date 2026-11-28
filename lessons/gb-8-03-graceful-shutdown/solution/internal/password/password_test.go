package password

import (
	"errors"
	"strings"
	"testing"
)

// Hash satri PHC shaklida va hozirgi parametrlar bilan yasalgan.
func TestHashFormat(t *testing.T) {
	h := Hash("tarvuz-osmon-qalam-ariq-27")
	if !strings.HasPrefix(h, "$argon2id$v=19$") {
		t.Fatalf("Hash %q: «$argon2id$v=19$» bilan boshlanishi kerak", h)
	}
	p, err := decode(h)
	if err != nil {
		t.Fatalf("decode(Hash(…)): %v", err)
	}
	if p.memory != memory || p.passes != passes || p.threads != threads {
		t.Errorf("satrdagi parametrlar m=%d,t=%d,p=%d, konstantalar m=%d,t=%d,p=%d",
			p.memory, p.passes, p.threads, memory, passes, threads)
	}
	if len(p.salt) != saltLen || len(p.key) != keyLen {
		t.Errorf("salt %d bayt, kalit %d bayt; kutilgan %d va %d", len(p.salt), len(p.key), saltLen, keyLen)
	}
}

// Bir xil parol ikki marta — ikki xil satr: har safar yangi salt.
func TestHashSalt(t *testing.T) {
	a, b := Hash("choy2024"), Hash("choy2024")
	if a == b {
		t.Fatalf("bir xil parol ikki marta bir xil satr berdi: salt tasodifiy emas\n%s", a)
	}
}

func TestVerify(t *testing.T) {
	h := Hash("tarvuz-osmon-qalam-ariq-27")
	if ok, err := Verify("tarvuz-osmon-qalam-ariq-27", h); err != nil || !ok {
		t.Errorf("toʻgʻri parol: Verify = %v, %v; kutilgan true, nil", ok, err)
	}
	if ok, err := Verify("tarvuz-osmon-qalam-ariq-28", h); err != nil || ok {
		t.Errorf("notoʻgʻri parol: Verify = %v, %v; kutilgan false, nil", ok, err)
	}
}

// Buzuq satr — ErrBadHash, panic emas va «har qanday parol mos» ham emas.
func TestVerifyBadHash(t *testing.T) {
	good := Hash("x")
	parts := strings.Split(good, "$")
	bad := map[string]string{
		"boʻsh":            "",
		"sha256 hex":       "385426486641a520e865c4e6b24d379c8f8afa1eb5de2f932ba90f0622f06289",
		"boshqa algoritm":  strings.Replace(good, "argon2id", "argon2i", 1),
		"t=0":              "$argon2id$v=19$m=19456,t=0,p=1$" + parts[4] + "$" + parts[5],
		"p=0":              "$argon2id$v=19$m=19456,t=2,p=0$" + parts[4] + "$" + parts[5],
		"boʻsh kalit":      "$argon2id$v=19$m=19456,t=2,p=1$" + parts[4] + "$",
		"padding bilan":    good + "=",
		"base64 emas salt": "$argon2id$v=19$m=19456,t=2,p=1$!!!$" + parts[5],
	}
	for name, h := range bad {
		ok, err := Verify("x", h)
		if ok || !errors.Is(err, ErrBadHash) {
			t.Errorf("%s: Verify = %v, %v; kutilgan false va ErrBadHash", name, ok, err)
		}
	}
}
