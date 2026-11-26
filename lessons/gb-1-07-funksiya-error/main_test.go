package main

// Bu fayl vazifani tekshiradi: `go test` uni ishga tushiradi.
// Test yozishni GB-1.08 darsida ochamiz; hozir uni oʻzgartirmang, faqat main.go'ni oʻzgartiring.

import (
	"bytes"
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// TestMain: agar BOZOR_ARGS berilgan boʻlsa, test fayli main()'ni shu argumentlar bilan
// ishga tushiradi (xuddi ./bozor add … kabi). Aks holda oddiy testlar ishlaydi.
func TestMain(m *testing.M) {
	if a, ok := os.LookupEnv("BOZOR_ARGS"); ok {
		os.Args = append([]string{"bozor"}, strings.Split(a, "\x1f")...)
		main()
		os.Exit(0)
	}
	os.Exit(m.Run())
}

// bozor dasturni berilgan argumentlar bilan ishga tushiradi: stdout, stderr va chiqish kodi.
func bozor(t *testing.T, args ...string) (stdout, stderr string, code int) {
	t.Helper()
	cmd := exec.Command(os.Args[0], "-test.run=^$")
	cmd.Env = append(os.Environ(), "BOZOR_ARGS="+strings.Join(args, "\x1f"))
	var o, e bytes.Buffer
	cmd.Stdout, cmd.Stderr = &o, &e
	err := cmd.Run()
	var ee *exec.ExitError
	if errors.As(err, &ee) {
		code = ee.ExitCode()
	} else if err != nil {
		t.Fatal(err)
	}
	return strings.TrimRight(o.String(), "\n"), strings.TrimRight(e.String(), "\n"), code
}

func buyruq(args []string) string { return "./bozor " + strings.Join(args, " ") }

func TestParseSoldBor(t *testing.T) {
	f, err := parser.ParseFile(token.NewFileSet(), "main.go", nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	for _, d := range f.Decls {
		fn, ok := d.(*ast.FuncDecl)
		if !ok || fn.Name.Name != "parseSold" {
			continue
		}
		p, r := fn.Type.Params.List, fn.Type.Results
		if len(p) == 1 && len(p[0].Names) == 1 && types(p[0].Type) == "string" &&
			r != nil && len(r.List) == 2 && types(r.List[0].Type) == "int" && types(r.List[1].Type) == "error" {
			return
		}
		t.Fatalf("parseSold bor, lekin imzosi boshqa: kerak «func parseSold(s string) (int, error)» (TASK.md)")
	}
	t.Fatalf("main.go'da «func parseSold(s string) (int, error)» funksiyasi yoʻq: miqdorni shu funksiya tekshirsin (TASK.md)")
}

func types(e ast.Expr) string {
	if id, ok := e.(*ast.Ident); ok {
		return id.Name
	}
	return "?"
}

func TestAddMiqdor(t *testing.T) {
	args := []string{"add", "olma", "12000", "5"}
	out, errOut, code := bozor(t, args...)
	want := strings.Join([]string{
		"bozor ishga tushdi",
		"versiya: 0.1.0",
		" 0 false",
		"olma 1200000 true 40 1.5",
		"olma 35",
		"olma, 5 dona       48000.00",
		"QQS                 5760.00",
		"yetkazish          15000.00",
		"jami               68760.00",
		"aksiya: 1 dona bepul",
	}, "\n")
	if code != 0 || out != want {
		t.Fatalf("«%s» 5 dona uchun chek chiqarishi va 0 kod bilan tugashi kerak.\n  kerak:\n%s\n  hozir (kod %d, xato oqimi «%s»):\n%s\n(TASK.md)",
			buyruq(args), want, code, errOut, out)
	}
}

// start: dastur har doim avval shu uch qatorni chiqaradi (ishga tushish qatorlari va GB-1.02'dagi zero value qatori).
const start = "bozor ishga tushdi\nversiya: 0.1.0\n 0 false"

func TestRadEtadi(t *testing.T) {
	cases := []struct {
		args []string
		has  string // xato oqimida boʻlishi kerak boʻlgan matn
	}{
		{[]string{"add", "olma", "12000", "besh"}, `xato: strconv.Atoi: parsing "besh": invalid syntax`},
		{[]string{"add", "olma", "12000", "0"}, "xato: miqdor kamida 1 boʻlishi kerak, keldi 0"},
		{[]string{"add", "olma", "12000", "-2"}, "xato: miqdor kamida 1 boʻlishi kerak, keldi -2"},
		{[]string{"add", "olma", "12k", "5"}, `xato: strconv.Atoi: parsing "12k": invalid syntax`},
		{[]string{"add", "olma", "-5", "5"}, "xato: narx musbat boʻlishi kerak, keldi -5"},
		{[]string{"add", "olma", "12000"}, "foydalanish: bozor add NOM NARX MIQDOR"},
	}
	for _, c := range cases {
		out, errOut, code := bozor(t, c.args...)
		if code != 1 || errOut != c.has || !strings.HasPrefix(out, start) || strings.Contains(out, "jami") {
			t.Errorf("«%s» rad etilishi kerak: xato oqimida «%s», kod 1, oddiy chiqishda chek yoʻq;\n  hozir: xato oqimi «%s», chiqish «%s», kod %d",
				buyruq(c.args), c.has, errOut, out, code)
		}
	}
}
