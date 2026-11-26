package ozim

import (
	"bytes"
	"io"
	"testing"
	"testing/iotest"
)

// Generator Read shartnomasiga rioya qiladimi? iotest.TestReader kichik buferlar, Read(0) va EOF'ni tekshiradi.
func TestGenerator(t *testing.T) {
	want, err := io.ReadAll(NewGenerator(100_000))
	if err != nil {
		t.Fatal(err)
	}
	if err := iotest.TestReader(NewGenerator(100_000), want); err != nil {
		t.Fatal(err)
	}
	if !bytes.HasPrefix(want, []byte("1,8919,1\n2,16838,1\n")) || want[len(want)-1] != '\n' {
		t.Fatalf("kutilmagan boshlanish yoki oxir: %q…%q", want[:20], want[len(want)-5:])
	}
}

func TestHisoblagich(t *testing.T) {
	var h Hisoblagich
	n, err := io.Copy(&h, iotest.OneByteReader(bytes.NewReader([]byte("a\nb\nc\n"))))
	if err != nil || n != 6 || h.Bayt != 6 || h.Qator != 3 {
		t.Fatalf("n=%d err=%v %+v", n, err, h)
	}
}
