// TS: type Ms = number — faqat taxallus (alias): istalgan number oʻtadi
type Ms = number
function sleep(ms: Ms) {}
const n: number = 2
sleep(n) // ✅ xato yoʻq. Go'da type Duration int64 esa alohida tur: int oʻtmaydi
