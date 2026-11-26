# Taqqoslash uchun: xuddi shu API Express bilan, rasmiy Node image'da
FROM node:24-slim
WORKDIR /app
COPY package*.json ./
RUN npm ci --omit=dev
COPY index.js ./
CMD ["node", "index.js"]
