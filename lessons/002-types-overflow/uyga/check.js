// Frontend kabi: JSON'ni JSON.parse bilan o'qiydi va qiymatlarni chiqaradi.
let input = ''
process.stdin.on('data', (d) => (input += d))
process.stdin.on('end', () => {
  const post = JSON.parse(input)
  for (const [k, v] of Object.entries(post)) console.log(k, typeof v, String(v))
})
