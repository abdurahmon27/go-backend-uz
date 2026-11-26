// Uy vazifasi javobi (Guru): fayl/main.go, lekin rename'dan oldin fayl, keyin papka fsync qilinadi.
// Jarayon crash'iga rename yetarli edi; elektr oʻchishi yoki OS crash'ida esa maʼlumot hali
// page cache'da boʻlishi mumkin — fsync uni diskka tushiradi.
//
//	go build -o bin/fsync ./homework/fsync && time ./bin/fsync -usul qulf -mijoz A
//
// Asl izoh:
// fayl: buyurtmalarni bitta JSON faylda saqlaydigan "baza" (DB-1.01).
// Har bir buyurtma: faylni oʻqi → roʻyxatga qoʻsh → faylni qayta yoz.
//
//	./bin/fayl -mijoz A -n 1000              A mijozdan 1000 ta buyurtma yozadi
//	./bin/fayl -usul rename -mijoz A         vaqtinchalik fayl + rename bilan
//	./bin/fayl -usul qulf -mijoz A           flock + rename bilan
//	./bin/fayl -sana                         faylda nechta buyurtma borligini sanaydi
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"syscall"
)

type Buyurtma struct {
	ID    int    `json:"id"`
	Mijoz string `json:"mijoz"`
	Summa int    `json:"summa"`
}

type Ombor struct {
	Buyurtmalar []Buyurtma `json:"buyurtmalar"`
}

func oqi(fayl string) (Ombor, error) {
	var o Ombor
	b, err := os.ReadFile(fayl)
	if errors.Is(err, fs.ErrNotExist) {
		return o, nil // fayl yoʻq: ombor boʻsh
	}
	if err != nil {
		return o, err
	}
	err = json.Unmarshal(b, &o)
	return o, err
}

func yoz(fayl string, o Ombor) error {
	b, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(fayl, b, 0o644) // avval faylni boʻshatadi, keyin yozadi
}

// yozRename: avval yonidagi vaqtinchalik faylga yozadi, keyin rename bilan almashtiradi.
// rename atomik: oʻquvchi eski yoki yangi faylni koʻradi, yarimini emas.
func yozRename(fayl string, o Ombor) error {
	b, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		return err
	}
	tmp := fmt.Sprintf("%s.%d.tmp", fayl, os.Getpid())
	f, err := os.OpenFile(tmp, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	if _, err := f.Write(b); err != nil {
		f.Close()
		return err
	}
	if err := f.Sync(); err != nil { // 1) maʼlumot diskda
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	if err := os.Rename(tmp, fayl); err != nil {
		return err
	}
	d, err := os.Open(filepath.Dir(fayl)) // 2) yangi nom ham diskda: papkani fsync qilamiz
	if err != nil {
		return err
	}
	defer d.Close()
	return d.Sync()
}

// qulfla: hamma jarayonlar uchun umumiy qulf (flock). Qulf faylni ushlab turadi, ombor faylini emas.
func qulfla(fayl string) (*os.File, error) {
	f, err := os.OpenFile(fayl+".lock", os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return nil, err
	}
	return f, syscall.Flock(int(f.Fd()), syscall.LOCK_EX)
}

// tashxis: buzuq faylning boshidagi butun JSON qismini topadi va ortida nima qolganini aytadi.
func tashxis(fayl string) {
	b, _ := os.ReadFile(fayl)
	if len(b) == 0 {
		fmt.Println("fayl 0 bayt: hamma buyurtma yoʻq")
		return
	}
	var o Ombor
	dec := json.NewDecoder(bytes.NewReader(b))
	if dec.Decode(&o) != nil {
		fmt.Printf("fayl %d bayt, oxiri yozilmay qolgan\n", len(b))
		return
	}
	qoldi := len(b) - int(dec.InputOffset())
	fmt.Printf("boshida %d ta buyurtma, ortida yana %d bayt eski yozuv\n", len(o.Buyurtmalar), qoldi)
}

func main() {
	fayl := flag.String("fayl", "buyurtmalar.json", "ombor fayli")
	mijoz := flag.String("mijoz", "A", "kim buyurtma beryapti")
	n := flag.Int("n", 1000, "nechta buyurtma (0 = toʻxtatilguncha)")
	sana := flag.Bool("sana", false, "faqat sanash")
	usul := flag.String("usul", "oddiy", "oddiy | rename | qulf")
	flag.Parse()

	if *sana {
		o, err := oqi(*fayl)
		if err != nil {
			fmt.Println("fayl buzuq:", err)
			tashxis(*fayl)
			os.Exit(1)
		}
		jami := map[string]int{}
		for _, b := range o.Buyurtmalar {
			jami[b.Mijoz]++
		}
		fmt.Printf("faylda: %d ta buyurtma (A: %d, B: %d)\n", len(o.Buyurtmalar), jami["A"], jami["B"])
		return
	}

	saqla := yoz
	if *usul != "oddiy" {
		saqla = yozRename
	}
	buzuq := 0
	for i := 0; *n == 0 || i < *n; {
		var q *os.File
		if *usul == "qulf" {
			var err error
			if q, err = qulfla(*fayl); err != nil {
				fmt.Println(*mijoz+":", err)
				os.Exit(1)
			}
		}
		o, err := oqi(*fayl)
		if err != nil {
			q.Close() // nil boʻlsa ham xavfsiz: qulfni boʻshatadi
			buzuq++   // yarim yozilgan faylni oʻqidik: qayta urinamiz
			if buzuq > 100000 {
				fmt.Println(*mijoz+":", "fayl tiklanmadi:", err)
				os.Exit(1)
			}
			continue
		}
		o.Buyurtmalar = append(o.Buyurtmalar, Buyurtma{ID: len(o.Buyurtmalar) + 1, Mijoz: *mijoz, Summa: 25000})
		if err := saqla(*fayl, o); err != nil {
			fmt.Println(*mijoz+":", err)
			os.Exit(1)
		}
		q.Close() // qulf boʻshaydi
		i++
	}
	fmt.Printf("%s: %d ta yozdim, %d marta buzuq faylni oʻqidim\n", *mijoz, *n, buzuq)
}
