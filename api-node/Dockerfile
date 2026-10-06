FROM node:24-alpine
WORKDIR /app

# Primero las dependencias, para aprovechar la caché de Docker
COPY package.json package-lock.json ./
RUN npm ci --omit=dev

COPY src ./src

ENV NODE_ENV=production
USER node
EXPOSE 3000
CMD ["node", "src/server.js"]