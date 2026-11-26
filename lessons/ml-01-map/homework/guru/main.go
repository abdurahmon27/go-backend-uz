// Guru vazifasi javobi: toʻgʻri chiziqni (least squares) noldan, faqat standart kutubxona bilan.
//
//	uv run homework/nuqtalar.py && go run ./homework/guru homework/nuqtalar.csv
//
// Formula (CS229 konspekti, 1.2 «The normal equations» ning bir oʻzgaruvchili holati):
//
//	b = Σ(x-x̄)(y-ȳ) / Σ(x-x̄)²,   a = ȳ - b·x̄
package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "ishlatish: go run ./homework/guru nuqtalar.csv")
		os.Exit(2)
	}
	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer f.Close()
	qatorlar, err := csv.NewReader(f).ReadAll()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	var xs, ys []float64
	for _, q := range qatorlar {
		x, err1 := strconv.ParseFloat(q[0], 64)
		y, err2 := strconv.ParseFloat(q[1], 64)
		if err1 != nil || err2 != nil {
			fmt.Fprintln(os.Stderr, "notoʻgʻri qator:", q)
			os.Exit(1)
		}
		xs, ys = append(xs, x), append(ys, y)
	}
	a, b := chiziq(xs, ys)
	fmt.Printf("go:    ball = %.12f + %.12f × soat\n", a, b)
}

func chiziq(xs, ys []float64) (a, b float64) {
	n := float64(len(xs))
	var mx, my float64
	for i := range xs {
		mx += xs[i]
		my += ys[i]
	}
	mx, my = mx/n, my/n
	var sxy, sxx float64
	for i := range xs {
		sxy += (xs[i] - mx) * (ys[i] - my)
		sxx += (xs[i] - mx) * (xs[i] - mx)
	}
	b = sxy / sxx
	return my - b*mx, b
}
