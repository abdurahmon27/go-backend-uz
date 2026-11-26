// olcha: mikrosekund va millisekund darajasidagi ishlarni BITTALAB o'lchaydi.
// Bunday ish soatga qarashdan (time.Now ≈ 30 ns) ming barobar uzoq, shuning uchun har birini
// alohida o'lchab, ko'p marta takrorlaymiz va natijalarni tartiblaymiz:
//
//	  odatda  = mediana: natijalarning yarmi undan tez, yarmi sekin;
//	  oraliq  = eng tez 10% va eng sekin 10% tashlab yuborilgandan keyin qolganlari (p10 … p90).
//
//		go run ./olcha disk       SSD'dan 4 KiB o'qish (O_DIRECT) va xotiradagi nusxadan (page cache)
//		go run ./olcha fsync      4 KiB yozish: oddiy va fsync bilan
//		go run ./olcha tarmoq     Wi-Fi orqali routergacha borib-kelish (round trip)
//		go run ./olcha internet   go.dev gacha borib-kelish va to'liq HTTPS so'rov bosqichlari
//		go run ./olcha hammasi    hammasi ketma-ket (olchov.sh shuni ishlatadi)
//
// Disk sinovi joriy papkada 1 GiB lik sinov.dat faylini yaratadi (birinchi marta bir necha soniya).
// O_DIRECT faqat Linux'da bor: boshqa tizimda "SSD" qatori xotiradagi nusxani ham o'lchab qo'yishi mumkin.
package main

import (
	"context"
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"io"
	"math/rand/v2"
	"net"
	"net/http"
	"net/http/httptrace"
	"os"
	"slices"
	"strings"
	"syscall"
	"time"
)

var (
	faylYoli = flag.String("fayl", "sinov.dat", "disk sinovi uchun fayl (1 GiB)")
	lanAddr  = flag.String("lan", "", "lokal tarmoqdagi manzil host:port (standart: default gateway, 53-port)")
	host     = flag.String("host", "go.dev", "internet sinovi uchun server")
	oraliq   = flag.Duration("oraliq", 100*time.Millisecond, "tarmoq o'lchovlari orasidagi pauza")
)

func main() {
	flag.Parse()
	switch flag.Arg(0) {
	case "disk":
		disk()
	case "fsync":
		fsync()
	case "tarmoq":
		tarmoq()
	case "internet":
		internet()
	case "hammasi":
		disk()
		fsync()
		tarmoq()
		internet()
	default:
		fmt.Fprintln(os.Stderr, "ishlatish: go run ./olcha disk|fsync|tarmoq|internet|hammasi")
		os.Exit(2)
	}
}

// olcha: ishni n marta bajarib, har birining vaqtini yozib boradi.
func olcha(nom string, n int, ish func(i int)) natija {
	r := natija{nom: nom}
	for i := range n {
		boshi := time.Now()
		ish(i)
		r.t = append(r.t, time.Since(boshi))
	}
	return r
}

// olchaTarmoq: olcha kabi, lekin javob kelmagan urinish (ok == false) vaqtga qo'shilmaydi, sanaladi.
func olchaTarmoq(nom string, n int, ish func() (ok bool)) natija {
	r := natija{nom: nom}
	for range n {
		boshi := time.Now()
		ok := ish()
		d := time.Since(boshi)
		if ok {
			r.t = append(r.t, d)
		} else {
			r.yoq++
		}
		time.Sleep(*oraliq)
	}
	return r
}

// natija: bitta o'lchovning hamma vaqtlari.
type natija struct {
	nom string
	t   []time.Duration
	yoq int // tarmoq: javobsiz qolgan urinishlar
}

// p: tartiblangan vaqtlarning q-qismi, eng yaqin o'rin bo'yicha (0.5 = mediana; juft n da yuqori o'rtadagi).
func (n natija) p(q float64) time.Duration {
	s := slices.Clone(n.t)
	slices.Sort(s)
	return s[int(q*float64(len(s)-1)+0.5)]
}

var sarlavha bool

func chiqar(n natija) {
	if len(n.t) == 0 {
		fmt.Printf("%-22s: birorta ham javob kelmadi\n", n.nom)
		return
	}
	if !sarlavha {
		fmt.Printf("%-22s %6s %8s   %s\n", "o'lchov", "marta", "odatda", "oraliq")
		sarlavha = true
	}
	fmt.Printf("%-22s %6d %8s   %s … %s\n", n.nom, len(n.t), qisqa(n.p(0.5)), qisqa(n.p(0.1)), qisqa(n.p(0.9)))
	if n.yoq > 0 {
		fmt.Printf("  (yana %d ta urinish 3 soniyada javob olmadi)\n", n.yoq)
	}
}

// qisqa: 74.2µs, 1.93ms, 512ns: uchta raqamgacha.
func qisqa(d time.Duration) string {
	f := float64(d)
	for _, u := range []struct {
		b float64
		s string
	}{{1e9, "s"}, {1e6, "ms"}, {1e3, "µs"}, {1, "ns"}} {
		if f >= u.b || u.s == "ns" {
			x := f / u.b
			switch {
			case x >= 100:
				return fmt.Sprintf("%.0f%s", x, u.s)
			case x >= 10:
				return fmt.Sprintf("%.1f%s", x, u.s)
			default:
				return fmt.Sprintf("%.2f%s", x, u.s)
			}
		}
	}
	return d.String()
}

func tekshir(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "xato:", err)
		os.Exit(1)
	}
}

const (
	kb4  = 4 << 10
	mib  = 1 << 20
	hajm = 1 << 30 // sinov fayli: 1 GiB
)

// tayyorla: 1 GiB lik faylni tasodifiy baytlar bilan to'ldiradi (bo'sh "teshik" emas, haqiqiy bloklar).
func tayyorla() {
	if st, err := os.Stat(*faylYoli); err == nil && st.Size() == hajm {
		return
	}
	f, err := os.Create(*faylYoli)
	tekshir(err)
	buf := make([]byte, mib)
	r := rand.New(rand.NewPCG(1, 2))
	for range hajm / mib {
		for i := 0; i < len(buf); i += 8 {
			v := r.Uint64()
			for j := range 8 {
				buf[i+j] = byte(v >> (8 * j))
			}
		}
		_, err := f.Write(buf)
		tekshir(err)
	}
	tekshir(f.Sync())
	tekshir(f.Close())
}

func disk() {
	tayyorla()
	r := rand.New(rand.NewPCG(3, 4))

	// 1) SSD'dan 4 KiB, tasodifiy joydan. O_DIRECT xotiradagi nusxani (page cache) chetlab o'tadi:
	// har bir o'qish haqiqatan SSD'ga boradi.
	f, buf, err := ochDirect(*faylYoli, kb4)
	tekshir(err)
	chiqar(olcha("SSD'dan 4 KiB", 5000, func(int) {
		_, err := f.ReadAt(buf, int64(r.IntN(hajm/kb4))*kb4)
		tekshir(err)
	}))
	f.Close()

	// 2) Xuddi shu 4 KiB, oddiy open bilan: avval 16 MiB ni bir marta o'qiymiz,
	// Linux uni operativ xotirada nusxa qilib saqlaydi, keyingi o'qishlar nusxadan.
	g, err := os.Open(*faylYoli)
	tekshir(err)
	defer g.Close()
	_, err = g.ReadAt(make([]byte, 16*mib), 0)
	tekshir(err)
	b := make([]byte, kb4)
	chiqar(olcha("nusxadan 4 KiB", 5000, func(int) {
		_, err := g.ReadAt(b, int64(r.IntN(16*mib/kb4))*kb4)
		tekshir(err)
	}))
}

func fsync() {
	f, err := os.Create("fsync.dat")
	tekshir(err)
	defer os.Remove("fsync.dat")
	defer f.Close()
	b := make([]byte, kb4)
	for i := range 256 { // 1 MiB joy: keyin shu joyga qayta-qayta yozamiz
		_, err := f.WriteAt(b, int64(i)*kb4)
		tekshir(err)
	}
	tekshir(f.Sync())

	chiqar(olcha("4 KiB yozish", 2000, func(i int) {
		_, err := f.WriteAt(b, int64(i%256)*kb4)
		tekshir(err)
	}))
	chiqar(olcha("4 KiB yozish + fsync", 300, func(i int) {
		_, err := f.WriteAt(b, int64(i%256)*kb4)
		tekshir(err)
		tekshir(f.Sync()) // Linux'da fsync(2): ma'lumot diskka yetib borguncha kutadi
	}))
}

func tarmoq() {
	// Wi-Fi orqali routergacha (bizda: telefon hotspot'i). TCP ulanish = SYN ketadi, SYN-ACK qaytadi:
	// bitta borib-kelish.
	addr := *lanAddr
	if addr == "" {
		gw, err := gateway()
		tekshir(err)
		addr = net.JoinHostPort(gw, "53")
	}
	ulan(addr) // isitish: hisobga kirmaydi
	chiqar(olchaTarmoq("Wi-Fi: borib-kelish", 100, func() bool { return ulan(addr) }))
}

// ulan: TCP ulanish o'rnatib, darhol yopadi (vaqtni olchaTarmoq o'lchaydi).
// Port yopiq bo'lsa, qurilma SYN'ga RST bilan javob beradi ("connection refused"):
// bu ham bitta to'liq borib-kelish. 3 soniyada javob kelmasa (paket yo'qoldi), false.
func ulan(addr string) bool {
	c, err := net.DialTimeout("tcp", addr, 3*time.Second)
	if errors.Is(err, syscall.ECONNREFUSED) {
		return true
	}
	var ne net.Error
	if errors.As(err, &ne) && ne.Timeout() {
		return false
	}
	tekshir(err)
	c.Close()
	return true
}

func internet() {
	// 1) Bitta borib-kelish: TCP ulanish. DNS bir marta, oldindan: o'lchovga kirmasin.
	// Faqat IPv4 (tcp4): hamma o'lchov bitta yo'ldan borsin.
	ips, err := net.DefaultResolver.LookupIP(context.Background(), "ip4", *host)
	tekshir(err)
	addr := net.JoinHostPort(ips[0].String(), "443")
	ulan(addr) // isitish: birinchi urinish (radio uyg'onishi va h.k.) hisobga kirmaydi
	chiqar(olchaTarmoq("internet: borib-kelish", 30, func() bool { return ulan(addr) }))

	// 2) To'liq HTTPS so'rov, har safar yangi ulanish bilan. httptrace har bosqichga vaqt belgisi qo'yadi.
	jami, dns, tcp, tlsN, javob := natija{nom: "sahifa so'rovi"}, natija{nom: "  manzil topish"},
		natija{nom: "  ulanish"}, natija{nom: "  shifrlash"}, natija{nom: "  javob kutish"}
	url := "https://" + *host + "/robots.txt"
	tr4 := func() *http.Transport {
		var d net.Dialer
		return &http.Transport{
			DisableKeepAlives: true,
			DialContext: func(ctx context.Context, _, addr string) (net.Conn, error) {
				return d.DialContext(ctx, "tcp4", addr)
			},
		}
	}
	for i := range 11 {
		// bosqichlar shu so'rovning o'zgaruvchilariga yoziladi, so'rov muvaffaqiyatli bo'lsagina saqlanadi
		var dns0, tcp0, tls0, yubordi time.Time
		var dD, tD, sD, jD time.Duration
		iz := &httptrace.ClientTrace{
			// DNS: go.dev nomidan IP manzilni topish
			DNSStart: func(httptrace.DNSStartInfo) { dns0 = time.Now() },
			DNSDone:  func(httptrace.DNSDoneInfo) { dD = time.Since(dns0) },
			// TCP: ulanish (bitta borib-kelish)
			ConnectStart: func(string, string) { tcp0 = time.Now() },
			ConnectDone:  func(string, string, error) { tD = time.Since(tcp0) },
			// TLS: shifrlangan kanal haqida kelishish
			TLSHandshakeStart: func() { tls0 = time.Now() },
			TLSHandshakeDone:  func(tls.ConnectionState, error) { sD = time.Since(tls0) },
			// so'rov yuborilgandan javobning birinchi baytigacha
			WroteRequest:         func(httptrace.WroteRequestInfo) { yubordi = time.Now() },
			GotFirstResponseByte: func() { jD = time.Since(yubordi) },
		}
		req, err := http.NewRequest("GET", url, nil)
		tekshir(err)
		req = req.WithContext(httptrace.WithClientTrace(req.Context(), iz))
		klient := &http.Client{Transport: tr4(), Timeout: 10 * time.Second}
		boshi := time.Now()
		res, err := klient.Do(req)
		if err != nil { // tarmoq uzildi: urinish sanaladi, vaqt yozilmaydi
			jami.yoq++
			time.Sleep(*oraliq)
			continue
		}
		io.Copy(io.Discard, res.Body)
		res.Body.Close()
		d := time.Since(boshi)
		if i > 0 { // 0-so'rov isitish uchun (sertifikatlar yuklanadi): hisobga kirmaydi
			jami.t = append(jami.t, d)
			dns.t, tcp.t, tlsN.t, javob.t = append(dns.t, dD), append(tcp.t, tD), append(tlsN.t, sD), append(javob.t, jD)
		}
		time.Sleep(*oraliq)
	}
	for _, n := range []natija{jami, dns, tcp, tlsN, javob} {
		chiqar(n)
	}
}

// gateway: /proc/net/route dan default route'ning shlyuzi (Linux; manzil hex'da, little-endian: amd64/arm64).
func gateway() (string, error) {
	d, err := os.ReadFile("/proc/net/route")
	if err != nil {
		return "", fmt.Errorf("gateway topilmadi (%v): -lan host:port bering", err)
	}
	for _, q := range strings.Split(string(d), "\n")[1:] {
		f := strings.Fields(q)
		if len(f) > 2 && f[1] == "00000000" {
			var h uint32
			fmt.Sscanf(f[2], "%X", &h)
			return net.IPv4(byte(h), byte(h>>8), byte(h>>16), byte(h>>24)).String(), nil
		}
	}
	return "", fmt.Errorf("default route yo'q: -lan host:port bering")
}
