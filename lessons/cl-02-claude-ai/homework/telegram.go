// Package homework: CL-1.02 uy vazifasining «Qiyin» javobi.
//
// Telegram bot tokeni: «bot raqami : kalit». Rasmiy Bot API hujjatidagi namuna:
// 123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11 (https://core.telegram.org/bots/api#authorizing-your-bot).
// Asosiy skanerga qoʻshish uchun main.go dagi qoidalar roʻyxatiga shu qatorni yozing:
//
//	{"telegram-token", true, TelegramToken, nil},
package homework

import "regexp"

// TelegramToken: 6–12 xonali bot raqami, ikki nuqta, keyin 30+ ta harf/raqam/_/-.
var TelegramToken = regexp.MustCompile(`\b\d{6,12}:[A-Za-z0-9_-]{30,}`)
