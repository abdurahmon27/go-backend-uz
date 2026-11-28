// Vazifa tekshiruvi (TASK.md, 1-qadam): muvaffaqiyatli login log'da bitta qator — "login", ichida foydalanuvchining
// id va role'i, email yoʻq. Haqiqiy Postgres kerak: docker compose up -d --wait, keyin go test ./cmd/api -run LoginLog
// Bu faylni oʻzgartirmang.
package main

import (
	"bytes"
	"strings"
	"sync"
	"testing"

	"github.com/abdurahmon27/bozor/internal/pgtest"
	"github.com/abdurahmon27/bozor/internal/store"
)

func TestLoginLogLine(t *testing.T) {
	var buf bytes.Buffer
	h := newHandler(store.New(pgtest.New(t)), logTo(&buf), new(sync.WaitGroup))
	rec := send(h, "POST", "/login", adminLogin, "")
	if rec.Code != 200 {
		t.Fatalf("admin login: status %d, kutilgan 200: %s", rec.Code, rec.Body)
	}
	id := rec.Header().Get("X-Request-Id")

	var line map[string]any
	for _, l := range logLines(t, &buf) {
		if l["msg"] == "login" {
			line = l
		}
	}
	if line == nil {
		t.Fatalf("log'da \"login\" qatori yoʻq (TASK.md, 1-qadam: login handler'ida, cookie yozilgandan keyin "+
			"soʻrovning logger'i bilan Info darajada: xabar \"login\", atribut \"user\"). Hozirgi log:\n%s", &buf)
	}
	if line["level"] != "INFO" || line["request_id"] != id {
		t.Errorf("login qatori: level %v, request_id %v; kutilgan INFO va %q (soʻrovning logger'i)", line["level"], line["request_id"], id)
	}
	if strings.Contains(buf.String(), "admin@bozor.test") {
		t.Errorf("log'da email bor (TASK.md, 2-qadam: store.User'ga LogValue — faqat id va role):\n%s", &buf)
	}
	u, _ := line["user"].(map[string]any)
	if u["id"] != float64(3) || u["role"] != "admin" {
		t.Errorf("login qatorida user = %v; kutilgan {\"id\":3,\"role\":\"admin\"}", line["user"])
	}
}
