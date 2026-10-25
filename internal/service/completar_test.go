package service

import (
	"errors"
	"slices"
	"testing"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/sprint"
	"github.com/branko-utn/seguimiento-medicion/internal/store/memory"
)

// prepararCompletar arma un sprint S1 en Demo con HU-1 (3 SP) y HU-2 (5 SP) asignadas, y HU-3 sin asignar.
func prepararCompletar(t *testing.T) (*Sprints, *memory.Sprints, *memory.Historias) {
	t.Helper()
	s, repoSprints, repoHistorias := armarServicioSprints(t)
	if _, err := s.Crear("P1", "Sprint", octubre(23), octubre(25)); err != nil {
		t.Fatal(err)
	}
	criterio := []backlog.Criterio{{Texto: "funciona"}}
	for _, h := range []backlog.Historia{
		{ProyectoID: "P1", Titulo: "Crear proyecto", Estado: backlog.Pendiente, StoryPoints: 3, Criterios: criterio},
		{ProyectoID: "P1", Titulo: "Crear sprint", Estado: backlog.Pendiente, StoryPoints: 5, Criterios: criterio},
		{ProyectoID: "P1", Titulo: "Registrar defecto", Estado: backlog.Pendiente, StoryPoints: 2, Criterios: criterio},
	} {
		if _, err := repoHistorias.Agregar(h); err != nil {
			t.Fatal(err)
		}
	}
	for _, id := range []domain.HistoriaID{"HU-1", "HU-2"} {
		if err := s.Asignar("S1", id); err != nil {
			t.Fatal(err)
		}
	}
	return s, repoSprints, repoHistorias
}

func TestSprints_CompletarHistoria(t *testing.T) {
	s, repoSprints, repoHistorias := prepararCompletar(t)

	if err := s.CompletarHistoria("S1", "HU-1"); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	sp, _ := repoSprints.Buscar("S1")
	if !slices.Equal(sp.Completadas, []domain.HistoriaID{"HU-1"}) {
		t.Errorf("completadas de S1 = %v, se esperaba [HU-1]", sp.Completadas)
	}
	if h, _ := repoHistorias.Buscar("HU-1"); h.Estado != backlog.Completada {
		t.Errorf("estado de HU-1 = %s, se esperaba Completada", h.Estado)
	}
}

func TestSprints_CompletarHistoria_desdeEnProgreso(t *testing.T) {
	// RN-12.4: se completa desde Pendiente o desde En progreso.
	s, _, repoHistorias := prepararCompletar(t)
	h, _ := repoHistorias.Buscar("HU-2")
	h.Estado = backlog.EnProgreso
	if err := repoHistorias.Guardar(h); err != nil {
		t.Fatal(err)
	}
	if err := s.CompletarHistoria("S1", "HU-2"); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
}

func TestSprints_CompletarHistoria_errores(t *testing.T) {
	casos := []struct {
		titulo   string
		preparar func(t *testing.T, repoSprints *memory.Sprints, s *Sprints)
		sprint   domain.SprintID
		historia domain.HistoriaID
		quiere   error
	}{
		{"sprint inexistente", nil, "S9", "HU-1", sprint.ErrNoExiste},
		{"historia inexistente", nil, "S1", "HU-9", backlog.ErrNoExiste},
		{"historia no asignada al sprint", nil, "S1", "HU-3", sprint.ErrHistoriaNoAsignada},
		{"sprint cerrado", func(t *testing.T, r *memory.Sprints, _ *Sprints) {
			sp, _ := r.Buscar("S1")
			sp.Estado = sprint.Cerrado
			if err := r.Actualizar(sp); err != nil {
				t.Fatal(err)
			}
		}, "S1", "HU-1", sprint.ErrSprintCerrado},
		{"historia ya completada", func(t *testing.T, _ *memory.Sprints, s *Sprints) {
			if err := s.CompletarHistoria("S1", "HU-1"); err != nil {
				t.Fatal(err)
			}
		}, "S1", "HU-1", sprint.ErrHistoriaCompletada},
	}
	for _, c := range casos {
		t.Run(c.titulo, func(t *testing.T) {
			s, repoSprints, repoHistorias := prepararCompletar(t)
			if c.preparar != nil {
				c.preparar(t, repoSprints, s)
			}
			spAntes, _ := repoSprints.Buscar(c.sprint)
			hAntes, _ := repoHistorias.Buscar(c.historia)

			err := s.CompletarHistoria(c.sprint, c.historia)

			if !errors.Is(err, c.quiere) {
				t.Fatalf("error = %v, se esperaba %v", err, c.quiere)
			}
			// Si falla, no se modifica ninguna de las dos entidades.
			if spDespues, _ := repoSprints.Buscar(c.sprint); !slices.Equal(spDespues.Completadas, spAntes.Completadas) {
				t.Errorf("el sprint cambió: %v → %v", spAntes.Completadas, spDespues.Completadas)
			}
			if hDespues, _ := repoHistorias.Buscar(c.historia); hDespues.Estado != hAntes.Estado {
				t.Errorf("la historia cambió de estado: %s → %s", hAntes.Estado, hDespues.Estado)
			}
		})
	}
}
