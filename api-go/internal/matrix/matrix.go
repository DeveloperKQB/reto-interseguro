// Package matrix contiene la lógica matemática: validación y factorización QR.
package matrix

import (
	"fmt"
	"math"
)

// MaxDim limita filas y columnas para evitar peticiones abusivas (QR es O(m·n²)).
const MaxDim = 100

// Validate comprueba que la matriz sea rectangular, no vacía, de tamaño
// razonable y con valores finitos.
func Validate(a [][]float64) error {
	if len(a) == 0 {
		return fmt.Errorf("la matriz no puede estar vacía")
	}
	cols := len(a[0])
	if cols == 0 {
		return fmt.Errorf("las filas no pueden estar vacías")
	}
	if len(a) > MaxDim || cols > MaxDim {
		return fmt.Errorf("la matriz excede el tamaño máximo de %dx%d", MaxDim, MaxDim)
	}
	for i, row := range a {
		if len(row) != cols {
			return fmt.Errorf("la fila %d tiene %d columnas, se esperaban %d", i, len(row), cols)
		}
		for j, v := range row {
			if math.IsNaN(v) || math.IsInf(v, 0) {
				return fmt.Errorf("valor no finito en la posición [%d][%d]", i, j)
			}
		}
	}
	return nil
}
