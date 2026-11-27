package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
)

// Load katalogni path faylidan (bozor.json) oʻqiydi va har mahsulotni validate'dan oʻtkazadi.
// Fayl hali yoʻq boʻlsa (bozor birinchi marta ishga tushdi), katalog boʻsh.
// Buzuq JSON yoki qoidaga zid mahsulot — xato, fayl nomi bilan %w orqali oʻralgan.
func Load(path string) (*Store, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return &Store{}, nil
	}
	if err != nil {
		return nil, err
	}
	var products []Product
	if err := json.Unmarshal(data, &products); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	seen := make(map[int64]bool)
	for i, p := range products {
		if err := validate(p); err != nil {
			return nil, fmt.Errorf("%s: %d-mahsulot: %w", path, i+1, err)
		}
		if seen[p.ID] {
			return nil, fmt.Errorf("%s: %d-mahsulot: %w: ID %d takrorlangan", path, i+1, ErrBadProduct, p.ID)
		}
		seen[p.ID] = true
	}
	return &Store{products: products}, nil
}

// Save katalogni path fayliga atomar yozadi: avval shu papkadagi vaqtinchalik faylga yozib,
// diskka tushiradi (Sync), keyin os.Rename bilan bozor.json'ning oʻrniga qoʻyadi.
// Yozish oʻrtasida dastur yoki kompyuter toʻxtasa, eski bozor.json butun qoladi.
func (s *Store) Save(path string) (err error) {
	s.mu.Lock()
	products := slices.Clone(s.products) // qulf ostida nusxa; faylga yozish qulfsiz
	s.mu.Unlock()
	if products == nil {
		products = []Product{} // nil slice JSON'da null boʻladi, boʻsh slice esa []
	}
	data, err := json.MarshalIndent(products, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')

	f, err := os.CreateTemp(filepath.Dir(path), ".bozor-*.json")
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			f.Close()
			os.Remove(f.Name()) // xato boʻlsa vaqtinchalik fayl qolmasin
		}
	}()
	if err = f.Chmod(0o644); err != nil { // CreateTemp 0600 bilan ochadi; bozor.json GB-3.02'dagidek 0644
		return err
	}
	if _, err = f.Write(data); err != nil {
		return err
	}
	if err = f.Sync(); err != nil {
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}
	return os.Rename(f.Name(), path)
}
