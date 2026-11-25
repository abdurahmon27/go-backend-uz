package main

import "testing"

// Prefiksli soxta kalitlar (sk-ant-, ghp_, AKIA) boʻlaklardan yigʻiladi: repo'da provayder
// formatidagi toʻliq satr turmaydi. Qolgan qiymatlar (Qovun2026, Xq7L…) hech qaysi formatga mos emas.
func TestQatorTekshir(t *testing.T) {
	jadval := []struct {
		nom, qator, tur string // tur == "" — topilma boʻlmasligi kerak
	}{
		{"anthropic", "KEY=" + "sk-ant-" + "api03-SOXTA7kalit9Qx2mN4vB8c", "anthropic-kalit"},
		{"github classic", "token " + "gh" + "p_" + "SOXTA0123456789abcdefABCDEF", "github-token"},
		{"github fine-grained", "github" + "_pat_" + "SOXTA0123456789abcdef", "github-token"},
		{"aws", "id=" + "AK" + "IA" + "SOXTA0123456789A", "aws-kalit"},
		{"pem", "-----BEGIN RSA " + "PRIVATE KEY-----", "maxfiy-kalit"},
		{"dsn", "postgres://bozor:Qovun2026@localhost/bozor", "url-parol"},
		{"parol=", `password = "Qovun2026"`, "parol"},
		{"parol yaml", "parol: Qovun2026", "parol"},
		{"DB_PASSWORD", "DB_PASSWORD=Qovun2026", "parol"},
		{"DB_PASS", "DB_PASS=Qovun2026", "parol"},
		{"email", "mijoz=aziz@example.com", "email"},
		{"telefon", "tel=+998 90 000 11 22", "telefon"},
		{"telefon chiziqcha", "tel=+998-90-000-11-22", "telefon"},
		{"telefon qavs", "tel=+998 (90) 000-11-22", "telefon"},
		{"karta luhn", "karta=4111 1111 1111 1111", "karta"},
		{"nomaʼlum kalit", "kalit=Xq7Lm2Pz9Rt4Vb8Nc3Kd6Wf1Hs5Jy0Ge", "nomaʼlum-kalit"},

		{"buyurtma raqami (luhn emas)", "id=2026112510140712", ""},
		{"nollar (luhn oʻtadi, lekin karta emas)", "0000 0000 0000 0000", ""},
		{"git commit hash", "commit 3f9a1c27e0b84d6a9e2c5b7f1d08a4c6e9b2f731", ""},
		{"uuid", "id=550e8400-e29b-41d4-a716-446655440000", ""},
		{"go nomi", "func http2ErrCodeEnhanceYourCalm() {}", ""},
		{"boʻsh qiymat", "ANTHROPIC_API_KEY=", ""},
		{"oddiy gap", "panic: runtime error: invalid memory address", ""},
	}
	for _, j := range jadval {
		t.Run(j.nom, func(t *testing.T) {
			got := qatorTekshir("test:1", j.qator)
			switch {
			case j.tur == "" && len(got) != 0:
				t.Fatalf("topilma boʻlmasligi kerak edi, keldi: %s %q", got[0].tur, got[0].qiymat)
			case j.tur != "" && (len(got) != 1 || got[0].tur != j.tur):
				t.Fatalf("kutilgan bitta %s, keldi: %+v", j.tur, got)
			}
		})
	}
}

func TestLuhn(t *testing.T) {
	for s, want := range map[string]bool{
		"4111 1111 1111 1111": true,  // mashhur test karta raqami
		"4111 1111 1111 1112": false, // oxirgi raqam oʻzgardi
		"2026112510140712":    false,
	} {
		if luhn(s) != want {
			t.Errorf("luhn(%q) = %v, kutilgan %v", s, !want, want)
		}
	}
}

func TestYashir(t *testing.T) {
	for in, want := range map[string]string{
		"Qovun2026":                     "*********",
		"sk-ant-" + "api03-SOXTA7kalit": "sk-a…it",
	} {
		if got := yashir(in); got != want {
			t.Errorf("yashir(%q) = %q, kutilgan %q", in, got, want)
		}
	}
}

func TestEnvFayl(t *testing.T) {
	for nom, want := range map[string]bool{".env": true, ".env.local": true, ".env.example": false, "env.go": false} {
		if envFayl(nom) != want {
			t.Errorf("envFayl(%q) != %v", nom, want)
		}
	}
}
