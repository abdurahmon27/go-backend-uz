package store

// Bu fayl — jamoadoshingizning PR'i, tuzatilgan (TASK.md javobi).

import (
	"errors"
	"fmt"
	"strings"
)

// Rename id'li mahsulotning nomini almashtiradi.
// Nom chetidagi boʻsh joylarsiz saqlanadi; boʻsh nom — xato.
func (s *Store) Rename(id int64, name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("nom boʻsh boʻlmasin")
	}
	for i := range s.products {
		p := &s.products[i]
		if p.ID == id {
			p.Name = name
			return nil
		}
	}
	return fmt.Errorf("mahsulot %d: %w", id, ErrNotFound)
}
