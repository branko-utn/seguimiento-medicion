package metrics

import "errors"

// ErrDatosInconsistentes indica que las cantidades recibidas no pueden formar un porcentaje válido.
var ErrDatosInconsistentes = errors.New("las historias completadas no pueden superar al total")

// PorcentajeCompletadas calcula completadas / total × 100 (RN-25.1).
func PorcentajeCompletadas(completadas, total int) (float64, error) {
	return 0, nil
}
