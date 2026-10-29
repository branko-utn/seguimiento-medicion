package project

import (
	"testing"
	"time"
)

// proyectoOctubre va del 21/10 al 02/11, como el proyecto de los escenarios de HU-04.
func proyectoOctubre() Proyecto {
	return Proyecto{
		Inicio: time.Date(2026, 10, 21, 0, 0, 0, 0, time.Local),
		Fin:    time.Date(2026, 11, 2, 0, 0, 0, 0, time.Local),
		Estado: Planificado,
	}
}

func TestProyecto_EstadoAl(t *testing.T) {
	casos := []struct {
		nombre string
		hoy    time.Time
		quiere Estado
	}{
		{"antes del inicio", time.Date(2026, 10, 20, 23, 59, 0, 0, time.Local), Planificado},
		{"el día de inicio", time.Date(2026, 10, 21, 0, 0, 0, 0, time.Local), EnCurso},
		{"el día de inicio con hora", time.Date(2026, 10, 21, 18, 30, 0, 0, time.Local), EnCurso},
		{"en el medio", time.Date(2026, 10, 27, 0, 0, 0, 0, time.Local), EnCurso},
		{"el día de fin con hora", time.Date(2026, 11, 2, 23, 0, 0, 0, time.Local), EnCurso},
		{"después del fin", time.Date(2026, 11, 3, 0, 0, 0, 0, time.Local), Finalizado},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			if got := proyectoOctubre().EstadoAl(c.hoy); got != c.quiere {
				t.Errorf("EstadoAl(%s) = %s, se esperaba %s", c.hoy.Format("02/01 15:04"), got, c.quiere)
			}
		})
	}
}

func TestProyecto_DiasRestantes(t *testing.T) {
	casos := []struct {
		nombre string
		hoy    time.Time
		quiere int
	}{
		{"antes del inicio", time.Date(2026, 10, 20, 0, 0, 0, 0, time.Local), 13},
		{"en el medio, sin contar hoy", time.Date(2026, 10, 27, 0, 0, 0, 0, time.Local), 6},
		{"la hora de hoy no cambia el resultado", time.Date(2026, 10, 27, 23, 59, 0, 0, time.Local), 6},
		{"el día anterior al fin", time.Date(2026, 11, 1, 0, 0, 0, 0, time.Local), 1},
		{"el día de fin", time.Date(2026, 11, 2, 12, 0, 0, 0, time.Local), 0},
		{"después del fin", time.Date(2026, 11, 10, 0, 0, 0, 0, time.Local), 0},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			if got := proyectoOctubre().DiasRestantes(c.hoy); got != c.quiere {
				t.Errorf("DiasRestantes(%s) = %d, se esperaba %d", c.hoy.Format("02/01 15:04"), got, c.quiere)
			}
		})
	}
}
