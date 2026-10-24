package sprint

import (
	"errors"
	"slices"
	"testing"
	"time"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/project"
)

func dia(d int) time.Time {
	// Días de octubre de 2026; 32 = 1 de noviembre, 33 = 2 de noviembre (time.Date normaliza).
	return time.Date(2026, time.October, d, 0, 0, 0, 0, time.Local)
}

func TestNuevo(t *testing.T) {
	casos := []struct {
		titulo      string
		goal        string
		inicio, fin time.Time
		quiereError error
		quiereGoal  string
	}{
		{titulo: "crea un sprint válido", goal: "Primer MVP", inicio: dia(23), fin: dia(25), quiereGoal: "Primer MVP"},
		{titulo: "quita los espacios sobrantes del Sprint Goal", goal: "  Primer MVP  ", inicio: dia(23), fin: dia(25), quiereGoal: "Primer MVP"},
		{titulo: "acepta un sprint de un solo día", goal: "Demo", inicio: dia(23), fin: dia(23), quiereGoal: "Demo"},
		{titulo: "rechaza un Sprint Goal vacío", goal: "", inicio: dia(23), fin: dia(25), quiereError: ErrGoalVacio},
		{titulo: "rechaza un Sprint Goal con solo espacios", goal: "  ", inicio: dia(23), fin: dia(25), quiereError: ErrGoalVacio},
		{titulo: "rechaza un fin anterior al inicio", goal: "Primer MVP", inicio: dia(25), fin: dia(23), quiereError: ErrFechasInvalidas},
	}

	for _, caso := range casos {
		t.Run(caso.titulo, func(t *testing.T) {
			s, err := Nuevo("P1", caso.goal, caso.inicio, caso.fin)

			if caso.quiereError != nil {
				if !errors.Is(err, caso.quiereError) {
					t.Fatalf("se esperaba el error %v, se obtuvo %v", caso.quiereError, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("error inesperado: %v", err)
			}
			if s.ProyectoID != "P1" || s.Goal != caso.quiereGoal {
				t.Errorf("Nuevo() = %+v, se esperaba proyecto P1 y goal %q", s, caso.quiereGoal)
			}
			if !s.Inicio.Equal(caso.inicio) || !s.Fin.Equal(caso.fin) {
				t.Errorf("fechas = %v al %v, se esperaba %v al %v", s.Inicio, s.Fin, caso.inicio, caso.fin)
			}
			if s.Estado != Planificado {
				t.Errorf("estado = %v, se esperaba %v", s.Estado, Planificado)
			}
			if len(s.Historias) != 0 || len(s.Completadas) != 0 {
				t.Errorf("un sprint nuevo no tiene historias: %+v", s)
			}
		})
	}
}

func TestValidarFechasEnProyecto(t *testing.T) {
	p := project.Proyecto{ID: "P1", Inicio: dia(21), Fin: dia(33)} // 21/10 al 02/11

	casos := []struct {
		titulo      string
		inicio, fin time.Time
		quiereError error
	}{
		{titulo: "dentro del proyecto", inicio: dia(23), fin: dia(25)},
		{titulo: "ocupa todo el proyecto", inicio: dia(21), fin: dia(33)},
		{titulo: "empieza antes del proyecto", inicio: dia(20), fin: dia(25), quiereError: ErrFueraDelProyecto},
		{titulo: "termina después del proyecto", inicio: dia(30), fin: dia(34), quiereError: ErrFueraDelProyecto},
	}

	for _, caso := range casos {
		t.Run(caso.titulo, func(t *testing.T) {
			s := Sprint{ProyectoID: "P1", Goal: "Goal", Inicio: caso.inicio, Fin: caso.fin, Estado: Planificado}

			err := s.ValidarFechasEnProyecto(p)
			if !errors.Is(err, caso.quiereError) {
				t.Errorf("se esperaba el error %v, se obtuvo %v", caso.quiereError, err)
			}
		})
	}
}

func TestValidarSinSuperposicion(t *testing.T) {
	s1 := Sprint{ID: "S1", Inicio: dia(23), Fin: dia(25)}

	casos := []struct {
		titulo      string
		inicio, fin time.Time
		otros       []Sprint
		quiereError error
	}{
		{titulo: "sin otros sprints", inicio: dia(23), fin: dia(25)},
		{titulo: "empieza el día siguiente al fin del otro", inicio: dia(26), fin: dia(27), otros: []Sprint{s1}},
		{titulo: "termina el día anterior al inicio del otro", inicio: dia(21), fin: dia(22), otros: []Sprint{s1}},
		{titulo: "empieza el mismo día que termina el otro", inicio: dia(25), fin: dia(27), otros: []Sprint{s1}, quiereError: ErrSuperpuesto},
		{titulo: "termina el mismo día que empieza el otro", inicio: dia(21), fin: dia(23), otros: []Sprint{s1}, quiereError: ErrSuperpuesto},
		{titulo: "queda dentro del otro", inicio: dia(24), fin: dia(24), otros: []Sprint{s1}, quiereError: ErrSuperpuesto},
		{titulo: "contiene al otro", inicio: dia(22), fin: dia(26), otros: []Sprint{s1}, quiereError: ErrSuperpuesto},
		{
			titulo: "se superpone con el segundo de varios",
			inicio: dia(29), fin: dia(30),
			otros:       []Sprint{s1, {ID: "S2", Inicio: dia(28), Fin: dia(29)}},
			quiereError: ErrSuperpuesto,
		},
	}

	for _, caso := range casos {
		t.Run(caso.titulo, func(t *testing.T) {
			s := Sprint{Inicio: caso.inicio, Fin: caso.fin, Estado: Planificado}

			err := s.ValidarSinSuperposicion(caso.otros)
			if !errors.Is(err, caso.quiereError) {
				t.Errorf("se esperaba el error %v, se obtuvo %v", caso.quiereError, err)
			}
		})
	}
}

func TestIniciar(t *testing.T) {
	casos := []struct {
		titulo      string
		estado      Estado
		quiereError error
		quiere      Estado
	}{
		{titulo: "un sprint planificado pasa a activo", estado: Planificado, quiere: Activo},
		{titulo: "un sprint activo no se vuelve a iniciar", estado: Activo, quiereError: ErrNoPlanificado, quiere: Activo},
		{titulo: "un sprint cerrado no se inicia", estado: Cerrado, quiereError: ErrNoPlanificado, quiere: Cerrado},
	}

	for _, caso := range casos {
		t.Run(caso.titulo, func(t *testing.T) {
			s := Sprint{ID: "S1", Estado: caso.estado}

			err := s.Iniciar()
			if !errors.Is(err, caso.quiereError) {
				t.Fatalf("se esperaba el error %v, se obtuvo %v", caso.quiereError, err)
			}
			if s.Estado != caso.quiere {
				t.Errorf("estado = %v, se esperaba %v", s.Estado, caso.quiere)
			}
		})
	}
}

// historiaLista devuelve una historia de P1 que cumple todo para entrar a un sprint: estimada, con criterios y pendiente.
func historiaLista() backlog.Historia {
	return backlog.Historia{
		ID:          "HU-1",
		ProyectoID:  "P1",
		Titulo:      "Crear proyecto",
		Estado:      backlog.Pendiente,
		StoryPoints: 3,
		Criterios:   []backlog.Criterio{{Texto: "se guarda el proyecto"}},
	}
}

func TestAsignar(t *testing.T) {
	casos := []struct {
		titulo string
		estado Estado
		h      func(h *backlog.Historia)
	}{
		{titulo: "a un sprint planificado", estado: Planificado},
		{titulo: "a un sprint activo", estado: Activo},
		// RN-11.3 (ya asignada) la valida Historia.AsignarASprint, no el sprint.
		{titulo: "no controla si la historia ya tiene sprint", estado: Planificado, h: func(h *backlog.Historia) { h.SprintID = "S2" }},
	}

	for _, caso := range casos {
		t.Run(caso.titulo, func(t *testing.T) {
			s := Sprint{ID: "S1", ProyectoID: "P1", Estado: caso.estado, Historias: []domain.HistoriaID{"HU-9"}}
			h := historiaLista()
			if caso.h != nil {
				caso.h(&h)
			}

			if err := s.Asignar(h); err != nil {
				t.Fatalf("error inesperado: %v", err)
			}
			if !slices.Equal(s.Historias, []domain.HistoriaID{"HU-9", "HU-1"}) {
				t.Errorf("historias = %v, se esperaba [HU-9 HU-1]", s.Historias)
			}
		})
	}
}

func TestAsignar_errores(t *testing.T) {
	sinEstimar := func(h *backlog.Historia) { h.StoryPoints = 0 }
	sinCriterios := func(h *backlog.Historia) { h.Criterios = nil }
	completada := func(h *backlog.Historia) { h.Estado = backlog.Completada; h.SprintID = "S2" }
	deOtroProyecto := func(h *backlog.Historia) { h.ProyectoID = "P2" }

	casos := []struct {
		titulo      string
		estado      Estado
		cambios     []func(h *backlog.Historia)
		quiereError error
	}{
		{"historia sin estimar", Planificado, []func(*backlog.Historia){sinEstimar}, ErrHistoriaSinEstimar},
		{"historia con Story Points negativos", Planificado, []func(*backlog.Historia){func(h *backlog.Historia) { h.StoryPoints = -1 }}, ErrHistoriaSinEstimar},
		{"historia sin criterios", Planificado, []func(*backlog.Historia){sinCriterios}, ErrHistoriaSinCriterios},
		{"sprint cerrado", Cerrado, nil, ErrSprintCerrado},
		{"historia completada", Activo, []func(*backlog.Historia){completada}, ErrHistoriaCompletada},
		{"historia de otro proyecto", Planificado, []func(*backlog.Historia){deOtroProyecto}, ErrOtroProyecto},

		// Orden de validación: primero el sprint, después la historia.
		{"primero se controla que el sprint no esté cerrado", Cerrado,
			[]func(*backlog.Historia){deOtroProyecto, completada, sinEstimar, sinCriterios}, ErrSprintCerrado},
		{"después que la historia sea del proyecto", Planificado,
			[]func(*backlog.Historia){deOtroProyecto, completada, sinEstimar, sinCriterios}, ErrOtroProyecto},
		{"después que no esté completada", Planificado,
			[]func(*backlog.Historia){completada, sinEstimar, sinCriterios}, ErrHistoriaCompletada},
		{"después que esté estimada", Planificado,
			[]func(*backlog.Historia){sinEstimar, sinCriterios}, ErrHistoriaSinEstimar},
	}

	for _, caso := range casos {
		t.Run(caso.titulo, func(t *testing.T) {
			s := Sprint{ID: "S1", ProyectoID: "P1", Estado: caso.estado, Historias: []domain.HistoriaID{"HU-9"}}
			h := historiaLista()
			for _, cambiar := range caso.cambios {
				cambiar(&h)
			}

			err := s.Asignar(h)
			if !errors.Is(err, caso.quiereError) {
				t.Fatalf("se esperaba el error %v, se obtuvo %v", caso.quiereError, err)
			}
			if !slices.Equal(s.Historias, []domain.HistoriaID{"HU-9"}) {
				t.Errorf("historias = %v, se esperaba [HU-9]: un error no debe cambiar el sprint", s.Historias)
			}
		})
	}
}

// TestEstadoString verifica el texto de cada estado (ADR 0003: el cero no es válido).
func TestEstadoString(t *testing.T) {
	casos := []struct {
		estado Estado
		quiere string
	}{
		{Planificado, "Planificado"},
		{Activo, "Activo"},
		{Cerrado, "Cerrado"},
		{Estado(0), ""},
	}

	for _, caso := range casos {
		if got := caso.estado.String(); got != caso.quiere {
			t.Errorf("Estado(%d).String() = %q, se esperaba %q", int(caso.estado), got, caso.quiere)
		}
	}
}
