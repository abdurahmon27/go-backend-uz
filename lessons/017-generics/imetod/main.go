// imetod: ATAYLAB kompilyatsiya boʻlmaydi. Interface metodi type parameter e'lon qila olmaydi.
// go build ./imetod
package main

type Aylantiruvchi interface {
	Aylantir[U any](f func(int64) U) []U
}

func main() {}
