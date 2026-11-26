enum Status { Unknown, New, Cancelled, Paid } // Cancelled oʻrtaga qoʻshildi
const saved: number = 2 // kecha bazaga yozilgan: oʻshanda 2 = Paid
console.log(saved === Status.Paid, saved === Status.Cancelled) // false true
