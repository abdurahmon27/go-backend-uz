package main

import (
	"fmt"
	"net/http"
	"time"
)

// sekin: har bir soʻrov 1 soniya, masalan toʻlov.
func sekin(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Second)
	fmt.Fprintln(w, "javob: Salom!")
}
