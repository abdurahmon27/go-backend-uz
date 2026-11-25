// Anthropic'ning ESKI modellar uchun chiqargan ochiq tokenizer'i (@anthropic-ai/tokenizer).
// Uning README'si: Claude 3 dan beri bu algoritm aniq emas, faqat juda taxminiy.
// Ishga tushirish: npm install && node eski_claude.mjs
import { countTokens } from "@anthropic-ai/tokenizer";
import { readFileSync } from "node:fs";

const eng = readFileSync(new URL("udhr_eng.txt", import.meta.url), "utf8");
const uzb = readFileSync(new URL("udhr_uzn_latn.txt", import.meta.url), "utf8");
const e = countTokens(eng), u = countTokens(uzb);
const r = (s, n) => String(s).padStart(n);
console.log("tokenlar: Claude 3 gacha".padEnd(34) + " " + r(e, 9) + " " + r(u, 9) + " " + r((u / e).toFixed(2) + "x", 6));
