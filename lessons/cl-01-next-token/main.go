// Kichik til modeli: korpusdagi soʻzlarni sanaydi va keyingi soʻzni bashorat qiladi.
// Claude ham har qadamda bitta keyingi tokenni tanlaydi; bu dastur oʻsha gʻoyaning eng
// kichik nusxasi: n-gram model. Token = soʻz, model = sanoq jadvali (korpusdan bir marta
// oʻrganiladi), kontekst oynasi = javob paytida model koʻradigan oxirgi -oyna ta token.
//
//	go run . -oyna 1 -t 0 python tilini
package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand/v2"
	"os"
	"sort"
	"strings"
)

// maxOyna: model oʻrganadigan eng uzun kontekst. -oyna undan katta boʻlolmaydi.
const maxOyna = 3

// model: kontekst (oxirgi k ta token) → keyingi token → necha marta uchradi.
type model struct {
	sanoq  map[string]map[string]int
	tartib map[string]int // token korpusda birinchi marta qayerda uchragan (teng holatda shu hal qiladi)
}

// oqit korpusdagi har bir oʻrin uchun 1..oyna uzunlikdagi kontekstlarni sanaydi.
func oqit(tokenlar []string, oyna int) model {
	m := model{sanoq: map[string]map[string]int{}, tartib: map[string]int{}}
	for i, tok := range tokenlar {
		if _, bor := m.tartib[tok]; !bor {
			m.tartib[tok] = i
		}
		for k := 1; k <= oyna && k <= i; k++ {
			kon := strings.Join(tokenlar[i-k:i], " ")
			if m.sanoq[kon] == nil {
				m.sanoq[kon] = map[string]int{}
			}
			m.sanoq[kon][tok]++
		}
	}
	return m
}

type nomzod struct {
	token string
	soni  int
	p     float64
}

// taqsimot: eng uzun tanish kontekstni topadi va temperature bilan ehtimolliklarni hisoblaydi.
// p ~ soni^(1/t). t = 0 boʻlsa modelning oʻz ehtimolligi (t = 1) koʻrsatiladi,
// tanlov esa har doim eng ehtimolli token (greedy).
func (m model) taqsimot(matn []string, oyna int, t float64) (string, []nomzod) {
	for k := min(oyna, len(matn)); k >= 1; k-- {
		kon := strings.Join(matn[len(matn)-k:], " ")
		keyingi, bor := m.sanoq[kon]
		if !bor {
			continue // bunday kontekst korpusda yoʻq: qisqaroq kontekstga oʻtamiz
		}
		var ns []nomzod
		for tok, soni := range keyingi {
			ns = append(ns, nomzod{token: tok, soni: soni})
		}
		sort.Slice(ns, func(i, j int) bool {
			if ns[i].soni != ns[j].soni {
				return ns[i].soni > ns[j].soni
			}
			return m.tartib[ns[i].token] < m.tartib[ns[j].token]
		})
		daraja := 1.0
		if t > 0 {
			daraja = 1 / t
		}
		var jami float64
		for i := range ns {
			ns[i].p = math.Pow(float64(ns[i].soni), daraja)
			jami += ns[i].p
		}
		for i := range ns {
			ns[i].p /= jami
		}
		return kon, ns
	}
	return "", nil
}

// tanla: ehtimollik boʻyicha bitta tokenni tasodifiy tanlaydi (sampling).
func tanla(ns []nomzod, r *rand.Rand) string {
	x := r.Float64()
	for _, n := range ns {
		if x < n.p {
			return n.token
		}
		x -= n.p
	}
	return ns[len(ns)-1].token
}

func main() {
	oyna := flag.Int("oyna", 1, "kontekst oynasi: model oxirgi nechta tokenni koʻradi")
	t := flag.Float64("t", 0, "temperature: 0 = har doim eng ehtimolli token")
	seed := flag.Uint64("seed", 1, "tasodif urugʻi (bir xil seed = bir xil natija)")
	n := flag.Int("n", 1, "nechta javob yaratish (n > 1 boʻlsa faqat javoblar chiqadi)")
	korpus := flag.String("korpus", "korpus.txt", "model oʻrganadigan matn")
	flag.Parse()

	matn, err := os.ReadFile(*korpus)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if *oyna < 1 || *oyna > maxOyna {
		fmt.Fprintf(os.Stderr, "-oyna 1..%d boʻlishi kerak\n", maxOyna)
		os.Exit(2)
	}
	m := oqit(strings.Fields(strings.ToLower(string(matn))), maxOyna)
	boshi := strings.Fields(strings.ToLower(strings.Join(flag.Args(), " ")))
	if len(boshi) == 0 {
		fmt.Fprintln(os.Stderr, "boshlanish soʻzlarini bering, masalan: go run . python tilini")
		os.Exit(2)
	}
	r := rand.New(rand.NewPCG(*seed, 0))

	for j := 0; j < *n; j++ {
		jumla := append([]string(nil), boshi...)
		for len(jumla) < 20 && jumla[len(jumla)-1] != "." {
			kon, ns := m.taqsimot(jumla, *oyna, *t)
			if ns == nil {
				fmt.Printf("%-18s → ?         korpusda yoʻq\n", "["+jumla[len(jumla)-1]+"]")
				break
			}
			var tok string
			if *t == 0 {
				tok = ns[0].token
			} else {
				tok = tanla(ns, r)
			}
			if *n == 1 {
				var qator []string
				for _, x := range ns {
					qator = append(qator, fmt.Sprintf("%s %.2f", x.token, x.p))
				}
				fmt.Printf("%-18s → %-9s %s\n", "["+kon+"]", tok, strings.Join(qator, "  "))
			}
			jumla = append(jumla, tok)
		}
		fmt.Println(strings.Join(jumla, " "))
	}
}
