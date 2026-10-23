package backlog

import (
	"errors"
	"testing"
)

func TestNueva_validaTituloYPrioridad(t *testing.T) {
	casos := []struct {
		nombre      string
		titulo      string
		descripcion string
		prioridad   Prioridad
		wantTitulo  string
		wantErr     error
	}{
		{"válida con descripción", "Crear proyecto", "con detalle", PrioridadAlta, "Crear proyecto", nil},
		{"sin descripción", "Login", "", PrioridadMedia, "Login", nil},
		{"prioridad Baja", "Login", "", PrioridadBaja, "Login", nil},
		{"recorta los espacios del título", "  Login  ", "", PrioridadAlta, "Login", nil},
		{"título vacío", "", "", PrioridadAlta, "", ErrTituloVacio},
		{"título con solo espacios", "   ", "", PrioridadAlta, "", ErrTituloVacio},
		{"prioridad sin valor (cero)", "Login", "", Prioridad(0), "", ErrPrioridadInvalida},
		{"prioridad fuera de la lista", "Login", "", Prioridad(99), "", ErrPrioridadInvalida},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			h, err := Nueva("P1", c.titulo, c.descripcion, c.prioridad)
			if !errors.Is(err, c.wantErr) {
				t.Fatalf("error = %v, se esperaba %v", err, c.wantErr)
			}
			if c.wantErr != nil {
				if h.Titulo != "" {
					t.Errorf("con error no debería crearse la historia; se obtuvo %+v", h)
				}
				return
			}
			if h.Titulo != c.wantTitulo {
				t.Errorf("título = %q, se esperaba %q", h.Titulo, c.wantTitulo)
			}
			if h.Descripcion != c.descripcion {
				t.Errorf("descripción = %q, se esperaba %q", h.Descripcion, c.descripcion)
			}
			if h.Prioridad != c.prioridad {
				t.Errorf("prioridad = %v, se esperaba %v", h.Prioridad, c.prioridad)
			}
		})
	}
}

func TestNueva_estadoInicial(t *testing.T) {
	// RN-05.4: una historia nueva queda en estado Pendiente, sin Story Points, sin criterios y sin sprint.
	h, err := Nueva("P1", "Login", "", PrioridadAlta)
	if err != nil {
		t.Fatal(err)
	}
	if h.Estado != Pendiente {
		t.Errorf("estado = %v, se esperaba Pendiente", h.Estado)
	}
	if h.StoryPoints != 0 {
		t.Errorf("StoryPoints = %d, se esperaba 0", h.StoryPoints)
	}
	if len(h.Criterios) != 0 {
		t.Errorf("Criterios = %v, se esperaba vacío", h.Criterios)
	}
	if h.SprintID != "" {
		t.Errorf("SprintID = %q, se esperaba vacío", h.SprintID)
	}
	if h.ProyectoID != "P1" {
		t.Errorf("ProyectoID = %q, se esperaba P1", h.ProyectoID)
	}
	if h.ID != "" {
		t.Errorf("ID = %q, se esperaba vacío: lo asigna el repositorio (RN-05.1)", h.ID)
	}
}

func TestPrioridad_String(t *testing.T) {
	casos := map[Prioridad]string{
		PrioridadAlta:  "Alta",
		PrioridadMedia: "Media",
		PrioridadBaja:  "Baja",
		Prioridad(0):   "",
		Prioridad(99):  "",
	}
	for p, want := range casos {
		if got := p.String(); got != want {
			t.Errorf("Prioridad(%d).String() = %q, se esperaba %q", int(p), got, want)
		}
	}
}

func TestEstado_String(t *testing.T) {
	casos := map[Estado]string{
		Pendiente:  "Pendiente",
		EnProgreso: "En progreso",
		Completada: "Completada",
		Estado(0):  "",
		Estado(99): "",
	}
	for e, want := range casos {
		if got := e.String(); got != want {
			t.Errorf("Estado(%d).String() = %q, se esperaba %q", int(e), got, want)
		}
	}
}

func TestAsignarASprint(t *testing.T) {
	t.Run("guarda el sprint asignado", func(t *testing.T) {
		h, err := Nueva("P1", "Login", "", PrioridadAlta)
		if err != nil {
			t.Fatal(err)
		}
		if err := h.AsignarASprint("S1"); err != nil {
			t.Fatalf("error inesperado: %v", err)
		}
		if h.SprintID != "S1" {
			t.Errorf("SprintID = %q, se esperaba S1", h.SprintID)
		}
	})

	t.Run("una historia ya asignada no se asigna a otro sprint (RN-11.3)", func(t *testing.T) {
		h, err := Nueva("P1", "Login", "", PrioridadAlta)
		if err != nil {
			t.Fatal(err)
		}
		if err := h.AsignarASprint("S1"); err != nil {
			t.Fatal(err)
		}
		if err := h.AsignarASprint("S2"); !errors.Is(err, ErrYaAsignada) {
			t.Fatalf("error = %v, se esperaba ErrYaAsignada", err)
		}
		if h.SprintID != "S1" {
			t.Errorf("SprintID = %q, no debería cambiar; se esperaba S1", h.SprintID)
		}
	})
}
