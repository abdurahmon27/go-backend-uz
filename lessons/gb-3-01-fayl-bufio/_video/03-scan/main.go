// Video, 03-scan: bufio.Scanner fayldan qatorlarni beradi.
package main

import (
	"bufio"
	"fmt"
	"os"
)

func run() error {
	f, err := os.Open("katalog.txt")
	if err != nil {
		return err
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		fmt.Printf("%q\n", sc.Text())
	}
	if err := sc.Err(); err != nil {
		return err
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "xato:", err)
		os.Exit(1)
	}
}
