package qiyin

import (
	"errors"
	"testing"
)

func TestMaydonXato(t *testing.T) {
	err := Tekshir("")
	if err == nil {
		t.Fatal("xato kutilgan edi")
	}
	if got := err.Error(); got != "buyurtma: manzil: qiymat boʻsh" {
		t.Errorf("matn: %q", got)
	}
	if !errors.Is(err, ErrBosh) {
		t.Error("errors.Is ErrBosh ni topmadi (Unwrap yoʻqmi?)")
	}
	mx, ok := errors.AsType[*MaydonXato](err)
	if !ok || mx.Maydon != "manzil" {
		t.Errorf("AsType: %v, %v", mx, ok)
	}
	if Tekshir("Toshkent") != nil {
		t.Error("toʻgʻri manzilda nil kutilgan edi")
	}
}
