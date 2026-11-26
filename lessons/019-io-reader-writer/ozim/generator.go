// Package ozim — bu darsda oʻzimiz yozgan io.Reader va io.Writer.
package ozim

import (
	"io"
	"strconv"
)

// Generator — oʻz io.Reader'imiz. Buyurtma qatorlarini oldindan saqlamaydi:
// Read chaqirilganda p ga sigʻadigan qatorlarni yasaydi va beradi.
type Generator struct {
	qoldi  int64  // yana nechta bayt yasash kerak
	id     int64  // oxirgi buyurtma raqami
	navbat []byte // yasalgan, lekin hali berilmagan baytlar
}

// NewGenerator kamida hajm bayt beradi (oxirgi qator butun qoladi, shuning uchun bir necha bayt ortiq boʻlishi mumkin).
func NewGenerator(hajm int64) *Generator { return &Generator{qoldi: hajm} }

func (g *Generator) Read(p []byte) (int, error) {
	if len(p) == 0 { // boʻsh chelak: hech narsa bermaymiz, lekin EOF ham emas
		return 0, nil
	}
	for len(g.navbat) < len(p) && g.qoldi > 0 {
		oldin := len(g.navbat)
		g.id++
		g.navbat = Yoz(g.navbat, g.id)
		g.qoldi -= int64(len(g.navbat) - oldin)
	}
	if len(g.navbat) == 0 {
		return 0, io.EOF
	}
	n := copy(p, g.navbat)
	g.navbat = g.navbat[:copy(g.navbat, g.navbat[n:])]
	return n, nil
}

// Yoz bitta "id,summa,tolandi" qatorini b ga qoʻshadi.
func Yoz(b []byte, id int64) []byte {
	b = strconv.AppendInt(b, id, 10)
	b = append(b, ',')
	b = strconv.AppendInt(b, 1000+id*7919%499000, 10)
	if id%4 == 0 {
		return append(b, ",0\n"...)
	}
	return append(b, ",1\n"...)
}
