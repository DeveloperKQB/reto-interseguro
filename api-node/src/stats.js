/**
 * Lógica pura de estadísticas (sin dependencias de Express).
 */

/** Tolerancia para considerar un valor como cero (residuos de punto flotante). */
export const DIAGONAL_EPS = 1e-10;

/**
 * Valida que `matrices` sea un array no vacío de matrices rectangulares
 * con números finitos. Devuelve un mensaje de error o null si es válido.
 * @param {unknown} matrices
 * @returns {string|null}
 */
export function validateMatrices(matrices) {
  if (!Array.isArray(matrices) || matrices.length === 0) {
    return 'se espera "matrices": un array no vacío de matrices';
  }
  for (let k = 0; k < matrices.length; k++) {
    const m = matrices[k];
    if (!Array.isArray(m) || m.length === 0 || !Array.isArray(m[0]) || m[0].length === 0) {
      return `la matriz ${k} está vacía o mal formada`;
    }
    const cols = m[0].length;
    for (let i = 0; i < m.length; i++) {
      if (!Array.isArray(m[i]) || m[i].length !== cols) {
        return `la matriz ${k}: la fila ${i} no tiene ${cols} columnas`;
      }
      for (let j = 0; j < cols; j++) {
        const v = m[i][j];
        if (typeof v !== 'number' || !Number.isFinite(v)) {
          return `la matriz ${k}: valor no numérico en [${i}][${j}]`;
        }
      }
    }
  }
  return null;
}

/**
 * Una matriz es diagonal si todo elemento fuera de la diagonal principal
 * (i !== j) es cero, con tolerancia. Se acepta la definición generalizada
 * para matrices rectangulares.
 * @param {number[][]} m
 * @param {number} [eps]
 */
export function isDiagonal(m, eps = DIAGONAL_EPS) {
  for (let i = 0; i < m.length; i++) {
    for (let j = 0; j < m[i].length; j++) {
      if (i !== j && Math.abs(m[i][j]) > eps) return false;
    }
  }
  return true;
}

/**
 * Calcula máximo, mínimo, promedio y suma de TODOS los valores de todas
 * las matrices en una sola pasada, e indica si alguna es diagonal.
 * @param {number[][][]} matrices (ya validadas)
 */
export function computeStats(matrices) {
  let max = -Infinity;
  let min = Infinity;
  let sum = 0;
  let count = 0;

  for (const m of matrices) {
    for (const row of m) {
      for (const v of row) {
        if (v > max) max = v;
        if (v < min) min = v;
        sum += v;
        count++;
      }
    }
  }

  const diagonalByMatrix = matrices.map((m) => isDiagonal(m));

  return {
    max,
    min,
    average: sum / count,
    sum,
    count,
    anyDiagonal: diagonalByMatrix.some(Boolean),
    diagonalByMatrix,
  };
}