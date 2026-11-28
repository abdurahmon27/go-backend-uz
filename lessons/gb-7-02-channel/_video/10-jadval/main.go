package main

// Chuqurroq: channel jadvali. Har holatni alohida ishga tushiring:
//   ./_video/run.sh 10-jadval go run . nil-send | nil-recv | nil-close | closed-close | closed-send | closed-recv
import (
	"fmt"
	"os"
)

func main() {
	var nilCh chan int64 // nol qiymati: nil channel
	closed := make(chan int64)
	close(closed)
	switch os.Args[1] {
	case "nil-send":
		nilCh <- 1
	case "nil-recv":
		<-nilCh
	case "nil-close":
		close(nilCh)
	case "closed-close":
		close(closed)
	case "closed-send":
		closed <- 1
	case "closed-recv":
		v, ok := <-closed
		fmt.Println(v, ok)
	}
}
