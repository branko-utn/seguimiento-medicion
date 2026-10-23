package sprint

import (
	"errors"
	"testing"
	"time"

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
