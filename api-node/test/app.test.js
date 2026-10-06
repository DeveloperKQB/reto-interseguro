import { test, describe } from 'node:test';
import assert from 'node:assert/strict';
import request from 'supertest';
import jwt from 'jsonwebtoken';
import { createApp } from '../src/app.js';

// Pruebas de integración: levantan la app en memoria, sin abrir un puerto real.
const SECRET = 'secreto-de-pruebas-de-al-menos-32-caracteres';
const app = createApp({ jwtSecret: SECRET });
const token = jwt.sign({ sub: 'tester' }, SECRET, { algorithm: 'HS256', expiresIn: '5m' });
const auth = `Bearer ${token}`;

describe('POST /api/stats', () => {
  test('200 con estadísticas para matrices válidas', async () => {
    const res = await request(app)
      .post('/api/stats')
      .set('Authorization', auth)
      .send({ matrices: [[[1, 0], [0, 1]], [[2, 0], [0, 3]]] });

    assert.equal(res.status, 200);
    assert.equal(res.body.max, 3);
    assert.equal(res.body.sum, 7);
    assert.equal(res.body.anyDiagonal, true);
  });

  test('400 si falta el campo matrices', async () => {
    const res = await request(app).post('/api/stats').set('Authorization', auth).send({});
    assert.equal(res.status, 400);
    assert.ok(res.body.error);
  });

  test('400 con JSON mal formado', async () => {
    const res = await request(app)
      .post('/api/stats')
      .set('Authorization', auth)
      .set('Content-Type', 'application/json')
      .send('{malo');
    assert.equal(res.status, 400);
    assert.ok(res.body.error);
  });

  test('no expone la cabecera X-Powered-By', async () => {
    const res = await request(app).get('/health');
    assert.equal(res.headers['x-powered-by'], undefined);
  });
});

describe('Seguridad JWT', () => {
  const body = { matrices: [[[1]]] };

  test('401 sin token', async () => {
    const res = await request(app).post('/api/stats').send(body);
    assert.equal(res.status, 401);
  });

  test('401 con token firmado con otro secreto', async () => {
    const falso = jwt.sign({ sub: 'x' }, 'otro-secreto-distinto-de-32-caracteres!!', { algorithm: 'HS256' });
    const res = await request(app).post('/api/stats').set('Authorization', `Bearer ${falso}`).send(body);
    assert.equal(res.status, 401);
  });

  test('401 con token expirado', async () => {
    const expirado = jwt.sign({ sub: 'x', exp: Math.floor(Date.now() / 1000) - 60 }, SECRET);
    const res = await request(app).post('/api/stats').set('Authorization', `Bearer ${expirado}`).send(body);
    assert.equal(res.status, 401);
  });

  test('/health sigue siendo público', async () => {
    const res = await request(app).get('/health');
    assert.equal(res.status, 200);
  });
});