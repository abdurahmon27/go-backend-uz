// Package store — bozor maʼlumotlari Postgres'da: foydalanuvchilar va session'lar (GB-6.03, GB-6.04),
// mahsulotlar va buyurtmalar (orders.go, GB-5.10). GB-6.06: har buyurtmaning egasi bor (orders.user_id) va
// buyurtmani oʻqiydigan yoki oʻzgartiradigan har metod egasini ham oladi; foydalanuvchida rol bor (users.role).
//
// Session — serverdagi yozuv: «shu token egasi — shu foydalanuvchi, shu vaqtgacha». Bazada tokenning
// oʻzi emas, faqat SHA-256 hash'i turadi (sessions.token_hash). Qatorni oʻchirish = session'ni bekor qilish.
package store

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// User — bozor foydalanuvchisi (GB-6.03). PasswordHash — password.Hash qaytargan argon2id satri;
// tag "-": bu maydon JSON'ga hech qachon chiqmaydi. Uni faqat UserByEmail toʻldiradi (login uchun).
// Role (GB-6.06) — RoleCustomer yoki RoleAdmin (users.role, CHECK bilan). Uni mijoz hech qachon bermaydi:
// yangi foydalanuvchi bazaning DEFAULT'i bilan xaridor boʻladi.
type User struct {
	ID           int64  `json:"id"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	PasswordHash string `json:"-"`
}

// LogValue — User'ning log'dagi koʻrinishi (GB-8.02 vazifasi; slog.LogValuer): faqat id va rol.
// Busiz JSON handler User'ni json.Marshal qiladi va json tag'li Email chiqadi. Email — shaxsiy maʼlumot,
// PasswordHash — secret (tag "-" uni JSON'dan yashiradi, lekin boshqa handler'dan emas): ikkalasi ham log'ga yozilmaydi.
func (u User) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Int64("id", u.ID),
		slog.String("role", u.Role),
	)
}

// Rollar — users.role ustunidagi qiymatlar.
const (
	RoleCustomer = "xaridor"
	RoleAdmin    = "admin"
)

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
// GB-8.02: xato matnida email yoʻq — 5xx xatoning matni log'ga toʻliq yoziladi (writeError), email esa shaxsiy
// maʼlumot (OWASP Logging Cheat Sheet). Qaysi soʻrov ekanini log'dagi request_id aytadi.
func (s *Store) AddUser(ctx context.Context, email, passwordHash string) (User, error) {
	u := User{Email: email}
	err := s.pool.QueryRow(ctx,
		"INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id, role",
		email, passwordHash).Scan(&u.ID, &u.Role)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return User{}, ErrEmailTaken
	}
	if err != nil {
		return User{}, fmt.Errorf("foydalanuvchi yozilmadi: %w", err)
	}
	return u, nil
}

// UserByEmail foydalanuvchini parol hash'i bilan qaytaradi (login tekshiruvi uchun). Xato matnida email yoʻq (GB-8.02).
func (s *Store) UserByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := s.pool.QueryRow(ctx,
		"SELECT id, email, role, password_hash FROM users WHERE lower(email) = lower($1)",
		email).Scan(&u.ID, &u.Email, &u.Role, &u.PasswordHash)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrUserNotFound
	}
	if err != nil {
		return User{}, fmt.Errorf("foydalanuvchi oʻqilmadi: %w", err)
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
		SELECT u.id, u.email, u.role
		FROM sessions s JOIN users u ON u.id = s.user_id
		WHERE s.token_hash = $1 AND s.expires_at > now()`,
		tokenHash).Scan(&u.ID, &u.Email, &u.Role)
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

// DeleteExpiredSessions muddati oʻtgan qatorlarni tozalaydi (GB-6.04'ning 🟡 vazifasi).
// UserBySession ularni baribir qabul qilmaydi; bu faqat jadval oʻsib ketmasligi uchun.
func (s *Store) DeleteExpiredSessions(ctx context.Context) (int64, error) {
	tag, err := s.pool.Exec(ctx, "DELETE FROM sessions WHERE expires_at <= now()")
	if err != nil {
		return 0, fmt.Errorf("eski session'lar oʻchirilmadi: %w", err)
	}
	return tag.RowsAffected(), nil
}

// Users — hamma foydalanuvchilar, ID tartibida (GET /admin/users uchun; parol hash'isiz).
func (s *Store) Users(ctx context.Context) ([]User, error) {
	rows, err := s.pool.Query(ctx, "SELECT id, email, role FROM users ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []User{}
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Email, &u.Role); err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, rows.Err()
}
