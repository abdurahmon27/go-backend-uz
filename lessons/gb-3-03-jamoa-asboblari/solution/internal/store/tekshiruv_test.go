package store

// Bu fayl vazifani (TASK.md) tekshiradi: go test ./... uni boshqa fayllar bilan birga ishga tushiradi.
// Uni oʻzgartirmang. Test store paketining ichida, shuning uchun products maydonini ham koʻradi.

import (
	"errors"
	"testing"
)

// katalog ikki mahsulotli yangi katalog qaytaradi: olma (ID 1) va non (ID 2).
func katalog() *Store {
	var st Store
	st.Add("olma", 1_200_000, 40)
	st.Add("non", 400_000, 100)
	return &st
}

// 1–2-qadam: Rename nomni almashtiradi va chetidagi boʻsh joylarni olib tashlaydi.
func TestVazifaRename(t *testing.T) {
	st := katalog()
	if err := st.Rename(1, "  oʻrik "); err != nil {
		t.Fatalf("Rename(1, «  oʻrik ») = %v; kutilgan nil (TASK.md)", err)
	}
	switch got := st.products[0].Name; got {
	case "oʻrik":
		// toʻgʻri
	case "olma":
		t.Error("Rename(1, «  oʻrik »)'dan keyin nom hali «olma»: nom almashmadi. go vet ./... nima deydi? (TASK.md, 1-qadam)")
	default:
		t.Errorf("Rename(1, «  oʻrik »)'dan keyin nom «%s»; kutilgan «oʻrik», chetidagi boʻsh joylarsiz. staticcheck ./... nima deydi? (TASK.md, 2-qadam)", got)
	}
	if got := st.products[1].Name; got != "non" {
		t.Errorf("Rename(1, …) ikkinchi mahsulotni ham oʻzgartirdi: «%s»; kutilgan «non»", got)
	}
}

// 2-qadam: faqat boʻsh joydan iborat nom — xato, katalog oʻzgarmaydi.
func TestVazifaRenameBosh(t *testing.T) {
	st := katalog()
	if err := st.Rename(2, "   "); err == nil {
		t.Error("Rename(2, «   ») = nil; kutilgan error: boʻsh joylar olib tashlansa, nom boʻsh qoladi (TASK.md, 2-qadam)")
	}
	if got := st.products[1].Name; got != "non" {
		t.Errorf("Rename(2, «   »)'dan keyin nom «%s»; kutilgan «non» (xatoda katalog oʻzgarmasin)", got)
	}
}

// Videodagi qoida: katalogda yoʻq ID — ErrNotFound (bu test boshidanoq oʻtadi).
func TestRenameTopilmadi(t *testing.T) {
	st := katalog()
	if err := st.Rename(9, "choy"); !errors.Is(err, ErrNotFound) {
		t.Errorf("Rename(9, «choy») = %v; kutilgan ErrNotFound", err)
	}
}
