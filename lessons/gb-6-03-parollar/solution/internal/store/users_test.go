package store

import (
	"errors"
	"testing"
)

func TestUsers(t *testing.T) {
	var u Users
	a, err := u.Add("dilnoza@bozor.test", "$argon2id$…")
	if err != nil || a.ID != 1 {
		t.Fatalf("Add: %+v, %v", a, err)
	}
	if _, err := u.Add("DILNOZA@bozor.test", "$argon2id$…"); !errors.Is(err, ErrEmailTaken) {
		t.Errorf("takroriy email (boshqa harf bilan): %v, kutilgan ErrEmailTaken", err)
	}
	if got, err := u.ByEmail("Dilnoza@Bozor.test"); err != nil || got.ID != 1 {
		t.Errorf("ByEmail: %+v, %v", got, err)
	}
	if _, err := u.ByEmail("jasur@bozor.test"); !errors.Is(err, ErrUserNotFound) {
		t.Errorf("ByEmail(yoʻq): %v, kutilgan ErrUserNotFound", err)
	}
}
