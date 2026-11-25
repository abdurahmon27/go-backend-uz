package orta

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaqla(t *testing.T) {
	nom := filepath.Join(t.TempDir(), "a.txt")
	if err := Saqla(nom, []byte("salom")); err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(nom)
	if string(b) != "salom" {
		t.Fatalf("fayl: %q", b)
	}
	if err := saqlaXato(nom, []byte("salom")); err != nil {
		t.Fatal(err)
	}
}

// /dev/full: yozish "no space left on device" bilan qaytadi.
func TestToliqDisk(t *testing.T) {
	if _, err := os.Stat("/dev/full"); err != nil {
		t.Skip("/dev/full yo'q")
	}
	if err := Saqla("/dev/full", []byte("salom")); err == nil {
		t.Fatal("xato kutilgan edi")
	} else {
		t.Log("Saqla:", err)
	}
}
