package metrics

import "errors"

// ErrHorasNegativas indica que las horas no pueden ser negativas.
var ErrHorasNegativas = errors.New("las horas no pueden ser negativas")

// ErrDesviacionNoCalculable indica que no es posible calcular la desviación sin horas estimadas.
var ErrDesviacionNoCalculable = errors.New("la desviación no es calculable sin horas estimadas")

// HorasReales suma las horas registradas (RN-26.1).
// Se suman todas las horas recibidas sin filtrar horas negativas,
// asumiendo que los datos ya fueron validados por el dominio.
func HorasReales(horas []float64) float64 {
	total := 0.0
	for _, h := range horas {
		total += h
	}
	return total
}

// Desviacion calcula la desviación porcentual (RN-26.2).
// Sin horas estimadas (0) devuelve ErrDesviacionNoCalculable (RN-26.3).
// Si las horas estimadas o reales son negativas (< 0), devuelve 0 y ErrHorasNegativas.
func Desviacion(estimadas, reales float64) (float64, error) {
	if estimadas < 0 || reales < 0 {
		return 0, ErrHorasNegativas
	}
	if estimadas == 0 {
		return 0, ErrDesviacionNoCalculable
	}
	return (reales - estimadas) / estimadas * 100, nil
}
