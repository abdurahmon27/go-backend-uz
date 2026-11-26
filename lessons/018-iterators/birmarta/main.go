// Bir martalik iterator: strings.SplitSeq ikkinchi marta boshidan boshlamaydi.
package main

import (
	"fmt"
	"strings"
)

func main() {
	qismlar := strings.SplitSeq("Aziz,Bobur,Dilnoz", ",")
	for q := range qismlar {
		fmt.Println("1-marta:", q)
	}
	for q := range qismlar {
		fmt.Println("2-marta:", q)
	}
}
