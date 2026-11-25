#!/usr/bin/env bash
# Serverni ishga tushiradi va stack trace'siz, faqat asosiy qatorlarni ko'rsatadi.
# CGO_ENABLED=0: net paketi bu mashinada C kompilyatorisiz yig'ilishi uchun.
cd "$(dirname "$0")/.." && CGO_ENABLED=0 go run ./server 2>&1 | grep -E '^(javob|http: panic|panic:|exit status)'
