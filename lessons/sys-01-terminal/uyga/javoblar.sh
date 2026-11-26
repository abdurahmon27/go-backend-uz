#!/usr/bin/env bash
# SYS-0.01 · Uy vazifasi javoblari. Avval oʻzingiz urinib koʻring!
# Ishga tushirish (sys-01-terminal papkasidan): bash uyga/javoblar.sh
set -u

echo "== 🟢 Oson: Malika nechta ishni qilmagan? (Aziz → Malika)"
grep -r TODO eslatmalar | grep Malika | wc -l

echo "== 🟡 Oʻrta: yil papkasida Boburning hamma yozuvlari (qilinganlari bilan)"
# birinchi ishchini (grep -r TODO) olib tashlaymiz, papkani grep Bobur'ga beramiz: Bobur bor hamma qator
grep -r Bobur yil | wc -l

echo "== 🔍 Qiziqish: yil papkasida kim eng koʻp ish qoldirgan?"
# 1-yoʻl: har bir ism uchun videodagi qatorni ishga tushirish
for kim in Aziz Malika Bobur Dilnoza; do
  echo "$kim $(grep -r TODO yil | grep "$kim" | wc -l)"
done
# 2-yoʻl: yana ikki ishchi. cut -d' ' -f2 qatorni boʻsh joydan boʻlib, 2-boʻlakni (ismni) oladi;
# sort bir xil ismlarni yonma-yon qoʻyadi, uniq -c yonma-yon turganlarini sanaydi, sort -rn kattasini tepaga chiqaradi
grep -rh TODO yil | cut -d' ' -f2 | tr -d : | sort | uniq -c | sort -rn
