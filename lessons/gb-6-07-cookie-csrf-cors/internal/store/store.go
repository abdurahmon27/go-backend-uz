// Package store — bozor foydalanuvchilari, session'lari va buyurtmalari (GB-6.07 laboratoriyasi).
//
// Bu laboratoriyada hammasi xotirada (map + mutex, GB-4.07): GB-6.07'ning mavzusi brauzer va server
// orasidagi qoidalar (cookie, CSRF, CORS, rate limit), bazaning oʻzi emas. Qoidalar GB-6.04/6.06'dagi
// Postgres versiyasi bilan bir xil:
//   - session — «shu token hash'ining egasi — shu foydalanuvchi, shu vaqtgacha»; tokenning oʻzi
//     saqlanmaydi, faqat SHA-256 hash'i (GB-6.04: sessions.token_hash);
//   - buyurtmani faqat egasi koʻradi va oʻzgartiradi; begona buyurtma yoʻq buyurtma bilan bir xil:
//     ErrOrderNotFound (GB-6.06: WHERE id = $1 AND user_id = $2).
//
// Server qayta ishga tushsa, xotira boshidan boshlanadi (Seed): videodagi hujumni qayta-qayta
// koʻrsatish uchun bu qulay — har ishga tushishda buyurtma 41 yana «yangi».
package store

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

// User — mijozga koʻrsatiladigan foydalanuvchi (GB-6.04). Parol hash'i JSON'ga hech qachon chiqmaydi.
type User struct {
	ID           int64  `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
}

// Order — buyurtma (GB-5.10/6.06'ning qisqa shakli: bu laboratoriyaga qatorlar va summa kerak emas).
type Order struct {
	ID     int64  `json:"id"`
	UserID int64  `json:"user_id"`
	Status string `json:"status"`
	Note   string `json:"note,omitempty"`
}

// Buyurtma holatlari (GB-6.06'dagi orders.status: DEFAULT 'yangi').
const (
	StatusNew       = "yangi"
	StatusCancelled = "bekor qilindi"
)

var (
	// ErrUserNotFound — bunday email'li foydalanuvchi yoʻq.
	ErrUserNotFound = errors.New("foydalanuvchi topilmadi")
	// ErrEmailTaken — bu email bilan foydalanuvchi allaqachon bor (GB-6.03).
	ErrEmailTaken = errors.New("bu email bilan foydalanuvchi bor")
	// ErrNoSession — bunday token'li session yoʻq yoki muddati oʻtgan (GB-6.04).
	ErrNoSession = errors.New("session yoʻq yoki muddati oʻtgan")
	// ErrOrderNotFound — bunday buyurtma yoʻq yoki u boshqa foydalanuvchiniki (GB-6.06: ikkalasi bir xil).
	ErrOrderNotFound = errors.New("buyurtma topilmadi")
	// ErrAlreadyCancelled — buyurtma allaqachon bekor qilingan.
	ErrAlreadyCancelled = errors.New("buyurtma allaqachon bekor qilingan")
)

type session struct {
	userID  int64
	expires time.Time
}

// Store — xotiradagi "jadvallar". Bir vaqtda koʻp soʻrov ishlatadi, shuning uchun mutex (GB-4.07).
// Metodlar GB-6.04/6.06'dagi Postgres Store bilan bir xil shaklda (ctx birinchi argument).
type Store struct {
	mu       sync.Mutex
	users    map[string]User    // lower(email) → foydalanuvchi (GB-6.04: users_email_lower_idx)
	sessions map[string]session // string(token hash) → session
	orders   map[int64]Order
	nextUser int64
}

// New — boʻsh Store.
func New() *Store {
	return &Store{
		users:    map[string]User{},
		sessions: map[string]session{},
		orders:   map[int64]Order{},
	}
}

// AddUser yangi foydalanuvchini yozadi (passwordHash — password.Hash natijasi). Email harf kattaligisiz
// noyob: «Dilnoza@bozor.test» va «dilnoza@bozor.test» — bitta foydalanuvchi.
func (s *Store) AddUser(ctx context.Context, email, passwordHash string) (User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := strings.ToLower(email)
	if _, ok := s.users[key]; ok {
		return User{}, fmt.Errorf("%s: %w", email, ErrEmailTaken)
	}
	s.nextUser++
	u := User{ID: s.nextUser, Email: email, PasswordHash: passwordHash}
	s.users[key] = u
	return u, nil
}

// UserByEmail foydalanuvchini parol hash'i bilan qaytaradi (login tekshiruvi uchun).
func (s *Store) UserByEmail(ctx context.Context, email string) (User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	u, ok := s.users[strings.ToLower(email)]
	if !ok {
		return User{}, fmt.Errorf("%s: %w", email, ErrUserNotFound)
	}
	return u, nil
}

// CreateSession yangi session yozadi: token hash'i, kimniki va qachongacha.
func (s *Store) CreateSession(ctx context.Context, tokenHash []byte, userID int64, expires time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[string(tokenHash)] = session{userID: userID, expires: expires}
	return nil
}

// UserBySession token hash'i boʻyicha session egasini qaytaradi. Muddati oʻtgan session topilmagan
// bilan bir xil: ErrNoSession (GB-6.04: WHERE token_hash = $1 AND expires_at > now()).
func (s *Store) UserBySession(ctx context.Context, tokenHash []byte) (User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	row, ok := s.sessions[string(tokenHash)]
	if !ok || !time.Now().Before(row.expires) {
		return User{}, ErrNoSession
	}
	for _, u := range s.users {
		if u.ID == row.userID {
			u.PasswordHash = ""
			return u, nil
		}
	}
	return User{}, ErrNoSession
}

// DeleteSession bitta session'ni oʻchiradi (logout). Yozuv yoʻq boʻlsa ham xato emas.
func (s *Store) DeleteSession(ctx context.Context, tokenHash []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, string(tokenHash))
	return nil
}

// DeleteUserSessions foydalanuvchining hamma session'larini oʻchiradi va nechtasi oʻchganini qaytaradi
// (GB-6.04 vazifasi: «hamma qurilmadan chiqish»).
func (s *Store) DeleteUserSessions(ctx context.Context, userID int64) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var n int64
	for k, row := range s.sessions {
		if row.userID == userID {
			delete(s.sessions, k) // range ichida delete mumkin (Go spec: For statements with range clause)
			n++
		}
	}
	return n, nil
}

// Order — foydalanuvchining buyurtmasi. Begona yoki yoʻq: ErrOrderNotFound (GB-6.06).
func (s *Store) Order(ctx context.Context, id, userID int64) (Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	o, ok := s.orders[id]
	if !ok || o.UserID != userID {
		return Order{}, fmt.Errorf("buyurtma %d, foydalanuvchi %d: %w", id, userID, ErrOrderNotFound)
	}
	return o, nil
}

// CancelOrder buyurtmani bekor qiladi — faqat egasi uchun (GB-6.06'dagi egalik sharti bilan).
// Postgres'da bu: UPDATE orders SET status = 'bekor qilindi' WHERE id = $1 AND user_id = $2 AND status = 'yangi'.
func (s *Store) CancelOrder(ctx context.Context, id, userID int64) (Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	o, ok := s.orders[id]
	if !ok || o.UserID != userID {
		return Order{}, fmt.Errorf("buyurtma %d, foydalanuvchi %d: %w", id, userID, ErrOrderNotFound)
	}
	if o.Status == StatusCancelled {
		return Order{}, fmt.Errorf("buyurtma %d: %w", id, ErrAlreadyCancelled)
	}
	o.Status = StatusCancelled
	s.orders[id] = o
	return o, nil
}
