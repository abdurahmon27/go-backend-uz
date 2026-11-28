package importer

import (
	"bufio"
	"fmt"
	"io"
)

// words — katalog nomlari; birinchi toʻrttasi GB-7.03 jadvalidagi tartibda: olma, non, oʻrik, choy.
var words = []string{"olma", "non", "oʻrik", "choy", "anor", "uzum", "gilos", "shaftoli", "qovun", "tarvuz",
	"yongʻoq", "bodom", "guruch", "sabzi", "piyoz", "kartoshka", "pomidor", "bodring", "qatiq", "sut"}

// Generate w'ga n qatorli sinov katalogini yozadi: «nom;narx;qoldiq», narx soʻmda. Nomlar takrorlanmaydi:
// dastlabki len(words) ta qator — soʻzning oʻzi, keyingilari «soʻz-raqam» (masalan «anor-25»).
func Generate(w io.Writer, n int) error {
	bw := bufio.NewWriter(w)
	head := [][2]int{{12000, 40}, {4000, 120}, {35000, 5}, {25000, 8}} // olma, non, oʻrik, choy
	for i := range n {
		name := words[i%len(words)]
		if i >= len(words) {
			name = fmt.Sprintf("%s-%d", name, i+1)
		}
		price, stock := 1000+(i*37)%50000, i%200
		if i < len(head) {
			price, stock = head[i][0], head[i][1]
		}
		if _, err := fmt.Fprintf(bw, "%s;%d;%d\n", name, price, stock); err != nil {
			return err
		}
	}
	return bw.Flush()
}
