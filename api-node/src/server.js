import { createApp } from './app.js';

const PORT = process.env.PORT ?? 3000;

createApp().listen(PORT, () => {
  console.log(`api-node escuchando en http://localhost:${PORT}`);
});