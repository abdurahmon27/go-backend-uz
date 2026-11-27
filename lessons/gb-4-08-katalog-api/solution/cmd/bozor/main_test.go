package main

// CLI testlari: main emas, run chaqiriladi. run katalog faylining yoʻlini va chiqish uchun
// io.Writer'ni parametr qilib oladi, shuning uchun test har safar oʻz vaqtinchalik papkasida
// (t.TempDir) ishlaydi va chiqishni strings.Builder'da oʻqiydi (GB-2.07).

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// bozor run'ni buyruq qatoridagidek chaqiradi: bozor(t, path, "order 1 3").
func bozor(t *testing.T, path, line string) (string, error) {
	t.Helper()
	var out strings.Builder
	err := run(path, strings.Fields(line), &out)
	return out.String(), err
}

// yangi videodagi toʻrt mahsulotli bozor.json'ni yangi vaqtinchalik papkada yasaydi.
func yangi(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "bozor.json")
	for _, line := range []string{"add olma 12000 40", "add non 4000 3", "add choy 25000 8", "add oʻrik 35000 2"} {
		if _, err := bozor(t, path, line); err != nil {
			t.Fatalf("%s: %v", line, err)
		}
	}
	return path
}

func TestAddList(t *testing.T) {
	path := yangi(t)
	out, err := bozor(t, path, "list")
	if err != nil {
		t.Fatal(err)
	}
	want := "" +
		" 1. olma                    12000 soʻm   40 dona\n" +
		" 2. non                      4000 soʻm    3 dona\n" +
		" 3. choy                    25000 soʻm    8 dona\n" +
		" 4. oʻrik                   35000 soʻm    2 dona\n"
	if out != want {
		t.Fatalf("list chiqishi:\n%s\nkutilgan:\n%s", out, want)
	}
}

// order qoldiqni kamaytiradi va bu bozor.json'ga yoziladi: keyingi list uni koʻradi.
func TestOrderSaved(t *testing.T) {
	path := yangi(t)
	out, err := bozor(t, path, "order 1 3")
	if err != nil {
		t.Fatal(err)
	}
	if want := "buyurtma: olma × 3 = 36000 soʻm, qoldiq 37\n"; out != want {
		t.Fatalf("order chiqishi %q, kutilgan %q", out, want)
	}
	out, _ = bozor(t, path, "list")
	if !strings.Contains(out, "olma                    12000 soʻm   37 dona") {
		t.Fatalf("order'dan keyin list'da olma qoldigʻi 37 emas:\n%s", out)
	}
}

func TestErrors(t *testing.T) {
	path := yangi(t)
	tests := []struct {
		line string
		want string // message(err) — foydalanuvchi koʻradigan matn; "" — foydalanish xatosi (errUsage)
	}{
		{"", ""},
		{"sell 1 1", ""},
		{"order 1", ""},
		{"order 9 1", "xato: bunday mahsulot yoʻq"},
		{"order 4 5", "xato: oʻrik: faqat 2 dona qoldi"},
		{"order 1 0", "xato: dona soni notoʻgʻri"},
		{"order bir 1", `xato: ID: strconv.ParseInt: parsing "bir": invalid syntax`},
		{"add anor 0 1", "xato: narx: narx musbat boʻlishi kerak, keldi 0"},
	}
	for _, tc := range tests {
		t.Run(tc.line, func(t *testing.T) {
			_, err := bozor(t, path, tc.line)
			if err == nil {
				t.Fatalf("%q xatosiz bajarildi", tc.line)
			}
			if tc.want == "" {
				if !errors.Is(err, errUsage) {
					t.Fatalf("%q: xato %v, kutilgan errUsage (chiqish kodi 2)", tc.line, err)
				}
				return
			}
			if got := message(err); got != tc.want {
				t.Fatalf("%q: xabar %q, kutilgan %q", tc.line, got, tc.want)
			}
		})
	}
	// Rad etilgan buyurtmalar katalogni oʻzgartirmadi.
	out, _ := bozor(t, path, "list")
	if !strings.Contains(out, "oʻrik                   35000 soʻm    2 dona") {
		t.Fatalf("xatolardan keyin katalog oʻzgardi:\n%s", out)
	}
}

// Buzuq bozor.json: dastur xato bilan chiqadi va faylga tegmaydi.
func TestBrokenFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bozor.json")
	writeFile(t, path, `[{"id":1,"name":"ol`)
	_, err := bozor(t, path, "add non 4000 3")
	if err == nil || !strings.Contains(err.Error(), "unexpected end of JSON input") {
		t.Fatalf("buzuq faylda xato %v, kutilgan «unexpected end of JSON input»", err)
	}
	if errors.Is(err, errUsage) {
		t.Fatal("buzuq fayl — foydalanish xatosi emas")
	}
	if data, _ := os.ReadFile(path); string(data) != `[{"id":1,"name":"ol` {
		t.Fatalf("buzuq fayl oʻzgardi: %q", data)
	}
}

// Saqlab boʻlmasa, «buyurtma: …» chiqmaydi: foydalanuvchi saqlanmagan natijani koʻrmasin.
func TestNoOutputWithoutSave(t *testing.T) {
	path := yangi(t)
	dir := filepath.Dir(path)
	if err := os.Chmod(dir, 0o500); err != nil { // papka faqat oʻqish uchun: Save yangi fayl ocholmaydi
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chmod(dir, 0o700) })
	out, err := bozor(t, path, "order 1 3")
	if err == nil {
		t.Skip("papka yozishdan himoyalanmadi (root?): bu tekshiruvni oʻtkazib boʻlmaydi")
	}
	if out != "" {
		t.Fatalf("Save xato berdi (%v), lekin chiqishda %q bor", err, out)
	}
}

func writeFile(t *testing.T, path, data string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
}
