package store

// Foydalanuvchilar (GB-6.03). Hozircha xotirada, katalog kabi bitta qulf ostida: bu dars parol
// saqlash haqida, Postgres emas. Postgres'da xuddi shu satr users jadvalining password_hash
// (text) ustunida turadi.

import (
	"errors"
	"strings"
	"sync"
)

// User — bozor foydalanuvchisi. PasswordHash — password.Hash qaytargan argon2id satri;
// parolning oʻzi hech qayerda saqlanmaydi. Tag "-": bu maydon JSON'ga hech qachon chiqmaydi.
type User struct {
	ID           int64  `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
}

// ErrEmailTaken — bu email bilan foydalanuvchi allaqachon bor.
var ErrEmailTaken = errors.New("bu email bilan foydalanuvchi bor")

// ErrUserNotFound — bunday email bilan foydalanuvchi yoʻq.
var ErrUserNotFound = errors.New("foydalanuvchi topilmadi")

// Users — foydalanuvchilar roʻyxati. Nol qiymati tayyor; nusxalamang, doim *Users bilan ishlang.
type Users struct {
	mu   sync.Mutex
	list []User
}

// Add yangi foydalanuvchi qoʻshadi: ID = roʻyxatdagi soni + 1. Email katta-kichik harfga qaramay
// solishtiriladi (Dilnoza@… va dilnoza@… — bitta odam).
func (u *Users) Add(email, passwordHash string) (User, error) {
	u.mu.Lock()
	defer u.mu.Unlock()
	for _, x := range u.list {
		if strings.EqualFold(x.Email, email) {
			return User{}, ErrEmailTaken
		}
	}
	user := User{ID: int64(len(u.list)) + 1, Email: email, PasswordHash: passwordHash}
	u.list = append(u.list, user)
	return user, nil
}

// ByEmail foydalanuvchini email boʻyicha topadi (login uchun, GB-6.04).
func (u *Users) ByEmail(email string) (User, error) {
	u.mu.Lock()
	defer u.mu.Unlock()
	for _, x := range u.list {
		if strings.EqualFold(x.Email, email) {
			return x, nil
		}
	}
	return User{}, ErrUserNotFound
}
