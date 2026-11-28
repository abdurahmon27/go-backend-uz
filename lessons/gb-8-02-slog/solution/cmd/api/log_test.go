// GB-8.02 videosidagi qoidalar: JSON log qatorlari, bitta soʻrovning hamma qatorida bitta request_id, 4xx qatorida
// mijoz yuborgan matn yoʻq, 5xx qatorida toʻliq xato, config'da parol yoʻq. Postgres kerak emas: go test ./cmd/api -run Slog
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// logTo — buf'ga JSON qator yozadigan logger (main'dagidek, faqat os.Stdout oʻrnida buf).
func logTo(buf *bytes.Buffer) *slog.Logger {
	return slog.New(slog.NewJSONHandler(buf, nil))
}

// logLines buf'dagi har qatorni JSON obyekt sifatida oʻqiydi.
func logLines(t *testing.T, buf *bytes.Buffer) []map[string]any {
	t.Helper()
	var out []map[string]any
	for _, line := range strings.Split(strings.TrimSpace(buf.String()), "\n") {
		var m map[string]any
		if err := json.Unmarshal([]byte(line), &m); err != nil {
			t.Fatalf("log qatori JSON emas: %q (%v)", line, err)
		}
		out = append(out, m)
	}
	return out
}

// leak — GB-8.02 videosidagi decode xatosi: mijozning parolidan besh harf (`\ulugb`) xato matnida.
var leak = fmt.Errorf("%w: %w", errBadJSON, errors.New("invalid escape sequence `\\ulugb` in string"))

func TestSlogRequestID(t *testing.T) {
	var buf bytes.Buffer
	h := requestID(logTo(&buf), logRequests(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeError(w, r, leak)
	})))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("POST", "/login", nil))
	id := rec.Header().Get("X-Request-Id")

	lines := logLines(t, &buf)
	if len(lines) != 2 {
		t.Fatalf("bitta soʻrov: %d ta log qatori, kutilgan 2 (avval xato, keyin soʻrov):\n%s", len(lines), &buf)
	}
	for i, want := range []struct{ level, msg string }{{"WARN", "mijoz xatosi"}, {"INFO", "soʻrov"}} {
		l := lines[i]
		if l["level"] != want.level || l["msg"] != want.msg || l["request_id"] != id {
			t.Errorf("%d-qator: level %v, msg %v, request_id %v; kutilgan %s, %q, %q", i+1, l["level"], l["msg"], l["request_id"], want.level, want.msg, id)
		}
	}
	if lines[0]["status"] != float64(400) || lines[1]["status"] != float64(400) {
		t.Errorf("status: %v va %v, kutilgan 400 va 400", lines[0]["status"], lines[1]["status"])
	}
	if strings.Contains(buf.String(), "ulugb") {
		t.Errorf("4xx qatorida mijozning matni (parol boʻlagi) bor:\n%s", &buf)
	}
}

func TestSlog5xxKeepsError(t *testing.T) {
	var buf bytes.Buffer
	h := requestID(logTo(&buf), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeError(w, r, fmt.Errorf("mahsulotlar oʻqilmadi: %w", errors.New("closed pool")))
	}))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/products", nil))
	l := logLines(t, &buf)[0]
	if rec.Code != 500 || l["level"] != "ERROR" || l["err"] != "mahsulotlar oʻqilmadi: closed pool" {
		t.Errorf("5xx: status %d, qator %v; kutilgan 500, ERROR va toʻliq err (tuzatish uchun kerak)", rec.Code, l)
	}
	if strings.Contains(rec.Body.String(), "closed pool") {
		t.Errorf("ichki xato mijozga chiqdi: %s", rec.Body)
	}
}

func TestSlogConfigValue(t *testing.T) {
	var buf bytes.Buffer
	logTo(&buf).Info("config", "config", config{port: 8782, databaseURL: labDB})
	if strings.Contains(buf.String(), "bozor-lab") {
		t.Fatalf("config qatorida parol bor: %s", &buf)
	}
	c, _ := logLines(t, &buf)[0]["config"].(map[string]any)
	if c["addr"] != ":8782" || c["db"] != "localhost:5582/bozor" {
		t.Errorf("config: %v; kutilgan {addr: \":8782\", db: \"localhost:5582/bozor\"} (LogValue)", c)
	}
}

func TestSlogFallback(t *testing.T) {
	if loggerFrom(context.Background()) != slog.Default() {
		t.Error("context'da logger yoʻq: loggerFrom slog.Default() qaytarishi kerak")
	}
}
