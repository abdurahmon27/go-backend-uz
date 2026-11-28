package store

import "context"

// Seed — GB-6.04'dagi seed.sql'ning ikki foydalanuvchisi (GB-6.03'dagi Dilnoza va Jasur), xuddi oʻsha
// argon2id satrlari bilan (internal/password.Hash: m=47104 KiB, t=1, p=1). Parollar faqat shu laboratoriya uchun:
//
//	dilnoza@bozor.test — tarvuz-osmon-qalam-ariq-27   (id 1)
//	jasur@bozor.test   — choy-non-qand-27             (id 2)
func (s *Store) Seed(ctx context.Context) error {
	seed := []struct{ email, hash string }{
		{"dilnoza@bozor.test", "$argon2id$v=19$m=47104,t=1,p=1$nOpo1eo3/LqYO44wrQWXCw$krxLcr34r7HmNQnM3UHQF6NHozzxymB8qQeRiGq+wQA"},
		{"jasur@bozor.test", "$argon2id$v=19$m=47104,t=1,p=1$ecjvB7e6NhL2MXutDLGfXQ$OkrQiftpgFvIuhs1dpTtHdLt8mGImGy4bFew0fqLJAc"},
	}
	for _, u := range seed {
		if _, err := s.AddUser(ctx, u.email, u.hash); err != nil {
			return err
		}
	}
	return nil
}
