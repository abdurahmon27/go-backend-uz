// brauzer — GB-6.07 videosidagi brauzer tajribalarini haqiqiy Chrome'da (headless) bajaradigan yordamchi.
// Chrome'ni chromedp kutubxonasi DevTools protokoli orqali boshqaradi. Har ishga tushishda yangi, bo'sh
// profil: cookie'lar faqat shu tajribaniki.
//
//	go run . httponly                     Dilnoza login → bozor sahifasida document.cookie va fetch("/me")
//	go run . csrf http://127.0.0.1:9090   Dilnoza login → begona sahifa → «Sovgʻani olish» tugmasi → buyurtma 41
//	go run . cors http://localhost:9090   Dilnoza login → begona sahifadagi skript /me'ni oʻqishga urinadi
//
// Oldin: bozor (go run ./cmd/api, :8765) va begona (go run ./cmd/begona, :9090) ishlab tursin.
// Bu kod darsning mavzusi emas (README): videoda faqat natijasi va server log'i koʻrsatiladi.
package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

const bozor = "http://localhost:8765"

// loginJS — bozor sahifasida (bir xil origin) POST /login: javobdagi Set-Cookie'ni brauzer oʻzi saqlaydi.
const loginJS = `fetch("/login", {method: "POST", headers: {"Content-Type": "application/json"},
  body: JSON.stringify({email: "dilnoza@bozor.test", password: "tarvuz-osmon-qalam-ariq-27"})})
  .then(async r => r.status + " " + (await r.text()).trim())`

// do — bitta qadam (natijasiz); eval — skript natijasi matn sifatida (Promise boʻlsa kutib).
func do(ctx context.Context, a chromedp.Action[chromedp.Void]) error {
	_, err := chromedp.Run(ctx, a)
	return err
}

func eval(ctx context.Context, js string) string {
	s, err := chromedp.Run(ctx, chromedp.Evaluate[string](js, chromedp.EvalAwaitPromise))
	if err != nil {
		return "xato: " + err.Error()
	}
	return s
}

func login(ctx context.Context) error {
	if err := do(ctx, chromedp.Navigate(bozor+"/health")); err != nil {
		return err
	}
	fmt.Println("login:", eval(ctx, loginJS))
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "foydalanish: go run . httponly | csrf URL | cors URL")
		os.Exit(2)
	}
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	if err := login(ctx); err != nil {
		fmt.Fprintln(os.Stderr, "xato:", err)
		os.Exit(1)
	}
	var err error
	switch os.Args[1] {
	case "httponly":
		err = httpOnly(ctx)
	case "csrf":
		err = csrf(ctx, os.Args[2])
	case "cors":
		err = corsRead(ctx, os.Args[2])
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "xato:", err)
		os.Exit(1)
	}
}

// httpOnly: sahifa skripti cookie'ni koʻrmaydi, lekin brauzer uni soʻrovga qoʻshadi.
func httpOnly(ctx context.Context) error {
	fmt.Printf("document.cookie: %q\n", eval(ctx, `document.cookie`))
	fmt.Println("fetch /me:", eval(ctx, `fetch("/me").then(async r => r.status + " " + (await r.text()).trim())`))
	return nil
}

// csrf: begona sahifada tugma bosiladi; brauzer formani bozor'ga yuboradi va javob sahifasini ochadi.
func csrf(ctx context.Context, page string) error {
	for _, a := range []chromedp.Action[chromedp.Void]{
		chromedp.Navigate(page),
		chromedp.WaitVisible("button"),
		chromedp.Click("button"),
		chromedp.Sleep(700 * time.Millisecond),
		chromedp.WaitReady("body"),
	} {
		if err := do(ctx, a); err != nil {
			return err
		}
	}
	fmt.Println("sahifa:", page, "→ «Sovgʻani olish» bosildi")
	fmt.Println("javob:", eval(ctx, `document.body.innerText.trim()`))
	if err := do(ctx, chromedp.Navigate(bozor+"/orders/41")); err != nil {
		return err
	}
	fmt.Println("GET /orders/41:", eval(ctx, `document.body.innerText.trim()`))
	return nil
}

// corsRead: begona sahifadagi skript GET /me javobini oʻqishga urinadi.
func corsRead(ctx context.Context, page string) error {
	if err := do(ctx, chromedp.Navigate(page+"/oqish")); err != nil {
		return err
	}
	if _, err := chromedp.Run(ctx, chromedp.Poll[chromedp.Void](`natija.textContent !== "kutilmoqda…"`)); err != nil {
		return err
	}
	fmt.Println("sahifa:", page+"/oqish")
	fmt.Println(strings.TrimSpace(eval(ctx, `natija.textContent`)))
	return nil
}
