package effort

import (
	"errors"
	"time"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
)

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

func NuevoRegistro(integranteID domain.IntegranteID, historiaID domain.HistoriaID, fecha time.Time, actividad string, horas float64, hoy time.Time) (Registro, error) {
	return Registro{}, nil
}

func ValidarHorasDiarias(existentes []Registro, nuevo Registro) error {
	return nil
}
