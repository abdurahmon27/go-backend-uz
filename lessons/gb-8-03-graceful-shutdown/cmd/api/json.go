package main

// JSON kirish va chiqish (GB-4.06): qat'iy oʻqish va limit. Xato javobi problem.go'da (GB-6.02).

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
	json.NewEncoder(w).Encode(v) // Encode xatosi: status allaqachon ketgan, odatda mijoz uzilgan (GB-8.02 README)
}

// maxBody — body'ning eng katta hajmi, baytda (GB-4.06).
const maxBody = 4096 // GB-5.10: savatda 20 qatorgacha

// decodeJSON body'ni dst'ga oʻqiydi: eng koʻpi maxBody bayt, nomaʼlum kalit yoʻq, bitta JSON qiymat.
// Xato boʻlsa javobni oʻzi yozadi va false qaytaradi: chaqiruvchi darhol return qiladi.
// Decode xatosi errBadJSON bilan birga oʻraladi (ikkita %w, GB-6.02): problemFor ikkalasini ham koʻradi.
func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, maxBody)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		writeError(w, r, fmt.Errorf("%w: %w", errBadJSON, err))
		return false
	}
	// Decode faqat birinchi JSON qiymatni oʻqiydi: undan keyin hech narsa boʻlmasligi kerak.
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		if err == nil {
			err = errors.New("body'da bitta JSON boʻlishi kerak")
		}
		writeError(w, r, fmt.Errorf("%w: %w", errBadJSON, err))
		return false
	}
	return true
}
