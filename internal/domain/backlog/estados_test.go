package backlog

import (
	"errors"
	"testing"
)

func TestCambiarEstado(t *testing.T) {
	casos := []struct {
		nombre    string
		estado    Estado // estado inicial
		conSprint bool   // si la historia está asignada a un sprint
		nuevo     Estado // estado pedido
		quiere    Estado // estado final esperado si no hay error
		quiereErr error
		quiereMsg string // mensaje exacto que verá el usuario (lo compara la BDD)
	}{
		// CA-08.1: Pendiente → En progreso → Completada.
		{"empezar una historia del sprint", Pendiente, true, EnProgreso, EnProgreso, nil, ""},
		{"volver a pendiente", EnProgreso, true, Pendiente, Pendiente, nil, ""},
		// Caso límite: En progreso → Pendiente no exige sprint (RN-08.2, sin condición).
		{"volver a pendiente sin sprint", EnProgreso, false, Pendiente, Pendiente, nil, ""},
		{"completar desde pendiente en un sprint", Pendiente, true, Completada, Completada, nil, ""},
		{"completar desde en progreso", EnProgreso, true, Completada, Completada, nil, ""},
		// RN-08.2: En progreso y Completada exigen estar en un sprint.
		{"empezar sin sprint", Pendiente, false, EnProgreso, Pendiente, ErrSinSprint, "la historia debe estar asignada a un sprint"},
		{"completar sin sprint", Pendiente, false, Completada, Pendiente, ErrSinSprint, "la historia debe estar asignada a un sprint"},
		// RN-08.4: cambiar al mismo estado es un error, antes que el requisito de sprint.
		{"mismo estado", Pendiente, false, Pendiente, Pendiente, ErrTransicionInvalida, "no se puede pasar de Pendiente a Pendiente"},
		{"mismo estado en progreso", EnProgreso, true, EnProgreso, EnProgreso, ErrTransicionInvalida, "no se puede pasar de En progreso a En progreso"},
		// RN-08.3: una historia completada no cambia de estado.
		{"una historia completada no cambia", Completada, true, EnProgreso, Completada, ErrTransicionInvalida, "no se puede pasar de Completada a En progreso"},
		// RN-08.1: el estado nuevo debe ser uno de los tres.
		{"estado inválido", Pendiente, false, Estado(99), Pendiente, ErrEstadoInvalido, "el estado debe ser Pendiente, En progreso o Completada"},
	}
	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			h, _ := Nueva("P1", "Login", "", PrioridadMedia)
			h.Estado = caso.estado
			if caso.conSprint {
				h.SprintID = "S1"
			}

			err := h.CambiarEstado(caso.nuevo)
			if !errors.Is(err, caso.quiereErr) {
				t.Fatalf("se esperaba el error %v, se obtuvo %v", caso.quiereErr, err)
			}
			if caso.quiereMsg != "" && err.Error() != caso.quiereMsg {
				t.Errorf("mensaje = %q, se esperaba %q", err.Error(), caso.quiereMsg)
			}
			// Con error, el estado no cambia; sin error, queda en el pedido.
			quiere := caso.quiere
			if caso.quiereErr != nil {
				quiere = caso.estado
			}
			if h.Estado != quiere {
				t.Errorf("estado = %v, se esperaba %v", h.Estado, quiere)
			}
		})
	}
}

func TestCompletar(t *testing.T) {
	t.Run("desde pendiente en un sprint", func(t *testing.T) {
		h, _ := Nueva("P1", "Login", "", PrioridadMedia)
		h.SprintID = "S1"

		if err := h.Completar(); err != nil {
			t.Fatalf("error inesperado: %v", err)
		}
		if h.Estado != Completada {
			t.Errorf("estado = %v, se esperaba Completada", h.Estado)
		}
	})

	t.Run("desde en progreso", func(t *testing.T) {
		h, _ := Nueva("P1", "Login", "", PrioridadMedia)
		h.SprintID = "S1"
		h.Estado = EnProgreso

		if err := h.Completar(); err != nil {
			t.Fatalf("error inesperado: %v", err)
		}
		if h.Estado != Completada {
			t.Errorf("estado = %v, se esperaba Completada", h.Estado)
		}
	})

	t.Run("sin sprint", func(t *testing.T) {
		h, _ := Nueva("P1", "Login", "", PrioridadMedia)

		if err := h.Completar(); !errors.Is(err, ErrSinSprint) {
			t.Fatalf("se esperaba ErrSinSprint, se obtuvo %v", err)
		}
		if h.Estado != Pendiente {
			t.Errorf("estado = %v, no debería cambiar", h.Estado)
		}
	})

	t.Run("ya completada", func(t *testing.T) {
		// La misma regla de CambiarEstado (spec HU-08): HU-12 valida este caso antes en el sprint, pero el
		// dominio debe cumplirlo por sí solo.
		h, _ := Nueva("P1", "Login", "", PrioridadMedia)
		h.SprintID = "S1"
		h.Estado = Completada

		err := h.Completar()
		if !errors.Is(err, ErrTransicionInvalida) {
			t.Fatalf("se esperaba ErrTransicionInvalida, se obtuvo %v", err)
		}
		if err.Error() != "no se puede pasar de Completada a Completada" {
			t.Errorf("mensaje = %q", err.Error())
		}
		if h.Estado != Completada {
			t.Errorf("estado = %v, no debería cambiar", h.Estado)
		}
	})
}
