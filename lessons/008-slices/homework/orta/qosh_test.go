package orta

import (
	"slices"
	"testing"
)

func TestQosh(t *testing.T) {
	asos := make([]string, 0, 4)
	asos = append(asos, "non", "sut", "suv")
	ali := qosh(asos, "tuxum")
	vali := qosh(asos, "choy")
	if !slices.Equal(ali, []string{"non", "sut", "suv", "tuxum"}) {
		t.Errorf("ali = %v", ali)
	}
	if !slices.Equal(vali, []string{"non", "sut", "suv", "choy"}) {
		t.Errorf("vali = %v", vali)
	}
	if got := asos[:4][3]; got != "" {
		t.Errorf("asos massivining 4-joyiga yozilgan: %q", got)
	}
}
