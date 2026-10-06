package matrix

import "math"

// cleanEps: valores con magnitud menor se redondean a 0 para eliminar
// residuos de punto flotante (p. ej. 1e-17) y el "-0".
const cleanEps = 1e-12

// QR calcula la factorización QR reducida de A (m×n) mediante reflexiones
// de Householder: A = Q·R, con k = min(m, n),
//   - Q (m×k) con columnas ortonormales,
//   - R (k×n) triangular superior con diagonal no negativa.
//
// Se asume que A ya fue validada con Validate.
func QR(a [][]float64) (q, r [][]float64) {
	m, n := len(a), len(a[0])
	k := min(m, n)

	// R arranca como copia de A; Q como la identidad m×m.
	R := make([][]float64, m)
	for i := range a {
		R[i] = append([]float64(nil), a[i]...)
	}
	Q := identity(m)

	v := make([]float64, m)
	for col := 0; col < min(m-1, n); col++ {
		// x = R[col:, col]; construimos v tal que H = I - 2vvᵀ anule x bajo la diagonal.
		norm := 0.0
		for i := col; i < m; i++ {
			norm += R[i][col] * R[i][col]
		}
		norm = math.Sqrt(norm)
		if norm == 0 {
			continue // columna ya nula: nada que reflejar
		}
		// Signo opuesto a x0 para evitar cancelación numérica.
		alpha := -norm
		if R[col][col] < 0 {
			alpha = norm
		}
		vNorm := 0.0
		for i := col; i < m; i++ {
			v[i] = R[i][col]
			if i == col {
				v[i] -= alpha
			}
			vNorm += v[i] * v[i]
		}
		vNorm = math.Sqrt(vNorm)
		if vNorm == 0 {
			continue
		}
		for i := col; i < m; i++ {
			v[i] /= vNorm
		}

		// R = H·R (solo afecta filas col..m-1 y columnas col..n-1).
		for j := col; j < n; j++ {
			dot := 0.0
			for i := col; i < m; i++ {
				dot += v[i] * R[i][j]
			}
			for i := col; i < m; i++ {
				R[i][j] -= 2 * v[i] * dot
			}
		}
		// Q = Q·H (acumula las reflexiones).
		for i := 0; i < m; i++ {
			dot := 0.0
			for l := col; l < m; l++ {
				dot += Q[i][l] * v[l]
			}
			for l := col; l < m; l++ {
				Q[i][l] -= 2 * dot * v[l]
			}
		}
	}

	// Normalización: diagonal de R no negativa => QR única (si A es de rango completo).
	for i := 0; i < k; i++ {
		if R[i][i] < 0 {
			for j := 0; j < n; j++ {
				R[i][j] = -R[i][j]
			}
			for row := 0; row < m; row++ {
				Q[row][i] = -Q[row][i]
			}
		}
	}

	// Forma reducida: Q[:, :k] y R[:k, :], con ceros exactos bajo la diagonal.
	q = make([][]float64, m)
	for i := 0; i < m; i++ {
		q[i] = clean(Q[i][:k])
	}
	r = make([][]float64, k)
	for i := 0; i < k; i++ {
		for j := 0; j < i && j < n; j++ {
			R[i][j] = 0
		}
		r[i] = clean(R[i])
	}
	return q, r
}

func identity(n int) [][]float64 {
	id := make([][]float64, n)
	for i := range id {
		id[i] = make([]float64, n)
		id[i][i] = 1
	}
	return id
}

func clean(row []float64) []float64 {
	out := make([]float64, len(row))
	for i, x := range row {
		if math.Abs(x) >= cleanEps {
			out[i] = x
		}
	}
	return out
}
