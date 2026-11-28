package main

// Login'ga rate limit (GB-6.07). Har login urinishi argon2id hisoblaydi (GB-6.03: 46 MiB xotira):
// cheklovsiz POST /login — parol taxmin qilish uchun ochiq eshik va serverni band qilishning arzon yoʻli.
//
// golang.org/x/time/rate — Go jamoasining token bucket'i: chelakda eng koʻpi `burst` ta token,
// har `every` vaqtda bittadan toʻladi; har soʻrov bitta token oladi; token yoʻq — 429 Too Many Requests
// (RFC 6585 §4) va Retry-After: necha soniyadan keyin qayta urinish mumkin (RFC 9110 §10.2.3).

import (
	"math"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// limiters — har kalit (IP yoki email) uchun alohida chelak. Map bir vaqtda koʻp soʻrovdan
// ishlatiladi, shuning uchun mutex (GB-4.07). Kalitlar oʻchirilmaydi: uzoq ishlaydigan serverda
// map oʻsib boradi — tozalash: sweep (🟡 ixtiyoriy vazifa).
type limiters struct {
	mu    sync.Mutex
	m     map[string]*rate.Limiter
	every rate.Limit
	burst int
}

// newLimiters — har kalitga: boshida burst ta urinish, keyin har every'da bittadan.
func newLimiters(every time.Duration, burst int) *limiters {
	return &limiters{m: map[string]*rate.Limiter{}, every: rate.Every(every), burst: burst}
}

// get — kalitning chelagi; birinchi murojaatda yangisi (toʻla) yaratiladi.
func (l *limiters) get(key string) *rate.Limiter {
	l.mu.Lock()
	defer l.mu.Unlock()
	lim, ok := l.m[key]
	if !ok {
		lim = rate.NewLimiter(l.every, l.burst)
		l.m[key] = lim
	}
	return lim
}

// wait — 0: urinish mumkin (token olindi). Aks holda — keyingi token qancha vaqtdan keyin.
// Reserve navbat oladi, Delay — shu navbatgacha kutish; biz kutmaymiz: Cancel tokenni qaytaradi.
func (l *limiters) wait(key string) time.Duration {
	res := l.get(key).Reserve()
	d := res.Delay()
	if d > 0 {
		res.Cancel()
	}
	return d
}

// limitByIP — middleware (GB-6.01 shakli): kalit — mijozning IP manzili. r.RemoteAddr — TCP ulanishning
// manzili («IP:port»), uni mijoz soxtalashtira olmaydi. X-Forwarded-For sarlavhasiga ishonmaymiz: uni
// har kim yozadi. Server proksi (Caddy, GB-8.06) ortida boʻlsa, RemoteAddr — proksining manzili: hamma mijoz bitta
// chelakka tushadi. U holda IP'ni faqat ishonchli proksi yozgan sarlavhadan olish kerak (GB-8.06).
// net/http serveri RemoteAddr'ni doim «IP:port» qiladi, shuning uchun SplitHostPort xatosi bu yerda kutilmaydi.
// Eslatma: IPv6'da bitta foydalanuvchida odatda butun /64 tarmoq boʻladi — production'da kalitni /64 boʻyicha olish mumkin.
func (l *limiters) limitByIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		if d := l.wait(ip); d > 0 {
			tooManyRequests(w, d)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// tooManyRequests — 429 problem va Retry-After butun soniyalarda, yuqoriga yaxlitlangan.
func tooManyRequests(w http.ResponseWriter, d time.Duration) {
	w.Header().Set("Retry-After", strconv.Itoa(int(math.Ceil(d.Seconds()))))
	writeProblem(w, newProblem(http.StatusTooManyRequests, "urinishlar koʻp: keyinroq qayta urinib koʻring"))
}

// sweep (🟡 ixtiyoriy vazifa) — toʻla chelaklarni map'dan oʻchiradi va nechtasi oʻchganini qaytaradi.
// Toʻla chelak — kalit uzoq vaqt urinmagan: yangisi baribir toʻla yaratiladi, natija oʻzgarmaydi.
// Serverda uni vaqti-vaqti bilan chaqirish kerak (masalan time.Ticker bilan goroutine'da, GB-7.01).
func (l *limiters) sweep(now time.Time) int {
	l.mu.Lock()
	defer l.mu.Unlock()
	n := 0
	for key, lim := range l.m {
		if lim.TokensAt(now) >= float64(l.burst) {
			delete(l.m, key) // range ichida delete mumkin (Go spec: For statements with range clause)
			n++
		}
	}
	return n
}
