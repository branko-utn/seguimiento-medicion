package metrics

import "errors"

// ErrDatosInconsistentes indica que las cantidades recibidas no pueden formar un porcentaje válido.
var ErrDatosInconsistentes = errors.New("las historias completadas no pueden superar al total")

// PorcentajeCompletadas calcula completadas / total × 100 (RN-25.1).
// Si el total es 0 devuelve 0 (RN-25.2); las cantidades negativas o las completadas mayores que el total
// devuelven ErrDatosInconsistentes (RN-25.3).
func PorcentajeCompletadas(completadas, total int) (float64, error) {
	if completadas < 0 || total < 0 || completadas > total {
		return 0, ErrDatosInconsistentes
	}
	if total == 0 {
		return 0, nil
	}
	return float64(completadas) / float64(total) * 100, nil
}
