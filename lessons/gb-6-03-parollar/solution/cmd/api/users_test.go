package main

// POST /users testlari (GB-6.03): 201 javobida parol ham, hash ham yoʻq; 409, 422.

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestRegister(t *testing.T) {
	captureLog(t)
	h := newHandler(testStore(t))

	rec := send(h, "POST", "/users", `{"email":"dilnoza@bozor.test","password":"tarvuz-osmon-qalam-ariq-27"}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("POST /users: status %d, kutilgan 201. Body: %s", rec.Code, rec.Body)
	}
	body := rec.Body.String()
	if strings.Contains(body, "tarvuz") || strings.Contains(body, "argon2id") || strings.Contains(body, "password") {
		t.Errorf("201 javobida parol yoki hash bor: %s", body)
	}
	var u struct {
		ID    int64  `json:"id"`
		Email string `json:"email"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &u); err != nil || u.ID != 1 || u.Email != "dilnoza@bozor.test" {
		t.Errorf("javob %s (%v); kutilgan id 1 va email", body, err)
	}

	// Xuddi shu email, boshqa harf bilan — 409.
	rec = send(h, "POST", "/users", `{"email":"Dilnoza@bozor.test","password":"boshqa-uzun-parol-2025"}`)
	if rec.Code != http.StatusConflict || rec.Header().Get("Content-Type") != "application/problem+json" {
		t.Errorf("takroriy email: %d %s, kutilgan 409 problem", rec.Code, rec.Body)
	}
}

func TestRegisterValidation(t *testing.T) {
	captureLog(t)
	h := newHandler(testStore(t))
	for _, c := range []struct{ body, field string }{
		{`{"email":"dilnoza@bozor.test","password":"choy2024"}`, "password"},                         // 8 belgi
		{`{"email":"dilnoza@bozor.test","password":"` + strings.Repeat("oʻ", 33) + `"}`, "password"}, // 66 belgi
		{`{"email":"dilnoza","password":"tarvuz-osmon-qalam-ariq-27"}`, "email"},
	} {
		rec := send(h, "POST", "/users", c.body)
		var p problem
		json.Unmarshal(rec.Body.Bytes(), &p)
		if rec.Code != http.StatusUnprocessableEntity || p.Errors[c.field] == "" {
			t.Errorf("%s: %d %s, kutilgan 422 va errors.%s", c.body, rec.Code, rec.Body, c.field)
		}
	}
	// 15 belgi, 30 bayt: belgilar baytda emas, rune'da sanaladi (GB-2.10).
	rec := send(h, "POST", "/users", `{"email":"jasur@bozor.test","password":"`+strings.Repeat("ё", 15)+`"}`)
	if rec.Code != http.StatusCreated {
		t.Errorf("15 ta «ё» (15 belgi, 30 bayt): %d %s, kutilgan 201", rec.Code, rec.Body)
	}
}
