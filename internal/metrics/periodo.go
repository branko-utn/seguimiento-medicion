package metrics

import "time"

// RegistroHoras es un registro de esfuerzo visto desde las métricas: solo la fecha y las horas (ADR 0002, sección 3).
// El servicio ya eligió los registros de las historias que corresponden (por ejemplo, las de un sprint).
type RegistroHoras struct {
	Fecha time.Time
	Horas float64
}

// HorasEnPeriodo suma las horas de los registros con fecha entre inicio y fin, ambos incluidos y comparando por día
// (RN-28.3 y RN-29.3, TEC-20). Sin registros en el período devuelve 0.
func HorasEnPeriodo(registros []RegistroHoras, inicio, fin time.Time) float64 {
	desde, hasta := soloDia(inicio), soloDia(fin)
	var horas []float64
	for _, r := range registros {
		if d := soloDia(r.Fecha); !d.Before(desde) && !d.After(hasta) {
			horas = append(horas, r.Horas)
		}
	}
	return HorasReales(horas)
}

// soloDia devuelve la fecha sin la hora.
func soloDia(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
