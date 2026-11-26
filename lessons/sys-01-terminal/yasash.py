#!/usr/bin/env python3
"""yil/ papkasini yasaydi: 12 oy × 30 qator jamoa eslatmalari (videoning «bir yillik» papkasi).

Qoida: har bir TODO — alohida, hali qilinmagan ish (bir xil TODO ikki marta yozilmaydi va hech qachon
«qilindi» deb yozilmaydi). Qilingan ishlar — kundalik ishlar, ular hech qachon TODO boʻlmaydi.
Shuning uchun «TODO qatorlari soni = qolgan ishlar soni». Natija har safar bir xil (random.seed).
Ishga tushirish (sys-01-terminal papkasidan): python3 yasash.py
"""
import random
from pathlib import Path

OYLAR = ["yanvar", "fevral", "mart", "aprel", "may", "iyun",
         "iyul", "avgust", "sentabr", "oktabr", "noyabr", "dekabr"]
ODAMLAR = ["Aziz", "Malika", "Bobur", "Dilnoza"]
QILINGAN = ["pochtani tekshirdi", "zaxira nusxa oldi", "narxlarni yangiladi", "mijozga javob yozdi",
            "hisobotni topshirdi", "yangi rasmlar yukladi", "parolni yangiladi", "saytni tekshirdi",
            "uchrashuvda qatnashdi", "menyu matnini oʻqidi", "jadvalni toʻldirdi", "kompyuterni yangiladi"]
SAHIFALAR = ["bosh", "aloqa", "narxlar", "yangiliklar", "savat", "toʻlov", "kirish", "yordam", "blog",
             "aksiya", "katalog", "buyurtma", "profil", "sharhlar", "yetkazish", "kafolat", "vakansiya",
             "hamkorlar", "galereya", "xarita", "obuna", "sozlamalar", "bildirishnoma", "qidiruv", "tarix"]
ISHLAR = ["{} sahifasini tuzatish", "{} sahifasiga rasm tayyorlash", "{} sahifasi matnini yozish",
          "{} sahifasiga tugma qoʻshish"]

random.seed(1964)
todo = [i.format(s) for s in SAHIFALAR for i in ISHLAR]  # 100 ta har xil ish
random.shuffle(todo)
papka = Path(__file__).resolve().parent / "yil"
papka.mkdir(exist_ok=True)
for i, oy in enumerate(OYLAR, 1):
    qatorlar = []
    for _ in range(30):
        kim = random.choice(ODAMLAR)
        if todo and random.random() < 0.25:
            qatorlar.append(f"TODO {kim}: {todo.pop()}")
        else:
            qatorlar.append(f"{kim}: {random.choice(QILINGAN)}")
    (papka / f"{i:02d}-{oy}.txt").write_text("\n".join(qatorlar) + "\n")
