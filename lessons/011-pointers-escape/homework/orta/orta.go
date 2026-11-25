// orta: summani fmt.Println'siz chiqaramiz, shunda hech narsa qochmaydi.
// go build -gcflags=-m ./homework/orta 2>&1 | grep -v inlin
// fmt.Println ...any oladi: int64 interface'ga joylanadi va "escapes to heap".
// strconv.AppendInt esa stack'dagi massivga yozadi.
package orta

import (
	"os"
	"strconv"
)

type Buyurtma struct {
	Summa   int64
	Tolandi bool
}

func chop(b Buyurtma) {
	var buf [24]byte // stack'da
	s := strconv.AppendInt(buf[:0], b.Summa, 10)
	s = append(s, '\n')
	os.Stdout.Write(s)
}
