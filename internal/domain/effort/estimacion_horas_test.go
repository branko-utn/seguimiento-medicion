package effort

import (
	"errors"
	"testing"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
)

func TestNuevaEstimacion(t *testing.T) {
	casos := []struct {
		titulo      string
		historiaID  domain.HistoriaID
		horas       float64
		quiereError error
		quiereHoras float64
	}{
		{
			titulo:      "crea una estimacion valida",
			historiaID:  domain.HistoriaID("HU-1"),
			horas:       4,
			quiereHoras: 4,
		},
		{
			titulo:      "acepta fracciones de hora",
			historiaID:  domain.HistoriaID("HU-1"),
			horas:       0.5,
			quiereHoras: 0.5,
		},
		{
			titulo:      "acepta 1.5 horas",
			historiaID:  domain.HistoriaID("HU-2"),
			horas:       1.5,
			quiereHoras: 1.5,
		},
		{
			titulo:      "rechaza horas iguales a cero",
			historiaID:  domain.HistoriaID("HU-1"),
			horas:       0,
			quiereError: ErrHorasInvalidas,
		},
		{
			titulo:      "rechaza horas negativas",
			historiaID:  domain.HistoriaID("HU-1"),
			horas:       -2,
			quiereError: ErrHorasInvalidas,
		},
		{
			titulo:      "rechaza horas negativas pequenas",
			historiaID:  domain.HistoriaID("HU-1"),
			horas:       -0.1,
			quiereError: ErrHorasInvalidas,
		},
	}

	for _, c := range casos {
		t.Run(c.titulo, func(t *testing.T) {
			est, err := NuevaEstimacion(c.historiaID, c.horas)
			if c.quiereError != nil {
				if !errors.Is(err, c.quiereError) {
					t.Fatalf("se esperaba error %v, se obtuvo %v", c.quiereError, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("no se esperaba error: %v", err)
			}
			if est.HistoriaID != c.historiaID {
				t.Fatalf("HistoriaID esperado %s, obtenido %s", c.historiaID, est.HistoriaID)
			}
			if est.Horas != c.quiereHoras {
				t.Fatalf("Horas esperadas %v, obtenidas %v", c.quiereHoras, est.Horas)
			}
		})
	}
}
