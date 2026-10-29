package metrics

import (
	"math"
	"testing"
	"time"
)

func TestHorasEnPeriodo(t *testing.T) {
	dia := func(d int) time.Time { return time.Date(2026, 10, d, 0, 0, 0, 0, time.Local) }
	inicio, fin := dia(23), dia(25)
	casos := []struct {
		nombre    string
		registros []RegistroHoras
		quiere    float64
	}{
		{"sin registros", nil, 0},
		{"todos dentro del período", []RegistroHoras{{Fecha: dia(23), Horas: 2}, {Fecha: dia(24), Horas: 1.5}}, 3.5},
		{"el día de inicio y el de fin cuentan", []RegistroHoras{{Fecha: dia(23), Horas: 1}, {Fecha: dia(25), Horas: 4}}, 5},
		{"la hora del registro no lo saca del día de fin", []RegistroHoras{{Fecha: dia(25).Add(20 * time.Hour), Horas: 3}}, 3},
		{"antes y después no se suman", []RegistroHoras{{Fecha: dia(22), Horas: 8}, {Fecha: dia(24), Horas: 2}, {Fecha: dia(26), Horas: 5}}, 2},
		{"ninguno dentro del período", []RegistroHoras{{Fecha: dia(21), Horas: 3}}, 0},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			got := HorasEnPeriodo(c.registros, inicio, fin)
			if math.Abs(got-c.quiere) > 1e-9 {
				t.Errorf("HorasEnPeriodo() = %v, se esperaba %v", got, c.quiere)
			}
		})
	}
}
