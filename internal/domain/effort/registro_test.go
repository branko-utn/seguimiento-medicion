package effort

import (
	"errors"
	"testing"
	"time"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
)

func dia(s string) time.Time {
	f, err := time.ParseInLocation("2006-01-02", s, time.Local)
	if err != nil {
		panic(err)
	}
	return f
}

func TestNuevoRegistro(t *testing.T) {
	hoy := dia("2026-10-24")
	casos := []struct {
		titulo          string
		fecha           time.Time
		actividad       string
		horas           float64
		quiereError     error
		quiereActividad string
		quiereFecha     time.Time
	}{
		{titulo: "registro válido", fecha: hoy, actividad: "Desarrollo", horas: 3, quiereActividad: "Desarrollo", quiereFecha: hoy},
		{titulo: "fecha de hoy es válida", fecha: hoy, actividad: "Tests", horas: 1, quiereActividad: "Tests", quiereFecha: hoy},
		{titulo: "fecha pasada es válida", fecha: dia("2026-10-23"), actividad: "Tests", horas: 1, quiereActividad: "Tests", quiereFecha: dia("2026-10-23")},
		{titulo: "la hora de la fecha no cuenta: hoy a las 23:00 es válido", fecha: hoy.Add(23 * time.Hour), actividad: "Tests", horas: 1, quiereActividad: "Tests", quiereFecha: hoy},
		{titulo: "acepta fracciones de hora", fecha: hoy, actividad: "Revisión de PR", horas: 0.5, quiereActividad: "Revisión de PR", quiereFecha: hoy},
		{titulo: "acepta 24 horas en un registro", fecha: hoy, actividad: "Desarrollo", horas: 24, quiereActividad: "Desarrollo", quiereFecha: hoy},
		{titulo: "recorta los espacios de la actividad", fecha: hoy, actividad: "  Tests ", horas: 1, quiereActividad: "Tests", quiereFecha: hoy},
		{titulo: "actividad vacía", fecha: hoy, actividad: "", horas: 3, quiereError: ErrActividadVacia},
		{titulo: "actividad con solo espacios", fecha: hoy, actividad: "  ", horas: 3, quiereError: ErrActividadVacia},
		{titulo: "horas en cero", fecha: hoy, actividad: "Desarrollo", horas: 0, quiereError: ErrHorasInvalidas},
		{titulo: "horas negativas", fecha: hoy, actividad: "Desarrollo", horas: -1, quiereError: ErrHorasInvalidas},
		{titulo: "fecha de mañana", fecha: dia("2026-10-25"), actividad: "Desarrollo", horas: 2, quiereError: ErrFechaFutura},
		{titulo: "fecha vacía (RN-20.1: es obligatoria)", fecha: time.Time{}, actividad: "Desarrollo", horas: 2, quiereError: ErrFechaVacia},
	}

	for _, c := range casos {
		t.Run(c.titulo, func(t *testing.T) {
			r, err := NuevoRegistro("I1", "HU-1", c.fecha, c.actividad, c.horas, hoy)
			if c.quiereError != nil {
				if !errors.Is(err, c.quiereError) {
					t.Fatalf("se esperaba el error %v, se obtuvo %v", c.quiereError, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("no se esperaba error: %v", err)
			}
			if r.IntegranteID != "I1" || r.HistoriaID != "HU-1" {
				t.Errorf("IDs = (%s, %s), se esperaba (I1, HU-1)", r.IntegranteID, r.HistoriaID)
			}
			if r.Actividad != c.quiereActividad {
				t.Errorf("Actividad = %q, se esperaba %q", r.Actividad, c.quiereActividad)
			}
			if r.Horas != c.horas {
				t.Errorf("Horas = %v, se esperaba %v", r.Horas, c.horas)
			}
			if !r.Fecha.Equal(c.quiereFecha) {
				t.Errorf("Fecha = %v, se esperaba %v", r.Fecha, c.quiereFecha)
			}
			if r.ID != "" {
				t.Errorf("ID = %q, el ID lo genera el repositorio", r.ID)
			}
		})
	}
}

func TestValidarHorasDiarias(t *testing.T) {
	reg := func(integrante domain.IntegranteID, fecha string, horas float64) Registro {
		return Registro{IntegranteID: integrante, HistoriaID: "HU-1", Fecha: dia(fecha), Actividad: "Desarrollo", Horas: horas}
	}
	casos := []struct {
		titulo      string
		existentes  []Registro
		nuevo       Registro
		quiereError error
	}{
		{titulo: "sin registros previos", nuevo: reg("I1", "2026-10-23", 8)},
		{titulo: "10 + 14 horas el mismo día: 24 justas", existentes: []Registro{reg("I1", "2026-10-23", 10)}, nuevo: reg("I1", "2026-10-23", 14)},
		{titulo: "24 horas y 0,5 más el mismo día", existentes: []Registro{reg("I1", "2026-10-23", 24)}, nuevo: reg("I1", "2026-10-23", 0.5), quiereError: ErrSuperaHorasDiarias},
		{titulo: "10 + 14 + 0,5 el mismo día", existentes: []Registro{reg("I1", "2026-10-23", 10), reg("I1", "2026-10-23", 14)}, nuevo: reg("I1", "2026-10-23", 0.5), quiereError: ErrSuperaHorasDiarias},
		{titulo: "24 horas en un día y 1 al día siguiente", existentes: []Registro{reg("I1", "2026-10-23", 24)}, nuevo: reg("I1", "2026-10-24", 1)},
		{titulo: "las horas de otro integrante no cuentan", existentes: []Registro{reg("I2", "2026-10-23", 24)}, nuevo: reg("I1", "2026-10-23", 8)},
	}

	for _, c := range casos {
		t.Run(c.titulo, func(t *testing.T) {
			err := ValidarHorasDiarias(c.existentes, c.nuevo)
			if !errors.Is(err, c.quiereError) {
				t.Fatalf("se esperaba el error %v, se obtuvo %v", c.quiereError, err)
			}
		})
	}
}

func TestMensajesDeError(t *testing.T) {
	casos := map[error]string{
		ErrActividadVacia:           "la actividad es obligatoria",
		ErrFechaFutura:              "la fecha no puede ser futura",
		ErrFechaVacia:               "la fecha es obligatoria",
		ErrSuperaHorasDiarias:       "el integrante superaría las 24 horas en ese día",
		ErrIntegranteDeOtroProyecto: "el integrante no pertenece al proyecto de la historia",
		ErrHorasInvalidas:           "las horas deben ser mayores que cero",
	}
	for err, mensaje := range casos {
		if err.Error() != mensaje {
			t.Errorf("mensaje = %q, la spec dice %q", err.Error(), mensaje)
		}
	}
}
