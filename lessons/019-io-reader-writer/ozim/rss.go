package ozim

import (
	"bytes"
	"os"
	"strconv"
)

// RSS — shu jarayon (dasturning oʻzi) egallagan eng katta operativ xotira (peak RSS), MiB. Faqat Linux:
// /proc/self/status dagi VmHWM. getrusage'ning Maxrss'i bu yerga yaramaydi: `go run` va shell
// dasturni fork/vfork + exec bilan ishga tushiradi, Maxrss esa ota jarayonning xotirasini ham "meros" oladi.
func RSS() float64 {
	b, err := os.ReadFile("/proc/self/status")
	if err != nil {
		return -1
	}
	_, s, _ := bytes.Cut(b, []byte("VmHWM:"))
	s, _, _ = bytes.Cut(s, []byte("kB"))
	kb, err := strconv.ParseInt(string(bytes.TrimSpace(s)), 10, 64)
	if err != nil {
		return -1
	}
	return float64(kb) / 1024
}
