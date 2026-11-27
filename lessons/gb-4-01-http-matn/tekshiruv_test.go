package main

// Bu fayl vazifani (TASK.md) tekshiradi: go test uni main.go bilan birga ishga tushiradi.
// Uni oʻzgartirmang.

import "testing"

// 1-qadam: parseRequestLine uch qismni ajratadi, notoʻgʻri qatorda xato qaytaradi.
func TestVazifaParseRequestLine(t *testing.T) {
	tests := []struct {
		name    string
		line    string
		want    Request
		wantErr bool
	}{
		{"get_mahsulot", "GET /products/1 HTTP/1.1", Request{"GET", "/products/1", "HTTP/1.1"}, false},
		{"post_royxat", "POST /products HTTP/1.1", Request{"POST", "/products", "HTTP/1.1"}, false},
		{"yol_ildiz", "GET / HTTP/1.1", Request{"GET", "/", "HTTP/1.1"}, false},
		{"yol_savol_bilan", "GET /products?limit=2 HTTP/1.1", Request{"GET", "/products?limit=2", "HTTP/1.1"}, false},
		{"qator_oxiri_crlf", "GET /products/1 HTTP/1.1\r\n", Request{Method: "GET", Path: "/products/1", Proto: "HTTP/1.1"}, false},
		{"qator_oxiri_cr", "POST /products HTTP/1.1\r", Request{Method: "POST", Path: "/products", Proto: "HTTP/1.1"}, false},
		{"bosh_qator", "", Request{}, true},
		{"ikki_qism", "GET /products/1", Request{}, true},
		{"tort_qism", "GET /products/1 HTTP/1.1 extra", Request{}, true},
		{"qoshaloq_probel", "GET  /products/1 HTTP/1.1", Request{}, true},
		{"yol_slashsiz", "GET products/1 HTTP/1.1", Request{}, true},
		{"versiya_http_emas", "GET /products/1 FTP/1.1", Request{}, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseRequestLine(tc.line)
			if tc.wantErr {
				if err == nil {
					t.Errorf("parseRequestLine(%q): xato kutilgan, nil qaytdi (TASK.md, 1-qadam)", tc.line)
				}
				return
			}
			if err != nil {
				t.Errorf("parseRequestLine(%q): xato kutilmagan, keldi: %v (TASK.md, 1-qadam)", tc.line, err)
				return
			}
			if got != tc.want {
				t.Errorf("parseRequestLine(%q) = %+v; kutilgan %+v (TASK.md, 1-qadam)", tc.line, got, tc.want)
			}
		})
	}
}

// 2-qadam: statusClass kodning birinchi raqamiga qarab sinfni qaytaradi.
func TestVazifaStatusClass(t *testing.T) {
	tests := []struct {
		code int
		want string
	}{
		{100, "axborot"},
		{200, "muvaffaqiyat"},
		{201, "muvaffaqiyat"},
		{299, "muvaffaqiyat"},
		{301, "yoʻnaltirish"},
		{400, "mijoz xatosi"},
		{404, "mijoz xatosi"},
		{415, "mijoz xatosi"},
		{500, "server xatosi"},
		{503, "server xatosi"},
		{99, "nomaʼlum"},
		{600, "nomaʼlum"},
	}
	for _, tc := range tests {
		if got := statusClass(tc.code); got != tc.want {
			t.Errorf("statusClass(%d) = %q; kutilgan %q (TASK.md, 2-qadam)", tc.code, got, tc.want)
		}
	}
}
