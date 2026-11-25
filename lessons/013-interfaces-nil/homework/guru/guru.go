// Uy vazifasi (guru): interface ichida typed nil bormi? go test -v ./homework/guru
package guru

import "reflect"

// bosh: a == nil yoki ichida nil pointer/map/slice/chan/func (typed nil) boʻlsa true.
func bosh(a any) bool {
	if a == nil {
		return true
	}
	v := reflect.ValueOf(a)
	switch v.Kind() {
	case reflect.Pointer, reflect.Map, reflect.Slice, reflect.Chan, reflect.Func:
		return v.IsNil()
	}
	return false
}
