package ozim

import "bytes"

// Hisoblagich — oʻz io.Writer'imiz: kelgan baytlar va qatorlarni sanaydi, chelakni saqlamaydi.
type Hisoblagich struct{ Bayt, Qator int64 } // ikkita sanoq

func (h *Hisoblagich) Write(p []byte) (int, error) {
	h.Bayt += int64(len(p))                        // kelgan baytlar soni
	h.Qator += int64(bytes.Count(p, []byte{'\n'})) // qator oxirlari
	return len(p), nil                             // hammasini oldim, muammo yoʻq
}
