package main

// JSON kirish va chiqish (GB-4.06): bitta xato shakli, qat'iy oʻqish, limit.

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// apiError — hamma xato javobining ichki shakli: {"error": {"code": …, "message": …, "fields": …}}.
type apiError struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

// writeJSON: sarlavha → status → body (GB-4.03 tartibi).
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v) // Encode xatosi: log'ga yozish GB-8.02'da
}

// writeError — barcha xatolar shu funksiyadan chiqadi.
func writeError(w http.ResponseWriter, status int, e apiError) {
	writeJSON(w, status, map[string]apiError{"error": e})
}

// maxBody — body'ning eng katta hajmi, baytda (GB-4.06).
const maxBody = 1024

// decodeJSON body'ni dst'ga oʻqiydi: eng koʻpi maxBody bayt, nomaʼlum kalit yoʻq, bitta JSON qiymat.
// Xato boʻlsa javobni oʻzi yozadi va false qaytaradi: chaqiruvchi darhol return qiladi.
func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, maxBody)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		writeDecodeError(w, err)
		return false
	}
	// Decode faqat birinchi JSON qiymatni oʻqiydi: undan keyin hech narsa boʻlmasligi kerak.
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		if err == nil {
			err = errors.New("body'da bitta JSON boʻlishi kerak")
		}
		writeDecodeError(w, err)
		return false
	}
	return true
}

// writeDecodeError: juda katta body — 413, qolgan hamma oʻqish xatosi — 400.
func writeDecodeError(w http.ResponseWriter, err error) {
	var tooBig *http.MaxBytesError
	if errors.As(err, &tooBig) {
		writeError(w, http.StatusRequestEntityTooLarge,
			apiError{Code: "too_large", Message: fmt.Sprintf("body %d baytdan oshmasin", tooBig.Limit)})
		return
	}
	writeError(w, http.StatusBadRequest, apiError{Code: "bad_json", Message: err.Error()})
}
