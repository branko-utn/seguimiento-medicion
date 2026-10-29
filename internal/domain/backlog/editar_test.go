package backlog

import (
	"errors"
	"reflect"
	"testing"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
)

// historiaParaEditar arma una historia con todos los campos cargados, para verificar que Editar solo toca título,
// descripción y prioridad.
func historiaParaEditar(estado Estado, sprint string) Historia {
	return Historia{
		ID:          "HU-1",
		ProyectoID:  "P1",
		Titulo:      "Crear proyecto",
		Descripcion: "Alta de proyecto",
		Prioridad:   PrioridadMedia,
		Estado:      estado,
		StoryPoints: 3,
		Criterios:   []Criterio{{Texto: "Valida el nombre"}},
		SprintID:    domain.SprintID(sprint),
	}
}

func TestEditar(t *testing.T) {
	casos := []struct {
		nombre      string
		estado      Estado
		sprint      string
		titulo      string
		descripcion string
		prioridad   Prioridad
		wantTitulo  string
		wantErr     error
	}{
		{"cambia título, descripción y prioridad", Pendiente, "", "Crear y listar proyectos", "Alta y listado", PrioridadAlta, "Crear y listar proyectos", nil},
		{"recorta los espacios del título", Pendiente, "", "  Listar proyectos  ", "", PrioridadBaja, "Listar proyectos", nil},
		{"sin cambiar nada se guarda igual", Pendiente, "", "Crear proyecto", "Alta de proyecto", PrioridadMedia, "Crear proyecto", nil},
		{"asignada a un sprint y en progreso se puede editar", EnProgreso, "S1", "Otro título", "", PrioridadAlta, "Otro título", nil},
		{"título vacío (RN-06.1)", Pendiente, "", "", "", PrioridadAlta, "", ErrTituloVacio},
		{"título con solo espacios (RN-06.1)", Pendiente, "", "   ", "", PrioridadAlta, "", ErrTituloVacio},
		{"prioridad sin valor (RN-06.1)", Pendiente, "", "Login", "", Prioridad(0), "", ErrPrioridadInvalida},
		{"prioridad fuera de la lista (RN-06.1)", Pendiente, "", "Login", "", Prioridad(99), "", ErrPrioridadInvalida},
		{"completada no se edita (RN-06.2)", Completada, "S1", "Otro título", "", PrioridadAlta, "", ErrCompletadaNoSeEdita},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			h := historiaParaEditar(c.estado, c.sprint)
			antes := historiaParaEditar(c.estado, c.sprint)

			err := h.Editar(c.titulo, c.descripcion, c.prioridad)
			if !errors.Is(err, c.wantErr) {
				t.Fatalf("error = %v, se esperaba %v", err, c.wantErr)
			}
			if c.wantErr != nil {
				// Si falla, la historia no cambia.
				if !reflect.DeepEqual(h, antes) {
					t.Errorf("con error la historia no debería cambiar:\nquedó   %+v\nantes   %+v", h, antes)
				}
				return
			}

			want := antes
			want.Titulo = c.wantTitulo
			want.Descripcion = c.descripcion
			want.Prioridad = c.prioridad
			// El ID, el estado, los Story Points, los criterios y el sprint no cambian.
			if !reflect.DeepEqual(h, want) {
				t.Errorf("historia editada:\nquedó   %+v\nesperada %+v", h, want)
			}
		})
	}
}

func TestPuedeEliminarse(t *testing.T) {
	casos := []struct {
		nombre  string
		estado  Estado
		sprint  string
		wantErr error
	}{
		{"pendiente sin sprint se puede eliminar", Pendiente, "", nil},
		{"en progreso sin sprint se puede eliminar", EnProgreso, "", nil},
		{"asignada a un sprint no se elimina (RN-06.3)", Pendiente, "S1", ErrAsignadaNoSeElimina},
		{"en progreso en un sprint no se elimina (RN-06.3)", EnProgreso, "S1", ErrAsignadaNoSeElimina},
		{"completada no se elimina (RN-06.3)", Completada, "", ErrCompletadaNoSeElimina},
		{"completada en un sprint informa que está completada (RN-06.3)", Completada, "S1", ErrCompletadaNoSeElimina},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			h := historiaParaEditar(c.estado, c.sprint)
			if err := h.PuedeEliminarse(); !errors.Is(err, c.wantErr) {
				t.Fatalf("error = %v, se esperaba %v", err, c.wantErr)
			}
		})
	}
}
