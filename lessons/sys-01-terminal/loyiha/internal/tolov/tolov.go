// Package tolov: toʻlov (namuna).
package tolov

import "errors"

// ErrRad toʻlov rad etilganda qaytadi.
var ErrRad = errors.New("toʻlov rad etildi")

// Tolash summani yechadi.
func Tolash(summa int64) error {
	if summa <= 0 {
		return ErrRad
	}
	// TODO: idempotency key
	return nil
}
