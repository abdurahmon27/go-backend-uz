package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func slowHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("boshladi")
	time.Sleep(time.Second)
	fmt.Println("tugadi")
}

func main() {
	http.HandleFunc("/slow", slowHandler)
	log.Fatal(http.ListenAndServe("localhost:8085", nil))
}
