// bozor — bozor CLI v0 (GB-3.06): katalog bozor.json faylida turadi.
//
//	bozor add NOM NARX QOLDIQ   mahsulot qoʻshadi (narx soʻmda)
//	bozor list                  katalogni chiqaradi
//	bozor order ID DONA         buyurtma: qoldiqni kamaytiradi
//
// Chiqish kodi: 0 — bajarildi, 1 — xato (mahsulot yoʻq, qoldiq yetmaydi, buzuq fayl),
// 2 — buyruq notoʻgʻri chaqirildi.
package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/abdurahmon27/bozor/internal/store"
)

const dataFile = "bozor.json"

const tiyinPerSom = 100        // 1 soʻm = 100 tiyin (GB-1.03)
const maxPriceSom = 10_000_000 // bitta mahsulot narxining yuqori chegarasi, soʻmda (GB-1.07)

var errUsage = errors.New("foydalanish: bozor add NOM NARX QOLDIQ | list | order ID DONA")

func main() {
	err := run(dataFile, os.Args[1:], os.Stdout)
	if err == nil {
		return
	}
	fmt.Fprintln(os.Stderr, message(err))
	if errors.Is(err, errUsage) {
		os.Exit(2)
	}
	os.Exit(1)
}

// run — butun CLI: katalogni path'dan yuklaydi, args[0] buyrugʻini bajaradi, katalog
// oʻzgargan boʻlsa avval saqlaydi, keyin natija qatorini w'ga yozadi (saqlanmagan narsa
// «bajarildi» deb chiqmaydi). w: main'da os.Stdout, testda strings.Builder.
func run(path string, args []string, w io.Writer) error {
	if len(args) == 0 {
		return errUsage
	}
	st, err := store.Load(path)
	if err != nil {
		return err
	}
	var msg string // oʻzgartiruvchi buyruqning natija qatori
	switch args[0] {
	case "list":
		return list(st, w)
	case "add":
		msg, err = add(st, args[1:])
	case "order":
		msg, err = order(st, args[1:])
	default:
		return errUsage
	}
	if err != nil {
		return err
	}
	if err := st.Save(path); err != nil {
		return err
	}
	_, err = fmt.Fprintln(w, msg)
	return err
}

// add: bozor add NOM NARX QOLDIQ.
func add(st *store.Store, args []string) (string, error) {
	if len(args) != 3 {
		return "", errUsage
	}
	price, err := parsePrice(args[1])
	if err != nil {
		return "", fmt.Errorf("narx: %w", err)
	}
	stock, err := strconv.Atoi(args[2])
	if err != nil {
		return "", fmt.Errorf("qoldiq: %w", err)
	}
	p, err := st.Add(args[0], price, stock)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("qoʻshildi: %d. %s, %d soʻm, %d dona", p.ID, p.Name, p.Price/tiyinPerSom, p.Stock), nil
}

// list: bozor list. %-20s ustunni baytlar emas, belgilar (rune) boʻyicha toʻldiradi (GB-2.10).
func list(st *store.Store, w io.Writer) error {
	n := 0
	for p := range st.All() {
		n++
		_, err := fmt.Fprintf(w, "%2d. %-20s %8d soʻm %4d dona\n",
			p.ID, p.Name, p.Price/tiyinPerSom, p.Stock)
		if err != nil {
			return err
		}
	}
	if n == 0 {
		_, err := fmt.Fprintln(w, "katalog boʻsh: bozor add NOM NARX QOLDIQ")
		return err
	}
	return nil
}

// order: bozor order ID DONA.
func order(st *store.Store, args []string) (string, error) {
	if len(args) != 2 {
		return "", errUsage
	}
	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return "", fmt.Errorf("ID: %w", err)
	}
	quantity, err := strconv.Atoi(args[1])
	if err != nil {
		return "", fmt.Errorf("dona: %w", err)
	}
	p, err := st.Sell(id, quantity)
	if err != nil {
		return "", err
	}
	total := p.Price * int64(quantity)
	return fmt.Sprintf("buyurtma: %s × %d = %d soʻm, qoldiq %d", p.Name, quantity, total/tiyinPerSom, p.Stock), nil
}

// message xatoni (err nil emas) foydalanuvchi oʻqiydigan matnga aylantiradi (GB-2.08).
func message(err error) string {
	var se *store.StockError
	switch {
	case errors.Is(err, errUsage):
		return err.Error()
	case errors.Is(err, store.ErrNotFound):
		return "xato: bunday mahsulot yoʻq"
	case errors.Is(err, store.ErrBadQuantity):
		return "xato: dona soni notoʻgʻri"
	case errors.As(err, &se):
		return fmt.Sprintf("xato: %s: faqat %d dona qoldi", se.Name, se.Stock)
	}
	return "xato: " + err.Error()
}

// parsePrice soʻmdagi narx matnini tiyinga oʻgiradi (GB-1.07).
func parsePrice(s string) (int64, error) {
	som, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	if som <= 0 {
		return 0, fmt.Errorf("narx musbat boʻlishi kerak, keldi %d", som)
	}
	if som > maxPriceSom {
		return 0, fmt.Errorf("narx %d soʻmdan oshmasin, keldi %d", maxPriceSom, som)
	}
	return int64(som) * tiyinPerSom, nil
}
