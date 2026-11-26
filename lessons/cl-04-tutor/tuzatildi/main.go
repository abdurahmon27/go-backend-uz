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

// holat: handler shu kod va matnni HTTP javobga yozadi
func holat(err error) (int, string) {
	if err == nil {
		return 200, "ok"
	}
	if errors.Is(err, ErrTopilmadi) {
		return 404, "topilmadi"
	}
	// zanjirdan *MaydonXato turidagi xatoni qidiradi va topilsa oʻzini qaytaradi (Go 1.26)
	if mx, ok := errors.AsType[*MaydonXato](err); ok {
		return 400, "notoʻgʻri maydon: " + mx.Maydon
	}
	return 500, "ichki xato"
}

func main() {
	for _, narx := range []int{12000, -5} {
		err := tekshir(narx)
		kod, matn := holat(err)
		fmt.Printf("POST /buyurtma narx=%d → %d %s\n", narx, kod, matn)
	}
}
