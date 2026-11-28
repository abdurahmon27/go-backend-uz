// GB-8.01 sozlamasi: PORT, DATABASE_URL, HOST / -host va parol String()da yoʻq. GB-8.01 vazifasi (HOST) bu yerda bajarilgan.
// Bu testlar Postgres'siz ishlaydi: go test ./cmd/api -run Config
package main

import (
	"os"
	"strings"
	"testing"
)

const labDB = "postgres://bozor:bozor-lab@localhost:5583/bozor"

// env — testdagi muhit: berilgan nomlar qoʻyiladi, "-" qiymatli nom muhitdan olib tashlanadi.
// t.Setenv test tugagach eski qiymatni qaytaradi.
func env(t *testing.T, kv map[string]string) {
	t.Helper()
	for k, v := range kv {
		t.Setenv(k, "")
		if v == "-" {
			os.Unsetenv(k)
		} else {
			os.Setenv(k, v)
		}
	}
}

func TestConfigVideo(t *testing.T) {
	env(t, map[string]string{"PORT": "-", "HOST": "-", "DATABASE_URL": labDB})
	cfg, err := loadConfig(nil)
	if err != nil {
		t.Fatalf("toʻgʻri sozlamada xato: %v", err)
	}
	if cfg.port != 8765 || cfg.addr() != ":8765" {
		t.Errorf("standart: port %d, addr() %q; kutilgan 8765 va \":8765\" (videodagi qoida: PORT yoʻq → 8765)", cfg.port, cfg.addr())
	}

	env(t, map[string]string{"PORT": "9090"})
	if cfg, err := loadConfig([]string{"-port", "8781"}); err != nil || cfg.port != 8781 {
		t.Errorf("PORT=9090 va -port 8781: port %d, xato %v; kutilgan 8781 (flag muhitdan ustun)", cfg.port, err)
	}

	env(t, map[string]string{"PORT": "80a", "DATABASE_URL": "-"})
	_, err = loadConfig(nil)
	if err == nil || !strings.Contains(err.Error(), `"80a"`) || !strings.Contains(err.Error(), "DATABASE_URL") {
		t.Errorf("PORT=80a va DATABASE_URL yoʻq: xato %v; kutilgan ikkala muammo bitta xatoda (errors.Join)", err)
	}
}

func TestConfigStringHidesPassword(t *testing.T) {
	cfg := config{port: 8765, databaseURL: labDB}
	s := cfg.String()
	if strings.Contains(s, "bozor-lab") {
		t.Errorf("String() = %q: parol log'ga chiqyapti", s)
	}
	if !strings.Contains(s, "localhost:5583/bozor") {
		t.Errorf("String() = %q: bazaning host'i va nomi koʻrinishi kerak (localhost:5583/bozor)", s)
	}
	if s := (config{databaseURL: "host=localhost password=bozor-lab"}).String(); strings.Contains(s, "bozor-lab") {
		t.Errorf("String() = %q: URL boʻlmagan DATABASE_URL ham yashirilishi kerak", s)
	}
}

func TestConfigHost(t *testing.T) {
	env(t, map[string]string{"PORT": "-", "HOST": "-", "DATABASE_URL": labDB})
	for _, c := range []struct {
		host string // HOST qiymati ("-" — qoʻyilmagan)
		args []string
		want string // cfg.addr()
	}{
		{"-", nil, ":8765"},
		{"127.0.0.1", nil, "127.0.0.1:8765"},
		{"localhost", nil, "localhost:8765"},
		{"::1", nil, "[::1]:8765"},
		{"0.0.0.0", []string{"-host", "127.0.0.1"}, "127.0.0.1:8765"},
	} {
		env(t, map[string]string{"HOST": c.host})
		cfg, err := loadConfig(c.args)
		if err != nil {
			t.Errorf("HOST=%s %v: xato %v\n  kutilgan addr() %q (TASK.md, 1-qadam: -host flag'i, standarti envOr(\"HOST\", \"\"))",
				c.host, c.args, err, c.want)
			continue
		}
		if got := cfg.addr(); got != c.want {
			t.Errorf("HOST=%s %v: addr() %q, kutilgan %q (TASK.md, 3-qadam: addr() host'ni ham qoʻshadi, net.JoinHostPort)",
				c.host, c.args, got, c.want)
		}
	}

	for _, bad := range []string{"127.0.0.300", "bozor.uz", " 127.0.0.1"} {
		env(t, map[string]string{"HOST": bad, "DATABASE_URL": "-"})
		_, err := loadConfig(nil)
		if err == nil || !strings.Contains(err.Error(), "host") || !strings.Contains(err.Error(), "DATABASE_URL") {
			t.Errorf("HOST=%q va DATABASE_URL yoʻq: xato %v\n  kutilgan ikkala muammo bitta xatoda: host %q … va DATABASE_URL … (TASK.md, 2-qadam)",
				bad, err, bad)
		}
	}
}
