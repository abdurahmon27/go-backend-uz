package main

import (
	"errors"
	"fmt"
)

// Request — HTTP soʻrovning birinchi qatoridan olingan uch qism (curl -v'dagi "> GET /products/1 HTTP/1.1").
type Request struct {
	Method string // "GET"
	Path   string // "/products/1"
	Proto  string // "HTTP/1.1"
}

// parseRequestLine "GET /products/1 HTTP/1.1" kabi qatorni qismlarga ajratadi (TASK.md, 1-qadam).
func parseRequestLine(line string) (Request, error) {
	return Request{}, errors.New("parseRequestLine hali yozilmagan")
}

// statusClass status kodi qaysi sinfga kirishini aytadi (TASK.md, 2-qadam).
func statusClass(code int) string {
	return ""
}

func main() {
	req, err := parseRequestLine("GET /products/1 HTTP/1.1")
	if err != nil {
		fmt.Println("xato:", err)
		return
	}
	fmt.Println(req.Method, req.Path, req.Proto)
	fmt.Println(404, statusClass(404))
}
