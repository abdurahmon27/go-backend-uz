package homework

import "testing"

func TestTelegramToken(t *testing.T) {
	for s, want := range map[string]bool{
		"BOT_TOKEN=" + "123456:" + "ABC-DEF1234ghIkl-zyx57W2v1u123ew11":     true, // hujjatdagi namuna
		"BOT_TOKEN=" + "7012345678:" + "SOXTA_kalit-0123456789abcdefABCDEF": true,
		"vaqt 10:14:02, port :8080":                                         false,
		"id=2026112510140712:ok":                                            false,
	} {
		if got := TelegramToken.MatchString(s); got != want {
			t.Errorf("%q: %v, kutilgan %v", s, got, want)
		}
	}
}
