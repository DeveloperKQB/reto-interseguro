import jwt from 'jsonwebtoken';

/**
 * Middleware que exige "Authorization: Bearer <token>" firmado con HS256
 * usando el secreto compartido con la API de Go.
 * @param {string} secret
 */
export function requireJwt(secret) {
  return (req, res, next) => {
    const header = req.get('Authorization') ?? '';
    const [scheme, token] = header.split(' ');

    if (scheme !== 'Bearer' || !token) {
      return res.status(401).json({ error: 'token ausente, inválido o expirado' });
    }
    try {
      // Solo HS256: evita ataques de confusión de algoritmo.
      req.user = jwt.verify(token, secret, { algorithms: ['HS256'] });
      next();
    } catch {
      res.status(401).json({ error: 'token ausente, inválido o expirado' });
    }
  };
}