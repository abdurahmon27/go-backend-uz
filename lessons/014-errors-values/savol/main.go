// savol: errors.Join(nil, nil) == nil ?
// go run ./savol
package main

import (
	"errors"
	"fmt"
)

func main() {
	err := errors.Join(nil, nil)
	fmt.Println("err == nil:", err == nil)
	fmt.Printf("ichida: %T\n", err)
}
