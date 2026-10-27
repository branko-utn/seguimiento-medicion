package service

import (
	"errors"
	"testing"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/sprint"
)

func TestSprints_CerrarSprint(t *testing.T) {
	s, repoSprints, repoHistorias := prepararCompletar(t) // S1 con HU-1 (3 SP) y HU-2 (5 SP)
	if err := s.Iniciar("S1"); err != nil {
		t.Fatal(err)
	}
	if err := s.CompletarHistoria("S1", "HU-1"); err != nil {
		t.Fatal(err)
	}
	h2, _ := repoHistorias.Buscar("HU-2")
	h2.Estado = backlog.EnProgreso
	if err := repoHistorias.Guardar(h2); err != nil {
		t.Fatal(err)
	}

	if err := s.CerrarSprint("S1"); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}

	sp, _ := repoSprints.Buscar("S1")
	if sp.Estado != sprint.Cerrado || sp.SPPlanificados != 8 || sp.SPCompletados != 3 {
		t.Errorf("sprint = %s, %d planificados, %d completados; se esperaba Cerrado, 8 y 3", sp.Estado, sp.SPPlanificados, sp.SPCompletados)
	}
	if h, _ := repoHistorias.Buscar("HU-2"); h.Estado != backlog.Pendiente || h.SprintID != "" {
		t.Errorf("HU-2 = %s en %q, se esperaba Pendiente y sin sprint (RN-13.3)", h.Estado, h.SprintID)
	}
	if h, _ := repoHistorias.Buscar("HU-1"); h.Estado != backlog.Completada || h.SprintID != "S1" {
		t.Errorf("HU-1 = %s en %q, se esperaba Completada en S1", h.Estado, h.SprintID)
	}
}

func TestSprints_CerrarSprint_sinHistorias(t *testing.T) {
	s, repoSprints, _ := armarServicioSprints(t)
	if _, err := s.Crear("P1", "Sprint", octubre(23), octubre(25)); err != nil {
		t.Fatal(err)
	}
	if err := s.Iniciar("S1"); err != nil {
		t.Fatal(err)
	}
	if err := s.CerrarSprint("S1"); err != nil {
		t.Fatal(err)
	}
	if sp, _ := repoSprints.Buscar("S1"); sp.Estado != sprint.Cerrado || sp.SPPlanificados != 0 || sp.SPCompletados != 0 {
		t.Errorf("sprint = %+v, se esperaba Cerrado con 0 y 0", sp)
	}
}

func TestSprints_CerrarSprint_errores(t *testing.T) {
	casos := []struct {
		titulo   string
		iniciar  bool
		cerrar   bool
		sprintID string
		quiere   error
	}{
		{"sprint inexistente", true, false, "S9", sprint.ErrNoExiste},
		{"sprint que no empezó", false, false, "S1", sprint.ErrNoActivo},
		{"sprint ya cerrado", true, true, "S1", sprint.ErrSprintCerrado},
	}
	for _, c := range casos {
		t.Run(c.titulo, func(t *testing.T) {
			s, _, repoHistorias := prepararCompletar(t)
			if c.iniciar {
				if err := s.Iniciar("S1"); err != nil {
					t.Fatal(err)
				}
			}
			if c.cerrar {
				if err := s.CerrarSprint("S1"); err != nil {
					t.Fatal(err)
				}
				// Después del primer cierre HU-1 y HU-2 volvieron al backlog: se reasigna una para ver que no cambie.
				h, _ := repoHistorias.Buscar("HU-1")
				h.SprintID = "S1"
				if err := repoHistorias.Guardar(h); err != nil {
					t.Fatal(err)
				}
			}

			err := s.CerrarSprint(domain.SprintID(c.sprintID))

			if !errors.Is(err, c.quiere) {
				t.Fatalf("error = %v, se esperaba %v", err, c.quiere)
			}
			// Si falla, ninguna historia vuelve al backlog.
			if h, _ := repoHistorias.Buscar("HU-1"); h.SprintID != "S1" {
				t.Errorf("con error, HU-1 no debería volver al backlog (sprint %q)", h.SprintID)
			}
		})
	}
}
