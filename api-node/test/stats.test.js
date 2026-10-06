import { test, describe } from 'node:test';
import assert from 'node:assert/strict';
import { computeStats, isDiagonal, validateMatrices } from '../src/stats.js';

describe('computeStats', () => {
  test('calcula max, min, suma y promedio sobre todas las matrices', () => {
    const stats = computeStats([
      [[1, 2], [3, 4]],
      [[-5, 10]],
    ]);
    assert.equal(stats.max, 10);
    assert.equal(stats.min, -5);
    assert.equal(stats.sum, 15);
    assert.equal(stats.count, 6);
    assert.equal(stats.average, 2.5);
  });

  test('detecta si alguna matriz es diagonal', () => {
    const stats = computeStats([
      [[1, 2], [3, 4]],
      [[2, 0], [0, 3]],
    ]);
    assert.equal(stats.anyDiagonal, true);
    assert.deepEqual(stats.diagonalByMatrix, [false, true]);
  });
});

describe('isDiagonal', () => {
  test('matriz diagonal cuadrada', () => {
    assert.equal(isDiagonal([[1, 0], [0, 5]]), true);
  });

  test('matriz no diagonal', () => {
    assert.equal(isDiagonal([[1, 2], [0, 5]]), false);
  });

  test('tolera residuos de punto flotante', () => {
    assert.equal(isDiagonal([[1, 1e-17], [-1e-15, 1]]), true);
  });

  test('matriz rectangular diagonal (definición generalizada)', () => {
    assert.equal(isDiagonal([[1, 0, 0], [0, 2, 0]]), true);
  });
});

describe('validateMatrices', () => {
  const casosInvalidos = {
    'no es array': 'hola',
    'array vacío': [],
    'matriz vacía': [[]],
    'fila irregular': [[[1, 2], [3]]],
    'valor no numérico': [[[1, 'a']]],
    'valor null': [[[1, null]]],
  };

  for (const [nombre, entrada] of Object.entries(casosInvalidos)) {
    test(`rechaza: ${nombre}`, () => {
      assert.notEqual(validateMatrices(entrada), null);
    });
  }

  test('acepta matrices válidas', () => {
    assert.equal(validateMatrices([[[1, 2], [3, 4]], [[5]]]), null);
  });
});