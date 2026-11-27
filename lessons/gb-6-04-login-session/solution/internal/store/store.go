// Package store — bozor foydalanuvchilari va session'lari Postgres'da (GB-6.03, GB-6.04).
//
// Session — serverdagi yozuv: «shu token egasi — shu foydalanuvchi, shu vaqtgacha». Bazada tokenning
// oʻzi emas, faqat SHA-256 hash'i turadi (sessions.token_hash): baza sizib chiqsa ham, undagi
// qiymatni cookie'ga qoʻyib kirib boʻlmaydi. Qatorni oʻchirish = session'ni bekor qilish (logout).
package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// User — bozor foydalanuvchisi (GB-6.03 bilan bir xil). PasswordHash — password.Hash qaytargan argon2id
// satri; tag "-": bu maydon JSON'ga hech qachon chiqmaydi. Uni faqat UserByEmail toʻldiradi (login uchun).
type User struct {
	ID           int64  `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
}

// ErrUserNotFound — bunday email'li foydalanuvchi yoʻq.
var ErrUserNotFound = errors.New("foydalanuvchi topilmadi")

// ErrEmailTaken — bu email bilan foydalanuvchi allaqachon bor (GB-6.03; bu yerda users_email_lower_idx).
var ErrEmailTaken = errors.New("bu email bilan foydalanuvchi bor")

// ErrNoSession — bunday token'li session yoʻq yoki muddati oʻtgan.
var ErrNoSession = errors.New("session yoʻq yoki muddati oʻtgan")

// Store — hamma SQL shu yerda. Pool'ni bir vaqtda koʻp goroutine ishlatishi mumkin (GB-5.05).
type Store struct {
	pool *pgxpool.Pool
}

// New — pool ustidagi Store.
func New(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool}
}

// AddUser yangi foydalanuvchini yozadi (POST /users, GB-6.03). passwordHash — password.Hash natijasi.
// Email band bo'lsa (unique indeks, SQLSTATE 23505) — ErrEmailTaken.
func (s *Store) AddUser(ctx context.Context, email, passwordHash string) (User, error) {
	u := User{Email: email}
	err := s.pool.QueryRow(ctx,
		"INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id",
		email, passwordHash).Scan(&u.ID)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return User{}, fmt.Errorf("email %q: %w", email, ErrEmailTaken)
	}
	if err != nil {
		return User{}, fmt.Errorf("foydalanuvchi %q yozilmadi: %w", email, err)
	}
	return u, nil
}

// UserByEmail foydalanuvchini parol hash'i bilan qaytaradi (login tekshiruvi uchun).
func (s *Store) UserByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := s.pool.QueryRow(ctx,
		"SELECT id, email, password_hash FROM users WHERE lower(email) = lower($1)",
		email).Scan(&u.ID, &u.Email, &u.PasswordHash)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, fmt.Errorf("email %q: %w", email, ErrUserNotFound)
	}
	if err != nil {
		return User{}, fmt.Errorf("email %q: %w", email, err)
	}
	return u, nil
}

// CreateSession yangi session yozadi: token hash'i, kimniki va qachongacha.
func (s *Store) CreateSession(ctx context.Context, tokenHash []byte, userID int64, expires time.Time) error {
	_, err := s.pool.Exec(ctx,
		"INSERT INTO sessions (token_hash, user_id, expires_at) VALUES ($1, $2, $3)",
		tokenHash, userID, expires)
	if err != nil {
		return fmt.Errorf("session yozilmadi: %w", err)
	}
	return nil
}

// UserBySession token hash'i boʻyicha session egasini qaytaradi. Muddati oʻtgan session
// topilmagan bilan bir xil: ErrNoSession. Vaqtni Postgres'ning now() funksiyasi beradi.
func (s *Store) UserBySession(ctx context.Context, tokenHash []byte) (User, error) {
	var u User
	err := s.pool.QueryRow(ctx, `
		SELECT u.id, u.email
		FROM sessions s JOIN users u ON u.id = s.user_id
		WHERE s.token_hash = $1 AND s.expires_at > now()`,
		tokenHash).Scan(&u.ID, &u.Email)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrNoSession
	}
	if err != nil {
		return User{}, fmt.Errorf("session oʻqilmadi: %w", err)
	}
	return u, nil
}

// DeleteSession bitta session'ni oʻchiradi (logout). Qator yoʻq boʻlsa ham xato emas:
// natija bir xil — bu token bilan endi kirib boʻlmaydi.
func (s *Store) DeleteSession(ctx context.Context, tokenHash []byte) error {
	_, err := s.pool.Exec(ctx, "DELETE FROM sessions WHERE token_hash = $1", tokenHash)
	if err != nil {
		return fmt.Errorf("session oʻchirilmadi: %w", err)
	}
	return nil
}

// DeleteUserSessions foydalanuvchining hamma session'larini oʻchiradi ("hamma qurilmadan chiqish")
// va nechta qator oʻchganini qaytaradi.
func (s *Store) DeleteUserSessions(ctx context.Context, userID int64) (int64, error) {
	tag, err := s.pool.Exec(ctx, "DELETE FROM sessions WHERE user_id = $1", userID)
	if err != nil {
		return 0, fmt.Errorf("foydalanuvchi %d session'lari oʻchirilmadi: %w", userID, err)
	}
	return tag.RowsAffected(), nil
}

// DeleteExpiredSessions muddati oʻtgan qatorlarni tozalaydi (🟡 ixtiyoriy vazifa).
// UserBySession ularni baribir qabul qilmaydi; bu faqat jadval oʻsib ketmasligi uchun.
func (s *Store) DeleteExpiredSessions(ctx context.Context) (int64, error) {
	tag, err := s.pool.Exec(ctx, "DELETE FROM sessions WHERE expires_at <= now()")
	if err != nil {
		return 0, fmt.Errorf("eski session'lar oʻchirilmadi: %w", err)
	}
	return tag.RowsAffected(), nil
}

// errOptional — 🟡 ixtiyoriy vazifa hali yozilmagan (starter'da DeleteExpiredSessions shuni qaytaradi).
var errOptional = errors.New("ixtiyoriy vazifa hali yozilmagan")
