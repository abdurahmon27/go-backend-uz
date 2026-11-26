# /// script
# requires-python = ">=3.10"
# dependencies = ["numpy==2.5.5"]
# ///
"""Guru vazifasi uchun: seed 1 train nuqtalarini nuqtalar.csv ga yozadi va numpy chizigʻini chiqaradi.

    uv run homework/nuqtalar.py        (lessons/ml-01-map papkasidan)
"""
import importlib.util
from pathlib import Path

from numpy.polynomial import Polynomial

papka = Path(__file__).resolve().parent
spec = importlib.util.spec_from_file_location("tajriba", papka.parent / "tajriba.py")
t = importlib.util.module_from_spec(spec)
spec.loader.exec_module(t)

_, (xs, ys), _ = t.tajriba(1)
with open(papka / "nuqtalar.csv", "w", encoding="utf-8") as f:
    for x, y in zip(xs, ys):
        f.write(f"{float(x)!r},{float(y)!r}\n")
a, b = Polynomial.fit(xs, ys, 1).convert().coef
print(f"numpy: ball = {a:.12f} + {b:.12f} × soat")
