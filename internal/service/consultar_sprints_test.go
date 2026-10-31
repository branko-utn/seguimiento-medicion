package service

import (
	"errors"
	"testing"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/project"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/sprint"
)

func TestSprints_Cerrados(t *testing.T) {
	s, _, repoHistorias := armarServicioSprints(t)
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
	for _, sp := range []struct {
		goal        string
		inicio, fin int
	}{
		{"Primer MVP", 23, 25},
		{"Segundo sprint", 26, 27},
		{"Tercer sprint", 28, 29},
	} {
		if _, err := s.Crear("P1", sp.goal, octubre(sp.inicio), octubre(sp.fin)); err != nil {
			t.Fatal(err)
		}
	}

	// S2 se cierra antes que S1: el listado igual tiene que ordenarse por fecha de inicio (RN-14.2).
	for _, id := range []domain.HistoriaID{"HU-2", "HU-3"} {
		if err := s.Asignar("S2", id); err != nil {
			t.Fatal(err)
		}
	}
	for _, paso := range []func() error{
		func() error { return s.Iniciar("S2") },
		func() error { return s.CompletarHistoria("S2", "HU-2") },
		func() error { return s.CerrarSprint("S2") },
		func() error { return s.Asignar("S1", "HU-1") },
		func() error { return s.Iniciar("S1") },
		func() error { return s.CompletarHistoria("S1", "HU-1") },
		func() error { return s.CerrarSprint("S1") },
		func() error { return s.Asignar("S3", "HU-3") },
		func() error { return s.Iniciar("S3") },
	} {
		if err := paso(); err != nil {
			t.Fatal(err)
		}
	}
	// RN-14.3: los SP del listado son los fijados al cerrar, aunque la historia cambie después.
	h, _ := repoHistorias.Buscar("HU-2")
	h.StoryPoints = 100
	if err := repoHistorias.Guardar(h); err != nil {
		t.Fatal(err)
	}

	cerrados, err := s.Cerrados("P1")
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(cerrados) != 2 {
		t.Fatalf("Cerrados(P1) = %+v, se esperaban 2 sprints", cerrados)
	}
	if cerrados[0].ID != "S1" || cerrados[0].Goal != "Primer MVP" || !cerrados[0].Inicio.Equal(octubre(23)) ||
		cerrados[0].SPPlanificados != 3 || cerrados[0].SPCompletados != 3 {
		t.Errorf("primer sprint = %+v, se esperaba S1 · Primer MVP · 3/3", cerrados[0])
	}
	if cerrados[1].ID != "S2" || cerrados[1].Goal != "Segundo sprint" || cerrados[1].SPPlanificados != 7 ||
		cerrados[1].SPCompletados != 5 {
		t.Errorf("segundo sprint = %+v, se esperaba S2 · Segundo sprint · 7/5", cerrados[1])
	}

	// Un proyecto sin sprints cerrados devuelve un listado vacío, no un error.
	if lista, err := s.Cerrados("P2"); err != nil || len(lista) != 0 {
		t.Errorf("Cerrados(P2) = %+v (error %v), se esperaba vacío", lista, err)
	}
}

func TestSprints_Cerrados_sinSprintsCerrados(t *testing.T) {
	s, _ := nuevoServicioSprints(t)
	if _, err := s.Crear("P1", "Activo", octubre(23), octubre(25)); err != nil {
		t.Fatal(err)
	}
	if err := s.Iniciar("S1"); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Crear("P1", "Planificado", octubre(26), octubre(27)); err != nil {
		t.Fatal(err)
	}

	cerrados, err := s.Cerrados("P1")
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(cerrados) != 0 {
		t.Errorf("Cerrados(P1) = %+v, se esperaba vacío: solo se listan los cerrados (RN-14.1)", cerrados)
	}
}

func TestSprints_Cerrados_proyectoInexistente(t *testing.T) {
	s, _ := nuevoServicioSprints(t)

	if _, err := s.Cerrados("P9"); !errors.Is(err, project.ErrNoExiste) {
		t.Errorf("error = %v, se esperaba %v", err, project.ErrNoExiste)
	}
}

func TestSprints_Detalle(t *testing.T) {
	s, _, repoHistorias := prepararCompletar(t) // S1 con HU-1 (3 SP) y HU-2 (5 SP)
	if err := s.Iniciar("S1"); err != nil {
		t.Fatal(err)
	}
	if err := s.CompletarHistoria("S1", "HU-1"); err != nil {
		t.Fatal(err)
	}
	if err := s.CerrarSprint("S1"); err != nil {
		t.Fatal(err)
	}

	detalle, err := s.Detalle("S1")
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if detalle.Sprint.ID != "S1" || detalle.Sprint.Estado != sprint.Cerrado ||
		detalle.Sprint.SPPlanificados != 8 || detalle.Sprint.SPCompletados != 3 {
		t.Errorf("sprint del detalle = %+v, se esperaba S1 cerrado con 8/3 SP", detalle.Sprint)
	}
	quiere := []HistoriaDetalle{
		{ID: "HU-1", Titulo: "Crear proyecto", StoryPoints: 3, Completada: true},
		// RN-14.4: HU-2 volvió al backlog al cerrar, pero el detalle la muestra como no completada.
		{ID: "HU-2", Titulo: "Crear sprint", StoryPoints: 5, Completada: false},
	}
	if len(detalle.Historias) != len(quiere) {
		t.Fatalf("el detalle tiene %d historias, se esperaban %d: %+v", len(detalle.Historias), len(quiere), detalle.Historias)
	}
	for i, h := range quiere {
		if detalle.Historias[i] != h {
			t.Errorf("historia %d = %+v, se esperaba %+v", i+1, detalle.Historias[i], h)
		}
	}
	if h, _ := repoHistorias.Buscar("HU-2"); h.SprintID != "" {
		t.Errorf("HU-2 quedó asignada al sprint %q: el detalle se arma con Sprint.Historias", h.SprintID)
	}
}

func TestSprints_Detalle_sinHistorias(t *testing.T) {
	s, _ := nuevoServicioSprints(t)
	if _, err := s.Crear("P1", "Sprint", octubre(23), octubre(25)); err != nil {
		t.Fatal(err)
	}
	if err := s.Iniciar("S1"); err != nil {
		t.Fatal(err)
	}
	if err := s.CerrarSprint("S1"); err != nil {
		t.Fatal(err)
	}

	detalle, err := s.Detalle("S1")
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(detalle.Historias) != 0 || detalle.Sprint.SPPlanificados != 0 || detalle.Sprint.SPCompletados != 0 {
		t.Errorf("detalle = %+v, se esperaba sin historias y con 0/0 SP", detalle)
	}
}

func TestSprints_Detalle_errores(t *testing.T) {
	casos := []struct {
		titulo   string
		iniciar  bool
		sprintID domain.SprintID
		quiere   error
	}{
		{"sprint inexistente", false, "S9", sprint.ErrNoExiste},
		{"sprint planificado", false, "S1", sprint.ErrNoCerrado},
		{"sprint activo", true, "S1", sprint.ErrNoCerrado},
	}
	for _, c := range casos {
		t.Run(c.titulo, func(t *testing.T) {
			s, _ := nuevoServicioSprints(t)
			if _, err := s.Crear("P1", "Sprint", octubre(23), octubre(25)); err != nil {
				t.Fatal(err)
			}
			if c.iniciar {
				if err := s.Iniciar("S1"); err != nil {
					t.Fatal(err)
				}
			}

			_, err := s.Detalle(c.sprintID)

			if !errors.Is(err, c.quiere) {
				t.Errorf("error = %v, se esperaba %v", err, c.quiere)
			}
		})
	}
}
