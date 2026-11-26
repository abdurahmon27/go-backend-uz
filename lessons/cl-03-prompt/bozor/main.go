package main

import (
	"errors"
	"fmt"
)

var ErrTopilmadi = errors.New("buyurtma topilmadi")

var buyurtmalar = map[string]string{"1": "olma, 2 kg"}

func topish(id string) (string, error) {
	b, ok := buyurtmalar[id]
	if !ok {
		return "", fmt.Errorf("topish %s: %w", id, ErrTopilmadi)
	}
	return b, nil
}

// holat: handler shu kodni HTTP javobga yozadi
func holat(err error) int {
	if err == nil {
		return 200
	}
	if err == ErrTopilmadi {
		return 404
	}
	return 500
}

func main() {
	for _, id := range []string{"1", "42"} {
		_, err := topish(id)
		fmt.Printf("GET /buyurtma/%s → %d %v\n", id, holat(err), err)
	}
}
