package metrics

import (
	"errors"
	"math"
	"testing"
)

func TestPorcentajeCompletadas(t *testing.T) {
	casos := []struct {
		nombre      string
		completadas int
		total       int
		quiere      float64
		quiereErr   error
	}{
		{"3 de 4", 3, 4, 75, nil},
		{"1 de 3", 1, 3, 100.0 / 3.0, nil},
		{"5 de 5", 5, 5, 100, nil},
		{"0 de 5", 0, 5, 0, nil},
		{"0 de 0", 0, 0, 0, nil},
		{"más completadas que el total", 6, 5, 0, ErrDatosInconsistentes},
		{"completadas negativas", -1, 5, 0, ErrDatosInconsistentes},
		{"total negativo", 1, -5, 0, ErrDatosInconsistentes},
		{"ambos negativos", -1, -5, 0, ErrDatosInconsistentes},
		{"completadas con total cero", 3, 0, 0, ErrDatosInconsistentes},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			got, err := PorcentajeCompletadas(c.completadas, c.total)
			if c.quiereErr != nil {
				if !errors.Is(err, c.quiereErr) {
					t.Fatalf("PorcentajeCompletadas(%d, %d) devolvió error %v, se esperaba %v", c.completadas, c.total, err, c.quiereErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("PorcentajeCompletadas(%d, %d) devolvió error inesperado: %v", c.completadas, c.total, err)
			}
			if math.Abs(got-c.quiere) > 1e-9 {
				t.Errorf("PorcentajeCompletadas(%d, %d) = %v, se esperaba %v", c.completadas, c.total, got, c.quiere)
			}
		})
	}
}

func TestErrDatosInconsistentes(t *testing.T) {
	const quiere = "las historias completadas no pueden superar al total"
	if got := ErrDatosInconsistentes.Error(); got != quiere {
		t.Errorf("ErrDatosInconsistentes = %q, se esperaba %q", got, quiere)
	}
}
