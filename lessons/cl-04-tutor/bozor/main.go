package main

import (
	"errors"
	"fmt"
)

var ErrTopilmadi = errors.New("buyurtma topilmadi")

// MaydonXato: buyurtmaning bitta maydoni notoʻgʻri
type MaydonXato struct {
	Maydon string
	Sabab  string
}

func (e *MaydonXato) Error() string { return e.Maydon + ": " + e.Sabab }

func tekshir(narx int) error {
	if narx <= 0 {
		return fmt.Errorf("tekshir: %w", &MaydonXato{Maydon: "narx", Sabab: "musbat boʻlishi kerak"})
	}
	return nil
}

// holat: handler shu kodni HTTP javobga yozadi
func holat(err error) int {
	if err == nil {
		return 200
	}
	if errors.Is(err, ErrTopilmadi) {
		return 404
	}
	// mening urinishim: oʻtgan safar errors.Is yordam bergan edi
	if errors.Is(err, &MaydonXato{}) {
		return 400
	}
	return 500
}

func main() {
	for _, narx := range []int{12000, -5} {
		err := tekshir(narx)
		fmt.Printf("POST /buyurtma narx=%d → %d %v\n", narx, holat(err), err)
	}
}
