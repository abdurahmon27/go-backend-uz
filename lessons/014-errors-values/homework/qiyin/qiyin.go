// Qiyin: oʻz xato turimiz boshqa xatoni oʻraydi.
// Unwrap metodi boʻlsa, errors.Is va errors.AsType uning ichiga ham qaraydi.
package qiyin

import (
	"errors"
	"fmt"
)

var ErrBosh = errors.New("qiymat boʻsh")

// MaydonXato qaysi maydonda xato boʻlganini aytadi va asl xatoni saqlaydi.
type MaydonXato struct {
	Maydon string
	Err    error
}

func (e *MaydonXato) Error() string { return e.Maydon + ": " + e.Err.Error() }
func (e *MaydonXato) Unwrap() error { return e.Err }

func Tekshir(manzil string) error {
	if manzil == "" {
		return fmt.Errorf("buyurtma: %w", &MaydonXato{Maydon: "manzil", Err: ErrBosh})
	}
	return nil
}
