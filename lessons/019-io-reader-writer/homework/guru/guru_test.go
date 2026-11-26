package guru

import (
	"io"
	"strings"
	"testing"
	"testing/iotest"
	"unicode/utf8"
)

func TestRuneSanagich(t *testing.T) {
	matn := strings.Repeat("7000 soʻm, gʻisht, 😀 ", 1000)
	var s RuneSanagich
	// OneByteReader har safar bitta bayt beradi: har bir koʻp baytli rune albatta kesiladi.
	if _, err := io.Copy(&s, iotest.OneByteReader(strings.NewReader(matn))); err != nil {
		t.Fatal(err)
	}
	if want := int64(utf8.RuneCountInString(matn)); s.Rune != want {
		t.Fatalf("got %d, want %d", s.Rune, want)
	}
}

func TestBuzuqOxir(t *testing.T) {
	var s RuneSanagich
	s.Write([]byte("so\xca"))
	s.Close()
	if want := int64(utf8.RuneCount([]byte("so\xca"))); s.Rune != want {
		t.Fatalf("got %d, want %d", s.Rune, want)
	}
}
