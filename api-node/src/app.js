import express from 'express';
import { validateMatrices, computeStats } from './stats.js';
import { requireJwt } from './auth.js';

/**
 * Crea la app de Express. Se separa de server.js para poder
 * probarla sin levantar el puerto (tests de integración).
 * @param {{ jwtSecret: string }} options
 */
export function createApp({ jwtSecret }) {
  const app = express();
  app.disable('x-powered-by');

  // Q y R de 100x100 suman ~20.000 números: subimos el límite por defecto (100kb).
  app.use(express.json({ limit: '1mb' }));

  // Ruta pública (usada por el healthcheck de Docker)
  app.get('/health', (req, res) => {
    res.json({ status: 'ok', service: 'api-node' });
  });

  /**
   * POST /api/stats  (protegida con JWT)
   * Body:  { "matrices": [ [[...]], [[...]] ] }
   * Resp:  { max, min, average, sum, count, anyDiagonal, diagonalByMatrix }
   */
  app.post('/api/stats', requireJwt(jwtSecret), (req, res) => {
    const matrices = req.body?.matrices;
    const error = validateMatrices(matrices);
    if (error) {
      return res.status(400).json({ error });
    }
    res.json(computeStats(matrices));
  });

  // JSON mal formado u otros errores: siempre respondemos JSON.
  app.use((err, req, res, next) => {
    const status = err.status ?? 500;
    res.status(status).json({ error: status === 500 ? 'error interno' : err.message });
  });

  return app;
}