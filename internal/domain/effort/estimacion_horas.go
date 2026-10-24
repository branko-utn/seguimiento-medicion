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

// NuevaEstimacion crea una estimación de horas válida para una historia.
// Las horas deben ser mayores que cero (RN-19.1).
func NuevaEstimacion(historiaID domain.HistoriaID, horas float64) (EstimacionHoras, error) {
	if horas <= 0 {
		return EstimacionHoras{}, ErrHorasInvalidas
	}
	return EstimacionHoras{
		HistoriaID: historiaID,
		Horas:      horas,
	}, nil
}
