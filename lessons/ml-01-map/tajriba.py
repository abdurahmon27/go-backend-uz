# /// script
# requires-python = ">=3.10"
# dependencies = ["numpy==2.5.5"]
# ///
"""Toʻgʻri chiziq yoki har nuqtadan oʻtadigan egri chiziq: qaysi biri YANGI talabalarda yaxshiroq?

Maʼlumot sintetik (oʻzimiz yasaymiz, shuning uchun haqiqiy qonuniyatni bilamiz):
    ball = 40 + 5 × soat + shovqin,   shovqin ~ Normal(0, 5)
10 ta talaba bilan oʻqitamiz (train), modellar koʻrmagan 10 ta yangi talabada tekshiramiz (test).
Xato = oʻrtacha mutlaq xato: model oʻrtacha necha ballga adashdi.

Faqat CPU va numpy. Bashorat bashorat.md faylida, u shu skript birinchi marta ishga tushishidan
OLDIN commit qilingan (git log bilan tekshiring).

    uv run tajriba.py          # seed 1: jadval va eng katta xato
    uv run tajriba.py seeds    # 1..5 seed: tasodifmi yoki qonuniyatmi?
    uv run tajriba.py seeds 100  # xuddi shu, lekin 100 ta train talaba (bashorat2.md)
    uv run tajriba.py ming     # 1000 seed: bashoratlar nega xato chiqdi (natijadan keyin qoʻshilgan)
"""
import sys

import numpy as np
from numpy.polynomial import Polynomial

N_TRAIN = 10
N_TEST = 10
MODELLAR = [("chiziq", 1), ("daraja 9", 9)]  # 10 nuqta, 10 parametr: daraja 9 har nuqtadan oʻtadi


def talabalar(rng, n):
    soat = rng.uniform(0, 10, n)
    ball = 40 + 5 * soat + rng.normal(0, 5, n)
    return soat, ball


def tajriba(seed, n_train=N_TRAIN):
    rng = np.random.default_rng(seed)
    xs, ys = talabalar(rng, n_train)
    xt, yt = talabalar(rng, N_TEST)
    natija = {}
    for nom, daraja in MODELLAR:
        p = Polynomial.fit(xs, ys, daraja)  # least squares; domain [min, max] ga moslanadi
        natija[nom] = {
            "train": float(np.mean(np.abs(p(xs) - ys))),
            "test": float(np.mean(np.abs(p(xt) - yt))),
            "p": p,
        }
    return natija, (xs, ys), (xt, yt)


def bitta():
    natija, (xs, _), (xt, yt) = tajriba(1)
    print(f"seed 1 · train: {N_TRAIN} talaba · test: {N_TEST} yangi talaba")
    print(f"{'model':<10}{'train xato':>12}{'test xato':>11}")
    for nom, _ in MODELLAR:
        r = natija[nom]
        print(f"{nom:<10}{r['train']:>12.1f}{r['test']:>11.1f}")
    p = natija["daraja 9"]["p"]
    i = int(np.argmax(np.abs(p(xt) - yt)))
    print(f"eng katta xato: {xt[i]:.1f} soat, haqiqiy {yt[i]:.0f}, model {p(xt[i]):.0f}")
    print(f"train soatlari: {xs.min():.1f} … {xs.max():.1f}")


def seeds(n_train=N_TRAIN):
    print(f"{'seed':<6}{'chiziq test':>13}{'daraja 9 test':>15}   train: {n_train}")
    yomon = 0
    for s in range(1, 6):
        natija, _, _ = tajriba(s, n_train)
        c, d = natija["chiziq"]["test"], natija["daraja 9"]["test"]
        yomon += d > c
        print(f"{s:<6}{c:>13.1f}{d:>15.1f}")
    print(f"daraja 9 testda yomonroq: {yomon} / 5")


def ming():
    """Bashoratlar nega xato chiqdi? Bitta seed emas, 1000 ta seed (natijadan keyin qoʻshilgan tahlil)."""
    R = [tajriba(s) for s in range(1, 1001)]
    tr = np.array([n["chiziq"]["train"] for n, _, _ in R])
    te = np.array([n["chiziq"]["test"] for n, _, _ in R])
    d9 = np.array([n["daraja 9"]["test"] for n, _, _ in R])
    ichida = 0
    for n, (xs, _), (xt, yt) in R:
        p = n["daraja 9"]["p"]
        i = int(np.argmax(np.abs(p(xt) - yt)))
        ichida += xs.min() <= xt[i] <= xs.max()
    print("1000 seed:")
    print(f"chiziq train xato, oʻrtacha: {tr.mean():.1f}")
    print(f"daraja 9 test xato, mediana: {np.median(d9):.0f}")
    print(f"daraja 9 test 10…100 oraligʻida: {np.sum((d9 >= 10) & (d9 <= 100))} / 1000")
    print(f"eng katta xato train oraligʻi ichida: {ichida} / 1000")
    print(f"daraja 9 testda yomonroq: {np.sum(d9 > te)} / 1000")


if __name__ == "__main__":
    rejim = sys.argv[1] if len(sys.argv) > 1 else ""
    if rejim == "seeds":
        seeds(int(sys.argv[2]) if len(sys.argv) > 2 else N_TRAIN)
    elif rejim == "ming":
        ming()
    else:
        bitta()
