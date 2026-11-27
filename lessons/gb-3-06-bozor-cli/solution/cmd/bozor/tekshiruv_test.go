package main

// Bu fayl vazifaning 3-qadamini (TASK.md) tekshiradi: bozor restock ID DONA buyrugʻi.
// go test ./... uni main.go bilan birga ishga tushiradi. Uni oʻzgartirmang.

import (
	"errors"
	"strings"
	"testing"
)

// 3-qadam: restock qoldiqni oshiradi, natijani chiqaradi va bozor.json'ga saqlaydi.
func TestVazifaRestockBuyrugi(t *testing.T) {
	path := yangi(t)
	out, err := bozor(t, path, "restock 2 10")
	if err != nil {
		t.Fatalf("bozor restock 2 10: xato «%s», kutilgan chiqish «toʻldirildi: non, qoldiq 13» (TASK.md, 3-qadam)", message(err))
	}
	if want := "toʻldirildi: non, qoldiq 13\n"; out != want {
		t.Fatalf("bozor restock 2 10 chiqishi %q, kutilgan %q (TASK.md, 3-qadam; ʻ — U+02BB, oddiy ' emas: TASK.md'dan nusxalang)", out, want)
	}
	out, _ = bozor(t, path, "list")
	if !strings.Contains(out, "non                      4000 soʻm   13 dona") {
		t.Fatalf("restock'dan keyin list'da non qoldigʻi 13 emas: katalogni saqlang (TASK.md, 3-qadam)\n%s", out)
	}
}

// 3-qadam: xatolar order'dagidek: message() matni va chiqish kodi.
func TestVazifaRestockXatolar(t *testing.T) {
	path := yangi(t)
	tests := []struct {
		line string
		want string // message(err); "" — foydalanish xatosi (errUsage, chiqish kodi 2)
	}{
		{"restock 9 5", "xato: bunday mahsulot yoʻq"},
		{"restock 2 0", "xato: dona soni notoʻgʻri"},
		{"restock 2", ""},
		{"restock 2 10 5", ""},
		{"restock non 10", "xato: …"}, // ID son emas: istalgan xato, foydalanish xatosi emas
	}
	for _, tc := range tests {
		t.Run(tc.line, func(t *testing.T) {
			_, err := bozor(t, path, tc.line)
			if err == nil {
				t.Fatalf("bozor %s xatosiz bajarildi (TASK.md, 3-qadam)", tc.line)
			}
			if tc.want == "" {
				if !errors.Is(err, errUsage) {
					t.Fatalf("bozor %s: xato «%s», kutilgan foydalanish xatosi errUsage (TASK.md, 3-qadam)", tc.line, message(err))
				}
				return
			}
			if tc.want == "xato: …" {
				if errors.Is(err, errUsage) || !strings.HasPrefix(message(err), "xato: ") {
					t.Fatalf("bozor %s: xabar «%s», kutilgan «xato: …» (ID son emas; TASK.md, 3-qadam)", tc.line, message(err))
				}
				return
			}
			if got := message(err); got != tc.want {
				t.Fatalf("bozor %s: xabar «%s», kutilgan «%s» (TASK.md, 3-qadam)", tc.line, got, tc.want)
			}
		})
	}
	if !strings.Contains(errUsage.Error(), "restock ID DONA") {
		t.Fatalf("foydalanish matnida restock yoʻq: %q; errUsage'ga «| restock ID DONA» qoʻshing (TASK.md, 3-qadam)", errUsage.Error())
	}
}
