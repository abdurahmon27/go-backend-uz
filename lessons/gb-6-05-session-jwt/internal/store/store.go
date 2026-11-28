// Package store — bozor foydalanuvchilari, session'lari va bekor qilingan JWT'lar (GB-6.05).
//
// Bu laboratoriyada hammasi xotirada (map + mutex): GB-6.05'ning mavzusi session va JWT'ning
// farqi, bazaning oʻzi emas. Metodlar va qoidalar GB-6.04'dagi Postgres Store bilan bir xil:
//   - session — serverdagi yozuv «shu token hash'i egasi — shu foydalanuvchi, shu vaqtgacha»;
//     tokenning oʻzi saqlanmaydi, faqat SHA-256 hash'i (GB-6.04: sessions.token_hash);
//   - yozuvni oʻchirish = session'ni bekor qilish: keyingi soʻrovning oʻzidayoq 401.
//
// Server qayta ishga tushsa, xotiradagi hamma narsa yoʻqoladi; GB-6.04'da ular jadvallarda saqlanadi.
package store

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

// User — bozor foydalanuvchisi (GB-6.04 bilan bir xil). PasswordHash JSON'ga hech qachon chiqmaydi.
// GB-6.05: Email boʻsh boʻlsa ham JSON'ga tushmaydi (omitempty) — JWT yoʻli faqat tokendagi sub'ni (ID) biladi.
type User struct {
	ID           int64  `json:"id"`
	Email        string `json:"email,omitempty"`
	PasswordHash string `json:"-"`
}

// ErrUserNotFound — bunday email'li foydalanuvchi yoʻq.
var ErrUserNotFound = errors.New("foydalanuvchi topilmadi")

// ErrEmailTaken — bu email bilan foydalanuvchi allaqachon bor (GB-6.03).
var ErrEmailTaken = errors.New("bu email bilan foydalanuvchi bor")

// ErrNoSession — bunday token'li session yoʻq yoki muddati oʻtgan.
var ErrNoSession = errors.New("session yoʻq yoki muddati oʻtgan")

type sessionRow struct {
	userID  int64
	expires time.Time
}

// Store — xotiradagi "jadvallar". Bir vaqtda koʻp soʻrov ishlatadi, shuning uchun mutex (GB-4.07).
type Store struct {
	mu        sync.Mutex
	users     map[string]User       // lower(email) → foydalanuvchi (GB-6.04: users_email_lower_idx)
	emailByID map[int64]string      // id → lower(email) (GB-6.04: users primary key)
	sessions  map[string]sessionRow // string(token hash) → session
	revoked   map[string]time.Time  // jti → token'ning exp vaqti (VAZIFA, TASK.md)
	nextID    int64
}

// New — boʻsh Store.
func New() *Store {
	return &Store{
		users:     map[string]User{},
		emailByID: map[int64]string{},
		sessions:  map[string]sessionRow{},
		revoked:   map[string]time.Time{},
	}
}

// AddUser yangi foydalanuvchini yozadi (POST /users, GB-6.03). passwordHash — password.Hash natijasi.
// Email band boʻlsa (katta-kichik harf farqi hisobga olinmaydi) — ErrEmailTaken.
func (s *Store) AddUser(ctx context.Context, email, passwordHash string) (User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	k := strings.ToLower(email)
	if _, ok := s.users[k]; ok {
		return User{}, fmt.Errorf("email %q: %w", email, ErrEmailTaken)
	}
	s.nextID++
	s.users[k] = User{ID: s.nextID, Email: email, PasswordHash: passwordHash}
	s.emailByID[s.nextID] = k
	return User{ID: s.nextID, Email: email}, nil
}

// UserByEmail foydalanuvchini parol hash'i bilan qaytaradi (login tekshiruvi uchun).
func (s *Store) UserByEmail(ctx context.Context, email string) (User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	u, ok := s.users[strings.ToLower(email)]
	if !ok {
		return User{}, fmt.Errorf("email %q: %w", email, ErrUserNotFound)
	}
	return u, nil
}

// CreateSession yangi session yozadi: token hash'i, kimniki va qachongacha.
func (s *Store) CreateSession(ctx context.Context, tokenHash []byte, userID int64, expires time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[string(tokenHash)] = sessionRow{userID: userID, expires: expires}
	return nil
}

// UserBySession token hash'i boʻyicha session egasini qaytaradi. Muddati oʻtgan session
// topilmagan bilan bir xil: ErrNoSession (eskirgan yozuv shu yerda oʻchiriladi). Bu — har himoyalangan
// soʻrovdagi qidiruv: token hash'i → session → foydalanuvchi, ikkalasi ham map'dan kalit boʻyicha
// (GB-6.04'da: sessions va users jadvallarining JOIN'li bitta SELECT'i).
func (s *Store) UserBySession(ctx context.Context, tokenHash []byte) (User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	row, ok := s.sessions[string(tokenHash)]
	if !ok {
		return User{}, ErrNoSession
	}
	if !time.Now().Before(row.expires) {
		delete(s.sessions, string(tokenHash))
		return User{}, ErrNoSession
	}
	u, ok := s.users[s.emailByID[row.userID]]
	if !ok {
		return User{}, ErrNoSession
	}
	return User{ID: u.ID, Email: u.Email}, nil
}

// DeleteSession bitta session'ni oʻchiradi (logout). Yozuv yoʻq boʻlsa ham xato emas:
// natija bir xil — bu token bilan endi kirib boʻlmaydi.
func (s *Store) DeleteSession(ctx context.Context, tokenHash []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, string(tokenHash))
	return nil
}

// DeleteUserSessions foydalanuvchining hamma session'larini oʻchiradi ("hamma qurilmadan chiqish", GB-6.04)
// va nechtasi oʻchganini qaytaradi.
func (s *Store) DeleteUserSessions(ctx context.Context, userID int64) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var n int64
	for h, row := range s.sessions {
		if row.userID == userID {
			delete(s.sessions, h) // range ichida delete mumkin (Go spec: For statements with range clause)
			n++
		}
	}
	return n, nil
}

// errTask — asosiy vazifa hali yozilmagan (TASK.md, 1-qadam).
var errTask = errors.New("vazifa: RevokeToken hali yozilmagan (TASK.md)")

// RevokeToken JWT'ni bekor qilinganlar roʻyxatiga (denylist) qoʻshadi: kalit — tokenning jti'si,
// qiymat — uning exp vaqti (shu vaqtdan keyin token baribir yaroqsiz, yozuvni tozalash mumkin).
//
// VAZIFA (TASK.md, 1-qadam): s.revoked'ga yozing (mutex ostida) va nil qaytaring.
func (s *Store) RevokeToken(ctx context.Context, jti string, expires time.Time) error {
	return errTask
}

// TokenRevoked jti denylist'da bormi? Bu — JWT'ni bekor qilishning narxi: har soʻrovda yana bitta qidiruv.
//
// VAZIFA (TASK.md, 1-qadam): s.revoked'dan comma-ok bilan oʻqing (GB-2.02).
func (s *Store) TokenRevoked(ctx context.Context, jti string) (bool, error) {
	return false, nil
}

// errOptional — 🟡 ixtiyoriy vazifa hali yozilmagan.
var errOptional = errors.New("ixtiyoriy vazifa hali yozilmagan")

// PurgeRevoked exp vaqti now'dan oldin boʻlgan denylist yozuvlarini oʻchiradi va nechtasi oʻchganini
// qaytaradi (🟡 ixtiyoriy vazifa). Muddati oʻtgan token Parse'dan baribir oʻtmaydi: yozuv keraksiz.
func (s *Store) PurgeRevoked(ctx context.Context, now time.Time) (int, error) {
	return 0, errOptional
}
