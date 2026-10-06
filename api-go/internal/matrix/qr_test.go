package matrix

import (
	"math"
	"testing"
)

const tol = 1e-9

func TestQR(t *testing.T) {
	cases := []struct {
		name string
		a    [][]float64
	}{
		{"cuadrada 3x3", [][]float64{{12, -51, 4}, {6, 167, -68}, {-4, 24, -41}}},
		{"alta 4x2", [][]float64{{1, 2}, {3, 4}, {5, 6}, {7, 8}}},
		{"ancha 2x4", [][]float64{{1, 2, 3, 4}, {5, 6, 7, 8}}},
		{"diagonal", [][]float64{{2, 0}, {0, 3}}},
		{"una fila", [][]float64{{-3, 4}}},
		{"columna nula", [][]float64{{0, 1}, {0, 2}}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			q, r := QR(tc.a)
			m, n := len(tc.a), len(tc.a[0])
			k := min(m, n)

			// 1) A = Q·R
			for i := 0; i < m; i++ {
				for j := 0; j < n; j++ {
					s := 0.0
					for l := 0; l < k; l++ {
						s += q[i][l] * r[l][j]
					}
					if math.Abs(s-tc.a[i][j]) > tol {
						t.Fatalf("A != QR en [%d][%d]: %v vs %v", i, j, s, tc.a[i][j])
					}
				}
			}
			// 2) Qᵀ·Q = I y 3) R triangular superior con diagonal >= 0
			for i := 0; i < k; i++ {
				for j := 0; j < k; j++ {
					s := 0.0
					for l := 0; l < m; l++ {
						s += q[l][i] * q[l][j]
					}
					want := 0.0
					if i == j {
						want = 1
					}
					if math.Abs(s-want) > tol {
						t.Fatalf("QᵀQ no es identidad en [%d][%d]: %v", i, j, s)
					}
				}
				if r[i][i] < 0 {
					t.Fatalf("diagonal negativa en R[%d][%d]", i, i)
				}
				for j := 0; j < i && j < n; j++ {
					if r[i][j] != 0 {
						t.Fatalf("R no es triangular en [%d][%d]", i, j)
					}
				}
			}
		})
	}
}

func TestQRValoresConocidos(t *testing.T) {
	// Ejemplo clásico: R esperada = [[14,21,-14],[0,175,-70],[0,0,35]]
	_, r := QR([][]float64{{12, -51, 4}, {6, 167, -68}, {-4, 24, -41}})
	want := [][]float64{{14, 21, -14}, {0, 175, -70}, {0, 0, 35}}
	for i := range want {
		for j := range want[i] {
			if math.Abs(r[i][j]-want[i][j]) > tol {
				t.Fatalf("R[%d][%d] = %v, se esperaba %v", i, j, r[i][j], want[i][j])
			}
		}
	}
}

func TestValidate(t *testing.T) {
	malas := map[string][][]float64{
		"vacía":      {},
		"fila vacía": {{}},
		"irregular":  {{1, 2}, {3}},
		"NaN":        {{math.NaN()}},
	}
	for nombre, a := range malas {
		if Validate(a) == nil {
			t.Errorf("%s: se esperaba error", nombre)
		}
	}
	if err := Validate([][]float64{{1, 2}, {3, 4}}); err != nil {
		t.Errorf("matriz válida rechazada: %v", err)
	}
}
