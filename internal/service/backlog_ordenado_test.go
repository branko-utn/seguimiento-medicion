package service

import (
	"errors"
	"reflect"
	"testing"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/project"
	"github.com/branko-utn/seguimiento-medicion/internal/store/memory"
)

// backlogOrdenado arma el servicio con dos proyectos:
//   - "Demo" (P1): HU-1 Baja, HU-2 Alta, HU-3 Media y HU-4 Alta, como en los Antecedentes de HU-09;
//   - "Otro" (P2): HU-5 Alta, para verificar que no se mezcla con Demo (RN-09.3).
//
// Devuelve también el repositorio, para dejar historias en otro estado sin pasar por HU-08.
func backlogOrdenado(t *testing.T) (*Backlog, *memory.Historias) {
	t.Helper()
	repoProyectos, _ := memory.NuevoProyectos("")
	repoHistorias, _ := memory.NuevoHistorias("")
	for _, nombre := range []string{"Demo", "Otro"} {
		if _, err := NuevoProyectos(repoProyectos).Crear(nombre, "", inicio, fin); err != nil {
			t.Fatal(err)
		}
	}
	s := NuevoBacklog(repoHistorias, repoProyectos, nil, nil, nil)
	for _, h := range []struct {
		proyecto  domain.ProyectoID
		prioridad backlog.Prioridad
	}{
		{"P1", backlog.PrioridadBaja}, {"P1", backlog.PrioridadAlta}, {"P1", backlog.PrioridadMedia},
		{"P1", backlog.PrioridadAlta}, {"P2", backlog.PrioridadAlta},
	} {
		if _, err := s.Crear(h.proyecto, "Historia", "", h.prioridad); err != nil {
			t.Fatal(err)
		}
	}
	return s, repoHistorias
}

// idsHistorias devuelve los IDs de las historias, en orden.
func idsHistorias(hs []backlog.Historia) []domain.HistoriaID {
	ids := []domain.HistoriaID{}
	for _, h := range hs {
		ids = append(ids, h.ID)
	}
	return ids
}

func TestBacklog_DelProyecto(t *testing.T) {
	pendiente, completada := backlog.Pendiente, backlog.Completada
	casos := []struct {
		nombre string
		estado *backlog.Estado
		want   []domain.HistoriaID
	}{
		// RN-09.1 y RN-09.3: ordenado por prioridad e ID, sin la historia de Otro.
		{"sin filtro", nil, []domain.HistoriaID{"HU-2", "HU-4", "HU-3", "HU-1"}},
		// RN-09.2: solo las del estado pedido, en el mismo orden.
		{"pendientes", &pendiente, []domain.HistoriaID{"HU-2", "HU-4", "HU-1"}},
		// Ninguna historia en el estado pedido: lista vacía, sin error.
		{"ninguna completada", &completada, []domain.HistoriaID{}},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			s, repo := backlogOrdenado(t)
			prepararHistoria(t, repo, "HU-3", backlog.EnProgreso, "S1")

			got, err := s.DelProyecto("P1", c.estado)
			if err != nil {
				t.Fatalf("error inesperado: %v", err)
			}
			if ids := idsHistorias(got); !reflect.DeepEqual(ids, c.want) {
				t.Errorf("backlog = %v, se esperaba %v", ids, c.want)
			}
		})
	}
}

func TestBacklog_DelProyecto_sinHistorias(t *testing.T) {
	repoProyectos, _ := memory.NuevoProyectos("")
	repoHistorias, _ := memory.NuevoHistorias("")
	if _, err := NuevoProyectos(repoProyectos).Crear("Vacío", "", inicio, fin); err != nil {
		t.Fatal(err)
	}
	s := NuevoBacklog(repoHistorias, repoProyectos, nil, nil, nil)

	// Backlog vacío: lista vacía, sin error.
	got, err := s.DelProyecto("P1", nil)
	if err != nil || len(got) != 0 {
		t.Errorf("DelProyecto = %v, %v; se esperaba una lista vacía sin error", got, err)
	}
}

func TestBacklog_DelProyecto_errores(t *testing.T) {
	invalido := backlog.Estado(99)
	casos := []struct {
		nombre   string
		proyecto domain.ProyectoID
		estado   *backlog.Estado
		want     error
	}{
		{"proyecto inexistente", "P9", nil, project.ErrNoExiste},
		{"estado invalido", "P1", &invalido, backlog.ErrEstadoInvalido},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			s, _ := backlogOrdenado(t)

			if _, err := s.DelProyecto(c.proyecto, c.estado); !errors.Is(err, c.want) {
				t.Errorf("error = %v, se esperaba %v", err, c.want)
			}
		})
	}
}
