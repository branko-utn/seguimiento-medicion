package defect

import (
	"errors"
	"testing"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
)

func defectoEn(estado Estado) Defecto {
	d := Defecto{ID: "D1", Descripcion: "Acepta fechas invertidas", Severidad: Alta, Estado: estado, HistoriaID: "HU-1", SprintDeteccion: "S2"}
	if estado == Resuelto {
		d.SprintResolucion = "S2"
	}
	return d
}

func TestCambiarEstado(t *testing.T) {
	casos := []struct {
		titulo        string
		desde, hacia  Estado
		quiereError   error
		quiereMensaje string
	}{
		{titulo: "Abierto → En progreso", desde: Abierto, hacia: EnProgreso},
		{titulo: "En progreso → Abierto", desde: EnProgreso, hacia: Abierto},
		{titulo: "Abierto → Abierto", desde: Abierto, hacia: Abierto, quiereError: ErrTransicionInvalida, quiereMensaje: "no se puede pasar de Abierto a Abierto"},
		{titulo: "En progreso → En progreso", desde: EnProgreso, hacia: EnProgreso, quiereError: ErrTransicionInvalida, quiereMensaje: "no se puede pasar de En progreso a En progreso"},
		{titulo: "Resuelto → En progreso", desde: Resuelto, hacia: EnProgreso, quiereError: ErrTransicionInvalida, quiereMensaje: "no se puede pasar de Resuelto a En progreso"},
		{titulo: "Resuelto → Abierto", desde: Resuelto, hacia: Abierto, quiereError: ErrTransicionInvalida, quiereMensaje: "no se puede pasar de Resuelto a Abierto"},
		{titulo: "Abierto → Resuelto sin sprint", desde: Abierto, hacia: Resuelto, quiereError: ErrResolverSinSprint, quiereMensaje: "para resolver un defecto hay que indicar el sprint de resolución"},
		{titulo: "En progreso → Resuelto sin sprint", desde: EnProgreso, hacia: Resuelto, quiereError: ErrResolverSinSprint, quiereMensaje: "para resolver un defecto hay que indicar el sprint de resolución"},
		{titulo: "Resuelto → Resuelto: primero la transición, como HU-08", desde: Resuelto, hacia: Resuelto, quiereError: ErrTransicionInvalida, quiereMensaje: "no se puede pasar de Resuelto a Resuelto"},
	}

	for _, c := range casos {
		t.Run(c.titulo, func(t *testing.T) {
			d := defectoEn(c.desde)
			antes := d
			err := d.CambiarEstado(c.hacia)
			if c.quiereError != nil {
				if !errors.Is(err, c.quiereError) {
					t.Fatalf("se esperaba el error %v, se obtuvo %v", c.quiereError, err)
				}
				if err.Error() != c.quiereMensaje {
					t.Errorf("mensaje = %q, la spec dice %q", err.Error(), c.quiereMensaje)
				}
				if d != antes {
					t.Errorf("el defecto cambió aunque hubo error: %+v", d)
				}
				return
			}
			if err != nil {
				t.Fatalf("no se esperaba error: %v", err)
			}
			if d.Estado != c.hacia {
				t.Errorf("Estado = %v, se esperaba %v", d.Estado, c.hacia)
			}
			if d.SprintResolucion != "" {
				t.Errorf("SprintResolucion = %q, se esperaba vacío", d.SprintResolucion)
			}
		})
	}
}

func TestResolver(t *testing.T) {
	casos := []struct {
		titulo        string
		desde         Estado
		sprint        domain.SprintID
		quiereError   error
		quiereMensaje string
	}{
		{titulo: "desde En progreso", desde: EnProgreso, sprint: "S3"},
		{titulo: "directamente desde Abierto", desde: Abierto, sprint: "S2"},
		{titulo: "sin sprint", desde: Abierto, sprint: "", quiereError: ErrResolverSinSprint, quiereMensaje: "para resolver un defecto hay que indicar el sprint de resolución"},
		{titulo: "un defecto resuelto no cambia", desde: Resuelto, sprint: "S3", quiereError: ErrTransicionInvalida, quiereMensaje: "no se puede pasar de Resuelto a Resuelto"},
	}

	for _, c := range casos {
		t.Run(c.titulo, func(t *testing.T) {
			d := defectoEn(c.desde)
			antes := d
			err := d.Resolver(c.sprint)
			if c.quiereError != nil {
				if !errors.Is(err, c.quiereError) {
					t.Fatalf("se esperaba el error %v, se obtuvo %v", c.quiereError, err)
				}
				if err.Error() != c.quiereMensaje {
					t.Errorf("mensaje = %q, la spec dice %q", err.Error(), c.quiereMensaje)
				}
				if d != antes {
					t.Errorf("el defecto cambió aunque hubo error: %+v", d)
				}
				return
			}
			if err != nil {
				t.Fatalf("no se esperaba error: %v", err)
			}
			if d.Estado != Resuelto || d.SprintResolucion != c.sprint {
				t.Errorf("defecto = %+v, se esperaba Resuelto en %s", d, c.sprint)
			}
		})
	}
}

func TestMensajesDeErrorHU23(t *testing.T) {
	casos := map[error]string{
		ErrNoExiste:                 "el defecto no existe",
		ErrResolverSinSprint:        "para resolver un defecto hay que indicar el sprint de resolución",
		ErrResolucionDeOtroProyecto: "el sprint de resolución no pertenece al proyecto del defecto",
		ErrResolucionAnterior:       "el sprint de resolución no puede ser anterior al de detección",
	}
	for err, mensaje := range casos {
		if err.Error() != mensaje {
			t.Errorf("mensaje = %q, la spec dice %q", err.Error(), mensaje)
		}
	}
}
