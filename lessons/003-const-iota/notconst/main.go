package main

import (
	"fmt"
	"math"
	"time"
)

func main() {
	const maxSize = 10 * 1024 * 1024 // faqat konstantalar va operatorlar: kompilyator hisoblaydi
	const root2 = math.Sqrt(2)       // ataylab xato: funksiya chaqiruvi — natija ishlash paytida
	const start = time.Now()         // ataylab xato: time.Now() faqat ishlash paytida maʼlum
	fmt.Println(maxSize, root2, start)
}
