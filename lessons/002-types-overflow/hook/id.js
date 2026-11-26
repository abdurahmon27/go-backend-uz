// Backend JSON'da 64 bitlik ID yubordi, frontend uni JSON.parse bilan son qilib o'qidi.
const body = '{"id":1234567890123456789}'
console.log(body)
console.log(JSON.parse(body).id)
