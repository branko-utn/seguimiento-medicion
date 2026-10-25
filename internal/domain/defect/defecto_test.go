package defect

import (
	"errors"
	"testing"
)

func TestNuevo(t *testing.T) {
	casos := []struct {
		titulo            string
		descripcion       string
		severidad         Severidad
		quiereError       error
		quiereDescripcion string
	}{
		{titulo: "defecto válido", descripcion: "Acepta fechas invertidas", severidad: Alta, quiereDescripcion: "Acepta fechas invertidas"},
		{titulo: "severidad Crítica", descripcion: "Se cae al guardar", severidad: Critica, quiereDescripcion: "Se cae al guardar"},
		{titulo: "severidad Media", descripcion: "Texto cortado", severidad: Media, quiereDescripcion: "Texto cortado"},
		{titulo: "severidad Baja", descripcion: "Falta una tilde", severidad: Baja, quiereDescripcion: "Falta una tilde"},
		{titulo: "recorta los espacios de la descripción", descripcion: "  Falta una tilde ", severidad: Baja, quiereDescripcion: "Falta una tilde"},
		{titulo: "descripción vacía", descripcion: "", severidad: Alta, quiereError: ErrDescripcionVacia},
		{titulo: "descripción con solo espacios", descripcion: "  ", severidad: Alta, quiereError: ErrDescripcionVacia},
		{titulo: "severidad cero (sin valor)", descripcion: "Acepta fechas invertidas", severidad: 0, quiereError: ErrSeveridadInvalida},
		{titulo: "severidad fuera de la lista", descripcion: "Acepta fechas invertidas", severidad: Baja + 1, quiereError: ErrSeveridadInvalida},
	}

	for _, c := range casos {
		t.Run(c.titulo, func(t *testing.T) {
			d, err := Nuevo(c.descripcion, c.severidad, "HU-1", "S1")
			if c.quiereError != nil {
				if !errors.Is(err, c.quiereError) {
					t.Fatalf("se esperaba el error %v, se obtuvo %v", c.quiereError, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("no se esperaba error: %v", err)
			}
			if d.Descripcion != c.quiereDescripcion {
				t.Errorf("Descripcion = %q, se esperaba %q", d.Descripcion, c.quiereDescripcion)
			}
			if d.Severidad != c.severidad {
				t.Errorf("Severidad = %v, se esperaba %v", d.Severidad, c.severidad)
			}
			if d.HistoriaID != "HU-1" || d.SprintDeteccion != "S1" {
				t.Errorf("historia y sprint = (%s, %s), se esperaba (HU-1, S1)", d.HistoriaID, d.SprintDeteccion)
			}
			// RN-22.3: un defecto nuevo queda Abierto y sin sprint de resolución.
			if d.Estado != Abierto {
				t.Errorf("Estado = %v, se esperaba Abierto", d.Estado)
			}
			if d.SprintResolucion != "" {
				t.Errorf("SprintResolucion = %q, se esperaba vacío", d.SprintResolucion)
			}
			if d.ID != "" {
				t.Errorf("ID = %q, el ID lo genera el repositorio", d.ID)
			}
		})
	}
}

func TestSeveridadString(t *testing.T) {
	casos := map[Severidad]string{Critica: "Crítica", Alta: "Alta", Media: "Media", Baja: "Baja", 0: "", Baja + 1: ""}
	for s, texto := range casos {
		if s.String() != texto {
			t.Errorf("Severidad(%d).String() = %q, se esperaba %q", int(s), s.String(), texto)
		}
	}
}

func TestEstadoString(t *testing.T) {
	casos := map[Estado]string{Abierto: "Abierto", EnProgreso: "En progreso", Resuelto: "Resuelto", 0: "", Resuelto + 1: ""}
	for e, texto := range casos {
		if e.String() != texto {
			t.Errorf("Estado(%d).String() = %q, se esperaba %q", int(e), e.String(), texto)
		}
	}
}

func TestMensajesDeError(t *testing.T) {
	casos := map[error]string{
		ErrDescripcionVacia:     "la descripción es obligatoria",
		ErrSeveridadInvalida:    "la severidad debe ser Crítica, Alta, Media o Baja",
		ErrSprintDeOtroProyecto: "el sprint de detección no pertenece al proyecto de la historia",
	}
	for err, mensaje := range casos {
		if err.Error() != mensaje {
			t.Errorf("mensaje = %q, la spec dice %q", err.Error(), mensaje)
		}
	}
}
