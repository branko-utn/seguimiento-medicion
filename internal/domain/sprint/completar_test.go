package sprint

import (
	"errors"
	"testing"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
)

// sprintConHistorias arma un sprint activo con HU-1 y HU-2 asignadas.
func sprintConHistorias() Sprint {
	return Sprint{ID: "S1", ProyectoID: "P1", Estado: Activo, Historias: []domain.HistoriaID{"HU-1", "HU-2"}}
}

func TestCompletar(t *testing.T) {
	s := sprintConHistorias()

	if err := s.Completar("HU-1"); err != nil {
		t.Fatalf("Completar(HU-1): %v", err)
	}
	if err := s.Completar("HU-2"); err != nil {
		t.Fatalf("Completar(HU-2): %v", err)
	}
	if len(s.Completadas) != 2 || s.Completadas[0] != "HU-1" || s.Completadas[1] != "HU-2" {
		t.Errorf("Completadas = %v, se esperaba [HU-1 HU-2] en orden", s.Completadas)
	}
}

func TestCompletar_sprintPlanificado(t *testing.T) {
	// La spec solo prohíbe completar en un sprint cerrado (RN-12.2).
	s := sprintConHistorias()
	s.Estado = Planificado
	if err := s.Completar("HU-1"); err != nil {
		t.Errorf("en un sprint planificado se puede completar: %v", err)
	}
}

func TestCompletar_errores(t *testing.T) {
	casos := []struct {
		nombre   string
		preparar func(*Sprint)
		historia domain.HistoriaID
		quiere   error
		mensaje  string
	}{
		{"historia no asignada (RN-12.1)", nil, "HU-3", ErrHistoriaNoAsignada, "la historia no está asignada a este sprint"},
		{"sprint cerrado (RN-12.2)", func(s *Sprint) { s.Estado = Cerrado }, "HU-1", ErrSprintCerrado, "el sprint está cerrado y no se puede modificar"},
		{"ya completada (RN-12.3)", func(s *Sprint) { s.Completadas = []domain.HistoriaID{"HU-1"} }, "HU-1", ErrHistoriaCompletada, "la historia ya está completada"},
		{"cerrado y no asignada: primero el sprint", func(s *Sprint) { s.Estado = Cerrado }, "HU-3", ErrSprintCerrado, "el sprint está cerrado y no se puede modificar"},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			s := sprintConHistorias()
			if c.preparar != nil {
				c.preparar(&s)
			}
			antes := len(s.Completadas)

			err := s.Completar(c.historia)

			if !errors.Is(err, c.quiere) {
				t.Fatalf("error = %v, se esperaba %v", err, c.quiere)
			}
			if err.Error() != c.mensaje {
				t.Errorf("mensaje = %q, se esperaba %q", err.Error(), c.mensaje)
			}
			if len(s.Completadas) != antes {
				t.Errorf("con error no debe cambiar Completadas: %v", s.Completadas)
			}
		})
	}
}
