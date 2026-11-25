package qiyin

import "testing"

// Xato faqat cap yetmaganda chiqadi, shuning uchun ikkala holatni tekshiramiz.
func TestQosh(t *testing.T) {
	jadval := []struct {
		nom string
		s   []Buyurtma
	}{
		{"cap yetadi", make([]Buyurtma, 1, 4)},
		{"cap yetmaydi", make([]Buyurtma, 1)},
	}
	for _, q := range jadval {
		t.Run(q.nom, func(t *testing.T) {
			s := qosh(q.s, Buyurtma{Summa: 300})
			if !s[0].Tolandi {
				t.Errorf("s[0].Tolandi = false, true kutilgan edi")
			}
		})
	}
}

// TestQosh'dagi jadvalni qoshXato bilan ishlatsangiz, "cap yetmaydi" holati
// yiqiladi. Bu test esa xato haqiqatan borligini koʻrsatadi (u oʻtadi):
// go test -run TestXatoBor -v ./homework/qiyin
func TestXatoBor(t *testing.T) {
	s := qoshXato(make([]Buyurtma, 1), Buyurtma{Summa: 300})
	if s[0].Tolandi {
		t.Errorf("kutilmagan: xato yoʻqolibdi")
	}
	t.Log("qoshXato: s[0].Tolandi = false — oʻzgarish eski massivda qoldi")
}
