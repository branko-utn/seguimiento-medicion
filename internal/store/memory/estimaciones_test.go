package memory

import (
	"testing"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/effort"
)

func TestEstimacionesGuardarBuscar(t *testing.T) {
	casos := []struct {
		titulo      string
		guardar     []effort.EstimacionHoras
		buscar      domain.HistoriaID
		quiereOk    bool
		quiereHoras float64
	}{
		{
			titulo:      "guarda una estimacion",
			guardar:     []effort.EstimacionHoras{{HistoriaID: domain.HistoriaID("HU-1"), Horas: 4}},
			buscar:      domain.HistoriaID("HU-1"),
			quiereOk:    true,
			quiereHoras: 4,
		},
		{
			titulo:      "reemplaza la estimacion de la misma historia",
			guardar:     []effort.EstimacionHoras{{HistoriaID: domain.HistoriaID("HU-1"), Horas: 4}, {HistoriaID: domain.HistoriaID("HU-1"), Horas: 6}},
			buscar:      domain.HistoriaID("HU-1"),
			quiereOk:    true,
			quiereHoras: 6,
		},
		{
			titulo:      "historias distintas conviven",
			guardar:     []effort.EstimacionHoras{{HistoriaID: domain.HistoriaID("HU-1"), Horas: 4}, {HistoriaID: domain.HistoriaID("HU-2"), Horas: 1.5}},
			buscar:      domain.HistoriaID("HU-2"),
			quiereOk:    true,
			quiereHoras: 1.5,
		},
		{
			titulo:   "buscar inexistente devuelve ok falso",
			guardar:  []effort.EstimacionHoras{{HistoriaID: domain.HistoriaID("HU-1"), Horas: 4}},
			buscar:   domain.HistoriaID("HU-99"),
			quiereOk: false,
		},
	}

	for _, c := range casos {
		t.Run(c.titulo, func(t *testing.T) {
			r, err := NuevoEstimaciones("")
			if err != nil {
				t.Fatalf("no se esperaba error: %v", err)
			}
			for _, e := range c.guardar {
				if err := r.Guardar(e); err != nil {
					t.Fatalf("no se esperaba error al guardar: %v", err)
				}
			}
			got, ok := r.Buscar(c.buscar)
			if ok != c.quiereOk {
				t.Fatalf("ok esperado %v, obtenido %v", c.quiereOk, ok)
			}
			if !c.quiereOk {
				return
			}
			if got.HistoriaID != c.buscar {
				t.Fatalf("HistoriaID esperado %s, obtenido %s", c.buscar, got.HistoriaID)
			}
			if got.Horas != c.quiereHoras {
				t.Fatalf("Horas esperadas %v, obtenidas %v", c.quiereHoras, got.Horas)
			}
		})
	}
}
