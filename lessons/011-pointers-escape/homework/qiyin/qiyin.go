// qiyin: yangi buyurtma qoʻshib, birinchisini "toʻlandi" deb belgilash.
package qiyin

type Buyurtma struct {
	Summa   int64
	Tolandi bool
}

// qoshXato: ataylab xato. p append'dan oldin olingan; append yangi massiv
// yaratsa, p eski massivga ishora qiladi va oʻzgarish yoʻqoladi.
func qoshXato(s []Buyurtma, yangi Buyurtma) []Buyurtma {
	p := &s[0]
	s = append(s, yangi)
	p.Tolandi = true
	return s
}

// qosh: toʻgʻri versiya — append'dan keyin indeks orqali.
func qosh(s []Buyurtma, yangi Buyurtma) []Buyurtma {
	s = append(s, yangi)
	s[0].Tolandi = true
	return s
}
