// orta: qaysi qatorlar kompilyatsiya boʻladi? Avval oʻzingiz javob bering, keyin tekshiring:
//
//	go build -o /dev/null ./homework/orta
//
// Bu papka ATAYLAB kompilyatsiya boʻlmaydi. Javob — README'da.
package main

type Nuqta struct{ X, Y int }

func (n Nuqta) Uzunlik() int { return n.X*n.X + n.Y*n.Y }
func (n *Nuqta) Sur(d int)   { n.X += d }

func yangi() Nuqta   { return Nuqta{1, 2} }
func ikki() [2]Nuqta { return [2]Nuqta{} }

func main() {
	nuqtalar := []Nuqta{{1, 2}}
	m := map[string]Nuqta{"a": {}}
	mp := map[string]*Nuqta{"a": {}}

	nuqtalar[0].Sur(1)    // 1
	mp["a"].Sur(1)        // 2
	_ = yangi().Uzunlik() // 3
	_ = m["a"].Uzunlik()  // 4
	m["a"].X = 5          // 5
	f := m["a"].Sur       // 6
	ikki()[0].Sur(1)      // 7
	mp["a"].X = 5         // 8
	_ = f
}
