package main

// Sozlama (GB-8.01). Kod hamma muhitda bir xil; muhitga qarab oʻzgaradigan narsa (port, baza manzili, kalitlar)
// tashqaridan keladi: muhit oʻzgaruvchilari va flag'lardan (Twelve-Factor App, III. Config).
// Ustunlik: flag > muhit oʻzgaruvchisi > standart qiymat. Hammasi ishga tushishda BIR MARTA tekshiriladi:
// biror qiymat notoʻgʻri boʻlsa server hech narsa ochmaydi, hamma xatolarni sanab, 1 kodi bilan chiqadi.
//
//	PORT          -port   HTTP port, 1–65535 (standart 8765)
//	DATABASE_URL  —       Postgres manzili. SECRET: ichida parol bor. Faqat muhitdan, flag emas:
//	                      buyruq qatori (ps, /proc/PID/cmdline) shu kompyuterdagi hamma foydalanuvchiga koʻrinadi.
//
// Qiymatlar: .env.example (gitda) → nusxasi .env (gitda emas, .gitignore).

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strconv"
)

// config — serverning ishga tushish sozlamalari.
type config struct {
	port        int
	databaseURL string // secret: ichida parol bor
}

// envOr — key qoʻyilgan boʻlsa uning qiymati (boʻsh boʻlsa ham), qoʻyilmagan boʻlsa def.
// «PORT=» (boʻsh) — xato sozlama: u standart qiymatga jim almashmaydi, tekshiruvda ushlanadi.
func envOr(key, def string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return def
}

// loadConfig — args (os.Args[1:]) va muhitdan config. Hamma xato bitta error'da (errors.Join), har biri alohida qatorda.
// flag.ContinueOnError: notoʻgʻri flag'da Parse dasturni toʻxtatmaydi, xato qaytaradi (testlar uchun ham).
func loadConfig(args []string) (config, error) {
	fs := flag.NewFlagSet("api", flag.ContinueOnError)
	port := fs.String("port", envOr("PORT", "8765"), "HTTP port (muhitda: PORT)")
	if err := fs.Parse(args); err != nil {
		return config{}, err
	}

	var errs []error
	n, err := strconv.Atoi(*port)
	if err != nil || n < 1 || n > 65535 {
		errs = append(errs, fmt.Errorf("port %q: 1–65535 oraligʻidagi butun son kerak", *port))
	}
	db := os.Getenv("DATABASE_URL")
	if db == "" {
		errs = append(errs, errors.New("DATABASE_URL qoʻyilmagan: .env.example'ga qarang"))
	}
	if len(errs) > 0 {
		return config{}, errors.Join(errs...)
	}
	return config{port: n, databaseURL: db}, nil
}

// addr — http.Server'ning Addr maydoni: ":8765" — shu portda, kompyuterning hamma tarmoq manzillarida.
func (c config) addr() string {
	return fmt.Sprintf(":%d", c.port)
}

// String — config'ning log'dagi koʻrinishi. fmt'ning %v, %+v va %s'i shu metodni chaqiradi (GB-2.07).
// Bazadan faqat host va baza nomi: foydalanuvchi, parol va ?parametrlar chiqmaydi (OWASP Logging Cheat Sheet:
// «Database connection strings» log'ga toʻgʻridan-toʻgʻri yozilmaydi). %#v String'ni chaqirmaydi: config'ni %#v bilan yozmang.
func (c config) String() string {
	db := "[yashirin]"
	if u, err := url.Parse(c.databaseURL); err == nil && u.Host != "" {
		db = u.Host + u.Path
	}
	return fmt.Sprintf("port=%d db=%s", c.port, db)
}
