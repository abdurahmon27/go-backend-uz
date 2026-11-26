type Item = { price: number };

const item: Item = JSON.parse('{"price": "2"}');
console.log(item.price + 2);
