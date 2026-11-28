package token

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

// payload — tekshiruv uchun claim'lar matni (exp 0 boʻlsa exp yoʻq).
func payload(sub, iss, aud string, exp int64) string {
	if exp == 0 {
		return fmt.Sprintf(`{"sub":%q,"iss":%q,"aud":%q}`, sub, iss, aud)
	}
	return fmt.Sprintf(`{"sub":%q,"iss":%q,"aud":%q,"exp":%d}`, sub, iss, aud, exp)
}

const hs256 = `{"alg":"HS256","typ":"JWT"}`

func TestQoldaYasalganTokenOtadi(t *testing.T) {
	key := NewKey()
	exp := time.Now().Add(time.Minute).Unix()
	tok := SignHS256(hs256, payload("7", Issuer, Audience, exp), key)
	if n := strings.Count(tok, "."); n != 2 {
		t.Fatalf("tokenda %d ta nuqta, kutilgan 2 (header.payload.imzo)", n)
	}
	c, err := Parse(tok, key)
	if err != nil || c.UserID != 7 {
		t.Fatalf("qoʻlda yasalgan token: UserID=%d err=%v, kutilgan 7 va nil", c.UserID, err)
	}
}

func TestParseRadEtadi(t *testing.T) {
	key := NewKey()
	exp := time.Now().Add(time.Minute).Unix()
	good := SignHS256(hs256, payload("1", Issuer, Audience, exp), key)
	parts := strings.Split(good, ".")
	cases := map[string]string{
		"boshqa kalit":   SignHS256(hs256, payload("1", Issuer, Audience, exp), NewKey()),
		"sub almashgan":  parts[0] + "." + b64([]byte(payload("2", Issuer, Audience, exp))) + "." + parts[2],
		"alg none":       b64([]byte(`{"alg":"none","typ":"JWT"}`)) + "." + parts[1] + ".",
		"alg HS384":      SignHS256(`{"alg":"HS384","typ":"JWT"}`, payload("1", Issuer, Audience, exp), key),
		"exp yoʻq":       SignHS256(hs256, payload("1", Issuer, Audience, 0), key),
		"muddati oʻtgan": SignHS256(hs256, payload("1", Issuer, Audience, time.Now().Add(-time.Minute).Unix()), key),
		"iss boshqa":     SignHS256(hs256, payload("1", "begona", Audience, exp), key),
		"aud boshqa":     SignHS256(hs256, payload("1", Issuer, "hisobot-api", exp), key),
		"sub son emas":   SignHS256(hs256, payload("dilnoza", Issuer, Audience, exp), key),
		"ikki qism":      parts[0] + "." + parts[1],
	}
	for name, tok := range cases {
		if _, err := Parse(tok, key); err == nil {
			t.Errorf("%s: Parse qabul qildi, rad etishi kerak edi", name)
		}
	}
}

func TestIssue(t *testing.T) {
	key := NewKey()
	tok, err := Issue(42, 15*time.Minute, key)
	if err != nil {
		t.Fatal(err)
	}
	c, err := Parse(tok, key)
	if err != nil {
		t.Fatalf("Issue bergan token Parse'dan oʻtmadi: %v", err)
	}
	if c.UserID != 42 || c.JTI == "" {
		t.Errorf("claims: %+v; kutilgan UserID 42 va boʻsh boʻlmagan JTI", c)
	}
	if d := time.Until(c.Expires); d <= 14*time.Minute || d > 15*time.Minute {
		t.Errorf("exp %v keyin, kutilgan ~15 daqiqa", d)
	}
	other, _ := Issue(42, 15*time.Minute, key)
	if c2, _ := Parse(other, key); c2.JTI == c.JTI {
		t.Error("ikki tokenning jti'si bir xil: har token oʻz jti'sini olishi kerak")
	}
	if _, err := Issue(1, time.Minute, []byte("qisqa")); err == nil {
		t.Error("5 baytli kalit qabul qilindi: HS256 kaliti kamida 32 bayt boʻlsin")
	}
}
