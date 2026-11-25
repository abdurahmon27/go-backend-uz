// O'rta: Close xatosini yo'qotmaydigan saqla.
package orta

import (
	"errors"
	"os"
)

// saqlaXato: defer f.Close() Close qaytargan xatoni tashlab yuboradi.
func saqlaXato(nom string, data []byte) error {
	f, err := os.Create(nom)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	return err
}

// Saqla: nomlangan natija + defer closure, Close xatosi ham qaytadi.
func Saqla(nom string, data []byte) (err error) {
	f, err := os.Create(nom)
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, f.Close())
	}()
	_, err = f.Write(data)
	return err
}
