package effort

import (
	"errors"
	"strings"
	"time"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
)

// MaxHorasPorDia es el máximo de horas que un integrante puede registrar en un mismo día (RN-20.3).
const MaxHorasPorDia = 24

// Registro es el esfuerzo que un integrante trabajó en una historia en un día (ADR 0002).
type Registro struct {
	ID           domain.RegistroID
	IntegranteID domain.IntegranteID
	HistoriaID   domain.HistoriaID
	Fecha        time.Time
	Actividad    string
	Horas        float64
}

var (
	ErrActividadVacia           = errors.New("la actividad es obligatoria")
	ErrFechaFutura              = errors.New("la fecha no puede ser futura")
	ErrSuperaHorasDiarias       = errors.New("el integrante superaría las 24 horas en ese día")
	ErrIntegranteDeOtroProyecto = errors.New("el integrante no pertenece al proyecto de la historia")
)

// NuevoRegistro crea un registro de esfuerzo válido: actividad no vacía (RN-20.1), horas mayores que cero (RN-20.2) y
// fecha no posterior a hoy (RN-20.4). La fecha de hoy se recibe como parámetro para no depender del reloj.
// Que el integrante y la historia existan, que sean del mismo proyecto y el tope diario los controla el servicio.
func NuevoRegistro(integranteID domain.IntegranteID, historiaID domain.HistoriaID, fecha time.Time, actividad string, horas float64, hoy time.Time) (Registro, error) {
	actividad = strings.TrimSpace(actividad)
	if actividad == "" {
		return Registro{}, ErrActividadVacia
	}
	if horas <= 0 {
		return Registro{}, ErrHorasInvalidas
	}
	fecha = soloDia(fecha)
	if fecha.After(soloDia(hoy)) {
		return Registro{}, ErrFechaFutura
	}
	return Registro{
		IntegranteID: integranteID,
		HistoriaID:   historiaID,
		Fecha:        fecha,
		Actividad:    actividad,
		Horas:        horas,
	}, nil
}

// ValidarHorasDiarias controla que el integrante del registro nuevo no supere MaxHorasPorDia en ese día, sumando sus
// registros existentes de ese día y el nuevo (RN-20.3).
func ValidarHorasDiarias(existentes []Registro, nuevo Registro) error {
	total := nuevo.Horas
	for _, r := range existentes {
		if r.IntegranteID == nuevo.IntegranteID && mismoDia(r.Fecha, nuevo.Fecha) {
			total += r.Horas
		}
	}
	if total > MaxHorasPorDia {
		return ErrSuperaHorasDiarias
	}
	return nil
}

// soloDia quita la hora de una fecha.
func soloDia(t time.Time) time.Time {
	a, m, d := t.Date()
	return time.Date(a, m, d, 0, 0, 0, 0, t.Location())
}

func mismoDia(a, b time.Time) bool {
	return soloDia(a).Equal(soloDia(b))
}
