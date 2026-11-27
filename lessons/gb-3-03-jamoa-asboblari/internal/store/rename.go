package store

// Bu fayl — jamoadoshingizning PR'i (TASK.md). U kompilyatsiya boʻladi va eski test oʻtadi,
// lekin go vet va staticcheck unda xato topadi.

import (
	"errors"
	"fmt"
	"strings"
)

// Rename id'li mahsulotning nomini almashtiradi.
// Nom chetidagi boʻsh joylarsiz saqlanadi; boʻsh nom — xato.
func (s *Store) Rename(id int64, name string) error {
	strings.TrimSpace(name)
	if name == "" {
		return errors.New("nom boʻsh boʻlmasin")
	}
	for i := range s.products {
		p := &s.products[i]
		if p.ID == id {
			p.Name = p.Name
			return nil
		}
	}
	return fmt.Errorf("mahsulot %d: %w", id, ErrNotFound)
}
