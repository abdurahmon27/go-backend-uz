// TS: const a = 3 → literal type 3; let b = a → number ga kengayadi ("widening")
const a = 3
let b = a
const ok: 3 = a
const bad: 3 = b // TS2322: Type 'number' is not assignable to type '3'.
