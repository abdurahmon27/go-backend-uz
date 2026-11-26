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
    uv run tajriba.py ming 100   # xuddi shu, 100 ta train talaba bilan
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


def xato(p, x, y):
    return float(np.mean(np.abs(p(x) - y)))  # oʻrtacha necha ballga adashdi


def tajriba(seed, n_train=N_TRAIN):
    rng = np.random.default_rng(seed)
    xs, ys = talabalar(rng, n_train)  # train: ball = 40 + 5·soat + shovqin
    xt, yt = talabalar(rng, N_TEST)  # test: modellar bularni koʻrmagan
    natija = {}
    for nom, daraja in MODELLAR:  # chiziq (1) va daraja 9
        p = Polynomial.fit(xs, ys, daraja)  # least squares; daraja 9 = 10 parametr
        natija[nom] = {"train": xato(p, xs, ys), "test": xato(p, xt, yt), "p": p}
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
    test = {nom: [] for nom, _ in MODELLAR}
    for s in range(1, 6):
        natija, _, _ = tajriba(s, n_train)
        for nom, _ in MODELLAR:
            test[nom].append(natija[nom]["test"])
    print(f"train: {n_train} talaba · test xato, seed 1…5")
    for nom, _ in MODELLAR:
        print(f"{nom:<9}" + "".join(f"{v:>9.1f}" for v in test[nom]))
    yomon = sum(d > c for c, d in zip(test["chiziq"], test["daraja 9"]))
    print(f"daraja 9 testda yomonroq: {yomon} / 5")


def ming(n_train=N_TRAIN):
    """Bashoratlar nega xato chiqdi? Bitta seed emas, 1000 ta seed (natijadan keyin qoʻshilgan tahlil)."""
    R = [tajriba(s, n_train) for s in range(1, 1001)]
    tr = np.array([n["chiziq"]["train"] for n, _, _ in R])
    te = np.array([n["chiziq"]["test"] for n, _, _ in R])
    d9 = np.array([n["daraja 9"]["test"] for n, _, _ in R])
    ichida = 0
    for n, (xs, _), (xt, yt) in R:
        p = n["daraja 9"]["p"]
        i = int(np.argmax(np.abs(p(xt) - yt)))
        ichida += xs.min() <= xt[i] <= xs.max()
    print(f"1000 seed · train: {n_train} talaba")
    print(f"chiziq train xato, oʻrtacha: {tr.mean():.1f}")
    m = float(np.median(d9))
    print(f"daraja 9 test xato, mediana: {m:.0f}" if m >= 100 else f"daraja 9 test xato, mediana: {m:.1f}")
    print(f"daraja 9 test 10…100 oraligʻida: {np.sum((d9 >= 10) & (d9 <= 100))} / 1000")
    print(f"eng katta xato train oraligʻi ichida: {ichida} / 1000")
    print(f"daraja 9 testda yomonroq: {np.sum(d9 > te)} / 1000")


if __name__ == "__main__":
    rejim = sys.argv[1] if len(sys.argv) > 1 else ""
    if rejim == "seeds":
        seeds(int(sys.argv[2]) if len(sys.argv) > 2 else N_TRAIN)
    elif rejim == "ming":
        ming(int(sys.argv[2]) if len(sys.argv) > 2 else N_TRAIN)
    else:
        bitta()
