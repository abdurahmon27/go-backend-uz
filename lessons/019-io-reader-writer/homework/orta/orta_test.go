package orta

import (
	"io"
	"strings"
	"testing"
	"testing/iotest"
)

func TestCRsiz(t *testing.T) {
	got, err := io.ReadAll(CRsiz{iotest.OneByteReader(strings.NewReader("7000 soʻm\r\n42000 soʻm\r\n"))})
	if err != nil || string(got) != "7000 soʻm\n42000 soʻm\n" {
		t.Fatalf("%q %v", got, err)
	}
	// iotest.TestReader Read shartnomasini tekshiradi (n, EOF, kichik buferlar).
	if err := iotest.TestReader(CRsiz{strings.NewReader("a\r\nb\r\n")}, []byte("a\nb\n")); err != nil {
		t.Fatal(err)
	}
}
