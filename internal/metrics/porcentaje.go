package metrics

import "errors"

// ErrDatosInconsistentes indica que las cantidades recibidas no pueden formar un porcentaje válido.
var ErrDatosInconsistentes = errors.New("las historias completadas no pueden superar al total")

// ErrDesviacionNoCalculable indica que no es posible calcular la desviación sin horas estimadas.
var ErrDesviacionNoCalculable = errors.New("la desviación no es calculable sin horas estimadas")

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

// HorasReales suma las horas registradas (RN-26.1).
func HorasReales(horas []float64) float64 {
	total := 0.0
	for _, h := range horas {
		if h < 0 {
			continue
		}
		total += h
	}
	return total
}

// Desviacion calcula la desviación porcentual (RN-26.2).
// Sin horas estimadas (0) devuelve ErrDesviacionNoCalculable (RN-26.3).
// Horas negativas devuelven ErrDatosInconsistentes (RN-26.4).
func Desviacion(estimadas, reales float64) (float64, error) {
	if estimadas < 0 || reales < 0 {
		return 0, ErrDatosInconsistentes
	}
	if estimadas == 0 {
		return 0, ErrDesviacionNoCalculable
	}
	return (reales - estimadas) / estimadas * 100, nil
}
