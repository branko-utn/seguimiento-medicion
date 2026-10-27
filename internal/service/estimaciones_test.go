package service

import (
	"errors"
	"testing"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/effort"
	"github.com/branko-utn/seguimiento-medicion/internal/store/memory"
)

func nuevoServicioEstimaciones(t *testing.T) *Estimaciones {
	t.Helper()
	repoProy, err := memory.NuevoProyectos("")
	if err != nil {
		t.Fatal(err)
	}
	repoHist, err := memory.NuevoHistorias("")
	if err != nil {
		t.Fatal(err)
	}
	repoEst, err := memory.NuevoEstimaciones("")
	if err != nil {
		t.Fatal(err)
	}
	repoSpr, err := memory.NuevoSprints("")
	if err != nil {
		t.Fatal(err)
	}
	proy := NuevoProyectos(repoProy)
	sprints := NuevoSprints(repoSpr, repoProy, repoHist)
	hist := NuevoBacklog(repoHist, repoProy, sprints)
	if _, err := proy.Crear("Demo", "", inicio, fin); err != nil {
		t.Fatal(err)
	}
	if _, err := hist.Crear(domain.ProyectoID("P1"), "Crear proyecto", "", backlog.PrioridadAlta); err != nil {
		t.Fatal(err)
	}
	return NuevoEstimaciones(repoEst, repoHist)
}

func TestEstimaciones_Estimar(t *testing.T) {
	s := nuevoServicioEstimaciones(t)

	e, err := s.Estimar(domain.HistoriaID("HU-1"), 4)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if e.Horas != 4 {
		t.Fatalf("horas esperadas 4, obtenidas %v", e.Horas)
	}

	e2, err := s.Estimar(domain.HistoriaID("HU-1"), 6)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if e2.Horas != 6 {
		t.Fatalf("horas esperadas 6, obtenidas %v", e2.Horas)
	}
	got, ok, err := s.Buscar(domain.HistoriaID("HU-1"))
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if !ok || got.Horas != 6 {
		t.Fatalf("se esperaba estimacion 6, ok=%v got=%+v", ok, got)
	}
}

func TestEstimaciones_Estimar_errores(t *testing.T) {
	casos := []struct {
		titulo      string
		historiaID  domain.HistoriaID
		horas       float64
		quiereError error
	}{
		{"horas iguales a cero", domain.HistoriaID("HU-1"), 0, effort.ErrHorasInvalidas},
		{"horas negativas", domain.HistoriaID("HU-1"), -2, effort.ErrHorasInvalidas},
		{"historia inexistente", domain.HistoriaID("HU-99"), 4, backlog.ErrNoExiste},
	}

	for _, c := range casos {
		t.Run(c.titulo, func(t *testing.T) {
			s := nuevoServicioEstimaciones(t)
			_, err := s.Estimar(c.historiaID, c.horas)
			if !errors.Is(err, c.quiereError) {
				t.Fatalf("se esperaba error %v, se obtuvo %v", c.quiereError, err)
			}
		})
	}
}
