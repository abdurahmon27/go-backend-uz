# /// script
# requires-python = ">=3.10"
# dependencies = ["tiktoken==0.14.0"]
# ///
"""Bir xil matn ingliz va oʻzbek tilida nechta token boʻladi?

Claude tokenizer'i ochiq emas, shuning uchun bu yerda OCHIQ tokenizer'lar ishlatiladi:
tiktoken'dagi o200k_base va cl100k_base (OpenAI). Natija Claude uchun aniq son emas,
faqat tendensiyani koʻrsatadi. Ishga tushirish: uv run tokenlar.py
"""
import sys
from pathlib import Path

import tiktoken

papka = Path(__file__).parent
eng = (papka / "udhr_eng.txt").read_text(encoding="utf-8")
uzb = (papka / "udhr_uzn_latn.txt").read_text(encoding="utf-8")

# Oʻzbek matnida apostrof uch xil yoziladi: ʻ (U+02BB, rasmiy), ' (klaviatura), ‘ (Word avtozamenasi).
ascii_ = uzb.replace("ʻ", "'").replace("ʼ", "'")
word = uzb.replace("ʻ", "‘").replace("ʼ", "’")


def sana(enc, matn):
    return len(enc.encode(matn))


rejim = sys.argv[1] if len(sys.argv) > 1 else "hammasi"
QATOR = "{:34} {:>9} {:>9} {:>6}"

if rejim in ("hammasi", "til"):
    o200k = tiktoken.get_encoding("o200k_base")
    print(QATOR.format("Inson huquqlari deklaratsiyasi", "inglizcha", "oʻzbekcha", "farq"))
    print(QATOR.format("belgilar", len(eng), len(uzb), f"{len(uzb) / len(eng):.2f}x"))
    e, u = sana(o200k, eng), sana(o200k, uzb)
    print(QATOR.format("tokenlar: o200k_base (OpenAI)", e, u, f"{u / e:.2f}x"))

if rejim in ("hammasi", "cl100k"):
    cl = tiktoken.get_encoding("cl100k_base")
    e, u = sana(cl, eng), sana(cl, uzb)
    print(QATOR.format("tokenlar: cl100k_base (OpenAI)", e, u, f"{u / e:.2f}x"))

if rejim in ("hammasi", "apostrof"):
    enc = tiktoken.get_encoding("o200k_base")
    print("oʻzbekcha matn, o200k_base, apostrof turi boʻyicha")
    for nom, matn in (("ʻ  U+02BB", uzb), ("'  ASCII", ascii_), ("‘  U+2018", word)):
        print(f"  {nom:10} {sana(enc, matn):5} token")

if rejim in ("hammasi", "soz"):
    enc = tiktoken.get_encoding("o200k_base")
    for s in ("we are learning", "oʻrganmoqdamiz", "o'rganmoqdamiz"):
        qism = [enc.decode([t]) for t in enc.encode(s)]
        print(f"{s:16} {len(qism)} token: " + " | ".join(qism))
