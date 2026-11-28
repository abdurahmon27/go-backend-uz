package store

import "context"

// Seed — GB-6.06'dagi seed.sql'ning ikki xaridori va ikki buyurtmasi, xuddi oʻsha argon2id satrlari bilan
// (GB-6.03, internal/password.Hash: m=47104 KiB, t=1, p=1). Parollar faqat shu laboratoriya uchun (README):
//
//	dilnoza@bozor.test — tarvuz-osmon-qalam-ariq-27   (id 1, buyurtma 41)
//	jasur@bozor.test   — choy-non-qand-27             (id 2, buyurtma 42)
func (s *Store) Seed(ctx context.Context) error {
	users := []struct{ email, hash string }{
		{"dilnoza@bozor.test", "$argon2id$v=19$m=47104,t=1,p=1$nOpo1eo3/LqYO44wrQWXCw$krxLcr34r7HmNQnM3UHQF6NHozzxymB8qQeRiGq+wQA"},
		{"jasur@bozor.test", "$argon2id$v=19$m=47104,t=1,p=1$ecjvB7e6NhL2MXutDLGfXQ$OkrQiftpgFvIuhs1dpTtHdLt8mGImGy4bFew0fqLJAc"},
	}
	for _, u := range users {
		if _, err := s.AddUser(ctx, u.email, u.hash); err != nil {
			return err
		}
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	// Buyurtma ID'lari 41 dan (GB-6.06): ishlab turgan doʻkondagidek, undan oldin ham buyurtmalar boʻlgan.
	s.orders[41] = Order{ID: 41, UserID: 1, Status: StatusNew, Note: "Chilonzor 9-kvartal, 12-uy"}
	s.orders[42] = Order{ID: 42, UserID: 2, Status: StatusNew, Note: "Yunusobod 4-kvartal, 17-uy, ertalab"}
	return nil
}
