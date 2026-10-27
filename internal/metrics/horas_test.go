package metrics

import (
	"errors"
	"math"
	"testing"
)

func TestHorasReales(t *testing.T) {
	casos := []struct {
		nombre string
		horas  []float64
		quiere float64
	}{
		{"varios registros", []float64{3, 2.5, 4}, 9.5},
		{"lista vacía", []float64{}, 0},
		{"nil", nil, 0},
		{"valores", []float64{1.5, 0, 2}, 3.5},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			got := HorasReales(c.horas)
			if math.Abs(got-c.quiere) > 1e-9 {
				t.Errorf("HorasReales() = %v, se esperaba %v", got, c.quiere)
			}
		})
	}
}

func TestDesviacion(t *testing.T) {
	casos := []struct {
		nombre    string
		estimadas float64
		reales    float64
		quiere    float64
		quiereErr error
	}{
		{"+20%", 100, 120, 20, nil},
		{"-20%", 100, 80, -20, nil},
		{"0%", 100, 100, 0, nil},
		{"-100% (reales 0)", 100, 0, -100, nil},
		{"estimadas 0 (ErrDesviacionNoCalculable)", 0, 50, 0, ErrDesviacionNoCalculable},
		{"estimadas 0 y reales 0 (ErrDesviacionNoCalculable)", 0, 0, 0, ErrDesviacionNoCalculable},
		{"horas negativas", -10, 50, 0, ErrHorasNegativas},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			got, err := Desviacion(c.estimadas, c.reales)
			if c.quiereErr != nil {
				if !errors.Is(err, c.quiereErr) {
					t.Fatalf("Desviacion(%v, %v) devolvió error %v, se esperaba %v", c.estimadas, c.reales, err, c.quiereErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Desviacion(%v, %v) devolvió error inesperado: %v", c.estimadas, c.reales, err)
			}
			if math.Abs(got-c.quiere) > 1e-9 {
				t.Errorf("Desviacion(%v, %v) = %v, se esperaba %v", c.estimadas, c.reales, got, c.quiere)
			}
		})
	}
}
