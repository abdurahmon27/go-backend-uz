// Vazifa tekshiruvi (TASK.md, 2-qadam): User log'da faqat id va role bilan chiqadi — email (shaxsiy maʼlumot) va
// parol hash'i (secret) yoʻq. Postgres kerak emas: go test ./internal/store -run LogValue
// Bu faylni oʻzgartirmang.
package store

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
)

func TestUserLogValue(t *testing.T) {
	u := User{ID: 3, Email: "admin@bozor.test", Role: RoleAdmin,
		PasswordHash: "$argon2id$v=19$m=47104,t=1,p=1$XUKfu6Xd7Fcse0J9GFRDuA$PO1/aBbjTlLgep39ILqaMPF8thT1IGnqNdRedkSTFjs"}
	var buf bytes.Buffer
	slog.New(slog.NewJSONHandler(&buf, nil)).Info("login", "user", u)
	out := strings.TrimSpace(buf.String())

	if _, ok := any(u).(slog.LogValuer); !ok {
		t.Errorf("User'da LogValue metodi yoʻq: JSON handler uni json.Marshal qiladi va email chiqadi (TASK.md, 2-qadam).\n"+
			"hozir: %s", out)
	}
	if strings.Contains(out, "admin@bozor.test") || strings.Contains(out, "argon2id") {
		t.Fatalf("log qatorida email yoki parol hash'i bor:\n%s", out)
	}
	var line struct {
		User map[string]any `json:"user"`
	}
	if err := json.Unmarshal([]byte(out), &line); err != nil {
		t.Fatalf("log qatori JSON emas: %v", err)
	}
	if len(line.User) != 2 || line.User["id"] != float64(3) || line.User["role"] != "admin" {
		t.Errorf("user = %v; kutilgan faqat ikki kalit: {\"id\":3,\"role\":\"admin\"} (slog.GroupValue)", line.User)
	}
}
