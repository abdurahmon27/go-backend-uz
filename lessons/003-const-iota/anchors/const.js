// JS const faqat nomni qayta bogʻlashni taqiqlaydi
const cfg = { port: 80 };
cfg.port = 81;            // ✅ obyekt ichi oʻzgaradi
console.log(cfg.port);    // 81
const start = Date.now(); // ✅ qiymat ishlash paytida olinadi
console.log(typeof start); // number
