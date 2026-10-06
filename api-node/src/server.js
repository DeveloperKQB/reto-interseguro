import { createApp } from './app.js';

const PORT = process.env.PORT ?? 3000;
const JWT_SECRET = process.env.JWT_SECRET;

// Falla al arrancar en lugar de funcionar con un secreto inseguro.
if (!JWT_SECRET || JWT_SECRET.length < 32) {
  console.error('JWT_SECRET es obligatorio y debe tener al menos 32 caracteres');
  process.exit(1);
}

createApp({ jwtSecret: JWT_SECRET }).listen(PORT, () => {
  console.log(`api-node escuchando en http://localhost:${PORT}`);
});