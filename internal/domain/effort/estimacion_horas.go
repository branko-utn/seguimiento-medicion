package effort

import (
	"errors"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
)

var (
	ErrHorasInvalidas = errors.New("las horas deben ser mayores que cero")
)

// EstimacionHoras representa las horas estimadas para una historia (ADR 0002).
type EstimacionHoras struct {
	HistoriaID domain.HistoriaID
	Horas      float64
}

// NuevaEstimacion crea una estimacion de horas valida para una historia.
func NuevaEstimacion(historiaID domain.HistoriaID, horas float64) (EstimacionHoras, error) {
	return EstimacionHoras{}, ErrHorasInvalidas
}
