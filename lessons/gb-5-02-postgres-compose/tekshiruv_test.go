package bozordb

// Vazifani tekshiradi: `go test` oʻtguncha ishlang. Bu faylni oʻzgartirmang.
// compose.yaml'ni `docker compose config` oʻqiydi (JSON ko'rinishida), keyin ishlayotgan
// konteynerlarga psql bilan so'rov yuboradi. Shuning uchun oldin: docker compose up -d.

import (
	"encoding/json"
	"os/exec"
	"strings"
	"testing"
	"time"
)

type composeConfig struct {
	Services map[string]struct {
		Image       string            `json:"image"`
		Environment map[string]string `json:"environment"`
		Ports       []struct {
			HostIP    string `json:"host_ip"`
			Target    int    `json:"target"`
			Published string `json:"published"`
		} `json:"ports"`
		Volumes []struct {
			Type   string `json:"type"`
			Source string `json:"source"`
			Target string `json:"target"`
		} `json:"volumes"`
	} `json:"services"`
}

func docker(t *testing.T, args ...string) string {
	t.Helper()
	if _, err := exec.LookPath("docker"); err != nil {
		t.Fatalf("docker topilmadi: Docker oʻrnatilgan va PATH'da boʻlishi kerak")
	}
	out, err := exec.Command("docker", args...).CombinedOutput()
	if err != nil {
		t.Fatalf("docker %s xato berdi: %v\n%s", strings.Join(args, " "), err, out)
	}
	return string(out)
}

func readConfig(t *testing.T) composeConfig {
	t.Helper()
	out := docker(t, "compose", "-f", "compose.yaml", "config", "--format", "json")
	var cfg composeConfig
	if err := json.Unmarshal([]byte(out), &cfg); err != nil {
		t.Fatalf("docker compose config JSON'ini oʻqib boʻlmadi: %v", err)
	}
	return cfg
}

func TestAsosiyServis(t *testing.T) {
	db, ok := readConfig(t).Services["db"]
	if !ok {
		t.Fatalf("compose.yaml'da db servisi yoʻq: videodagi faylni qaytaring")
	}
	if !strings.HasPrefix(db.Image, "postgres:") {
		t.Errorf("db image'i postgres:… boʻlishi kerak, hozir %q", db.Image)
	}
	if len(db.Ports) != 1 || db.Ports[0].HostIP != "127.0.0.1" || db.Ports[0].Target != 5432 || db.Ports[0].Published != "5482" {
		t.Errorf(`db porti "127.0.0.1:5482:5432" boʻlishi kerak, hozir %+v`, db.Ports)
	}
	if len(db.Volumes) != 1 || db.Volumes[0].Type != "volume" || db.Volumes[0].Source != "pgdata" || db.Volumes[0].Target != "/var/lib/postgresql" {
		t.Errorf("db'da pgdata nomli volume /var/lib/postgresql'ga ulangan boʻlishi kerak, hozir %+v", db.Volumes)
	}
}

func TestIkkinchiServis(t *testing.T) {
	cfg := readConfig(t)
	test, ok := cfg.Services["db-test"]
	if !ok {
		t.Fatalf("compose.yaml'da db-test servisi yoʻq: db blokini nusxalab, nomini db-test qiling va TASK.md'dagi 4 narsani oʻzgartiring")
	}
	if !strings.HasPrefix(test.Image, "postgres:") {
		t.Errorf("db-test image'i postgres:… boʻlishi kerak, hozir %q", test.Image)
	}
	if got := test.Environment["POSTGRES_DB"]; got != "bozor_test" {
		t.Errorf("db-test'da POSTGRES_DB bozor_test boʻlishi kerak, hozir %q", got)
	}
	if len(test.Ports) != 1 || test.Ports[0].HostIP != "127.0.0.1" || test.Ports[0].Target != 5432 || test.Ports[0].Published != "5483" {
		t.Errorf(`db-test porti "127.0.0.1:5483:5432" boʻlishi kerak (chap raqam xostniki), hozir %+v`, test.Ports)
	}
	if len(test.Volumes) != 1 || test.Volumes[0].Type != "volume" || test.Volumes[0].Target != "/var/lib/postgresql" {
		t.Fatalf("db-test'da nomli volume /var/lib/postgresql'ga ulangan boʻlishi kerak, hozir %+v", test.Volumes)
	}
	if test.Volumes[0].Source == "pgdata" {
		t.Errorf("db-test oʻz volume'iga ega boʻlishi kerak (masalan pgtest); pgdata asosiy bazaniki")
	}
}

// psql xizmat ichida bitta so'rovni bajaradi. Postgres hali ishga tushayotgan bo'lsa (up'dan keyingi birinchi soniyalar),
// 30 soniyagacha kutib qayta urinadi; xizmat umuman ishlamayotgan bo'lsa, darhol yiqiladi.
func psql(t *testing.T, service, database, query string) string {
	t.Helper()
	var out []byte
	var err error
	for deadline := time.Now().Add(30 * time.Second); ; time.Sleep(time.Second) {
		out, err = exec.Command("docker", "compose", "-f", "compose.yaml", "exec", "-T", service,
			"psql", "-U", "bozor", "-d", database, "-tAc", query).CombinedOutput()
		if err == nil {
			return strings.TrimSpace(string(out))
		}
		if strings.Contains(string(out), "is not running") || time.Now().After(deadline) {
			break
		}
	}
	t.Fatalf("%s servisiga psql bilan ulanib boʻlmadi (%v): oldin docker compose up -d\n%s", service, err, out)
	return ""
}

func TestAsosiyBazaIshlayapti(t *testing.T) {
	if got := psql(t, "db", "bozor", "SELECT count(*) FROM products"); got != "3" {
		t.Errorf("products jadvalida 3 ta qator boʻlishi kerak, hozir %q: oldin docker compose up -d, "+
			"keyin docker compose exec -T db psql -U bozor -d bozor < bozor.sql", got)
	}
}

func TestTestBazaIshlayapti(t *testing.T) {
	if got := psql(t, "db-test", "bozor_test", "SELECT 1"); got != "1" {
		t.Errorf("db-test bazasi SELECT 1'ga 1 qaytarishi kerak, hozir %q: docker compose up -d", got)
	}
}
