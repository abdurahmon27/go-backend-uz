package qiyin

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"
)

func TestChekli(t *testing.T) {
	var buf bytes.Buffer
	n, err := io.Copy(&Chekli{W: &buf, Qoldi: 5}, strings.NewReader("salom dunyo"))
	if n != 5 || !errors.Is(err, ErrToldi) || buf.String() != "salom" {
		t.Fatalf("n=%d err=%v buf=%q", n, err, buf.String())
	}
}
