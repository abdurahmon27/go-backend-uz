package main

// JSON kirish va chiqish (GB-4.06): qat'iy oʻqish va limit. GB-6.02: xato javobi endi problem.go'da.

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// writeJSON: sarlavha → status → body (GB-4.03 tartibi).
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v) // Encode xatosi: log'ga yozish GB-8.02'da
}

// maxBody — body'ning eng katta hajmi, baytda (GB-4.06).
const maxBody = 1024

// decodeJSON body'ni dst'ga oʻqiydi: eng koʻpi maxBody bayt, nomaʼlum kalit yoʻq, bitta JSON qiymat.
// Xato boʻlsa javobni oʻzi yozadi va false qaytaradi: chaqiruvchi darhol return qiladi.
// GB-6.02: Decode xatosi errBadJSON bilan birga oʻraladi (ikkita %w): problemFor ikkalasini ham koʻradi,
// asl matn esa faqat log'ga boradi.
func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, maxBody)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		writeError(w, fmt.Errorf("%w: %w", errBadJSON, err))
		return false
	}
	// Decode faqat birinchi JSON qiymatni oʻqiydi: undan keyin hech narsa boʻlmasligi kerak.
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		if err == nil {
			err = errors.New("body'da bitta JSON boʻlishi kerak")
		}
		writeError(w, fmt.Errorf("%w: %w", errBadJSON, err))
		return false
	}
	return true
}
