package main

// Bu fayl vazifani (TASK.md) tekshiradi: go test uni main.go bilan birga ishga tushiradi.
// Uni oʻzgartirmang. Bu yerda loadCatalog fayl bilan emas, matn bilan chaqiriladi:
//   strings.NewReader(s) — s matnidan oʻqiydigan io.Reader (videodagi kabi)
//   strings.HasPrefix(s, p) — s matni p bilan boshlanadimi
//   strings.Repeat(s, n) — s'ni n marta takrorlagan matn
//   *strconv.NumError — strconv.Atoi qaytaradigan xato turi («invalid syntax»);
//   errors.As uni oʻralgan xato ichidan topadi (GB-2.08)
//   bufio.ErrTooLong — Scanner'ning «token too long» xatosi (sentinel, errors.Is bilan)
//   t.Helper() — checkProducts yordamchi funksiya: xato xabari shu yerni emas, uni chaqirgan qatorni koʻrsatadi
//   check func(err error) string — struct maydoni ham funksiya boʻlishi mumkin (GB-2.09'dagi fn func() kabi);
//   har bir holat oʻz tekshiruvini shu maydonda olib keladi, boʻsh matn qaytarsa — hammasi joyida

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
	"testing"
)

func checkProducts(t *testing.T, got, want []Product, step string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("%d ta mahsulot yuklandi, kerak %d: %v (%s)", len(got), len(want), got, step)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("%d-mahsulot %+v, kerak %+v (%s)", i+1, got[i], want[i], step)
		}
	}
}

var videoCatalog = []Product{
	{ID: 1, Name: "olma", Price: 1_200_000, Stock: 40},
	{ID: 2, Name: "oʻrik", Price: 3_500_000, Stock: 5},
	{ID: 3, Name: "non", Price: 400_000, Stock: 100},
	{ID: 4, Name: "choy", Price: 2_500_000, Stock: 8},
}

// Videodagi loadCatalog: fayldan ham, matndan ham bir xil oʻqiydi (bu test boshidanoq oʻtadi).
func TestVideo(t *testing.T) {
	f, err := os.Open("katalog.txt")
	if err != nil {
		t.Fatalf("katalog.txt ochilmadi: %v", err)
	}
	defer f.Close()
	got, err := loadCatalog(f)
	if err != nil {
		t.Fatalf("loadCatalog(katalog.txt) xato qaytardi: %v", err)
	}
	checkProducts(t, got, videoCatalog, "videodagi katalog.txt")

	r := strings.NewReader("1;olma;12000;40\n2;oʻrik;35000;5\n3;non;4000;100\n4;choy;25000;8")
	got, err = loadCatalog(r)
	if err != nil {
		t.Fatalf("loadCatalog(strings.NewReader(…)) xato qaytardi: %v", err)
	}
	checkProducts(t, got, videoCatalog, "oxirgi qatorda \\n yoʻq — Scanner uni ham beradi")
}

// 1-qadam: boʻsh qatorlar va # bilan boshlangan izohlar tashlab oʻtiladi.
func TestVazifaIzoh(t *testing.T) {
	input := "# bozor katalogi, 2026-11-27\n" +
		"1;olma;12000;40\n" +
		"\n" +
		"   \n" +
		"  # ichimliklar\n" +
		"4;choy;25000;8\n"
	got, err := loadCatalog(strings.NewReader(input))
	if err != nil {
		t.Fatalf("loadCatalog izohli katalogda xato qaytardi («%v»): boʻsh qatorlarni va # bilan boshlangan qatorlarni tashlab oʻting (TASK.md, 1-qadam)", err)
	}
	checkProducts(t, got, []Product{videoCatalog[0], videoCatalog[3]}, "TASK.md, 1-qadam")
}

// 2-qadam: xato qaysi qatorda ekani aytiladi, asl xato esa %w bilan ichida qoladi.
func TestVazifaQatorRaqami(t *testing.T) {
	long := "3;non " + strings.Repeat("yangi uzilgan issiq non, ", 2800) + ";4000;100\n"
	tests := []struct {
		name   string
		input  string
		prefix string
		check  func(err error) string // boʻsh qator — hammasi joyida
	}{
		{
			name:   "narx_son_emas",
			input:  "1;olma;12000;40\n2;non;abc;100\n",
			prefix: "2-qator: ",
			check: func(err error) string {
				var ne *strconv.NumError
				if !errors.As(err, &ne) {
					return "ichida strconv xatosi yoʻq: xatoni %w bilan oʻrang, %v bilan emas"
				}
				return ""
			},
		},
		{
			name:   "izohdan_keyin",
			input:  "# bozor katalogi\n1;olma;12000;40\n\n4;non;abc;100\n",
			prefix: "4-qator: ",
			check: func(err error) string {
				var ne *strconv.NumError
				if !errors.As(err, &ne) {
					return "ichida strconv xatosi yoʻq: xatoni %w bilan oʻrang, %v bilan emas"
				}
				return ""
			},
		},
		{
			name:   "maydon_yetmaydi",
			input:  "1;olma;12000;40\n2;oʻrik;35000;5\n3;non;4000\n",
			prefix: "3-qator: ",
			check: func(err error) string {
				if !strings.Contains(err.Error(), "4 ta maydon kerak, bor 3") {
					return "parseLine'ning xato matni saqlanmagan"
				}
				return ""
			},
		},
		{
			name:   "uzun_qator",
			input:  "1;olma;12000;40\n2;oʻrik;35000;5\n" + long + "4;choy;25000;8\n",
			prefix: "3-qator: ",
			check: func(err error) string {
				if !errors.Is(err, bufio.ErrTooLong) {
					return "errors.Is(err, bufio.ErrTooLong) false: sc.Err()'ni ham %w bilan oʻrang"
				}
				return ""
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := loadCatalog(strings.NewReader(tc.input))
			if err == nil {
				t.Fatalf("xato kutilgan, nil qaytdi (%d ta mahsulot) (TASK.md, 2-qadam)", len(got))
			}
			msg := err.Error()
			if !strings.HasPrefix(msg, tc.prefix) {
				if r := []rune(msg); len(r) > 80 { // rune chegarasida kesamiz (GB-2.10)
					msg = string(r[:80]) + "…"
				}
				t.Fatalf("xato «%s» bilan boshlanishi kerak, hozir «%s» (TASK.md, 2-qadam)", tc.prefix, msg)
			}
			if why := tc.check(err); why != "" {
				t.Errorf("%s (TASK.md, 2-qadam)", why)
			}
		})
	}
}
