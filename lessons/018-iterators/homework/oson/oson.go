// Oson: Oraliq(a, b) a dan b gacha (b ham kiradi) sonlarni beradi.
package oson

import "iter"

func Oraliq(a, b int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := a; i <= b; i++ {
			if !yield(i) {
				return
			}
		}
	}
}
