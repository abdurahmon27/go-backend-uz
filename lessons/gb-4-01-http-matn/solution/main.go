package main

import (
	"errors"
	"fmt"
	"strings"
)

// Request — HTTP soʻrovning birinchi qatoridan olingan uch qism (curl -v'dagi "> GET /products/1 HTTP/1.1").
type Request struct {
	Method string // "GET"
	Path   string // "/products/1"
	Proto  string // "HTTP/1.1"
}

// parseRequestLine "GET /products/1 HTTP/1.1" kabi qatorni qismlarga ajratadi (TASK.md, 1-qadam).
func parseRequestLine(line string) (Request, error) {
	line = strings.TrimRight(line, "\r\n") // HTTP qatori \r\n bilan tugaydi
	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return Request{}, fmt.Errorf("soʻrov qatori uch qismdan iborat boʻlishi kerak, hozir %d: %q", len(parts), line)
	}
	method, path, proto := parts[0], parts[1], parts[2]
	if method == "" {
		return Request{}, errors.New("method boʻsh")
	}
	if !strings.HasPrefix(path, "/") {
		return Request{}, fmt.Errorf("yoʻl / bilan boshlanishi kerak: %q", path)
	}
	if !strings.HasPrefix(proto, "HTTP/") {
		return Request{}, fmt.Errorf("versiya HTTP/ bilan boshlanishi kerak: %q", proto)
	}
	return Request{Method: method, Path: path, Proto: proto}, nil
}

// statusClass status kodi qaysi sinfga kirishini aytadi (TASK.md, 2-qadam).
func statusClass(code int) string {
	switch {
	case code >= 100 && code < 200:
		return "axborot"
	case code >= 200 && code < 300:
		return "muvaffaqiyat"
	case code >= 300 && code < 400:
		return "yoʻnaltirish"
	case code >= 400 && code < 500:
		return "mijoz xatosi"
	case code >= 500 && code < 600:
		return "server xatosi"
	}
	return "nomaʼlum"
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
